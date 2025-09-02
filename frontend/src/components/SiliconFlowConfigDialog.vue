<script setup lang="ts">
import { VisuallyHidden } from 'reka-ui'
import { useI18n } from 'vue-i18n'
import SiliconFlowCard from '@/components/SiliconFlowCard.vue'
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog'

const emit = defineEmits<{
  configured: [provider: string]
}>()

const open = defineModel<boolean>('open')

const { t } = useI18n()
</script>

<template>
  <Dialog v-model:open="open">
    <DialogContent class="sm:max-w-lg p-0 border-none">
      <VisuallyHidden>
        <DialogHeader>
          <DialogTitle>{{ t('siliconFlow') }}</DialogTitle>
          <DialogDescription>
            {{ t('siliconFlowDescription') }}
          </DialogDescription>
        </DialogHeader>
      </VisuallyHidden>
      <SiliconFlowCard show-title @configured="(provider) => { open = false; emit('configured', provider) }" />
    </DialogContent>
  </Dialog>
</template>

<i18n lang="yaml">
en-US:
  siliconFlow: SiliconFlow
  siliconFlowDescription: Purchase and configure API through SiliconFlow

zh-CN:
  siliconFlow: 硅基流动
  siliconFlowDescription: 通过 硅基流动 购买和配置 API
</i18n>
