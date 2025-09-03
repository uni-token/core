import type { Provider, ProviderUserInfo } from './index'
import { createSharedComposable } from '@vueuse/core'
import { ref } from 'vue'
import SiliconFlowLoginCard from '@/components/SiliconFlowLoginCard.vue'
import { useServiceStore } from '@/stores'
import { useI18n } from '../locals'

export const useSiliconFlowProvider = createSharedComposable((): Provider => {
  const { t } = useI18n({
    'zh-CN': {
      providerName: '硅基流动',
      mainlandIdCard: '中国大陆二代居民身份证',
      hkMacaoPass: '港澳居民来往内地通行证',
      taiwanPass: '台湾居民来往内地通行证',
      hkMacaoResidence: '港澳居民居住证',
      taiwanResidence: '台湾居民居住证',
      foreignerPermit: '外国人永久居留证',
      otherType: '其他类型用户',
    },
    'en-US': {
      providerName: 'SiliconFlow',
      mainlandIdCard: 'Mainland China ID Card',
      hkMacaoPass: 'Hong Kong/Macau Resident Travel Permit',
      taiwanPass: 'Taiwan Resident Travel Permit',
      hkMacaoResidence: 'Hong Kong/Macau Residence Permit',
      taiwanResidence: 'Taiwan Residence Permit',
      foreignerPermit: 'Foreigner Permanent Residence Permit',
      otherType: 'Other Type',
    },
  })
  const { fetch } = useServiceStore()

  const user = ref<null | ProviderUserInfo>()

  return {
    id: 'siliconflow',
    get name() {
      return t('providerName')
    },

    get user() {
      return user.value
    },
    async refreshUser() {
      const res = await fetch('siliconflow/status', {
        method: 'GET',
      })

      if (res.ok) {
        const data = await res.json()
        if (data.code === 20000 && data.status && data.data) {
          user.value = {
            name: data.data.name,
            verified: data.data.auth === 1,
            phone: data.data.phone,
            email: data.data.email,
            balance: data.data.balance,
          }
          return
        }
      }
      user.value = null
    },

    Login: SiliconFlowLoginCard,
    async logout() {
      const res = await fetch('siliconflow/logout', {
        method: 'POST',
      })

      if (!res.ok) {
        throw new Error('Logout failed')
      }
    },

    verification: {
      async check() {
        const res = await fetch('siliconflow/auth/info', {
          method: 'GET',
        })

        if (res.ok) {
          const data = await res.json()
          if (data.code === 20000 && data.status && data.data?.auth) {
            let name = data.data.username
            if (name.length > 1) {
              name = name.charAt(0) + '*'.repeat(name.length - 1)
            }

            let cardId = data.data.cardId
            if (cardId.length > 8) {
              cardId = cardId.substring(0, 6) + '*'.repeat(cardId.length - 10) + cardId.substring(cardId.length - 4)
            }

            return {
              name,
              cardId,
              time: data.data.authTime?.seconds ? data.data.authTime.seconds * 1000 : undefined,
            }
          }
        }
        return null
      },

      get cardTypes() {
        return [
          { value: 1, label: t('mainlandIdCard') },
          { value: 2, label: t('hkMacaoPass') },
          { value: 3, label: t('taiwanPass') },
          { value: 4, label: t('hkMacaoResidence') },
          { value: 5, label: t('taiwanResidence') },
          { value: 6, label: t('foreignerPermit') },
          { value: 100, label: t('otherType') },
        ]
      },

      async submit(data) {
        const res = await fetch('siliconflow/auth/save', {
          body: JSON.stringify({
            username: data.name.trim(),
            cardType: data.cardType,
            cardId: data.cardId.trim(),
            authType: 0,
            update: false,
            industry: '其他',
            authOperationType: 1,
          }),
          method: 'POST',
        })

        if (res.ok) {
          const data = await res.json()
          if (data.code === 20000 && data.status && data.data) {
            return {
              qrcUrl: data.data.authUrl,
            }
          }
        }
        return 'failed'
      },
    },

    payment: {
      async createWeChatPay(options) {
        const res = await fetch('siliconflow/payment/create', {
          body: JSON.stringify({
            platform: 'wx',
            amount: String(options.amount),
          }),
          method: 'POST',
        })

        if (res.ok) {
          const data = await res.json()
          if (data.code === 20000 && data.status && data.data) {
            return {
              orderId: data.data.order,
              qrcUrl: data.data.codeUrl,
              qrcTimeout: 120 * 1000,
            }
          }
        }
        throw new Error('QR code generation failed')
      },

      async checkWeChatPay(options) {
        const res = await fetch(`siliconflow/payment/status?order=${options.orderId}`, {
          method: 'GET',
        })

        if (res.ok) {
          const data = await res.json()
          if (data.code === 20000 && data.status && data.data) {
            const payStatus = data.data.payStatus
            if (payStatus === 1) {
              return 'success'
            }
            else if (payStatus === 2) {
              return 'wait'
            }
          }
        }
        throw new Error('Payment status check failed')
      },
    },

    baseURL: 'https://api.siliconflow.cn/v1',
    async createKey() {
      const res = await fetch('siliconflow/apikey/create', {
        body: JSON.stringify({
          description: 'Generated by UniToken',
        }),
        method: 'POST',
      })

      if (res.ok) {
        const data = await res.json()
        return data.data.secretKey
      }
      else {
        throw new Error('API Key creation failed')
      }
    },
  }
})
