import type { ProviderUserInfo } from './index'
import { defineAsyncComponent, ref } from 'vue'
import { toast } from 'vue-sonner'
import { useServiceStore } from '@/stores'
import { useI18n } from '../locals'
import { defineProvider, useProviderSession } from './index'

export const useSiliconFlowProvider = defineProvider(() => {
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
      invalidCode: '验证码无效，请重试',
      loginFailed: '登录失败，请稍后重试',
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
      invalidCode: 'Invalid code, please try again',
      loginFailed: 'Login failed, please try again later',
    },
  })
  const { proxy } = useServiceStore()

  const session = useProviderSession<{
    cookie: string
    subjectID: string
  }>('siliconflow')

  const user = ref<null | ProviderUserInfo>()

  const commonHeaders = {
    'Accept': '*/*',
    'Accept-Language': 'zh-CN,zh;q=0.9',
    'Priority': 'u=1, i',
    'Sec-CH-UA': `"Not)A;Brand";v="8", "Chromium";v="138", "Microsoft Edge";v="138"`,
    'Sec-CH-UA-Mobile': '?0',
    'Sec-CH-UA-Platform': `"Linux"`,
    'Sec-Fetch-Dest': 'empty',
    'Sec-Fetch-Mode': 'cors',
    'Sec-Fetch-Site': 'same-origin',
    'User-Agent': 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36 Edg/138.0.0.0',
    'Origin': 'https://cloud.siliconflow.cn',
  }

  return {
    id: 'siliconflow',
    get name() {
      return t('providerName')
    },
    homepage: 'https://www.siliconflow.cn/',

    get user() {
      return user.value
    },
    async refreshUser() {
      const s = await session.get()
      if (!s) {
        user.value = null
        return
      }

      const { ok, json } = await proxy('https://cloud.siliconflow.cn/biz-server/api/v1/user/info', {
        headers: {
          ...commonHeaders,
          'Referer': 'https://cloud.siliconflow.cn/me/account/info',
          'X-Subject-ID': s.subjectID,
          'Cookie': s.cookie,
        },
      })

      if (!ok) {
        user.value = null
        return
      }

      if (json.code === 20000 && json.status && json.data) {
        user.value = {
          name: json.data.name,
          verified: json.data.auth === 1,
          phone: json.data.phone,
          email: json.data.email,
          balance: json.data.balance,
        }
      }
    },

    Login: defineAsyncComponent(() => import('@/components/SiliconFlowLoginCard.vue')),
    async logout() {
      await session.delete()
    },

    verification: {
      async check() {
        const s = await session.get()
        if (!s) {
          throw new Error('No session found')
        }
        const { ok, json } = await proxy('https://cloud.siliconflow.cn/biz-server/api/v1/subject/auth/info', {
          headers: {
            ...commonHeaders,
            'Referer': 'https://cloud.siliconflow.cn/me/account/authentication/personal',
            'X-Subject-ID': s.subjectID,
            'Cookie': s.cookie,
          },
        })

        if (ok && json.code === 20000 && json.status && json.data?.auth) {
          let name = json.data.username
          if (name.length > 1) {
            name = name.charAt(0) + '*'.repeat(name.length - 1)
          }

          let cardId = json.data.cardId
          if (cardId.length > 8) {
            cardId = cardId.substring(0, 6) + '*'.repeat(cardId.length - 10) + cardId.substring(cardId.length - 4)
          }

          return {
            name,
            cardId,
            time: json.data.authTime?.seconds ? json.data.authTime.seconds * 1000 : undefined,
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
        const s = await session.get()
        if (!s) {
          throw new Error('No session found')
        }

        const { ok, json } = await proxy('https://cloud.siliconflow.cn/biz-server/api/v1/subject/auth/save', {
          method: 'POST',
          headers: {
            ...commonHeaders,
            'Referer': 'https://cloud.siliconflow.cn/me/account/authentication/personal',
            'X-Subject-ID': s.subjectID,
            'Cookie': s.cookie,
          },
          body: JSON.stringify({
            username: data.name.trim(),
            cardType: data.cardType,
            cardId: data.cardId.trim(),
            authType: 0,
            update: false,
            industry: '其他',
            authOperationType: 1,
          }),
        })
        if (ok && json.code === 20000 && json.status && json.data) {
          return {
            qrcUrl: json.data.authUrl,
          }
        }
        return 'failed'
      },
    },

    payment: {
      async createWeChatPay(options) {
        const s = await session.get()
        if (!s) {
          throw new Error('No session found')
        }
        const { ok, json } = await proxy('https://cloud.siliconflow.cn/biz-server/api/v1/pay/transactions', {
          method: 'POST',
          body: JSON.stringify({
            platform: 'wx',
            amount: String(options.amount),
          }),
          headers: {
            ...commonHeaders,
            'Referer': 'https://cloud.siliconflow.cn/me/account/recharge',
            'X-Subject-ID': s.subjectID,
            'Cookie': s.cookie,
          },
        })

        if (ok && json.code === 20000 && json.status && json.data) {
          return {
            orderId: json.data.order,
            qrcUrl: json.data.codeUrl,
            qrcTimeout: 120 * 1000,
          }
        }
        throw new Error('QR code generation failed')
      },

      async checkWeChatPay(options) {
        const s = await session.get()
        if (!s) {
          throw new Error('No session found')
        }
        const { ok, json } = await proxy(`https://cloud.siliconflow.cn/biz-server/api/v1/pay/status?order=${options.orderId}`, {
          method: 'GET',
          headers: {
            ...commonHeaders,
            'X-Subject-ID': s.subjectID,
            'Cookie': s.cookie,
          },
        })

        if (ok && json.code === 20000 && json.status && json.data) {
          const payStatus = json.data.payStatus
          if (payStatus === 1) {
            return 'success'
          }
          else if (payStatus === 2) {
            return 'wait'
          }
        }
        throw new Error('Payment status check failed')
      },
    },

    baseURL: 'https://api.siliconflow.cn/v1',
    async createKey() {
      const s = await session.get()
      if (!s) {
        throw new Error('No session found')
      }
      const { ok, json } = await proxy('https://cloud.siliconflow.cn/biz-server/api/v1/apikey/create', {
        body: JSON.stringify({
          description: 'Generated by UniToken',
        }),
        method: 'POST',
        headers: {
          ...commonHeaders,
          'Referer': 'https://cloud.siliconflow.cn/me/account/ak',
          'X-Subject-ID': s.subjectID,
          'Cookie': s.cookie,
        },
      })

      if (ok) {
        return json.data.secretKey
      }
      else {
        throw new Error('API Key creation failed')
      }
    },

    async sendSms(payload: string) {
      return proxy('https://account.siliconflow.cn/api/open/sms', {
        method: 'POST',
        headers: commonHeaders,
        body: payload,
      })
    },

    async sendEmail(payload: string) {
      return proxy('https://account.siliconflow.cn/api/open/email', {
        method: 'POST',
        headers: commonHeaders,
        body: payload,
      })
    },

    async loginViaSms(payload: any) {
      const { ok, status, headers } = await proxy('https://account.siliconflow.cn/api/open/login/user', {
        method: 'POST',
        headers: commonHeaders,
        body: JSON.stringify(payload),
      })

      if (!ok) {
        toast.error(status === 401 ? t('invalidCode') : t('loginFailed'))
        throw new Error(`Login request failed with status ${status}`)
      }

      const setCookieHeader = headers.get('Set-Cookie')
      if (!setCookieHeader) {
        throw new Error('No Set-Cookie header found')
      }
      const cookie = setCookieHeader.split(';').map(c => c.trim()).join('; ')

      const meRes = await proxy('https://cloud.siliconflow.cn/me', {
        headers: {
          ...commonHeaders,
          'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7',
          'Referer': 'https://account.siliconflow.cn/',
          'Sec-Fetch-Dest': 'document',
          'Sec-Fetch-Mode': 'navigate',
          'Sec-Fetch-Site': 'same-site',
          'Sec-Fetch-User': '?1',
          'Upgrade-Insecure-Requests': '1',
          'Cookie': cookie,
        },
      })
      const subjectID = meRes.headers.get('X-Subject-ID')!
      await session.put({ cookie, subjectID })
    },

    async loginViaEmail(payload: any) {
      const { ok, status, headers } = await proxy('https://account.siliconflow.cn/api/open/login/email', {
        method: 'POST',
        headers: commonHeaders,
        body: JSON.stringify(payload),
      })

      if (!ok) {
        toast.error(status === 401 ? t('invalidCode') : t('loginFailed'))
        throw new Error(`Login request failed with status ${status}`)
      }

      const setCookieHeader = headers.get('Set-Cookie')
      if (!setCookieHeader) {
        throw new Error('No Set-Cookie header found')
      }
      const cookie = setCookieHeader.split(';').map(c => c.trim()).join('; ')

      const meRes = await proxy('https://cloud.siliconflow.cn/me', {
        headers: {
          ...commonHeaders,
          'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7',
          'Referer': 'https://account.siliconflow.cn/',
          'Sec-Fetch-Dest': 'document',
          'Sec-Fetch-Mode': 'navigate',
          'Sec-Fetch-Site': 'same-site',
          'Sec-Fetch-User': '?1',
          'Upgrade-Insecure-Requests': '1',
          'Cookie': cookie,
        },
      })
      const subjectID = meRes.headers.get('X-Subject-ID')!
      await session.put({ cookie, subjectID })
    },
  }
})
