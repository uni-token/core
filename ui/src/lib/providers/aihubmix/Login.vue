<script setup lang="ts">
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue-sonner'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { useKeysStore } from '@/stores'
import { useAIHubMixProvider } from '.'

const { t } = useI18n()
const keysStore = useKeysStore()
const provider = useAIHubMixProvider()

const isRegisterMode = ref(false)
const username = ref('')
const password = ref('')
const confirmPassword = ref('')
const isLoading = ref(false)

const passwordError = computed(() => {
  if (isRegisterMode.value && password.value && confirmPassword.value) {
    if (password.value !== confirmPassword.value)
      return t('passwordsDoNotMatch')
  }
  if (isRegisterMode.value && password.value) {
    if (password.value.length < 8 || password.value.length > 20)
      return t('passwordLengthError')
  }
  return ''
})

const canSubmit = computed(() => {
  if (isLoading.value)
    return false
  if (isRegisterMode.value)
    return username.value && password.value && confirmPassword.value && !passwordError.value
  else
    return username.value && password.value
})

async function submitForm() {
  if (!canSubmit.value)
    return

  try {
    isLoading.value = true
    if (isRegisterMode.value) {
      // Register logic
      await provider.apis.register(username.value, password.value)
      toast.success(t('registerSuccess'))
      // Switch to login mode after successful registration
      isRegisterMode.value = false
    }
    else {
      // Login logic
      await provider.apis.login(username.value, password.value)
      toast.success(t('loginSuccess'))
    }

    await provider.refreshUser()

    // Automatically create API key after successful login
    if (!isRegisterMode.value)
      await keysStore.createAndAddKey(provider)

    // Clear form
    username.value = ''
    password.value = ''
    confirmPassword.value = ''
  }
  catch (error) {
    console.error('Auth error:', error)
    toast.error(isRegisterMode.value ? t('registerFailed') : t('loginFailed'))
  }
  finally {
    isLoading.value = false
  }
}
</script>

<template>
  <Card>
    <CardHeader class="pb-4">
      <CardTitle class="text-lg">
        {{ isRegisterMode ? t('registerTitle') : t('loginTitle') }}
      </CardTitle>
      <CardDescription class="text-sm">
        {{ isRegisterMode ? t('registerDescription') : t('loginDescription') }}
      </CardDescription>
    </CardHeader>
    <CardContent class="space-y-4">
      <!-- Username Input -->
      <div class="flex rounded-md border border-input bg-background">
        <Input
          id="username" v-model="username" :placeholder="t('usernamePlaceholder')" type="text"
          class="border-0 focus-visible:ring-0 focus-visible:ring-offset-0 h-10"
        />
      </div>

      <!-- Password Input -->
      <div class="flex rounded-md border border-input bg-background">
        <Input
          id="password" v-model="password" :placeholder="t('passwordPlaceholder')" type="password"
          class="border-0 focus-visible:ring-0 focus-visible:ring-offset-0 h-10"
        />
      </div>

      <!-- Confirm Password Input (Register Mode) -->
      <div v-if="isRegisterMode" class="flex flex-col">
        <div class="flex rounded-md border border-input bg-background">
          <Input
            id="confirmPassword" v-model="confirmPassword" :placeholder="t('confirmPasswordPlaceholder')"
            type="password" class="border-0 focus-visible:ring-0 focus-visible:ring-offset-0 h-10"
          />
        </div>
        <p v-if="passwordError" class="text-xs text-red-500 mt-1.5">
          {{ passwordError }}
        </p>
      </div>
    </CardContent>
    <CardFooter class="flex flex-col space-y-3 pt-3 items-start">
      <Button class="w-full h-10" :disabled="!canSubmit" @click="submitForm()">
        <span v-if="isLoading">{{ isRegisterMode ? t('registering') : t('loggingIn') }}</span>
        <span v-else>{{ isRegisterMode ? t('register') : t('login') }}</span>
      </Button>
      <div class="w-full grid grid-cols-1">
        <Button variant="link" class="w-full h-10 text-sm" @click="isRegisterMode = !isRegisterMode">
          {{ isRegisterMode ? t('switchToLogin') : t('switchToRegister') }}
        </Button>
      </div>
    </CardFooter>
  </Card>
</template>

<i18n lang="yaml">
en-US:
  loginTitle: Login to AIHubMix
  loginDescription: Enter your username and password to log in
  registerTitle: Create an Account
  registerDescription: Create a new AIHubMix account
  usernamePlaceholder: Username
  passwordPlaceholder: Password
  confirmPasswordPlaceholder: Confirm Password
  passwordsDoNotMatch: Passwords do not match
  passwordLengthError: Password must be 8-20 characters long
  loggingIn: Logging in...
  registering: Creating account...
  login: Login
  register: Register
  switchToLogin: Already have an account? Login
  switchToRegister: Don't have an account? Register
  loginSuccess: Login successful!
  registerSuccess: Registration successful! Please log in.
  loginFailed: Login failed. Please check your credentials.
  registerFailed: Registration failed. The username may already exist.

zh-CN:
  loginTitle: 登录 AIHubMix
  loginDescription: 输入您的账户名和密码进行登录
  registerTitle: 创建新账户
  registerDescription: 创建一个新的 AIHubMix 账户
  usernamePlaceholder: 账户名
  passwordPlaceholder: 密码
  confirmPasswordPlaceholder: 确认密码
  passwordsDoNotMatch: 两次输入的密码不一致
  passwordLengthError: 密码长度必须为 8-20 个字符
  loggingIn: 登录中...
  registering: 注册中...
  login: 登录
  register: 注册
  switchToLogin: 已有账户？前往登录
  switchToRegister: 没有账户？前往注册
  loginSuccess: 登录成功！
  registerSuccess: 注册成功！请登录。
  loginFailed: 登录失败，请检查您的账户名和密码。
  registerFailed: 注册失败，账户名可能已被占用。
</i18n>
