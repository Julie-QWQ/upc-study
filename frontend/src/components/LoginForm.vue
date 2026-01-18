<script setup lang="ts">
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import { useLoginForm } from '@/composables/useValidator'
import { useEmailVerification } from '@/composables/useEmailVerification'
import type { LoginRequest } from '@/types'

const emit = defineEmits<{
  success: []
}>()

const { login } = useAuth()
const router = useRouter()
const {
  sendVerificationCode,
  loginWithCode,
  countdown,
  isSending,
  error: verificationError
} = useEmailVerification()

// 登录方式：username 或 email
const loginMethod = ref<'username' | 'email'>('username')

const {
  username,
  password,
  errors,
  isSubmitting,
  handleSubmit
} = useLoginForm(loginMethod)

// 验证码
const verificationCode = ref('')

const resetLoginFields = () => {
  username.value = ''
  password.value = ''
  verificationCode.value = ''
  if (errors.value) {
    errors.value.username = ''
    errors.value.password = ''
    errors.value.verificationCode = ''
  }
  verificationError.value = null
  countdown.value = 0
}

watch(loginMethod, () => {
  resetLoginFields()
})

// 发送验证码
const handleSendCode = async () => {
  if (!username.value) {
    errors.value.username = '请先输入邮箱'
    ElMessage.error(errors.value.username)
    return
  }

  // 简单的邮箱格式验证
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(username.value)) {
    errors.value.username = '请输入有效的邮箱地址'
    ElMessage.error(errors.value.username)
    return
  }

  const success = await sendVerificationCode({
    email: username.value,
    purpose: 'login'
  })

  if (success) {
    ElMessage.success('验证码已发送，请查收邮箱')
  } else if (verificationError.value) {
    ElMessage.error(verificationError.value)
  } else {
    ElMessage.error('发送验证码失败')
  }
}

const onSubmit = handleSubmit(async () => {
  // 如果选择邮箱登录，需要验证码
  if (loginMethod.value === 'email') {
    if (!verificationCode.value) {
      errors.value.verificationCode = '请输入验证码'
      ElMessage.error(errors.value.verificationCode)
      return
    }

    if (verificationCode.value.length !== 6) {
      errors.value.verificationCode = '验证码为6位数字'
      ElMessage.error(errors.value.verificationCode)
      return
    }

    // 使用验证码登录（不需要密码）
    const success = await loginWithCode({
      email: username.value,
      code: verificationCode.value
    })

    if (success) {
      emit('success')
    }
  } else {
    if (!username.value) {
      errors.value.username = '请输入用户名或邮箱'
      ElMessage.error(errors.value.username)
      return
    }
    if (!password.value) {
      errors.value.password = '请输入密码'
      ElMessage.error(errors.value.password)
      return
    }
    // 使用用户名登录，不需要验证码
    const credentials: LoginRequest = {
      username: username.value,
      password: password.value
    }

    const success = await login(credentials)
    if (success) {
      emit('success')
    }
  }
}, () => {
  ElMessage.error('请检查登录信息')
})
</script>

<template>
  <form @submit="onSubmit" class="auth-form">
    <!-- 登录方式选择 -->
    <div class="login-method-toggle">
      <button
        type="button"
        :class="['method-button', { active: loginMethod === 'username' }]"
        @click="loginMethod = 'username'"
      >
        用户名登录
      </button>
      <button
        type="button"
        :class="['method-button', { active: loginMethod === 'email' }]"
        @click="loginMethod = 'email'"
      >
        邮箱登录
      </button>
    </div>

    <div class="form-group">
      <label for="username">{{ loginMethod === 'username' ? '用户名' : '邮箱' }}</label>
      <input
        id="username"
        v-model="username"
        :type="loginMethod === 'username' ? 'text' : 'email'"
        :placeholder="loginMethod === 'username' ? '输入用户名' : '输入邮箱'"
        autocomplete="username"
        :disabled="isSubmitting"
      />
      <p v-if="errors.username" class="error-message">{{ errors.username }}</p>
    </div>

    <!-- 只有用户名登录时才显示密码输入 -->
    <div v-if="loginMethod === 'username'" class="form-group">
      <label for="password">密码</label>
      <input
        id="password"
        v-model="password"
        type="password"
        placeholder="输入密码"
        autocomplete="current-password"
        :disabled="isSubmitting"
      />
      <p v-if="errors.password" class="error-message">{{ errors.password }}</p>
    </div>

    <!-- 只有选择邮箱登录时才显示验证码输入 -->
    <div v-if="loginMethod === 'email'" class="form-group">
      <label for="verificationCode">邮箱验证码</label>
      <div class="verification-input-group">
        <input
          id="verificationCode"
          v-model="verificationCode"
          type="text"
          placeholder="输入6位验证码"
          maxlength="6"
          :disabled="isSubmitting"
        />
        <button
          type="button"
          class="send-code-button"
          :disabled="isSending || countdown > 0"
          @click="handleSendCode"
        >
          {{ countdown > 0 ? `${countdown}秒后重试` : '发送验证码' }}
        </button>
      </div>
      <p v-if="errors.verificationCode" class="error-message">{{ errors.verificationCode }}</p>
      <p v-if="verificationError" class="error-message">{{ verificationError }}</p>
    </div>

    <div class="form-actions">
      <div class="forgot-password">
        <router-link to="/forgot-password">忘记密码？</router-link>
      </div>
    </div>

    <button type="submit" class="submit-button" :disabled="isSubmitting">
      {{ isSubmitting ? '登录中...' : '登录' }}
    </button>

  </form>
</template>

<style scoped lang="scss">
@import '@/assets/styles/auth.scss';

// 登录方式切换按钮
.login-method-toggle {
  display: flex;
  gap: 8px;
  padding: 6px;
  margin-bottom: 20px;
  border-radius: 999px;
  background: rgba(15, 118, 110, 0.08);
  border: 1px solid rgba(15, 118, 110, 0.18);

  .method-button {
    flex: 1;
    padding: 10px 16px;
    height: 42px;
    background: transparent;
    border: none;
    color: rgba(15, 23, 42, 0.7);
    border-radius: 999px;
    cursor: pointer;
    font-size: 0.95rem;
    font-weight: 600;
    transition: all 0.2s ease;

    &:hover {
      color: rgba(15, 23, 42, 0.9);
    }

    &.active {
      color: #ffffff;
      background: linear-gradient(120deg, #0f766e, #14b8a6);
      box-shadow: 0 8px 18px rgba(15, 118, 110, 0.25);
    }
  }
}

.verification-input-group {
  display: flex;
  gap: 10px;

  input {
    flex: 1;
  }

  .send-code-button {
    position: relative;
    flex-shrink: 0;
    padding: 0 20px;
    height: 44px;
    // 使用与背景协调的青绿色渐变
    background: linear-gradient(
      135deg,
      rgba(15, 118, 110, 0.9) 0%,
      rgba(20, 184, 166, 0.9) 100%
    );
    backdrop-filter: blur(10px);
    -webkit-backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.3);
    color: white;
    border-radius: 8px;
    cursor: pointer;
    font-size: 14px;
    font-weight: 500;
    white-space: nowrap;
    transition: all 0.3s ease;
    overflow: hidden;
    box-shadow:
      0 4px 12px rgba(15, 118, 110, 0.3),
      inset 0 1px 0 rgba(255, 255, 255, 0.2);

    // 磨砂玻璃反光效果
    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: -100%;
      width: 100%;
      height: 100%;
      background: linear-gradient(
        90deg,
        transparent,
        rgba(255, 255, 255, 0.4),
        transparent
      );
      transition: left 0.5s ease;
    }

    // 鼠标悬停时的反光效果
    &:hover:not(:disabled) {
      background: linear-gradient(
        135deg,
        rgba(15, 118, 110, 1) 0%,
        rgba(20, 184, 166, 1) 100%
      );
      border-color: rgba(255, 255, 255, 0.4);
      box-shadow:
        0 6px 20px rgba(15, 118, 110, 0.4),
        inset 0 1px 0 rgba(255, 255, 255, 0.3);
      transform: translateY(-2px);

      &::before {
        left: 100%;
      }
    }

    // 点击效果
    &:active:not(:disabled) {
      transform: translateY(0);
      box-shadow:
        0 2px 8px rgba(15, 118, 110, 0.3),
        inset 0 1px 0 rgba(255, 255, 255, 0.2);
    }

    // 禁用状态
    &:disabled {
      background: rgba(144, 147, 153, 0.3);
      border-color: rgba(144, 147, 153, 0.2);
      color: rgba(255, 255, 255, 0.4);
      cursor: not-allowed;
      box-shadow: none;
      transform: none;

      &::before {
        display: none;
      }
    }

    // 倒计时状态
    &:not(:disabled):hover {
      animation: shimmer 1.5s infinite;
    }
  }
}

@keyframes shimmer {
  0% {
    box-shadow:
      0 4px 12px rgba(15, 118, 110, 0.2),
      inset 0 1px 0 rgba(255, 255, 255, 0.2);
  }
  50% {
    box-shadow:
      0 6px 20px rgba(15, 118, 110, 0.4),
      inset 0 1px 0 rgba(255, 255, 255, 0.3);
  }
  100% {
    box-shadow:
      0 4px 12px rgba(15, 118, 110, 0.2),
      inset 0 1px 0 rgba(255, 255, 255, 0.2);
  }
}
</style>
