<script lang="ts">
import { useScriptTag } from '@vueuse/core'
import { VisuallyHidden } from 'reka-ui'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { Dialog, DialogContent, DialogDescription, DialogTitle } from './ui/dialog'

const state = ref<'loading' | 'init' | 'pending' | 'success'>('loading')

const window_ = window as unknown as {
  initSMCaptcha: any
  _smConf: any
  _smReadyFuncs: any[]
  SMSdk: any
}

useScriptTag('https://castatic.fengkongcloud.cn/pr/v1.0.4/smcp.min.js', () => {
  state.value = 'init'
})
window_._smReadyFuncs = []
window_.SMSdk = {
  ready(fn: any) {
    fn && window_._smReadyFuncs.push(fn)
  },
}
window_._smConf = {
  organization: 'P9usCUBauxft8eAmUXaZ',
  appId: 'default',
  publicKey: 'MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDetfEgYD4aE1ZjmWJ6/jnPurhzI+yeRoJHWrnNtQMte3stQ4VjG3yu21FuN75E6cDpA9KtDXwcB2M/FiGUAe3G0rNotbWI8+SjZfUbW/OILFTzY0uaeEkmVGW5WyJ6weQbbr1xTCPa2OO3YIMeZljWUYHG5h21WAm/PATg8im8cQIDAQAB',
  staticHost: 'static.portal101.cn',
  protocol: 'https',
  apiHost: 'fp-it-acc.portal101.cn',
}
useScriptTag('https://static.portal101.cn/dist/web/v3.0.0/fp.min.js')

export const deviceId = new Promise<string>((resolve) => {
  window_.SMSdk.ready(() => {
    resolve(window_.SMSdk.getDeviceId?.() || '')
  })
})
</script>

<script setup lang="ts">
const props = defineProps<{
  enabled: boolean
  config: any
}>()
const emits = defineEmits<{
  next: [result: any]
}>()

const { t } = useI18n()

const captcha = ref<any>(null)

function onClick() {
  window_.initSMCaptcha({
    ...props.config,
    appendTo: 'sm-captcha',
  }, (c: any) => {
    captcha.value = c
    state.value = 'pending'
    c.onSuccess(async (e: any) => {
      if (e.pass) {
        state.value = 'success'
        emits('next', e.rid)
      }
      else {
        console.error('Captcha verification failed')
      }
    })
  })
}

const statusText = computed(() => ({
  loading: t('loading'),
  init: t('getCode'),
  pending: t('getCode'),
  success: t('smsSent'),
})[state.value])
</script>

<template>
  <button @click="onClick">
    {{ statusText }}
    <Dialog :open="state === 'pending'" @update:open="captcha?.reset(); state = 'init'">
      <DialogContent>
        <VisuallyHidden>
          <DialogTitle>
            Captcha
          </DialogTitle>
          <DialogDescription>
            Captcha
          </DialogDescription>
        </VisuallyHidden>
        <div class="relative m-4">
          <div id="sm-captcha" />
        </div>
      </DialogContent>
    </Dialog>
  </button>
</template>

<i18n lang="yaml">
en-US:
  loading: Loading...
  getCode: Get Code
  smsSent: Code Sent

zh-CN:
  loading: 加载中...
  getCode: 获取验证码
  smsSent: 已发送
</i18n>
