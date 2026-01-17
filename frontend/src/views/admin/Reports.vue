<template>
  <div class="reports-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1>举报处理</h1>
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

    <!-- 筛选标签 -->
    <div class="filter-tabs">
      <button
        :class="['filter-tab', { active: filterStatus === '' }]"
        @click="setFilterStatus('')"
      >
        全部
      </button>
      <button
        :class="['filter-tab', { active: filterStatus === 'pending' }]"
        @click="setFilterStatus('pending')"
      >
        待处理
      </button>
      <button
        :class="['filter-tab', { active: filterStatus === 'approved' }]"
        @click="setFilterStatus('approved')"
      >
        已成立
      </button>
      <button
        :class="['filter-tab', { active: filterStatus === 'rejected' }]"
        @click="setFilterStatus('rejected')"
      >
        不成立
      </button>
    </div>

    <!-- 举报列表 -->
    <div v-loading="loading" class="reports-list">
      <div v-for="report in reports" :key="report.id" class="report-card">
        <div class="report-icon">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"></path>
            <line x1="12" y1="9" x2="12" y2="13"></line>
            <line x1="12" y1="17" x2="12.01" y2="17"></line>
          </svg>
        </div>

        <div class="report-info">
          <div class="report-material">
            {{ report.material?.title }}
          </div>
          <div class="report-meta">
            <span class="reason-tag" :class="report.reason">
              {{ getReasonText(report.reason) }}
            </span>
            <span class="meta-separator">·</span>
            <span class="reporter-name">
              {{ report.reporter?.real_name || report.reporter?.username }}
            </span>
            <span class="meta-separator">·</span>
            <span class="report-time">{{ format.datetime(report.created_at) }}</span>
          </div>
          <div v-if="report.description" class="report-description">
            {{ report.description }}
          </div>
        </div>

        <div class="report-status">
          <span :class="['status-badge', report.status]">
            {{ getStatusText(report.status) }}
          </span>
        </div>

        <div class="report-actions">
          <button
            v-if="report.status === 'pending'"
            class="action-btn primary"
            @click="handleHandle(report)"
          >
            处理
          </button>
          <button class="action-btn" @click="handleViewDetail(report)">
            查看详情
          </button>
        </div>
      </div>

      <div v-if="!loading && reports.length === 0" class="empty-state">
        <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"></path>
          <line x1="12" y1="9" x2="12" y2="13"></line>
          <line x1="12" y1="17" x2="12.01" y2="17"></line>
        </svg>
        <p>暂无举报记录</p>
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

    <!-- 处理对话框 -->
    <HandleReport
      v-model="handleDialogVisible"
      :report="currentReport"
      @success="handleHandleSuccess"
    />

    <!-- 详情对话框 -->
    <div v-if="detailDialogVisible" class="dialog-overlay" @click.self="detailDialogVisible = false">
      <div class="dialog large">
        <div class="dialog-header">
          <h3>举报详情</h3>
          <button class="close-btn" @click="detailDialogVisible = false">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 6 6 18"></path>
              <path d="m6 6 12 12"></path>
            </svg>
          </button>
        </div>

        <div v-if="currentReport" class="dialog-body">
          <!-- 举报信息 -->
          <div class="detail-section">
            <h4>举报信息</h4>
            <div class="detail-group">
              <label>被举报资料</label>
              <button class="link-btn" @click="handleViewMaterial(currentReport.material)">
                {{ currentReport.material?.title }}
              </button>
            </div>

            <div class="detail-row">
              <div class="detail-group">
                <label>举报人</label>
                <div class="detail-value">
                  {{ currentReport.reporter?.real_name || currentReport.reporter?.username }}
                </div>
              </div>
              <div class="detail-group">
                <label>举报时间</label>
                <div class="detail-value">{{ format.datetime(currentReport.created_at) }}</div>
              </div>
            </div>

            <div class="detail-row">
              <div class="detail-group">
                <label>举报原因</label>
                <span :class="['reason-tag', currentReport.reason]">
                  {{ getReasonText(currentReport.reason) }}
                </span>
              </div>
              <div class="detail-group">
                <label>处理状态</label>
                <span :class="['status-badge', currentReport.status]">
                  {{ getStatusText(currentReport.status) }}
                </span>
              </div>
            </div>

            <div class="detail-group">
              <label>详细描述</label>
              <div class="detail-value">{{ currentReport.description }}</div>
            </div>

            <div v-if="currentReport.handled_at" class="detail-row">
              <div class="detail-group">
                <label>处理时间</label>
                <div class="detail-value">{{ format.datetime(currentReport.handled_at) }}</div>
              </div>
              <div v-if="currentReport.handler" class="detail-group">
                <label>处理人</label>
                <div class="detail-value">
                  {{ currentReport.handler.real_name || currentReport.handler.username }}
                </div>
              </div>
            </div>

            <div v-if="currentReport.handle_note" class="detail-group">
              <label>处理说明</label>
              <div class="detail-value">{{ currentReport.handle_note }}</div>
            </div>
          </div>

          <!-- 资料信息 -->
          <div v-if="currentReport.material" class="detail-section">
            <h4>资料信息</h4>
            <div class="detail-group">
              <label>资料名称</label>
              <div class="detail-value">{{ currentReport.material.title }}</div>
            </div>

            <div v-if="currentReport.material.description" class="detail-group">
              <label>资料描述</label>
              <div class="detail-value">{{ currentReport.material.description }}</div>
            </div>

            <div class="detail-row">
              <div class="detail-group">
                <label>分类</label>
                <span class="detail-tag">
                  {{ getCategoryText(currentReport.material.category) }}
                </span>
              </div>
              <div class="detail-group">
                <label>课程名称</label>
                <div class="detail-value">{{ currentReport.material.course_name }}</div>
              </div>
            </div>

            <div class="detail-row">
              <div class="detail-group">
                <label>上传者</label>
                <div class="detail-value">
                  {{ currentReport.material.uploader?.real_name || currentReport.material.uploader?.username }}
                </div>
              </div>
              <div class="detail-group">
                <label>上传时间</label>
                <div class="detail-value">
                  {{ format.datetime(currentReport.material.created_at) }}
                </div>
              </div>
            </div>

            <div class="detail-row">
              <div class="detail-group">
                <label>文件大小</label>
                <div class="detail-value">
                  {{ format.fileSize(currentReport.material.file_size) }}
                </div>
              </div>
              <div class="detail-group">
                <label>下载次数</label>
                <div class="detail-value">{{ currentReport.material.download_count }}</div>
              </div>
            </div>

            <div v-if="currentReport.material.tags && currentReport.material.tags.length > 0" class="detail-group">
              <label>标签</label>
              <div class="tags-list">
                <span
                  v-for="tag in currentReport.material.tags"
                  :key="tag"
                  class="tag-item"
                >
                  {{ tag }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <div class="dialog-footer">
          <button class="btn btn-secondary" @click="detailDialogVisible = false">关闭</button>
          <button
            v-if="currentReport?.status === 'pending'"
            class="btn btn-primary"
            @click="handleHandleFromDetail"
          >
            立即处理
          </button>
        </div>
      </div>
    </div>

    <!-- 资料详情对话框 -->
    <div v-if="materialDialogVisible" class="dialog-overlay" @click.self="materialDialogVisible = false">
      <div class="dialog">
        <div class="dialog-header">
          <h3>资料详情</h3>
          <button class="close-btn" @click="materialDialogVisible = false">
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

          <div v-if="currentMaterial.description" class="detail-group">
            <label>资料描述</label>
            <div class="detail-value">{{ currentMaterial.description }}</div>
          </div>

          <div class="detail-row">
            <div class="detail-group">
              <label>分类</label>
              <span class="detail-tag">
                {{ getCategoryText(currentMaterial.category) }}
              </span>
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

          <div v-if="currentMaterial.tags && currentMaterial.tags.length > 0" class="detail-group">
            <label>标签</label>
            <div class="tags-list">
              <span v-for="tag in currentMaterial.tags" :key="tag" class="tag-item">
                {{ tag }}
              </span>
            </div>
          </div>
        </div>

        <div class="dialog-footer">
          <button class="btn btn-secondary" @click="materialDialogVisible = false">关闭</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { materialApi, reportApi } from '@/api/material'
import { useMaterialCategoryStore } from '@/stores/materialCategory'
import { format } from '@/utils/format'
import HandleReport from '@/components/HandleReport.vue'
import type { Report, Material, ReportListParams } from '@/types'

const materialCategoryStore = useMaterialCategoryStore()

const loading = ref(false)
const filterStatus = ref('')
const page = ref(1)
const size = ref(20)
const total = ref(0)

const reports = ref<Report[]>([])

const handleDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const materialDialogVisible = ref(false)
const currentReport = ref<Report | null>(null)
const currentMaterial = ref<Material | null>(null)

const loadReports = async () => {
  loading.value = true
  try {
    const params: ReportListParams = {
      page: page.value,
      size: size.value,
      status: filterStatus.value || undefined
    }

    const response = await reportApi.getReports(params)

    if (response.code === 0 && response.data) {
      reports.value = response.data.list || []
      total.value = response.data.total
    }
  } finally {
    loading.value = false
  }
}

const handleRefresh = () => {
  loadReports()
}

const setFilterStatus = (status: string) => {
  filterStatus.value = status
  page.value = 1
  loadReports()
}

const handlePageChange = (newPage: number) => {
  page.value = newPage
  loadReports()
}

const handleHandle = (report: Report) => {
  currentReport.value = report
  handleDialogVisible.value = true
}

const handleViewDetail = (report: Report) => {
  currentReport.value = report
  detailDialogVisible.value = true
}

const handleHandleFromDetail = () => {
  detailDialogVisible.value = false
  handleDialogVisible.value = true
}

const handleHandleSuccess = () => {
  loadReports()
}

const handleViewMaterial = async (material?: Material) => {
  if (!material) return

  try {
    const response = await materialApi.getMaterial(material.id)
    if (response.code === 0 && response.data) {
      currentMaterial.value = response.data
      materialDialogVisible.value = true
    }
  } catch (error: any) {
    ElMessage.error(error.message || '获取资料详情失败')
  }
}

const getStatusText = (status: string) => {
  const textMap: Record<string, string> = {
    pending: '待处理',
    approved: '已成立',
    rejected: '不成立'
  }
  return textMap[status] || status
}

const getReasonText = (reason: string) => {
  const textMap: Record<string, string> = {
    inappropriate: '内容不当',
    copyright: '版权问题',
    wrong_category: '分类错误',
    low_quality: '质量差',
    other: '其他'
  }
  return textMap[reason] || reason
}

const getCategoryText = (category: string) => {
  return materialCategoryStore.getCategoryName(category)
}

onMounted(async () => {
  // 加载分类数据
  await materialCategoryStore.fetchActiveCategories()
  loadReports()
})
</script>

<style scoped lang="scss">
.reports-page {
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

.filter-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 24px;
  border-bottom: 1px solid #e5e7eb;
  padding-bottom: 0;
}

.filter-tab {
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
}

.reports-list {
  min-height: 300px;
}

.report-card {
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

    .report-status {
      display: none;
    }
  }
}

.report-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  background: #fef3c7;
  color: #b45309;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.report-info {
  min-width: 0;
}

.report-material {
  font-size: 15px;
  font-weight: 600;
  color: #111827;
  margin-bottom: 6px;
}

.report-meta {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 6px;
  flex-wrap: wrap;
}

.reason-tag {
  padding: 2px 8px;
  font-size: 12px;
  font-weight: 500;
  border-radius: 4px;

  &.inappropriate {
    background: #fee2e2;
    color: #b91c1c;
  }

  &.copyright {
    background: #fef3c7;
    color: #b45309;
  }

  &.wrong_category,
  &.low_quality {
    background: #dbeafe;
    color: #1d4ed8;
  }

  &.other {
    background: #e5e7eb;
    color: #6b7280;
  }
}

.meta-separator {
  color: #d1d5db;
}

.reporter-name {
  font-size: 13px;
  color: #6b7280;
}

.report-time {
  font-size: 13px;
  color: #9ca3af;
}

.report-description {
  font-size: 13px;
  color: #6b7280;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-height: 1.5;
  max-height: 2.7em;
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

.report-actions {
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

  &.large {
    max-width: 800px;
  }
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

.detail-section {
  margin-bottom: 24px;
  padding-bottom: 24px;
  border-bottom: 1px solid #e5e7eb;

  &:last-child {
    margin-bottom: 0;
    padding-bottom: 0;
    border-bottom: none;
  }

  h4 {
    margin: 0 0 16px 0;
    font-size: 16px;
    font-weight: 600;
    color: #111827;
  }
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

.link-btn {
  background: none;
  border: none;
  padding: 0;
  color: #2563eb;
  font-size: 14px;
  cursor: pointer;
  text-decoration: underline;
  text-align: left;

  &:hover {
    color: #1d4ed8;
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
