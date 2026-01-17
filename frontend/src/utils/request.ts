import axios, { type AxiosInstance, type AxiosRequestConfig, type AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'
import { API_BASE_URL, RESPONSE_CODE, STORAGE_KEYS } from './constants'
import { storage } from './storage'
import type { ApiResponse } from '@/types'

// Mock 支持 (仅在开发环境)
const USE_MOCK = import.meta.env.USE_MOCK === 'true'
let mockHandler: typeof import('@/mock/handler') | null = null

if (USE_MOCK && import.meta.env.DEV) {
  import('@/mock/handler').then(module => {
    mockHandler = module
    console.log('[Mock] Mock 数据已启用')
  })
}

/**
 * 创建 axios 实例
 */
const service: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json'
  },
  paramsSerializer: {
    serialize: (params) => {
      // 过滤掉 undefined 和 null 值,防止查询参数污染
      const filteredParams = Object.keys(params).reduce((acc, key) => {
        if (params[key] !== undefined && params[key] !== null) {
          acc[key] = params[key]
        }
        return acc
      }, {} as Record<string, any>)
      return new URLSearchParams(filteredParams).toString()
    }
  }
})

/**
 * 请求拦截器
 */
service.interceptors.request.use(
  config => {
    // 添加 Token
    const token = storage.getToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }

    return config
  },
  error => {
    return Promise.reject(error)
  }
)

/**
 * 响应拦截器
 */
service.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    const { code, message, data } = response.data

    // 成功
    if (code === RESPONSE_CODE.SUCCESS) {
      return response.data
    }

    // 业务错误
    if (code !== RESPONSE_CODE.USER_DISABLED) {
      ElMessage.error(message || '请求失败')
    }
    const error = new Error(message || '请求失败') as Error & { code?: number }
    error.code = code
    return Promise.reject(error)
  },
  error => {
    // HTTP 错误
    if (error.response) {
      const { status } = error.response

      switch (status) {
        case 401:
          ElMessage.error('登录已过期，请重新登录')
          storage.clearAuth()
          window.location.href = '/login'
          break
        case 403:
          ElMessage.error('没有权限访问')
          break
        case 404:
          ElMessage.error('请求的资源不存在')
          break
        case 500:
          ElMessage.error('服务器内部错误')
          break
        default:
          ElMessage.error(error.response.data?.message || '请求失败')
      }
    } else if (error.request) {
      ElMessage.error('网络错误，请检查网络连接')
    } else {
      ElMessage.error('请求配置错误')
    }

    return Promise.reject(error)
  }
)

/**
 * 封装请求方法
 */
export const request = {
  /**
   * GET 请求
   */
  async get<T = any>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    // Mock 处理
    if (USE_MOCK && mockHandler) {
      const response = await mockHandler.handleMockRequest('GET', url, config?.params)
      return response as ApiResponse<T>
    }
    return service.get(url, config)
  },

  /**
   * POST 请求
   */
  async post<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    // Mock 处理
    if (USE_MOCK && mockHandler) {
      const response = await mockHandler.handleMockRequest('POST', url, data)
      return response as ApiResponse<T>
    }
    return service.post(url, data, config)
  },

  /**
   * PUT 请求
   */
  async put<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    if (USE_MOCK && mockHandler) {
      const response = await mockHandler.handleMockRequest('PUT', url, data)
      return response as ApiResponse<T>
    }
    return service.put(url, data, config)
  },

  /**
   * DELETE 请求
   */
  async delete<T = any>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    if (USE_MOCK && mockHandler) {
      const response = await mockHandler.handleMockRequest('DELETE', url, config?.params)
      return response as ApiResponse<T>
    }
    return service.delete(url, config)
  },

  /**
   * PATCH 请求
   */
  async patch<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    if (USE_MOCK && mockHandler) {
      const response = await mockHandler.handleMockRequest('PATCH', url, data)
      return response as ApiResponse<T>
    }
    return service.patch(url, data, config)
  },

  /**
   * 文件上传
   */
  async upload<T = any>(url: string, formData: FormData, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    if (USE_MOCK && mockHandler) {
      const response = await mockHandler.handleMockRequest('POST', url, formData)
      return response as ApiResponse<T>
    }
    return service.post(url, formData, {
      ...config,
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  }
}

export default service
