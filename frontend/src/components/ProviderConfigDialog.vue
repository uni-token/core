<script setup lang="ts">
import type { Provider } from '@/lib/providers'
import { VisuallyHidden } from 'reka-ui'
import { useI18n } from 'vue-i18n'
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import ProviderCard from './ProviderCard.vue'

const props = defineProps<{
  provider: Provider
}>()

const emit = defineEmits<{
  close: []
  configured: [key: string]
}>()

const { t } = useI18n()
</script>

<template>
  <Dialog open @update:open="emit('close')">
    <DialogContent class="sm:max-w-lg p-0 border-none">
      <VisuallyHidden>
        <DialogHeader>
          <DialogTitle>{{ provider.name }}</DialogTitle>
          <DialogDescription>
            {{ t('description', [provider.name]) }}
          </DialogDescription>
        </DialogHeader>
      </VisuallyHidden>
      <ProviderCard :provider="props.provider" @configured="(key) => { emit('close'); emit('configured', key) }" />
    </DialogContent>
  </Dialog>
</template>

<i18n lang="yaml">
en-US:
  description: Purchase and configure API through {0}

zh-CN:
  description: 通过 {0} 购买和配置 API
</i18n>
