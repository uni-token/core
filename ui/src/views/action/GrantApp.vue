<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import KeySelector from '@/components/KeySelector.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { useKeysStore } from '@/stores'
import { useAppStore } from '@/stores/app'
import LogoSvg from '/logo.svg?raw'

const { t } = useI18n()
const router = useRouter()

const query = new URLSearchParams(window.location.search)

const selectedKey = ref<string>('')
const appStore = useAppStore()
const keysStore = useKeysStore()

async function registerAction(granted: boolean) {
  const appId = query.get('appId')

  if (!appId) {
    toast.error(t('missingAppId'))
    return
  }

  await appStore.toggleAppAuthorization(appId, granted, selectedKey.value)
  router.replace('/')
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center p-6 bg-gradient-to-br from-background to-muted/20">
    <Card class="w-full max-w-2xl shadow-lg">
      <CardHeader class="text-center space-y-4 pb-8">
        <div class="flex justify-center">
          <div class="rounded-2xl bg-primary p-3 text-primary-foreground shadow-md">
            <div class="h-10 w-10" v-html="LogoSvg" />
          </div>
        </div>
        <div>
          <CardTitle class="text-3xl font-bold">
            {{ t('appPermissionRequest') }}
          </CardTitle>
          <CardDescription class="text-base mt-2">
            {{ t('appPermissionDescription') }}
          </CardDescription>
        </div>
      </CardHeader>

      <CardContent class="space-y-6">
        <div class="space-y-4 bg-muted/30 rounded-lg p-6">
          <div class="flex items-baseline gap-3">
            <span class="text-sm font-medium text-muted-foreground min-w-20">{{ t('appName') }}</span>
            <span class="text-lg font-semibold">{{ query.get('appName') || '-' }}</span>
          </div>
          <div class="flex items-baseline gap-3">
            <span class="text-sm font-medium text-muted-foreground min-w-20">{{ t('appDescription') }}</span>
            <span class="text-base">{{ query.get('appDescription') || '-' }}</span>
          </div>
        </div>

        <KeySelector v-model="selectedKey" compact />

        <div class="flex gap-3 pt-4">
          <Button
            variant="outline"
            size="lg"
            class="flex-1"
            @click="registerAction(false)"
          >
            {{ t('deny') }}
          </Button>
          <Button
            variant="default"
            size="lg"
            class="flex-1"
            :disabled="!selectedKey || keysStore.keys.length === 0"
            @click="registerAction(true)"
          >
            {{ t('approve') }}
          </Button>
        </div>
      </CardContent>
    </Card>
  </div>
</template>

<i18n lang="yaml">
en-US:
  missingAppId: Missing application ID
  appPermissionRequest: Application Authorization
  appPermissionDescription: This application is requesting access to use your AI provider keys
  appName: App Name
  appDescription: Description
  deny: Deny
  approve: Approve

zh-CN:
  missingAppId: 缺少应用ID
  appPermissionRequest: 应用授权请求
  appPermissionDescription: 该应用请求使用您的 AI 提供商密钥
  appName: 应用名称
  appDescription: 应用描述
  deny: 拒绝
  approve: 同意
</i18n>
