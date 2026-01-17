package repository

import (
	"github.com/study-upc/backend/internal/model"
	"gorm.io/gorm"
)

// AdminRepository 管理员数据访问接口
type AdminRepository interface {
	// 系统配置管理
	GetSystemConfig(key string) (*model.SystemConfig, error)
	ListSystemConfigs(req *model.SystemConfigListRequest) ([]model.SystemConfig, int64, error)
	UpdateSystemConfig(key, value string) error
	CreateSystemConfig(config *model.SystemConfig) error
	DeleteSystemConfig(key string) error

	// 用户管理
	ListUsers(req *model.UserListRequest) ([]model.User, int64, error)
	GetUserByID(id uint) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	UpdateUser(user *model.User) error
	UpdateUserStatus(id uint, status, reason string) error
	DeleteUser(id uint) error
	CountUserDownloads(userID uint) (int64, error)
	CountUserUploads(userID uint) (int64, error)
	CountUserFavorites(userID uint) (int64, error)
}

type adminRepository struct {
	db *gorm.DB
}

// NewAdminRepository 创建管理员仓库
func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db: db}
}

// ============ 系统配置管理 ============

// GetSystemConfig 获取单个系统配置
func (r *adminRepository) GetSystemConfig(key string) (*model.SystemConfig, error) {
	var config model.SystemConfig
	err := r.db.Where("config_key = ?", key).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// ListSystemConfigs 获取系统配置列表
func (r *adminRepository) ListSystemConfigs(req *model.SystemConfigListRequest) ([]model.SystemConfig, int64, error) {
	var configs []model.SystemConfig
	var total int64

	query := r.db.Model(&model.SystemConfig{})

	// 分类筛选
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}

	// 关键词搜索
	if req.Keyword != "" {
		query = query.Where("config_key LIKE ? OR description LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("category, config_key").Find(&configs).Error; err != nil {
		return nil, 0, err
	}

	return configs, total, nil
}

// UpdateSystemConfig 更新系统配置
func (r *adminRepository) UpdateSystemConfig(key, value string) error {
	return r.db.Model(&model.SystemConfig{}).
		Where("config_key = ?", key).
		Update("config_value", value).Error
}

// CreateSystemConfig 创建系统配置
func (r *adminRepository) CreateSystemConfig(config *model.SystemConfig) error {
	return r.db.Create(config).Error
}

// DeleteSystemConfig 删除系统配置
func (r *adminRepository) DeleteSystemConfig(key string) error {
	return r.db.Where("config_key = ?", key).Delete(&model.SystemConfig{}).Error
}

// ============ 用户管理 ============

// ListUsers 获取用户列表
func (r *adminRepository) ListUsers(req *model.UserListRequest) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.Model(&model.User{})

	// 关键词搜索 (用户名、真实姓名、邮箱)
	if req.Keyword != "" {
		query = query.Where("username LIKE ? OR real_name LIKE ? OR email LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 角色筛选
	if req.Role != "" {
		query = query.Where("role = ?", req.Role)
	}

	// 状态筛选
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// 专业筛选
	if req.Major != "" {
		query = query.Where("major LIKE ?", "%"+req.Major+"%")
	}

	// 班级筛选
	if req.Class != "" {
		query = query.Where("class LIKE ?", "%"+req.Class+"%")
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	orderBy := "created_at DESC"
	if req.SortBy != "" {
		direction := "DESC"
		if req.SortOrder == "asc" {
			direction = "ASC"
		}
		orderBy = req.SortBy + " " + direction
	}

	// 分页
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order(orderBy).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetUserByID 根据ID获取用户
func (r *adminRepository) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func (r *adminRepository) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func (r *adminRepository) UpdateUser(user *model.User) error {
	return r.db.Save(user).Error
}

// UpdateUserStatus 更新用户状态
func (r *adminRepository) UpdateUserStatus(id uint, status, reason string) error {
	updates := map[string]any{
		"status": status,
	}
	if reason != "" {
		updates["ban_reason"] = reason
	} else {
		updates["ban_reason"] = nil
	}

	return r.db.Model(&model.User{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// DeleteUser 删除用户 (软删除)
func (r *adminRepository) DeleteUser(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

func (r *adminRepository) CountUserDownloads(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.DownloadRecord{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *adminRepository) CountUserUploads(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Material{}).Where("uploader_id = ?", userID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *adminRepository) CountUserFavorites(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Favorite{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
