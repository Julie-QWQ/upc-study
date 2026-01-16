package model

import (
	"time"
)

// AccessLog 访问日志模型
type AccessLog struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UserID    *uint      `gorm:"index" json:"user_id,omitempty"`     // 用户ID (可为空,游客访问)
	IPAddress string     `gorm:"type:varchar(50)" json:"ip_address"`  // IP地址
	Path      string     `gorm:"type:varchar(255);not null" json:"path"`     // 访问路径
	Method    string     `gorm:"type:varchar(10);not null" json:"method"`    // 请求方法
	UserAgent string     `gorm:"type:text" json:"user_agent"`                 // 用户代理
	Referer   string     `gorm:"type:varchar(500)" json:"referer"`            // 来源页面
}

// TableName 指定表名
func (AccessLog) TableName() string {
	return "access_logs"
}

// LoginLog 登录日志模型
type LoginLog struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UserID    uint      `gorm:"index:idx_login_user_date;not null" json:"user_id"`    // 用户ID
	IPAddress string    `gorm:"type:varchar(50)" json:"ip_address"`                    // IP地址
	UserAgent string    `gorm:"type:text" json:"user_agent"`                          // 用户代理
	Success   bool      `gorm:"index;not null" json:"success"`                         // 是否登录成功
}

// TableName 指定表名
func (LoginLog) TableName() string {
	return "login_logs"
}

// SystemConfig 系统配置模型
type SystemConfig struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ConfigKey   string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"config_key"`   // 配置键
	ConfigValue string    `gorm:"type:text" json:"config_value"`                                 // 配置值
	Description string    `gorm:"type:varchar(255)" json:"description"`                          // 配置说明
	Category    string    `gorm:"type:varchar(50);default:'general'" json:"category"`            // 配置分类
}

// TableName 指定表名
func (SystemConfig) TableName() string {
	return "system_configs"
}

// OverviewStatistics 概览统计
type OverviewStatistics struct {
	Users         UserStatistics         `json:"users"`
	Materials     MaterialStatistics     `json:"materials"`
	Downloads     DownloadStatistics     `json:"downloads"`
	Applications  ApplicationStatistics  `json:"applications"`
	Visits        VisitStatistics        `json:"visits"`
}

// UserStatistics 用户统计
type UserStatistics struct {
	Total      int64 `json:"total"`
	Today      int64 `json:"today"`       // 今日新增用户
	Week       int64 `json:"week"`
	Month      int64 `json:"month"`
	Active     int64 `json:"active"`      // 今日活跃用户（登录用户数）
	ByRole     map[string]int64 `json:"by_role"`    // 按角色统计
}

// MaterialStatistics 资料统计
type MaterialStatistics struct {
	Total     int64 `json:"total"`
	Approved  int64 `json:"approved"`
	Pending   int64 `json:"pending"`
	Rejected  int64 `json:"rejected"`
	Offline   int64 `json:"offline"`
	Today     int64 `json:"today"`
	Week      int64 `json:"week"`
	ByCategory map[string]int64 `json:"by_category"` // 按分类统计
}

// DownloadStatistics 下载统计
type DownloadStatistics struct {
	Total     int64 `json:"total"`
	Today     int64 `json:"today"`
	Week      int64 `json:"week"`
	Month     int64 `json:"month"`
}

// ApplicationStatistics 学委申请统计
type ApplicationStatistics struct {
	Total    int64 `json:"total"`
	Pending  int64 `json:"pending"`
	Approved int64 `json:"approved"`
	Rejected int64 `json:"rejected"`
}

// VisitStatistics 访问统计
type VisitStatistics struct {
	Total    int64 `json:"total"`
	Today    int64 `json:"today"`
	Week     int64 `json:"week"`
	Month    int64 `json:"month"`
	Unique   int64 `json:"unique"`     // 独立访客
}

// TrendData 趋势数据 (用于图表展示)
type TrendData struct {
	Date  string `json:"date"`           // 日期
	Count int64  `json:"count"`          // 数量
	Value int64  `json:"value"`          // 值 (用于流量等)
}

// UserListRequest 用户列表请求
type UserListRequest struct {
	Page      int    `form:"page" json:"page"`
	PageSize  int    `form:"page_size" json:"page_size"`
	Keyword   string `form:"keyword" json:"keyword"`     // 搜索关键词
	Role      string `form:"role" json:"role"`           // 角色筛选
	Status    string `form:"status" json:"status"`       // 状态筛选
	Major     string `form:"major" json:"major"`         // 专业筛选
	Class     string `form:"class" json:"class"`         // 班级筛选
	SortBy    string `form:"sort_by" json:"sort_by"`     // 排序字段
	SortOrder string `form:"sort_order" json:"sort_order"` // 排序方向
}

// UserDetailResponse 用户详情响应
type UserDetailResponse struct {
	User            User             `json:"user"`
	Statistics      UserStatistics   `json:"statistics"`
	RecentActivity  []ActivityLog    `json:"recent_activity"`
	DownloadTotal   int64            `json:"download_total"`
	UploadTotal     int64            `json:"upload_total"`
	FavoriteTotal   int64            `json:"favorite_total"`
}

// ActivityLog 活动日志
type ActivityLog struct {
	Action      string    `json:"action"`       // 操作类型
	Resource    string    `json:"resource"`     // 资源类型
	Description string    `json:"description"`  // 描述
	CreatedAt   time.Time `json:"created_at"`   // 时间
}

// SystemConfigListRequest 系统配置列表请求
type SystemConfigListRequest struct {
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"page_size" json:"page_size"`
	Category string `form:"category" json:"category"` // 配置分类
	Keyword  string `form:"keyword" json:"keyword"`   // 搜索关键词
}

// UpdateSystemConfigRequest 更新系统配置请求
type UpdateSystemConfigRequest struct {
	ConfigKey   string `json:"config_key" binding:"required"`
	ConfigValue string `json:"config_value" binding:"required"`
}

// UpdateUserStatusRequest 更新用户状态请求
type UpdateUserStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active inactive banned"`
	Reason string `json:"reason"`
}

// PageViewRequest 页面浏览请求
type PageViewRequest struct {
	Path    string `json:"path" binding:"required"`    // 页面路径
	Referer string `json:"referer"`                   // 来源页面
}
