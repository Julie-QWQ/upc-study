import request from '@/utils/request'

// 用户管理相关类型

export interface User {
  id: number
  username: string
  email: string
  real_name: string
  role: string
  status: string
  avatar: string
  phone: string
  major: string
  class: string
  last_login_at: string
  created_at: string
}

export interface UserListRequest {
  page?: number
  page_size?: number
  keyword?: string
  role?: string
  status?: string
  major?: string
  class?: string
  sort_by?: string
  sort_order?: string
}

export interface UserListResponse {
  items: User[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}

export interface UpdateUserStatusRequest {
  status: string
  reason?: string
}

export interface UserDetailResponse {
  user: User
  statistics?: UserStatistics
  recent_activity?: ActivityLog[]
  download_total?: number
  upload_total?: number
  favorite_total?: number
}

export interface ActivityLog {
  action: string
  resource: string
  description: string
  created_at: string
}

export interface UserStatistics {
  total?: number
  today?: number
  week?: number
  month?: number
  active?: number
  by_role?: Record<string, number>
}

// 系统配置相关类型

export interface SystemConfig {
  id: number
  config_key: string
  config_value: string
  description: string
  category: string
  created_at: string
  updated_at: string
}

export interface SystemConfigListRequest {
  page?: number
  page_size?: number
  category?: string
  keyword?: string
}

export interface UpdateSystemConfigRequest {
  config_key: string
  config_value: string
}

/**
 * 获取用户列表
 */
export function getUserList(params: UserListRequest) {
  return request<UserListResponse>({
    url: '/admin/users',
    method: 'get',
    params
  })
}

/**
 * 获取用户详情
 */
export function getUserDetail(id: number) {
  return request<UserDetailResponse>({
    url: `/admin/users/${id}`,
    method: 'get'
  })
}

/**
 * 更新用户信息
 */
export function updateUserInfo(id: number, data: Partial<User>) {
  return request({
    url: `/admin/users/${id}`,
    method: 'put',
    data
  })
}

/**
 * 更新用户状态
 */
export function updateUserStatus(id: number, data: UpdateUserStatusRequest) {
  return request({
    url: `/admin/users/${id}/status`,
    method: 'put',
    data
  })
}

/**
 * 删除用户
 */
export function deleteUser(id: number) {
  return request({
    url: `/admin/users/${id}`,
    method: 'delete'
  })
}

/**
 * 获取系统配置列表
 */
export function getSystemConfigList(params: SystemConfigListRequest) {
  return request<{
    list: SystemConfig[]
    total: number
    page: number
    page_size: number
  }>({
    url: '/admin/configs',
    method: 'get',
    params
  })
}

/**
 * 获取单个系统配置
 */
export function getSystemConfig(key: string) {
  return request<SystemConfig>({
    url: `/admin/configs/${key}`,
    method: 'get'
  })
}

/**
 * 创建系统配置
 */
export function createSystemConfig(data: Omit<SystemConfig, 'id' | 'created_at' | 'updated_at'>) {
  return request({
    url: '/admin/configs',
    method: 'post',
    data
  })
}

/**
 * 更新系统配置
 */
export function updateSystemConfig(data: UpdateSystemConfigRequest) {
  return request({
    url: '/admin/configs',
    method: 'put',
    data
  })
}

/**
 * 删除系统配置
 */
export function deleteSystemConfig(key: string) {
  return request({
    url: `/admin/configs/${key}`,
    method: 'delete'
  })
}
