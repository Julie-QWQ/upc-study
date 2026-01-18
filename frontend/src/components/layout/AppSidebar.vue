<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useNavigationStore } from '@/stores/navigation'
import { ElMessage } from 'element-plus'
import { authApi } from '@/api/auth'
import SiteName from '@/components/SiteName.vue'

const router = useRouter()
const authStore = useAuthStore()
const navigationStore = useNavigationStore()

// 用户信息
const userName = computed(() => authStore.userName)
const userAvatar = computed(() => authStore.userAvatar)

// 生成默认头像
const defaultAvatar = computed(() => {
  return userName.value ? userName.value.charAt(0).toUpperCase() : 'U'
})

// 用户角色标签和样式
const roleLabel = computed(() => {
  if (authStore.isAdmin) return '管理员'
  if (authStore.isCommittee) return '学委'
  return '学生'
})

const roleClass = computed(() => {
  if (authStore.isAdmin) return 'role-admin'
  if (authStore.isCommittee) return 'role-committee'
  return 'role-student'
})

// 用户菜单项
const userNavItems = [
  {
    name: '首页',
    path: '/materials',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
      <polyline points="9 22 9 12 15 12 15 22"/>
    </svg>`
  },
  {
    name: '资料搜索',
    path: '/search',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <circle cx="11" cy="11" r="8"/>
      <path d="m21 21-4.35-4.35"/>
      <path d="M11 8a3 3 0 0 0-3 3"/>
      <path d="M14 11a3 3 0 0 0-3-3"/>
    </svg>`
  },
  {
    name: '热门资料',
    path: '/hot-materials',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M6 9H4.5a2.5 2.5 0 0 1 0-5H6"/>
      <path d="M18 9h1.5a2.5 2.5 0 0 0 0-5H18"/>
      <path d="M4 22h16"/>
      <path d="M10 14.66V17c0 .55-.47.98-.97 1.21C7.85 18.75 7 20.24 7 22"/>
      <path d="M14 14.66V17c0 .55.47.98.97 1.21C16.15 18.75 17 20.24 17 22"/>
      <path d="M18 2H6v7a6 6 0 0 0 12 0V2Z"/>
    </svg>`
  },
  {
    name: '上传资料',
    path: '/materials/upload',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
      <polyline points="17 8 12 3 7 8"/>
      <line x1="12" y1="3" x2="12" y2="15"/>
    </svg>`,
    roles: ['committee', 'admin']
  },
  {
    name: '我的收藏',
    path: '/favorites',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z"/>
    </svg>`
  },
  {
    name: '下载历史',
    path: '/downloads',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
      <polyline points="7 10 12 15 17 10"/>
      <line x1="12" y1="15" x2="12" y2="3"/>
    </svg>`
  },
  {
    name: '申请学委',
    path: '/committee/apply',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/>
      <circle cx="9" cy="7" r="4"/>
      <line x1="19" y1="8" x2="19" y2="14"/>
      <line x1="22" y1="11" x2="16" y2="11"/>
    </svg>`,
    roles: ['student']
  },
  {
    name: '我的申请',
    path: '/committee/applications',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
      <polyline points="14 2 14 8 20 8"/>
      <line x1="16" y1="13" x2="8" y2="13"/>
      <line x1="16" y1="17" x2="8" y2="17"/>
      <polyline points="10 9 9 9 8 9"/>
    </svg>`
  }
]

// 管理员菜单项
const adminNavItems = [
  {
    name: '仪表盘',
    path: '/admin/dashboard',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <rect x="3" y="3" width="7" height="9"/>
      <rect x="14" y="3" width="7" height="5"/>
      <rect x="14" y="12" width="7" height="9"/>
      <rect x="3" y="16" width="7" height="5"/>
    </svg>`
  },
  {
    name: '用户管理',
    path: '/admin/users',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
      <circle cx="9" cy="7" r="4"/>
      <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
      <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
    </svg>`
  },
  {
    name: '资料审核',
    path: '/admin/materials',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
      <polyline points="14 2 14 8 20 8"/>
      <path d="M9 15l2 2 4-4"/>
    </svg>`
  },
  {
    name: '资料类型管理',
    path: '/admin/material-categories',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M9 11l3 3L22 4"/>
      <path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"/>
      <path d="M9 11h6"/>
      <path d="M9 15h6"/>
      <path d="M9 7h6"/>
    </svg>`
  },
  {
    name: '学委申请审核',
    path: '/admin/applications',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"/>
      <rect x="8" y="2" width="8" height="4" rx="1" ry="1"/>
      <path d="M9 14h6"/>
      <path d="M9 10h6"/>
      <path d="M9 18h6"/>
    </svg>`
  },
  {
    name: '举报处理',
    path: '/admin/reports',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/>
      <line x1="12" y1="9" x2="12" y2="13"/>
      <line x1="12" y1="17" x2="12.01" y2="17"/>
    </svg>`
  },
  {
    name: '公告管理',
    path: '/admin/announcements',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M2 12h5"/>
      <path d="M4 12v6a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-6"/>
      <path d="M5 12V7a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2v5"/>
      <path d="M15 12v-2"/>
      <path d="M15 12h5"/>
      <path d="M18 12v-2"/>
    </svg>`
  },
  {
    name: '系统设置',
    path: '/admin/settings',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <circle cx="12" cy="12" r="3"/>
      <path d="M12 1v6m0 6v6"/>
      <path d="m4.93 4.93 4.24 4.24m5.66 5.66 4.24 4.24M1 12h6m6 0h6"/>
      <path d="m4.93 19.07 4.24-4.24m5.66-5.66 4.24-4.24"/>
    </svg>`
  }
]

// 根据用户角色显示不同的菜单项
const navItems = computed(() => {
  // 过滤用户菜单项，根据角色显示
  const filteredUserItems = userNavItems.filter(item => {
    if (!item.roles) return true // 没有角色限制，所有人可见
    return item.roles.some(role => {
      if (role === 'committee') return authStore.isCommittee || authStore.isAdmin
      if (role === 'admin') return authStore.isAdmin
      if (role === 'student') return !authStore.isCommittee && !authStore.isAdmin
      return true
    })
  })

  // 修改"我的申请"的路径和名称(学委显示我的资料,学生显示申请)
  const items = filteredUserItems.map(item => {
    if (item.name === '我的申请') {
      return {
        ...item,
        name: authStore.isCommittee ? '我的资料' : '我的申请',
        path: authStore.isCommittee ? '/materials/my' : '/committee/applications'
      }
    }
    return item
  })

  if (authStore.isAdmin) {
    items.push(...adminNavItems)
  }
  return items
})

// 判断当前激活的菜单项（使用导航状态，响应更快速）
const isActive = (path: string) => {
  // 使用导航状态而不是路由状态，避免等待路由切换
  return navigationStore.isPathActive(path)
}

// 导航跳转
const navigate = (path: string) => {
  router.push(path)
}

// 退出登录
const handleLogout = async () => {
  // 立即清除本地认证状态
  authStore.clearAuth()

  // 显示成功消息
  ElMessage.success('已退出登录')

  // 异步调用 logout API(不等待,避免阻塞跳转)
  authApi.logout().catch(error => {
    console.error('Logout API failed:', error)
  })

  // 立即跳转到登录页
  router.push('/login')
}
</script>

<template>
  <aside class="app-sidebar">
    <div class="sidebar-header">
      <h1 class="app-title"><SiteName /></h1>

      <!-- 用户信息卡片 -->
      <div class="user-info-card">
        <div class="user-avatar">
          <img v-if="userAvatar" :src="userAvatar" :alt="userName" />
          <span v-else class="default-avatar">{{ defaultAvatar }}</span>
        </div>
        <div class="user-details">
          <div class="user-name">{{ userName || '未登录' }}</div>
          <div class="user-role-badge" :class="roleClass">
            {{ roleLabel }}
          </div>
        </div>
      </div>
    </div>

    <nav class="sidebar-nav">
      <button
        v-for="item in navItems"
        :key="item.path"
        :class="['nav-item', { active: isActive(item.path) }]"
        @click="navigate(item.path)"
      >
        <span class="nav-icon" v-html="item.icon"></span>
        <span class="nav-label">{{ item.name }}</span>
      </button>
    </nav>

    <div class="sidebar-footer">
      <button class="logout-btn" @click="handleLogout">
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
          <polyline points="16 17 21 12 16 7"/>
          <line x1="21" y1="12" x2="9" y2="12"/>
        </svg>
        <span>退出登录</span>
      </button>
    </div>
  </aside>
</template>

<style scoped lang="scss">
.app-sidebar {
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  width: 260px;
  background: #ffffff;
  border-right: 1px solid #f2f2f2;
  display: flex;
  flex-direction: column;
  padding: 24px 16px;
  z-index: 100;
}

.sidebar-header {
  padding: 0 8px 32px;
  border-bottom: 1px solid #f2f2f2;
  margin-bottom: 24px;

  .app-title {
    font-size: 24px;
    font-weight: 700;
    color: #1a1a1a;
    margin: 0 0 16px 0;
    letter-spacing: -0.02em;
  }

  .user-info-card {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px;
    background: #fafafa;
    border-radius: 12px;
    border: 1px solid #f2f2f2;

    .user-avatar {
      width: 48px;
      height: 48px;
      border-radius: 50%;
      overflow: hidden;
      flex-shrink: 0;
      background: #e5e5e5;
      display: flex;
      align-items: center;
      justify-content: center;

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }

      .default-avatar {
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: #ffffff;
        font-size: 20px;
        font-weight: 600;
      }
    }

    .user-details {
      flex: 1;
      min-width: 0;

      .user-name {
        font-size: 14px;
        font-weight: 600;
        color: #1a1a1a;
        margin-bottom: 4px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      .user-role-badge {
        display: inline-block;
        padding: 2px 8px;
        font-size: 11px;
        font-weight: 600;
        border-radius: 10px;
        text-transform: uppercase;
        letter-spacing: 0.05em;

        &.role-student {
          background: #f0f9ff;
          color: #0284c7;
        }

        &.role-committee {
          background: #fef3c7;
          color: #d97706;
        }

        &.role-admin {
          background: #fee2e2;
          color: #dc2626;
        }
      }
    }
  }
}

.sidebar-nav {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
  overflow-y: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: none;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s;
  font-size: 15px;
  color: #333;
  text-align: left;
  width: 100%;

  .nav-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    color: #666;
    transition: color 0.15s;
  }

  .nav-label {
    font-weight: 400;
  }

  &:hover {
    background: #f5f5f5;
  }

  &.active {
    background: #1a1a1a;
    color: #ffffff;

    .nav-icon {
      color: #ffffff;
    }

    .nav-label {
      font-weight: 500;
    }
  }
}

.sidebar-footer {
  padding-top: 16px;
  border-top: 1px solid #f2f2f2;
}

.logout-btn {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: none;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s;
  font-size: 15px;
  color: #666;
  width: 100%;

  svg {
    color: #666;
  }

  &:hover {
    background: #fef2f2;
    color: #dc2626;

    svg {
      color: #dc2626;
    }
  }
}

// 响应式
@media (max-width: 1024px) {
  .app-sidebar {
    width: 80px;
    padding: 16px 12px;

    .sidebar-header {
      padding: 0 4px 16px;

      .app-title {
        font-size: 18px;
        text-align: center;
        margin-bottom: 12px;
      }

      .user-info-card {
        flex-direction: column;
        padding: 8px;

        .user-avatar {
          width: 40px;
          height: 40px;
        }

        .user-details {
          display: none;
        }
      }
    }

    .nav-item {
      padding: 12px;
      justify-content: center;

      .nav-label {
        display: none;
      }
    }

    .logout-btn {
      padding: 12px;
      justify-content: center;

      span {
        display: none;
      }
    }
  }
}

@media (max-width: 768px) {
  .app-sidebar {
    display: none;
  }
}
</style>
