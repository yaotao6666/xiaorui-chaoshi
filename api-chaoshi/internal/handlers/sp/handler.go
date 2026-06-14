package sp

import (
	"chaoshi_api/internal/models"
	"chaoshi_api/internal/storage"
	"chaoshi_api/internal/utils"
	"chaoshi_api/pkg/database"
	"chaoshi_api/pkg/response"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func parseModelJSON(value models.JSON) any {
	if len(value) == 0 {
		return nil
	}
	var out any
	if err := json.Unmarshal([]byte(value), &out); err != nil {
		return nil
	}
	return out
}

func parseStringSlice(value any) []string {
	rawList, ok := value.([]any)
	if !ok {
		return []string{}
	}
	result := make([]string, 0, len(rawList))
	for _, item := range rawList {
		if s, ok := item.(string); ok && s != "" {
			result = append(result, s)
		}
	}
	return result
}

func buildAccessibleMerchantAsset(resource string) string {
	service := storage.GetService()
	if service == nil {
		return resource
	}

	return service.BuildURL(resource)
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError, "参数错误")
		return
	}

	var adminUser models.AdminUser
	if err := database.DB.Where("username = ? AND status = ?", req.Username, 1).First(&adminUser).Error; err != nil {
		response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "用户名或密码错误")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(adminUser.Password), []byte(req.Password)); err != nil {
		response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "用户名或密码错误")
		return
	}

	now := time.Now()
	database.DB.Model(&models.AdminUser{}).Where("id = ?", adminUser.ID).Update("last_login_at", now)

	token, _ := utils.GenerateToken(adminUser.ID, "admin", adminUser.Username)
	database.DB.Preload("AdminProfile").First(&adminUser, adminUser.ID)
	adminProfileName := ""
	if adminUser.AdminProfile != nil {
		adminProfileName = adminUser.AdminProfile.Name
	}
	if adminProfileName == "" {
		adminProfileName = adminUser.Username
	}
	displayName := adminUser.Name
	if displayName == "" {
		displayName = adminUser.Username
	}
	response.Success(c, gin.H{
		"token": token,
		"admin": gin.H{
			"id":           adminUser.AdminProfileID,
			"name":         adminProfileName,
			"display_name": displayName,
		},
	})
}

func Logout(c *gin.Context) {
	response.Success(c, gin.H{"message": "退出成功"})
}

func GetDashboard(c *gin.Context) {
	if _, ok := getCurrentAdminUserID(c); !ok {
		response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "后台身份无效")
		return
	}

	var totalMerchants int64
	var todayOrders int64
	var todayRevenue float64
	merchantIDsQuery := database.DB.Model(&models.Merchant{}).Select("id")

	database.DB.Model(&models.Merchant{}).Count(&totalMerchants)
	database.DB.Model(&models.Order{}).
		Where("merchant_id IN (?) AND DATE(created_at) = CURDATE()", merchantIDsQuery).
		Count(&todayOrders)
	database.DB.Model(&models.Order{}).
		Where("merchant_id IN (?) AND DATE(created_at) = CURDATE() AND status >= 2", merchantIDsQuery).
		Select("COALESCE(SUM(pay_amount), 0)").
		Scan(&todayRevenue)

	var distribution []struct {
		Category string `json:"category"`
		Count    int64  `json:"count"`
	}
	database.DB.Model(&models.Merchant{}).
		Select("business_category as category, count(*) as count").
		Group("business_category").
		Scan(&distribution)

	var trend []struct {
		Date   string `json:"date"`
		Orders int64  `json:"orders"`
	}
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		var orders int64
		database.DB.Model(&models.Order{}).
			Where("merchant_id IN (?) AND DATE(created_at) = ?", merchantIDsQuery, date).
			Count(&orders)
		trend = append(trend, struct {
			Date   string `json:"date"`
			Orders int64  `json:"orders"`
		}{Date: date, Orders: orders})
	}

	response.Success(c, gin.H{
		"total_merchants":   totalMerchants,
		"pending_merchants": 0,
		"today_orders":      todayOrders,
		"today_revenue":     todayRevenue,
		"distribution":      distribution,
		"trend":             trend,
	})
}

func GetMerchantDetail(c *gin.Context) {
	merchantID := c.Param("merchant_id")
	id, _ := strconv.ParseUint(merchantID, 10, 64)
	if _, ok := getCurrentAdminUserID(c); !ok {
		response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "后台身份无效")
		return
	}

	var merchant models.Merchant
	if err := database.DB.Where("id = ?", id).First(&merchant).Error; err != nil {
		response.Fail(c, http.StatusNotFound, response.CodeNotFound, "商家不存在")
		return
	}

	var totalOrders int64
	database.DB.Model(&models.Order{}).Where("merchant_id = ?", id).Count(&totalOrders)

	var totalUsers int64
	database.DB.Model(&models.UserVisit{}).Where("merchant_id = ?", id).Distinct("user_id").Count(&totalUsers)

	var totalAmount float64
	database.DB.Model(&models.Order{}).
		Where("merchant_id = ? AND status >= 2", id).
		Select("COALESCE(SUM(pay_amount), 0)").
		Scan(&totalAmount)

	var adminStaff models.MerchantStaff
	_ = database.DB.Select("id", "username", "name", "phone", "status").
		Where("merchant_id = ? AND role = ?", id, "owner").
		Order("status DESC, id ASC").
		First(&adminStaff).Error

	response.Success(c, gin.H{
		"id":                merchant.ID,
		"name":              merchant.Name,
		"logo":              buildAccessibleMerchantAsset(merchant.Logo),
		"cover_image":       buildAccessibleMerchantAsset(merchant.CoverImage),
		"contact_name":      merchant.ContactName,
		"contact_phone":     merchant.ContactPhone,
		"contact_email":     merchant.ContactEmail,
		"address":           merchant.Address,
		"business_category": merchant.BusinessCategory,
		"business_hours":    merchant.BusinessHours,
		"announcement":      merchant.Announcement,
		"status":            merchant.Status,
		"qrcode_url":        "",
		"created_at":        merchant.CreatedAt,
		"admin_staff_id":    adminStaff.ID,
		"admin_username":    adminStaff.Username,
		"admin_name":        adminStaff.Name,
		"admin_phone":       adminStaff.Phone,
		"admin_status":      adminStaff.Status,
		"total_orders":      totalOrders,
		"total_amount":      totalAmount,
		"total_users":       totalUsers,
	})
}

type UpdateMerchantAssetsRequest struct {
	Logo       string `json:"logo"`
	CoverImage string `json:"cover_image"`
}

func UpdateMerchantAssets(c *gin.Context) {
	merchantID := c.Param("merchant_id")
	id, _ := strconv.ParseUint(merchantID, 10, 64)
	if _, ok := getCurrentAdminUserID(c); !ok {
		response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "后台身份无效")
		return
	}

	var req UpdateMerchantAssetsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError, "参数错误")
		return
	}

	updates := map[string]interface{}{}
	if req.Logo != "" {
		updates["logo"] = req.Logo
	}
	if req.CoverImage != "" {
		updates["cover_image"] = req.CoverImage
	}
	if len(updates) == 0 {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError, "未提供可更新内容")
		return
	}

	if err := database.DB.Model(&models.Merchant{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, "更新商家图片失败")
		return
	}

	GetMerchantDetail(c)
}

func GetMerchantDistribution(c *gin.Context) {
	if _, ok := getCurrentAdminUserID(c); !ok {
		response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "后台身份无效")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "6"))
	sortBy := c.DefaultQuery("sort_by", "order_amount")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 6
	}

	metrics, totals := buildSpMerchantMetrics()

	validSortFields := map[string]string{
		"visit_rate":       "visit_rate",
		"order_rate":       "order_rate",
		"order_amount":     "order_amount",
		"avg_order_amount": "avg_order_amount",
		"visit_users":      "visit_users",
		"order_users":      "order_users",
		"paid_orders":      "paid_orders",
	}
	if _, valid := validSortFields[sortBy]; !valid {
		sortBy = "order_amount"
	}
	if sortOrder != "asc" {
		sortOrder = "desc"
	}

	sortField := validSortFields[sortBy]
	sortMultiplier := 1.0
	if sortOrder == "desc" {
		sortMultiplier = -1.0
	}
	sortSlice(metrics, func(item gin.H) float64 {
		if val, ok := item[sortField]; ok {
			switch v := val.(type) {
			case float64:
				return v * sortMultiplier
			case int64:
				return float64(v) * sortMultiplier
			case int:
				return float64(v) * sortMultiplier
			}
		}
		return 0
	})

	total := len(metrics)
	start := (page - 1) * pageSize
	if start > total {
		start = total
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	pagedMetrics := metrics[start:end]

	response.Success(c, gin.H{
		"merchants": pagedMetrics,
		"totals":    totals,
		"pagination": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func GetMerchantList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")
	status := c.Query("status")
	if _, ok := getCurrentAdminUserID(c); !ok {
		response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "后台身份无效")
		return
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	query := database.DB.Model(&models.Merchant{})

	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if status != "" {
		statusInt, _ := strconv.Atoi(status)
		query = query.Where("status = ?", statusInt)
	}

	var total int64
	query.Count(&total)

	var merchants []models.Merchant
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&merchants).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, "获取商家列表失败")
		return
	}

	list := make([]gin.H, 0, len(merchants))
	for _, m := range merchants {
		var totalUsers int64
		database.DB.Model(&models.UserVisit{}).
			Where("merchant_id = ?", m.ID).
			Distinct("user_id").
			Count(&totalUsers)

		var totalOrders int64
		database.DB.Model(&models.Order{}).
			Where("merchant_id = ?", m.ID).
			Count(&totalOrders)

		var totalAmount float64
		database.DB.Model(&models.Order{}).
			Where("merchant_id = ? AND status >= 2", m.ID).
			Select("COALESCE(SUM(pay_amount), 0)").
			Scan(&totalAmount)

		list = append(list, gin.H{
			"id":                m.ID,
			"name":              m.Name,
			"logo":              buildAccessibleMerchantAsset(m.Logo),
			"cover_image":       buildAccessibleMerchantAsset(m.CoverImage),
			"contact_name":      m.ContactName,
			"contact_phone":     m.ContactPhone,
			"contact_email":     m.ContactEmail,
			"address":           m.Address,
			"business_category": m.BusinessCategory,
			"business_hours":    m.BusinessHours,
			"announcement":      m.Announcement,
			"status":            m.Status,
			"created_at":        m.CreatedAt,
			"total_users":       totalUsers,
			"total_orders":      totalOrders,
			"total_amount":      totalAmount,
		})
	}

	response.Success(c, gin.H{
		"list": list,
		"pagination": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func GetOrderAnalytics(c *gin.Context) {
	if _, ok := getCurrentAdminUserID(c); !ok {
		response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "后台身份无效")
		return
	}

	location := time.Now().Location()
	now := time.Now().In(location)

	response.Success(c, gin.H{
		"day": buildSpOrderBuckets(now, 7, func(base time.Time, offset int) (time.Time, time.Time, string) {
			start := time.Date(base.Year(), base.Month(), base.Day(), 0, 0, 0, 0, location).AddDate(0, 0, -(6 - offset))
			end := start.Add(24 * time.Hour)
			return start, end, start.Format("01-02")
		}),
		"week": buildSpOrderBuckets(now, 8, func(base time.Time, offset int) (time.Time, time.Time, string) {
			weekdayOffset := (int(base.Weekday()) + 6) % 7
			weekStart := time.Date(base.Year(), base.Month(), base.Day(), 0, 0, 0, 0, location).AddDate(0, 0, -weekdayOffset)
			start := weekStart.AddDate(0, 0, -7*(7-offset))
			end := start.AddDate(0, 0, 7)
			year, week := start.ISOWeek()
			return start, end, fmt.Sprintf("%d-W%02d", year, week)
		}),
		"month": buildSpOrderBuckets(now, 12, func(base time.Time, offset int) (time.Time, time.Time, string) {
			start := time.Date(base.Year(), base.Month(), 1, 0, 0, 0, 0, location).AddDate(0, -(11 - offset), 0)
			end := start.AddDate(0, 1, 0)
			return start, end, start.Format("2006-01")
		}),
		"year": buildSpOrderBuckets(now, 5, func(base time.Time, offset int) (time.Time, time.Time, string) {
			start := time.Date(base.Year()-(4-offset), 1, 1, 0, 0, 0, 0, location)
			end := start.AddDate(1, 0, 0)
			return start, end, start.Format("2006")
		}),
	})
}

func GetAmountAnalytics(c *gin.Context) {
	if _, ok := getCurrentAdminUserID(c); !ok {
		response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "后台身份无效")
		return
	}

	days := c.DefaultQuery("days", "7")
	daysInt, _ := strconv.Atoi(days)

	merchantIDsQuery := database.DB.Model(&models.Merchant{}).Select("id")

	var trends []struct {
		Date   string  `json:"date"`
		Amount float64 `json:"amount"`
	}

	for i := daysInt - 1; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		var amount float64
		database.DB.Model(&models.Order{}).
			Where("merchant_id IN (?) AND DATE(created_at) = ? AND status >= 2", merchantIDsQuery, date).
			Select("COALESCE(SUM(pay_amount), 0)").Scan(&amount)
		trends = append(trends, struct {
			Date   string  `json:"date"`
			Amount float64 `json:"amount"`
		}{Date: date, Amount: amount})
	}

	response.Success(c, gin.H{
		"trends": trends,
	})
}

func GetTopMerchants(c *gin.Context) {
	if _, ok := getCurrentAdminUserID(c); !ok {
		response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "后台身份无效")
		return
	}

	limit := c.DefaultQuery("limit", "10")
	limitInt, _ := strconv.Atoi(limit)
	if limitInt <= 0 {
		limitInt = 10
	}
	metric := c.DefaultQuery("metric", "order_amount")

	metrics, _ := buildSpMerchantMetrics()
	sort.Slice(metrics, func(i, j int) bool {
		left := getSpMerchantMetricValue(metrics[i], metric)
		right := getSpMerchantMetricValue(metrics[j], metric)
		if left == right {
			return metrics[i]["merchant_name"].(string) < metrics[j]["merchant_name"].(string)
		}
		return left > right
	})

	if len(metrics) > limitInt {
		metrics = metrics[:limitInt]
	}

	list := make([]gin.H, 0, len(metrics))
	for index, item := range metrics {
		list = append(list, gin.H{
			"rank":             index + 1,
			"metric":           metric,
			"merchant_id":      item["merchant_id"],
			"merchant_name":    item["merchant_name"],
			"merchant_logo":    item["merchant_logo"],
			"visit_rate":       item["visit_rate"],
			"order_rate":       item["order_rate"],
			"order_amount":     item["order_amount"],
			"avg_order_amount": item["avg_order_amount"],
			"visit_users":      item["visit_users"],
			"order_users":      item["order_users"],
			"paid_orders":      item["paid_orders"],
		})
	}

	response.Success(c, list)
}

func buildSpMerchantMetrics() ([]gin.H, gin.H) {
	var merchants []models.Merchant
	database.DB.Order("created_at DESC").Find(&merchants)

	totalVisitUsers := int64(0)
	totalOrderUsers := int64(0)
	totalPaidOrders := int64(0)
	totalOrderAmount := 0.0

	metrics := make([]gin.H, 0, len(merchants))
	for _, merchant := range merchants {
		var visitUsers int64
		var orderUsers int64
		var paidOrders int64
		var orderAmount float64

		database.DB.Model(&models.UserVisit{}).
			Where("merchant_id = ?", merchant.ID).
			Distinct("user_id").
			Count(&visitUsers)
		database.DB.Model(&models.Order{}).
			Where("merchant_id = ? AND status >= 2", merchant.ID).
			Distinct("user_id").
			Count(&orderUsers)
		database.DB.Model(&models.Order{}).
			Where("merchant_id = ? AND status >= 2", merchant.ID).
			Count(&paidOrders)
		database.DB.Model(&models.Order{}).
			Where("merchant_id = ? AND status >= 2", merchant.ID).
			Select("COALESCE(SUM(pay_amount), 0)").
			Scan(&orderAmount)

		totalVisitUsers += visitUsers
		totalOrderUsers += orderUsers
		totalPaidOrders += paidOrders
		totalOrderAmount += orderAmount

		avgOrderAmount := 0.0
		if paidOrders > 0 {
			avgOrderAmount = orderAmount / float64(paidOrders)
		}

		metrics = append(metrics, gin.H{
			"merchant_id":      merchant.ID,
			"merchant_name":    merchant.Name,
			"merchant_logo":    buildAccessibleMerchantAsset(merchant.Logo),
			"visit_users":      visitUsers,
			"order_users":      orderUsers,
			"paid_orders":      paidOrders,
			"order_amount":     orderAmount,
			"avg_order_amount": avgOrderAmount,
			"visit_rate":       0.0,
			"order_rate":       0.0,
		})
	}

	for _, item := range metrics {
		visitUsers := item["visit_users"].(int64)
		orderUsers := item["order_users"].(int64)
		if totalVisitUsers > 0 {
			item["visit_rate"] = float64(visitUsers) / float64(totalVisitUsers) * 100
		}
		if visitUsers > 0 {
			item["order_rate"] = float64(orderUsers) / float64(visitUsers) * 100
		}
	}

	return metrics, gin.H{
		"merchant_count": totalCountInt64(len(metrics)),
		"visit_users":    totalVisitUsers,
		"order_users":    totalOrderUsers,
		"paid_orders":    totalPaidOrders,
		"order_amount":   totalOrderAmount,
	}
}

func totalCountInt64(value int) int64 {
	return int64(value)
}

func sortSlice(items []gin.H, keyFunc func(gin.H) float64) {
	n := len(items)
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			if keyFunc(items[i]) > keyFunc(items[j]) {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
}

func getSpMerchantMetricValue(item gin.H, metric string) float64 {
	switch metric {
	case "visit_rate":
		return item["visit_rate"].(float64)
	case "order_rate":
		return item["order_rate"].(float64)
	case "avg_order_amount":
		return item["avg_order_amount"].(float64)
	default:
		return item["order_amount"].(float64)
	}
}

func buildSpOrderBuckets(
	base time.Time,
	count int,
	rangeBuilder func(base time.Time, offset int) (time.Time, time.Time, string),
) []gin.H {
	merchantIDsQuery := database.DB.Model(&models.Merchant{}).Select("id")
	result := make([]gin.H, 0, count)
	for index := 0; index < count; index++ {
		start, end, label := rangeBuilder(base, index)
		var orderCount int64
		database.DB.Model(&models.Order{}).
			Where("merchant_id IN (?) AND status >= 2 AND created_at >= ? AND created_at < ?", merchantIDsQuery, start, end).
			Count(&orderCount)
		result = append(result, gin.H{
			"label":       label,
			"order_count": orderCount,
		})
	}
	return result
}

func GetMerchantQRCode(c *gin.Context) {
	merchantID := c.Param("merchant_id")
	id, _ := strconv.ParseUint(merchantID, 10, 64)

	if _, ok := getCurrentAdminUserID(c); !ok {
		response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "后台身份无效")
		return
	}

	var merchant models.Merchant
	if err := database.DB.Where("id = ?", id).Select("id", "name").First(&merchant).Error; err != nil {
		response.Fail(c, http.StatusNotFound, response.CodeNotFound, "商家不存在")
		return
	}

	qrCode, err := utils.GenerateMerchantStoreQRCode(id, 280)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, "生成微信小程序码失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"merchant_id":   id,
		"merchant_name": merchant.Name,
		"qrcode_url":    qrCode.QRCodeURL,
		"page_path":     qrCode.Page,
		"scene":         qrCode.Scene,
		"placeholder":   false,
	})
}

func GetRefunds(c *gin.Context) {
	if _, ok := getCurrentAdminUserID(c); !ok {
		response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "后台身份无效")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	merchantID := c.Query("merchant_id")
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	merchantIDsQuery := database.DB.Model(&models.Merchant{}).Select("id")

	query := database.DB.Model(&models.Refund{}).Preload("Order").
		Joins("JOIN orders ON orders.id = refunds.order_id").
		Where("orders.merchant_id IN (?)", merchantIDsQuery)

	if merchantID != "" {
		id, _ := strconv.ParseUint(merchantID, 10, 64)
		query = query.Where("orders.merchant_id = ?", id)
	}
	if status != "" {
		query = query.Where("refunds.status = ?", status)
	}

	var total int64
	query.Count(&total)

	var refunds []models.Refund
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&refunds).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, "获取退款列表失败")
		return
	}

	response.Success(c, gin.H{
		"list": refunds,
		"pagination": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

type changePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func ChangePassword(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "未登录")
		return
	}
	spUserID, ok := userIDValue.(uint64)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "未登录")
		return
	}

	var req changePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError, "参数错误")
		return
	}

	var adminUser models.AdminUser
	if err := database.DB.First(&adminUser, spUserID).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, "修改失败")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(adminUser.Password), []byte(req.OldPassword)); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError, "旧密码错误")
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, "修改失败")
		return
	}

	if err := database.DB.Model(&models.AdminUser{}).Where("id = ?", adminUser.ID).Update("password", string(hashed)).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, "修改失败")
		return
	}

	response.Success(c, gin.H{"message": "修改成功"})
}

func GetActivities(c *gin.Context) {
	if _, ok := getCurrentAdminUserID(c); !ok {
		response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "后台身份无效")
		return
	}

	var banners []models.Activity
	var announcements []models.Activity

	database.DB.Where("type = ? AND status = ?", "banner", 1).Order("sort ASC, created_at DESC").Find(&banners)
	database.DB.Where("type = ? AND status = ?", "announcement", 1).Order("sort ASC, created_at DESC").Find(&announcements)

	response.Success(c, gin.H{
		"banners":       banners,
		"announcements": announcements,
	})
}

type CreateActivityRequest struct {
	Type      string `json:"type" binding:"required,oneof=banner announcement"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Image     string `json:"image"`
	LinkType  string `json:"link_type" binding:"omitempty,oneof=merchant webview none"`
	LinkValue string `json:"link_value"`
	Sort      uint   `json:"sort"`
	Status    uint8  `json:"status" binding:"omitempty,oneof=0 1"`
}

func CreateActivity(c *gin.Context) {
	if _, ok := getCurrentAdminUserID(c); !ok {
		response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "后台身份无效")
		return
	}

	var req CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError, "参数错误")
		return
	}

	activity := models.Activity{
		Type:      req.Type,
		Title:     req.Title,
		Content:   req.Content,
		Image:     req.Image,
		LinkType:  req.LinkType,
		LinkValue: req.LinkValue,
		Sort:      req.Sort,
		Status:    1,
	}
	if req.Status > 0 {
		activity.Status = req.Status
	}

	if err := database.DB.Create(&activity).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, "创建活动失败")
		return
	}

	response.Success(c, gin.H{"id": activity.ID, "message": "创建成功"})
}

func UpdateActivity(c *gin.Context) {
	id := c.Param("id")
	activityID, _ := strconv.ParseUint(id, 10, 64)

	if _, ok := getCurrentAdminUserID(c); !ok {
		response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "后台身份无效")
		return
	}

	var req CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError, "参数错误")
		return
	}

	var activity models.Activity
	if err := database.DB.Where("id = ?", activityID).First(&activity).Error; err != nil {
		response.Fail(c, http.StatusNotFound, response.CodeNotFound, "活动不存在")
		return
	}

	updates := map[string]interface{}{}
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Image != "" {
		updates["image"] = req.Image
	}
	if req.LinkType != "" {
		updates["link_type"] = req.LinkType
	}
	if req.LinkValue != "" {
		updates["link_value"] = req.LinkValue
	}
	if req.Sort > 0 {
		updates["sort"] = req.Sort
	}
	if req.Status > 0 {
		updates["status"] = req.Status
	}

	if err := database.DB.Model(&activity).Updates(updates).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, "更新活动失败")
		return
	}

	database.DB.First(&activity, activityID)
	response.Success(c, activity)
}

func DeleteActivity(c *gin.Context) {
	id := c.Param("id")
	activityID, _ := strconv.ParseUint(id, 10, 64)

	if _, ok := getCurrentAdminUserID(c); !ok {
		response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "后台身份无效")
		return
	}

	var activity models.Activity
	if err := database.DB.Where("id = ?", activityID).First(&activity).Error; err != nil {
		response.Fail(c, http.StatusNotFound, response.CodeNotFound, "活动不存在")
		return
	}

	if err := database.DB.Model(&activity).Update("status", 0).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, "删除活动失败")
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}
