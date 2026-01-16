package repository

import "github.com/study-upc/backend/internal/model"

// SystemConfigRepository 提供系统配置的读取能力
type SystemConfigRepository interface {
	GetSystemConfig(key string) (*model.SystemConfig, error)
	CreateSystemConfig(config *model.SystemConfig) error
}
