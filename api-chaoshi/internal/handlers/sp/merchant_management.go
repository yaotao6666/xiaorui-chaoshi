package sp

import (
	"net/http"
	"strconv"
	"strings"
	"chaoshi_api/internal/models"
	"chaoshi_api/pkg/database"
	"chaoshi_api/pkg/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CreateMerchantRequest struct {
	Name                 string  `json:"name" binding:"required"`
	ContactName          string  `json:"contact_name"`
	ContactPhone         string  `json:"contact_phone"`
	ContactEmail         string  `json:"contact_email"`
	Address              string  `json:"address"`
	BusinessCategory     string  `json:"business_category"`
	BusinessHours        string  `json:"business_hours"`
	Announcement         string  `json:"announcement"`
	Username             string  `json:"username" binding:"required"`
	Password             string  `json:"password" binding:"required,min=6"`
	StaffName            string  `json:"staff_name"`
	StaffPhone           string  `json:"staff_phone"`
}

type UpdateMerchantRequest struct {
	Name             *string  `json:"name"`
	ContactName      *string  `json:"contact_name"`
	ContactPhone     *string  `json:"contact_phone"`
	ContactEmail     *string  `json:"contact_email"`
	Address          *string  `json:"address"`
	BusinessCategory *string  `json:"business_category"`
	BusinessHours    *string  `json:"business_hours"`
	Announcement     *string  `json:"announcement"`
	Status           *uint8   `json:"status"`
}

type ResetMerchantAdminPasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func CreateMerchant(c *gin.Context) {
	var req CreateMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError, "参数错误")
		return
	}

	if _, ok := getCurrentAdminUserID(c); !ok {
		response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "后台身份无效")
		return
	}

	var existing models.MerchantStaff
	if err := database.DB.Where("username = ?", strings.TrimSpace(req.Username)).First(&existing).Error; err == nil {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError, "商家登录账号已存在")
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, "生成密码失败")
		return
	}

	merchant := models.Merchant{
		Name:              strings.TrimSpace(req.Name),
		ContactName:       strings.TrimSpace(req.ContactName),
		ContactPhone:      strings.TrimSpace(req.ContactPhone),
		ContactEmail:      strings.TrimSpace(req.ContactEmail),
		Address:           strings.TrimSpace(req.Address),
		BusinessCategory:  strings.TrimSpace(req.BusinessCategory),
		BusinessHours:     strings.TrimSpace(req.BusinessHours),
		Announcement:      strings.TrimSpace(req.Announcement),
		Status:            1,
	}

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&merchant).Error; err != nil {
			return err
		}

		staff := models.MerchantStaff{
			MerchantID: merchant.ID,
			Username:   strings.TrimSpace(req.Username),
			Password:   string(passwordHash),
			Name:       strings.TrimSpace(req.StaffName),
			Phone:      strings.TrimSpace(req.StaffPhone),
			Role:       "owner",
			Status:     1,
		}
		if staff.Name == "" {
			staff.Name = merchant.ContactName
		}
		if staff.Phone == "" {
			staff.Phone = merchant.ContactPhone
		}
		return tx.Create(&staff).Error
	}); err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, "创建商家失败")
		return
	}

	c.Params = append(c.Params, gin.Param{Key: "merchant_id", Value: strconv.FormatUint(merchant.ID, 10)})
	GetMerchantDetail(c)
}

func UpdateMerchant(c *gin.Context) {
	merchantID, _ := strconv.ParseUint(c.Param("merchant_id"), 10, 64)
	if _, ok := getCurrentAdminUserID(c); !ok {
		response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "后台身份无效")
		return
	}

	var merchant models.Merchant
	if err := database.DB.Where("id = ?", merchantID).First(&merchant).Error; err != nil {
		response.Fail(c, http.StatusNotFound, response.CodeNotFound, "商家不存在")
		return
	}

	var req UpdateMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError, "参数错误")
		return
	}

	updates := map[string]any{}
	if req.Name != nil {
		updates["name"] = strings.TrimSpace(*req.Name)
	}
	if req.ContactName != nil {
		updates["contact_name"] = strings.TrimSpace(*req.ContactName)
	}
	if req.ContactPhone != nil {
		updates["contact_phone"] = strings.TrimSpace(*req.ContactPhone)
	}
	if req.ContactEmail != nil {
		updates["contact_email"] = strings.TrimSpace(*req.ContactEmail)
	}
	if req.Address != nil {
		updates["address"] = strings.TrimSpace(*req.Address)
	}
	if req.BusinessCategory != nil {
		updates["business_category"] = strings.TrimSpace(*req.BusinessCategory)
	}
	if req.BusinessHours != nil {
		updates["business_hours"] = strings.TrimSpace(*req.BusinessHours)
	}
	if req.Announcement != nil {
		updates["announcement"] = strings.TrimSpace(*req.Announcement)
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) == 0 {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError, "未提供可更新内容")
		return
	}

	if err := database.DB.Model(&merchant).Updates(updates).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, "更新商家失败")
		return
	}

	GetMerchantDetail(c)
}

func ResetMerchantAdminPassword(c *gin.Context) {
	merchantID, _ := strconv.ParseUint(c.Param("merchant_id"), 10, 64)
	if _, ok := getCurrentAdminUserID(c); !ok {
		response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "后台身份无效")
		return
	}

	var merchant models.Merchant
	if err := database.DB.Select("id").
		Where("id = ?", merchantID).
		First(&merchant).Error; err != nil {
		response.Fail(c, http.StatusNotFound, response.CodeNotFound, "商家不存在")
		return
	}

	var req ResetMerchantAdminPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError, "参数错误")
		return
	}

	var adminStaff models.MerchantStaff
	if err := database.DB.
		Where("merchant_id = ? AND role = ?", merchantID, "owner").
		Order("status DESC, id ASC").
		First(&adminStaff).Error; err != nil {
		response.Fail(c, http.StatusNotFound, response.CodeNotFound, "商家管理员不存在")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, "密码加密失败")
		return
	}

	if err := database.DB.Model(&models.MerchantStaff{}).
		Where("id = ?", adminStaff.ID).
		Update("password", string(hashedPassword)).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, "重置管理员密码失败")
		return
	}

	response.Success(c, gin.H{
		"staff_id": adminStaff.ID,
		"username": adminStaff.Username,
		"message":  "管理员密码已重置",
	})
}
