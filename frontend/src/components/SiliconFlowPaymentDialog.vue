<script setup lang="ts">
import { useI18n } from 'vue-i18n'
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
import { useSiliconFlowStore } from '@/stores'

const { t } = useI18n()
const siliconFlowStore = useSiliconFlowStore()

const model = defineModel<boolean>()

function handleOpenChange(open: boolean) {
  if (!open) {
    siliconFlowStore.closePaymentDialog()
  }
}
</script>

<template>
  <Dialog v-model:open="model" @update:open="handleOpenChange">
    <DialogContent class="max-w-sm">
      <DialogHeader>
        <DialogTitle>{{ t('accountRecharge') }}</DialogTitle>
        <DialogDescription>
          {{ t('accountRechargeDescription') }}
        </DialogDescription>
        <DialogClose />
      </DialogHeader>
      <div class="space-y-4">
        <div class="space-y-2">
          <label for="modal-amount" class="text-sm font-medium">{{ t('rechargeAmount') }}</label>
          <div class="flex items-center space-x-2">
            <Input
              id="modal-amount"
              v-model="siliconFlowStore.payment.amount"
              :placeholder="t('rechargeAmountPlaceholder')"
              type="number"
              step="0.01"
              min="0.01"
              class="flex-1"
              :disabled="!!siliconFlowStore.payment.qrCode"
            />
            <span class="text-sm text-muted-foreground">{{ t('currency') }}</span>
          </div>
        </div>

        <!-- Quick Amount Buttons -->
        <div v-if="!siliconFlowStore.payment.qrCode" class="grid grid-cols-3 gap-2">
          <Button
            v-for="amount in siliconFlowStore.quickAmounts"
            :key="amount"
            variant="outline"
            size="sm"
            class="h-8"
            @click="siliconFlowStore.setPaymentAmount(amount.toString())"
          >
            {{ amount }}{{ t('currency') }}
          </Button>
        </div>

        <!-- Payment QR Code -->
        <div v-if="siliconFlowStore.payment.qrCode" class="space-y-3">
          <div class="text-center">
            <div class="inline-block p-4 bg-white rounded-lg border-2 border-dashed border-gray-300">
              <div class="mx-auto w-40" v-html="siliconFlowStore.payment.qrCodeSVG" />
            </div>
          </div>
          <div class="text-center space-y-2">
            <p class="text-sm text-muted-foreground">
              {{ t('scanQRCodeMessage') }}
            </p>
            <div v-if="siliconFlowStore.payment.expired" class="text-center text-sm text-red-500">
              {{ t('qrCodeExpired') }}
            </div>
          </div>
          <div class="flex gap-2 justify-center">
            <Button v-if="siliconFlowStore.payment.expired" variant="default" @click="siliconFlowStore.regeneratePayment">
              {{ t('regenerateQRCode') }}
            </Button>

            <Button v-else variant="default" :disabled="siliconFlowStore.payment.checking" @click="siliconFlowStore.manualCheckPayment">
              <span v-if="siliconFlowStore.payment.checking">{{ t('checking') }}</span>
              <span v-else>{{ t('paymentCompleted') }}</span>
            </Button>
          </div>
        </div>

        <!-- Generate Payment Button -->
        <div v-else class="space-y-3">
          <Button
            :disabled="!siliconFlowStore.canCreatePayment"
            class="w-full"
            @click="siliconFlowStore.createPayment"
          >
            <span v-if="siliconFlowStore.payment.loading">{{ t('generating') }}</span>
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
  accountRechargeDescription: Recharge to your SiliconFlow account
  rechargeAmount: Recharge Amount
  rechargeAmountPlaceholder: Please enter recharge amount
  currency: Yuan
  scanQRCodeMessage: Please use your phone to open WeChat and scan the QR code to complete payment
  qrCodeExpired: QR code has expired, please regenerate
  regenerateQRCode: Regenerate QR Code
  checking: Checking...
  paymentCompleted: Payment Completed
  generating: Generating...
  generatePaymentQRCode: Generate Payment QR Code

zh-CN:
  accountRecharge: 账户充值
  accountRechargeDescription: 向硅基流动账户充值
  rechargeAmount: 充值金额
  rechargeAmountPlaceholder: 请输入充值金额
  currency: 元
  scanQRCodeMessage: 请使用手机打开微信扫描二维码完成支付
  qrCodeExpired: 二维码已过期，请重新生成
  regenerateQRCode: 重新生成二维码
  checking: 检查中...
  paymentCompleted: 已完成支付
  generating: 生成中...
  generatePaymentQRCode: 生成支付二维码
</i18n>
