<template>
  <div class="auth-split">
    <div class="auth-split-grid">
      <section class="auth-intro">
        <div>
          <h1 class="auth-brand-title"><SiteName /></h1>
          <p class="auth-brand-subtitle">{{ siteDescription }}</p>
        </div>
        <div>
          <h2 class="auth-intro-heading">创建你的学习账号。</h2>
          <p class="auth-intro-text">
            上传、整理并共享学习资料，让学习更高效。
          </p>
          <ul class="auth-intro-list">
            <li>快速上传与分类资料</li>
            <li>随时查看下载与收藏记录</li>
            <li>获取公告与审核通知</li>
          </ul>
        </div>
      </section>

      <section class="auth-panel">
        <div class="auth-card">
          <h2 class="auth-title">注册</h2>
          <p class="auth-description">创建账号，开始使用学习资料平台</p>

          <RegisterForm @success="handleRegisterSuccess" />

          <div class="auth-footer">
            已有账号？
            <router-link to="/login">立即登录</router-link>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import RegisterForm from '@/components/RegisterForm.vue'
import SiteName from '@/components/SiteName.vue'
import { useSystemStore } from '@/stores/system'

const systemStore = useSystemStore()
const siteDescription = computed(() => systemStore.getConfig('site_description', '学院学习资料托管平台'))

const handleRegisterSuccess = () => {
  // 注册成功后由 composable 处理跳转
}
</script>

<style scoped lang="scss">
@import '@/assets/styles/auth.scss';
</style>
