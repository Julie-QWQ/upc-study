import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/api/auth'
import { storage } from '@/utils/storage'
import type { UserInfo, LoginRequest, RegisterRequest, ChangePasswordRequest } from '@/types'
import { ElMessage, ElMessageBox } from 'element-plus'
import { RESPONSE_CODE } from '@/utils/constants'

export const useAuthStore = defineStore('auth', () => {
  // 状态
  const user = ref<UserInfo | null>(null)
  const accessToken = ref<string | null>(null)
  const refreshToken = ref<string | null>(null)
  const tokenExpireTime = ref<number | null>(null)

  // 计算属性
  const isLoggedIn = computed(() => !!user.value && !!accessToken.value)
  const userRole = computed(() => user.value?.role || null)
  const userName = computed(() => user.value?.real_name || user.value?.username || '')
  const userAvatar = computed(() => user.value?.avatar || '')

  // 初始化：从 localStorage 恢复状态
  const initAuth = () => {
    const savedUser = storage.getUser()
    const savedToken = storage.getToken()
    const savedRefreshToken = storage.getRefreshToken()
    const savedExpireTime = storage.getTokenExpireTime()

    if (savedUser && savedToken) {
      user.value = savedUser
      accessToken.value = savedToken
      refreshToken.value = savedRefreshToken
      tokenExpireTime.value = savedExpireTime

      // 检查 token 是否即将过期（剩余时间小于 5 分钟）
      if (savedExpireTime && Date.now() > savedExpireTime - 5 * 60 * 1000) {
        refreshAccessToken().catch(() => {
          // 刷新失败，清除认证状态
          clearAuth()
        })
      }
    }
  }

  // 保存认证状态到 localStorage
  const saveAuth = (userInfo: UserInfo, token: string, refresh: string, expiresIn: number) => {
    user.value = userInfo
    accessToken.value = token
    refreshToken.value = refresh
    tokenExpireTime.value = Date.now() + expiresIn * 1000

    storage.setUser(userInfo)
    storage.setToken(token)
    storage.setRefreshToken(refresh)
    storage.setTokenExpireTime(tokenExpireTime.value)
  }

  // 清除认证状态
  const clearAuth = () => {
    user.value = null
    accessToken.value = null
    refreshToken.value = null
    tokenExpireTime.value = null

    storage.clearAuth()
  }

  // 登录
  const login = async (credentials: LoginRequest) => {
    try {
      const response = await authApi.login(credentials)
      const { user: userInfo, access_token, refresh_token, expires_in } = response.data

      saveAuth(userInfo, access_token, refresh_token, expires_in)

      ElMessage.success('登录成功')

      // 异步记录访问日志（登录成功时记录一次访问）
      fetch('/api/v1/statistics/page-view', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${access_token}`
        },
        body: JSON.stringify({
          path: '/login',
          referer: undefined
        }),
        keepalive: true
      }).catch((error) => {
        console.error('[登录访问记录失败]', error)
      })

      return true
    } catch (error: any) {
      if (error?.code === RESPONSE_CODE.USER_DISABLED) {
        const message = error?.message || '该账号已被封禁，请联系管理员。'
        await ElMessageBox.alert(message, '无法登录', {
          type: 'error',
          confirmButtonText: '知道了'
        })
      } else {
        ElMessage.error(error.message || '登录失败')
      }
      return false
    }
  }

  // 注册
  const register = async (data: RegisterRequest) => {
    try {
      const response = await authApi.register(data)
      // 注册成功，不自动登录
      ElMessage.success('注册成功，请登录')
      return true
    } catch (error: any) {
      ElMessage.error(error.message || '注册失败')
      return false
    }
  }

  // 登出
  const logout = async () => {
    // 立即清除本地认证状态
    clearAuth()
    ElMessage.success('已退出登录')

    // 异步调用 logout API（不等待）
    authApi.logout().catch(error => {
      console.error('Logout API failed:', error)
    })
  }

  // 刷新访问令牌
  const refreshAccessToken = async () => {
    if (!refreshToken.value) {
      throw new Error('没有刷新令牌')
    }

    try {
      const response = await authApi.refreshToken({
        refresh_token: refreshToken.value
      })
      const { user: userInfo, access_token, refresh_token, expires_in } = response.data

      saveAuth(userInfo, access_token, refresh_token, expires_in)
      return true
    } catch (error) {
      clearAuth()
      throw error
    }
  }

  // 修改密码
  const changePassword = async (data: ChangePasswordRequest) => {
    try {
      await authApi.changePassword(data)
      ElMessage.success('密码修改成功，请重新登录')
      await logout()
      return true
    } catch (error: any) {
      ElMessage.error(error.message || '密码修改失败')
      return false
    }
  }

  // 获取用户信息
  const fetchUserInfo = async () => {
    try {
      const response = await authApi.getUserInfo()
      user.value = response.data
      storage.setUser(response.data)
      return response.data
    } catch (error: any) {
      ElMessage.error(error.message || '获取用户信息失败')
      throw error
    }
  }

  // 检查权限
  const hasRole = (roles: string[]) => {
    if (!user.value) return false
    return roles.includes(user.value.role)
  }

  const isAdmin = computed(() => user.value?.role === 'admin')
  const isCommittee = computed(() => user.value?.role === 'committee' || user.value?.role === 'admin')

  return {
    // 状态
    user,
    accessToken,
    isLoggedIn,
    userRole,
    userName,
    userAvatar,
    isAdmin,
    isCommittee,

    // 方法
    initAuth,
    login,
    register,
    logout,
    refreshAccessToken,
    changePassword,
    fetchUserInfo,
    hasRole,
    clearAuth
  }
})
