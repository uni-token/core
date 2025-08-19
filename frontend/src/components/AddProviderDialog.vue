<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import ManualConfigCard from './ManualConfigCard.vue'
import SiliconFlowCard from './SiliconFlowCard.vue'

const emit = defineEmits<{
  configured: [provider: string]
}>()

const open = defineModel<boolean>('open')

const { t } = useI18n()

async function handleConfigured(provider: string) {
  emit('configured', provider)
  open.value = false
}
</script>

<template>
  <Dialog v-model:open="open">
    <DialogContent class="sm:max-w-2xl max-h-[90vh] overflow-y-auto">
      <DialogHeader>
        <DialogTitle>{{ t('providers.addProviderDialog.title') }}</DialogTitle>
        <DialogDescription>
          {{ t('providers.addProviderDialog.description') }}
        </DialogDescription>
      </DialogHeader>

      <div class="mt-6 space-y-6">
        <SiliconFlowCard show-title @configured="handleConfigured" />
        <ManualConfigCard @configured="handleConfigured" />
      </div>
    </DialogContent>
  </Dialog>
</template>
