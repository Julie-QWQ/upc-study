package oss

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrInvalidFileType 无效的文件类型
	ErrInvalidFileType = errors.New("无效的文件类型")
	// ErrFileTooLarge 文件过大
	ErrFileTooLarge = errors.New("文件大小超过限制")
	// ErrInvalidFileName 无效的文件名
	ErrInvalidFileName = errors.New("无效的文件名")
)

// OSSClient OSS 客户端接口
type OSSClient interface {
	// TestConnection ?? OSS/MinIO ??
	TestConnection(ctx context.Context) error
	// GeneratePresignedUploadURL 生成预签名上传 URL
	GeneratePresignedUploadURL(ctx context.Context, fileKey string, expiresIn time.Duration) (string, error)
	// GeneratePresignedDownloadURL 生成预签名下载 URL
	GeneratePresignedDownloadURL(ctx context.Context, fileKey string, expiresIn time.Duration) (string, error)
	// DeleteFile 删除文件
	DeleteFile(ctx context.Context, fileKey string) error
	// GetFile 获取文件
	GetFile(ctx context.Context, fileKey string) ([]byte, error)
	// FileExists 检查文件是否存在
	FileExists(ctx context.Context, fileKey string) (bool, error)
}

// FileValidator 文件验证器
type FileValidator struct {
	// MaxFileSize 最大文件大小（字节）
	MaxFileSize int64
	// AllowedMimeTypes 允许的 MIME 类型
	AllowedMimeTypes map[string]bool
	// AllowedExtensions 允许的文件扩展名
	AllowedExtensions map[string]bool
}

// NewFileValidator 创建文件验证器
func NewFileValidator(maxFileSize int64) *FileValidator {
	return NewFileValidatorWithTypes(maxFileSize, getDefaultAllowedTypes())
}

// NewFileValidatorWithTypes 使用指定的文件类型创建验证器
func NewFileValidatorWithTypes(maxFileSize int64, allowedTypes []string) *FileValidator {
	allowedMimeTypes := make(map[string]bool)
	allowedExtensions := make(map[string]bool)

	// MIME 类型映射表
	mimeMap := map[string]string{
		"pdf":  "application/pdf",
		"doc":  "application/msword",
		"docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"xls":  "application/vnd.ms-excel",
		"xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"ppt":  "application/vnd.ms-powerpoint",
		"pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		"txt":  "text/plain",
		"md":   "text/markdown",
		"csv":  "text/csv",
		"zip":  "application/zip",
		"rar":  "application/x-rar-compressed",
		"7z":   "application/x-7z-compressed",
		"tar":  "application/x-tar",
		"jpg":  "image/jpeg",
		"jpeg": "image/jpeg",
		"png":  "image/png",
		"gif":  "image/gif",
		"bmp":  "image/bmp",
		"webp": "image/webp",
	}

	// 根据允许的类型构建 MIME 类型和扩展名映射
	for _, fileType := range allowedTypes {
		ext := "." + strings.ToLower(strings.TrimSpace(fileType))
		if mime, ok := mimeMap[strings.ToLower(fileType)]; ok {
			allowedMimeTypes[mime] = true
			allowedExtensions[ext] = true
		}
	}

	return &FileValidator{
		MaxFileSize:       maxFileSize,
		AllowedMimeTypes:  allowedMimeTypes,
		AllowedExtensions: allowedExtensions,
	}
}

// getDefaultAllowedTypes 获取默认允许的文件类型
func getDefaultAllowedTypes() []string {
	return []string{"pdf", "docx", "doc", "pptx", "ppt", "txt", "md", "zip", "rar"}
}

// ValidateFile 验证文件
func (v *FileValidator) ValidateFile(fileName string, fileSize int64, mimeType string) error {
	// 检查文件大小
	if fileSize > v.MaxFileSize {
		return ErrFileTooLarge
	}

	// 检查文件名
	if fileName == "" {
		return ErrInvalidFileName
	}

	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(fileName))
	if !v.AllowedExtensions[ext] {
		return ErrInvalidFileType
	}

	// 检查 MIME 类型
	if !v.AllowedMimeTypes[mimeType] {
		return ErrInvalidFileType
	}

	return nil
}

// UploadSignatureResult 上传签名结果
type UploadSignatureResult struct {
	UploadID    string    `json:"upload_id"`
	FileKey     string    `json:"file_key"`
	UploadURL   string    `json:"upload_url"`
	ExpiresAt   time.Time `json:"expires_at"`
	MaxFileSize int64     `json:"max_file_size"`
}

// OSSService OSS 服务接口
type OSSService interface {
	// GenerateUploadSignature 生成上传签名
	GenerateUploadSignature(ctx context.Context, userID uint, fileName string, fileSize int64, mimeType string) (*UploadSignatureResult, error)
	// GenerateDownloadSignature 生成下载签名
	GenerateDownloadSignature(ctx context.Context, fileKey string) (string, error)
	// DeleteFile 删除文件
	DeleteFile(ctx context.Context, fileKey string) error
	// ValidateFile 验证文件
	ValidateFile(fileName string, fileSize int64, mimeType string) error
	// UpdateConfig 更新配置
	UpdateConfig(maxFileSize int64, allowedTypes []string)
	// GetAllowedExtensions 获取允许的文件扩展名
	GetAllowedExtensions() []string
	// GetMaxFileSize 获取最大文件大小
	GetMaxFileSize() int64
}

// ossService OSS 服务实现
type ossService struct {
	client           OSSClient
	validator        *FileValidator
	uploadExpireIn   time.Duration
	downloadExpireIn time.Duration
}

// NewOSSService 创建 OSS 服务实例
func NewOSSService(client OSSClient, maxFileSize int64, uploadExpireIn, downloadExpireIn time.Duration) OSSService {
	return &ossService{
		client:           client,
		validator:        NewFileValidator(maxFileSize),
		uploadExpireIn:   uploadExpireIn,
		downloadExpireIn: downloadExpireIn,
	}
}

// UpdateConfig 更新文件验证器配置
func (s *ossService) UpdateConfig(maxFileSize int64, allowedTypes []string) {
	s.validator = NewFileValidatorWithTypes(maxFileSize, allowedTypes)
}

// GetAllowedExtensions 获取允许的文件扩展名列表
func (s *ossService) GetAllowedExtensions() []string {
	extensions := make([]string, 0, len(s.validator.AllowedExtensions))
	for ext := range s.validator.AllowedExtensions {
		// 移除点号
		extensions = append(extensions, strings.TrimPrefix(ext, "."))
	}
	return extensions
}

// GetMaxFileSize 获取最大文件大小
func (s *ossService) GetMaxFileSize() int64 {
	return s.validator.MaxFileSize
}

// GenerateUploadSignature 生成上传签名
func (s *ossService) GenerateUploadSignature(ctx context.Context, userID uint, fileName string, fileSize int64, mimeType string) (*UploadSignatureResult, error) {
	// 验证文件
	if err := s.validator.ValidateFile(fileName, fileSize, mimeType); err != nil {
		return nil, fmt.Errorf("文件验证失败: %w", err)
	}

	// 生成文件存储键
	fileKey := s.generateFileKey(userID, fileName)

	// 生成预签名上传 URL
	uploadURL, err := s.client.GeneratePresignedUploadURL(ctx, fileKey, s.uploadExpireIn)
	if err != nil {
		return nil, fmt.Errorf("生成预签名上传 URL 失败: %w", err)
	}

	// 生成上传任务 ID
	uploadID := uuid.New().String()

	return &UploadSignatureResult{
		UploadID:    uploadID,
		FileKey:     fileKey,
		UploadURL:   uploadURL,
		ExpiresAt:   time.Now().Add(s.uploadExpireIn),
		MaxFileSize: s.validator.MaxFileSize,
	}, nil
}

// GenerateDownloadSignature 生成下载签名
func (s *ossService) GenerateDownloadSignature(ctx context.Context, fileKey string) (string, error) {
	downloadURL, err := s.client.GeneratePresignedDownloadURL(ctx, fileKey, s.downloadExpireIn)
	if err != nil {
		return "", fmt.Errorf("生成预签名下载 URL 失败: %w", err)
	}
	return downloadURL, nil
}

// DeleteFile 删除文件
func (s *ossService) DeleteFile(ctx context.Context, fileKey string) error {
	return s.client.DeleteFile(ctx, fileKey)
}

// ValidateFile 验证文件
func (s *ossService) ValidateFile(fileName string, fileSize int64, mimeType string) error {
	return s.validator.ValidateFile(fileName, fileSize, mimeType)
}

// generateFileKey 生成文件存储键
func (s *ossService) generateFileKey(userID uint, fileName string) string {
	// 格式: materials/{userID}/{uuid}_{fileName}
	uid := uuid.New().String()
	return fmt.Sprintf("materials/%d/%s_%s", userID, uid, fileName)
}
