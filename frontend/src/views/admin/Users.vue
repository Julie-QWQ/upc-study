<template>
  <div class="users-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1>用户管理</h1>
      <div class="header-info">
        <span class="total-count">共 {{ pagination.total }} 位用户</span>
      </div>
    </div>

    <!-- 搜索栏 -->
    <div class="search-bar">
      <div class="search-input">
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="11" cy="11" r="8"></circle>
          <path d="m21 21-4.35-4.35"></path>
        </svg>
        <input
          v-model="searchForm.keyword"
          type="text"
          placeholder="搜索用户名、姓名或邮箱"
          @keyup.enter="handleSearch"
        />
      </div>

      <select v-model="searchForm.role" class="filter-select" @change="handleSearch">
        <option value="">全部角色</option>
        <option value="student">学生</option>
        <option value="committee">学委</option>
        <option value="admin">管理员</option>
      </select>

      <select v-model="searchForm.status" class="filter-select" @change="handleSearch">
        <option value="">全部状态</option>
        <option value="active">正常</option>
        <option value="banned">已封禁</option>
      </select>

      <button v-if="hasSearch" class="reset-btn" @click="handleReset">
        重置
      </button>
    </div>

    <!-- 用户列表 -->
    <div v-loading="loading" class="users-list">
      <div v-for="user in userList" :key="user.id" class="user-card">
        <div class="user-avatar">
          <div class="avatar-placeholder">
            {{ (user.real_name || user.username).charAt(0).toUpperCase() }}
          </div>
        </div>

        <div class="user-info">
          <div class="user-name">
            {{ user.real_name || user.username }}
            <span v-if="user.role" :class="['role-badge', user.role]">
              {{ getRoleText(user.role) }}
            </span>
          </div>
          <div class="user-details">
            <span class="detail-item">@{{ user.username }}</span>
            <span class="detail-separator">·</span>
            <span class="detail-item">{{ user.email }}</span>
          </div>
          <div v-if="user.major || user.class" class="user-meta">
            <span v-if="user.major">{{ user.major }}</span>
            <span v-if="user.class">{{ user.class }}</span>
          </div>
        </div>

        <div class="user-status">
          <span :class="['status-badge', user.status]">
            {{ getStatusText(user.status) }}
          </span>
        </div>

        <div class="user-time">
          {{ formatDate(user.created_at) }}
        </div>

        <div class="user-actions">
          <button class="action-btn" @click="handleViewDetail(user)">查看</button>
          <button class="action-btn" @click="handleUpdateStatus(user)">状态</button>
          <button class="action-btn danger" @click="handleDelete(user)">删除</button>
        </div>
      </div>

      <div v-if="!loading && userList.length === 0" class="empty-state">
        <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
          <circle cx="9" cy="7" r="4"></circle>
          <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
          <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
        </svg>
        <p>暂无用户数据</p>
      </div>
    </div>

    <!-- 分页 -->
    <div v-if="pagination.total > 0" class="pagination">
      <span class="pagination-info">
        显示 {{ (pagination.page - 1) * pagination.page_size + 1 }}-
        {{ Math.min(pagination.page * pagination.page_size, pagination.total) }}
        / 共 {{ pagination.total }} 条
      </span>
      <div class="pagination-controls">
        <button
          class="page-btn"
          :disabled="pagination.page === 1"
          @click="changePage(pagination.page - 1)"
        >
          上一页
        </button>
        <button
          class="page-btn"
          :disabled="pagination.page * pagination.page_size >= pagination.total"
          @click="changePage(pagination.page + 1)"
        >
          下一页
        </button>
      </div>
    </div>

    <!-- 状态更新对话框 -->
    <!-- 用户详情 -->
    <div v-if="detailDialogVisible" class="dialog-overlay" @click.self="detailDialogVisible = false">
      <div class="dialog dialog-large">
        <div class="dialog-header">
          <h3>用户详情</h3>
          <button class="close-btn" @click="detailDialogVisible = false">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 6 6 18"></path>
              <path d="m6 6 12 12"></path>
            </svg>
          </button>
        </div>

        <div class="dialog-body">
          <div v-if="detailLoading" class="detail-loading">加载中...</div>
          <div v-else-if="!userDetail" class="detail-empty">暂无用户详情</div>
          <div v-else class="detail-grid">
            <div class="detail-section">
              <h4>基本信息</h4>
              <div class="detail-row">
                <span class="detail-label">用户名</span>
                <span class="detail-value">{{ userDetail.user.username || '-' }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">真实姓名</span>
                <span class="detail-value">{{ userDetail.user.real_name || '-' }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">邮箱</span>
                <span class="detail-value">{{ userDetail.user.email || '-' }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">手机号</span>
                <span class="detail-value">{{ userDetail.user.phone || '-' }}</span>
              </div>
            </div>

            <div class="detail-section">
              <h4>身份信息</h4>
              <div class="detail-row">
                <span class="detail-label">角色</span>
                <span class="detail-value">{{ getRoleText(userDetail.user.role) }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">状态</span>
                <span class="detail-value">{{ getStatusText(userDetail.user.status) }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">专业</span>
                <span class="detail-value">{{ userDetail.user.major || '-' }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">班级</span>
                <span class="detail-value">{{ userDetail.user.class || '-' }}</span>
              </div>
            </div>

            <div class="detail-section">
              <h4>记录统计</h4>
              <div class="detail-row">
                <span class="detail-label">下载总次数</span>
                <span class="detail-value">{{ userDetail.download_total ?? 0 }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">上传资料</span>
                <span class="detail-value">{{ userDetail.upload_total ?? 0 }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">收藏数量</span>
                <span class="detail-value">{{ userDetail.favorite_total ?? 0 }}</span>
              </div>
            </div>

            <div class="detail-section">
              <h4>时间信息</h4>
              <div class="detail-row">
                <span class="detail-label">注册时间</span>
                <span class="detail-value">{{ formatDate(userDetail.user.created_at) }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">最近登录</span>
                <span class="detail-value">{{ formatDate(userDetail.user.last_login_at) }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="dialog-footer">
          <button class="btn btn-secondary" @click="detailDialogVisible = false">关闭</button>
        </div>
      </div>
    </div>

    <!-- 状态更新对话框 -->
    <div v-if="statusDialogVisible" class="dialog-overlay" @click.self="statusDialogVisible = false">
      <div class="dialog">
        <div class="dialog-header">
          <h3>更新用户状态</h3>
          <button class="close-btn" @click="statusDialogVisible = false">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 6 6 18"></path>
              <path d="m6 6 12 12"></path>
            </svg>
          </button>
        </div>

        <div class="dialog-body">
          <div class="form-group">
            <label>用户</label>
            <div class="user-tag">{{ currentUser?.username }}</div>
          </div>

          <div class="form-group">
            <label>当前状态</label>
            <span :class="['status-badge', currentUser?.status]">
              {{ getStatusText(currentUser?.status) }}
            </span>
          </div>

          <div class="form-group">
            <label>新状态</label>
            <select v-model="statusForm.status" class="form-select">
              <option value="active">正常</option>
              <option value="banned">已封禁</option>
            </select>
          </div>

          <div class="form-group">
            <label>原因</label>
            <textarea
              v-model="statusForm.reason"
              class="form-textarea"
              rows="3"
              placeholder="请输入原因(选填)"
            ></textarea>
          </div>
        </div>

        <div class="dialog-footer">
          <button class="btn btn-secondary" @click="statusDialogVisible = false">取消</button>
          <button class="btn btn-primary" @click="confirmUpdateStatus">确定</button>
        </div>
      </div>
    </div>

    <!-- 删除确认对话框 -->
    <div v-if="deleteDialogVisible" class="dialog-overlay" @click.self="deleteDialogVisible = false">
      <div class="dialog dialog-small">
        <div class="dialog-header">
          <h3>确认删除</h3>
          <button class="close-btn" @click="deleteDialogVisible = false">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 6 6 18"></path>
              <path d="m6 6 12 12"></path>
            </svg>
          </button>
        </div>

        <div class="dialog-body">
          <div class="delete-confirmation">
            <div class="warning-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M21.73 18l-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"></path>
                <line x1="12" y1="9" x2="12" y2="13"></line>
                <line x1="12" y1="17" x2="12.01" y2="17"></line>
              </svg>
            </div>
            <p class="delete-message">
              确定要删除用户 <strong>{{ currentUser?.username }}</strong> 吗？
            </p>
            <p class="delete-warning">此操作不可恢复，该用户的所有数据将被永久删除。</p>
          </div>
        </div>

        <div class="dialog-footer">
          <button class="btn btn-secondary" @click="deleteDialogVisible = false">取消</button>
          <button class="btn btn-danger" @click="confirmDelete">删除</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getUserList, getUserDetail, updateUserStatus, deleteUser } from '@/api/admin'
import type { User, UserDetailResponse } from '@/api/admin'

const loading = ref(false)
const userList = ref<User[]>([])

const searchForm = reactive({
  keyword: '',
  role: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const statusDialogVisible = ref(false)
const deleteDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const detailLoading = ref(false)
const currentUser = ref<User>()
const userDetail = ref<UserDetailResponse | null>(null)
const statusForm = reactive({
  status: '',
  reason: ''
})

const hasSearch = computed(() => {
  return searchForm.keyword || searchForm.role || searchForm.status
})

async function loadUserList() {
  loading.value = true
  try {
    const { data } = await getUserList({
      page: pagination.page,
      page_size: pagination.page_size,
      ...searchForm
    })
    userList.value = data.list || []
    pagination.total = data.total || 0
  } catch (error: any) {
    ElMessage.error(error.message || '获取用户列表失败')
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  pagination.page = 1
  loadUserList()
}

function handleReset() {
  Object.assign(searchForm, {
    keyword: '',
    role: '',
    status: ''
  })
  handleSearch()
}

function changePage(newPage: number) {
  pagination.page = newPage
  loadUserList()
}

async function handleViewDetail(user: User) {
  detailDialogVisible.value = true
  detailLoading.value = true
  userDetail.value = null
  try {
    const { data } = await getUserDetail(user.id)
    userDetail.value = data
  } catch (error: any) {
    ElMessage.error(error.message || '获取用户详情失败')
  } finally {
    detailLoading.value = false
  }
}

function handleUpdateStatus(user: User) {
  currentUser.value = user
  statusForm.status = user.status
  statusForm.reason = ''
  statusDialogVisible.value = true
}

async function confirmUpdateStatus() {
  if (!currentUser.value) return

  try {
    await updateUserStatus(currentUser.value.id, statusForm)
    ElMessage.success('更新状态成功')
    statusDialogVisible.value = false
    loadUserList()
  } catch (error: any) {
    ElMessage.error(error.message || '更新状态失败')
  }
}

function handleDelete(user: User) {
  currentUser.value = user
  deleteDialogVisible.value = true
}

async function confirmDelete() {
  if (!currentUser.value) return

  try {
    await deleteUser(currentUser.value.id)
    ElMessage.success('删除成功')
    deleteDialogVisible.value = false
    loadUserList()
  } catch (error: any) {
    ElMessage.error(error.message || '删除失败')
  }
}

function formatDate(date: string) {
  const d = new Date(date)
  if (Number.isNaN(d.getTime())) {
    return '-'
  }
  const now = new Date()
  const diff = now.getTime() - d.getTime()
  if (diff < 0) return '刚刚'
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (days === 0) return '今天'
  if (days === 1) return '昨天'
  if (days < 7) return `${days} 天前`
  if (days < 30) return `${Math.floor(days / 7)} 周前`
  return d.toLocaleDateString('zh-CN')
}

function getRoleText(role: string) {
  const map: Record<string, string> = {
    admin: '管理员',
    committee: '学委',
    student: '学生'
  }
  return map[role] || role
}

function getStatusText(status?: string) {
  const map: Record<string, string> = {
    active: '正常',
    banned: '已封禁'
  }
  return map[status || ''] || status || ''
}

onMounted(() => {
  loadUserList()
})
</script>

<style scoped lang="scss">
.users-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px;
}

.page-header {
  margin-bottom: 24px;

  h1 {
    font-size: 28px;
    font-weight: 700;
    color: #111827;
    margin: 0 0 8px 0;
  }

  .header-info {
    .total-count {
      font-size: 14px;
      color: #6b7280;
    }
  }
}

.search-bar {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
  flex-wrap: wrap;
}

.search-input {
  flex: 1;
  min-width: 250px;
  position: relative;
  display: flex;
  align-items: center;
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 0 12px;

  svg {
    color: #9ca3af;
    flex-shrink: 0;
  }

  input {
    flex: 1;
    border: none;
    outline: none;
    padding: 10px 8px;
    font-size: 14px;
    color: #111827;
    background: transparent;

    &::placeholder {
      color: #9ca3af;
    }
  }

  &:focus-within {
    border-color: #111827;
  }
}

.filter-select {
  padding: 10px 14px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  background: #ffffff;
  color: #111827;
  font-size: 14px;
  cursor: pointer;
  outline: none;
  transition: all 0.2s;

  &:focus {
    border-color: #111827;
  }
}

.reset-btn {
  padding: 10px 16px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  background: #ffffff;
  color: #6b7280;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;

  &:hover {
    background: #f9fafb;
    border-color: #d1d5db;
  }
}

.users-list {
  min-height: 300px;
}

.user-card {
  display: grid;
  grid-template-columns: auto 1fr auto auto auto;
  gap: 16px;
  align-items: center;
  padding: 20px;
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  margin-bottom: 12px;
  transition: all 0.2s;

  &:hover {
    border-color: #d1d5db;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  }

  @media (max-width: 768px) {
    grid-template-columns: auto 1fr;
    gap: 12px;

    .user-status,
    .user-time {
      display: none;
    }
  }
}

.user-avatar {
  .avatar-placeholder {
    width: 48px;
    height: 48px;
    border-radius: 50%;
    background: #f3f4f6;
    color: #6b7280;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 18px;
    font-weight: 600;
  }
}

.user-info {
  min-width: 0;
}

.user-name {
  font-size: 15px;
  font-weight: 600;
  color: #111827;
  margin-bottom: 4px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.role-badge {
  padding: 2px 8px;
  font-size: 12px;
  font-weight: 500;
  border-radius: 4px;

  &.admin {
    background: #fee2e2;
    color: #b91c1c;
  }

  &.committee {
    background: #fef3c7;
    color: #b45309;
  }

  &.student {
    background: #e5e7eb;
    color: #6b7280;
  }
}

.user-details {
  font-size: 13px;
  color: #6b7280;
  margin-bottom: 4px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.detail-separator {
  color: #d1d5db;
}

.user-meta {
  font-size: 12px;
  color: #9ca3af;
  display: flex;
  gap: 8px;

  span:not(:last-child)::after {
    content: '·';
    margin-left: 8px;
  }
}

.status-badge {
  padding: 4px 10px;
  font-size: 12px;
  font-weight: 500;
  border-radius: 6px;
  white-space: nowrap;

  &.active {
    background: #dcfce7;
    color: #15803d;
  }

  &.banned {
    background: #fee2e2;
    color: #b91c1c;
  }
}

.user-time {
  font-size: 13px;
  color: #9ca3af;
  white-space: nowrap;
}

.user-actions {
  display: flex;
  gap: 8px;

  @media (max-width: 768px) {
    grid-column: 2;
    justify-content: flex-end;
  }
}

.action-btn {
  padding: 6px 12px;
  border: none;
  border-radius: 6px;
  background: #f3f4f6;
  color: #374151;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;

  &:hover {
    background: #e5e7eb;
  }

  &.danger {
    color: #b91c1c;
    background: #fef2f2;

    &:hover {
      background: #fee2e2;
    }
  }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: #9ca3af;

  svg {
    margin-bottom: 16px;
  }

  p {
    margin: 0;
    font-size: 14px;
  }
}

.pagination {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid #e5e7eb;
}

.pagination-info {
  font-size: 13px;
  color: #6b7280;
}

.pagination-controls {
  display: flex;
  gap: 8px;
}

.page-btn {
  padding: 8px 16px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  background: #ffffff;
  color: #111827;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;

  &:hover:not(:disabled) {
    background: #f9fafb;
    border-color: #d1d5db;
  }

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
}

.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.dialog {
  background: #ffffff;
  border-radius: 12px;
  width: 90%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
}

.dialog-large {
  max-width: 760px;
}

.detail-loading,
.detail-empty {
  padding: 12px 0;
  color: #6b7280;
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.detail-section {
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  padding: 12px 14px;
  background: #fafafa;
}

.detail-section h4 {
  margin: 0 0 10px 0;
  font-size: 14px;
  font-weight: 600;
  color: #111827;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  font-size: 13px;
  color: #4b5563;
  padding: 4px 0;
}

.detail-label {
  color: #6b7280;
  flex-shrink: 0;
}

.detail-value {
  color: #111827;
  text-align: right;
  word-break: break-all;
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #e5e7eb;

  h3 {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
    color: #111827;
  }

  .close-btn {
    background: none;
    border: none;
    padding: 4px;
    cursor: pointer;
    color: #6b7280;
    transition: color 0.2s;

    &:hover {
      color: #111827;
    }
  }
}

.dialog-body {
  padding: 20px;
}

.form-group {
  margin-bottom: 20px;

  &:last-child {
    margin-bottom: 0;
  }

  label {
    display: block;
    font-size: 13px;
    font-weight: 600;
    color: #374151;
    margin-bottom: 8px;
  }
}

.user-tag {
  display: inline-block;
  padding: 6px 12px;
  background: #f3f4f6;
  border-radius: 6px;
  font-size: 14px;
  color: #111827;
}

.form-select {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  color: #111827;
  outline: none;
  transition: all 0.2s;

  &:focus {
    border-color: #111827;
  }
}

.form-textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  color: #111827;
  outline: none;
  resize: vertical;
  font-family: inherit;
  transition: all 0.2s;

  &:focus {
    border-color: #111827;
  }

  &::placeholder {
    color: #9ca3af;
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 20px;
  border-top: 1px solid #e5e7eb;
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;

  &.btn-secondary {
    background: #f3f4f6;
    color: #111827;

    &:hover {
      background: #e5e7eb;
    }
  }

  &.btn-primary {
    background: #111827;
    color: #ffffff;

    &:hover {
      background: #000000;
    }
  }

  &.btn-danger {
    background: #b91c1c;
    color: #ffffff;

    &:hover {
      background: #991b1b;
    }
  }
}

.dialog-small {
  max-width: 420px;
}

.delete-confirmation {
  text-align: center;
  padding: 12px 0;
}

.warning-icon {
  display: inline-flex;
  justify-content: center;
  align-items: center;
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: #fef3c7;
  color: #b45309;
  margin-bottom: 20px;
}

.delete-message {
  font-size: 16px;
  font-weight: 500;
  color: #111827;
  margin: 0 0 12px 0;
  line-height: 1.5;

  strong {
    font-weight: 600;
    color: #b91c1c;
  }
}

.delete-warning {
  font-size: 14px;
  color: #6b7280;
  margin: 0;
  line-height: 1.5;
}

@media (max-width: 768px) {
  .detail-grid {
    grid-template-columns: 1fr;
  }
}

</style>
