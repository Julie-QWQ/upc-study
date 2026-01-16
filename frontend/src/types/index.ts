// 通用响应结构
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

// 分页数据结构
export interface PaginateData<T = any> {
  total: number
  page: number
  size: number
  list: T[]
}

// 分页查询参数
export interface PaginationParams {
  page?: number
  size?: number
}

// 用户角色
export type UserRole = 'student' | 'committee' | 'admin'

// 用户状态
export type UserStatus = 'active' | 'banned'

// 用户信息
export interface UserInfo {
  id: number
  username: string
  email: string
  real_name: string
  role: UserRole
  status: UserStatus
  avatar: string
  phone: string
  major: string
  class: string
  created_at: string
}

// 登录请求
export interface LoginRequest {
  username: string
  password: string
}

// 登录响应
export interface LoginResponse {
  access_token: string
  refresh_token: string
  expires_in: number
  user: UserInfo
}

// 注册请求
export interface RegisterRequest {
  username: string
  email: string
  password: string
  real_name: string
  major: string
  class: string
}

// 修改密码请求
export interface ChangePasswordRequest {
  old_password: string
  new_password: string
}

// 刷新 Token 请求
export interface RefreshTokenRequest {
  refresh_token: string
}

// 路由元信息
export interface RouteMeta {
  title?: string
  requiresAuth?: boolean
  roles?: UserRole[]
}

// ==================== 资料相关类型 ====================

// 资料分类（动态配置，不再使用硬编码）
export type MaterialCategory = string

// 资料分类配置
export interface MaterialCategoryConfig {
  id: number
  code: string
  name: string
  name_zh?: string
  name_en?: string
  description?: string
  icon?: string
  sort_order: number
  is_active: boolean
  created_at: string
  updated_at: string
}

// 资料分类配置请求
export interface MaterialCategoryConfigRequest {
  code: string
  name: string
  description?: string
  icon?: string
  sort_order?: number
  is_active?: boolean
}

// 资料状态
export type MaterialStatus = 'pending' | 'approved' | 'rejected' | 'deleted'

// 资料信息
export interface Material {
  id: number
  title: string
  description: string
  category: MaterialCategory
  status: MaterialStatus
  course_name: string
  file_name: string
  file_size: number
  file_key: string
  mime_type: string
  download_count: number
  favorite_count: number
  view_count: number
  uploader_id: number
  tags: string[]
  reviewer_id?: number
  reviewed_at?: string
  rejection_reason?: string
  created_at: string
  updated_at: string
  uploader?: UserInfo
  is_favorited?: boolean
}

// 资料列表查询参数
export interface MaterialListParams extends PaginationParams {
  category?: MaterialCategory
  status?: MaterialStatus
  course_name?: string
  keyword?: string
  sort_by?: 'created_at' | 'updated_at' | 'download_count' | 'view_count' | 'favorite_count'
  order?: 'asc' | 'desc'
  uploader_id?: number // 上传者ID筛选,用于"我的资料"查询
  reviewed_only?: boolean // 管理员专用:只获取已审核资料(approved + rejected)
}

// 创建资料请求
export interface CreateMaterialRequest {
  title: string
  description: string
  category: MaterialCategory
  course_name: string
  file_name: string
  file_size: number
  file_key: string
  mime_type: string
  tags?: string[]
}

// 更新资料请求
export interface UpdateMaterialRequest {
  title?: string
  description?: string
  category?: MaterialCategory
  course_name?: string
  tags?: string[]
}

// 上传签名响应
export interface UploadSignature {
  upload_url: string
  file_key: string
  expires_at: string
}

// 上传签名请求
export interface UploadSignatureRequest {
  file_name: string
  file_size: number
  mime_type: string
}

// 下载签名响应
export interface DownloadSignature {
  download_url: string
  expires_at: string
}

// ==================== 收藏相关类型 ====================

// 收藏信息
export interface Favorite {
  id: number
  user_id: number
  material_id: number
  material?: Material
  created_at: string
}

// 收藏列表参数
export interface FavoriteListParams extends PaginationParams {
  category?: MaterialCategory
  course_name?: string
}

// ==================== 下载记录相关类型 ====================

// 下载记录
export interface DownloadRecord {
  id: number
  user_id: number
  material_id: number
  material?: Material
  created_at: string
}

// 下载记录列表参数
export interface DownloadRecordListParams extends PaginationParams {
  category?: MaterialCategory
  course_name?: string
}

// ==================== 举报相关类型 ====================

// 举报原因
export type ReportReason =
  | 'inappropriate'   // 内容不当
  | 'copyright'       // 版权问题
  | 'wrong_category'  // 分类错误
  | 'low_quality'     // 质量差
  | 'other'           // 其他

// 举报状态
export type ReportStatus = 'pending' | 'approved' | 'rejected'

// 举报信息
export interface Report {
  id: number
  material_id: number
  user_id: number
  reason: ReportReason
  description: string
  status: ReportStatus
  handler_id?: number
  handled_at?: string
  handle_note?: string
  created_at: string
  updated_at: string
  material?: Material
  reporter?: UserInfo
  handler?: UserInfo
}

// 举报请求
export interface CreateReportRequest {
  reason: ReportReason
  description: string
}

// 举报列表参数
export interface ReportListParams extends PaginationParams {
  status?: ReportStatus
  reason?: ReportReason
}

// 处理举报请求
export interface HandleReportRequest {
  approved: boolean
  note: string
}

// ==================== 学委申请相关类型 ====================

// 申请状态
export type ApplicationStatus = 'pending' | 'approved' | 'rejected' | 'cancelled'

// 学委申请信息
export interface CommitteeApplication {
  id: number
  user_id: number
  status: ApplicationStatus
  reason: string
  reviewer_id?: number
  review_comment?: string
  reviewed_at?: string
  created_at: string
  updated_at: string
  user?: UserInfo
  reviewer?: UserInfo
}

// 学委申请请求
export interface CreateCommitteeApplicationRequest {
  reason: string
}

// 学委申请列表参数
export interface CommitteeApplicationListParams extends PaginationParams {
  status?: ApplicationStatus
}

// 审核学委申请请求
export interface ReviewCommitteeApplicationRequest {
  approved: boolean
  comment?: string
}

// ==================== 审核记录相关类型 ====================

// 审核目标类型
export type ReviewTarget = 'material' | 'committee' | 'report'

// 审核动作
export type ReviewAction = 'approved' | 'rejected'

// 审核记录
export interface ReviewRecord {
  id: number
  reviewer_id: number
  target_type: ReviewTarget
  target_id: number
  action: ReviewAction
  comment?: string
  original_data: Record<string, any>
  created_at: string
  reviewer?: UserInfo
}

// 审核历史查询参数
export interface ReviewHistoryParams extends PaginationParams {
  target_type?: ReviewTarget
  target_id?: number
  reviewer_id?: number
  action?: ReviewAction
}

// ==================== 通知相关类型 ====================

// 通知类型
export type NotificationType =
  | 'system'       // 系统通知
  | 'material'     // 资料审核通知
  | 'committee'    // 学委申请通知
  | 'report'       // 举报处理通知

// 通知状态
export type NotificationStatus = 'unread' | 'read'

// 通知信息
export interface Notification {
  id: number
  user_id: number
  type: NotificationType
  title: string
  content: string
  status: NotificationStatus
  link?: string
  read_at?: string
  created_at: string
  updated_at: string
}

// 通知列表参数
export interface NotificationListParams extends PaginationParams {
  status?: NotificationStatus
  type?: NotificationType
}

// 标记已读请求
export interface MarkAsReadRequest {
  notification_ids?: number[]
}

// ==================== 公告相关 ====================

// 公告优先级
export type AnnouncementPriority = 'normal' | 'high'

// 公告信息
export interface Announcement {
  id: number
  title: string
  content: string
  priority: AnnouncementPriority
  author?: UserInfo
  is_active: boolean
  published_at?: string
  expires_at?: string
  created_at: string
  updated_at: string
}

// 公告列表参数
export interface AnnouncementListParams extends PaginationParams {
  priority?: AnnouncementPriority
  is_active?: boolean
  author_id?: number
}

// 创建公告请求
export interface CreateAnnouncementRequest {
  title: string
  content: string
  priority: AnnouncementPriority
  is_active: boolean
  published_at?: string
  expires_at?: string
}

// 更新公告请求
export interface UpdateAnnouncementRequest {
  title: string
  content: string
  priority: AnnouncementPriority
  is_active: boolean
  published_at?: string
  expires_at?: string
}
