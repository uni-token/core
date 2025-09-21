<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import LanguageSelector from '@/components/LanguageSelector.vue'
import ServiceStatus from '@/components/ServiceStatus.vue'
import StartServiceButton from '@/components/StartServiceButton.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { useAuthStore, useServiceStore } from '@/stores'

const props = defineProps<{
  register?: boolean
}>()

const { t } = useI18n()
const serverStore = useServiceStore()
const authStore = useAuthStore()

const username = ref('')
const password = ref('')
const confirmPassword = ref('')
const errorMessage = ref('')

async function handleSubmit() {
  if (!username.value || !password.value) {
    errorMessage.value = t('fillCompleteInfo')
    return
  }

  if (props.register) {
    if (!confirmPassword.value) {
      errorMessage.value = t('confirmPasswordRequired')
      return
    }
    if (password.value !== confirmPassword.value) {
      errorMessage.value = t('passwordMismatch')
      return
    }
  }

  errorMessage.value = ''

  try {
    const result = props.register
      ? await authStore.register(username.value, password.value)
      : await authStore.login(username.value, password.value)

    if (result.status !== 'success') {
      errorMessage.value = result.message || t('operationFailed')
    }
  }
  catch (error) {
    errorMessage.value = `${t('networkError')}: ${error instanceof Error ? error.message : t('unknownError')}`
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-background p-4">
    <Card class="w-full max-w-md">
      <CardHeader class="text-center">
        <CardTitle class="text-2xl font-bold">
          {{ serverStore.serverConnected || authStore.hasSavedCredentials ? props.register ? t('register') : t('login') : `${t('login')}/${t('register')}` }}
        </CardTitle>
        <CardDescription>
          {{ props.register ? t('createAccount') : t('enterCredentials') }}
        </CardDescription>
        <!-- Service Status -->
        <div class="pt-2">
          <ServiceStatus variant="compact" />
        </div>
      </CardHeader>
      <CardContent>
        <form
          class="space-y-4"
          :aria-label="props.register ? t('registerForm') : t('loginForm')"
          novalidate
          @submit.prevent="handleSubmit"
        >
          <div class="space-y-2">
            <label for="username" class="text-sm font-medium">{{ t('username') }}</label>
            <Input
              id="username"
              v-model="username"
              type="text"
              :placeholder="t('usernamePlaceholder')"
              required
              :disabled="authStore.isLoading"
              autocomplete="username"
              :aria-invalid="errorMessage ? 'true' : 'false'"
              :aria-describedby="errorMessage ? 'error-message' : undefined"
            />
          </div>
          <div class="space-y-2">
            <label for="password" class="text-sm font-medium">{{ t('password') }}</label>
            <Input
              id="password"
              v-model="password"
              type="password"
              :placeholder="t('passwordPlaceholder')"
              required
              :disabled="authStore.isLoading"
              :autocomplete="props.register ? 'new-password' : 'current-password'"
              :aria-invalid="errorMessage ? 'true' : 'false'"
              :aria-describedby="errorMessage ? 'error-message' : undefined"
            />
          </div>
          <div v-if="props.register" class="space-y-2">
            <label for="confirmPassword" class="text-sm font-medium">{{ t('confirmPassword') }}</label>
            <Input
              id="confirmPassword"
              v-model="confirmPassword"
              type="password"
              :placeholder="t('confirmPasswordPlaceholder')"
              required
              :disabled="authStore.isLoading"
              autocomplete="new-password"
              :aria-invalid="errorMessage ? 'true' : 'false'"
              :aria-describedby="errorMessage ? 'error-message' : undefined"
              :aria-label="t('confirmPassword')"
            />
          </div>
          <Button
            v-if="authStore.isLoading || serverStore.serverConnected"
            type="submit"
            class="w-full"
            :disabled="authStore.isLoading"
            :aria-label="authStore.isLoading ? t('processingWait') : (props.register ? t('submitRegister') : t('submitLogin'))"
          >
            <div v-if="authStore.isLoading" class="flex items-center gap-2">
              <div
                class="w-4 h-4 border-2 border-current border-t-transparent rounded-full animate-spin"
                aria-hidden="true"
              />
              {{ t('processing') }}
            </div>
            <span v-else>{{ props.register ? t('register') : t('login') }}</span>
          </Button>
          <StartServiceButton v-else>
            <Button
              type="button"
              class="w-full opacity-50"
            >
              {{ t('disconnected') }}
            </Button>
          </StartServiceButton>
          <div
            v-if="errorMessage"
            id="error-message"
            class="text-sm text-destructive text-center"
            role="alert"
            aria-live="polite"
          >
            {{ errorMessage }}
          </div>
        </form>
      </CardContent>
    </Card>

    <div class="fixed right-4 bottom-4">
      <LanguageSelector compact />
    </div>
  </div>
</template>

<i18n lang="yaml">
zh-CN:
  login: 登录
  register: 注册
  createAccount: 创建新账户
  enterCredentials: 输入您的凭据
  loginForm: 登录表单
  registerForm: 注册表单
  username: 用户名
  usernamePlaceholder: 输入用户名
  password: 密码
  passwordPlaceholder: 输入密码
  confirmPassword: 确认密码
  confirmPasswordPlaceholder: 再次输入密码
  processing: 处理中...
  processingWait: 处理中，请稍候
  submitLogin: 提交登录
  submitRegister: 提交注册
  fillCompleteInfo: 请填写完整信息
  confirmPasswordRequired: 请输入确认密码
  passwordMismatch: 两次输入的密码不一致
  operationFailed: 操作失败
  networkError: 网络错误
  unknownError: 未知错误
en-US:
  login: Login
  register: Register
  createAccount: Create a new account
  enterCredentials: Enter your credentials
  loginForm: Login form
  registerForm: Register form
  username: Username
  usernamePlaceholder: Enter username
  password: Password
  passwordPlaceholder: Enter password
  confirmPassword: Confirm Password
  confirmPasswordPlaceholder: Enter password again
  processing: Processing...
  processingWait: Processing, please wait
  submitLogin: Submit login
  submitRegister: Submit registration
  fillCompleteInfo: Please fill in complete information
  confirmPasswordRequired: Please enter confirmation password
  passwordMismatch: Passwords do not match
  operationFailed: Operation failed
  networkError: Network error
  unknownError: Unknown error
</i18n>
