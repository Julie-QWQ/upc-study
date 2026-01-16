<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useSearchStore } from '@/stores/search'
import { useMaterialStore } from '@/stores/material'
import { useMaterialCategoryStore } from '@/stores/materialCategory'
import MaterialCard from '@/components/material/MaterialCard.vue'
import type { MaterialListParams, MaterialCategory } from '@/types'

const router = useRouter()
const route = useRoute()
const searchStore = useSearchStore()
const materialStore = useMaterialStore()
const materialCategoryStore = useMaterialCategoryStore()

const filters = ref<MaterialListParams>({
  page: 1,
  size: 20,
  sort_by: 'created_at',
  order: 'desc'
})

const selectedCategory = ref<MaterialCategory | ''>('')
const selectedSort = ref<MaterialListParams['sort_by']>('created_at')
const courseInput = ref('')

// 分类选项 - 从 materialCategoryStore 动态加载
const categories = computed(() => {
  const defaultCategories = [
    { label: '全部分类', value: '' },
    { label: '课件', value: 'courseware' },
    { label: '教材', value: 'textbook' },
    { label: '试卷', value: 'exam_paper' },
    { label: '笔记', value: 'note' },
    { label: '习题', value: 'exercise' },
    { label: '实验', value: 'experiment' },
    { label: '论文', value: 'thesis' },
    { label: '参考资料', value: 'reference' },
    { label: '其他', value: 'other' }
  ]

  try {
    if (!materialCategoryStore) {
      console.warn('materialCategoryStore 未初始化')
      return defaultCategories
    }

    const activeCats = materialCategoryStore.activeCategories || []

    // 如果没有加载到分类,使用默认列表
    if (!Array.isArray(activeCats) || activeCats.length === 0) {
      return defaultCategories
    }

    // 动态分类列表
    const dynamicCategories = [
      { label: '全部分类', value: '' },
      ...activeCats.map(cat => ({
        label: cat.name_zh || cat.name || cat.code,
        value: cat.code
      }))
    ]

    return dynamicCategories
  } catch (error) {
    console.error('获取分类列表出错:', error)
    return defaultCategories
  }
})

// 排序选项
const sortOptions = [
  { label: '最新上传', value: 'created_at' as const },
  { label: '最多下载', value: 'download_count' as const },
  { label: '最多收藏', value: 'favorite_count' as const },
  { label: '最多浏览', value: 'view_count' as const }
]

// 搜索关键词
const keyword = ref('')

// 获取分类标签
const categoryLabel = (value: MaterialCategory | string) => {
  const cat = categories.value.find(c => c.value === value)
  return cat?.label || value
}

// 加载搜索结果
const loadResults = async () => {
  try {
    await materialStore.fetchMaterials({
      page: filters.value.page,
      size: filters.value.size,
      sort_by: filters.value.sort_by,
      order: filters.value.order,
      keyword: keyword.value || undefined,
      category: selectedCategory.value || undefined,
      course_name: courseInput.value || undefined
    })
  } catch (error) {
    console.error('加载搜索结果失败:', error)
  }
}

// 处理搜索
const handleSearch = async () => {
  filters.value.page = 1
  await loadResults()
}

// 处理分类变化
const handleCategoryChange = async () => {
  filters.value.page = 1
  await loadResults()
}

// 处理分类点击
const handleCategoryClick = async (value: string) => {
  selectedCategory.value = value as MaterialCategory | ''
  await handleCategoryChange()
}

// 处理排序变化
const handleSortChange = async () => {
  filters.value.page = 1
  filters.value.sort_by = selectedSort.value
  await loadResults()
}

// 处理排序点击
const handleSortClick = async (value: MaterialListParams['sort_by']) => {
  selectedSort.value = value
  await handleSortChange()
}

// 处理课程名称变化
const handleCourseChange = async () => {
  filters.value.page = 1
  await loadResults()
}

// 处理页码变化
const handlePageChange = async (page: number) => {
  filters.value.page = page
  await loadResults()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

// 重置筛选
const handleReset = async () => {
  selectedCategory.value = ''
  selectedSort.value = 'created_at'
  courseInput.value = ''
  keyword.value = ''
  filters.value = {
    page: 1,
    size: 20,
    sort_by: 'created_at',
    order: 'desc'
  }
  await loadResults()
}

onMounted(async () => {
  // 加载资料分类（如果失败就使用默认分类）
  // 注意：这个 API 需要认证，未登录时会失败，但会使用默认分类列表
  materialCategoryStore.fetchActiveCategories().catch(err => {
    console.log('加载动态分类失败,使用默认分类列表')
  }).finally(() => {
    // 无论成功失败，都继续加载搜索结果
    const queryKeyword = route.query.keyword as string
    if (queryKeyword) {
      keyword.value = queryKeyword
    }
    loadResults()
  })
})
</script>

<template>
  <div class="search-page">
    <div class="page-container">
      <main class="content-stream">
        <!-- 页面标题 -->
        <header class="stream-header">
          <h1 class="page-title">资料搜索</h1>
        </header>

        <!-- 左右布局 -->
        <div class="search-layout">
          <!-- 左侧结果列表 -->
          <div class="results-area">
            <!-- 当前筛选条件标签 -->
            <div v-if="selectedCategory || courseInput || keyword" class="active-filters">
              <span class="filters-label">当前筛选:</span>
              <span v-if="keyword" class="filter-tag">
                关键词: {{ keyword }}
                <button class="tag-remove" @click="keyword = ''; handleSearch()">
                  <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M18 6L6 18M6 6l12 12"/>
                  </svg>
                </button>
              </span>
              <span v-if="selectedCategory" class="filter-tag">
                分类: {{ categoryLabel(selectedCategory) }}
                <button class="tag-remove" @click="selectedCategory = ''; handleCategoryChange()">
                  <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M18 6L6 18M6 6l12 12"/>
                  </svg>
                </button>
              </span>
              <span v-if="courseInput" class="filter-tag">
                课程: {{ courseInput }}
                <button class="tag-remove" @click="courseInput = ''; handleCourseChange()">
                  <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M18 6L6 18M6 6l12 12"/>
                  </svg>
                </button>
              </span>
            </div>

            <!-- 加载状态 -->
            <div v-if="materialStore.loading" class="loading-state">
              <div v-for="i in 6" :key="i" class="skeleton-card"></div>
            </div>

            <!-- 空状态 -->
            <div v-else-if="materialStore.materials.length === 0" class="empty-state">
              <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <circle cx="11" cy="11" r="8"/>
                <path d="m21 21-4.35-4.35"/>
              </svg>
              <h2>未找到相关资料</h2>
              <p>尝试调整搜索关键词或筛选条件</p>
              <button class="reset-btn-large" @click="handleReset">清除所有筛选</button>
            </div>

            <!-- 搜索结果 -->
            <div v-else>
              <div class="results-header">
                <h2 class="results-title">搜索结果</h2>
                <span class="results-count">共 {{ materialStore.total }} 条</span>
              </div>

              <div class="materials-list">
                <MaterialCard
                  v-for="material in materialStore.materials"
                  :key="material.id"
                  :material="material"
                />
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
          </div>

          <!-- 右侧筛选面板 -->
          <aside class="filters-sidebar">
            <div class="filter-panel">
              <!-- 课程名称 -->
              <div class="filter-group">
                <label class="filter-label">课程名称</label>
                <input
                  v-model="courseInput"
                  type="text"
                  placeholder="输入课程名称"
                  class="course-filter-input"
                  @change="handleCourseChange"
                />
              </div>

              <!-- 分类筛选 -->
              <div class="filter-group">
                <label class="filter-label">资料分类</label>
                <div class="category-list">
                  <button
                    v-for="cat in categories"
                    :key="cat.value"
                    type="button"
                    :class="['category-item', { active: selectedCategory === cat.value }]"
                    @click="handleCategoryClick(cat.value)"
                  >
                    {{ cat.label }}
                  </button>
                </div>
              </div>

              <!-- 搜索框 -->
              <div class="filter-group">
                <label class="filter-label">关键词搜索</label>
                <div class="search-input-wrapper">
                  <svg class="search-icon" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="11" cy="11" r="8"/>
                    <path d="m21 21-4.35-4.35"/>
                  </svg>
                  <input
                    v-model="keyword"
                    type="text"
                    placeholder="搜索标题、描述..."
                    class="search-input"
                    @keyup.enter="handleSearch"
                  />
                </div>
                <button class="search-submit-btn" @click="handleSearch">
                  搜索
                </button>
              </div>

              <!-- 排序方式 -->
              <div class="filter-group">
                <label class="filter-label">排序方式</label>
                <div class="sort-list">
                  <button
                    v-for="sort in sortOptions"
                    :key="sort.value"
                    type="button"
                    :class="['sort-item', { active: selectedSort === sort.value }]"
                    @click="handleSortClick(sort.value)"
                  >
                    {{ sort.label }}
                  </button>
                </div>
              </div>

              <!-- 重置按钮 -->
              <button v-if="selectedCategory || courseInput || keyword" class="reset-filters-btn" @click="handleReset">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 12"/>
                  <path d="M3 3v9h9"/>
                </svg>
                清除筛选
              </button>
            </div>
          </aside>
        </div>
      </main>
    </div>
  </div>
</template>

<style scoped lang="scss">
.search-page {
  width: 100%;
  background: #ffffff;
  min-height: 100vh;
}

.page-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 32px 24px;
}

.content-stream {
  max-width: 100%;
}

.stream-header {
  margin-bottom: 24px;

  .page-title {
    font-size: 28px;
    font-weight: 700;
    color: #111827;
    margin: 0;
  }
}

.search-layout {
  display: grid;
  grid-template-columns: 1fr 320px;
  gap: 32px;
  align-items: start;
}

// 右侧筛选面板
.filters-sidebar {
  position: sticky;
  top: 24px;

  .filter-panel {
    background: #ffffff;
    border: 1px solid #e5e7eb;
    border-radius: 12px;
    padding: 20px;

    .filter-group {
      margin-bottom: 24px;

      &:last-child {
        margin-bottom: 0;
      }

      .filter-label {
        display: block;
        font-size: 13px;
        font-weight: 600;
        color: #374151;
        margin-bottom: 10px;
      }

      // 搜索框
      .search-input-wrapper {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 10px 12px;
        background: #ffffff;
        border: 1px solid #e5e7eb;
        border-radius: 8px;
        transition: all 0.2s;

        &:focus-within {
          border-color: #111827;
        }

        .search-icon {
          flex-shrink: 0;
          color: #9ca3af;
        }

        .search-input {
          flex: 1;
          border: none;
          background: transparent;
          font-size: 14px;
          color: #111827;
          outline: none;

          &::placeholder {
            color: #9ca3af;
          }
        }
      }

      .search-submit-btn {
        width: 100%;
        margin-top: 8px;
        padding: 10px;
        background: #111827;
        color: #ffffff;
        border: none;
        border-radius: 8px;
        font-size: 14px;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.2s;

        &:hover {
          background: #000000;
        }
      }

      // 分类列表
      .category-list {
        display: flex;
        flex-direction: column;
        gap: 4px;

        .category-item {
          padding: 8px 12px;
          background: #ffffff;
          border: 1px solid #e5e7eb;
          border-radius: 8px;
          font-size: 14px;
          color: #6b7280;
          cursor: pointer;
          transition: all 0.2s;
          text-align: left;
          width: 100%;
          display: block;
          position: relative;
          z-index: 1;

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

      // 课程输入框
      .course-filter-input {
        width: 100%;
        padding: 10px 12px;
        background: #ffffff;
        border: 1px solid #e5e7eb;
        border-radius: 8px;
        font-size: 14px;
        color: #111827;
        outline: none;
        transition: all 0.2s;

        &:focus {
          border-color: #111827;
        }

        &::placeholder {
          color: #9ca3af;
        }
      }

      // 排序列表
      .sort-list {
        display: flex;
        flex-direction: column;
        gap: 4px;

        .sort-item {
          padding: 8px 12px;
          background: #ffffff;
          border: 1px solid #e5e7eb;
          border-radius: 8px;
          font-size: 14px;
          color: #6b7280;
          cursor: pointer;
          transition: all 0.2s;
          text-align: left;
          width: 100%;
          display: block;
          position: relative;
          z-index: 1;

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
    }

    .reset-filters-btn {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 6px;
      width: 100%;
      padding: 10px;
      background: #ffffff;
      border: 1px solid #dc2626;
      color: #dc2626;
      border-radius: 8px;
      font-size: 14px;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.2s;

      &:hover {
        background: #dc2626;
        color: #ffffff;
      }

      svg {
        width: 16px;
        height: 16px;
      }
    }
  }
}

// 左侧结果区域
.results-area {
  flex: 1;
  min-width: 0;

  .active-filters {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-bottom: 20px;
    padding: 12px 16px;
    background: #f3f4f6;
    border-radius: 8px;

    .filters-label {
      font-size: 13px;
      font-weight: 600;
      color: #6b7280;
      align-self: center;
    }

    .filter-tag {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      padding: 4px 10px;
      background: #111827;
      color: #ffffff;
      border-radius: 6px;
      font-size: 13px;
      font-weight: 500;

      .tag-remove {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 14px;
        height: 14px;
        background: rgba(255, 255, 255, 0.2);
        border: none;
        border-radius: 2px;
        cursor: pointer;
        transition: background 0.2s;

        &:hover {
          background: rgba(255, 255, 255, 0.3);
        }

        svg {
          width: 10px;
          height: 10px;
        }
      }
    }
  }

  .loading-state {
    display: flex;
    flex-direction: column;
    gap: 16px;

    .skeleton-card {
      height: 200px;
      background: linear-gradient(90deg, #f5f5f5 25%, #e8e8e8 50%, #f5f5f5 75%);
      background-size: 200% 100%;
      animation: shimmer 1.5s infinite;
      border-radius: 8px;
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
      margin: 0 0 24px 0;
    }

    .reset-btn-large {
      padding: 10px 24px;
      background: #111827;
      color: #ffffff;
      border: none;
      border-radius: 8px;
      font-size: 14px;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.2s;

      &:hover {
        background: #000000;
      }
    }
  }

  .results-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 16px;

    .results-title {
      font-size: 20px;
      font-weight: 600;
      color: #111827;
      margin: 0;
    }

    .results-count {
      font-size: 14px;
      color: #9ca3af;
    }
  }

  .materials-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }
}

.pagination-section {
  margin-top: 40px;
  padding-top: 24px;
  border-top: 1px solid #f2f2f2;
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

// 响应式设计
@media (max-width: 1024px) {
  .search-layout {
    grid-template-columns: 1fr;
  }

  .filters-sidebar {
    position: static;

    .filter-panel {
      margin-bottom: 24px;
    }
  }
}

@media (max-width: 768px) {
  .page-container {
    padding: 24px 16px;
  }

  .stream-header {
    .page-title {
      font-size: 24px;
    }
  }

  .results-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }
}
</style>
