import type { ProviderUserInfo } from '@/lib/providers'
import { defineAsyncComponent, ref } from 'vue'
import { useI18n } from '@/lib/locals'
import { defineProvider, useProviderSession } from '@/lib/providers'
import { useServiceStore } from '@/stores'

export const useDeepSeekProvider = defineProvider(() => {
  const { t } = useI18n({
    'zh-CN': {
      providerName: 'DeepSeek',
      mainlandIdCard: '中国大陆二代居民身份证',
    },
    'en-US': {
      providerName: 'DeepSeek',
      mainlandIdCard: 'Mainland China ID Card',
    },
  })
  const { proxy, getSimpleProxyUrl } = useServiceStore()

  const session = useProviderSession<{
    token: string
  }>('deepseek')
  const user = ref<null | ProviderUserInfo>()

  async function makeHeaders(requireSession: boolean, extraHeaders: Record<string, string> = {}) {
    const token = requireSession ? (await session.get())?.token : null
    if (requireSession && !token) {
      throw new Error('No session token available')
    }

    return {
      'accept': '*/*',
      'accept-language': 'zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6',
      'priority': 'u=1, i',
      'referer': 'https://platform.deepseek.com/top_up',
      'sec-ch-ua': `"Not)A;Brand";v="8", "Chromium";v="138", "Microsoft Edge";v="138"`,
      'sec-ch-ua-mobile': '?0',
      'sec-ch-ua-platform': `"Linux"`,
      'sec-fetch-dest': 'empty',
      'sec-fetch-mode': 'cors',
      'sec-fetch-site': 'same-origin',
      'x-app-version': '20240425.0',
      'user-agent': 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36 Edg/138.0.0.0',
      'content-type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...extraHeaders,
    }
  }

  return {
    id: 'deepseek',
    get name() {
      return t('providerName')
    },
    get logo() {
      return import.meta.resolve('./logo.png')
    },
    homepage: 'https://www.deepseek.com/',

    get user() {
      return user.value
    },
    async refreshUser() {
      if (!(await session.get())?.token) {
        user.value = null
        return
      }

      const { ok, json } = await proxy('https://platform.deepseek.com/auth-api/v0/users/current', {
        method: 'GET',
        headers: await makeHeaders(true),
      })

      if (ok && json.code === 0 && json.data) {
        const { ok, json: balanceJson } = await proxy('https://platform.deepseek.com/api/v0/users/get_user_summary', {
          method: 'GET',
          headers: await makeHeaders(true),
        })
        if (!ok || balanceJson.code !== 0 || !balanceJson.data) {
          throw new Error('Failed to fetch user balance')
        }

        user.value = {
          name: json.data.id_profile?.name || json.data.mobile_number || 'Unknown',
          verified: !!json.data.identity_verification_id,
          phone: json.data.mobile_number,
          email: json.data.email,
          balance: +balanceJson.data.normal_wallets[0].balance,
        }
      }
      else {
        user.value = null
      }
    },

    Login: defineAsyncComponent(() => import('@/lib/providers/deepseek/Login.vue')),
    async logout() {
      await session.delete()
      user.value = null
    },

    verification: {
      async check() {
        const { ok, json } = await proxy('https://platform.deepseek.com/api/v1/my_identity_verification', {
          method: 'GET',
          headers: await makeHeaders(true),
        })

        if (ok && json.code === 0 && json.data) {
          return {
            name: json.data.data.name_desensitized || '',
            cardId: json.data.data.id_card_number_desensitized || '',
            time: undefined, // DeepSeek doesn't provide verification time
          }
        }
        else {
          return null
        }
      },

      get cardTypes() {
        return [
          { value: 1, label: t('mainlandIdCard') },
        ]
      },

      async submit(data) {
        const { ok, json } = await proxy('https://platform.deepseek.com/api/v1/identity_verify', {
          method: 'POST',
          headers: await makeHeaders(true),
          body: JSON.stringify({
            name: data.name.trim(),
            id: data.cardId.trim(),
          }),
        })

        if (ok && json.code === 0) {
          return 'success'
        }
        else {
          return 'failed'
        }
      },
    },

    payment: {
      async createWeChatPay(options) {
        // Generate a random UUID for request_id
        const requestId = crypto.randomUUID()

        const { ok, json } = await proxy('https://platform.deepseek.com/api/v1/payments', {
          method: 'POST',
          headers: await makeHeaders(true),
          body: JSON.stringify({
            order_info: {
              payment_method_type: 'WECHAT',
              amount: options.amount,
              currency: 'CNY',
              request_id: requestId,
            },
          }),
        })

        if (ok && json.code === 0 && json.data) {
          return {
            orderId: json.data.payment_order_id,
            qrcUrl: json.data.url,
            qrcTimeout: 15 * 60 * 1000, // 15 minutes timeout
          }
        }

        throw new Error('QR code generation failed')
      },

      async checkWeChatPay(options) {
        const { ok, json } = await proxy(`https://platform.deepseek.com/api/v1/payments/${options.orderId}/capture`, {
          method: 'GET',
          headers: await makeHeaders(true),
        })

        if (ok && json.code === 0 && json.data?.order) {
          const status = json.data.order.status
          if (status === 'SUCCESS') {
            return 'success'
          }
          else if (status === 'CREATED' || status === 'PENDING') {
            return 'wait'
          }
          else {
            return 'canceled'
          }
        }

        throw new Error('Payment status check failed')
      },
    },

    baseURL: 'https://api.deepseek.com/v1',
    async createKey() {
      const { ok, json } = await proxy('https://platform.deepseek.com/api/v0/users/edit_api_keys', {
        method: 'POST',
        headers: await makeHeaders(true),
        body: JSON.stringify({
          action: 'create',
          name: 'Generated by UniToken',
          redacted_key: null,
          created_at: null,
        }),
      })

      if (ok && json.code === 0 && json.data?.api_key?.sensitive_id) {
        return json.data.api_key.sensitive_id
      }
      throw new Error('API Key creation failed')
    },

    apis: {
      async sendSMS(payload: unknown) {
        const { ok, json } = await proxy('https://platform.deepseek.com/auth-api/v0/users/create_sms_verification_code', {
          method: 'POST',
          headers: await makeHeaders(false),
          body: JSON.stringify(payload),
        })
        if (!ok) {
          throw new Error('SMS request failed')
        }
        return json
      },

      async loginWithSMS(payload: unknown) {
        const answer = await getChallengeAndSolve('/auth-api/v0/users/login_by_mobile_sms')
        const { ok, json } = await proxy('https://platform.deepseek.com/auth-api/v0/users/login_by_mobile_sms', {
          method: 'POST',
          headers: await makeHeaders(false, {
            'X-DS-Guest-PoW-Response': btoa(JSON.stringify({
              salt: answer.challengeResponse.salt,
              answer: answer.challengeResponse.answer,
            })),
          }),
          body: JSON.stringify(payload),
        })
        if (!ok) {
          throw new Error('Login request failed')
        }
        await session.put({ token: json.data.user.token })
      },
    },
  }

  async function getChallengeAndSolve(target_path: string) {
    const { ok, json } = await proxy('https://platform.deepseek.com/auth-api/v0/users/create_guest_challenge', {
      method: 'POST',
      headers: await makeHeaders(false),
      body: JSON.stringify({
        target_path,
      }),
    })
    if (!ok) {
      throw new Error('Login request failed')
    }

    if (!window.Worker)
      throw new Error('Worker is not supported')

    const { ok: ok2, text: workerCode } = await proxy('https://static.deepseek.com/platform/static/282.cf61fc9b07.js', { method: 'GET' })
    if (!ok2 || !workerCode) {
      throw new Error('Failed to load PoW worker script')
    }
    const { guest_challenge } = json.data.biz_data

    const newBase = getSimpleProxyUrl('https://static.deepseek.com/platform')
    const newWorkerCode = workerCode.replace('https://static.deepseek.com/platform', newBase)

    const blob = new Blob([newWorkerCode], { type: 'application/javascript' })
    const workerUrl = URL.createObjectURL(blob)
    const worker = new Worker(workerUrl)
    URL.revokeObjectURL(workerUrl)

    worker.postMessage({
      type: 'pow-challenge',
      challenge: {
        ...guest_challenge,
        expireAt: guest_challenge.expire_at,
        expireAfter: guest_challenge.expire_after,
      },
    })
    const start = performance.now()
    const result = await new Promise<any>((resolve, reject) => {
      worker.onmessage = (r) => {
        if (r.data.type === 'pow-error') {
          reject(r.data.error)
          worker.terminate()
          return
        }
        if (r.data.type === 'pow-answer') {
          resolve(r.data)
          worker.terminate()
        }
      }
      worker.onerror = (e) => {
        reject(new Error('Worker error: '.concat(e.message, ' ').concat(e.filename, ':').concat(`${e.lineno}`, ':').concat(`${e.colno}`)))
        worker.terminate()
      }
    },
    )
    const a = performance.now()
    return {
      challengeResponse: result.answer,
      duration: a - start,
    }
  }
})
