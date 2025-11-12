<script setup lang="ts">
import type { Provider } from '@/lib/providers'
import { ListPlusIcon } from 'lucide-vue-next'
import { shallowRef } from 'vue'
import { useI18n } from 'vue-i18n'
import ProviderName from '@/components/ProviderName.vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'
import { useProviders } from '@/lib/providers'
import ManualConfigDialog from './ManualConfigDialog.vue'
import ProviderConfigDialog from './ProviderConfigDialog.vue'

const { t } = useI18n()

const providers = useProviders()

const showProviderDialog = shallowRef<Provider | null>(null)
</script>

<template>
  <Dialog>
    <DialogTrigger>
      <Card class="relative gap-2 hover:bg-secondary text-left">
        <CardHeader>
          <CardTitle>
            {{ t('title') }}
          </CardTitle>
        </CardHeader>
        <CardContent>
          123
        </CardContent>
      </Card>
    </DialogTrigger>
    <DialogContent class="sm:max-w-lg">
      <DialogHeader>
        <DialogTitle>{{ t('title') }}</DialogTitle>
        <DialogDescription>
          {{ t('description') }}
        </DialogDescription>
      </DialogHeader>

      <div class="flex flex-col gap-2">
        <div v-for="provider in providers" :key="provider.id" class="p-4 border rounded-lg hover:bg-secondary" @click="showProviderDialog = provider">
          <ProviderName :provider="provider" class="font-medium" />
        </div>
        <ManualConfigDialog>
          <template #trigger>
            <div class="p-4 border rounded-lg hover:bg-secondary">
              <div class="flex gap-2 items-center">
                <ListPlusIcon class="w-5 mb-1 rounded" />
                {{ t('manual') }}
              </div>
            </div>
          </template>
        </ManualConfigDialog>
      </div>
    </DialogContent>
  </Dialog>

  <ProviderConfigDialog
    v-if="showProviderDialog != null"
    :provider="showProviderDialog"
    @close="showProviderDialog = null"
  />
</template>

<i18n lang="yaml">
zh-CN:
  title: 其他提供商
  description: 请输入您的 API 配置信息
  manual: 手动配置
en-US:
  title: Manual API Configuration
  description: Please enter your API configuration information
  providerName: Provider Name
  providerNamePlaceholder: e.g., OpenAI, Claude, etc.
  baseUrl: Base URL
  baseUrlPlaceholder: e.g., https://api.openai.com/v1
  apiKey: API Key
  apiKeyPlaceholder: Enter your API Key
  save: Save Configuration
</i18n>
