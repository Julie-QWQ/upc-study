package model

import (
	"time"

	"gorm.io/gorm"
)

// MaterialCategoryType 资料分类类型(存储代码)
type MaterialCategoryType string

// MaterialCategory 资料分类模型
type MaterialCategory struct {
	ID          uint             `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	DeletedAt   gorm.DeletedAt   `gorm:"index" json:"-"`
	Code        MaterialCategoryType `gorm:"type:varchar(50);not null;uniqueIndex" json:"code"`        // 分类代码
	Name        string           `gorm:"type:varchar(100);not null" json:"name"`                       // 分类名称(中文)
	Description string           `gorm:"type:text" json:"description"`                                // 分类描述
	Icon        string           `gorm:"type:varchar(100)" json:"icon"`                               // 图标
	SortOrder   int              `gorm:"not null;default:0;index" json:"sort_order"`                  // 排序
	IsActive    bool             `gorm:"not null;default:true;index" json:"is_active"`                // 是否启用
}

// TableName 指定表名
func (MaterialCategory) TableName() string {
	return "material_categories"
}

// MaterialStatus 资料状态
type MaterialStatus string

const (
	StatusPending   MaterialStatus = "pending"   // 待审核
	StatusApproved  MaterialStatus = "approved"  // 已通过
	StatusRejected  MaterialStatus = "rejected"  // 已拒绝
	StatusDeleted   MaterialStatus = "deleted"   // 已删除
)

// ReportStatus 举报状态
type ReportStatus string

const (
	ReportStatusPending  ReportStatus = "pending"  // 待处理
	ReportStatusApproved ReportStatus = "approved" // 已通过
	ReportStatusRejected ReportStatus = "rejected" // 已驳回
)

// Material 资料模型
type Material struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Title           string                `gorm:"type:varchar(200);not null;index:idx_title" json:"title"`                      // 资料标题
	Description     string                `gorm:"type:text" json:"description"`                                                // 资料描述
	Category        MaterialCategoryType  `gorm:"type:varchar(50);not null;index:idx_category" json:"category"`                // 分类代码
	CategoryInfo    *MaterialCategory     `gorm:"foreignKey:Category;references:Code" json:"category_info,omitempty"`          // 分类信息
	CourseName      string                `gorm:"type:varchar(100);index" json:"course_name"`                                   // 课程名称
	UploaderID      uint                  `gorm:"not null;index:idx_uploader" json:"uploader_id"`                               // 上传者ID
	Uploader        *User                 `gorm:"foreignKey:UploaderID" json:"uploader,omitempty"`                             // 上传者信息
	Status          MaterialStatus        `gorm:"type:varchar(20);not null;default:'pending';index:idx_status" json:"status"`  // 状态
	FileName        string                `gorm:"type:varchar(255);not null" json:"file_name"`                                  // 原始文件名
	FileSize        int64                 `gorm:"not null" json:"file_size"`                                                    // 文件大小（字节）
	FileKey         string                `gorm:"type:varchar(500);not null;uniqueIndex" json:"file_key"`                      // OSS 存储键
	MimeType        string                `gorm:"type:varchar(100);not null" json:"mime_type"`                                 // MIME 类型
	DownloadCount   int                   `gorm:"not null;default:0" json:"download_count"`                                    // 下载次数
	FavoriteCount   int                   `gorm:"not null;default:0" json:"favorite_count"`                                    // 收藏次数
	ViewCount       int                   `gorm:"not null;default:0" json:"view_count"`                                        // 浏览次数
	ReviewerID      *uint                 `gorm:"index" json:"reviewer_id,omitempty"`                                         // 审核人ID
	Reviewer        *User                 `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`                            // 审核人信息
	ReviewedAt      *time.Time            `json:"reviewed_at,omitempty"`                                                      // 审核时间
	RejectionReason string                `gorm:"type:text" json:"rejection_reason,omitempty"`                                // 拒绝原因
	SearchVector    string                `gorm:"type:tsvector;index:idx_search,gin" json:"-"`                                // 全文搜索向量
}

// TableName 指定表名
func (Material) TableName() string {
	return "materials"
}

// Favorite 收藏模型
type Favorite struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID     uint      `gorm:"not null;uniqueIndex:idx_user_material" json:"user_id"`
	User       *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	MaterialID uint      `gorm:"not null;uniqueIndex:idx_user_material;index:idx_material" json:"material_id"`
	Material   *Material `gorm:"foreignKey:MaterialID" json:"material,omitempty"`
}

// TableName 指定表名
func (Favorite) TableName() string {
	return "favorites"
}

// DownloadRecord 下载记录模型
type DownloadRecord struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID     uint      `gorm:"not null;index:idx_user_download" json:"user_id"`
	User       *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	MaterialID uint      `gorm:"not null;index:idx_material_download" json:"material_id"`
	Material   *Material `gorm:"foreignKey:MaterialID" json:"material,omitempty"`
}

// TableName 指定表名
func (DownloadRecord) TableName() string {
	return "download_records"
}

// Report 举报模型
type Report struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID       uint           `gorm:"not null;index:idx_user_report" json:"user_id"`
	User         *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	MaterialID   uint           `gorm:"not null;index:idx_material_report" json:"material_id"`
	Material     *Material      `gorm:"foreignKey:MaterialID" json:"material,omitempty"`
	Reason       string         `gorm:"type:varchar(50);not null" json:"reason"`                     // 举报原因
	Description  string         `gorm:"type:text" json:"description"`                                // 详细描述
	Status       ReportStatus   `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`   // 处理状态
	HandlerID    *uint          `gorm:"index" json:"handler_id,omitempty"`                                // 处理人ID
	Handler      *User          `gorm:"foreignKey:HandlerID" json:"handler,omitempty"`                    // 处理人信息
	HandledAt    *time.Time     `gorm:"column:handled_at" json:"handled_at,omitempty"`                   // 处理时间
	HandleNote   string         `gorm:"column:handle_comment;type:text" json:"handle_note,omitempty"`   // 处理备注
}

// TableName 指定表名
func (Report) TableName() string {
	return "reports"
}

// 举报原因常量
const (
	ReportReasonInappropriate = "inappropriate" // 内容不当
	ReportReasonCopyright     = "copyright"     // 版权问题
	ReportReasonWrong         = "wrong"         // 内容错误
	ReportReasonDuplicate     = "duplicate"     // 重复内容
	ReportReasonOther         = "other"         // 其他原因
)
