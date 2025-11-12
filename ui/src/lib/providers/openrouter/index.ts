import type { ProviderUserInfo } from '@/lib/providers'
import { defineAsyncComponent, ref } from 'vue'
import { useI18n } from '@/lib/locals'
import { defineProvider, useProviderSession } from '@/lib/providers'

export const useOpenRouterProvider = defineProvider(() => {
  const { t } = useI18n({
    'zh-CN': {
      providerName: 'OpenRouter',
    },
    'en-US': {
      providerName: 'OpenRouter',
    },
  })

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
    get logo() {
      return import.meta.resolve('./logo.png')
    },
    homepage: 'https://openrouter.com/',

    get user() {
      return user.value
    },
    async refreshUser() {
      const s = await session.get()
      if (!s) {
        user.value = null
        return
      }

      const res = await fetch('https://openrouter.ai/api/v1/credits', {
        method: 'GET',
        headers: {
          Authorization: `Bearer ${s.key}`,
        },
      })
      if (!res.ok) {
        user.value = null
      }

      const data = await res.json()
      user.value = {
        name: s.userId.slice(8),
        balance: {
          amount: data.data.total_credits,
          currency: 'USD',
        },
      }
    },

    Login: defineAsyncComponent(() => import('@/lib/providers/openrouter/Login.vue')),
    async logout() {
      await session.delete()
      user.value = null
    },

    payment: {
      websiteURL: 'https://openrouter.ai/settings/credits',
    },

    baseURL: 'https://openrouter.ai/api/v1',
    async createKey() {
      const s = await session.get()
      if (!s) {
        throw new Error('No session')
      }
      return s.key
    },

    apis: {
      session,
    },
  }
})
