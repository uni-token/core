import { computed } from 'vue'
import { createI18n } from 'vue-i18n'
import enUS from './en-US.json'
import zhCN from './zh-CN.json'

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

const messages = {
  'zh-CN': zhCN,
  'en-US': enUS,
}

const i18n = createI18n({
  legacy: false,
  locale: getInitialLocale(),
  fallbackLocale: 'zh-CN',
  messages,
})

export function setLocale(locale: string) {
  if (['zh-CN', 'en-US'].includes(locale)) {
    i18n.global.locale.value = locale as 'zh-CN' | 'en-US'
    localStorage.setItem('locale', locale)
    document.documentElement.lang = locale
  }
}

export function getCurrentLocale() {
  return computed(() => i18n.global.locale.value)
}

export default i18n
