import type { Provider, ProviderUserInfo } from './index'
import { createSharedComposable } from '@vueuse/core'
import { markRaw, ref } from 'vue'
import OpenRouterLoginCard from '@/components/OpenRouterLoginCard.vue'
import { useServiceStore } from '@/stores'
import { useI18n } from '../locals'

export const useOpenRouterProvider = createSharedComposable((): Provider => {
  const { t } = useI18n({
    'zh-CN': {
      providerName: 'OpenRouter',
    },
    'en-US': {
      providerName: 'OpenRouter',
    },
  })
  const { fetch } = useServiceStore()

  const user = ref<null | ProviderUserInfo>()

  return {
    id: 'openrouter',
    get name() {
      return t('providerName')
    },
    homepage: 'https://openrouter.com/',

    get user() {
      return user.value
    },
    async refreshUser() {
      const res = await fetch('openrouter/status', {
        method: 'GET',
      })

      if (res.ok) {
        const data = await res.json()
        if (data.userId) {
          user.value = {
            name: data.userId.slice(8),
            balance: data.credits,
          }
          return
        }
      }
      user.value = null
    },

    Login: markRaw(OpenRouterLoginCard),
    async logout() {
      const res = await fetch('openrouter/logout', {
        method: 'POST',
      })

      if (!res.ok) {
        throw new Error('Logout failed')
      }
      user.value = null
    },

    baseURL: 'https://openrouter.ai/api/v1',
    async createKey() {
      const res = await fetch('openrouter/key')

      if (res.ok) {
        const data = await res.json()
        return data.key
      }
      throw new Error('API Key creation failed')
    },
  }
})
