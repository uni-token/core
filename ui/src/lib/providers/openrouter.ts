import type { ProviderUserInfo } from './index'
import { defineAsyncComponent, ref } from 'vue'
import { useServiceStore } from '@/stores'
import { useI18n } from '../locals'
import { defineProvider, useProviderSession } from './index'

export const useOpenRouterProvider = defineProvider(() => {
  const { t } = useI18n({
    'zh-CN': {
      providerName: 'OpenRouter',
    },
    'en-US': {
      providerName: 'OpenRouter',
    },
  })
  const { proxy } = useServiceStore()

  const session = useProviderSession<{
    key: string
    userId: string
  }>('openrouter')

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
      const s = await session.get()
      if (!s) {
        return
      }
      const res = await proxy('https://openrouter.ai/api/v1/credits', {
        method: 'GET',
        headers: {
          Authorization: `Bearer ${s.key}`,
        },
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

    Login: defineAsyncComponent(() => import('@/components/OpenRouterLoginCard.vue')),
    async logout() {
      await session.delete()
      user.value = null
    },

    baseURL: 'https://openrouter.ai/api/v1',
    async createKey() {
      const s = await session.get()
      if (!s) {
        throw new Error('No session')
      }
      return s.key
    },

    apis: {},
  }
})
