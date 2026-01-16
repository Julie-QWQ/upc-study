<template>
  <div class="auth-split">
    <div class="auth-split-grid">
      <section class="auth-intro">
        <div>
          <h1 class="auth-brand-title"><SiteName /></h1>
          <p class="auth-brand-subtitle">{{ siteDescription }}</p>
        </div>
        <div>
          <h2 class="auth-intro-heading">把资料整理成你的学习基地</h2>
          <p class="auth-intro-text">
            课程资料、试卷与实验一站管理，快速搜索与分享。
          </p>
          <ul class="auth-intro-list">
            <li>搜索课件、试卷与实验资料</li>
            <li>收藏常用资源并追踪下载</li>
            <li>与同学共享高质量资料</li>
          </ul>
        </div>
      </section>

      <section class="auth-panel">
        <div class="auth-card">
          <h2 class="auth-title">登录</h2>
          <p class="auth-description">欢迎回来，请登录你的账号</p>

          <LoginForm @success="handleLoginSuccess" />

          <div class="auth-footer">
            还没有账号？
            <router-link to="/register">立即注册</router-link>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import LoginForm from '@/components/LoginForm.vue'
import SiteName from '@/components/SiteName.vue'
import { useSystemStore } from '@/stores/system'

const systemStore = useSystemStore()
const siteDescription = computed(() => systemStore.getConfig('site_description', '学院学习资料托管平台'))

const handleLoginSuccess = () => {
  // 登录成功后由 composable 处理跳转
}
</script>

<style scoped lang="scss">
@import '@/assets/styles/auth.scss';
</style>

