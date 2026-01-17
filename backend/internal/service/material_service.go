package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/study-upc/backend/internal/model"
	"github.com/study-upc/backend/internal/pkg/oss"
	"github.com/study-upc/backend/internal/repository"
	"gorm.io/gorm"
)

var (
	// ErrMaterialNotFound 资料不存在
	ErrMaterialNotFound = errors.New("资料不存在")
	// ErrAccessDenied 无权访问
	ErrAccessDenied = errors.New("无权访问该资源")
	// ErrMaterialAlreadyApproved 资料已审核通过
	ErrMaterialAlreadyApproved = errors.New("资料已审核通过，无法修改")
	// ErrMaterialAlreadyReviewed 资料已审核
	ErrMaterialAlreadyReviewed = errors.New("资料已审核")
	// ErrInvalidMaterialStatus 无效的资料状态
	ErrInvalidMaterialStatus = errors.New("无效的资料状态")
	// ErrDownloadLimitExceeded 超过每日下载上限
	ErrDownloadLimitExceeded = errors.New("已达到每日下载上限")
)

const (
	downloadDailyLimitKey = "download_daily_limit"
	defaultDailyLimit     = 20
)

// MaterialService 资料服务接口
type MaterialService interface {
	// CreateMaterial 创建资料
	CreateMaterial(ctx context.Context, userID uint, req *model.CreateMaterialRequest) (*model.MaterialResponse, error)
	// UpdateMaterial 更新资料
	UpdateMaterial(ctx context.Context, materialID, userID uint, userRole string, req *model.UpdateMaterialRequest) (*model.MaterialResponse, error)
	// GetMaterial 获取资料详情
	GetMaterial(ctx context.Context, materialID, currentUserID uint) (*model.MaterialResponse, error)
	// ListMaterials 获取资料列表
	ListMaterials(ctx context.Context, req *model.MaterialListRequest, currentUserID uint, currentUserRole string) (*model.MaterialListResponse, error)
	// DeleteMaterial 删除资料
	DeleteMaterial(ctx context.Context, materialID uint) error
	// ReviewMaterial 审核资料
	ReviewMaterial(ctx context.Context, materialID, reviewerID uint, req *model.ReviewMaterialRequest) error
	// GetUploadSignature 获取上传签名
	GetUploadSignature(ctx context.Context, userID uint, req *model.UploadSignatureRequest) (*model.UploadSignatureResponse, error)
	// GetDownloadURL 获取下载链接
	GetDownloadURL(ctx context.Context, materialID, userID uint) (string, error)
	// SearchMaterials 搜索资料
	SearchMaterials(ctx context.Context, keyword string, page, pageSize int) (*model.MaterialListResponse, error)
	// DeleteUploadedFile 删除已上传但未创建记录的文件
	DeleteUploadedFile(ctx context.Context, userID uint, fileKey string) error
}

// materialService 资料服务实现
type materialService struct {
	materialRepo      repository.MaterialRepository
	favoriteRepo      repository.FavoriteRepository
	downloadRepo      repository.DownloadRecordRepository
	categoryRepo      *repository.MaterialCategoryRepository
	configRepo        repository.SystemConfigRepository
	ossService        oss.OSSService
	redisClient       *redis.Client
	cacheTTL          time.Duration
}

// NewMaterialService 创建资料服务实例
func NewMaterialService(
	materialRepo repository.MaterialRepository,
	favoriteRepo repository.FavoriteRepository,
	downloadRepo repository.DownloadRecordRepository,
	categoryRepo *repository.MaterialCategoryRepository,
	configRepo repository.SystemConfigRepository,
	ossService oss.OSSService,
	redisClient *redis.Client,
) MaterialService {
	return &materialService{
		materialRepo: materialRepo,
		favoriteRepo: favoriteRepo,
		downloadRepo: downloadRepo,
		categoryRepo: categoryRepo,
		configRepo:   configRepo,
		ossService:   ossService,
		redisClient:  redisClient,
		cacheTTL:     10 * time.Minute, // 默认缓存 10 分钟
	}
}

// CreateMaterial 创建资料
func (s *materialService) CreateMaterial(ctx context.Context, userID uint, req *model.CreateMaterialRequest) (*model.MaterialResponse, error) {
	// 验证文件
	if err := s.ossService.ValidateFile(req.FileName, req.FileSize, req.MimeType); err != nil {
		return nil, fmt.Errorf("文件验证失败: %w", err)
	}

	// 验证资料类型(动态验证)
	_, err := s.categoryRepo.GetByCode(string(req.Category))
	if err != nil {
		return nil, fmt.Errorf("无效的资料类型: %s", req.Category)
	}

	// 创建资料记录
	material := &model.Material{
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		CourseName:  req.CourseName,
		UploaderID:  userID,
		Status:      model.StatusPending, // 默认待审核
		FileName:    req.FileName,
		FileSize:    req.FileSize,
		FileKey:     req.FileKey,
		MimeType:    req.MimeType,
		DownloadCount: 0,
		FavoriteCount: 0,
		ViewCount:    0,
	}

	// 生成搜索向量（全文搜索）
	material.SearchVector = s.generateSearchVector(req.Title, req.Description, req.CourseName)

	if err := s.materialRepo.Create(ctx, material); err != nil {
		return nil, fmt.Errorf("创建资料失败: %w", err)
	}

	return material.ToMaterialResponse(), nil
}

// UpdateMaterial 更新资料
func (s *materialService) UpdateMaterial(ctx context.Context, materialID, userID uint, userRole string, req *model.UpdateMaterialRequest) (*model.MaterialResponse, error) {
	// 验证资料类型(动态验证)
	_, err := s.categoryRepo.GetByCode(string(req.Category))
	if err != nil {
		return nil, fmt.Errorf("无效的资料类型: %s", req.Category)
	}

	// 获取资料
	material, err := s.materialRepo.FindByID(ctx, materialID)
	if err != nil {
		if errors.Is(err, repository.ErrMaterialNotFound) {
			return nil, ErrMaterialNotFound
		}
		return nil, fmt.Errorf("获取资料失败: %w", err)
	}

	// 检查权限
	isAdmin := userRole == "admin"
	isUploader := material.UploaderID == userID

	if !isAdmin && !isUploader {
		return nil, ErrAccessDenied
	}

	// 非管理员只能修改待审核或已拒绝的资料
	if !isAdmin && material.Status == model.StatusApproved {
		return nil, ErrMaterialAlreadyApproved
	}

	// 更新资料
	material.Title = req.Title
	material.Description = req.Description
	material.Category = req.Category
	material.CourseName = req.CourseName

	// 管理员修改不改变状态,学委修改重新提交审核
	if !isAdmin {
		material.Status = model.StatusPending
	}

	// 更新搜索向量
	material.SearchVector = s.generateSearchVector(req.Title, req.Description, req.CourseName)

	if err := s.materialRepo.Update(ctx, material); err != nil {
		return nil, fmt.Errorf("更新资料失败: %w", err)
	}

	// 清除缓存
	s.clearMaterialCache(ctx, materialID)

	return material.ToMaterialResponse(), nil
}

// GetMaterial 获取资料详情
func (s *materialService) GetMaterial(ctx context.Context, materialID, currentUserID uint) (*model.MaterialResponse, error) {
	// 尝试从缓存获取
	cacheKey := s.getMaterialCacheKey(materialID)
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		// 缓存命中，这里简化处理，实际应该序列化/反序列化
		_ = cached
	}

	// 从数据库获取
	material, err := s.materialRepo.FindByIDWithUploader(ctx, materialID)
	if err != nil {
		if errors.Is(err, repository.ErrMaterialNotFound) {
			return nil, ErrMaterialNotFound
		}
		return nil, fmt.Errorf("获取资料失败: %w", err)
	}

	// 权限检查：只有已通过审核的资料对所有用户可见
	if material.Status != model.StatusApproved {
		// 只有上传者和管理员可以查看未审核的资料
		if material.UploaderID != currentUserID {
			return nil, ErrAccessDenied
		}
	}

	// 增加浏览次数
	if err := s.materialRepo.IncrementViewCount(ctx, materialID); err != nil {
		// 记录错误但不影响主流程
		fmt.Printf("增加浏览次数失败: %v\n", err)
	}

	response := material.ToMaterialResponse()

	// 检查当前用户是否已收藏
	if currentUserID > 0 {
		favorited, _ := s.favoriteRepo.Exists(ctx, currentUserID, materialID)
		response.IsFavorited = favorited
	}

	return response, nil
}

// ListMaterials 获取资料列表
func (s *materialService) ListMaterials(ctx context.Context, req *model.MaterialListRequest, currentUserID uint, currentUserRole string) (*model.MaterialListResponse, error) {
	// 构建查询选项
	opts := &repository.MaterialListOptions{
		CourseName: req.CourseName,
		Keyword:    req.Keyword,
		SortBy:     req.SortBy,
		SortOrder:  req.SortOrder,
	}

	if req.Category != "" {
		opts.Category = &req.Category
	}

	// 设置上传者ID筛选(用于"我的资料"查询)
	if req.UploaderID != nil {
		opts.UploaderID = req.UploaderID
	}

	// 根据请求参数设置状态过滤
	if req.ReviewedOnly {
		// 管理员查询已审核资料(approved + rejected)
		opts.Statuses = []model.MaterialStatus{model.StatusApproved, model.StatusRejected}
	} else if req.Status != "" {
		// 如果明确指定了状态参数,使用指定的状态(用于"我的资料"页面筛选)
		opts.Status = &req.Status
	} else if req.UploaderID == nil {
		// 非查询"我的资料"时,默认只显示已审核通过的资料(公共展示区)
		status := model.StatusApproved
		opts.Status = &status
	}
	// 查询"我的资料"时(指定了 uploader_id)且未指定状态,显示该用户上传的所有状态资料

	// 获取资料列表
	materials, total, err := s.materialRepo.List(ctx, req.Page, req.PageSize, opts)
	if err != nil {
		return nil, fmt.Errorf("获取资料列表失败: %w", err)
	}

	// 转换为响应格式
	responses := make([]*model.MaterialResponse, 0, len(materials))
	for _, material := range materials {
		response := material.ToMaterialResponse()
		// 检查是否已收藏
		if currentUserID > 0 {
			favorited, _ := s.favoriteRepo.Exists(ctx, currentUserID, material.ID)
			response.IsFavorited = favorited
		}
		responses = append(responses, response)
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &model.MaterialListResponse{
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
		Materials:  responses,
	}, nil
}

// DeleteMaterial 删除资料
func (s *materialService) DeleteMaterial(ctx context.Context, materialID uint) error {
	// 获取资料
	material, err := s.materialRepo.FindByID(ctx, materialID)
	if err != nil {
		if errors.Is(err, repository.ErrMaterialNotFound) {
			return ErrMaterialNotFound
		}
		return fmt.Errorf("获取资料失败: %w", err)
	}

	// 删除 OSS 文件
	if err := s.ossService.DeleteFile(ctx, material.FileKey); err != nil {
		// 记录错误但继续删除数据库记录
		fmt.Printf("删除 OSS 文件失败: %v\n", err)
	}

	// 删除资料记录
	if err := s.materialRepo.Delete(ctx, materialID); err != nil {
		return fmt.Errorf("删除资料失败: %w", err)
	}

	// 清除缓存
	s.clearMaterialCache(ctx, materialID)

	return nil
}

// ReviewMaterial 审核资料
func (s *materialService) ReviewMaterial(ctx context.Context, materialID, reviewerID uint, req *model.ReviewMaterialRequest) error {
	// 获取资料
	material, err := s.materialRepo.FindByID(ctx, materialID)
	if err != nil {
		if errors.Is(err, repository.ErrMaterialNotFound) {
			return ErrMaterialNotFound
		}
		return fmt.Errorf("获取资料失败: %w", err)
	}

	// 检查状态
	if material.Status != model.StatusPending {
		return ErrMaterialAlreadyReviewed
	}

	// 如果是拒绝状态，必须填写拒绝原因
	if req.Status == model.StatusRejected && req.RejectionReason == "" {
		return errors.New("拒绝时必须填写拒绝原因")
	}

	// 更新审核状态
	var rejectionReason string
	if req.Status == model.StatusRejected {
		rejectionReason = req.RejectionReason
	}

	if err := s.materialRepo.UpdateReviewStatus(ctx, materialID, req.Status, &reviewerID, rejectionReason); err != nil {
		return fmt.Errorf("更新审核状态失败: %w", err)
	}

	// 清除缓存
	s.clearMaterialCache(ctx, materialID)

	return nil
}

// GetUploadSignature 获取上传签名
func (s *materialService) GetUploadSignature(ctx context.Context, userID uint, req *model.UploadSignatureRequest) (*model.UploadSignatureResponse, error) {
	result, err := s.ossService.GenerateUploadSignature(ctx, userID, req.FileName, req.FileSize, req.MimeType)
	if err != nil {
		return nil, fmt.Errorf("生成上传签名失败: %w", err)
	}

	return &model.UploadSignatureResponse{
		UploadID:    result.UploadID,
		FileKey:     result.FileKey,
		UploadURL:   result.UploadURL,
		ExpiresIn:   int64(result.ExpiresAt.Sub(time.Now()).Seconds()),
		MaxFileSize: result.MaxFileSize,
	}, nil
}

// GetDownloadURL 获取下载链接
func (s *materialService) GetDownloadURL(ctx context.Context, materialID, userID uint) (string, error) {
	// 获取资料
	material, err := s.materialRepo.FindByID(ctx, materialID)
	if err != nil {
		if errors.Is(err, repository.ErrMaterialNotFound) {
			return "", ErrMaterialNotFound
		}
		return "", fmt.Errorf("获取资料失败: %w", err)
	}

	// 检查状态：只能下载已审核通过的资料
	if material.Status != model.StatusApproved {
		return "", ErrAccessDenied
	}

	dailyLimit := s.getDailyDownloadLimit(ctx)
	if dailyLimit > 0 {
		now := time.Now()
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		downloaded, err := s.downloadRepo.CountByUserSince(ctx, userID, startOfDay)
		if err != nil {
			return "", fmt.Errorf("统计下载次数失败: %w", err)
		}
		if downloaded >= int64(dailyLimit) {
			return "", ErrDownloadLimitExceeded
		}
	}

	// 生成下载签名
	downloadURL, err := s.ossService.GenerateDownloadSignature(ctx, material.FileKey)
	if err != nil {
		return "", fmt.Errorf("生成下载签名失败: %w", err)
	}

	// 创建下载记录
	record := &model.DownloadRecord{
		UserID:     userID,
		MaterialID: materialID,
	}
	if err := s.downloadRepo.Create(ctx, record); err != nil {
		// 记录错误但不影响下载
		fmt.Printf("创建下载记录失败: %v\n", err)
	}

	// 增加下载次数
	if err := s.materialRepo.IncrementDownloadCount(ctx, materialID); err != nil {
		// 记录错误但不影响下载
		fmt.Printf("增加下载次数失败: %v\n", err)
	}

	// 清除缓存
	s.clearMaterialCache(ctx, materialID)

	return downloadURL, nil
}

func (s *materialService) getDailyDownloadLimit(ctx context.Context) int {
	if s.configRepo == nil {
		return defaultDailyLimit
	}

	config, err := s.configRepo.GetSystemConfig(downloadDailyLimitKey)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = s.configRepo.CreateSystemConfig(&model.SystemConfig{
				ConfigKey:   downloadDailyLimitKey,
				ConfigValue: strconv.Itoa(defaultDailyLimit),
				Description: "每个用户每日最大下载次数",
				Category:    "download",
			})
			return defaultDailyLimit
		}
		return defaultDailyLimit
	}

	limit, err := strconv.Atoi(strings.TrimSpace(config.ConfigValue))
	if err != nil || limit < 0 {
		return defaultDailyLimit
	}
	return limit
}

// SearchMaterials 搜索资料
func (s *materialService) SearchMaterials(ctx context.Context, keyword string, page, pageSize int) (*model.MaterialListResponse, error) {
	if keyword == "" {
		return nil, errors.New("搜索关键词不能为空")
	}

	// 尝试从缓存获取搜索结果
	cacheKey := s.getSearchCacheKey(keyword, page, pageSize)
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		// 缓存命中
		_ = cached
	}

	// 全文搜索
	materials, total, err := s.materialRepo.SearchByKeyword(ctx, keyword, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("搜索资料失败: %w", err)
	}

	// 转换为响应格式
	responses := make([]*model.MaterialResponse, 0, len(materials))
	for _, material := range materials {
		response := material.ToMaterialResponse()
		responses = append(responses, response)
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	result := &model.MaterialListResponse{
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		Materials:  responses,
	}

	// 缓存搜索结果
	s.cacheSearchResult(ctx, cacheKey, result)

	return result, nil
}

// generateSearchVector 生成搜索向量
func (s *materialService) generateSearchVector(title, description, courseName string) string {
	// 将标题、描述、课程名称组合成搜索向量
	parts := make([]string, 0, 3)
	if title != "" {
		parts = append(parts, title)
	}
	if description != "" {
		parts = append(parts, description)
	}
	if courseName != "" {
		parts = append(parts, courseName)
	}
	return strings.Join(parts, " ")
}

// getMaterialCacheKey 获取资料缓存键
func (s *materialService) getMaterialCacheKey(materialID uint) string {
	return fmt.Sprintf("material:%d", materialID)
}

// getSearchCacheKey 获取搜索缓存键
func (s *materialService) getSearchCacheKey(keyword string, page, pageSize int) string {
	return fmt.Sprintf("search:%s:%d:%d", keyword, page, pageSize)
}

// clearMaterialCache 清除资料缓存
func (s *materialService) clearMaterialCache(ctx context.Context, materialID uint) {
	cacheKey := s.getMaterialCacheKey(materialID)
	if err := s.redisClient.Del(ctx, cacheKey).Err(); err != nil {
		fmt.Printf("清除缓存失败: %v\n", err)
	}
}

// cacheSearchResult 缓存搜索结果
func (s *materialService) cacheSearchResult(ctx context.Context, cacheKey string, result *model.MaterialListResponse) {
	// 这里简化处理，实际应该序列化结果
	if err := s.redisClient.Set(ctx, cacheKey, result, s.cacheTTL).Err(); err != nil {
		fmt.Printf("缓存搜索结果失败: %v\n", err)
	}
}

// DeleteUploadedFile 删除已上传但未创建记录的文件
func (s *materialService) DeleteUploadedFile(ctx context.Context, userID uint, fileKey string) error {
	// 验证 fileKey 格式和权限
	// fileKey 格式应该是: materials/{userID}/{uuid}_{fileName}
	parts := strings.Split(fileKey, "/")
	if len(parts) != 3 || parts[0] != "materials" {
		return errors.New("无效的文件路径")
	}

	// 验证文件是否属于该用户
	// parts[1] 应该是用户ID的字符串形式
	if fmt.Sprintf("%d", userID) != parts[1] {
		return errors.New("无权删除该文件")
	}

	// 删除 OSS 文件
	if err := s.ossService.DeleteFile(ctx, fileKey); err != nil {
		return fmt.Errorf("删除OSS文件失败: %w", err)
	}

	return nil
}
