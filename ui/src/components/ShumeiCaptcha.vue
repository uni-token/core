<script setup lang="ts">
import { useScriptTag } from '@vueuse/core'
import { VisuallyHidden } from 'reka-ui'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { Dialog, DialogContent, DialogDescription, DialogTitle } from './ui/dialog'

const props = defineProps<{
  enabled: boolean
  config: any
}>()
const emits = defineEmits<{
  next: [result: any]
}>()

const { t } = useI18n()

const state = ref<'loading' | 'init' | 'pending' | 'success'>('loading')
const captcha = ref<any>(null)

useScriptTag('https://castatic.fengkongcloud.cn/pr/v1.0.4/smcp.min.js', () => {
  state.value = 'init'
})

function onClick() {
  // @ts-expect-error global
  window.initSMCaptcha({
    ...props.config,
    appendTo: 'sm-captcha',
  }, (c: any) => {
    captcha.value = c
    state.value = 'pending'
    c.onSuccess((e: any) => {
      if (e.pass) {
        state.value = 'success'
        emits('next', e.rid)
      }
      else {
        // Failed
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
