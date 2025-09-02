import { computed } from 'vue'
import { createI18n } from 'vue-i18n'

function getBrowserLocale(): string {
  const browserLocale = navigator.language || navigator.languages[0]

  const supportedLocales = ['zh-CN', 'en-US']

  if (supportedLocales.includes(browserLocale)) {
    return browserLocale
  }

  const languageCode = browserLocale.split('-')[0]
  const matchedLocale = supportedLocales.find(locale =>
    locale.startsWith(languageCode),
  )

  return matchedLocale || 'zh-CN'
}

function getInitialLocale(): string {
  const savedLocale = localStorage.getItem('locale')
  if (savedLocale && ['zh-CN', 'en-US'].includes(savedLocale)) {
    return savedLocale
  }
  return getBrowserLocale()
}

const i18n = createI18n({
  legacy: false,
  locale: getInitialLocale(),
  fallbackLocale: 'zh-CN',
})

export function setLocale(locale: string) {
  if (['zh-CN', 'en-US'].includes(locale)) {
    i18n.global.locale.value = locale as 'zh-CN' | 'en-US'
    localStorage.setItem('locale', locale)
    document.documentElement.lang = locale
  }
}

export function getCurrentLocale() {
  return computed(() => i18n.global.locale.value as 'zh-CN' | 'en-US')
}

export default i18n

export function useI18n<T extends Record<string, string>>(messages: {
  [L in 'zh-CN' | 'en-US']: T
}) {
  return {
    t: (key: keyof T) => messages[getCurrentLocale().value][key],
  }
}
