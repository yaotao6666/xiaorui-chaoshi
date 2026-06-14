package sp

import (
	"chaoshi_api/internal/models"
	"chaoshi_api/pkg/database"
	"chaoshi_api/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getCurrentAdminUserID(c *gin.Context) (uint64, bool) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	adminUserID, ok := userIDValue.(uint64)
	if !ok {
		return 0, false
	}

	var adminUser models.AdminUser
	if err := database.DB.Select("id").First(&adminUser, adminUserID).Error; err != nil {
		return 0, false
	}
	return adminUser.ID, true
}

func GetAnnouncements(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	if _, ok := getCurrentAdminUserID(c); !ok {
		response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "后台身份无效")
		return
	}

	var total int64
	database.DB.Model(&models.Announcement{}).Count(&total)

	var announcements []models.Announcement
	offset := (page - 1) * pageSize
	if err := database.DB.
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&announcements).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, "获取公告列表失败")
		return
	}

	response.Success(c, gin.H{
		"list": announcements,
		"pagination": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

type CreateAnnouncementRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Status  uint8  `json:"status" binding:"omitempty,oneof=0 1"`
}

func GetAnnouncementDetail(c *gin.Context) {
	id := c.Param("id")
	announcementID, _ := strconv.ParseUint(id, 10, 64)

	if _, ok := getCurrentAdminUserID(c); !ok {
		response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "后台身份无效")
		return
	}

	var announcement models.Announcement
	if err := database.DB.Where("id = ?", announcementID).First(&announcement).Error; err != nil {
		response.Fail(c, http.StatusNotFound, response.CodeNotFound, "公告不存在")
		return
	}

	response.Success(c, announcement)
}

func CreateAnnouncement(c *gin.Context) {
	var req CreateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError, "参数错误")
		return
	}

	if _, ok := getCurrentAdminUserID(c); !ok {
		response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "后台身份无效")
		return
	}

	status := uint8(1)
	if req.Status == 0 {
		status = 0
	}

	announcement := models.Announcement{
		Title:   req.Title,
		Content: req.Content,
		Status:  status,
	}

	if err := database.DB.Create(&announcement).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, "创建公告失败")
		return
	}

	response.Success(c, gin.H{"id": announcement.ID, "message": "创建成功"})
}

func UpdateAnnouncement(c *gin.Context) {
	id := c.Param("id")
	announcementID, _ := strconv.ParseUint(id, 10, 64)

	var req CreateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError, "参数错误")
		return
	}

	if _, ok := getCurrentAdminUserID(c); !ok {
		response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "后台身份无效")
		return
	}

	var announcement models.Announcement
	if err := database.DB.Where("id = ?", announcementID).First(&announcement).Error; err != nil {
		response.Fail(c, http.StatusNotFound, response.CodeNotFound, "公告不存在")
		return
	}

	updates := map[string]interface{}{
		"title":   req.Title,
		"content": req.Content,
		"status":  req.Status,
	}

	if err := database.DB.Model(&announcement).Updates(updates).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, "更新公告失败")
		return
	}

	database.DB.First(&announcement, announcementID)
	response.Success(c, announcement)
}

func DeleteAnnouncement(c *gin.Context) {
	id := c.Param("id")
	announcementID, _ := strconv.ParseUint(id, 10, 64)

	if _, ok := getCurrentAdminUserID(c); !ok {
		response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "后台身份无效")
		return
	}

	var announcement models.Announcement
	if err := database.DB.Where("id = ?", announcementID).First(&announcement).Error; err != nil {
		response.Fail(c, http.StatusNotFound, response.CodeNotFound, "公告不存在")
		return
	}

	if err := database.DB.Model(&announcement).Update("status", 0).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeServerError, "删除公告失败")
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}
