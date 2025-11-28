<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import ProvidersSection from '@/components/ProvidersSection.vue'

const { t } = useI18n()
</script>

<template>
  <div class="p-6 flex flex-col gap-6 h-full">
    <div class="flex items-center justify-between">
      <h2 class="text-2xl font-bold">
        {{ t('title') }}
      </h2>
    </div>
    <ProvidersSection />
  </div>
</template>

<i18n lang="yaml">
zh-CN:
  title: 模型管理

en-US:
  title: Model Providers
</i18n>
