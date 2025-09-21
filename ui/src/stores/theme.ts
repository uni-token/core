import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export type Theme = 'system' | 'light' | 'dark'

export const useThemeStore = defineStore('theme', () => {
  const theme = ref<Theme>('system')
  const isDark = ref(false)

  // Initialize theme from localStorage or default to system
  function initTheme() {
    const savedTheme = localStorage.getItem('theme') as Theme | null
    if (savedTheme && ['system', 'light', 'dark'].includes(savedTheme)) {
      theme.value = savedTheme
    }
    updateTheme()
  }

  // Update the actual theme based on the selected theme
  function updateTheme() {
    if (theme.value === 'system') {
      isDark.value = window.matchMedia('(prefers-color-scheme: dark)').matches
    }
    else {
      isDark.value = theme.value === 'dark'
    }

    // Apply or remove dark class from document element
    if (isDark.value) {
      document.documentElement.classList.add('dark')
    }
    else {
      document.documentElement.classList.remove('dark')
    }
  }

  // Set theme and save to localStorage
  function setTheme(newTheme: Theme) {
    theme.value = newTheme
    localStorage.setItem('theme', newTheme)
    updateTheme()
  }

  // Listen for system theme changes
  function setupSystemThemeListener() {
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
    mediaQuery.addEventListener('change', () => {
      if (theme.value === 'system') {
        updateTheme()
      }
    })
  }

  // Watch for theme changes
  watch(theme, updateTheme)

  return {
    theme,
    isDark,
    setTheme,
    initTheme,
    setupSystemThemeListener,
  }
})
