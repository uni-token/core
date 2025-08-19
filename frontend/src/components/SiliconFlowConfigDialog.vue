<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import SiliconFlow from '@/components/SiliconFlowCard.vue'
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog'

const emit = defineEmits<{
  configured: [provider: string]
}>()

const { t } = useI18n()
const open = defineModel<boolean>('open')
</script>

<template>
  <Dialog v-model:open="open">
    <DialogContent class="sm:max-w-lg">
      <DialogHeader>
        <DialogTitle>{{ t('providers.configureSiliconFlow') }}</DialogTitle>
        <DialogDescription>
          {{ t('providers.siliconFlowDescription') }}
        </DialogDescription>
      </DialogHeader>

      <div class="my-2">
        <SiliconFlow @configured="(provider) => emit('configured', provider)" />
      </div>
    </DialogContent>
  </Dialog>
</template>
