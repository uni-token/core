<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, shallowRef } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps<{
  enabled: boolean
  config: any
}>()
const emits = defineEmits<{
  next: [result: any]
}>()
const { t } = useI18n()
const state = ref<'loading' | 'init' | 'pending' | 'success'>('loading')

declare global {
  interface Window {
    initGeetest4: any
  }
}

const captcha = shallowRef<any>(null)

onMounted(async () => {
  if (!window.initGeetest4) {
    await new Promise((resolve) => {
      const script = document.getElementById('gt4-script-tag')
      if (script) {
        script.addEventListener('load', resolve)
      }
      else {
        throw new Error('Geetest script not found')
      }
    })
  }
  window.initGeetest4(
    props.config,
    (e: any) => {
      state.value = 'init'
      captcha.value = e
      e.onReady(() => {
        setTimeout(() => {
          document.querySelectorAll('.geetest_captcha').forEach((el) => {
            ; (el as HTMLElement).style.pointerEvents = 'auto'
            el.addEventListener('mousedown', (ev) => {
              ev.stopPropagation()
            })
            el.addEventListener('pointerdown', (ev) => {
              ev.stopPropagation()
            })
          })
        }, 300)
        e.onSuccess(() => {
          state.value = 'success'
          emits('next', e.getValidate())
        })
      })
    },
  )
})

onUnmounted(() => {
  if (captcha.value) {
    captcha.value.destroy()
    captcha.value = null
  }
})

function onClick() {
  if (!props.enabled || state.value === 'success' || !captcha.value)
    return

  state.value = 'pending'
  captcha.value.showCaptcha()
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
