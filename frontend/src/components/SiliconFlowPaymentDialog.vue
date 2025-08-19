<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogClose,
  DialogContent,
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
        <DialogTitle>{{ t('siliconFlow.accountRecharge') }}</DialogTitle>
        <DialogClose />
      </DialogHeader>
      <div class="space-y-4">
        <div class="space-y-2">
          <label for="modal-amount" class="text-sm font-medium">{{ t('siliconFlow.rechargeAmount') }}</label>
          <div class="flex items-center space-x-2">
            <Input
              id="modal-amount"
              v-model="siliconFlowStore.payment.amount"
              :placeholder="t('siliconFlow.rechargeAmountPlaceholder')"
              type="number"
              step="0.01"
              min="0.01"
              class="flex-1"
              :disabled="!!siliconFlowStore.payment.qrCode"
            />
            <span class="text-sm text-muted-foreground">{{ t('siliconFlow.currency') }}</span>
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
            {{ amount }}{{ t('siliconFlow.currency') }}
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
              {{ t('siliconFlow.scanQRCodeMessage') }}
            </p>
            <div v-if="siliconFlowStore.payment.expired" class="text-center text-sm text-red-500">
              {{ t('siliconFlow.qrCodeExpired') }}
            </div>
          </div>
          <div class="flex gap-2 justify-center">
            <Button v-if="siliconFlowStore.payment.expired" variant="default" @click="siliconFlowStore.regeneratePayment">
              {{ t('siliconFlow.regenerateQRCode') }}
            </Button>

            <Button v-else variant="default" :disabled="siliconFlowStore.payment.checking" @click="siliconFlowStore.manualCheckPayment">
              <span v-if="siliconFlowStore.payment.checking">{{ t('siliconFlow.checking') }}</span>
              <span v-else>{{ t('siliconFlow.paymentCompleted') }}</span>
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
            <span v-if="siliconFlowStore.payment.loading">{{ t('siliconFlow.generating') }}</span>
            <span v-else>{{ t('siliconFlow.generatePaymentQRCode') }}</span>
          </Button>
        </div>
      </div>
    </DialogContent>
  </Dialog>
</template>
