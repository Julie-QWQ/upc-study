<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useMaterialStore } from '@/stores/material'
import { useAuthStore } from '@/stores/auth'
import MaterialCard from '@/components/material/MaterialCard.vue'
import AnnouncementBoard from '@/components/layout/AnnouncementBoard.vue'
import type { MaterialListParams } from '@/types'

const materialStore = useMaterialStore()
const authStore = useAuthStore()

// 当前筛选条件
const filters = ref<MaterialListParams>({
  page: 1,
  size: 20,
  sort_by: 'created_at',
  order: 'desc',
  uploader_id: undefined // 明确设置为 undefined,防止从"我的资料"切换过来时残留
})

// 判断用户是否可以上传资料(学委或管理员)
const canUploadMaterial = computed(() => authStore.isCommittee)

// 加载资料列表
const loadMaterials = async () => {
  await materialStore.fetchMaterials(filters.value)
}

// 处理页码变化
const handlePageChange = (page: number) => {
  filters.value.page = page
  loadMaterials()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => {
  loadMaterials()
})
</script>

<template>
  <div class="materials-page">
    <div class="page-container">
      <!-- 中间内容流 -->
      <main class="content-stream">
        <!-- 页面标题和操作 -->
        <header class="stream-header">
          <h1 class="page-title">学习资料</h1>
          <el-button
            v-if="canUploadMaterial"
            type="primary"
            @click="$router.push('/materials/upload')"
            class="upload-btn"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 5v14M5 12h14"/>
            </svg>
            上传资料
          </el-button>
        </header>

        <!-- 加载状态 -->
        <div v-if="materialStore.loading" class="loading-state">
          <div v-for="i in 6" :key="i" class="skeleton-card"></div>
        </div>

        <!-- 空状态 -->
        <div v-else-if="materialStore.materials.length === 0" class="empty-state">
          <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
            <polyline points="14 2 14 8 20 8"/>
          </svg>
          <h2>暂无资料</h2>
          <p>还没有上传任何学习资料</p>
          <button v-if="canUploadMaterial" @click="$router.push('/materials/upload')" class="upload-first-btn">
            上传第一个资料
          </button>
        </div>

        <!-- 资料列表 -->
        <div v-else class="materials-list">
          <MaterialCard
            v-for="material in materialStore.materials"
            :key="material.id"
            :material="material"
          />
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

      <!-- 右侧边栏 -->
      <aside class="right-sidebar">
        <!-- 公告栏 -->
        <AnnouncementBoard />
      </aside>
    </div>
  </div>
</template>

<style scoped lang="scss">
.materials-page {
  width: 100%;
  background: #ffffff;
  min-height: 100vh;
}

.page-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 32px 24px;
  display: grid;
  grid-template-columns: 1fr 340px;
  gap: 48px;
}

// main content area
.content-stream {
  max-width: 100%;
}

.stream-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
  padding-bottom: 24px;
  border-bottom: 1px solid #f2f2f2;

  .page-title {
    font-size: 32px;
    font-weight: 700;
    color: #1a1a1a;
    margin: 0;
    letter-spacing: -0.02em;
  }

  .upload-btn {
    background: #1a1a1a;
    color: #ffffff;
    border: none;
    padding: 10px 20px;
    font-size: 14px;
    font-weight: 500;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s;
    display: flex;
    align-items: center;
    gap: 6px;

    &:hover {
      background: #333;
    }

    svg {
      width: 16px;
      height: 16px;
    }
  }
}

// material list
.materials-list {
  display: flex;
  flex-direction: column;
  gap: 0;
}

// loading state
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

// empty state
.empty-state {
  text-align: center;
  padding: 80px 40px;
  background: #fafafa;
  border-radius: 12px;

  svg {
    color: #ccc;
    margin-bottom: 24px;
  }

  h2 {
    font-size: 22px;
    font-weight: 600;
    color: #1a1a1a;
    margin: 0 0 8px 0;
  }

  p {
    font-size: 15px;
    color: #999;
    margin: 0 0 24px 0;
  }

  .upload-first-btn {
    background: #1a1a1a;
    color: #ffffff;
    border: none;
    padding: 12px 24px;
    font-size: 14px;
    font-weight: 500;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s;

    &:hover {
      background: #333;
    }
  }
}

// pagination
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
    color: #666;

    &:hover {
      color: #1a1a1a;
    }

    &.active {
      color: #1a1a1a;
      font-weight: 600;
    }
  }
}

// right sidebar
.right-sidebar {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

// responsive
@media (max-width: 1024px) {
  .page-container {
    grid-template-columns: 1fr;
    gap: 32px;
  }

  .right-sidebar {
    order: -1;
  }
}

@media (max-width: 768px) {
  .page-container {
    padding: 24px 16px;
  }

  .stream-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;

    .page-title {
      font-size: 26px;
    }
  }
}
</style>


