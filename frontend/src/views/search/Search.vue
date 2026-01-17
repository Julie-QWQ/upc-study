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
          <p class="page-subtitle">探索海量学习资源</p>
        </header>

        <!-- 搜索区域 -->
        <section class="search-section">
          <!-- 课程名称筛选 -->
          <div class="course-filter">
            <div class="course-search-bar">
              <svg class="search-icon" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="11" cy="11" r="8"/>
                <path d="m21 21-4.35-4.35"/>
              </svg>
              <input
                v-model="courseInput"
                type="text"
                placeholder="输入课程名称筛选"
                class="search-input"
                @change="handleCourseChange"
              />
            </div>
            <button class="course-search-btn" @click="handleSearch">
              搜索
            </button>
          </div>

          <!-- 快速筛选标签 -->
          <div class="quick-filters">
            <div class="filter-group">
              <span class="filter-group-label">分类:</span>
              <div class="filter-chips">
                <button
                  v-for="cat in categories.slice(0, 8)"
                  :key="cat.value"
                  :class="['chip', { active: selectedCategory === cat.value }]"
                  @click="handleCategoryClick(cat.value)"
                >
                  {{ cat.label }}
                </button>
                <el-dropdown trigger="click" @command="handleCategoryClick">
                  <button class="chip chip-more">
                    更多
                    <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="m6 9 6 6 6-6"/>
                    </svg>
                  </button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item
                        v-for="cat in categories.slice(8)"
                        :key="cat.value"
                        :command="cat.value"
                        :class="{ 'is-active': selectedCategory === cat.value }"
                      >
                        {{ cat.label }}
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </div>
            </div>

            <div class="filter-group">
              <span class="filter-group-label">排序:</span>
              <div class="filter-chips">
                <button
                  v-for="sort in sortOptions"
                  :key="sort.value"
                  :class="['chip', { active: selectedSort === sort.value }]"
                  @click="handleSortClick(sort.value)"
                >
                  {{ sort.label }}
                </button>
              </div>
            </div>
          </div>

          <!-- 关键词搜索 -->
          <div class="keyword-filter">
            <label class="keyword-label">关键词搜索</label>
            <input
              v-model="keyword"
              type="text"
              placeholder="搜索资料标题、描述..."
              class="keyword-input"
              @keyup.enter="handleSearch"
            />
          </div>

          <!-- 当前筛选条件 -->
          <div v-if="selectedCategory || courseInput || keyword" class="active-filters">
            <div class="active-tags">
              <span v-if="courseInput" class="tag">
                课程: {{ courseInput }}
                <button class="tag-close" @click="courseInput = ''; handleCourseChange()">
                  <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M18 6L6 18M6 6l12 12"/>
                  </svg>
                </button>
              </span>
              <span v-if="selectedCategory" class="tag">
                分类: {{ categoryLabel(selectedCategory) }}
                <button class="tag-close" @click="selectedCategory = ''; handleCategoryChange()">
                  <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M18 6L6 18M6 6l12 12"/>
                  </svg>
                </button>
              </span>
              <span v-if="keyword" class="tag">
                关键词: {{ keyword }}
                <button class="tag-close" @click="keyword = ''; handleSearch()">
                  <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M18 6L6 18M6 6l12 12"/>
                  </svg>
                </button>
              </span>
            </div>
            <button class="clear-all-btn" @click="handleReset">
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 12"/>
                <path d="M3 3v9h9"/>
              </svg>
              清除筛选
            </button>
          </div>
        </section>

        <!-- 搜索结果 -->
        <section class="results-section">
          <!-- 加载状态 -->
          <div v-if="materialStore.loading" class="loading-state">
            <div v-for="i in 6" :key="i" class="skeleton-card"></div>
          </div>

          <!-- 空状态 -->
          <div v-else-if="materialStore.materials.length === 0" class="empty-state">
            <div class="empty-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="80" height="80" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <circle cx="11" cy="11" r="8"/>
                <path d="m21 21-4.35-4.35"/>
              </svg>
            </div>
            <h2 class="empty-title">未找到相关资料</h2>
            <p class="empty-desc">尝试调整搜索关键词或筛选条件</p>
            <button class="empty-action" @click="handleReset">清除所有筛选</button>
          </div>

          <!-- 结果列表 -->
          <div v-else class="results-content">
            <div class="results-info">
              <h2 class="results-title">搜索结果</h2>
              <span class="results-count">{{ materialStore.total }} 条资料</span>
            </div>

            <div class="materials-list">
              <MaterialCard
                v-for="material in materialStore.materials"
                :key="material.id"
                :material="material"
              />
            </div>

            <!-- 分页 -->
            <div v-if="materialStore.total > materialStore.size" class="pagination-section">
              <el-pagination
                :current-page="materialStore.page"
                :page-size="materialStore.size"
                :total="materialStore.total"
                layout="prev, pager, next"
                @current-change="handlePageChange"
              />
            </div>
          </div>
        </section>
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
  max-width: 1200px;
  margin: 0 auto;
  padding: 32px 24px;
}

.content-stream {
  max-width: 100%;
}

// 页面头部
.stream-header {
  margin-bottom: 32px;
  padding-bottom: 24px;
  border-bottom: 1px solid #f2f2f2;

  .page-title {
    font-size: 32px;
    font-weight: 700;
    color: #1a1a1a;
    margin: 0 0 8px 0;
    letter-spacing: -0.02em;
  }

  .page-subtitle {
    font-size: 15px;
    color: #6b7280;
    margin: 0;
  }
}

// 搜索区域
.search-section {
  background: #ffffff;
  border: 1px solid #f2f2f2;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 32px;
}

// 课程筛选
.course-filter {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;

  .course-search-bar {
    flex: 1;
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 16px;
    background: #f9fafb;
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    transition: all 0.2s;

    &:focus-within {
      background: #ffffff;
      border-color: #1a1a1a;
      box-shadow: 0 0 0 3px rgba(26, 26, 26, 0.05);
    }

    .search-icon {
      flex-shrink: 0;
      color: #9ca3af;
    }

    .search-input {
      flex: 1;
      border: none;
      background: transparent;
      font-size: 15px;
      color: #1a1a1a;
      outline: none;

      &::placeholder {
        color: #9ca3af;
      }
    }
  }

  .course-search-btn {
    padding: 10px 24px;
    background: #1a1a1a;
    color: #ffffff;
    border: none;
    border-radius: 8px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
    white-space: nowrap;

    &:hover {
      background: #000000;
    }

    &:active {
      transform: scale(0.98);
    }
  }
}

// 快速筛选
.quick-filters {
  margin-bottom: 24px;

  .filter-group {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 16px;

    &:last-child {
      margin-bottom: 0;
    }

    .filter-group-label {
      font-size: 14px;
      font-weight: 600;
      color: #374151;
      white-space: nowrap;
    }

    .filter-chips {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
      flex: 1;
    }
  }

  .chip {
    padding: 6px 12px;
    background: #f9fafb;
    border: 1px solid #e5e7eb;
    border-radius: 6px;
    font-size: 13px;
    color: #6b7280;
    cursor: pointer;
    transition: all 0.2s;
    font-weight: 500;
    display: inline-flex;
    align-items: center;
    gap: 4px;

    &:hover {
      background: #ffffff;
      border-color: #1a1a1a;
      color: #1a1a1a;
    }

    &.active {
      background: #1a1a1a;
      border-color: #1a1a1a;
      color: #ffffff;
    }

    &.chip-more {
      color: #1a1a1a;

      svg {
        width: 10px;
        height: 10px;
      }
    }
  }
}

// 关键词搜索
.keyword-filter {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 24px;

  .keyword-label {
    font-size: 14px;
    font-weight: 600;
    color: #374151;
    white-space: nowrap;
  }

  .keyword-input {
    flex: 1;
    padding: 10px 14px;
    background: #ffffff;
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    font-size: 14px;
    color: #1a1a1a;
    outline: none;
    transition: all 0.2s;

    &:focus {
      border-color: #1a1a1a;
      box-shadow: 0 0 0 3px rgba(26, 26, 26, 0.05);
    }

    &::placeholder {
      color: #9ca3af;
    }
  }
}

// 当前筛选标签
.active-filters {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 4px;

  .active-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    flex: 1;

    .tag {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      padding: 6px 10px;
      background: #1a1a1a;
      color: #ffffff;
      border-radius: 6px;
      font-size: 13px;
      font-weight: 500;

      .tag-close {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 16px;
        height: 16px;
        background: rgba(255, 255, 255, 0.15);
        border: none;
        border-radius: 4px;
        cursor: pointer;
        transition: all 0.2s;

        &:hover {
          background: rgba(255, 255, 255, 0.25);
        }

        svg {
          width: 10px;
          height: 10px;
        }
      }
    }
  }

  .clear-all-btn {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 8px 16px;
    background: #ffffff;
    border: 1px solid #dc2626;
    color: #dc2626;
    border-radius: 6px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
    white-space: nowrap;
    flex-shrink: 0;

    &:hover {
      background: #dc2626;
      color: #ffffff;
    }

    svg {
      width: 14px;
      height: 14px;
    }
  }
}

// 搜索结果区域
.results-section {
  background: #ffffff;
}

// 加载状态
.loading-state {
  display: flex;
  flex-direction: column;
  gap: 16px;

  .skeleton-card {
    height: 160px;
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

// 空状态
.empty-state {
  text-align: center;
  padding: 80px 20px;
  background: #f9fafb;
  border-radius: 12px;
  border: 1px solid #f2f2f2;

  .empty-icon {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 80px;
    height: 80px;
    background: #ffffff;
    border-radius: 50%;
    margin-bottom: 24px;
    color: #d1d5db;
    border: 2px solid #f2f2f2;

    svg {
      width: 40px;
      height: 40px;
    }
  }

  .empty-title {
    font-size: 20px;
    font-weight: 600;
    color: #1a1a1a;
    margin: 0 0 8px 0;
  }

  .empty-desc {
    font-size: 14px;
    color: #6b7280;
    margin: 0 0 24px 0;
  }

  .empty-action {
    padding: 10px 24px;
    background: #1a1a1a;
    color: #ffffff;
    border: none;
    border-radius: 6px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;

    &:hover {
      background: #000000;
    }
  }
}

// 结果内容
.results-content {
  .results-info {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 20px;
    padding-bottom: 16px;
    border-bottom: 1px solid #f2f2f2;

    .results-title {
      font-size: 20px;
      font-weight: 600;
      color: #1a1a1a;
      margin: 0;
    }

    .results-count {
      font-size: 14px;
      color: #6b7280;
      background: #f9fafb;
      padding: 6px 12px;
      border-radius: 6px;
      border: 1px solid #f2f2f2;
    }
  }

  .materials-list {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }
}

// 分页
.pagination-section {
  margin-top: 40px;
  padding-top: 24px;
  border-top: 1px solid #f2f2f2;
  display: flex;
  justify-content: center;
}

// Element Plus 下拉菜单 Substack 风格
:deep(.el-dropdown-menu) {
  background: #ffffff !important;
  border: 1px solid #f2f2f2 !important;
  border-radius: 8px !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08) !important;
  padding: 6px !important;
  min-width: 140px !important;

  .el-dropdown-menu__item,
  .el-dropdown-menu-item {
    padding: 8px 12px !important;
    border-radius: 6px !important;
    font-size: 14px !important;
    color: #1a1a1a !important;
    transition: all 0.15s !important;
    line-height: 1.5 !important;

    &:hover {
      background: #f9fafb !important;
      color: #1a1a1a !important;
    }

    &.is-active {
      background: #1a1a1a !important;
      color: #ffffff !important;
      font-weight: 500 !important;
    }
  }
}

// Element Plus 下拉菜单项的额外样式覆盖
:deep(.el-dropdown-menu li) {
  padding: 8px 12px !important;
  border-radius: 6px !important;
  font-size: 14px !important;
  color: #1a1a1a !important;
  transition: all 0.15s !important;
  line-height: 1.5 !important;

  &:hover {
    background: #f9fafb !important;
    color: #1a1a1a !important;
  }

  &.is-active {
    background: #1a1a1a !important;
    color: #ffffff !important;
    font-weight: 500 !important;
  }
}

// Element Plus 分页 Substack 风格
:deep(.el-pagination) {
  .btn-prev,
  .btn-next,
  .el-pager li {
    background: #ffffff;
    border: 1px solid #e5e7eb;
    border-radius: 6px;
    font-weight: 500;
    color: #6b7280;
    transition: all 0.2s;

    &:hover {
      border-color: #1a1a1a;
      color: #1a1a1a;
    }

    &.active {
      background: #1a1a1a;
      border-color: #1a1a1a;
      color: #ffffff;
      font-weight: 600;
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .page-container {
    padding: 24px 16px;
  }

  .stream-header {
    .page-title {
      font-size: 28px;
    }

    .page-subtitle {
      font-size: 14px;
    }
  }

  .search-section {
    padding: 20px;
  }

  .course-filter {
    flex-direction: column;

    .course-search-btn {
      width: 100%;
    }
  }

  .keyword-filter {
    flex-direction: column;
    align-items: flex-start;

    .keyword-input {
      width: 100%;
      font-size: 16px; // 防止 iOS 缩放
    }
  }

  .quick-filters {
    .filter-group {
      flex-direction: column;
      align-items: flex-start;
      gap: 8px;

      .filter-chips {
        width: 100%;
      }
    }

    .chip {
      font-size: 12px;
      padding: 6px 10px;
    }
  }

  .active-filters {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;

    .clear-all-btn {
      width: 100%;
      justify-content: center;
    }
  }

  .results-info {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
}
</style>

<style lang="scss">
// Element Plus 下拉菜单 Substack 风格 - 全局样式
.el-dropdown-menu {
  background: #ffffff !important;
  border: 1px solid #f2f2f2 !important;
  border-radius: 8px !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08) !important;
  padding: 6px !important;
  min-width: 140px !important;

  .el-dropdown-menu__item {
    padding: 8px 12px !important;
    border-radius: 6px !important;
    font-size: 14px !important;
    color: #1a1a1a !important;
    transition: all 0.15s !important;
    line-height: 1.5 !important;

    &:hover {
      background: #f9fafb !important;
      color: #1a1a1a !important;
    }

    &.is-active {
      background: #1a1a1a !important;
      color: #ffffff !important;
      font-weight: 500 !important;
    }
  }
}
</style>
