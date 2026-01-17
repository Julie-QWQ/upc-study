package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/study-upc/backend/internal/model"
	"github.com/study-upc/backend/internal/pkg/utils"
	"github.com/study-upc/backend/internal/repository"

	"github.com/redis/go-redis/v9"
)

var (
	// ErrInvalidCredentials 无效的登录凭证
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	// ErrUserDisabled 用户已被禁用
	ErrUserDisabled = errors.New("用户已被禁用")
	// ErrUserInactive 用户未激活
	ErrUserInactive = errors.New("用户未激活")
	// ErrWrongPassword 旧密码错误
	ErrWrongPassword = errors.New("旧密码错误")
	// ErrTokenInBlacklist Token 已在黑名单中
	ErrTokenInBlacklist = errors.New("Token 已失效")
	// ErrTokenExpired Token 已过期
	ErrTokenExpired = utils.ErrTokenExpired
	// ErrTokenInvalid Token 无效
	ErrTokenInvalid = utils.ErrTokenInvalid
)

// AuthService 认证服务接口
type AuthService interface {
	// Register 用户注册
	Register(ctx context.Context, req *model.RegisterRequest) (*model.UserInfo, error)
	// Login 用户登录
	Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error)
	// Logout 用户登出
	Logout(ctx context.Context, accessToken string) error
	// RefreshToken 刷新 Token
	RefreshToken(ctx context.Context, refreshToken string) (*model.LoginResponse, error)
	// ChangePassword 修改密码
	ChangePassword(ctx context.Context, userID uint, req *model.ChangePasswordRequest) error
	// GetUserInfo 获取用户信息
	GetUserInfo(ctx context.Context, userID uint) (*model.UserInfo, error)
}

// authService 认证服务实现
type authService struct {
	userRepo      repository.UserRepository
	jwtManager    *utils.JWTManager
	redisClient   *redis.Client
	tokenBlacklistPrefix string
}

// NewAuthService 创建认证服务实例
func NewAuthService(
	userRepo repository.UserRepository,
	jwtManager *utils.JWTManager,
	redisClient *redis.Client,
) AuthService {
	return &authService{
		userRepo:             userRepo,
		jwtManager:           jwtManager,
		redisClient:          redisClient,
		tokenBlacklistPrefix: "auth:blacklist:",
	}
}

// Register 用户注册
func (s *authService) Register(ctx context.Context, req *model.RegisterRequest) (*model.UserInfo, error) {
	// 检查用户名是否已存在
	exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("检查用户名失败: %w", err)
	}
	if exists {
		return nil, repository.ErrUserAlreadyExists
	}

	// 检查邮箱是否已存在
	exists, err = s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("检查邮箱失败: %w", err)
	}
	if exists {
		return nil, repository.ErrUserAlreadyExists
	}

	// 密码加密
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	// 创建用户
	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		RealName:     req.RealName,
		Role:         model.RoleStudent, // 默认角色为学生
		Status:       model.StatusActive, // 默认状态为激活
		Major:        req.Major,
		Class:        req.Class,
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	userInfo := user.ToUserInfo()
	return &userInfo, nil
}

// Login 用户登录
func (s *authService) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	// 查找用户（支持用户名或邮箱登录）
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			// 尝试用邮箱登录
			user, err = s.userRepo.FindByEmail(ctx, req.Username)
			if err != nil {
				return nil, ErrInvalidCredentials
			}
		} else {
			return nil, fmt.Errorf("查找用户失败: %w", err)
		}
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	// 检查用户状态
	if user.Status == model.StatusBanned {
		reason := strings.TrimSpace(user.BanReason)
		if reason != "" {
			return nil, fmt.Errorf("%w: %s", ErrUserDisabled, reason)
		}
		return nil, ErrUserDisabled
	}

	// 生成 Token
	accessToken, refreshToken, err := s.jwtManager.GenerateTokenPair(user.ID, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("生成 Token 失败: %w", err)
	}

	// 更新最后登录时间
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		// 记录错误但不影响登录流程
		fmt.Printf("更新最后登录时间失败: %v\n", err)
	}

	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.jwtManager.GetAccessTTL(),
		User:         user.ToUserInfo(),
	}, nil
}

// Logout 用户登出
func (s *authService) Logout(ctx context.Context, accessToken string) error {
	// 解析 Token 获取过期时间
	claims, err := s.jwtManager.ParseToken(accessToken)
	if err != nil {
		return fmt.Errorf("解析 Token 失败: %w", err)
	}

	// 计算 Token 剩余有效期
	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl <= 0 {
		// Token 已过期，无需加入黑名单
		return nil
	}

	// 将 Token 加入黑名单
	key := s.tokenBlacklistPrefix + accessToken
	if err := s.redisClient.Set(ctx, key, "1", ttl).Err(); err != nil {
		return fmt.Errorf("加入 Token 黑名单失败: %w", err)
	}

	return nil
}

// RefreshToken 刷新 Token
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*model.LoginResponse, error) {
	// 验证刷新 Token
	claims, err := s.jwtManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("无效的刷新 Token: %w", err)
	}

	// 检查 Token 是否在黑名单中
	key := s.tokenBlacklistPrefix + refreshToken
	exists, err := s.redisClient.Exists(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("检查 Token 黑名单失败: %w", err)
	}
	if exists > 0 {
		return nil, ErrTokenInBlacklist
	}

	// 获取用户信息
	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 检查用户状态
	if user.Status == model.StatusBanned {
		reason := strings.TrimSpace(user.BanReason)
		if reason != "" {
			return nil, fmt.Errorf("%w: %s", ErrUserDisabled, reason)
		}
		return nil, ErrUserDisabled
	}

	// 生成新的 Token 对
	accessToken, newRefreshToken, err := s.jwtManager.GenerateTokenPair(user.ID, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("生成 Token 失败: %w", err)
	}

	// 将旧的刷新 Token 加入黑名单
	oldTTL := time.Until(claims.ExpiresAt.Time)
	if oldTTL > 0 {
		oldKey := s.tokenBlacklistPrefix + refreshToken
		if err := s.redisClient.Set(ctx, oldKey, "1", oldTTL).Err(); err != nil {
			fmt.Printf("将旧刷新 Token 加入黑名单失败: %v\n", err)
		}
	}

	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    s.jwtManager.GetAccessTTL(),
		User:         user.ToUserInfo(),
	}, nil
}

// ChangePassword 修改密码
func (s *authService) ChangePassword(ctx context.Context, userID uint, req *model.ChangePasswordRequest) error {
	// 获取用户信息
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 验证旧密码
	if !utils.CheckPassword(req.OldPassword, user.PasswordHash) {
		return ErrWrongPassword
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 更新密码
	if err := s.userRepo.UpdatePassword(ctx, userID, hashedPassword); err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	return nil
}

// GetUserInfo 获取用户信息
func (s *authService) GetUserInfo(ctx context.Context, userID uint) (*model.UserInfo, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	userInfo := user.ToUserInfo()
	return &userInfo, nil
}
