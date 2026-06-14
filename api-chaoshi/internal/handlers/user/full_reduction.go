package user

import (
	"net/http"
	"strconv"

	"chaoshi_api/internal/services/fullreduction"
	"chaoshi_api/pkg/response"

	"github.com/gin-gonic/gin"
)

func GetStoreFullReductionRules(c *gin.Context) {
	merchantID, _ := strconv.ParseUint(c.Param("merchant_id"), 10, 64)
	if merchantID == 0 {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError, "商家参数错误")
		return
	}

	rules, err := fullreduction.GetActiveRulesByMerchantID(merchantID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, "获取满减规则失败")
		return
	}

	response.Success(c, gin.H{
		"rules": rules,
	})
}
