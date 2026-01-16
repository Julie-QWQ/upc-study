package repository

import (
	"context"
	"errors"
	"time"

	"github.com/study-upc/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	// ErrMaterialNotFound 资料不存在错误
	ErrMaterialNotFound = errors.New("资料不存在")
	// ErrMaterialAlreadyExists 资料已存在错误
	ErrMaterialAlreadyExists = errors.New("资料已存在")
	// ErrFavoriteAlreadyExists 收藏已存在错误
	ErrFavoriteAlreadyExists = errors.New("已收藏该资料")
	// ErrFavoriteNotFound 收藏不存在错误
	ErrFavoriteNotFound = errors.New("未收藏该资料")
	// ErrReportAlreadyExists 举报已存在错误
	ErrReportAlreadyExists = errors.New("已举报该资料")
	// ErrReportNotFound 举报不存在错误
	ErrReportNotFound = errors.New("举报不存在")
)

// MaterialListOptions 资料列表查询选项
type MaterialListOptions struct {
	Category      *model.MaterialCategory
	CourseName    string
	Status        *model.MaterialStatus
	Statuses      []model.MaterialStatus // 支持多个状态查询(用于管理员查询"已审核"资料)
	Keyword       string
	SortBy        string // created_at, download_count, favorite_count, view_count, title
	SortOrder     string // asc, desc
	UploaderID    *uint  // 上传者ID筛选,用于"我的资料"查询
}

// MaterialRepository 资料数据访问层接口
type MaterialRepository interface {
	// Create 创建资料
	Create(ctx context.Context, material *model.Material) error
	// FindByID 根据ID查找资料
	FindByID(ctx context.Context, id uint) (*model.Material, error)
	// FindByIDWithUploader 根据ID查找资料（包含上传者信息）
	FindByIDWithUploader(ctx context.Context, id uint) (*model.Material, error)
	// FindByIDForUpdate 根据ID查找资料（加锁）
	FindByIDForUpdate(ctx context.Context, id uint) (*model.Material, error)
	// Update 更新资料
	Update(ctx context.Context, material *model.Material) error
	// Delete 删除资料（软删除）
	Delete(ctx context.Context, id uint) error
	// List 分页获取资料列表
	List(ctx context.Context, page, pageSize int, opts *MaterialListOptions) ([]*model.Material, int64, error)
	// IncrementViewCount 增加浏览次数
	IncrementViewCount(ctx context.Context, id uint) error
	// IncrementDownloadCount 增加下载次数
	IncrementDownloadCount(ctx context.Context, id uint) error
	// IncrementFavoriteCount 增加收藏次数
	IncrementFavoriteCount(ctx context.Context, id uint) error
	// DecrementFavoriteCount 减少收藏次数
	DecrementFavoriteCount(ctx context.Context, id uint) error
	// UpdateReviewStatus 更新审核状态
	UpdateReviewStatus(ctx context.Context, id uint, status model.MaterialStatus, reviewerID *uint, rejectionReason string) error
	// FindByFileKey 根据文件存储键查找资料
	FindByFileKey(ctx context.Context, fileKey string) (*model.Material, error)
	// SearchByKeyword 全文搜索
	SearchByKeyword(ctx context.Context, keyword string, page, pageSize int) ([]*model.Material, int64, error)
}

// FavoriteRepository 收藏数据访问层接口
type FavoriteRepository interface {
	// Create 创建收藏
	Create(ctx context.Context, favorite *model.Favorite) error
	// Delete 删除收藏
	Delete(ctx context.Context, userID, materialID uint) error
	// FindByUserAndMaterial 查找用户对指定资料的收藏记录
	FindByUserAndMaterial(ctx context.Context, userID, materialID uint) (*model.Favorite, error)
	// Exists 检查是否已收藏
	Exists(ctx context.Context, userID, materialID uint) (bool, error)
	// ListByUser 分页获取用户的收藏列表
	ListByUser(ctx context.Context, userID uint, page, pageSize int) ([]*model.Favorite, int64, error)
	// CountByMaterial 统计资料的收藏数
	CountByMaterial(ctx context.Context, materialID uint) (int64, error)
}

// DownloadRecordRepository 下载记录数据访问层接口
type DownloadRecordRepository interface {
	// Create 创建下载记录
	Create(ctx context.Context, record *model.DownloadRecord) error
	// FindByUserAndMaterial 查找用户对指定资料的下载记录
	FindByUserAndMaterial(ctx context.Context, userID, materialID uint) (*model.DownloadRecord, error)
	// ListByUser 分页获取用户的下载记录
	ListByUser(ctx context.Context, userID uint, page, pageSize int) ([]*model.DownloadRecord, int64, error)
	// CountByMaterial 统计资料的下载次数
	CountByMaterial(ctx context.Context, materialID uint) (int64, error)
	// ListRecentByMaterial 获取资料最近的下载记录
	ListRecentByMaterial(ctx context.Context, materialID uint, limit int) ([]*model.DownloadRecord, error)
	// CountByUserSince 统计用户从指定时间起的下载次数
	CountByUserSince(ctx context.Context, userID uint, since time.Time) (int64, error)
}

// ReportRepository 举报数据访问层接口
type ReportRepository interface {
	// Create 创建举报
	Create(ctx context.Context, report *model.Report) error
	// FindByID 根据ID查找举报
	FindByID(ctx context.Context, id uint) (*model.Report, error)
	// Update 更新举报
	Update(ctx context.Context, report *model.Report) error
	// List 分页获取举报列表
	List(ctx context.Context, page, pageSize int, status *model.ReportStatus) ([]*model.Report, int64, error)
	// ListByMaterial 获取资料的举报列表
	ListByMaterial(ctx context.Context, materialID uint) ([]*model.Report, error)
	// FindPendingByUserAndMaterial 查找用户对指定资料的待处理举报
	FindPendingByUserAndMaterial(ctx context.Context, userID, materialID uint) (*model.Report, error)
	// CountPending 统计待处理举报数
	CountPending(ctx context.Context) (int64, error)
}

// materialRepository 资料数据访问层实现
type materialRepository struct {
	db *gorm.DB
}

// NewMaterialRepository 创建资料数据访问层实例
func NewMaterialRepository(db *gorm.DB) MaterialRepository {
	return &materialRepository{db: db}
}

// Create 创建资料
func (r *materialRepository) Create(ctx context.Context, material *model.Material) error {
	result := r.db.WithContext(ctx).Create(material)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrMaterialAlreadyExists
		}
		return result.Error
	}
	return nil
}

// FindByID 根据ID查找资料
func (r *materialRepository) FindByID(ctx context.Context, id uint) (*model.Material, error) {
	var material model.Material
	result := r.db.WithContext(ctx).First(&material, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrMaterialNotFound
		}
		return nil, result.Error
	}
	return &material, nil
}

// FindByIDWithUploader 根据ID查找资料（包含上传者信息）
func (r *materialRepository) FindByIDWithUploader(ctx context.Context, id uint) (*model.Material, error) {
	var material model.Material
	result := r.db.WithContext(ctx).Preload("Uploader").Preload("Reviewer").First(&material, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrMaterialNotFound
		}
		return nil, result.Error
	}
	return &material, nil
}

// FindByIDForUpdate 根据ID查找资料（加锁）
func (r *materialRepository) FindByIDForUpdate(ctx context.Context, id uint) (*model.Material, error) {
	var material model.Material
	result := r.db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).First(&material, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrMaterialNotFound
		}
		return nil, result.Error
	}
	return &material, nil
}

// Update 更新资料
func (r *materialRepository) Update(ctx context.Context, material *model.Material) error {
	result := r.db.WithContext(ctx).Save(material)
	return result.Error
}

// Delete 删除资料（软删除）
func (r *materialRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&model.Material{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrMaterialNotFound
	}
	return nil
}

// List 分页获取资料列表
func (r *materialRepository) List(ctx context.Context, page, pageSize int, opts *MaterialListOptions) ([]*model.Material, int64, error) {
	var materials []*model.Material
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Material{})

	// 应用筛选条件
	if opts != nil {
		if opts.Category != nil {
			query = query.Where("category = ?", *opts.Category)
		}
		if opts.CourseName != "" {
			query = query.Where("course_name LIKE ?", "%"+opts.CourseName+"%")
		}
		// 支持单状态和多状态查询
		if opts.Status != nil {
			query = query.Where("status = ?", *opts.Status)
		} else if len(opts.Statuses) > 0 {
			query = query.Where("status IN ?", opts.Statuses)
		}
		if opts.Keyword != "" {
			// 模糊全文搜索 - 使用 plainto_tsquery 支持分词和模糊匹配
			query = query.Where("search_vector @@ plainto_tsquery('simple', ?)", opts.Keyword)
		}
		if opts.UploaderID != nil {
			query = query.Where("uploader_id = ?", *opts.UploaderID)
		}
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 应用排序
	orderBy := "created_at DESC"
	if opts != nil && opts.SortBy != "" {
		order := "DESC"
		if opts.SortOrder == "asc" {
			order = "ASC"
		}
		orderBy = opts.SortBy + " " + order
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	result := query.Order(orderBy).Offset(offset).Limit(pageSize).Preload("Uploader").Preload("Reviewer").Find(&materials)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return materials, total, nil
}

// IncrementViewCount 增加浏览次数
func (r *materialRepository) IncrementViewCount(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Model(&model.Material{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + 1"))
	return result.Error
}

// IncrementDownloadCount 增加下载次数
func (r *materialRepository) IncrementDownloadCount(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Model(&model.Material{}).Where("id = ?", id).UpdateColumn("download_count", gorm.Expr("download_count + 1"))
	return result.Error
}

// IncrementFavoriteCount 增加收藏次数
func (r *materialRepository) IncrementFavoriteCount(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Model(&model.Material{}).Where("id = ?", id).UpdateColumn("favorite_count", gorm.Expr("favorite_count + 1"))
	return result.Error
}

// DecrementFavoriteCount 减少收藏次数
func (r *materialRepository) DecrementFavoriteCount(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Model(&model.Material{}).Where("id = ?", id).UpdateColumn("favorite_count", gorm.Expr("favorite_count - 1"))
	return result.Error
}

// UpdateReviewStatus 更新审核状态
func (r *materialRepository) UpdateReviewStatus(ctx context.Context, id uint, status model.MaterialStatus, reviewerID *uint, rejectionReason string) error {
	updates := map[string]interface{}{
		"status":     status,
		"reviewed_at": gorm.Expr("NOW()"),
	}
	if reviewerID != nil {
		updates["reviewer_id"] = *reviewerID
	}
	if rejectionReason != "" {
		updates["rejection_reason"] = rejectionReason
	}

	result := r.db.WithContext(ctx).Model(&model.Material{}).Where("id = ?", id).Updates(updates)
	return result.Error
}

// FindByFileKey 根据文件存储键查找资料
func (r *materialRepository) FindByFileKey(ctx context.Context, fileKey string) (*model.Material, error) {
	var material model.Material
	result := r.db.WithContext(ctx).Where("file_key = ?", fileKey).First(&material)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrMaterialNotFound
		}
		return nil, result.Error
	}
	return &material, nil
}

// SearchByKeyword 全文搜索（支持模糊匹配，只搜索已审核通过的资料）
func (r *materialRepository) SearchByKeyword(ctx context.Context, keyword string, page, pageSize int) ([]*model.Material, int64, error) {
	var materials []*model.Material
	var total int64

	// 只搜索已审核通过的资料
	query := r.db.WithContext(ctx).Model(&model.Material{}).Where(
		"search_vector @@ plainto_tsquery('simple', ?) AND status = ?",
		keyword,
		model.StatusApproved,
	)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	result := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Preload("Uploader").Preload("Reviewer").Find(&materials)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return materials, total, nil
}

// favoriteRepository 收藏数据访问层实现
type favoriteRepository struct {
	db *gorm.DB
}

// NewFavoriteRepository 创建收藏数据访问层实例
func NewFavoriteRepository(db *gorm.DB) FavoriteRepository {
	return &favoriteRepository{db: db}
}

// Create 创建收藏
func (r *favoriteRepository) Create(ctx context.Context, favorite *model.Favorite) error {
	// 先检查是否存在已软删除的记录
	var existing model.Favorite
	result := r.db.WithContext(ctx).Unscoped().Where("user_id = ? AND material_id = ?", favorite.UserID, favorite.MaterialID).First(&existing)

	if result.Error == nil {
		// 找到了记录（可能是软删除的）
		if existing.DeletedAt.Valid {
			// 如果是软删除的记录，恢复它
			result = r.db.WithContext(ctx).Unscoped().Model(&existing).Updates(map[string]interface{}{
				"deleted_at": nil,
				"created_at": favorite.CreatedAt,
				"updated_at": favorite.UpdatedAt,
			})
			if result.Error != nil {
				return result.Error
			}
			// 更新传入的 favorite 的 ID
			favorite.ID = existing.ID
			return nil
		} else {
			// 如果未删除，说明已经收藏了
			return ErrFavoriteAlreadyExists
		}
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 其他错误
		return result.Error
	}

	// 没有找到任何记录，创建新的
	result = r.db.WithContext(ctx).Create(favorite)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrFavoriteAlreadyExists
		}
		return result.Error
	}
	return nil
}

// Delete 删除收藏
func (r *favoriteRepository) Delete(ctx context.Context, userID, materialID uint) error {
	result := r.db.WithContext(ctx).Where("user_id = ? AND material_id = ?", userID, materialID).Delete(&model.Favorite{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrFavoriteNotFound
	}
	return nil
}

// FindByUserAndMaterial 查找用户对指定资料的收藏记录
func (r *favoriteRepository) FindByUserAndMaterial(ctx context.Context, userID, materialID uint) (*model.Favorite, error) {
	var favorite model.Favorite
	result := r.db.WithContext(ctx).Where("user_id = ? AND material_id = ?", userID, materialID).First(&favorite)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrFavoriteNotFound
		}
		return nil, result.Error
	}
	return &favorite, nil
}

// Exists 检查是否已收藏（排除已软删除的记录）
func (r *favoriteRepository) Exists(ctx context.Context, userID, materialID uint) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Unscoped().Model(&model.Favorite{}).Where("user_id = ? AND material_id = ? AND deleted_at IS NULL", userID, materialID).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

// ListByUser 分页获取用户的收藏列表
func (r *favoriteRepository) ListByUser(ctx context.Context, userID uint, page, pageSize int) ([]*model.Favorite, int64, error) {
	var favorites []*model.Favorite
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Favorite{}).Where("user_id = ?", userID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	result := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Preload("Material.Uploader").Preload("Material.Reviewer").Find(&favorites)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return favorites, total, nil
}

// CountByMaterial 统计资料的收藏数
func (r *favoriteRepository) CountByMaterial(ctx context.Context, materialID uint) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&model.Favorite{}).Where("material_id = ?", materialID).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

// downloadRecordRepository 下载记录数据访问层实现
type downloadRecordRepository struct {
	db *gorm.DB
}

// NewDownloadRecordRepository 创建下载记录数据访问层实例
func NewDownloadRecordRepository(db *gorm.DB) DownloadRecordRepository {
	return &downloadRecordRepository{db: db}
}

// Create 创建下载记录
func (r *downloadRecordRepository) Create(ctx context.Context, record *model.DownloadRecord) error {
	result := r.db.WithContext(ctx).Create(record)
	return result.Error
}

// FindByUserAndMaterial 查找用户对指定资料的下载记录
func (r *downloadRecordRepository) FindByUserAndMaterial(ctx context.Context, userID, materialID uint) (*model.DownloadRecord, error) {
	var record model.DownloadRecord
	result := r.db.WithContext(ctx).Where("user_id = ? AND material_id = ?", userID, materialID).First(&record)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}
	return &record, nil
}

// ListByUser 分页获取用户的下载记录
func (r *downloadRecordRepository) ListByUser(ctx context.Context, userID uint, page, pageSize int) ([]*model.DownloadRecord, int64, error) {
	var records []*model.DownloadRecord
	var total int64

	query := r.db.WithContext(ctx).Model(&model.DownloadRecord{}).Where("user_id = ?", userID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	result := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Preload("Material.Uploader").Find(&records)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return records, total, nil
}

// CountByMaterial 统计资料的下载次数
func (r *downloadRecordRepository) CountByMaterial(ctx context.Context, materialID uint) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&model.DownloadRecord{}).Where("material_id = ?", materialID).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

// ListRecentByMaterial 获取资料最近的下载记录
func (r *downloadRecordRepository) ListRecentByMaterial(ctx context.Context, materialID uint, limit int) ([]*model.DownloadRecord, error) {
	var records []*model.DownloadRecord
	result := r.db.WithContext(ctx).Where("material_id = ?", materialID).Order("created_at DESC").Limit(limit).Preload("User").Find(&records)
	if result.Error != nil {
		return nil, result.Error
	}
	return records, nil
}

func (r *downloadRecordRepository) CountByUserSince(ctx context.Context, userID uint, since time.Time) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&model.DownloadRecord{}).
		Where("user_id = ? AND created_at >= ?", userID, since).
		Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

// reportRepository 举报数据访问层实现
type reportRepository struct {
	db *gorm.DB
}

// NewReportRepository 创建举报数据访问层实例
func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

// Create 创建举报
func (r *reportRepository) Create(ctx context.Context, report *model.Report) error {
	result := r.db.WithContext(ctx).Create(report)
	return result.Error
}

// FindByID 根据ID查找举报
func (r *reportRepository) FindByID(ctx context.Context, id uint) (*model.Report, error) {
	var report model.Report
	result := r.db.WithContext(ctx).
		Preload("User").
		Preload("Material", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped() // 包含已删除的资料
		}).
		Preload("Material.Uploader").
		Preload("Handler").
		First(&report, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrReportNotFound
		}
		return nil, result.Error
	}
	return &report, nil
}

// Update 更新举报
func (r *reportRepository) Update(ctx context.Context, report *model.Report) error {
	result := r.db.WithContext(ctx).Save(report)
	return result.Error
}

// List 分页获取举报列表
func (r *reportRepository) List(ctx context.Context, page, pageSize int, status *model.ReportStatus) ([]*model.Report, int64, error) {
	var reports []*model.Report
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Report{})

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	result := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Preload("User").
		Preload("Material", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped() // 包含已删除的资料
		}).
		Preload("Material.Uploader").
		Preload("Handler").
		Find(&reports)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return reports, total, nil
}

// ListByMaterial 获取资料的举报列表
func (r *reportRepository) ListByMaterial(ctx context.Context, materialID uint) ([]*model.Report, error) {
	var reports []*model.Report
	result := r.db.WithContext(ctx).Where("material_id = ?", materialID).Preload("User").Find(&reports)
	if result.Error != nil {
		return nil, result.Error
	}
	return reports, nil
}

// FindPendingByUserAndMaterial 查找用户对指定资料的待处理举报
func (r *reportRepository) FindPendingByUserAndMaterial(ctx context.Context, userID, materialID uint) (*model.Report, error) {
	var report model.Report
	result := r.db.WithContext(ctx).Where("user_id = ? AND material_id = ? AND status = ?", userID, materialID, model.ReportStatusPending).First(&report)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}
	return &report, nil
}

// CountPending 统计待处理举报数
func (r *reportRepository) CountPending(ctx context.Context) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&model.Report{}).Where("status = ?", model.ReportStatusPending).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}
