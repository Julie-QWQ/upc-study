<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useMaterialStore } from '@/stores/material'
import { useAuth } from '@/composables/useAuth'
import type { Material } from '@/types'

interface Props {
  materialId: number | string
}

const props = defineProps<Props>()
const router = useRouter()
const materialStore = useMaterialStore()
const { isAdmin, isCommittee, user } = useAuth()

const loading = computed(() => materialStore.loading)
const material = computed(() => materialStore.currentMaterial)
const showDropdown = ref(false)
const showReportModal = ref(false)
const deleteDialogVisible = ref(false)
const reportReason = ref('')

// 资料分类映射
const categoryMap: Record<string, string> = {
  textbook: '教材',
  reference: '参考书',
  exam_paper: '试卷',
  note: '笔记',
  exercise: '习题',
  experiment: '实验指导',
  thesis: '论文',
  other: '其他'
}

// 资料状态映射
const statusMap: Record<string, { text: string; type: any }> = {
  pending: { text: '待审核', type: 'warning' },
  approved: { text: '已通过', type: 'success' },
  rejected: { text: '已拒绝', type: 'danger' },
  deleted: { text: '已删除', type: 'info' }
}

// 格式化文件大小
const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
}

// 资料分类文本
const categoryText = computed(() => {
  return material.value ? categoryMap[material.value.category] : '其他'
})

// 资料状态信息
const statusInfo = computed(() => {
  return material.value ? statusMap[material.value.status] || { text: '未知', type: 'info' } : null
})

// 是否是上传者
const isUploader = computed(() => {
  return material.value && user.value ? material.value.uploader_id === user.value.id : false
})

// 加载资料详情
const loadDetail = async () => {
  const id = Number(props.materialId)

  // 验证ID是否为有效数字
  if (isNaN(id) || id <= 0 || !Number.isInteger(id)) {
    ElMessage.error('无效的资料 ID')
    router.push('/materials')
    return
  }

  const response = await materialStore.fetchMaterial(id)

  // 检查响应是否成功
  if (response.code !== 0 || !material.value) {
    ElMessage.error(response.message || '资料不存在')
    router.push('/materials')
  }
}

// 下载资料
const handleDownload = async () => {
  if (!material.value) return

  try {
    const response = await materialStore.getDownloadUrl(material.value.id)

    if (response.code === 0 && response.data) {
      // 打开下载链接
      window.open(response.data.download_url, '_blank')

      // 更新下载次数
      if (material.value) {
        material.value.download_count++
      }
    } else {
      ElMessage.error(response.message || '获取下载链接失败')
    }
  } catch (error: any) {
    ElMessage.error(error.message || '下载失败')
  }
}

// 收藏/取消收藏
const toggleFavorite = async () => {
  if (!material.value) return

  try {
    if (material.value.is_favorited) {
      await materialStore.removeFavorite(material.value.id)
      ElMessage.success('已取消收藏')
    } else {
      const response = await materialStore.addFavorite(material.value.id)
      // 检查响应消息,如果已经收藏则提示
      if (response.message === '已收藏') {
        ElMessage.info('您已收藏过该资料')
      } else {
        ElMessage.success('收藏成功')
      }
    }
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  }
}

// 举报资料
const openReportModal = () => {
  showReportModal.value = true
  reportReason.value = ''
}

const closeReportModal = () => {
  showReportModal.value = false
  reportReason.value = ''
}

const submitReport = async () => {
  if (!material.value) return

  if (!reportReason.value.trim()) {
    ElMessage.warning('请输入举报原因')
    return
  }

  try {
    await materialStore.createReport(material.value.id, 'other', reportReason.value)
    ElMessage.success('举报成功，我们将尽快处理')
    closeReportModal()
  } catch (error: any) {
    ElMessage.error(error.message || '举报失败')
  }
}

// 编辑资料
const handleEdit = () => {
  router.push(`/materials/${material.value?.id}/edit`)
}

// 删除资料
const handleDelete = () => {
  deleteDialogVisible.value = true
}

// 确认删除
const confirmDelete = async () => {
  if (!material.value) return

  try {
    await materialStore.deleteMaterial(material.value.id)
    ElMessage.success('删除成功')
    deleteDialogVisible.value = false
    router.push('/materials')
  } catch (error: any) {
    ElMessage.error(error.message || '删除失败')
  }
}

// 切换下拉菜单
const toggleDropdown = () => {
  showDropdown.value = !showDropdown.value
}

// 关闭下拉菜单
const closeDropdown = () => {
  showDropdown.value = false
}

// 格式化日期
const formatDate = (dateString: string): string => {
  const date = new Date(dateString)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (days === 0) return '今天'
  if (days === 1) return '昨天'
  if (days < 7) return `${days} 天前`
  if (days < 30) return `${Math.floor(days / 7)} 周前`
  if (days < 365) return `${Math.floor(days / 30)} 月前`

  return date.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
}

// 点击外部关闭下拉菜单
onMounted(() => {
  loadDetail()

  document.addEventListener('click', (e) => {
    const target = e.target as HTMLElement
    if (!target.closest('.dropdown-wrapper')) {
      closeDropdown()
    }
  })
})
</script>

<template>
  <div class="substack-detail" v-loading="loading">
    <article v-if="material" class="detail-article">
      <!-- 顶部导航 -->
      <header class="article-header">
        <button class="back-button" @click="$router.back()">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M19 12H5M12 19l-7-7 7-7"/>
          </svg>
          返回
        </button>

        <div class="header-actions">
          <button
            :class="['action-button', { active: material.is_favorited }]"
            @click="toggleFavorite"
            :title="material.is_favorited ? '取消收藏' : '收藏'"
            data-testid="favorite-button"
          >
            <svg v-if="material.is_favorited" width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
            </svg>
            <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
            </svg>
          </button>

          <div v-if="isAdmin || isUploader" class="dropdown-wrapper">
            <button class="action-button" @click="toggleDropdown">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="1"/>
                <circle cx="12" cy="5" r="1"/>
                <circle cx="12" cy="19" r="1"/>
              </svg>
            </button>
            <div :class="['dropdown-menu', { show: showDropdown }]">
              <button v-if="isCommittee" @click="handleEdit">编辑</button>
              <button v-if="isAdmin" class="danger" @click="handleDelete">删除</button>
            </div>
          </div>
        </div>
      </header>

      <!-- 文章主体 -->
      <div class="article-content">
        <!-- 分类和状态标签 -->
        <div class="meta-tags">
          <span class="category-tag">{{ categoryText }}</span>
          <span :class="['status-tag', material.status]">{{ statusInfo?.text }}</span>
        </div>

        <!-- 标题 -->
        <h1 class="article-title">{{ material.title }}</h1>

        <!-- 作者信息 -->
        <div class="author-info">
          <div class="author-avatar">
            {{ material.uploader?.real_name?.[0] || material.uploader?.username?.[0] || 'U' }}
          </div>
          <div class="author-details">
            <div class="author-name">
              {{ material.uploader?.real_name || material.uploader?.username || '未知' }}
            </div>
            <div class="post-meta">
              {{ formatDate(material.created_at) }}
              <span class="separator">·</span>
              {{ material.course_name || '未分类' }}
            </div>
          </div>
        </div>

        <!-- 资料描述 -->
        <div v-if="material.description" class="article-body">
          <p>{{ material.description }}</p>
        </div>

        <!-- 文件信息卡片 -->
        <div class="file-info-card">
          <div class="file-icon">📄</div>
          <div class="file-details">
            <div class="file-name">{{ material.file_name }}</div>
            <div class="file-meta">{{ formatFileSize(material.file_size) }}</div>
          </div>
        </div>

        <!-- 统计信息 -->
        <div class="stats-grid">
          <div class="stat-item">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
              <polyline points="7 10 12 15 17 10"/>
              <line x1="12" y1="15" x2="12" y2="3"/>
            </svg>
            <span>{{ material.download_count }} 次下载</span>
          </div>
          <div class="stat-item">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
              <circle cx="12" cy="12" r="3"/>
            </svg>
            <span>{{ material.view_count }} 次浏览</span>
          </div>
          <div class="stat-item">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
            </svg>
            <span>{{ material.favorite_count }} 次收藏</span>
          </div>
        </div>

        <!-- 标签 -->
        <div v-if="material.tags && material.tags.length > 0" class="tags-section">
          <div v-for="tag in material.tags" :key="tag" class="tag-pill">#{{ tag }}</div>
        </div>

        <!-- 拒绝原因 -->
        <div v-if="material.status === 'rejected' && material.rejection_reason" class="rejection-notice">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="12" y1="8" x2="12" y2="12"/>
            <line x1="12" y1="16" x2="12.01" y2="16"/>
          </svg>
          <div>
            <strong>未通过审核</strong>
            <p>{{ material.rejection_reason }}</p>
          </div>
        </div>
      </div>

      <!-- 底部操作栏 -->
      <footer class="article-footer">
        <div class="footer-actions">
          <button class="primary-button" @click="handleDownload">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
              <polyline points="7 10 12 15 17 10"/>
              <line x1="12" y1="15" x2="12" y2="3"/>
            </svg>
            下载资料
          </button>
          <button class="secondary-button" @click="openReportModal">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/>
              <line x1="12" y1="9" x2="12" y2="13"/>
              <line x1="12" y1="17" x2="12.01" y2="17"/>
            </svg>
            举报
          </button>
        </div>
      </footer>
    </article>

    <!-- 举报模态框 -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="showReportModal" class="modal-overlay" @click="closeReportModal">
          <div class="modal-container" @click.stop>
            <div class="modal-header">
              <h2 class="modal-title">举报资料</h2>
              <button class="modal-close" @click="closeReportModal">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="18" y1="6" x2="6" y2="18"/>
                  <line x1="6" y1="6" x2="18" y2="18"/>
                </svg>
              </button>
            </div>

            <div class="modal-body">
              <p class="modal-description">请说明您举报此资料的原因：</p>
              <textarea
                v-model="reportReason"
                class="modal-textarea"
                placeholder="请详细描述举报原因，例如：资料内容错误、版权问题、不当内容等..."
                rows="5"
                @keydown.ctrl.enter="submitReport"
              ></textarea>
            </div>

            <div class="modal-footer">
              <button class="modal-button secondary" @click="closeReportModal">
                取消
              </button>
              <button class="modal-button primary" @click="submitReport">
                提交举报
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- 删除确认对话框 -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="deleteDialogVisible" class="modal-overlay" @click.self="deleteDialogVisible = false">
          <div class="delete-modal-container" @click.stop>
            <div class="modal-header">
              <h2 class="modal-title">确认删除</h2>
              <button class="modal-close" @click="deleteDialogVisible = false">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="18" y1="6" x2="6" y2="18"/>
                  <line x1="6" y1="6" x2="18" y2="18"/>
                </svg>
              </button>
            </div>

            <div class="modal-body">
              <div class="delete-confirmation">
                <div class="warning-icon">
                  <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <path d="M21.73 18l-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"></path>
                    <line x1="12" y1="9" x2="12" y2="13"></line>
                    <line x1="12" y1="17" x2="12.01" y2="17"></line>
                  </svg>
                </div>
                <p class="delete-message">
                  确定要删除这个资料吗？
                </p>
                <p class="delete-warning">此操作不可恢复，该资料将被永久删除。</p>
              </div>
            </div>

            <div class="modal-footer">
              <button class="modal-button secondary" @click="deleteDialogVisible = false">
                取消
              </button>
              <button class="modal-button danger" @click="confirmDelete">
                删除
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<style scoped lang="scss">
.substack-detail {
  max-width: 900px;
  margin: 0 auto;
  padding: 20px;
}

.detail-article {
  background: #ffffff;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

// 顶部导航
.article-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #f2f2f2;

  .back-button {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 16px;
    border: none;
    background: transparent;
    color: #6b6b6b;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    border-radius: 6px;
    transition: all 0.2s;

    &:hover {
      background: #f7f7f7;
      color: #1a1a1a;
    }

    svg {
      transition: transform 0.2s;
    }

    &:hover svg {
      transform: translateX(-2px);
    }
  }

  .header-actions {
    display: flex;
    gap: 8px;
  }

  .action-button {
    width: 40px;
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 1px solid #e3e3e3;
    background: #ffffff;
    border-radius: 6px;
    cursor: pointer;
    color: #6b6b6b;
    transition: all 0.2s;

    &:hover {
      background: #f7f7f7;
      border-color: #d0d0d0;
      color: #1a1a1a;
    }

    &.active {
      background: #fff8e1;
      border-color: #ffd54f;
      color: #ff8f00;
    }
  }

  .dropdown-wrapper {
    position: relative;
  }

  .dropdown-menu {
    position: absolute;
    top: 100%;
    right: 0;
    margin-top: 8px;
    background: #ffffff;
    border: 1px solid #e3e3e3;
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    min-width: 120px;
    opacity: 0;
    visibility: hidden;
    transform: translateY(-8px);
    transition: all 0.2s;
    z-index: 10;

    &.show {
      opacity: 1;
      visibility: visible;
      transform: translateY(0);
    }

    button {
      display: block;
      width: 100%;
      padding: 12px 16px;
      border: none;
      background: transparent;
      text-align: left;
      font-size: 14px;
      color: #1a1a1a;
      cursor: pointer;
      transition: background 0.15s;

      &:first-child {
        border-radius: 8px 8px 0 0;
      }

      &:last-child {
        border-radius: 0 0 8px 8px;
      }

      &:hover {
        background: #f7f7f7;
      }

      &.danger {
        color: #dc2626;

        &:hover {
          background: #fef2f2;
        }
      }
    }
  }
}

// 文章主体
.article-content {
  padding: 40px 48px;
  max-width: 720px;
  margin: 0 auto;

  @media (max-width: 768px) {
    padding: 32px 24px;
  }
}

.meta-tags {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
}

.category-tag,
.status-tag {
  display: inline-block;
  padding: 6px 14px;
  font-size: 13px;
  font-weight: 500;
  border-radius: 20px;
  letter-spacing: 0.02em;
}

.category-tag {
  background: #f7f7f7;
  color: #6b6b6b;
}

.status-tag {
  &.pending {
    background: #fff8e1;
    color: #f57c00;
  }

  &.approved {
    background: #e8f5e9;
    color: #43a047;
  }

  &.rejected {
    background: #ffebee;
    color: #e53935;
  }

  &.deleted {
    background: #f5f5f5;
    color: #9e9e9e;
  }
}

.article-title {
  font-size: 40px;
  font-weight: 700;
  line-height: 1.2;
  color: #1a1a1a;
  margin: 0 0 32px;
  letter-spacing: -0.02em;

  @media (max-width: 768px) {
    font-size: 32px;
    margin-bottom: 24px;
  }
}

.author-info {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 40px;
  padding-bottom: 32px;
  border-bottom: 1px solid #f2f2f2;
}

.author-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #ffffff;
  font-size: 18px;
  font-weight: 600;
  flex-shrink: 0;
}

.author-details {
  flex: 1;
}

.author-name {
  font-size: 16px;
  font-weight: 600;
  color: #1a1a1a;
  margin-bottom: 4px;
}

.post-meta {
  font-size: 14px;
  color: #6b6b6b;

  .separator {
    margin: 0 8px;
  }
}

.article-body {
  margin: 32px 0;
  font-size: 18px;
  line-height: 1.8;
  color: #1a1a1a;

  p {
    margin: 0;
    white-space: pre-wrap;
  }

  @media (max-width: 768px) {
    font-size: 16px;
  }
}

.file-info-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: #f7f7f7;
  border-radius: 8px;
  margin: 32px 0;

  .file-icon {
    font-size: 32px;
    flex-shrink: 0;
  }

  .file-details {
    flex: 1;
    min-width: 0;
  }

  .file-name {
    font-size: 15px;
    font-weight: 500;
    color: #1a1a1a;
    margin-bottom: 4px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .file-meta {
    font-size: 13px;
    color: #6b6b6b;
  }
}

.stats-grid {
  display: flex;
  gap: 24px;
  margin: 32px 0;
  padding: 24px 0;
  border-top: 1px solid #f2f2f2;
  border-bottom: 1px solid #f2f2f2;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #6b6b6b;

  svg {
    flex-shrink: 0;
  }
}

.tags-section {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin: 32px 0;
}

.tag-pill {
  padding: 8px 16px;
  background: #f7f7f7;
  color: #6b6b6b;
  font-size: 14px;
  border-radius: 20px;
  transition: all 0.2s;
  cursor: pointer;

  &:hover {
    background: #e8e8e8;
    color: #1a1a1a;
  }
}

.rejection-notice {
  display: flex;
  gap: 16px;
  padding: 20px;
  background: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: 8px;
  margin: 32px 0;

  svg {
    flex-shrink: 0;
    color: #dc2626;
    margin-top: 2px;
  }

  div {
    flex: 1;
  }

  strong {
    display: block;
    font-size: 16px;
    font-weight: 600;
    color: #dc2626;
    margin-bottom: 8px;
  }

  p {
    font-size: 14px;
    color: #991b1b;
    margin: 0;
    line-height: 1.6;
  }
}

// 底部操作栏
.article-footer {
  padding: 32px 48px;
  border-top: 1px solid #f2f2f2;
  background: #ffffff;

  @media (max-width: 768px) {
    padding: 24px;
  }
}

.footer-actions {
  display: flex;
  gap: 16px;
  justify-content: center;
  max-width: 720px;
  margin: 0 auto;
}

.primary-button,
.secondary-button {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 14px 28px;
  font-size: 15px;
  font-weight: 600;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  border: none;

  svg {
    flex-shrink: 0;
  }
}

.primary-button {
  background: #1a1a1a;
  color: #ffffff;

  &:hover {
    background: #333333;
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  }

  &:active {
    transform: translateY(0);
  }
}

.secondary-button {
  background: #ffffff;
  color: #1a1a1a;
  border: 1px solid #e3e3e3;

  &:hover {
    background: #f7f7f7;
    border-color: #d0d0d0;
  }
}

// 模态框样式
.modal-overlay {
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
  padding: 20px;
  backdrop-filter: blur(4px);
}

.modal-container {
  background: #ffffff;
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
  max-width: 520px;
  width: 100%;
  animation: modalSlideIn 0.3s ease-out;
}

@keyframes modalSlideIn {
  from {
    opacity: 0;
    transform: translateY(-20px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px 28px 20px;
  border-bottom: 1px solid #f2f2f2;
}

.modal-title {
  font-size: 20px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0;
  letter-spacing: -0.01em;
}

.modal-close {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: transparent;
  color: #6b6b6b;
  cursor: pointer;
  border-radius: 6px;
  transition: all 0.2s;

  &:hover {
    background: #f7f7f7;
    color: #1a1a1a;
  }

  svg {
    flex-shrink: 0;
  }
}

.modal-body {
  padding: 24px 28px;
}

.modal-description {
  font-size: 15px;
  color: #6b6b6b;
  margin: 0 0 16px 0;
  line-height: 1.5;
}

.modal-textarea {
  width: 100%;
  padding: 12px 16px;
  font-size: 15px;
  line-height: 1.6;
  color: #1a1a1a;
  background: #fafafa;
  border: 1px solid #e3e3e3;
  border-radius: 8px;
  resize: vertical;
  font-family: inherit;
  transition: all 0.2s;

  &::placeholder {
    color: #999;
  }

  &:focus {
    outline: none;
    background: #ffffff;
    border-color: #1a1a1a;
    box-shadow: 0 0 0 3px rgba(26, 26, 26, 0.05);
  }
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 28px 24px;
}

.modal-button {
  padding: 10px 20px;
  font-size: 15px;
  font-weight: 500;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  border: none;

  &.secondary {
    background: #ffffff;
    color: #6b6b6b;
    border: 1px solid #e3e3e3;

    &:hover {
      background: #f7f7f7;
      border-color: #d0d0d0;
      color: #1a1a1a;
    }
  }

  &.primary {
    background: #1a1a1a;
    color: #ffffff;

    &:hover {
      background: #333333;
      transform: translateY(-1px);
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    }

    &:active {
      transform: translateY(0);
    }
  }

  &.danger {
    background: #dc2626;
    color: #ffffff;

    &:hover {
      background: #b91c1c;
      transform: translateY(-1px);
      box-shadow: 0 4px 12px rgba(220, 38, 38, 0.3);
    }

    &:active {
      transform: translateY(0);
    }
  }
}

.delete-modal-container {
  background: #ffffff;
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
  max-width: 420px;
  width: 100%;
  animation: modalSlideIn 0.3s ease-out;
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

  svg {
    width: 32px;
    height: 32px;
  }
}

.delete-message {
  font-size: 16px;
  font-weight: 500;
  color: #1a1a1a;
  margin: 0 0 12px 0;
  line-height: 1.5;
}

.delete-warning {
  font-size: 14px;
  color: #6b6b6b;
  margin: 0;
  line-height: 1.5;
}

// 模态框过渡动画
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;

  .modal-container {
    transition: all 0.3s ease;
  }
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;

  .modal-container {
    opacity: 0;
    transform: translateY(-20px) scale(0.95);
  }
}
</style>
