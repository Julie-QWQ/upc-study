package handler

import (
	"strconv"

	"github.com/study-upc/backend/internal/model"
	"github.com/study-upc/backend/internal/pkg/response"
	"github.com/study-upc/backend/internal/repository"
	"github.com/study-upc/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// MaterialHandler 资料处理器
type MaterialHandler struct {
	materialService  service.MaterialService
	favoriteService  service.FavoriteService
	reportService    service.ReportService
	downloadRecordRepo repository.DownloadRecordRepository
}

// NewMaterialHandler 创建资料处理器实例
func NewMaterialHandler(
	materialService service.MaterialService,
	favoriteService service.FavoriteService,
	reportService service.ReportService,
	downloadRecordRepo repository.DownloadRecordRepository,
) *MaterialHandler {
	return &MaterialHandler{
		materialService:  materialService,
		favoriteService:  favoriteService,
		reportService:    reportService,
		downloadRecordRepo: downloadRecordRepo,
	}
}

// CreateMaterial 创建资料
// @Summary 创建资料
// @Description 学委上传学习资料
// @Tags 资料
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.CreateMaterialRequest true "资料信息"
// @Success 200 {object} response.Response{data=model.MaterialResponse}
// @Router /api/v1/materials [post]
func (h *MaterialHandler) CreateMaterial(c *gin.Context) {
	var req model.CreateMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrInvalidParams, err.Error())
		return
	}

	// 从上下文获取用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.ErrUnauthorized, "未认证")
		return
	}

	materialResp, err := h.materialService.CreateMaterial(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		switch err {
		case service.ErrInvalidMaterialStatus:
			response.Error(c, response.ErrInvalidParams, err.Error())
		default:
			response.Error(c, response.ErrInternal, err.Error())
		}
		return
	}

	response.Success(c, materialResp)
}

// UpdateMaterial 更新资料
// @Summary 更新资料
// @Description 学委更新自己上传的资料
// @Tags 资料
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "资料ID"
// @Param request body model.UpdateMaterialRequest true "资料信息"
// @Success 200 {object} response.Response{data=model.MaterialResponse}
// @Router /api/v1/materials/{id} [put]
func (h *MaterialHandler) UpdateMaterial(c *gin.Context) {
	materialID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.ErrInvalidParams, "无效的资料ID")
		return
	}

	var req model.UpdateMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrInvalidParams, err.Error())
		return
	}

	// 从上下文获取用户 ID 和角色
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.ErrUnauthorized, "未认证")
		return
	}

	userRole, exists := c.Get("user_role")
	if !exists {
		response.Error(c, response.ErrUnauthorized, "未认证")
		return
	}

	materialResp, err := h.materialService.UpdateMaterial(c.Request.Context(), uint(materialID), userID.(uint), userRole.(string), &req)
	if err != nil {
		switch err {
		case service.ErrMaterialNotFound:
			response.Error(c, response.ErrNotFound, err.Error())
		case service.ErrAccessDenied:
			response.Error(c, response.ErrForbidden, err.Error())
		case service.ErrMaterialAlreadyApproved:
			response.Error(c, response.ErrInvalidParams, err.Error())
		default:
			response.Error(c, response.ErrInternal, err.Error())
		}
		return
	}

	response.Success(c, materialResp)
}

// GetMaterial 获取资料详情
// @Summary 获取资料详情
// @Description 获取资料详细信息
// @Tags 资料
// @Produce json
// @Param id path int true "资料ID"
// @Success 200 {object} response.Response{data=model.MaterialResponse}
// @Router /api/v1/materials/{id} [get]
func (h *MaterialHandler) GetMaterial(c *gin.Context) {
	materialID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.ErrInvalidParams, "无效的资料ID")
		return
	}

	// 获取当前用户 ID（可能为空）
	var currentUserID uint
	if userID, exists := c.Get("user_id"); exists {
		currentUserID = userID.(uint)
	}

	materialResp, err := h.materialService.GetMaterial(c.Request.Context(), uint(materialID), currentUserID)
	if err != nil {
		switch err {
		case service.ErrMaterialNotFound:
			response.Error(c, response.ErrNotFound, err.Error())
		case service.ErrAccessDenied:
			response.Error(c, response.ErrForbidden, err.Error())
		default:
			response.Error(c, response.ErrInternal, err.Error())
		}
		return
	}

	response.Success(c, materialResp)
}

// ListMaterials 获取资料列表
// @Summary 获取资料列表
// @Description 分页获取资料列表，支持筛选和排序
// @Tags 资料
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param category query string false "分类"
// @Param course_name query string false "课程名称"
// @Param status query string false "状态"
// @Param keyword query string false "搜索关键词"
// @Param sort_by query string false "排序字段" default(created_at)
// @Param sort_order query string false "排序方向" default(desc)
// @Success 200 {object} response.Response{data=model.MaterialListResponse}
// @Router /api/v1/materials [get]
func (h *MaterialHandler) ListMaterials(c *gin.Context) {
	var req model.MaterialListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, response.ErrInvalidParams, err.Error())
		return
	}

	// 获取当前用户 ID（可能为空）
	var currentUserID uint
	if userID, exists := c.Get("user_id"); exists {
		currentUserID = userID.(uint)
	}

	// 获取当前用户角色（可能为空）
	var currentUserRole string
	if role, exists := c.Get("user_role"); exists {
		currentUserRole = role.(string)
	}

	listResp, err := h.materialService.ListMaterials(c.Request.Context(), &req, currentUserID, currentUserRole)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, listResp)
}

// ListReviewedMaterials 获取已审核资料列表(管理员专用)
// @Summary 获取已审核资料列表
// @Description 管理员获取已审核资料列表
// @Tags 资料
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param sort_by query string false "排序字段" default(created_at)
// @Param sort_order query string false "排序方向" default(desc)
// @Success 200 {object} response.Response{data=model.MaterialListResponse}
// @Router /api/v1/admin/materials/reviewed [get]
func (h *MaterialHandler) ListReviewedMaterials(c *gin.Context) {
	var req model.MaterialListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, response.ErrInvalidParams, err.Error())
		return
	}
	req.ReviewedOnly = true

	// 获取当前用户 ID（可能为空）
	var currentUserID uint
	if userID, exists := c.Get("user_id"); exists {
		currentUserID = userID.(uint)
	}

	// 获取当前用户角色（可能为空）
	var currentUserRole string
	if role, exists := c.Get("user_role"); exists {
		currentUserRole = role.(string)
	}

	listResp, err := h.materialService.ListMaterials(c.Request.Context(), &req, currentUserID, currentUserRole)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, listResp)
}

// ListPendingMaterials 获取待审核资料列表(管理员专用)
// @Summary 获取待审核资料列表
// @Description 管理员获取待审核的资料列表
// @Tags 资料
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param sort_by query string false "排序字段" default(created_at)
// @Param sort_order query string false "排序方向" default(desc)
// @Success 200 {object} response.Response{data=model.MaterialListResponse}
// @Router /api/v1/admin/materials/pending [get]
func (h *MaterialHandler) ListPendingMaterials(c *gin.Context) {
	var req model.MaterialListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, response.ErrInvalidParams, err.Error())
		return
	}

	// 强制设置状态为 pending
	req.Status = model.StatusPending

	// 获取当前用户 ID(管理员)
	var currentUserID uint
	if userID, exists := c.Get("user_id"); exists {
		currentUserID = userID.(uint)
	}

	listResp, err := h.materialService.ListMaterials(c.Request.Context(), &req, currentUserID, "admin")
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, listResp)
}

// DeleteMaterial 删除资料
// @Summary 删除资料
// @Description 管理员删除资料
// @Tags 资料
// @Produce json
// @Security BearerAuth
// @Param id path int true "资料ID"
// @Success 200 {object} response.Response
// @Router /api/v1/materials/{id} [delete]
func (h *MaterialHandler) DeleteMaterial(c *gin.Context) {
	materialID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.ErrInvalidParams, "无效的资料ID")
		return
	}

	if err := h.materialService.DeleteMaterial(c.Request.Context(), uint(materialID)); err != nil {
		switch err {
		case service.ErrMaterialNotFound:
			response.Error(c, response.ErrNotFound, err.Error())
		default:
			response.Error(c, response.ErrInternal, err.Error())
		}
		return
	}

	response.Success(c, nil)
}

// ReviewMaterial 审核资料
// @Summary 审核资料
// @Description 管理员审核资料
// @Tags 资料
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "资料ID"
// @Param request body model.ReviewMaterialRequest true "审核信息"
// @Success 200 {object} response.Response
// @Router /api/v1/materials/{id}/review [post]
func (h *MaterialHandler) ReviewMaterial(c *gin.Context) {
	materialID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.ErrInvalidParams, "无效的资料ID")
		return
	}

	var req model.ReviewMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrInvalidParams, err.Error())
		return
	}

	// 从上下文获取用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.ErrUnauthorized, "未认证")
		return
	}

	if err := h.materialService.ReviewMaterial(c.Request.Context(), uint(materialID), userID.(uint), &req); err != nil {
		switch err {
		case service.ErrMaterialNotFound:
			response.Error(c, response.ErrNotFound, err.Error())
		case service.ErrMaterialAlreadyReviewed:
			response.Error(c, response.ErrInvalidParams, err.Error())
		default:
			response.Error(c, response.ErrInternal, err.Error())
		}
		return
	}

	response.Success(c, nil)
}

// GetUploadSignature 获取上传签名
// @Summary 获取上传签名
// @Description 获取文件上传的预签名URL
// @Tags 资料
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.UploadSignatureRequest true "文件信息"
// @Success 200 {object} response.Response{data=model.UploadSignatureResponse}
// @Router /api/v1/materials/upload-signature [post]
func (h *MaterialHandler) GetUploadSignature(c *gin.Context) {
	var req model.UploadSignatureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrInvalidParams, err.Error())
		return
	}

	// 从上下文获取用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.ErrUnauthorized, "未认证")
		return
	}

	signatureResp, err := h.materialService.GetUploadSignature(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, signatureResp)
}

// DeleteUploadedFile 删除已上传但未创建记录的文件
// @Summary 删除已上传文件
// @Description 删除已上传到OSS但未创建资料记录的文件
// @Tags 资料
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.DeleteUploadedFileRequest true "文件信息"
// @Success 200 {object} response.Response
// @Router /api/v1/materials/delete-uploaded-file [post]
func (h *MaterialHandler) DeleteUploadedFile(c *gin.Context) {
	var req model.DeleteUploadedFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrInvalidParams, err.Error())
		return
	}

	// 从上下文获取用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.ErrUnauthorized, "未认证")
		return
	}

	// 验证fileKey是否属于当前用户(格式: materials/{userID}/{uuid}_{fileName})
	// 这里简单验证前缀
	if err := h.materialService.DeleteUploadedFile(c.Request.Context(), userID.(uint), req.FileKey); err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetDownloadURL 获取下载链接
// @Summary 获取下载链接
// @Description 获取资料下载的预签名URL
// @Tags 资料
// @Produce json
// @Security BearerAuth
// @Param id path int true "资料ID"
// @Success 200 {object} response.Response{data=map[string]string}
// @Router /api/v1/materials/{id}/download [get]
func (h *MaterialHandler) GetDownloadURL(c *gin.Context) {
	materialID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.ErrInvalidParams, "无效的资料ID")
		return
	}

	// 从上下文获取用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.ErrUnauthorized, "未认证")
		return
	}

	downloadURL, err := h.materialService.GetDownloadURL(c.Request.Context(), uint(materialID), userID.(uint))
	if err != nil {
		switch err {
		case service.ErrMaterialNotFound:
			response.Error(c, response.ErrNotFound, err.Error())
		case service.ErrAccessDenied:
			response.Error(c, response.ErrForbidden, err.Error())
		case service.ErrDownloadLimitExceeded:
			response.Error(c, response.ErrForbidden, err.Error())
		default:
			response.Error(c, response.ErrInternal, err.Error())
		}
		return
	}

	response.Success(c, map[string]string{"download_url": downloadURL})
}

// SearchMaterials 搜索资料
// @Summary 搜索资料
// @Description 全文搜索资料
// @Tags 资料
// @Produce json
// @Param keyword query string true "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=model.MaterialListResponse}
// @Router /api/v1/materials/search [get]
func (h *MaterialHandler) SearchMaterials(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		response.Error(c, response.ErrInvalidParams, "搜索关键词不能为空")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	searchResp, err := h.materialService.SearchMaterials(c.Request.Context(), keyword, page, pageSize)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, searchResp)
}

// AddFavorite 添加收藏
// @Summary 添加收藏
// @Description 收藏资料
// @Tags 收藏
// @Produce json
// @Security BearerAuth
// @Param id path int true "资料ID"
// @Success 200 {object} response.Response
// @Router /api/v1/materials/{id}/favorite [post]
func (h *MaterialHandler) AddFavorite(c *gin.Context) {
	materialID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.ErrInvalidParams, "无效的资料ID")
		return
	}

	// 从上下文获取用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.ErrUnauthorized, "未认证")
		return
	}

	if err := h.favoriteService.AddFavorite(c.Request.Context(), userID.(uint), uint(materialID)); err != nil {
		switch err {
		case service.ErrMaterialNotFound:
			response.Error(c, response.ErrNotFound, err.Error())
		case service.ErrAlreadyFavorited:
			response.Error(c, response.ErrDuplicate, err.Error())
		case service.ErrAccessDenied:
			response.Error(c, response.ErrForbidden, err.Error())
		default:
			response.Error(c, response.ErrInternal, err.Error())
		}
		return
	}

	response.Success(c, nil)
}

// RemoveFavorite 取消收藏
// @Summary 取消收藏
// @Description 取消收藏资料
// @Tags 收藏
// @Produce json
// @Security BearerAuth
// @Param id path int true "资料ID"
// @Success 200 {object} response.Response
// @Router /api/v1/materials/{id}/favorite [delete]
func (h *MaterialHandler) RemoveFavorite(c *gin.Context) {
	materialID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.ErrInvalidParams, "无效的资料ID")
		return
	}

	// 从上下文获取用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.ErrUnauthorized, "未认证")
		return
	}

	if err := h.favoriteService.RemoveFavorite(c.Request.Context(), userID.(uint), uint(materialID)); err != nil {
		switch err {
		case service.ErrFavoriteNotFound:
			response.Error(c, response.ErrNotFound, err.Error())
		default:
			response.Error(c, response.ErrInternal, err.Error())
		}
		return
	}

	response.Success(c, nil)
}

// ListFavorites 获取收藏列表
// @Summary 获取收藏列表
// @Description 获取当前用户的收藏列表
// @Tags 收藏
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=[]model.FavoriteResponse}
// @Router /api/v1/favorites [get]
func (h *MaterialHandler) ListFavorites(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 从上下文获取用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.ErrUnauthorized, "未认证")
		return
	}

	favorites, total, err := h.favoriteService.ListFavorites(c.Request.Context(), userID.(uint), page, pageSize)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	// 提取资料信息用于前端展示，转换为 map 以便添加额外字段
	materials := make([]map[string]interface{}, 0, len(favorites))
	for _, fav := range favorites {
		if fav.Material != nil {
			materialMap := map[string]interface{}{
				"id":               fav.Material.ID,
				"title":            fav.Material.Title,
				"description":      fav.Material.Description,
				"category":         fav.Material.Category,
				"course_name":      fav.Material.CourseName,
				"uploader_id":      fav.Material.UploaderID,
				"uploader":         fav.Material.Uploader,
				"status":           fav.Material.Status,
				"file_name":        fav.Material.FileName,
				"file_size":        fav.Material.FileSize,
				"mime_type":        fav.Material.MimeType,
				"download_count":   fav.Material.DownloadCount,
				"favorite_count":   fav.Material.FavoriteCount,
				"view_count":       fav.Material.ViewCount,
				"reviewer_id":      fav.Material.ReviewerID,
				"reviewer":         fav.Material.Reviewer,
				"reviewed_at":      fav.Material.ReviewedAt,
				"rejection_reason": fav.Material.RejectionReason,
				"created_at":       fav.Material.CreatedAt,
				"updated_at":       fav.Material.UpdatedAt,
				"is_favorited":     true,
				"favorited_at":     fav.CreatedAt, // 添加收藏时间
			}
			materials = append(materials, materialMap)
		}
	}

	response.SuccessWithPaginate(c, total, page, pageSize, materials)
}

// CreateReport 创建举报
// @Summary 创建举报
// @Description 举报资料
// @Tags 举报
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "资料ID"
// @Param request body model.ReportRequest true "举报信息"
// @Success 200 {object} response.Response
// @Router /api/v1/materials/{id}/report [post]
func (h *MaterialHandler) CreateReport(c *gin.Context) {
	materialID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.ErrInvalidParams, "无效的资料ID")
		return
	}

	var req model.ReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrInvalidParams, err.Error())
		return
	}

	// 从上下文获取用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.ErrUnauthorized, "未认证")
		return
	}

	if err := h.reportService.CreateReport(c.Request.Context(), userID.(uint), uint(materialID), &req); err != nil {
		switch err {
		case service.ErrMaterialNotFound:
			response.Error(c, response.ErrNotFound, err.Error())
		case service.ErrAccessDenied:
			response.Error(c, response.ErrForbidden, err.Error())
		default:
			response.Error(c, response.ErrInternal, err.Error())
		}
		return
	}

	response.Success(c, nil)
}

// HandleReport 处理举报
// @Summary 处理举报
// @Description 管理员处理举报
// @Tags 举报
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "举报ID"
// @Param request body model.HandleReportRequest true "处理信息"
// @Success 200 {object} response.Response
// @Router /api/v1/reports/{id}/handle [post]
func (h *MaterialHandler) HandleReport(c *gin.Context) {
	reportID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.ErrInvalidParams, "无效的举报ID")
		return
	}

	var req model.HandleReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrInvalidParams, err.Error())
		return
	}

	// 从上下文获取用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.ErrUnauthorized, "未认证")
		return
	}

	if err := h.reportService.HandleReport(c.Request.Context(), uint(reportID), userID.(uint), &req); err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, nil)
}

// ListReports 获取举报列表
// @Summary 获取举报列表
// @Description 管理员获取举报列表
// @Tags 举报
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param status query string false "状态筛选"
// @Success 200 {object} response.Response{data=[]model.ReportResponse}
// @Router /api/v1/reports [get]
func (h *MaterialHandler) ListReports(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	statusStr := c.Query("status")

	var status *model.ReportStatus
	if statusStr != "" {
		s := model.ReportStatus(statusStr)
		status = &s
	}

	reports, total, err := h.reportService.ListReports(c.Request.Context(), page, pageSize, status)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"total":       total,
		"page":        page,
		"size":        pageSize,
		"total_pages": (int(total) + pageSize - 1) / pageSize,
		"list":        reports,
	})
}

// GetReport 获取举报详情
// @Summary 获取举报详情
// @Description 获取举报详细信息
// @Tags 举报
// @Produce json
// @Security BearerAuth
// @Param id path int true "举报ID"
// @Success 200 {object} response.Response{data=model.ReportResponse}
// @Router /api/v1/reports/{id} [get]
func (h *MaterialHandler) GetReport(c *gin.Context) {
	reportID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.ErrInvalidParams, "无效的举报ID")
		return
	}

	report, err := h.reportService.GetReport(c.Request.Context(), uint(reportID))
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, report)
}

// ListDownloadRecords 获取下载记录列表
// @Summary 获取下载记录列表
// @Description 获取当前用户的下载记录列表
// @Tags 下载记录
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=[]model.MaterialResponse}
// @Router /api/v1/downloads [get]
func (h *MaterialHandler) ListDownloadRecords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 从上下文获取用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.ErrUnauthorized, "未认证")
		return
	}

	records, total, err := h.downloadRecordRepo.ListByUser(c.Request.Context(), userID.(uint), page, pageSize)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	// 提取资料信息用于前端展示
	materials := make([]map[string]interface{}, 0, len(records))
	for _, record := range records {
		if record.Material != nil {
			materialMap := map[string]interface{}{
				"id":               record.Material.ID,
				"title":            record.Material.Title,
				"description":      record.Material.Description,
				"category":         record.Material.Category,
				"course_name":      record.Material.CourseName,
				"uploader_id":      record.Material.UploaderID,
				"uploader":         record.Material.Uploader,
				"status":           record.Material.Status,
				"file_name":        record.Material.FileName,
				"file_size":        record.Material.FileSize,
				"mime_type":        record.Material.MimeType,
				"download_count":   record.Material.DownloadCount,
				"favorite_count":   record.Material.FavoriteCount,
				"view_count":       record.Material.ViewCount,
				"reviewer_id":      record.Material.ReviewerID,
				"reviewer":         record.Material.Reviewer,
				"reviewed_at":      record.Material.ReviewedAt,
				"rejection_reason": record.Material.RejectionReason,
				"created_at":       record.Material.CreatedAt,
				"updated_at":       record.Material.UpdatedAt,
				"downloaded_at":    record.CreatedAt.Format("2006-01-02 15:04:05"), // 添加下载时间
			}
			materials = append(materials, materialMap)
		}
	}

	response.SuccessWithPaginate(c, total, page, pageSize, materials)
}
