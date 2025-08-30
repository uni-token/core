<script setup lang="ts">
import SiliconFlowCard from '@/components/SiliconFlowCard.vue'
import { Dialog, DialogContent } from '@/components/ui/dialog'

const emit = defineEmits<{
  configured: [provider: string]
}>()

const open = defineModel<boolean>('open')
</script>

<template>
  <Dialog v-model:open="open">
    <DialogContent class="sm:max-w-lg p-0 border-none">
      <SiliconFlowCard show-title @configured="(provider) => emit('configured', provider)" />
    </DialogContent>
  </Dialog>
</template>
