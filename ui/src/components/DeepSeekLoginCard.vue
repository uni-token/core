<script setup lang="ts">
import { computed, onUnmounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue-sonner'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { useDeepSeekProvider } from '@/lib/providers/deepseek'
import { useKeysStore, useServiceStore } from '@/stores'
import ShumeiCaptcha from './ShumeiCaptcha.vue'

const { t, locale } = useI18n()
const keysStore = useKeysStore()
const provider = useDeepSeekProvider()
const { api: fetch } = useServiceStore()

const captchaConfig = {
  organization: 'P9usCUBauxft8eAmUXaZ',
  maskBindClose: !1,
  lang: locale.value === 'zh-CN' ? 'zh-cn' : 'en',
  mode: 'spatial_select',
  domains: ['captcha.fengkongcloud.com'],
}

const phoneNumber = ref('')
const smsCode = ref('')
const agreed = ref(true)
const isLoading = ref(false)
const isSendingCode = ref(false)
const countdown = ref(0)

const canSendCode = computed(() => {
  return phoneNumber.value && !isSendingCode.value && countdown.value === 0
})

const canLogin = computed(() => {
  return phoneNumber.value && smsCode.value && agreed.value && !isLoading.value
})

let countdownInterval: NodeJS.Timeout | null = null

function startCountdown() {
  countdown.value = 60
  countdownInterval = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      clearInterval(countdownInterval!)
      countdownInterval = null
    }
  }, 1000)
}

async function sendSMS(rid: string, device_id: string) {
  if (!canSendCode.value)
    return

  try {
    isSendingCode.value = true

    // Send SMS with captcha verification result
    const res = await fetch('deepseek/sms', {
      body: JSON.stringify({
        locale: locale.value === 'zh-CN' ? 'zh_CN' : 'en_US',
        turnstile_token: '',
        shumei_verification: { region: 'CN', rid },
        device_id,
        scenario: 'login',
        mobile_number: phoneNumber.value,
      }),
      method: 'POST',
    })

    if (res.ok) {
      const { data } = await res.json()
      if (data.code === 0) {
        toast.success(t('smsSent'))
        startCountdown()
      }
      else {
        toast.error(t('smsFailure'))
      }
    }
    else {
      toast.error(t('smsFailure'))
    }
  }
  catch (error) {
    console.error('Send SMS error:', error)
    toast.error(t('smsFailure'))
  }
  finally {
    isSendingCode.value = false
  }
}

async function login() {
  if (!canLogin.value)
    return

  try {
    isLoading.value = true

    const res = await fetch('deepseek/login', {
      body: JSON.stringify({
        phone: phoneNumber.value,
        code: smsCode.value,
        area_code: '+86',
      }),
      method: 'POST',
    })

    if (res.ok) {
      await provider.refreshUser()
      // Clear login form
      phoneNumber.value = ''
      smsCode.value = ''
      toast.success(t('loginSuccess'))

      // Automatically create API key after successful login
      await keysStore.createAndAddKey(provider)
    }
    else {
      const errorData = await res.json()
      toast.error(errorData.message || t('loginFailure'))
    }
  }
  catch (error) {
    console.error('Login error:', error)
    toast.error(t('networkError'))
  }
  finally {
    isLoading.value = false
  }
}

function cleanup() {
  if (countdownInterval) {
    clearInterval(countdownInterval)
    countdownInterval = null
  }
}

// Cleanup on unmount
onUnmounted(cleanup)
</script>

<template>
  <Card>
    <CardHeader class="pb-4">
      <CardTitle class="text-lg">
        {{ t('loginTitle') }}
      </CardTitle>
      <CardDescription class="text-sm">
        {{ t('loginDescription') }}
      </CardDescription>
    </CardHeader>
    <CardContent class="space-y-4">
      <!-- Phone Number Input -->
      <div class="flex rounded-md border border-input bg-background">
        <div class="flex items-center px-3 border-r border-input bg-muted/50 rounded-l-md">
          <span class="text-sm font-medium text-muted-foreground">+86</span>
        </div>
        <Input
          id="phone"
          v-model="phoneNumber"
          :placeholder="t('phoneNumber')"
          type="tel"
          class="border-0 rounded-l-none focus-visible:ring-0 focus-visible:ring-offset-0 h-10"
          :disabled="isLoading"
        />
      </div>

      <!-- SMS Code Input -->
      <div class="flex rounded-md border border-input bg-background">
        <Input
          id="sms"
          v-model="smsCode"
          :placeholder="t('smsCode')"
          type="text"
          maxlength="6"
          class="w-fit flex-grow border-0 rounded-r-none focus-visible:ring-0 focus-visible:ring-offset-0 h-10"
          :disabled="isLoading"
        />
        <div class="border-l border-input">
          <ShumeiCaptcha
            :enabled="phoneNumber.length > 0 && !isSendingCode && countdown === 0"
            :config="captchaConfig"
            class="h-10 px-4 bg-muted/50 rounded-r-md border-0 text-xs text-primary hover:bg-muted/70 transition-colors disabled:opacity-50"
            @next="sendSMS"
          >
            {{ isSendingCode ? t('sending') : countdown > 0 ? `${countdown}s` : t('getCode') }}
          </ShumeiCaptcha>
        </div>
      </div>
    </CardContent>
    <CardFooter class="flex flex-col space-y-3 pt-3 items-start">
      <div class="flex items-center space-x-2 text-xs text-muted-foreground">
        <input
          id="agree"
          v-model="agreed"
          type="checkbox"
          class="h-3 w-3 rounded border border-input"
          :disabled="isLoading"
        >
        <label for="agree" class="flex items-center gap-1 cursor-pointer">
          <span>{{ t('agreeToTerms') }}</span>
          <a
            href="https://www.deepseek.com/zh/terms"
            target="_blank"
            class="text-primary hover:underline"
          >{{ t('userAgreement') }}</a>
          <span>{{ t('and') }}</span>
          <a
            href="https://www.deepseek.com/zh/privacy"
            target="_blank"
            class="text-primary hover:underline"
          >{{ t('privacyPolicy') }}</a>
        </label>
      </div>
      <Button
        class="w-full h-10"
        :disabled="!canLogin"
        @click="login"
      >
        <span v-if="isLoading">{{ t('loggingIn') }}</span>
        <span v-else>{{ t('registerLogin') }}</span>
      </Button>
    </CardFooter>
  </Card>
</template>

<i18n lang="yaml">
en-US:
  loginTitle: Login to DeepSeek
  loginDescription: Login using phone number and SMS verification code
  phoneNumber: Your phone number
  smsCode: SMS verification code
  agreeToTerms: I agree to DeepSeek's
  userAgreement: User Agreement
  and: and
  privacyPolicy: Privacy Policy
  loggingIn: Logging in...
  registerLogin: Register/Login
  getCode: Get Code
  sending: Sending...
  smsSent: SMS sent successfully
  smsFailure: Failed to send SMS
  loginSuccess: Login successful
  loginFailure: Login failed
  networkError: Network error

zh-CN:
  loginTitle: 登录 DeepSeek
  loginDescription: 使用手机号码和短信验证码登录
  phoneNumber: 您的手机号
  smsCode: 短信验证码
  agreeToTerms: 我同意 DeepSeek 的
  userAgreement: 用户协议
  and: 和
  privacyPolicy: 隐私政策
  loggingIn: 登录中...
  registerLogin: 注册/登录
  getCode: 获取验证码
  sending: 发送中...
  smsSent: 短信发送成功
  smsFailure: 短信发送失败
  loginSuccess: 登录成功
  loginFailure: 登录失败
  networkError: 网络错误
</i18n>
