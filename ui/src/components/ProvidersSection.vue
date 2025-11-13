<script setup lang="ts">
import type { Provider } from '@/lib/providers'
import { shallowRef } from 'vue'
import { useI18n } from 'vue-i18n'
import ProviderConfigDialog from '@/components/ProviderConfigDialog.vue'
import ProviderName from '@/components/ProviderName.vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useProviders } from '@/lib/providers'
import ProvidersListCard from './ProvidersListCard.vue'

const { t } = useI18n()
const providers = useProviders()

const showProviderDialog = shallowRef<Provider | null>(null)
</script>

<template>
  <div class="grid gap-4 grid-cols-1 sm:grid-cols-2 lg:grid-cols-4">
    <Card v-for="provider in providers" :key="provider.id" class="relative gap-2 hover:bg-secondary" @click="showProviderDialog = provider">
      <CardHeader>
        <div class="flex items-center justify-between flex-wrap gap-2">
          <CardTitle class="flex text-lg min-w-34">
            <ProviderName :provider="provider" />
          </CardTitle>
          <div v-if="provider.user !== undefined" class="flex items-center gap-2">
            <span
              class="inline-flex items-center rounded-full px-2 py-1 text-xs font-medium text-accent-foreground select-none"
              :class="{
                'bg-emerald-200 dark:bg-green-500/60': !!provider.user,
                'bg-accent': !provider.user,
              }"
            >
              {{ provider.user ? t('loggedIn') : t('loggedOut') }}
            </span>
          </div>
        </div>
      </CardHeader>

      <CardContent>
        <div class="text-sm text-muted-foreground">
          <p>
            {{ provider.description }}
          </p>
        </div>
      </CardContent>
    </Card>

    <ProvidersListCard />
  </div>

  <ProviderConfigDialog
    v-if="showProviderDialog != null"
    :provider="showProviderDialog"
    @close="showProviderDialog = null"
  />
</template>

<i18n lang="yaml">
zh-CN:
  title: 模型供应商
  loadFailed: 加载失败
  description1: 通过
  description2: 购买和配置 API
  loggedIn: 已登录
  loggedOut: 未登录

en-US:
  title: Model Providers
  loadFailed: Failed to load
  description1: Purchase and configure API through
  description2: ''
  loggedIn: Configured
  loggedOut: Configure
</i18n>
