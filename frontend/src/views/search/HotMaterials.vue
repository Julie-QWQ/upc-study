<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useSearchStore } from '@/stores/search'
import { useMaterialStore } from '@/stores/material'
import MaterialCard from '@/components/material/MaterialCard.vue'
import type { MaterialListParams } from '@/types'

const router = useRouter()
const searchStore = useSearchStore()
const materialStore = useMaterialStore()

const currentTab = ref<'hot' | 'recommended'>('hot')
const loading = ref(false)

// 热门资料筛选条件
const hotFilters = ref<MaterialListParams>({
  page: 1,
  size: 20,
  sort_by: 'download_count',
  order: 'desc'
})

// 推荐资料筛选条件
const recommendedFilters = ref<MaterialListParams>({
  page: 1,
  size: 20,
  sort_by: 'favorite_count',
  order: 'desc'
})

// 当前使用的筛选条件
const currentFilters = computed(() => {
  return currentTab.value === 'hot' ? hotFilters.value : recommendedFilters.value
})

// 获取热门资料
const fetchHotMaterials = async () => {
  loading.value = true
  try {
    await materialStore.fetchMaterials(hotFilters.value)
  } finally {
    loading.value = false
  }
}

// 获取推荐资料
const fetchRecommendedMaterials = async () => {
  loading.value = true
  try {
    await materialStore.fetchMaterials(recommendedFilters.value)
  } finally {
    loading.value = false
  }
}

// 处理标签页切换
const handleTabChange = async () => {
  if (currentTab.value === 'hot') {
    await fetchHotMaterials()
  } else {
    await fetchRecommendedMaterials()
  }
}

// 处理页码变化
const handlePageChange = async (page: number) => {
  currentFilters.value.page = page
  if (currentTab.value === 'hot') {
    await fetchHotMaterials()
  } else {
    await fetchRecommendedMaterials()
  }
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

// 跳转到资料详情
const goToDetail = (materialId: number) => {
  router.push(`/materials/${materialId}`)
}

// 计算排名数字(考虑分页)
const getRankNumber = (index: number) => {
  const pageSize = currentFilters.value.size || 20
  const currentPage = currentFilters.value.page || 1
  return (currentPage - 1) * pageSize + index + 1
}

// 获取排名样式类
const getRankClass = (index: number) => {
  const rank = getRankNumber(index)
  if (rank === 1) return 'first'
  if (rank === 2) return 'second'
  if (rank === 3) return 'third'
  return 'normal'
}

onMounted(async () => {
  await fetchHotMaterials()
})
</script>

<template>
  <div class="hot-materials-page">
    <div class="page-container">
      <main class="content-stream">
        <!-- 页面标题 -->
        <header class="stream-header">
          <div class="header-content">
            <svg class="trophy-icon" xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M6 9H4.5a2.5 2.5 0 0 1 0-5H6"/>
              <path d="M18 9h1.5a2.5 2.5 0 0 0 0-5H18"/>
              <path d="M4 22h16"/>
              <path d="M10 14.66V17c0 .55-.47.98-.97 1.21C7.85 18.75 7 20.24 7 22"/>
              <path d="M14 14.66V17c0 .55.47.98.97 1.21C16.15 18.75 17 20.24 17 22"/>
              <path d="M18 2H6v7a6 6 0 0 0 12 0V2Z"/>
            </svg>
            <h1 class="page-title">热门资料</h1>
          </div>
        </header>

        <!-- 标签页切换 -->
        <div class="tabs-section">
          <button
            :class="['tab-button', { active: currentTab === 'hot' }]"
            @click="currentTab = 'hot'; handleTabChange()"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M6 9H4.5a2.5 2.5 0 0 1 0-5H6"/>
              <path d="M18 9h1.5a2.5 2.5 0 0 0 0-5H18"/>
              <path d="M4 22h16"/>
              <path d="M10 14.66V17c0 .55-.47.98-.97 1.21C7.85 18.75 7 20.24 7 22"/>
              <path d="M14 14.66V17c0 .55.47.98.97 1.21C16.15 18.75 17 20.24 17 22"/>
              <path d="M18 2H6v7a6 6 0 0 0 12 0V2Z"/>
            </svg>
            热门排行
          </button>
          <button
            :class="['tab-button', { active: currentTab === 'recommended' }]"
            @click="currentTab = 'recommended'; handleTabChange()"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/>
            </svg>
            为您推荐
          </button>
        </div>

        <!-- 加载状态 -->
        <div v-if="loading" class="loading-state">
          <div v-for="i in 6" :key="i" class="skeleton-card"></div>
        </div>

        <!-- 空状态 -->
        <div v-else-if="materialStore.materials.length === 0" class="empty-state">
          <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M6 9H4.5a2.5 2.5 0 0 1 0-5H6"/>
            <path d="M18 9h1.5a2.5 2.5 0 0 0 0-5H18"/>
            <path d="M4 22h16"/>
            <path d="M10 14.66V17c0 .55-.47.98-.97 1.21C7.85 18.75 7 20.24 7 22"/>
            <path d="M14 14.66V17c0 .55.47.98.97 1.21C16.15 18.75 17 20.24 17 22"/>
            <path d="M18 2H6v7a6 6 0 0 0 12 0V2Z"/>
          </svg>
          <h2>暂无热门资料</h2>
          <p>{{ currentTab === 'hot' ? '还没有热门资料' : '暂无推荐内容' }}</p>
        </div>

        <!-- 资料列表 -->
        <div v-else class="materials-list">
          <div
            v-for="(material, index) in materialStore.materials"
            :key="material.id"
            class="ranked-material-item"
          >
            <div class="rank-number" :class="`rank-${getRankClass(index)}`" @click="goToDetail(material.id)">
              {{ getRankNumber(index) }}
            </div>
            <div class="material-card-wrapper" @click="goToDetail(material.id)">
              <MaterialCard
                :material="material"
                @click.stop
              />
            </div>
          </div>
        </div>

        <!-- 分页 -->
        <div v-if="materialStore.total > 0" class="pagination-section">
          <el-pagination
            :current-page="materialStore.page"
            :page-size="materialStore.size"
            :total="materialStore.total"
            layout="prev, pager, next"
            @current-change="handlePageChange"
          />
        </div>
      </main>
    </div>
  </div>
</template>

<style scoped lang="scss">
.hot-materials-page {
  width: 100%;
  background: #ffffff;
  min-height: 100vh;
}

.page-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 32px 24px;
}

.content-stream {
  max-width: 100%;
}

.stream-header {
  margin-bottom: 24px;

  .header-content {
    display: flex;
    align-items: center;
    gap: 12px;

    .trophy-icon {
      flex-shrink: 0;
      color: #f59e0b;
    }

    .page-title {
      font-size: 28px;
      font-weight: 700;
      color: #111827;
      margin: 0;
    }
  }
}

.tabs-section {
  display: flex;
  gap: 8px;
  margin-bottom: 24px;

  .tab-button {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 20px;
    background: #ffffff;
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    font-size: 14px;
    color: #6b7280;
    cursor: pointer;
    transition: all 0.2s;

    svg {
      width: 18px;
      height: 18px;
    }

    &:hover {
      border-color: #111827;
      color: #111827;
    }

    &.active {
      background: #111827;
      border-color: #111827;
      color: #ffffff;
      font-weight: 500;
    }
  }
}

.materials-list {
  display: flex;
  flex-direction: column;
  gap: 0;

  .ranked-material-item {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 16px 0;
    border-bottom: 1px solid #f2f2f2;

    &:last-child {
      border-bottom: none;
    }

    .rank-number {
      flex-shrink: 0;
      width: 48px;
      height: 48px;
      display: flex;
      align-items: center;
      justify-content: center;
      border-radius: 10px;
      font-size: 20px;
      font-weight: 700;
      color: #6b7280;
      background: #f3f4f6;
      cursor: pointer;

      &.rank-first {
        background: linear-gradient(135deg, #ffd700 0%, #ffed4e 100%);
        color: #92400e;
        box-shadow: 0 2px 8px rgba(255, 215, 0, 0.3);
        font-size: 24px;
      }

      &.rank-second {
        background: linear-gradient(135deg, #c0c0c0 0%, #e5e7eb 100%);
        color: #4b5563;
        box-shadow: 0 2px 8px rgba(192, 192, 192, 0.3);
        font-size: 22px;
      }

      &.rank-third {
        background: linear-gradient(135deg, #cd7f32 0%, #daa06d 100%);
        color: #78350f;
        box-shadow: 0 2px 8px rgba(205, 127, 50, 0.3);
        font-size: 22px;
      }

      &.rank-normal {
        background: #f3f4f6;
        color: #9ca3af;
      }
    }

    .material-card-wrapper {
      flex: 1;
      min-width: 0;
      cursor: pointer;
    }
  }
}

.loading-state {
  display: flex;
  flex-direction: column;
  gap: 0;

  .skeleton-card {
    height: 200px;
    background: linear-gradient(90deg, #f5f5f5 25%, #e8e8e8 50%, #f5f5f5 75%);
    background-size: 200% 100%;
    animation: shimmer 1.5s infinite;
    border-bottom: 1px solid #f2f2f2;

    &:last-child {
      border-bottom: none;
    }
  }
}

@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

.empty-state {
  text-align: center;
  padding: 80px 40px;
  background: #f9fafb;
  border-radius: 12px;

  svg {
    color: #d1d5db;
    margin-bottom: 24px;
  }

  h2 {
    font-size: 22px;
    font-weight: 600;
    color: #111827;
    margin: 0 0 8px 0;
  }

  p {
    font-size: 15px;
    color: #9ca3af;
    margin: 0;
  }
}

.pagination-section {
  margin-top: 40px;
  padding-top: 24px;
  border-top: 1px solid #e5e7eb;
  display: flex;
  justify-content: center;
}

:deep(.el-pagination) {
  .btn-prev,
  .btn-next,
  .el-pager li {
    background: none;
    border: none;
    font-weight: 500;
    color: #6b7280;

    &:hover {
      color: #111827;
    }

    &.active {
      color: #111827;
      font-weight: 600;
    }
  }
}

@media (max-width: 768px) {
  .page-container {
    padding: 24px 16px;
  }

  .stream-header {
    .header-content {
      .page-title {
        font-size: 24px;
      }
    }
  }

  .tabs-section {
    .tab-button {
      flex: 1;
      justify-content: center;
    }
  }

  .materials-list {
    .ranked-material-item {
      gap: 12px;
      padding: 12px 0;

      .rank-number {
        width: 36px;
        height: 36px;
        font-size: 16px;

        &.rank-first {
          font-size: 20px;
        }

        &.rank-second,
        &.rank-third {
          font-size: 18px;
        }
      }
    }
  }
}
</style>
