<script setup lang="ts">
import { renderSVG } from 'uqr'
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue-sonner'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { useSiliconFlowStore } from '@/stores'

const { t } = useI18n()
const siliconFlowStore = useSiliconFlowStore()

const CARD_TYPE_OPTIONS = computed(() => [
  { value: 1, label: t('siliconFlow.realName.cardTypes.mainlandIdCard') },
  { value: 2, label: t('siliconFlow.realName.cardTypes.hkMacaoPass') },
  { value: 3, label: t('siliconFlow.realName.cardTypes.taiwanPass') },
  { value: 4, label: t('siliconFlow.realName.cardTypes.hkMacaoResidence') },
  { value: 5, label: t('siliconFlow.realName.cardTypes.taiwanResidence') },
  { value: 6, label: t('siliconFlow.realName.cardTypes.foreignerPermit') },
  { value: 100, label: t('siliconFlow.realName.cardTypes.otherType') },
])

const open = defineModel({
  type: Boolean,
  default: false,
})
const realName = ref('')
const idCardNumber = ref('')
const cardType = ref(1)
const loading = ref(false)
const verificationStep = ref<'form' | 'qrcode' | 'completed'>('form')
const authUrl = ref('')
const authUrlQR = ref('')
const checkingAuth = ref(false)

const canSubmit = computed(() => {
  return realName.value.trim() && idCardNumber.value.trim() && !loading.value
})

const isAuthenticated = computed(() => {
  return siliconFlowStore.authInfo?.auth === 1
})

const authDisplayName = computed(() => {
  if (isAuthenticated.value && siliconFlowStore.authInfo?.username) {
    const name = siliconFlowStore.authInfo.username
    if (name.length > 1) {
      return name.charAt(0) + '*'.repeat(name.length - 1)
    }
    return name
  }
  return ''
})

const authDisplayIdCard = computed(() => {
  if (isAuthenticated.value && siliconFlowStore.authInfo?.cardId) {
    const idCard = siliconFlowStore.authInfo.cardId
    if (idCard.length > 8) {
      return idCard.substring(0, 6) + '*'.repeat(idCard.length - 10) + idCard.substring(idCard.length - 4)
    }
    return idCard
  }
  return ''
})

const selectedCardTypeLabel = computed(() => {
  const option = CARD_TYPE_OPTIONS.value.find(opt => opt.value === cardType.value)
  return option?.label || ''
})

function validateIdCard(idCard: string, cardType: number): boolean {
  if (!idCard)
    return false

  switch (cardType) {
    case 1: {
      const mainlandIdRegex = /^[1-9]\d{5}(?:18|19|20)\d{2}(?:0[1-9]|1[0-2])(?:[0-2][1-9]|10|20|30|31)\d{3}[0-9X]$/i
      return mainlandIdRegex.test(idCard)
    }
    case 2:
    case 3: {
      const hkMacaoTaiwanRegex = /^[HMC]\d{8,10}$|^\d{8,11}$/
      return hkMacaoTaiwanRegex.test(idCard)
    }
    case 4:
    case 5: {
      const residencePermitRegex = /^\d{18}$/
      return residencePermitRegex.test(idCard)
    }
    case 6: {
      const foreignerPermitRegex = /^[A-Z0-9]{8,20}$/i
      return foreignerPermitRegex.test(idCard)
    }
    case 100: {
      return idCard.length >= 5 && idCard.length <= 30
    }
    default:
      return false
  }
}

async function submitRealNameAuth() {
  if (!canSubmit.value)
    return

  if (!validateIdCard(idCardNumber.value, cardType.value)) {
    toast.error(t('siliconFlow.realName.invalidIdCard'))
    return
  }

  try {
    loading.value = true
    const result = await siliconFlowStore.submitRealNameAuth({
      username: realName.value.trim(),
      cardType: cardType.value,
      cardId: idCardNumber.value.trim(),
      authType: 0,
      update: false,
      industry: '其他',
      authOperationType: 1,
    })

    if (result.success && result.data?.authUrl) {
      authUrl.value = result.data.authUrl
      authUrlQR.value = renderSVG(result.data.authUrl, {})
      verificationStep.value = 'qrcode'
      toast.success(t('siliconFlow.realName.authRequestSubmitted'))
    }
    else {
      toast.error(result.message || t('siliconFlow.realName.authRequestFailed'))
    }
  }
  catch (error) {
    console.error('Real name auth error:', error)
    toast.error(t('siliconFlow.realName.authRequestFailed'))
  }
  finally {
    loading.value = false
  }
}

async function checkAuthStatus() {
  try {
    checkingAuth.value = true
    await siliconFlowStore.checkAuthInfo()

    if (isAuthenticated.value) {
      verificationStep.value = 'completed'
      toast.success(t('siliconFlow.realName.authCompleted'))
      // Delay closing dialog to allow user to see success message
      setTimeout(() => {
        open.value = false
        resetForm()
      }, 2000)
    }
    else {
      toast.info(t('siliconFlow.realName.authNotCompleted'))
    }
  }
  catch (error) {
    console.error('Check auth status error:', error)
    toast.error(t('siliconFlow.realName.checkAuthFailed'))
  }
  finally {
    checkingAuth.value = false
  }
}

function resetForm() {
  verificationStep.value = 'form'
  realName.value = ''
  idCardNumber.value = ''
  cardType.value = 1
  authUrl.value = ''
  authUrlQR.value = ''
  loading.value = false
  checkingAuth.value = false
}

function closeDialog() {
  open.value = false
  resetForm()
}

onMounted(() => {
  siliconFlowStore.checkAuthInfo()
})
</script>

<template>
  <Dialog v-model:open="open">
    <DialogContent class="sm:max-w-[425px]">
      <DialogHeader>
        <DialogTitle>{{ t('siliconFlow.realName.realNameAuth') }}</DialogTitle>
        <DialogDescription>
          {{ isAuthenticated ? t('siliconFlow.realName.authStatusDescription') : t('siliconFlow.realName.description') }}
        </DialogDescription>
      </DialogHeader>

      <div v-if="isAuthenticated" class="space-y-4">
        <Card>
          <CardHeader class="pb-3">
            <CardTitle class="text-sm flex items-center gap-2 text-green-700">
              <div class="h-2 w-2 bg-green-500 rounded-full" />
              {{ t('siliconFlow.realName.verified') }}
            </CardTitle>
          </CardHeader>
          <CardContent class="space-y-2">
            <div class="flex justify-between items-center text-sm">
              <span class="text-muted-foreground">{{ t('siliconFlow.realName.realName') }}</span>
              <span class="font-medium">{{ authDisplayName }}</span>
            </div>
            <div class="flex justify-between items-center text-sm">
              <span class="text-muted-foreground">{{ t('siliconFlow.realName.idCard') }}</span>
              <span class="font-medium">{{ authDisplayIdCard }}</span>
            </div>
            <div v-if="siliconFlowStore.authInfo?.authTime" class="flex justify-between items-center text-sm">
              <span class="text-muted-foreground">{{ t('siliconFlow.realName.authTime') }}</span>
              <span class="font-medium">{{ new Date(siliconFlowStore.authInfo.authTime.seconds * 1000).toLocaleDateString() }}</span>
            </div>
          </CardContent>
        </Card>
      </div>

      <div v-else-if="verificationStep === 'form'" class="space-y-4">
        <div class="space-y-2">
          <label for="realName" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">{{ t('siliconFlow.realName.realName') }}</label>
          <Input
            id="realName"
            v-model="realName"
            :placeholder="t('siliconFlow.realName.realNamePlaceholder')"
            :disabled="loading"
          />
        </div>

        <div class="space-y-2">
          <label for="cardType" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">{{ t('siliconFlow.realName.cardType') }}</label>
          <Select
            v-model="cardType"
            :disabled="loading"
          >
            <SelectTrigger>
              <SelectValue>
                {{ CARD_TYPE_OPTIONS.find(option => option.value === cardType)?.label }}
              </SelectValue>
            </SelectTrigger>
            <SelectContent>
              <SelectItem v-for="option in CARD_TYPE_OPTIONS" :key="option.value" :value="option.value">
                {{ option.label }}
              </SelectItem>
            </SelectContent>
          </Select>
        </div>

        <div class="space-y-2">
          <label for="idCardNumber" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">{{ selectedCardTypeLabel }}</label>
          <Input
            id="idCardNumber"
            v-model="idCardNumber"
            :placeholder="t('siliconFlow.realName.idCardPlaceholder')"
            :disabled="loading"
            :maxlength="cardType === 100 ? 30 : (cardType === 6 ? 20 : 18)"
          />
        </div>
      </div>

      <div v-else-if="verificationStep === 'qrcode'" class="space-y-4">
        <div class="text-center space-y-4">
          <div class="inline-block p-4 bg-white rounded-lg border-2 border-dashed border-gray-300">
            <div class="mx-auto w-48" v-html="authUrlQR" />
          </div>
          <p class="text-sm text-muted-foreground">
            {{ t('siliconFlow.realName.scanAlipayQR') }}
          </p>
        </div>
      </div>

      <div v-else-if="verificationStep === 'completed'" class="space-y-4">
        <div class="text-center space-y-4">
          <div class="mx-auto w-16 h-16 bg-green-100 rounded-full flex items-center justify-center">
            <svg class="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <p class="text-lg font-medium">
            {{ t('siliconFlow.realName.authSuccess') }}
          </p>
          <p class="text-sm text-muted-foreground">
            {{ t('siliconFlow.realName.authSuccessDescription') }}
          </p>
        </div>
      </div>

      <DialogFooter>
        <div v-if="verificationStep === 'form' && !isAuthenticated" class="flex gap-2 w-full">
          <Button variant="outline" class="flex-1" @click="closeDialog">
            {{ t('common.cancel') }}
          </Button>
          <Button :disabled="!canSubmit" class="flex-1" @click="submitRealNameAuth">
            <span v-if="loading">{{ t('siliconFlow.realName.submitting') }}</span>
            <span v-else>{{ t('common.submit') }}</span>
          </Button>
        </div>

        <div v-else-if="verificationStep === 'qrcode'" class="flex gap-2 w-full">
          <Button variant="outline" class="flex-1" @click="closeDialog">
            {{ t('common.cancel') }}
          </Button>
          <Button :disabled="checkingAuth" class="flex-1" @click="checkAuthStatus">
            <span v-if="checkingAuth">{{ t('siliconFlow.realName.checking') }}</span>
            <span v-else>{{ t('siliconFlow.realName.completedFaceAuth') }}</span>
          </Button>
        </div>

        <div v-else class="flex w-full">
          <Button class="w-full" @click="closeDialog">
            {{ t('common.close') }}
          </Button>
        </div>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
