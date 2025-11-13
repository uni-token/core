<script setup lang="ts">
import type { Provider, ProviderPaymentQRC } from '@/lib/providers'
import { CircleQuestionMarkIcon } from 'lucide-vue-next'
import { renderSVG } from 'uqr'
import { computed, onUnmounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue-sonner'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip'

const props = defineProps<{
  provider: Provider
}>()

const open = defineModel<boolean>()

const { t } = useI18n()

const payment = computed(() => props.provider.payment as ProviderPaymentQRC)
const quickAmounts = [10, 50, 100, 200, 500, 1000]
const formAmount = ref('')
const creating = ref(false)
const canCreatePayment = computed(() =>
  !creating.value
  && formAmount.value
  && Number.parseFloat(formAmount.value) > 0,
)
const qrcSvg = ref('')
const orderId = ref('')
const expired = ref(false)
const checking = ref(false)
let intervalId: any = null
let timeoutId: any = null

async function createPayment() {
  if (!canCreatePayment.value)
    return

  creating.value = true

  try {
    const { orderId: newOrderId, qrcUrl, interval, timeout } = await payment.value.create({
      amount: formAmount.value,
    })
    qrcSvg.value = renderSVG(qrcUrl, {})
    orderId.value = newOrderId
    expired.value = false

    if (interval) {
      intervalId = setInterval(() => {
        checkPayment()
      }, interval)
    }

    if (timeout) {
      timeoutId = setTimeout(() => {
        expired.value = true
        clearTimers()
      }, timeout)
    }
  }
  catch (error) {
    toast.error(t('paymentCreationFailed'))
    console.error(error)
  }
  finally {
    creating.value = false
  }
}

async function checkPayment(manual?: boolean) {
  if (manual) {
    checking.value = true
  }
  try {
    const result = await payment.value.check({ orderId: orderId.value })

    if (result === 'success') {
      toast.success(t('paymentCompleted'))
      open.value = false
      clearTimers()
    }
    else if (result === 'canceled') {
      toast.error(t('paymentCanceled'))
      open.value = false
      clearTimers()
    }
    else if (result === 'wait') {
      if (manual) {
        toast.info(t('paymentPending'))
      }
    }
  }
  catch (error) {
    console.error('Payment status check error:', error)
  }
  finally {
    if (manual) {
      checking.value = false
    }
  }
}

function clearTimers() {
  if (intervalId) {
    clearInterval(intervalId)
    intervalId = null
  }
  if (timeoutId) {
    clearTimeout(timeoutId)
    timeoutId = null
  }
}

onUnmounted(() => {
  clearTimers()
})

function clear() {
  clearTimers()
  formAmount.value = ''
  qrcSvg.value = ''
  orderId.value = ''
  expired.value = false
  checking.value = false
}
</script>

<template>
  <Dialog v-model:open="open" @update:open="open => !open && clear()">
    <DialogContent class="max-w-sm">
      <DialogHeader>
        <DialogTitle>{{ t('accountRecharge') }}</DialogTitle>
        <DialogDescription class="flex">
          {{ t('accountRechargeDescription', [provider.name]) }}
          <TooltipProvider>
            <Tooltip ignore-non-keyboard-focus>
              <TooltipTrigger>
                <CircleQuestionMarkIcon class="w-4 h-4  ml-1" />
              </TooltipTrigger>
              <TooltipContent>
                {{ t('accountRechargeTooltip', [provider.name]) }}
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
          <div class="flex-grow" />
          <div v-if="qrcSvg" class="text-sm text-muted-foreground">
            {{ payment.currency === 'CNY' ? formAmount + t('yuan') : `$${formAmount}` }}
          </div>
        </DialogDescription>
        <DialogClose />
      </DialogHeader>
      <div class="space-y-4">
        <div v-if="!qrcSvg" class="space-y-2">
          <label for="modal-amount" class="text-sm font-medium">{{ t('rechargeAmount') }}</label>
          <div class="flex items-center space-x-2">
            <Input
              id="modal-amount"
              v-model="formAmount"
              :placeholder="t('rechargeAmountPlaceholder')"
              type="number"
              step="0.01"
              min="0.01"
              class="flex-1"
            />
            <span class="text-sm text-muted-foreground">{{ payment.currency === 'CNY' ? t('yuan') : t('usd') }}</span>
          </div>
        </div>

        <!-- Quick Amount Buttons -->
        <div v-if="!qrcSvg" class="grid grid-cols-3 gap-2">
          <Button
            v-for="amount in quickAmounts"
            :key="amount"
            variant="outline"
            size="sm"
            class="h-8"
            @click="formAmount = amount.toString()"
          >
            {{ amount }}{{ payment.currency === 'CNY' ? t('yuan') : t('usd') }}
          </Button>
        </div>

        <!-- Payment QR Code -->
        <div v-if="qrcSvg" class="space-y-3">
          <div class="text-center">
            <div class="inline-block p-4 bg-white rounded-lg border-2 border-dashed border-gray-300">
              <div class="mx-auto w-40" v-html="qrcSvg" />
            </div>
          </div>
          <div class="text-center space-y-2">
            <p class="text-sm text-muted-foreground">
              {{ t('scanQRCodeMessage', [payment.platform]) }}
            </p>
            <div v-if="expired" class="text-center text-sm text-red-500">
              {{ t('qrCodeExpired') }}
            </div>
          </div>
          <div class="flex gap-2 justify-center">
            <Button v-if="expired" variant="default" @click="createPayment">
              {{ t('regenerateQRCode') }}
            </Button>

            <Button v-else variant="default" :disabled="checking" @click="checkPayment(true)">
              <span v-if="checking">{{ t('checking') }}</span>
              <span v-else>{{ t('check') }}</span>
            </Button>
          </div>
        </div>

        <!-- Generate Payment Button -->
        <div v-else class="space-y-3">
          <Button
            :disabled="!canCreatePayment"
            class="w-full"
            @click="createPayment"
          >
            <span v-if="creating">{{ t('generating') }}</span>
            <span v-else>{{ t('generatePaymentQRCode') }}</span>
          </Button>
        </div>
      </div>
    </DialogContent>
  </Dialog>
</template>

<i18n lang="yaml">
en-US:
  accountRecharge: Account Recharge
  accountRechargeDescription: Recharge to your {0} account
  rechargeAmount: Recharge Amount
  rechargeAmountPlaceholder: Please enter recharge amount
  yuan: Yuan
  usd: USD
  scanQRCodeMessage: Please use your phone to open {0} and scan the QR code to complete payment
  qrCodeExpired: QR code has expired, please regenerate
  regenerateQRCode: Regenerate QR Code
  checking: Checking...
  check: Complete
  paymentCreationFailed: Payment creation failed
  paymentCompleted: Payment Completed
  paymentCanceled: Payment Canceled
  paymentPending: Payment Not Completed
  generating: Generating...
  generatePaymentQRCode: Generate Payment QR Code
  accountRechargeTooltip: The funds will be directly recharged to your {0} account without going through the UniToken platform, ensuring the safety of your funds.

zh-CN:
  accountRecharge: 账户充值
  accountRechargeDescription: 向{0}账户充值
  rechargeAmount: 充值金额
  rechargeAmountPlaceholder: 请输入充值金额
  yuan: 元
  usd: 美元
  scanQRCodeMessage: 请使用手机打开{0}扫描二维码完成支付
  qrCodeExpired: 二维码已过期，请重新生成
  regenerateQRCode: 重新生成二维码
  checking: 检查中...
  check: 我已支付
  paymentCreationFailed: 支付创建失败
  paymentCompleted: 已完成
  paymentCanceled: 已取消支付
  paymentPending: 支付未完成
  generating: 生成中...
  generatePaymentQRCode: 生成支付二维码
  accountRechargeTooltip: 钱款会直接充值到您的{0}账户，不会经过 UniToken 平台，保证您的资金安全。
</i18n>
