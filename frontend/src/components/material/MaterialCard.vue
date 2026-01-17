<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useMaterialStore } from '@/stores/material'
import { useMaterialCategoryStore } from '@/stores/materialCategory'
import type { Material } from '@/types'

interface Props {
  material: Material
}

const props = defineProps<Props>()
const router = useRouter()
const materialStore = useMaterialStore()
const materialCategoryStore = useMaterialCategoryStore()

// 资料状态映射
const statusConfig: Record<string, { text: string; color: string }> = {
  pending: { text: '待审核', color: '#f59e0b' },
  approved: { text: '已通过', color: '#10b981' },
  rejected: { text: '已拒绝', color: '#ef4444' },
  deleted: { text: '已删除', color: '#9ca3af' }
}

// 格式化文件大小
const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
}

// 格式化时间
const formatTime = (dateStr: string): string => {
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (days === 0) return '今天'
  if (days === 1) return '昨天'
  if (days < 7) return `${days} 天前`
  if (days < 30) return `${Math.floor(days / 7)} 周前`
  return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' })
}

// 资料分类文本(从动态配置中获取)
const categoryText = computed(() => materialCategoryStore.getCategoryName(props.material.category))

// 资料状态配置
const statusInfo = computed(() => statusConfig[props.material.status] || { text: '未知', color: '#9ca3af' })

// 点击卡片跳转到详情页
const goToDetail = () => {
  router.push(`/materials/${props.material.id}`)
}

// 收藏/取消收藏
const toggleFavorite = async (event: MouseEvent) => {
  event.stopPropagation()

  // 使用 store 中的最新状态，而不是 props
  const material = materialStore.materials.find(m => m.id === props.material.id)
  const isFavorited = material?.is_favorited ?? props.material.is_favorited

  try {
    if (isFavorited) {
      await materialStore.removeFavorite(props.material.id)
      ElMessage.success('已取消收藏')
    } else {
      const response = await materialStore.addFavorite(props.material.id)
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
</script>

<template>
  <article class="material-item" @click="goToDetail">
    <div class="material-icon">
      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
        <polyline points="14 2 14 8 20 8"></polyline>
      </svg>
    </div>

    <div class="material-info">
      <div class="material-title">{{ material.title }}</div>
      <div class="material-meta">
        <span class="category-tag">{{ categoryText }}</span>
        <span class="separator">·</span>
        <span v-if="material.course_name" class="course-name">{{ material.course_name }}</span>
      </div>
      <div class="material-details">
        <span>{{ formatFileSize(material.file_size) }}</span>
        <span class="separator">·</span>
        <span>{{ material.uploader?.real_name || material.uploader?.username || '未知' }}</span>
        <span class="separator">·</span>
        <span>{{ formatTime(material.created_at) }}</span>
      </div>
    </div>

    <div class="material-status">
      <span :class="['status-badge', material.status]">
        {{ statusInfo.text }}
      </span>
    </div>

    <div class="material-actions">
      <button
        :class="['action-btn', { active: material.is_favorited }]"
        @click="toggleFavorite"
        title="收藏"
        data-testid="favorite-button"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="16"
          height="16"
          viewBox="0 0 24 24"
          :fill="material.is_favorited ? 'currentColor' : 'none'"
          stroke="currentColor"
          stroke-width="2"
        >
          <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/>
        </svg>
        {{ material.favorite_count || 0 }}
      </button>
      <div class="stat" title="下载次数">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
          <polyline points="7 10 12 15 17 10"/>
          <line x1="12" y1="15" x2="12" y2="3"/>
        </svg>
        {{ material.download_count }}
      </div>
    </div>
  </article>
</template>

<style scoped lang="scss">
.material-item {
  display: grid;
  grid-template-columns: auto 1fr auto auto;
  gap: 16px;
  align-items: center;
  padding: 20px 0;
  background: transparent;
  border: none;
  border-bottom: 1px solid #f2f2f2;
  border-radius: 0;
  transition: all 0.2s;
  cursor: pointer;

  &:hover {
    background: #f9fafb;
  }

  &:first-child {
    border-top: 1px solid #f2f2f2;
  }

  @media (max-width: 768px) {
    grid-template-columns: auto 1fr;
    gap: 12px;

    .material-status {
      display: none;
    }

    .material-actions {
      grid-column: 2;
      justify-content: flex-end;
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

.separator {
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

  &.deleted {
    background: #f3f4f6;
    color: #6b7280;
  }
}

.material-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 14px;
  border: none;
  border-radius: 6px;
  background: #f3f4f6;
  color: #374151;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;

  svg {
    width: 16px;
    height: 16px;
  }

  &:hover {
    background: #e5e7eb;
  }

  &.active {
    color: #f59e0b;

    &:hover {
      background: #fef3c7;
    }
  }
}

.stat {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #9ca3af;
  font-size: 13px;

  svg {
    width: 16px;
    height: 16px;
  }
}

// 响应式
@media (max-width: 768px) {
  .material-title {
    font-size: 14px;
  }
}
</style>
