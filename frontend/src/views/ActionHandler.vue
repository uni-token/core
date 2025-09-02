<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue-sonner'
import ProviderSelector from '@/components/ProviderSelector.vue'
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { useProvidersStore } from '@/stores'
import { useAppStore } from '@/stores/app'

const { t } = useI18n()
const params = new URLSearchParams(window.location.search)
let actionType = params.get('action')

const open = ref(!!actionType)
const selectedProvider = ref<string>('')
const appStore = useAppStore()
const providersStore = useProvidersStore()

async function registerAction(granted: boolean) {
  const appId = params.get('appId')

  if (!appId) {
    toast.error(t('actionHandler.missingAppId'))
    return
  }

  await appStore.toggleAppAuthorization(appId, granted, selectedProvider.value)
  actionType = null
  open.value = false
}

watch(open, (open) => {
  if (!open) {
    if (actionType === 'register') {
      registerAction(false)
    }
  }
})
</script>

<template>
  <Dialog v-model:open="open">
    <DialogContent v-if="actionType === 'register'" class="sm:max-w-xl">
      <DialogHeader>
        <DialogTitle>{{ t('actionHandler.appPermissionRequest') }}</DialogTitle>
        <DialogDescription>
          {{ t('actionHandler.appPermissionDescription') }}
        </DialogDescription>
      </DialogHeader>

      <div class="space-y-4">
        <div class="space-y-2">
          <div>
            <strong class="mr-4">{{ t('actionHandler.appName') }}</strong>
            <span class="text-lg">{{ params.get('appName') }}</span>
          </div>
          <div>
            <strong class="mr-4">{{ t('actionHandler.appDescription') }}</strong>
            <span class="text-lg">{{ params.get('appDescription') }}</span>
          </div>
          <ProviderSelector v-model="selectedProvider" class="pt-1" />
        </div>

        <div v-if="providersStore.providers.length > 0" class="flex justify-end mt-6 gap-2">
          <Button variant="outline" @click="registerAction(false)">
            {{ t('actionHandler.deny') }}
          </Button>
          <div
            :class="!selectedProvider ? 'cursor-not-allowed!' : ''"
            :title="!selectedProvider ? t('actionHandler.pleaseSelectProvider') : ''"
          >
            <Button
              variant="default"
              :disabled="!selectedProvider"
              @click="registerAction(true)"
            >
              {{ t('actionHandler.approve') }}
            </Button>
          </div>
        </div>
      </div>
    </DialogContent>
    <DialogContent v-else-if="actionType">
      <DialogHeader>
        <DialogTitle>{{ t('actionHandler.unrecognizedAction') }}</DialogTitle>
        <DialogDescription>
          {{ t('actionHandler.unrecognizedActionDescription', { actionType }) }}
        </DialogDescription>
      </DialogHeader>
    </DialogContent>
  </Dialog>
</template>
