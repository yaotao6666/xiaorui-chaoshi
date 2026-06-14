package upload

import (
	"fmt"
	"net/http"

	"chaoshi_api/internal/middleware"
	"chaoshi_api/internal/storage"
	"chaoshi_api/pkg/response"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

func (h *UploadHandler) UploadFile(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError, "请选择要上传的文件")
		return
	}

	service := storage.GetService()
	if service == nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, "文件存储服务未初始化")
		return
	}

	scope := resolveUploadScope(c)
	path, url, err := service.SaveUploadedFile(fileHeader, scope)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, "保存上传文件失败")
		return
	}

	response.Success(c, gin.H{
		"path":     path,
		"url":      url,
		"filename": fileHeader.Filename,
		"scope":    scope,
	})
}

func resolveUploadScope(c *gin.Context) string {
	userType := middleware.GetUserType(c)
	userID := middleware.GetUserID(c)

	switch userType {
	case "merchant":
		return fmt.Sprintf("store/%d", userID)
	case "sp":
		return "admin"
	case "user":
		return fmt.Sprintf("user/%d", userID)
	default:
		return "common"
	}
}
