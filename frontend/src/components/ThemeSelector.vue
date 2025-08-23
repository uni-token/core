<script setup lang="ts">
import type { Theme } from '@/stores/theme'
import { Moon, Sun } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { useThemeStore } from '@/stores/theme'

const { t } = useI18n()
const themeStore = useThemeStore()

const themeOptions = [
  { value: 'light' as const, label: t('settings.lightMode'), icon: Sun },
  { value: 'dark' as const, label: t('settings.darkMode'), icon: Moon },
  // { value: 'system' as const, label: t('settings.systemMode'), icon: Monitor },
]

function handleThemeChange(newTheme: unknown) {
  if (typeof newTheme === 'string' && ['light', 'dark', 'system'].includes(newTheme)) {
    themeStore.setTheme(newTheme as Theme)
  }
}
</script>

<template>
  <div class="flex items-center justify-between">
    <div>
      <h4 class="text-sm font-medium">
        {{ t('settings.theme') }}
      </h4>
      <p class="text-sm text-muted-foreground">
        {{ t('settings.themeDescription') }}
      </p>
    </div>
    <Select :model-value="themeStore.theme" @update:model-value="handleThemeChange">
      <SelectTrigger class="w-32">
        <SelectValue />
      </SelectTrigger>
      <SelectContent>
        <SelectItem
          v-for="option in themeOptions"
          :key="option.value"
          :value="option.value"
        >
          <div class="flex items-center gap-2">
            <component :is="option.icon" class="h-4 w-4" />
            <span>{{ option.label }}</span>
          </div>
        </SelectItem>
      </SelectContent>
    </Select>
  </div>
</template>
