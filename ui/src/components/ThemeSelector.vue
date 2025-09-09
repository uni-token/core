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
  { value: 'light' as const, label: t('lightMode'), icon: Sun },
  { value: 'dark' as const, label: t('darkMode'), icon: Moon },
  // { value: 'system' as const, label: t('systemMode'), icon: Monitor },
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
        {{ t('theme') }}
      </h4>
      <p class="text-sm text-muted-foreground">
        {{ t('themeDescription') }}
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

<i18n lang="yaml">
zh-CN:
  theme: 主题
  themeDescription: 选择您偏好的颜色主题
  lightMode: 浅色模式
  darkMode: 深色模式
  systemMode: 跟随系统
en-US:
  theme: Theme
  themeDescription: Choose your preferred color theme
  lightMode: Light Mode
  darkMode: Dark Mode
  systemMode: System Mode
</i18n>
