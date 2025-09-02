<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import ManualConfigCard from './ManualConfigCard.vue'
import SiliconFlowCard from './SiliconFlowCard.vue'

const emit = defineEmits<{
  configured: [key: string]
}>()

const open = defineModel<boolean>('open')

const { t } = useI18n()

async function handleConfigured(key: string) {
  emit('configured', key)
  open.value = false
}
</script>

<template>
  <Dialog v-model:open="open">
    <DialogContent class="sm:max-w-2xl max-h-[90vh] overflow-y-auto">
      <DialogHeader>
        <DialogTitle>{{ t('title') }}</DialogTitle>
        <DialogDescription>
          {{ t('description') }}
        </DialogDescription>
      </DialogHeader>

      <div class="mt-6 space-y-6">
        <SiliconFlowCard show-title @configured="handleConfigured" />
        <ManualConfigCard @configured="handleConfigured" />
      </div>
    </DialogContent>
  </Dialog>
</template>

<i18n lang="yaml">
zh-CN:
  title: 添加 API Key
  description: 选择一个提供商来添加 API Key
en-US:
  title: Add API Key
  description: Choose a provider to add API Key
</i18n>
