<script setup lang="ts">
import type { Provider } from '@/lib/providers'
import { shallowRef } from 'vue'
import { useI18n } from 'vue-i18n'
import ManualConfigCard from '@/components/ManualConfigCard.vue'
import ProviderConfigDialog from '@/components/ProviderConfigDialog.vue'
import ProviderName from '@/components/ProviderName.vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useProviders } from '@/lib/providers'

const { t } = useI18n()
const providers = useProviders()

const showProviderDialog = shallowRef<Provider | null>(null)
</script>

<template>
  <div class="space-y-6 flex-grow flex flex-col min-h-100">
    <div class="flex items-center justify-between">
      <h2 class="text-2xl font-bold">
        {{ t('title') }}
      </h2>
    </div>

    <div class="flex-grow flex flex-col">
      <div class="grid gap-4 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3">
        <Card v-for="provider in providers" :key="provider.id" class="relative gap-2 hover:bg-secondary" @click="showProviderDialog = provider">
          <CardHeader>
            <div class="flex items-center justify-between flex-wrap gap-2">
              <CardTitle class="flex ">
                <ProviderName :provider="provider" />
              </CardTitle>
              <div v-if="provider.user !== undefined" class="flex items-center gap-2 self-end">
                <span
                  class="inline-flex items-center rounded-full px-2 py-1 text-xs font-medium text-accent-foreground"
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
                {{ t('description1') }}
                <a :href="provider.homepage" target="_blank" class="text-blue-900 dark:text-blue-200 hover:underline" @click.stop>
                  {{ provider.name }}
                </a>
                {{ t('description2') }}
              </p>
            </div>
          </CardContent>
        </Card>

        <ManualConfigCard class="relative gap-2" />
      </div>
    </div>

    <ProviderConfigDialog
      v-if="showProviderDialog != null"
      :provider="showProviderDialog"
      @close="showProviderDialog = null"
    />
  </div>
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
  loggedIn: Logged In
  loggedOut: Logged Out
</i18n>
