<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog'

const { t } = useI18n()

const query = new URLSearchParams(window.location.search)
const open = ref(true)

onMounted(() => {
  const bc = new BroadcastChannel('openrouter-auth')
  bc.postMessage({ code: query.get('code') })
  setTimeout(() => {
    window.close()
  }, 100)
})
</script>

<template>
  <Dialog v-model:open="open">
    <DialogContent>
      <DialogHeader>
        <DialogTitle>{{ t('authorizationSuccessful') }}</DialogTitle>
        <DialogDescription>
          {{ t('authorizationSuccessfulDescription') }}
        </DialogDescription>
      </DialogHeader>
      <div class="flex justify-center mt-6">
        <Button @click="open = false">
          {{ t('close') }}
        </Button>
      </div>
    </DialogContent>
  </Dialog>
</template>

<i18n lang="yaml">
en-US:
  authorizationSuccessful: Authorization Successful
  authorizationSuccessfulDescription: Authorization completed successfully. Please close this page.
  close: Close

zh-CN:
  authorizationSuccessful: 授权成功
  authorizationSuccessfulDescription: 授权成功。请关闭该页面。
  close: 关闭
</i18n>
