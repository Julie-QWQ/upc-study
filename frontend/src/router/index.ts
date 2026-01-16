import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import DefaultLayout from '@/layouts/DefaultLayout.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/Login.vue'),
    meta: { title: '登录', requiresAuth: false, guest: true }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/auth/Register.vue'),
    meta: { title: '注册', requiresAuth: false, guest: true }
  },
  {
    path: '/forgot-password',
    name: 'ForgotPassword',
    component: () => import('@/views/auth/ForgotPassword.vue'),
    meta: { title: '忘记密码', requiresAuth: false, guest: true }
  },
  // 使用 DefaultLayout 的路由组（带侧边栏）
  {
    path: '/',
    component: DefaultLayout,
    children: [
      // 首页（重定向到资料列表）
      {
        path: '',
        redirect: '/materials'
      },
      // 资料管理相关路由
      {
        path: 'materials',
        name: 'Materials',
        component: () => import('@/views/material/List.vue'),
        meta: { title: '资料列表', requiresAuth: true }
      },
      {
        path: 'materials/:id',
        name: 'MaterialDetail',
        component: () => import('@/views/material/Detail.vue'),
        props: true, // 将路由参数 :id 作为 props 传递给组件
        meta: { title: '资料详情', requiresAuth: true }
      },
      {
        path: 'materials/:id/edit',
        name: 'MaterialEdit',
        component: () => import('@/views/material/Edit.vue'),
        props: true,
        meta: { title: '编辑资料', requiresAuth: true, roles: ['committee', 'admin'] }
      },
      {
        path: 'materials/upload',
        name: 'MaterialUpload',
        component: () => import('@/views/material/Upload.vue'),
        meta: { title: '上传资料', requiresAuth: true, roles: ['committee', 'admin'] }
      },
      {
        path: 'materials/my',
        name: 'MyMaterials',
        component: () => import('@/views/material/MyMaterials.vue'),
        meta: { title: '我的资料', requiresAuth: true, roles: ['committee', 'admin'] }
      },
      {
        path: 'favorites',
        name: 'Favorites',
        component: () => import('@/views/material/Favorites.vue'),
        meta: { title: '我的收藏', requiresAuth: true }
      },
      {
        path: 'downloads',
        name: 'Downloads',
        component: () => import('@/views/material/Downloads.vue'),
        meta: { title: '下载历史', requiresAuth: true }
      },
      // 搜索相关路由
      {
        path: 'search',
        name: 'Search',
        component: () => import('@/views/search/Search.vue'),
        meta: { title: '资料搜索', requiresAuth: true }
      },
      {
        path: 'hot-materials',
        name: 'HotMaterials',
        component: () => import('@/views/search/HotMaterials.vue'),
        meta: { title: '热门资料', requiresAuth: true }
      },
      // 学委申请相关路由
      {
        path: 'committee/apply',
        name: 'CommitteeApply',
        component: () => import('@/views/committee/Apply.vue'),
        meta: { title: '申请学委', requiresAuth: true }
      },
      {
        path: 'committee/applications',
        name: 'CommitteeApplications',
        component: () => import('@/views/committee/Applications.vue'),
        meta: { title: '我的申请', requiresAuth: true }
      },
      // 公告相关路由
      {
        path: 'announcements',
        name: 'Announcements',
        component: () => import('@/views/announcement/List.vue'),
        meta: { title: '公告列表', requiresAuth: true }
      },
      // 管理后台相关路由
      {
        path: 'admin/dashboard',
        name: 'AdminDashboard',
        component: () => import('@/views/admin/Dashboard.vue'),
        meta: { title: '仪表盘', requiresAuth: true, roles: ['admin'] }
      },
      {
        path: 'admin/materials',
        name: 'AdminMaterials',
        component: () => import('@/views/admin/Materials.vue'),
        meta: { title: '资料审核', requiresAuth: true, roles: ['admin'] }
      },
      {
        path: 'admin/applications',
        name: 'AdminApplications',
        component: () => import('@/views/admin/Applications.vue'),
        meta: { title: '学委申请审核', requiresAuth: true, roles: ['admin'] }
      },
      {
        path: 'admin/reports',
        name: 'AdminReports',
        component: () => import('@/views/admin/Reports.vue'),
        meta: { title: '举报处理', requiresAuth: true, roles: ['admin'] }
      },
      {
        path: 'admin/users',
        name: 'AdminUsers',
        component: () => import('@/views/admin/Users.vue'),
        meta: { title: '用户管理', requiresAuth: true, roles: ['admin'] }
      },
      {
        path: 'admin/settings',
        name: 'AdminSettings',
        component: () => import('@/views/admin/Settings.vue'),
        meta: { title: '系统设置', requiresAuth: true, roles: ['admin'] }
      },
      {
        path: 'admin/material-categories',
        name: 'AdminMaterialCategories',
        component: () => import('@/views/admin/MaterialCategories.vue'),
        meta: { title: '资料类型管理', requiresAuth: true, roles: ['admin'] }
      },
      {
        path: 'admin/announcements',
        name: 'AdminAnnouncements',
        component: () => import('@/views/admin/Announcements.vue'),
        meta: { title: '公告管理', requiresAuth: true, roles: ['admin'] }
      },
      // 通知相关路由
      {
        path: 'notifications',
        name: 'Notifications',
        component: () => import('@/views/notification/List.vue'),
        meta: { title: '通知中心', requiresAuth: true }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/error/404.vue'),
    meta: { title: '404' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  document.title = `${to.meta.title || 'UPC-DocHub'} | UPC-DocHub`

  // 初始化认证状态
  const authStore = useAuthStore()
  authStore.initAuth()

  // 检查是否需要认证
  const requiresAuth = to.meta.requiresAuth !== false
  const isLoggedIn = authStore.isLoggedIn

  if (requiresAuth && !isLoggedIn) {
    // 需要认证但未登录，跳转到登录页
    next({
      name: 'Login',
      query: { redirect: to.fullPath }
    })
    return
  }

  // 检查角色权限
  const requiredRoles = to.meta.roles as string[] | undefined
  if (requiredRoles && requiredRoles.length > 0) {
    const userRole = authStore.userRole
    if (!userRole || !requiredRoles.includes(userRole)) {
      // 权限不足，跳转到首页
      next({ name: 'Materials' })
      return
    }
  }

  // 允许所有导航（包括已登录用户访问游客页面）
  next()
})

export default router

