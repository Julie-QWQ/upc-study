// API 基础地址
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api'

// 本地存储键
export const STORAGE_KEYS = {
  TOKEN: 'study_upc_token',
  REFRESH_TOKEN: 'study_upc_refresh_token',
  TOKEN_EXPIRE_TIME: 'study_upc_token_expire_time',
  USER_INFO: 'study_upc_user_info',
  THEME: 'study_upc_theme'
} as const

// 用户角色
export const USER_ROLES = {
  STUDENT: 'student',
  COMMITTEE: 'committee',
  ADMIN: 'admin'
} as const

// 响应码
export const RESPONSE_CODE = {
  SUCCESS: 0,
  INVALID_PARAMS: 10001,
  UNAUTHORIZED: 10002,
  FORBIDDEN: 10003,
  NOT_FOUND: 10004,
  SERVER_ERROR: 10005,
  DUPLICATE: 10006,
  DATABASE_ERROR: 10007,
  USER_DISABLED: 10102
} as const

// 分页默认值
export const PAGINATION = {
  DEFAULT_PAGE: 1,
  DEFAULT_SIZE: 20,
  PAGE_SIZES: [10, 20, 50, 100]
} as const

// 文件上传限制
export const UPLOAD_LIMITS = {
  MAX_SIZE: 100 * 1024 * 1024, // 100MB
  ALLOWED_TYPES: [
    'application/pdf',
    'application/msword',
    'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
    'application/vnd.ms-excel',
    'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
    'application/vnd.ms-powerpoint',
    'application/vnd.openxmlformats-officedocument.presentationml.presentation',
    'application/zip',
    'application/x-rar-compressed',
    'text/plain'
  ]
} as const
