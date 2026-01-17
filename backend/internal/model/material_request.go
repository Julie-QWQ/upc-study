package model

import (
	"fmt"

	"github.com/google/uuid"
)

// CreateMaterialRequest 创建资料请求
type CreateMaterialRequest struct {
	Title       string           `json:"title" binding:"required,min=2,max=200"`
	Description string           `json:"description" binding:"max=2000"`
	Category    MaterialCategoryType `json:"category" binding:"required"` // 移除 oneof,改为动态验证
	CourseName  string           `json:"course_name" binding:"required,max=100"`
	FileName    string           `json:"file_name" binding:"required,max=255"`
	FileSize    int64            `json:"file_size" binding:"required,min=1,max=536870912"` // 最大 512MB
	MimeType    string           `json:"mime_type" binding:"required"`
	FileKey     string           `json:"file_key" binding:"required,max=500"` // OSS 存储键
}

// UpdateMaterialRequest 更新资料请求
type UpdateMaterialRequest struct {
	Title       string           `json:"title" binding:"required,min=2,max=200"`
	Description string           `json:"description" binding:"max=2000"`
	Category    MaterialCategoryType `json:"category" binding:"required"` // 移除 oneof,改为动态验证
	CourseName  string           `json:"course_name" binding:"required,max=100"`
}

// MaterialListRequest 资料列表查询请求
type MaterialListRequest struct {
	Page         int              `form:"page,default=1" binding:"min=1"`
	PageSize     int              `form:"page_size,default=20" binding:"min=1,max=100"`
	Category     MaterialCategoryType `form:"category" binding:"omitempty"` // 移除 oneof,改为动态验证
	CourseName   string           `form:"course_name" binding:"omitempty,max=100"`
	Status       MaterialStatus   `form:"status" binding:"omitempty,oneof=pending approved rejected deleted"`
	Keyword      string           `form:"keyword" binding:"omitempty,max=100"`
	SortBy       string           `form:"sort_by,default=created_at" binding:"omitempty,oneof=created_at download_count favorite_count view_count title"`
	SortOrder    string           `form:"sort_order,default=desc" binding:"omitempty,oneof=asc desc"`
	UploaderID   *uint            `form:"uploader_id" binding:"omitempty"` // 可选的上传者ID筛选,用于"我的资料"查询
	ReviewedOnly bool             `form:"reviewed_only"`                  // 管理员专用:只获取已审核资料(approved + rejected)
}

// ReviewMaterialRequest 审核资料请求
type ReviewMaterialRequest struct {
	Status          MaterialStatus `json:"status" binding:"required,oneof=approved rejected"`
	RejectionReason string         `json:"rejection_reason" binding:"omitempty,max=500"`
}

// MaterialResponse 资料响应
type MaterialResponse struct {
	ID              uint             `json:"id"`
	Title           string           `json:"title"`
	Description     string           `json:"description"`
	Category        MaterialCategoryType `json:"category"`
	CourseName      string           `json:"course_name"`
	UploaderID      uint             `json:"uploader_id"`
	Uploader        *UserInfo        `json:"uploader,omitempty"`
	Status          MaterialStatus   `json:"status"`
	FileName        string           `json:"file_name"`
	FileSize        int64            `json:"file_size"`
	MimeType        string           `json:"mime_type"`
	DownloadCount   int              `json:"download_count"`
	FavoriteCount   int              `json:"favorite_count"`
	ViewCount       int              `json:"view_count"`
	ReviewerID      *uint            `json:"reviewer_id,omitempty"`
	Reviewer        *UserInfo        `json:"reviewer,omitempty"`
	ReviewedAt      *string          `json:"reviewed_at,omitempty"`
	RejectionReason string           `json:"rejection_reason,omitempty"`
	CreatedAt       string           `json:"created_at"`
	UpdatedAt       string           `json:"updated_at"`
	IsFavorited     bool             `json:"is_favorited,omitempty"` // 当前用户是否已收藏
}

// UploadSignatureRequest 获取上传签名请求
type UploadSignatureRequest struct {
	FileName string `json:"file_name" binding:"required,max=255"`
	FileSize int64  `json:"file_size" binding:"required,min=1,max=536870912"` // 最大 512MB
	MimeType string `json:"mime_type" binding:"required"`
}

// UploadSignatureResponse 上传签名响应
type UploadSignatureResponse struct {
	UploadID    string `json:"upload_id"`              // 上传任务ID（UUID）
	FileKey     string `json:"file_key"`               // OSS 存储键
	UploadURL   string `json:"upload_url"`             // 上传URL（预签名URL）
	ExpiresIn   int64  `json:"expires_in"`             // 过期时间（秒）
	MaxFileSize int64  `json:"max_file_size"`          // 最大文件大小
}

// DeleteUploadedFileRequest 删除已上传文件请求
type DeleteUploadedFileRequest struct {
	FileKey string `json:"file_key" binding:"required"` // OSS 存储键
}

// ReportRequest 举报请求
type ReportRequest struct {
	Reason      string `json:"reason" binding:"required,oneof=inappropriate copyright wrong duplicate other"`
	Description string `json:"description" binding:"max=500"`
}

// HandleReportRequest 处理举报请求
type HandleReportRequest struct {
	Status     ReportStatus `json:"status" binding:"required,oneof=approved rejected"`
	HandleNote string       `json:"handle_note" binding:"omitempty,max=500"`
}

// ReportResponse 举报响应
type ReportResponse struct {
	ID          uint         `json:"id"`
	UserID      uint         `json:"user_id"`
	Reporter    *UserInfo    `json:"reporter,omitempty"` // 举报人信息(前端期望reporter字段)
	MaterialID  uint         `json:"material_id"`
	Material    *MaterialResponse `json:"material,omitempty"`
	Reason      string       `json:"reason"`
	Description string       `json:"description"`
	Status      ReportStatus `json:"status"`
	HandlerID   *uint        `json:"handler_id,omitempty"`
	Handler     *UserInfo    `json:"handler,omitempty"`
	HandledAt   *string      `json:"handled_at,omitempty"`
	HandleNote  string       `json:"handle_note,omitempty"`
	CreatedAt   string       `json:"created_at"`
}

// MaterialListResponse 资料列表响应
type MaterialListResponse struct {
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
	Materials  []*MaterialResponse `json:"materials"`
}

// FavoriteResponse 收藏响应
type FavoriteResponse struct {
	ID         uint             `json:"id"`
	UserID     uint             `json:"user_id"`
	MaterialID uint             `json:"material_id"`
	Material   *MaterialResponse `json:"material,omitempty"`
	CreatedAt  string           `json:"created_at"`
}

// DownloadRecordResponse 下载记录响应
type DownloadRecordResponse struct {
	ID         uint             `json:"id"`
	UserID     uint             `json:"user_id"`
	User       *UserInfo        `json:"user,omitempty"`
	MaterialID uint             `json:"material_id"`
	Material   *MaterialResponse `json:"material,omitempty"`
	CreatedAt  string           `json:"created_at"`
}

// ToMaterialResponse 将 Material 转换为 MaterialResponse
func (m *Material) ToMaterialResponse() *MaterialResponse {
	response := &MaterialResponse{
		ID:              m.ID,
		Title:           m.Title,
		Description:     m.Description,
		Category:        m.Category,
		CourseName:      m.CourseName,
		UploaderID:      m.UploaderID,
		Status:          m.Status,
		FileName:        m.FileName,
		FileSize:        m.FileSize,
		MimeType:        m.MimeType,
		DownloadCount:   m.DownloadCount,
		FavoriteCount:   m.FavoriteCount,
		ViewCount:       m.ViewCount,
		ReviewerID:      m.ReviewerID,
		RejectionReason: m.RejectionReason,
		CreatedAt:       m.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       m.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	// 转换上传者信息
	if m.Uploader != nil {
		uploaderInfo := m.Uploader.ToUserInfo()
		response.Uploader = &uploaderInfo
	}

	// 转换审核人信息
	if m.Reviewer != nil {
		reviewerInfo := m.Reviewer.ToUserInfo()
		response.Reviewer = &reviewerInfo
	}

	// 转换审核时间
	if m.ReviewedAt != nil {
		reviewedAt := m.ReviewedAt.Format("2006-01-02 15:04:05")
		response.ReviewedAt = &reviewedAt
	}

	return response
}

// GenerateFileKey 生成 OSS 文件存储键
func GenerateFileKey(userID uint, fileName string) string {
	// 格式: materials/{userID}/{uuid}_{fileName}
	uid := uuid.New().String()
	return fmt.Sprintf("materials/%d/%s_%s", userID, uid, fileName)
}
