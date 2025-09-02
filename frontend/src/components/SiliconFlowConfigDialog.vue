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
          <DialogTitle>{{ t('providers.siliconFlow') }}</DialogTitle>
          <DialogDescription>
            {{ t('providers.siliconFlowDescription') }}
          </DialogDescription>
        </DialogHeader>
      </VisuallyHidden>
      <SiliconFlowCard show-title @configured="(provider) => { open = false; emit('configured', provider) }" />
    </DialogContent>
  </Dialog>
</template>
