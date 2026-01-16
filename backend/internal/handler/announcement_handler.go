package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/study-upc/backend/internal/model"
	"github.com/study-upc/backend/internal/middleware"
	"github.com/study-upc/backend/internal/pkg/response"
	"github.com/study-upc/backend/internal/service"
)

// AnnouncementHandler 公告处理器
type AnnouncementHandler struct {
	announcementService service.AnnouncementService
}

// NewAnnouncementHandler 创建公告处理器实例
func NewAnnouncementHandler(announcementService service.AnnouncementService) *AnnouncementHandler {
	return &AnnouncementHandler{
		announcementService: announcementService,
	}
}

// CreateAnnouncement 创建公告
// @Summary 创建公告
// @Description 管理员创建新公告
// @Tags 公告管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.CreateAnnouncementRequest true "公告信息"
// @Success 200 {object} response.Response{data=model.AnnouncementResponse}
// @Router /api/v1/announcements [post]
func (h *AnnouncementHandler) CreateAnnouncement(c *gin.Context) {
	var req model.CreateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrInvalidParams, err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.ErrUnauthorized, "未授权")
		return
	}

	announcementResp, err := h.announcementService.CreateAnnouncement(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, announcementResp)
}

// GetAnnouncement 获取公告详情
// @Summary 获取公告详情
// @Description 根据ID获取公告详情
// @Tags 公告管理
// @Accept json
// @Produce json
// @Param id path int true "公告ID"
// @Success 200 {object} response.Response{data=model.AnnouncementResponse}
// @Router /api/v1/announcements/{id} [get]
func (h *AnnouncementHandler) GetAnnouncement(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.ErrInvalidParams, "无效的公告ID")
		return
	}

	announcementResp, err := h.announcementService.GetAnnouncement(c.Request.Context(), uint(id))
	if err != nil {
		if err == service.ErrAnnouncementNotFound {
			response.Error(c, response.ErrNotFound, "公告不存在")
			return
		}
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, announcementResp)
}

// GetActiveAnnouncements 获取启用的公告列表
// @Summary 获取启用的公告列表
// @Description 获取当前启用且有效的公告列表（用于首页公告栏）
// @Tags 公告管理
// @Accept json
// @Produce json
// @Param limit query int false "限制数量" default(5)
// @Success 200 {object} response.Response{data=[]model.Announcement}
// @Router /api/v1/announcements/active [get]
func (h *AnnouncementHandler) GetActiveAnnouncements(c *gin.Context) {
	limit := 5
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	announcements, err := h.announcementService.GetActiveAnnouncements(c.Request.Context(), limit)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, announcements)
}

// ListAnnouncements 查询公告列表
// @Summary 查询公告列表
// @Description 管理员查询公告列表（支持筛选、分页）
// @Tags 公告管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param priority query string false "优先级"
// @Param is_active query bool false "是否启用"
// @Param author_id query int false "发布者ID"
// @Success 200 {object} response.Response{data=response.PaginateData}
// @Router /api/v1/announcements [get]
func (h *AnnouncementHandler) ListAnnouncements(c *gin.Context) {
	var req model.AnnouncementListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, response.ErrInvalidParams, err.Error())
		return
	}

	announcements, total, err := h.announcementService.ListAnnouncements(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.SuccessWithPaginate(c, total, req.Page, req.PageSize, announcements)
}

// UpdateAnnouncement 更新公告
// @Summary 更新公告
// @Description 管理员更新公告
// @Tags 公告管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "公告ID"
// @Param request body model.UpdateAnnouncementRequest true "公告信息"
// @Success 200 {object} response.Response{data=model.AnnouncementResponse}
// @Router /api/v1/announcements/{id} [put]
func (h *AnnouncementHandler) UpdateAnnouncement(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.ErrInvalidParams, "无效的公告ID")
		return
	}

	var req model.UpdateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrInvalidParams, err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.ErrUnauthorized, "未授权")
		return
	}

	announcementResp, err := h.announcementService.UpdateAnnouncement(c.Request.Context(), uint(id), userID.(uint), middleware.IsAdmin(c), &req)
	if err != nil {
		if err == service.ErrAnnouncementNotFound {
			response.Error(c, response.ErrNotFound, "公告不存在")
			return
		}
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, announcementResp)
}

// DeleteAnnouncement 删除公告
// @Summary 删除公告
// @Description 管理员删除公告
// @Tags 公告管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "公告ID"
// @Success 200 {object} response.Response
// @Router /api/v1/announcements/{id} [delete]
func (h *AnnouncementHandler) DeleteAnnouncement(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.ErrInvalidParams, "无效的公告ID")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.ErrUnauthorized, "未授权")
		return
	}

	err = h.announcementService.DeleteAnnouncement(c.Request.Context(), uint(id), userID.(uint), middleware.IsAdmin(c))
	if err != nil {
		if err == service.ErrAnnouncementNotFound {
			response.Error(c, response.ErrNotFound, "公告不存在")
			return
		}
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, nil)
}
