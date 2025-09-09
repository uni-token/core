<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { getCurrentLocale, setLocale } from '@/lib/locals'

const props = defineProps<{
  compact?: boolean
}>()

const { t } = useI18n()

const currentLocale = getCurrentLocale()

const languages = [
  { value: 'zh-CN', label: '中文（简体）' },
  { value: 'en-US', label: 'English (US)' },
]

function handleLanguageChange(value: any) {
  if (value && typeof value === 'string') {
    setLocale(value)
  }
}
</script>

<template>
  <div class="space-y-2">
    <label v-if="!props.compact" class="text-sm font-medium">
      {{ t('language') }}
    </label>
    <Select
      :model-value="currentLocale"
      @update:model-value="handleLanguageChange"
    >
      <SelectTrigger class="w-full">
        <SelectValue>
          {{ languages.find(lang => lang.value === currentLocale)?.label }}
        </SelectValue>
      </SelectTrigger>
      <SelectContent>
        <SelectItem
          v-for="lang in languages"
          :key="lang.value"
          :value="lang.value"
        >
          {{ lang.label }}
        </SelectItem>
      </SelectContent>
    </Select>
  </div>
</template>

<i18n lang="yaml">
zh-CN:
  language: 语言
en-US:
  language: Language
</i18n>
