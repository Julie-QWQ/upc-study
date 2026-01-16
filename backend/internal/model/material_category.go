package model

import (
	"time"

	"gorm.io/gorm"
)

// MaterialCategoryConfig 资料类型配置
type MaterialCategoryConfig struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Code        string  `gorm:"type:varchar(50);not null;uniqueIndex:idx_code" json:"code"`                    // 类型代码
	Name        string  `gorm:"type:varchar(100);not null" json:"name"`                                        // 类型名称
	Description string  `gorm:"type:text" json:"description"`                                                   // 描述
	Icon        string  `gorm:"type:varchar(100)" json:"icon"`                                                   // 图标
	SortOrder   int     `gorm:"not null;default:0;index:idx_sort_order" json:"sort_order"`                      // 排序
	IsActive    bool    `gorm:"not null;default:true;index:idx_is_active" json:"is_active"`                     // 是否启用
	NameZh      string  `gorm:"-" json:"name_zh"`                                                                // 中文名称
	NameEn      string  `gorm:"-" json:"name_en"`                                                                // 英文名称
}

// TableName 指定表名
func (MaterialCategoryConfig) TableName() string {
	return "material_categories"
}

// MaterialCategoryRequest 创建/更新资料类型请求
type MaterialCategoryRequest struct {
	Code        string `json:"code" binding:"required,max=50"`                // 类型代码
	Name        string `json:"name" binding:"required,max=100"`              // 类型名称
	Description string `json:"description"`                                  // 描述
	Icon        string `json:"icon"`                                         // 图标
	SortOrder   int    `json:"sort_order"`                                   // 排序
	IsActive    *bool  `json:"is_active"`                                    // 是否启用
}
