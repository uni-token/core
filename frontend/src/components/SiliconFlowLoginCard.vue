<script setup lang="ts">
import { useScriptTag } from '@vueuse/core'
import { useI18n } from 'vue-i18n'
import Captcha from '@/components/Captcha.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Dialog, DialogClose, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { useSiliconFlowStore, useThemeStore } from '@/stores'

const { t, locale } = useI18n()
const siliconFlowStore = useSiliconFlowStore()
const themeStore = useThemeStore()

useScriptTag('http://res.wx.qq.com/connect/zh_CN/htmledition/js/wxLogin.js')

function wxLogin() {
  setTimeout(() => {
    // @ts-expect-error WeChat SDK
    // eslint-disable-next-line no-new
    new window.WxLogin({
      self_redirect: !1,
      id: 'SF_wx_login_qr_code_f',
      appid: 'wx637ec58e4e15a258',
      scope: 'snsapi_login',
      style: themeStore.isDark ? 'white' : 'black',
      lang: locale.value === 'zh-CN' ? 'cn' : 'en',
      // stylelite: 1,
      // fast_login: 1,
      redirect_uri: encodeURIComponent('https://account.siliconflow.cn/api/open/weixin'),
    })
  }, 500)
}
</script>

<template>
  <Card>
    <CardHeader class="pb-4">
      <CardTitle class="text-lg">
        {{ t('loginTitle') }}
      </CardTitle>
      <CardDescription class="text-sm">
        {{ siliconFlowStore.isEmailLogin ? t('emailLoginDescription') : t('loginDescription') }}
      </CardDescription>
    </CardHeader>
    <CardContent class="space-y-4">
      <!-- Phone Number Input -->
      <div v-if="!siliconFlowStore.isEmailLogin" class="flex rounded-md border border-input bg-background">
        <div class="flex items-center px-3 border-r border-input bg-muted/50 rounded-l-md">
          <span class="text-sm font-medium text-muted-foreground">+86</span>
        </div>
        <Input
          id="phone" v-model="siliconFlowStore.phoneNumber" :placeholder="t('phoneNumber')" type="tel"
          class="border-0 rounded-l-none focus-visible:ring-0 focus-visible:ring-offset-0 h-10"
        />
      </div>

      <!-- Email Input -->
      <div v-if="siliconFlowStore.isEmailLogin" class="flex rounded-md border border-input bg-background">
        <Input
          id="email" v-model="siliconFlowStore.email" :placeholder="t('emailAddressPlaceholder')" type="email"
          class="border-0 focus-visible:ring-0 focus-visible:ring-offset-0 h-10"
        />
      </div>

      <!-- SMS/Email Code Input -->
      <div class="flex rounded-md border border-input bg-background">
        <Input
          id="sms" v-model="siliconFlowStore.smsCode"
          :placeholder="siliconFlowStore.isEmailLogin ? t('emailCode') : t('smsCode')"
          type="text" maxlength="6"
          class="w-fit flex-grow border-0 rounded-r-none focus-visible:ring-0 focus-visible:ring-offset-0 h-10"
        />
        <div class="border-l border-input">
          <Captcha
            :enabled="siliconFlowStore.isEmailLogin ? siliconFlowStore.email.length > 0 : siliconFlowStore.phoneNumber.length > 0"
            :config="siliconFlowStore.captchaConfig"
            class="h-10 px-4 bg-muted/50 rounded-r-md border-0 text-xs text-primary hover:bg-muted/70 transition-colors disabled:opacity-50"
            @next="siliconFlowStore.sendSMS"
          />
        </div>
      </div>
    </CardContent>
    <CardFooter class="flex flex-col space-y-3 pt-3 items-start">
      <div class="flex items-center space-x-2 text-xs text-muted-foreground">
        <input id="agree" v-model="siliconFlowStore.agreed" type="checkbox" class="h-3 w-3 rounded border border-input">
        <label for="agree" class="flex items-center gap-1 cursor-pointer">
          <span>{{ t('agreeToTerms') }}</span>
          <a
            href="https://docs.siliconflow.cn/docs/user-agreement" target="_blank"
            class="text-primary hover:underline"
          >{{ t('userAgreement') }}</a>
          <span>{{ t('and') }}</span>
          <a
            href="https://docs.siliconflow.cn/docs/privacy-policy" target="_blank"
            class="text-primary hover:underline"
          >{{ t('privacyPolicy') }}</a>
        </label>
      </div>
      <Button class="w-full h-10" :disabled="!siliconFlowStore.canLogin" @click="siliconFlowStore.login()">
        <span v-if="siliconFlowStore.isLoading">{{ t('loggingIn') }}</span>
        <span v-else>{{ siliconFlowStore.isEmailLogin ? t('login') : t('registerLogin') }}</span>
      </Button>
      <div class="w-full grid grid-cols-2 gap-3">
        <Dialog>
          <DialogTrigger>
            <Button variant="outline" class="w-full h-10" @click="wxLogin">
              {{ t('wechatLogin') }}
            </Button>
          </DialogTrigger>
          <DialogContent class="max-w-sm">
            <DialogHeader>
              <DialogTitle>{{ t('wechatLogin') }}</DialogTitle>
              <DialogClose />
            </DialogHeader>
            <div id="SF_wx_login_qr_code_f" class="w-full flex justify-center" />
          </DialogContent>
        </Dialog>
        <Button variant="outline" class="w-full h-10" @click="siliconFlowStore.isEmailLogin = !siliconFlowStore.isEmailLogin">
          {{ siliconFlowStore.isEmailLogin ? t('phoneLogin') : t('emailLogin') }}
        </Button>
      </div>
      <!-- <div class="flex items-center space-x-2 text-xs text-muted-foreground">
        <input id="keep" v-model="siliconFlowStore.keepLogin" type="checkbox" class="h-3 w-3 rounded border border-input" checked>
        <label for="keep" class="cursor-pointer">{{ t('keepLoggedIn30Days') }}</label>
      </div> -->
    </CardFooter>
  </Card>
</template>

<i18n lang="yaml">
en-US:
  loginTitle: Login to SiliconFlow
  loginDescription: Login using phone number and SMS verification code
  emailLoginDescription: Login using email address and verification code
  phoneNumber: Your phone number
  emailAddressPlaceholder: Your email address
  smsCode: SMS verification code
  emailCode: Email verification code
  agreeToTerms: I agree to SiliconFlow's
  userAgreement: User Agreement
  and: and
  privacyPolicy: Privacy Policy
  loggingIn: Logging in...
  login: Login
  registerLogin: Register/Login
  wechatLogin: WeChat
  phoneLogin: Phone Login
  emailLogin: Email Login
  keepLoggedIn30Days: Keep logged in for 30 days

zh-CN:
  loginTitle: 登录 硅基流动
  loginDescription: 使用手机号码和短信验证码登录
  emailLoginDescription: 使用邮箱地址和验证码登录
  phoneNumber: 您的手机号
  emailAddressPlaceholder: 您的邮箱地址
  smsCode: 短信验证码
  emailCode: 邮箱验证码
  agreeToTerms: 我同意硅基流动的
  userAgreement: 用户协议
  and: 和
  privacyPolicy: 隐私政策
  loggingIn: 登录中...
  login: 登录
  registerLogin: 注册/登录
  wechatLogin: 微信登录
  phoneLogin: 短信登录
  emailLogin: 邮箱登录
  keepLoggedIn30Days: 30天内保持登录
</i18n>
