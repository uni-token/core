<script setup lang="ts">
import type { Provider } from '@/lib/providers'

defineProps<{
  provider: Provider
}>()
</script>

<template>
  <div class="flex gap-2 items-center">
    <img :src="provider.logo" :alt="provider.name" class="w-5 mb-1 rounded">
    {{ provider.name }}
  </div>
</template>
