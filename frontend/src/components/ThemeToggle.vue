<script setup lang="ts">
import type { Theme } from '@/stores/theme'
import { Monitor, Moon, Sun } from 'lucide-vue-next'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { Button } from '@/components/ui/button'
import { useThemeStore } from '@/stores/theme'

const { t } = useI18n()
const themeStore = useThemeStore()

const themeOptions: Theme[] = ['light', 'dark']

function cycleTheme() {
  const currentIndex = themeOptions.indexOf(themeStore.theme)
  const nextIndex = (currentIndex + 1) % themeOptions.length
  themeStore.setTheme(themeOptions[nextIndex])
}

const currentThemeIcon = computed(() => {
  switch (themeStore.theme) {
    case 'light':
      return Sun
    case 'dark':
      return Moon
    case 'system':
      return Monitor
    default:
      return Sun
  }
})

const currentThemeLabel = computed(() => {
  switch (themeStore.theme) {
    case 'light':
      return t('lightMode')
    case 'dark':
      return t('darkMode')
    case 'system':
      return t('systemMode')
    default:
      return t('lightMode')
  }
})
</script>

<template>
  <Button variant="ghost" size="icon" class="h-8 w-8" :title="currentThemeLabel" @click="cycleTheme">
    <component :is="currentThemeIcon" class="h-4 w-4" />
    <span class="sr-only">{{ t('theme') }}</span>
  </Button>
</template>

<i18n lang="yaml">
zh-CN:
  theme: 主题
  lightMode: 浅色模式
  darkMode: 深色模式
  systemMode: 跟随系统
en-US:
  theme: Theme
  lightMode: Light Mode
  darkMode: Dark Mode
  systemMode: System Mode
</i18n>
