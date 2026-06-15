package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"chaoshi_api/internal/models"
	"chaoshi_api/internal/services/payment/jsbank"
	"chaoshi_api/pkg/database"
	"chaoshi_api/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PrepareOrderPaymentRequest struct {
	ReturnPath string `json:"return_path"`
	Source     string `json:"source"`
}

func PrepareOrderPayment(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("order_id"), 10, 64)
	if err != nil || orderID == 0 {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError, "订单参数错误")
		return
	}

	var req PrepareOrderPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError, "参数错误")
		return
	}

	if strings.TrimSpace(req.Source) != "xcx_shell" {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError, "当前仅支持小程序壳发起支付")
		return
	}

	userID := GetCurrentUserID(c)
	if userID == 0 {
		response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "请先登录")
		return
	}

	var order models.Order
	if err := database.DB.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, http.StatusNotFound, response.CodeNotFound, "订单不存在")
			return
		}
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, "获取订单失败")
		return
	}

	if order.Status != 1 {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError, "当前订单状态不可发起支付")
		return
	}
	if order.PayAmount <= 0 {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError, "当前订单无需支付")
		return
	}

	var currentUser models.User
	if err := database.DB.First(&currentUser, userID).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, "获取用户信息失败")
		return
	}
	if strings.TrimSpace(currentUser.OpenID) == "" {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError, "当前用户缺少 openid")
		return
	}

	client, err := jsbank.NewClient()
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, err.Error())
		return
	}

	payResult, err := client.CreateMiniProgramPay(jsbank.MiniProgramPayRequest{
		OrderNo: order.OrderNo,
		Amount:  order.PayAmount,
		OpenID:  currentUser.OpenID,
	})
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"order_id":      order.ID,
		"order_no":      order.OrderNo,
		"merchant_id":   order.MerchantID,
		"pay_amount":    order.PayAmount,
		"channel":       "jsbank_wechat_mini",
		"pay_params":    payResult.PayParams,
		"return_target": buildOrderReturnTarget(req.ReturnPath, order.ID, order.MerchantID),
	})
}

func HandleJsBankPayNotify(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		c.String(http.StatusBadRequest, "fail")
		return
	}

	params := make(map[string]string)
	for key, values := range c.Request.PostForm {
		if len(values) == 0 {
			continue
		}
		params[key] = values[0]
	}
	if len(params) == 0 {
		for key, values := range c.Request.Form {
			if len(values) == 0 {
				continue
			}
			params[key] = values[0]
		}
	}
	if len(params) == 0 {
		c.String(http.StatusBadRequest, "fail")
		return
	}

	client, err := jsbank.NewClient()
	if err != nil {
		log.Printf("江苏银行回调初始化失败: %v", err)
		c.String(http.StatusInternalServerError, "fail")
		return
	}
	if err := client.VerifyNotify(params); err != nil {
		log.Printf("江苏银行回调验签失败: %v", err)
		c.String(http.StatusBadRequest, "fail")
		return
	}

	orderNo := firstNonEmpty(
		params["outTradeNo"],
		params["extfld2"],
		params["orderNo"],
	)
	if orderNo == "" {
		log.Printf("江苏银行回调缺少商户订单号: %+v", params)
		c.String(http.StatusBadRequest, "fail")
		return
	}

	if !isJsBankNotifySuccess(params) {
		log.Printf("江苏银行回调未返回成功状态: %+v", params)
		c.String(http.StatusOK, "success")
		return
	}

	payloadJSON, _ := json.Marshal(params)
	transactionID := firstNonEmpty(
		params["tradeNo"],
		params["transactionId"],
		params["payNo"],
		params["bankOrderNo"],
	)
	now := time.Now()

	tx := database.DB.Begin()

	var order models.Order
	if err := tx.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		tx.Rollback()
		log.Printf("江苏银行回调未找到订单: %s, err=%v", orderNo, err)
		c.String(http.StatusBadRequest, "fail")
		return
	}

	if order.Status != 2 {
		updates := map[string]interface{}{
			"status":             2,
			"paid_at":            &now,
			"pay_notify_payload": models.JSON(payloadJSON),
		}
		if transactionID != "" {
			updates["transaction_id"] = transactionID
		}
		if err := tx.Model(&models.Order{}).Where("id = ?", order.ID).Updates(updates).Error; err != nil {
			tx.Rollback()
			log.Printf("江苏银行回调更新订单失败: order_no=%s err=%v", orderNo, err)
			c.String(http.StatusInternalServerError, "fail")
			return
		}

		if err := tx.Model(&models.User{}).Where("id = ?", order.UserID).Updates(map[string]interface{}{
			"has_paid":      true,
			"first_paid_at": gorm.Expr("CASE WHEN first_paid_at IS NULL THEN ? ELSE first_paid_at END", now),
		}).Error; err != nil {
			tx.Rollback()
			log.Printf("江苏银行回调更新用户支付状态失败: order_no=%s err=%v", orderNo, err)
			c.String(http.StatusInternalServerError, "fail")
			return
		}
	} else {
		if err := tx.Model(&models.Order{}).Where("id = ?", order.ID).Update("pay_notify_payload", models.JSON(payloadJSON)).Error; err != nil {
			tx.Rollback()
			log.Printf("江苏银行回调写入原始报文失败: order_no=%s err=%v", orderNo, err)
			c.String(http.StatusInternalServerError, "fail")
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Printf("江苏银行回调提交事务失败: order_no=%s err=%v", orderNo, err)
		c.String(http.StatusInternalServerError, "fail")
		return
	}

	c.String(http.StatusOK, "success")
}

func buildOrderReturnTarget(returnPath string, orderID, merchantID uint64) string {
	trimmedPath := strings.TrimSpace(returnPath)
	if trimmedPath == "" {
		trimmedPath = "/pages/store/order-detail"
	}
	if !strings.HasPrefix(trimmedPath, "/") {
		trimmedPath = "/" + trimmedPath
	}

	parsedURL, err := url.Parse(trimmedPath)
	if err != nil {
		return fmt.Sprintf("/pages/store/order-detail?order_id=%d&merchant_id=%d", orderID, merchantID)
	}

	query := parsedURL.Query()
	query.Set("order_id", strconv.FormatUint(orderID, 10))
	query.Set("merchant_id", strconv.FormatUint(merchantID, 10))
	parsedURL.RawQuery = query.Encode()
	return parsedURL.String()
}

func isJsBankNotifySuccess(params map[string]string) bool {
	respCode := strings.TrimSpace(params["respCode"])
	orderStatus := strings.TrimSpace(params["orderStatus"])
	if respCode != "" && respCode != "000000" {
		return false
	}
	if orderStatus == "" {
		return true
	}
	return orderStatus == "1"
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func GetCurrentUserID(c *gin.Context) uint64 {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		return 0
	}

	switch typedValue := userIDValue.(type) {
	case uint64:
		return typedValue
	case uint:
		return uint64(typedValue)
	case int:
		if typedValue > 0 {
			return uint64(typedValue)
		}
	case float64:
		if typedValue > 0 {
			return uint64(typedValue)
		}
	}

	return 0
}
