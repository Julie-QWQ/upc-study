package model

import (
	"time"

	"gorm.io/gorm"
)

// UserRole 用户角色类型
type UserRole string

const (
	RoleStudent   UserRole = "student"   // 学生
	RoleCommittee UserRole = "committee" // 学委
	RoleAdmin     UserRole = "admin"     // 管理员
)

// UserStatus 用户状态
type UserStatus string

const (
	StatusActive UserStatus = "active" // 正常
	StatusBanned UserStatus = "banned" // 封禁
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Username    string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`       // 用户名
	Email       string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`         // 邮箱
	PasswordHash string     `gorm:"column:password_hash;type:varchar(255);not null" json:"-"`  // 密码哈希（不返回给前端）
	RealName    string     `gorm:"type:varchar(50)" json:"real_name"`                            // 真实姓名
	Role     UserRole   `gorm:"type:varchar(20);not null;default:'student'" json:"role"`      // 角色
	Status   UserStatus `gorm:"type:varchar(20);not null;default:'inactive'" json:"status"`   // 状态
	BanReason string     `gorm:"type:text" json:"ban_reason,omitempty"`                         // ????
	Avatar   string     `gorm:"type:varchar(255)" json:"avatar"`                               // 头像URL
	Phone    string     `gorm:"type:varchar(20)" json:"phone"`                                 // 联系电话
	Major    string     `gorm:"type:varchar(100)" json:"major"`                                // 专业
	Class    string     `gorm:"type:varchar(50)" json:"class"`                                 // 班级
	LastLoginAt *time.Time `json:"last_login_at"`                                              // 最后登录时间
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
