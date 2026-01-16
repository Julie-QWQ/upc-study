<template>
  <div class="home-page">
    <div class="page-container">
      <!-- 欢迎区域 -->
      <section class="hero-section">
        <h1 class="hero-title"><SiteName /></h1>
        <p class="hero-subtitle">{{ siteDescription }}</p>
        <div class="hero-stats">
          <div class="stat-item">
            <span class="stat-value">{{ overviewStats?.materials.total || 0 }}</span>
            <span class="stat-label">份资料</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">{{ overviewStats?.users.total || 0 }}</span>
            <span class="stat-label">位用户</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">{{ overviewStats?.downloads.total || 0 }}</span>
            <span class="stat-label">次下载</span>
          </div>
        </div>
      </section>

      <!-- 搜索入口 -->
      <section class="search-section">
        <div class="search-card" @click="$router.push('/search')">
          <svg class="search-icon" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="11" cy="11" r="8"/>
            <path d="m21 21-4.35-4.35"/>
          </svg>
          <div class="search-content">
            <div class="search-title">搜索资料</div>
            <div class="search-hint">查找课件、试卷、实验资料...</div>
          </div>
          <svg class="arrow-icon" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="9 18 15 12 9 6"/>
          </svg>
        </div>
      </section>

      <!-- 快捷入口 -->
      <section class="shortcuts-section">
        <h2 class="section-title">快速访问</h2>
        <div class="shortcuts-grid">
          <div class="shortcut-card" @click="$router.push('/materials/hot')">
            <div class="shortcut-icon hot">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M6 9H4.5a2.5 2.5 0 0 1 0-5H6"/>
                <path d="M18 9h1.5a2.5 2.5 0 0 0 0-5H18"/>
                <path d="M4 22h16"/>
                <path d="M10 14.66V17c0 .55-.47.98-.97 1.21C7.85 18.75 7 20.24 7 22"/>
                <path d="M14 14.66V17c0 .55.47.98.97 1.21C16.15 18.75 17 20.24 17 22"/>
                <path d="M18 2H6v7a6 6 0 0 0 12 0V2Z"/>
              </svg>
            </div>
            <div class="shortcut-content">
              <div class="shortcut-title">热门资料</div>
              <div class="shortcut-desc">查看最受欢迎的资料</div>
            </div>
          </div>

          <div class="shortcut-card" @click="$router.push('/materials/favorites')">
            <div class="shortcut-icon favorites">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z"/>
              </svg>
            </div>
            <div class="shortcut-content">
              <div class="shortcut-title">我的收藏</div>
              <div class="shortcut-desc">查看收藏的资料</div>
            </div>
          </div>

          <div class="shortcut-card" @click="$router.push('/materials/my')">
            <div class="shortcut-icon upload">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                <polyline points="17 8 12 3 7 8"/>
                <line x1="12" y1="3" x2="12" y2="15"/>
              </svg>
            </div>
            <div class="shortcut-content">
              <div class="shortcut-title">我的资料</div>
              <div class="shortcut-desc">管理上传的资料</div>
            </div>
          </div>

          <div v-if="authStore.isCommittee" class="shortcut-card" @click="$router.push('/materials/upload')">
            <div class="shortcut-icon upload">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M12 5v14M5 12h14"/>
              </svg>
            </div>
            <div class="shortcut-content">
              <div class="shortcut-title">上传资料</div>
              <div class="shortcut-desc">分享学习资料</div>
            </div>
          </div>
        </div>
      </section>

      <!-- 最新资料 -->
      <section class="materials-section">
        <div class="section-header">
          <h2 class="section-title">最新资料</h2>
          <button class="view-all-btn" @click="$router.push('/search')">
            查看全部
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="9 18 15 12 9 6"/>
            </svg>
          </button>
        </div>

        <div v-if="loading" class="loading-state">
          <div v-for="i in 4" :key="i" class="skeleton-card"></div>
        </div>

        <div v-else-if="recentMaterials.length === 0" class="empty-state">
          <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
            <polyline points="14 2 14 8 20 8"/>
          </svg>
          <p>暂无资料</p>
        </div>

        <div v-else class="materials-list">
          <div v-for="material in recentMaterials" :key="material.id" class="material-card" @click="$router.push(`/materials/${material.id}`)">
            <div class="material-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                <polyline points="14 2 14 8 20 8"/>
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
                <span>{{ formatRelativeTime(material.created_at) }}</span>
              </div>
            </div>
            <div class="material-stats">
              <div class="stat">
                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                  <polyline points="7 10 12 15 17 10"/>
                  <line x1="12" y1="15" x2="12" y2="3"/>
                </svg>
                <span>{{ material.download_count || 0 }}</span>
              </div>
              <div class="stat">
                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                  <circle cx="12" cy="12" r="3"/>
                </svg>
                <span>{{ material.view_count || 0 }}</span>
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useMaterialStore } from '@/stores/material'
import { useSystemStore } from '@/stores/system'
import { getOverviewStatistics } from '@/api/statistics'
import SiteName from '@/components/SiteName.vue'
import type { OverviewStatistics } from '@/api/statistics'
import type { MaterialCategory } from '@/types'

const systemStore = useSystemStore()
const siteDescription = computed(() => systemStore.getConfig('site_description', '学习资料托管平台'))

const authStore = useAuthStore()
const materialStore = useMaterialStore()

const loading = ref(false)
const overviewStats = ref<OverviewStatistics>()
const recentMaterials = ref<any[]>([])

// 分类文本映射
const getCategoryText = (category: MaterialCategory) => {
  const map: Record<MaterialCategory, string> = {
    courseware: '课件',
    exam: '试卷',
    experiment: '实验',
    exercise: '习题',
    reference: '参考资料',
    other: '其他'
  }
  return map[category] || category
}

// 格式化相对时间
const formatRelativeTime = (date: string) => {
  const d = new Date(date)
  const now = new Date()
  const diff = now.getTime() - d.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (days === 0) return '今天'
  if (days === 1) return '昨天'
  if (days < 7) return `${days} 天前`
  if (days < 30) return `${Math.floor(days / 7)} 周前`
  return d.toLocaleDateString('zh-CN')
}

// 加载统计数据
const loadStatistics = async () => {
  try {
    const { data } = await getOverviewStatistics()
    overviewStats.value = data
  } catch (error) {
    console.error('Failed to load statistics:', error)
  }
}

// 加载最新资料
const loadRecentMaterials = async () => {
  loading.value = true
  try {
    await materialStore.fetchMaterials({
      page: 1,
      size: 5,
      sort_by: 'created_at',
      order: 'desc'
    })
    recentMaterials.value = materialStore.materials
  } catch (error) {
    console.error('Failed to load recent materials:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadStatistics()
  loadRecentMaterials()
})
</script>

<style scoped lang="scss">
.home-page {
  width: 100%;
  background: #ffffff;
  min-height: 100vh;
}

.page-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 32px 24px;
}

// Hero ??????
.hero-section {
  text-align: center;
  padding: 60px 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 16px;
  margin-bottom: 32px;
  color: #ffffff;

  .hero-title {
    font-size: 48px;
    font-weight: 700;
    margin: 0 0 12px 0;
    letter-spacing: -0.02em;
  }

  .hero-subtitle {
    font-size: 20px;
    opacity: 0.9;
    margin: 0 0 32px 0;
  }

  .hero-stats {
    display: flex;
    justify-content: center;
    gap: 48px;

    .stat-item {
      display: flex;
      flex-direction: column;
      gap: 4px;

      .stat-value {
        font-size: 32px;
        font-weight: 700;
      }

      .stat-label {
        font-size: 14px;
        opacity: 0.9;
      }
    }
  }
}

// ????????????
.search-section {
  margin-bottom: 32px;

  .search-card {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 20px 24px;
    background: #fafafa;
    border: 1px solid #e5e7eb;
    border-radius: 12px;
    cursor: pointer;
    transition: all 0.2s;

    &:hover {
      border-color: #111827;
      background: #ffffff;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
    }

    .search-icon {
      flex-shrink: 0;
      color: #6b7280;
    }

    .search-content {
      flex: 1;

      .search-title {
        font-size: 16px;
        font-weight: 600;
        color: #111827;
        margin-bottom: 2px;
      }

      .search-hint {
        font-size: 14px;
        color: #9ca3af;
      }
    }

    .arrow-icon {
      flex-shrink: 0;
      color: #9ca3af;
      transition: transform 0.2s;
    }

    &:hover .arrow-icon {
      transform: translateX(4px);
      color: #111827;
    }
  }
}

// ????????????
.shortcuts-section {
  margin-bottom: 32px;

  .section-title {
    font-size: 24px;
    font-weight: 700;
    color: #111827;
    margin: 0 0 20px 0;
  }

  .shortcuts-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: 16px;
  }

  .shortcut-card {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 20px;
    background: #ffffff;
    border: 1px solid #e5e7eb;
    border-radius: 12px;
    cursor: pointer;
    transition: all 0.2s;

    &:hover {
      border-color: #111827;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
    }

    .shortcut-icon {
      width: 48px;
      height: 48px;
      border-radius: 10px;
      display: flex;
      align-items: center;
      justify-content: center;
      flex-shrink: 0;

      &.hot {
        background: #fef3c7;
        color: #f59e0b;
      }

      &.favorites {
        background: #fee2e2;
        color: #ef4444;
      }

      &.upload {
        background: #dcfce7;
        color: #22c55e;
      }
    }

    .shortcut-content {
      flex: 1;
      min-width: 0;

      .shortcut-title {
        font-size: 15px;
        font-weight: 600;
        color: #111827;
        margin-bottom: 2px;
      }

      .shortcut-desc {
        font-size: 13px;
        color: #9ca3af;
      }
    }
  }
}

// ????????????.materials-section {
  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    .section-title {
      font-size: 24px;
      font-weight: 700;
      color: #111827;
      margin: 0;
    }

    .view-all-btn {
      display: flex;
      align-items: center;
      gap: 6px;
      background: none;
      border: none;
      color: #6b7280;
      font-size: 14px;
      font-weight: 500;
      cursor: pointer;
      transition: color 0.2s;

      &:hover {
        color: #111827;
      }
    }
  }

  .loading-state {
    display: flex;
    flex-direction: column;
    gap: 12px;

    .skeleton-card {
      height: 80px;
      background: linear-gradient(90deg, #f5f5f5 25%, #e8e8e8 50%, #f5f5f5 75%);
      background-size: 200% 100%;
      animation: shimmer 1.5s infinite;
      border-radius: 12px;
    }
  }

  @keyframes shimmer {
    0% { background-position: 200% 0; }
    100% { background-position: -200% 0; }
  }

  .empty-state {
    text-align: center;
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

  .materials-list {
    display: flex;
    flex-direction: column;
    gap: 0;
  }

  .material-card {
    display: grid !important;
    grid-template-columns: auto 1fr auto !important;
    gap: 16px !important;
    align-items: center !important;
    padding: 20px 0 !important;
    background: transparent !important;
    border: none !important;
    border-top: none !important;
    border-bottom: none !important;
    cursor: pointer;
    transition: all 0.2s;

    &:hover {
      opacity: 0.8 !important;
      background: transparent !important;
    }

    &:active {
      opacity: 0.6 !important;
    }

    &:first-child {
      border-top: none !important;
    }

    .material-icon {
      width: 40px;
      height: 40px;
      border-radius: 8px;
      background: #f3f4f6;
      color: #6b7280;
      display: flex;
      align-items: center;
      justify-content: center;
      flex-shrink: 0;
    }

    .material-info {
      min-width: 0;

      .material-title {
        font-size: 15px;
        font-weight: 600;
        color: #111827;
        margin-bottom: 6px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      .material-meta {
        font-size: 13px;
        color: #6b7280;
        margin-bottom: 4px;
        display: flex;
        align-items: center;
        gap: 6px;

        .category-tag {
          padding: 2px 8px;
          background: #f3f4f6;
          border-radius: 4px;
          font-size: 12px;
          font-weight: 500;
        }

        .meta-separator {
          color: #d1d5db;
        }
      }

      .material-details {
        font-size: 12px;
        color: #9ca3af;
        display: flex;
        align-items: center;
        gap: 6px;

        .separator {
          color: #d1d5db;
        }
      }
    }

    .material-stats {
      display: flex;
      gap: 12px;
      flex-shrink: 0;

      .stat {
        display: flex;
        align-items: center;
        gap: 4px;
        font-size: 13px;
        color: #9ca3af;

        svg {
          width: 14px;
          height: 14px;
        }
      }
    }
  }
}

// ?????????@media (max-width: 768px) {
  .page-container {
    padding: 24px 16px;
  }

  .hero-section {
    padding: 40px 20px;

    .hero-title {
      font-size: 36px;
    }

    .hero-subtitle {
      font-size: 16px;
    }

    .hero-stats {
      gap: 24px;

      .stat-item {
        .stat-value {
          font-size: 24px;
        }
      }
    }
  }

  .shortcuts-section {
    .shortcuts-grid {
      grid-template-columns: 1fr;
    }
  }

  .material-card {
    grid-template-columns: auto 1fr;
    gap: 12px;

    .material-stats {
      grid-column: 2;
      justify-content: flex-start;
    }
  }
}
</style>


