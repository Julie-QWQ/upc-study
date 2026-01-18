<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import MaterialUpload from '@/components/material/MaterialUpload.vue'
import { useAuthStore } from '@/stores/auth'
import { onMounted } from 'vue'

const router = useRouter()
const authStore = useAuthStore()

// 权限检查
onMounted(() => {
  if (!authStore.isCommittee) {
    ElMessage.warning('您没有权限访问此页面')
    router.push('/materials')
  }
})

// 上传成功
const handleSuccess = () => {
  ElMessage.success('资料上传成功,等待管理员审核')
  router.push('/materials/my')
}

// 取消上传
const handleCancel = () => {
  router.back()
}
</script>

<template>
  <div class="upload-page">
    <div class="page-container">
      <main class="content-stream">
        <!-- 页面标题 -->
        <header class="stream-header">
          <h1 class="page-title">上传资料</h1>
        </header>

        <!-- 上传表单 -->
        <div class="upload-content">
          <MaterialUpload @success="handleSuccess" @cancel="handleCancel" />
        </div>
      </main>
    </div>
  </div>
</template>

<style scoped lang="scss">
.upload-page {
  width: 100%;
  background: #ffffff;
  min-height: 100vh;
}

.page-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 32px 24px;
}

.content-stream {
  max-width: 100%;
}

.stream-header {
  display: flex;
  align-items: center;
  gap: 16px;
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
}

.upload-content {
  // 内容样式在 MaterialUpload 组件中定义
}

@media (max-width: 768px) {
  .page-container {
    padding: 24px 16px;
  }

  .stream-header {
    flex-direction: column;
    align-items: flex-start;

    .page-title {
      font-size: 26px;
    }
  }
}
</style>
