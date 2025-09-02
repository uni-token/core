import { defineStore } from 'pinia'
import { renderSVG } from 'uqr'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue-sonner'
import { useKeysStore } from './keys'
import { useServiceStore } from './service'

export interface SiliconFlowUserInfo {
  isLoggedIn: boolean
  data?: {
    id?: string
    name?: string
    phone?: string
    email?: string
    balance?: string
    totalBalance?: string
    auth: number // 0: Not authed, 1: authed
    [key: string]: any
  }
  message?: string
}

export interface AuthInfo {
  auth: number // 0: Not authed, 1: authed
  authType: number
  username: string
  cardId: string
  cardType: number
  authTime: {
    seconds: number
    nanos: number
  } | null
  industry: string
  authLog: any
}

export interface RealNameAuthRequest {
  username: string
  cardType: number
  cardId: string
  authType: number
  update: boolean
  industry: string
  authOperationType: number
}

export interface RealNameAuthResponse {
  success: boolean
  message?: string
  data?: {
    authType: string
    certifyId: string
    authUrl: string
    code: string
    msg: string
    sub_code: string
    sub_msg: string
  }
}

export interface PaymentInfo {
  amount: string
  loading: boolean
  qrCode: string
  order: string
  qrCodeSVG: string
  countdown: number
  expired: boolean
  checking: boolean
}

export const useSiliconFlowStore = defineStore('siliconflow', () => {
  const { fetch } = useServiceStore()
  const keysStore = useKeysStore()
  const { t } = useI18n()

  // User related state
  const userInfo = ref<SiliconFlowUserInfo | null>(null)
  const authInfo = ref<AuthInfo | null>(null)
  const phoneNumber = ref('')
  const email = ref('')
  const smsCode = ref('')
  const agreed = ref(true)
  const keepLogin = ref(true)
  const isLoading = ref(true)
  const isEmailLogin = ref(false)

  // Payment related state
  const payment = ref<PaymentInfo>({
    amount: '',
    loading: false,
    qrCode: '',
    order: '',
    qrCodeSVG: '',
    countdown: 0,
    expired: false,
    checking: false,
  })

  const showPaymentDialog = ref(false)
  const quickAmounts = [10, 50, 100, 200, 500, 1000]

  // Intervals
  let paymentCheckInterval: any | null = null
  let qrCountdownInterval: any | null = null

  // Computed
  const isLoggedIn = computed(() => userInfo.value?.isLoggedIn || false)
  const userBalance = computed(() => userInfo.value?.data?.totalBalance || '0')
  const canLogin = computed(() => {
    const hasContact = isEmailLogin.value ? email.value : phoneNumber.value
    return hasContact && smsCode.value && agreed.value && !isLoading.value
  })
  const canCreatePayment = computed(() =>
    !payment.value.loading
    && payment.value.amount
    && Number.parseFloat(payment.value.amount) > 0,
  )

  // Captcha config
  const captchaConfig = {
    captchaId: '592ad182314270f0c1442d9aa82d3ac2',
    product: 'bind',
    language: 'zho',
    riskType: 'nine',
    protocol: 'https://',
  }

  // Actions
  async function checkLoginStatus() {
    try {
      const res = await fetch('siliconflow/status', {
        method: 'GET',
      })

      if (res.ok) {
        const data = await res.json()
        if (data.code === 20000 && data.status && data.data) {
          userInfo.value = {
            isLoggedIn: true,
            data: data.data,
            message: data.message,
          }
        }
        else {
          userInfo.value = {
            isLoggedIn: false,
            message: data.message || 'Not logged in',
          }
        }
      }
      else {
        userInfo.value = { isLoggedIn: false }
      }
    }
    catch (error) {
      console.error('Failed to check login status:', error)
      userInfo.value = { isLoggedIn: false }
    }
    finally {
      isLoading.value = false
    }
  }

  async function sendSMS(result: any) {
    if (isEmailLogin.value) {
      // Send email verification code
      await fetch('siliconflow/email', {
        body: JSON.stringify({
          email: email.value,
          ...result,
        }),
        method: 'POST',
      })
    }
    else {
      // Send SMS verification code
      await fetch('siliconflow/sms', {
        body: JSON.stringify({
          area: '+86',
          phone: phoneNumber.value,
          ...result,
        }),
        method: 'POST',
      })
    }
  }

  async function login() {
    if (isLoading.value)
      return

    try {
      isLoading.value = true
      const res = isEmailLogin.value
        ? await fetch('siliconflow/login/email', {
            body: JSON.stringify({
              email: email.value,
              code: smsCode.value,
              agree: agreed.value,
              keep: keepLogin.value,
              area: '+86',
            }),
            method: 'POST',
          })
        : await fetch('siliconflow/login', {
            body: JSON.stringify({
              phone: phoneNumber.value,
              code: smsCode.value,
              shareCode: '',
              agree: agreed.value,
              keep: keepLogin.value,
              area: '+86',
            }),
            method: 'POST',
          })

      if (res.ok) {
        await checkLoginStatus()
        // Clear login form
        phoneNumber.value = ''
        email.value = ''
        smsCode.value = ''
        toast.success(t('stores.siliconflow.loginSuccess'))

        // Automatically create API key after successful login
        await createApiKeyAndApply()
      }
      else {
        const errorData = await res.json()
        if (res.status === 401) {
          toast.error(errorData.message || t('stores.siliconflow.invalidCode'))
        }
        else {
          toast.error(errorData.message || t('stores.siliconflow.loginFailed'))
        }
      }
    }
    catch (error) {
      console.error('Login error:', error)
      toast.error(t('stores.siliconflow.networkError'))
    }
    finally {
      isLoading.value = false
    }
  }

  async function logout() {
    try {
      const res = await fetch('siliconflow/logout', {
        method: 'POST',
      })

      if (res.ok) {
        userInfo.value = { isLoggedIn: false }
        toast.success(t('stores.siliconflow.logoutSuccess'))
      }
      else {
        toast.error(t('stores.siliconflow.logoutFailed'))
      }
    }
    catch (error) {
      console.error('Failed to logout:', error)
      toast.error(t('stores.siliconflow.logoutFailedRetry'))
    }
  }

  async function createApiKeyAndApply() {
    try {
      const res = await fetch('siliconflow/apikey/create', {
        body: JSON.stringify({
          description: 'Generated by UniToken',
        }),
        method: 'POST',
      })

      if (res.ok) {
        const data = await res.json()
        const result = await saveApiKey(data.data.secretKey)
        toast.success(t('stores.siliconflow.apiKeyCreated'))
        return result
      }
      else {
        const errorData = await res.json()
        toast.error(errorData.message || t('stores.siliconflow.apiKeyCreateFailed'))
      }
    }
    catch (error) {
      console.error('Failed to create API key:', error)
      toast.error(t('stores.siliconflow.apiKeyCreateFailed'))
    }
  }

  async function saveApiKey(apiKey: string) {
    function getName(index: number): string {
      return `${t('stores.siliconflow.keyName')}${index > 0 ? ` ${index}` : ''}`
    }

    let count = 0
    while (true) {
      if (keysStore.keys.some(p => p.name === getName(count))) {
        count++
      }
      else {
        break
      }
    }

    return await keysStore.addKey({
      name: getName(count),
      type: 'siliconflow',
      protocol: 'openai',
      baseUrl: 'https://api.siliconflow.cn/v1',
      token: apiKey,
    })
  }

  // Payment actions
  function openPaymentDialog() {
    showPaymentDialog.value = true
    resetPayment()
  }

  function closePaymentDialog() {
    showPaymentDialog.value = false
    resetPayment()
    stopPaymentCheck()
    stopQRCountdown()
  }

  function resetPayment() {
    payment.value.qrCode = ''
    payment.value.order = ''
    payment.value.qrCodeSVG = ''
    payment.value.amount = ''
    payment.value.expired = false
    payment.value.countdown = 0
  }

  function setPaymentAmount(amount: string) {
    payment.value.amount = amount
  }

  async function createPayment() {
    if (!canCreatePayment.value)
      return

    try {
      payment.value.loading = true
      const res = await fetch('siliconflow/payment/create', {
        body: JSON.stringify({
          platform: 'wx',
          amount: String(payment.value.amount),
        }),
        method: 'POST',
      })

      if (res.ok) {
        const data = await res.json()
        if (data.code === 20000 && data.status && data.data) {
          payment.value.qrCode = data.data.codeUrl
          payment.value.order = data.data.order

          // Generate QR code SVG
          payment.value.qrCodeSVG = renderSVG(data.data.codeUrl, {})

          // Start payment checking and countdown
          startPaymentCheck()
          startQRCountdown()

          toast.success(t('stores.siliconflow.qrCodeGenerated'))
        }
        else {
          toast.error(data.message || t('stores.siliconflow.qrCodeGenerationFailed'))
        }
      }
      else {
        const errorData = await res.json()
        toast.error(errorData.message || t('stores.siliconflow.qrCodeGenerationFailed'))
      }
    }
    catch (error) {
      console.error('Payment creation error:', error)
      toast.error(t('stores.siliconflow.networkErrorRetry'))
    }
    finally {
      payment.value.loading = false
    }
  }

  function startPaymentCheck() {
    if (paymentCheckInterval) {
      clearInterval(paymentCheckInterval)
    }

    paymentCheckInterval = setInterval(async () => {
      if (!payment.value.order || payment.value.expired) {
        stopPaymentCheck()
        return
      }

      try {
        const res = await fetch(`siliconflow/payment/status?order=${payment.value.order}`, {
          method: 'GET',
        })

        if (res.ok) {
          const data = await res.json()
          if (data.code === 20000 && data.status && data.data) {
            const payStatus = data.data.payStatus
            if (payStatus === 1) {
              toast.success(t('stores.siliconflow.paymentSuccess'))
              await checkLoginStatus()
              closePaymentDialog()
            }
          }
        }
      }
      catch (error) {
        console.error('Payment status check error:', error)
      }
    }, 3000)
  }

  function stopPaymentCheck() {
    if (paymentCheckInterval) {
      clearInterval(paymentCheckInterval)
      paymentCheckInterval = null
    }
  }

  function startQRCountdown() {
    payment.value.countdown = 120
    payment.value.expired = false

    if (qrCountdownInterval) {
      clearInterval(qrCountdownInterval)
    }

    qrCountdownInterval = setInterval(() => {
      payment.value.countdown--
      if (payment.value.countdown <= 0) {
        payment.value.expired = true
        stopQRCountdown()
        stopPaymentCheck()
        toast.warning(t('stores.siliconflow.qrCodeExpired'))
      }
    }, 1000)
  }

  function stopQRCountdown() {
    if (qrCountdownInterval) {
      clearInterval(qrCountdownInterval)
      qrCountdownInterval = null
    }
  }

  async function manualCheckPayment() {
    if (!payment.value.order)
      return

    try {
      payment.value.checking = true
      const res = await fetch(`siliconflow/payment/status?order=${payment.value.order}`, {
        method: 'GET',
      })

      if (res.ok) {
        const data = await res.json()
        if (data.code === 20000 && data.status && data.data) {
          const payStatus = data.data.payStatus
          if (payStatus === 1) {
            toast.success(t('stores.siliconflow.paymentSuccess'))
            await checkLoginStatus()
            closePaymentDialog()
          }
          else {
            toast.info(t('stores.siliconflow.paymentNotCompleted'))
          }
        }
        else {
          toast.error(t('stores.siliconflow.checkPaymentFailed'))
        }
      }
    }
    catch (error) {
      console.error('Manual payment check error:', error)
      toast.error(t('stores.siliconflow.checkPaymentFailed'))
    }
    finally {
      payment.value.checking = false
    }
  }

  function regeneratePayment() {
    if (payment.value.amount) {
      createPayment()
    }
    else {
      resetPayment()
    }
  }

  // Real name authentication actions
  async function checkAuthInfo(): Promise<void> {
    checkLoginStatus()
    try {
      const res = await fetch('siliconflow/auth/info', {
        method: 'GET',
      })

      if (res.ok) {
        const data = await res.json()
        if (data.code === 20000 && data.status && data.data) {
          authInfo.value = data.data
        }
        else {
          authInfo.value = null
        }
      }
      else {
        authInfo.value = null
      }
    }
    catch (error) {
      console.error('Failed to check auth info:', error)
      toast.error(t('stores.siliconflow.getRealNameInfoFailed'))
      authInfo.value = null
    }
  }

  async function submitRealNameAuth(request: RealNameAuthRequest): Promise<RealNameAuthResponse> {
    try {
      const res = await fetch('siliconflow/auth/save', {
        body: JSON.stringify(request),
        method: 'POST',
      })

      if (res.ok) {
        const data = await res.json()
        if (data.code === 20000 && data.status && data.data) {
          return {
            success: true,
            data: data.data,
          }
        }
        else {
          return {
            success: false,
            message: data.message || t('stores.siliconflow.realNameAuthSubmitFailed'),
          }
        }
      }
      else {
        const errorData = await res.json()
        return {
          success: false,
          message: errorData.message || t('stores.siliconflow.realNameAuthSubmitFailed'),
        }
      }
    }
    catch (error) {
      console.error('Real name auth error:', error)
      toast.error(t('stores.siliconflow.networkErrorRetry'))
      return {
        success: false,
        message: 'Network error, please try again later',
      }
    }
  }

  return {
    // State
    userInfo,
    authInfo,
    phoneNumber,
    email,
    smsCode,
    agreed,
    keepLogin,
    isLoading,
    payment,
    showPaymentDialog,
    quickAmounts,
    captchaConfig,
    isEmailLogin,

    // Computed
    isLoggedIn,
    userBalance,
    canLogin,
    canCreatePayment,
    authed: computed(() => userInfo.value?.data?.auth === 1),
    keyId: computed(() => keysStore.keys.find(p => p.type === 'siliconflow')?.id || null),

    // Actions
    checkLoginStatus,
    sendSMS,
    login,
    logout,
    createApiKeyAndApply,
    saveApiKey,
    openPaymentDialog,
    closePaymentDialog,
    resetPayment,
    setPaymentAmount,
    createPayment,
    manualCheckPayment,
    regeneratePayment,
    checkAuthInfo,
    submitRealNameAuth,

    // Initialize method
    async init() {
      await checkLoginStatus()
    },
  }
})
