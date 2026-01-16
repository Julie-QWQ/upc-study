<script setup lang="ts">
import { useAuth } from '@/composables/useAuth'
import { useRegisterForm } from '@/composables/useValidator'
import type { RegisterRequest } from '@/types'

const emit = defineEmits<{
  success: []
}>()

const { register } = useAuth()

const {
  username,
  email,
  password,
  confirmPassword,
  real_name,
  major,
  className,
  errors,
  isSubmitting,
  handleSubmit
} = useRegisterForm()

const onSubmit = handleSubmit.withControlled(async () => {
  const data: RegisterRequest = {
    username: username.value,
    email: email.value,
    password: password.value,
    real_name: real_name.value,
    major: major.value,
    class: className.value || ''
  }

  const success = await register(data)
  if (success) {
    emit('success')
  }
})
</script>

<template>
  <form @submit="onSubmit" class="auth-form auth-form--register">
    <div class="form-grid">
      <div class="form-group">
        <label for="username">用户名</label>
        <input
          id="username"
          v-model="username"
          type="text"
          placeholder="输入用户名（字母、数字）"
          autocomplete="username"
          :disabled="isSubmitting"
        />
        <p v-if="errors.username" class="error-message">{{ errors.username }}</p>
      </div>

      <div class="form-group">
        <label for="email">邮箱</label>
        <input
          id="email"
          v-model="email"
          type="email"
          placeholder="输入邮箱地址"
          autocomplete="email"
          :disabled="isSubmitting"
        />
        <p v-if="errors.email" class="error-message">{{ errors.email }}</p>
      </div>

      <div class="form-group">
        <label for="real_name">真实姓名</label>
        <input
          id="real_name"
          v-model="real_name"
          type="text"
          placeholder="输入真实姓名"
          autocomplete="name"
          :disabled="isSubmitting"
        />
        <p v-if="errors.real_name" class="error-message">{{ errors.real_name }}</p>
      </div>

      <div class="form-group">
        <label for="major">专业</label>
        <input
          id="major"
          v-model="major"
          type="text"
          placeholder="输入专业"
          :disabled="isSubmitting"
        />
        <p v-if="errors.major" class="error-message">{{ errors.major }}</p>
      </div>

      <div class="form-group">
        <label for="class">班级</label>
        <input
          id="class"
          v-model="className"
          type="text"
          placeholder="输入班级"
          :disabled="isSubmitting"
        />
        <p v-if="errors.class" class="error-message">{{ errors.class }}</p>
      </div>

      <div class="form-group full-width">
        <label for="password">密码</label>
        <input
          id="password"
          v-model="password"
          type="password"
          placeholder="输入密码（至少6位）"
          autocomplete="new-password"
          :disabled="isSubmitting"
        />
        <p v-if="errors.password" class="error-message">{{ errors.password }}</p>
      </div>

      <div class="form-group full-width">
        <label for="confirmPassword">确认密码</label>
        <input
          id="confirmPassword"
          v-model="confirmPassword"
          type="password"
          placeholder="再次输入密码"
          autocomplete="new-password"
          :disabled="isSubmitting"
        />
        <p v-if="errors.confirmPassword" class="error-message">{{ errors.confirmPassword }}</p>
      </div>
    </div>

    <button type="submit" class="submit-button" :disabled="isSubmitting">
      {{ isSubmitting ? '注册中...' : '注册' }}
    </button>

    <div class="divider">
      <span>或</span>
    </div>

    <div class="social-login">
      <button type="button" class="social-button" @click="$router.push('/login')">
        返回登录
      </button>
    </div>
  </form>
</template>

<style scoped lang="scss">
@import '@/assets/styles/auth.scss';
</style>
