package service

import (
	"errors"

	"github.com/study-upc/backend/internal/model"
	"github.com/study-upc/backend/internal/repository"
	"gorm.io/gorm"
)

var (
	// ErrLastAdmin 不能禁用/删除最后一个管理员
	ErrLastAdmin = errors.New("不能禁用/删除最后一个管理员")
	// ErrCannotDeleteAdmin 不能删除管理员
	ErrCannotDeleteAdmin = errors.New("不能删除管理员账户")
)

const (
	downloadDailyLimitKeyConfig = "download_daily_limit"
	defaultDailyLimitValue      = "20"
)

// AdminService 管理员服务接口
type AdminService interface {
	// 系统配置管理
	GetSystemConfig(key string) (*model.SystemConfig, error)
	ListSystemConfigs(req *model.SystemConfigListRequest) ([]model.SystemConfig, int64, error)
	UpdateSystemConfig(key, value string) error
	CreateSystemConfig(config *model.SystemConfig) error
	DeleteSystemConfig(key string) error

	// 用户管理
	ListUsers(req *model.UserListRequest) ([]model.User, int64, error)
	GetUserDetail(id uint) (*model.UserDetailResponse, error)
	UpdateUserStatus(id uint, status, reason string) error
	UpdateUserInfo(id uint, updates map[string]interface{}) error
	DeleteUser(id uint) error
}

type adminService struct {
	adminRepo  repository.AdminRepository
	userRepo   repository.UserRepository
	materialRepo repository.MaterialRepository
}

// NewAdminService 创建管理员服务
func NewAdminService(
	adminRepo repository.AdminRepository,
	userRepo repository.UserRepository,
	materialRepo repository.MaterialRepository,
) AdminService {
	return &adminService{
		adminRepo:  adminRepo,
		userRepo:   userRepo,
		materialRepo: materialRepo,
	}
}

// ============ 系统配置管理 ============

// GetSystemConfig 获取单个系统配置
func (s *adminService) GetSystemConfig(key string) (*model.SystemConfig, error) {
	return s.adminRepo.GetSystemConfig(key)
}

// ListSystemConfigs 获取系统配置列表
func (s *adminService) ListSystemConfigs(req *model.SystemConfigListRequest) ([]model.SystemConfig, int64, error) {
	s.ensureDownloadDailyLimitConfig()
	return s.adminRepo.ListSystemConfigs(req)
}

// UpdateSystemConfig 更新系统配置
func (s *adminService) UpdateSystemConfig(key, value string) error {
	// 检查配置是否存在
	config, err := s.adminRepo.GetSystemConfig(key)
	if err != nil {
		return err
	}

	// 更新配置值
	return s.adminRepo.UpdateSystemConfig(config.ConfigKey, value)
}

// CreateSystemConfig 创建系统配置
func (s *adminService) CreateSystemConfig(config *model.SystemConfig) error {
	return s.adminRepo.CreateSystemConfig(config)
}

// DeleteSystemConfig 删除系统配置
func (s *adminService) DeleteSystemConfig(key string) error {
	return s.adminRepo.DeleteSystemConfig(key)
}

func (s *adminService) ensureDownloadDailyLimitConfig() {
	_, err := s.adminRepo.GetSystemConfig(downloadDailyLimitKeyConfig)
	if err == nil {
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}
	_ = s.adminRepo.CreateSystemConfig(&model.SystemConfig{
		ConfigKey:   downloadDailyLimitKeyConfig,
		ConfigValue: defaultDailyLimitValue,
		Description: "每个用户每日最大下载次数",
		Category:    "download",
	})
}

// ============ 用户管理 ============

// ListUsers 获取用户列表
func (s *adminService) ListUsers(req *model.UserListRequest) ([]model.User, int64, error) {
	return s.adminRepo.ListUsers(req)
}

// GetUserDetail 获取用户详情
func (s *adminService) GetUserDetail(id uint) (*model.UserDetailResponse, error) {
	// 获取用户信息
	user, err := s.adminRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	response := &model.UserDetailResponse{
		User: *user,
	}

	downloadTotal, err := s.adminRepo.CountUserDownloads(id)
	if err != nil {
		return nil, err
	}
	uploadTotal, err := s.adminRepo.CountUserUploads(id)
	if err != nil {
		return nil, err
	}
	favoriteTotal, err := s.adminRepo.CountUserFavorites(id)
	if err != nil {
		return nil, err
	}

	response.DownloadTotal = downloadTotal
	response.UploadTotal = uploadTotal
	response.FavoriteTotal = favoriteTotal

	return response, nil
}

// UpdateUserStatus 更新用户状态
func (s *adminService) UpdateUserStatus(id uint, status, reason string) error {
	// 检查用户是否存在
	user, err := s.adminRepo.GetUserByID(id)
	if err != nil {
		return err
	}

	// 不能禁用最后一个管理员
	if user.Role == model.RoleAdmin && status == string(model.StatusBanned) {
		// 检查是否还有其他管理员
		adminCount, err := s.userRepo.CountByRole(model.RoleAdmin)
		if err != nil {
			return err
		}
		if adminCount <= 1 {
			return ErrLastAdmin
		}
	}

	// 更新用户状态
	return s.adminRepo.UpdateUserStatus(id, status)
}

// UpdateUserInfo 更新用户信息
func (s *adminService) UpdateUserInfo(id uint, updates map[string]interface{}) error {
	// 获取用户
	user, err := s.adminRepo.GetUserByID(id)
	if err != nil {
		return err
	}

	// 更新字段
	if realName, ok := updates["real_name"].(string); ok {
		user.RealName = realName
	}
	if email, ok := updates["email"].(string); ok {
		user.Email = email
	}
	if phone, ok := updates["phone"].(string); ok {
		user.Phone = phone
	}
	if major, ok := updates["major"].(string); ok {
		user.Major = major
	}
	if class, ok := updates["class"].(string); ok {
		user.Class = class
	}
	if avatar, ok := updates["avatar"].(string); ok {
		user.Avatar = avatar
	}

	return s.adminRepo.UpdateUser(user)
}

// DeleteUser 删除用户
func (s *adminService) DeleteUser(id uint) error {
	// 检查用户是否存在
	user, err := s.adminRepo.GetUserByID(id)
	if err != nil {
		return err
	}

	// 不能删除管理员
	if user.Role == model.RoleAdmin {
		return ErrCannotDeleteAdmin
	}

	// 软删除用户
	return s.adminRepo.DeleteUser(id)
}
