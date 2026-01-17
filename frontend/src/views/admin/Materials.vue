<template>
  <div class="materials-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1>资料审核</h1>
      <div class="header-actions">
        <span v-if="pendingCount > 0" class="pending-badge">{{ pendingCount }} 待审核</span>
        <button class="refresh-btn" @click="handleRefresh">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 12a9 9 0 0 0-9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"></path>
            <path d="M3 3v5h5"></path>
            <path d="M3 12a9 9 0 0 0 9 9 9.75 9.75 0 0 0 6.74-2.74L21 16"></path>
            <path d="M16 21h5v-5"></path>
          </svg>
          刷新
        </button>
      </div>
    </div>

    <!-- 标签切换 -->
    <div class="tabs">
      <button
        :class="['tab-btn', { active: activeTab === 'pending' }]"
        @click="switchTab('pending')"
      >
        待审核
        <span v-if="pendingCount > 0" class="tab-count">{{ pendingCount }}</span>
      </button>
      <button
        :class="['tab-btn', { active: activeTab === 'reviewed' }]"
        @click="switchTab('reviewed')"
      >
        已审核
      </button>
    </div>

    <!-- 资料列表 -->
    <div v-loading="loading" class="materials-list">
      <div v-for="material in materials" :key="material.id" class="material-card">
        <div class="material-icon">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
            <polyline points="14 2 14 8 20 8"></polyline>
          </svg>
        </div>

        <div class="material-info">
          <div class="material-title">{{ material.title }}</div>
          <div class="material-meta">
            <span class="category-tag">{{ getCategoryText(material.category) }}</span>
            <span class="meta-separator">·</span>
            <span class="course-name">{{ material.course_name }}</span>
          </div>
          <div class="material-details">
            <span>{{ material.uploader?.real_name || material.uploader?.username }}</span>
            <span class="separator">·</span>
            <span>{{ format.fileSize(material.file_size) }}</span>
            <span class="separator">·</span>
            <span>{{ format.datetime(material.created_at) }}</span>
          </div>
        </div>

        <div class="material-status">
          <span :class="['status-badge', material.status]">
            {{ getStatusText(material.status) }}
          </span>
        </div>

        <div class="material-actions">
          <button class="action-btn primary" @click="handleReview(material)">
            审核资料
          </button>
          <button class="action-btn" @click="handleViewDetail(material)">
            查看详情
          </button>
        </div>
      </div>

      <div v-if="!loading && materials.length === 0" class="empty-state">
        <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
          <polyline points="14 2 14 8 20 8"></polyline>
        </svg>
        <p>{{ activeTab === 'pending' ? '暂无待审核资料' : '暂无已审核资料' }}</p>
      </div>
    </div>

    <!-- 分页 -->
    <div v-if="total > 0" class="pagination">
      <span class="pagination-info">
        显示 {{ (page - 1) * size + 1 }}-{{ Math.min(page * size, total) }} / 共 {{ total }} 条
      </span>
      <div class="pagination-controls">
        <button class="page-btn" :disabled="page === 1" @click="handlePageChange(page - 1)">
          上一页
        </button>
        <button class="page-btn" :disabled="page * size >= total" @click="handlePageChange(page + 1)">
          下一页
        </button>
      </div>
    </div>

    <!-- 审核对话框 -->
    <ReviewMaterial
      v-model="reviewDialogVisible"
      :material="currentMaterial"
      @success="handleReviewSuccess"
    />

    <!-- 详情对话框 -->
    <div v-if="detailDialogVisible" class="dialog-overlay" @click.self="detailDialogVisible = false">
      <div class="dialog">
        <div class="dialog-header">
          <h3>资料详情</h3>
          <button class="close-btn" @click="detailDialogVisible = false">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 6 6 18"></path>
              <path d="m6 6 12 12"></path>
            </svg>
          </button>
        </div>

        <div v-if="currentMaterial" class="dialog-body">
          <div class="detail-group">
            <label>资料名称</label>
            <div class="detail-value">{{ currentMaterial.title }}</div>
          </div>

          <div class="detail-group">
            <label>资料描述</label>
            <div class="detail-value">{{ currentMaterial.description || '-' }}</div>
          </div>

          <div class="detail-row">
            <div class="detail-group">
              <label>分类</label>
              <span class="detail-tag">{{ getCategoryText(currentMaterial.category) }}</span>
            </div>
            <div class="detail-group">
              <label>课程名称</label>
              <div class="detail-value">{{ currentMaterial.course_name }}</div>
            </div>
          </div>

          <div class="detail-row">
            <div class="detail-group">
              <label>上传者</label>
              <div class="detail-value">
                {{ currentMaterial.uploader?.real_name || currentMaterial.uploader?.username }}
              </div>
            </div>
            <div class="detail-group">
              <label>文件大小</label>
              <div class="detail-value">{{ format.fileSize(currentMaterial.file_size) }}</div>
            </div>
          </div>

          <div class="detail-group">
            <label>文件名</label>
            <div class="detail-value">{{ currentMaterial.file_name }}</div>
          </div>

          <div class="detail-group">
            <label>标签</label>
            <div class="tags-list">
              <span
                v-for="tag in currentMaterial.tags"
                :key="tag"
                class="tag-item"
              >
                {{ tag }}
              </span>
              <span v-if="!currentMaterial.tags || currentMaterial.tags.length === 0">-</span>
            </div>
          </div>

          <div class="detail-row">
            <div class="detail-group">
              <label>下载次数</label>
              <div class="detail-value">{{ currentMaterial.download_count }}</div>
            </div>
            <div class="detail-group">
              <label>收藏次数</label>
              <div class="detail-value">{{ currentMaterial.favorite_count }}</div>
            </div>
          </div>

          <div class="detail-group">
            <label>上传时间</label>
            <div class="detail-value">{{ format.datetime(currentMaterial.created_at) }}</div>
          </div>
        </div>

        <div class="dialog-footer">
          <button class="btn btn-secondary" @click="detailDialogVisible = false">关闭</button>
          <button class="btn btn-primary" @click="handleReviewFromDetail">审核资料</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useReviewStore } from '@/stores/review'
import { useMaterialStore } from '@/stores/material'
import { useMaterialCategoryStore } from '@/stores/materialCategory'
import { format } from '@/utils/format'
import ReviewMaterial from '@/components/ReviewMaterial.vue'
import type { Material } from '@/types'

const reviewStore = useReviewStore()
const materialStore = useMaterialStore()
const materialCategoryStore = useMaterialCategoryStore()

const loading = ref(false)
const activeTab = ref('pending')
const page = ref(1)
const size = ref(20)
const pendingCount = ref(0)

const reviewDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const currentMaterial = ref<Material | null>(null)

const materials = computed(() => {
  return activeTab.value === 'pending'
    ? reviewStore.pendingMaterials
    : materialStore.materials
})

const total = computed(() => {
  return activeTab.value === 'pending' ? reviewStore.total : materialStore.total
})

const loadPendingMaterials = async () => {
  loading.value = true
  try {
    await reviewStore.fetchPendingMaterials({
      page: page.value,
      size: size.value
    })
    pendingCount.value = reviewStore.pendingMaterials.length
  } finally {
    loading.value = false
  }
}

const loadReviewedMaterials = async () => {
  loading.value = true
  try {
    await materialStore.fetchReviewedMaterials({
      page: page.value,
      size: size.value,
      reviewed_only: true
    })
  } finally {
    loading.value = false
  }
}

const handleRefresh = () => {
  if (activeTab.value === 'pending') {
    loadPendingMaterials()
  } else {
    loadReviewedMaterials()
  }
}

const switchTab = (tab: string) => {
  activeTab.value = tab
  page.value = 1
  handleRefresh()
}

const handlePageChange = (newPage: number) => {
  page.value = newPage
  handleRefresh()
}

const handleReview = (material: Material) => {
  currentMaterial.value = material
  reviewDialogVisible.value = true
}

const handleViewDetail = (material: Material) => {
  currentMaterial.value = material
  detailDialogVisible.value = true
}

const handleReviewFromDetail = () => {
  detailDialogVisible.value = false
  reviewDialogVisible.value = true
}

const handleReviewSuccess = () => {
  loadPendingMaterials()
}

const getCategoryText = (category: string) => {
  return materialCategoryStore.getCategoryName(category)
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    pending: '待审核',
    approved: '已通过',
    rejected: '已拒绝'
  }
  return statusMap[status] || status
}

onMounted(async () => {
  // 加载分类数据
  await materialCategoryStore.fetchActiveCategories()
  loadPendingMaterials()
})
</script>

<style scoped lang="scss">
.materials-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;

  h1 {
    font-size: 28px;
    font-weight: 700;
    color: #111827;
    margin: 0;
  }

  .header-actions {
    display: flex;
    align-items: center;
    gap: 12px;
  }
}

.pending-badge {
  padding: 6px 12px;
  background: #fef3c7;
  color: #b45309;
  font-size: 13px;
  font-weight: 600;
  border-radius: 6px;
}

.refresh-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  background: #ffffff;
  color: #111827;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;

  &:hover {
    background: #f9fafb;
    border-color: #d1d5db;
  }
}

.tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 24px;
  border-bottom: 1px solid #e5e7eb;
  padding-bottom: 0;
}

.tab-btn {
  position: relative;
  padding: 12px 16px;
  border: none;
  background: none;
  color: #6b7280;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  border-bottom: 2px solid transparent;
  margin-bottom: -1px;

  &:hover {
    color: #111827;
  }

  &.active {
    color: #111827;
    border-bottom-color: #111827;
  }

  .tab-count {
    display: inline-block;
    margin-left: 6px;
    padding: 2px 6px;
    background: #f3f4f6;
    border-radius: 4px;
    font-size: 12px;
    font-weight: 600;
  }
}

.materials-list {
  min-height: 300px;
}

.material-card {
  display: grid;
  grid-template-columns: auto 1fr auto auto;
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

    .material-status {
      display: none;
    }
  }
}

.material-icon {
  width: 48px;
  height: 48px;
  border-radius: 10px;
  background: #f3f4f6;
  color: #6b7280;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.material-info {
  min-width: 0;
}

.material-title {
  font-size: 15px;
  font-weight: 600;
  color: #111827;
  margin-bottom: 6px;
}

.material-meta {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 6px;
}

.category-tag {
  padding: 2px 8px;
  background: #e5e7eb;
  color: #6b7280;
  font-size: 12px;
  font-weight: 500;
  border-radius: 4px;
}

.meta-separator {
  color: #d1d5db;
}

.course-name {
  font-size: 13px;
  color: #6b7280;
}

.material-details {
  font-size: 13px;
  color: #9ca3af;
  display: flex;
  align-items: center;
  gap: 6px;

  .separator {
    color: #d1d5db;
  }
}

.status-badge {
  padding: 4px 10px;
  font-size: 12px;
  font-weight: 500;
  border-radius: 6px;
  white-space: nowrap;

  &.pending {
    background: #fef3c7;
    color: #b45309;
  }

  &.approved {
    background: #dcfce7;
    color: #15803d;
  }

  &.rejected {
    background: #fee2e2;
    color: #b91c1c;
  }
}

.material-actions {
  display: flex;
  gap: 8px;

  @media (max-width: 768px) {
    grid-column: 2;
    justify-content: flex-end;
  }
}

.action-btn {
  padding: 8px 14px;
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

  &.primary {
    background: #111827;
    color: #ffffff;

    &:hover {
      background: #000000;
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
  max-width: 700px;
  max-height: 90vh;
  overflow-y: auto;
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

.detail-group {
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

.detail-value {
  font-size: 14px;
  color: #111827;
}

.detail-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  margin-bottom: 20px;

  @media (max-width: 600px) {
    grid-template-columns: 1fr;
  }
}

.detail-tag {
  display: inline-block;
  padding: 4px 10px;
  background: #f3f4f6;
  border-radius: 6px;
  font-size: 13px;
  color: #111827;
}

.tags-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;

  .tag-item {
    padding: 4px 10px;
    background: #f3f4f6;
    border-radius: 6px;
    font-size: 13px;
    color: #111827;
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
}
</style>
