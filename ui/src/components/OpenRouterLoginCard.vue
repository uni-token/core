<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { toast } from 'vue-sonner'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { useOpenRouterProvider } from '@/lib/providers/openrouter'
import { useKeysStore, useServiceStore } from '@/stores'

const { t } = useI18n()
const { api: fetch } = useServiceStore()
const provider = useOpenRouterProvider()
const keysStore = useKeysStore()

function handleLogin() {
  const callbackUrl = `${window.location.protocol}//${window.location.host}/action/openrouter-auth`
  const authUrl = `https://openrouter.ai/auth?callback_url=${encodeURIComponent(callbackUrl)}`
  window.open(authUrl, '_blank', 'width=600,height=600,noopener=yes,noreferrer=yes')

  const bc = new BroadcastChannel('openrouter-auth')
  bc.onmessage = async (event) => {
    if (event.data.code) {
      toast.success(t('authorizationSuccessful'))

      const response = await window.fetch('https://openrouter.ai/api/v1/auth/keys', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          code: event.data.code,
        }),
      })
      const { key, user_id } = await response.json()

      await fetch('openrouter/authed', {
        method: 'POST',
        body: JSON.stringify({ key, userId: user_id }),
      })
      await provider.refreshUser()
      await keysStore.createAndAddKey(provider)
    }
  }
}
</script>

<template>
  <Card>
    <CardHeader class="pb-4">
      <CardTitle class="text-lg">
        {{ t('loginTitle') }}
      </CardTitle>
      <CardDescription class="text-sm">
        {{ t('loginDescription') }}
      </CardDescription>
    </CardHeader>
    <CardContent class="space-y-4">
      <Button
        class="w-full h-10"
        @click="handleLogin"
      >
        {{ t('loginButton') }}
      </Button>
    </CardContent>
  </Card>
</template>

<i18n lang="yaml">
en-US:
  loginTitle: Login to OpenRouter
  loginDescription: Connect your OpenRouter account to access AI models
  loginButton: Login with OpenRouter
  popupBlocked: Popup was blocked. Please allow popups for this site.
  authorizationSuccessful: Authorization Successful

zh-CN:
  loginTitle: 登录 OpenRouter
  loginDescription: 连接您的 OpenRouter 账户以访问 AI 模型
  loginButton: 使用 OpenRouter 登录
  popupBlocked: 弹出窗口被阻止。请允许该网站弹出窗口。
  authorizationSuccessful: 授权成功
</i18n>
