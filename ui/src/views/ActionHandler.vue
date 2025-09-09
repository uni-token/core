<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue-sonner'
import KeySelector from '@/components/KeySelector.vue'
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { useKeysStore } from '@/stores'
import { useAppStore } from '@/stores/app'

const { t } = useI18n()
const params = new URLSearchParams(window.location.search)
let actionType = params.get('action')

const open = ref(!!actionType)
const selectedKey = ref<string>('')
const appStore = useAppStore()
const keysStore = useKeysStore()

async function registerAction(granted: boolean) {
  const appId = params.get('appId')

  if (!appId) {
    toast.error(t('missingAppId'))
    return
  }

  await appStore.toggleAppAuthorization(appId, granted, selectedKey.value)
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
        <DialogTitle>{{ t('appPermissionRequest') }}</DialogTitle>
        <DialogDescription>
          {{ t('appPermissionDescription') }}
        </DialogDescription>
      </DialogHeader>

      <div class="space-y-4">
        <div class="space-y-2">
          <div>
            <strong class="mr-4">{{ t('appName') }}</strong>
            <span class="text-lg">{{ params.get('appName') }}</span>
          </div>
          <div>
            <strong class="mr-4">{{ t('appDescription') }}</strong>
            <span class="text-lg">{{ params.get('appDescription') }}</span>
          </div>
          <KeySelector v-model="selectedKey" class="pt-1" />
        </div>

        <div v-if="keysStore.keys.length > 0" class="flex justify-end mt-6 gap-2">
          <Button variant="outline" @click="registerAction(false)">
            {{ t('deny') }}
          </Button>
          <div
            :class="!selectedKey ? 'cursor-not-allowed!' : ''"
            :title="!selectedKey ? t('pleaseSelectKey') : ''"
          >
            <Button
              variant="default"
              :disabled="!selectedKey"
              @click="registerAction(true)"
            >
              {{ t('approve') }}
            </Button>
          </div>
        </div>
      </div>
    </DialogContent>
    <DialogContent v-else-if="actionType">
      <DialogHeader>
        <DialogTitle>{{ t('unrecognizedAction') }}</DialogTitle>
        <DialogDescription>
          {{ t('unrecognizedActionDescription', { actionType }) }}
        </DialogDescription>
      </DialogHeader>
    </DialogContent>
  </Dialog>
</template>

<i18n lang="yaml">
en-US:
  missingAppId: Missing application ID
  appPermissionRequest: Application Permission Request
  appPermissionDescription: This application is requesting access to your account information. Please confirm the following information is correct.
  appName: Application Name
  appDescription: Application Description
  deny: Deny
  pleaseSelectKey: Please select a provider
  approve: Approve
  unrecognizedAction: Unrecognized Action
  unrecognizedActionDescription: 'The current action type "{actionType}" is not supported or recognized.'

zh-CN:
  missingAppId: 缺少应用ID
  appPermissionRequest: 应用权限请求
  appPermissionDescription: 该应用请求访问您的账户信息。请确认以下信息是否正确。
  appName: 应用名称
  appDescription: 应用描述
  deny: 拒绝
  pleaseSelectKey: 请选择一个提供商
  approve: 同意
  unrecognizedAction: 无法识别的操作
  unrecognizedActionDescription: '当前操作类型 "{actionType}" 未被支持或识别。'
</i18n>
