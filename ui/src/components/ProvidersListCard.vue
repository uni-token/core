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
      <Card class="relative gap-2 hover:bg-secondary text-left h-full">
        <CardHeader>
          <CardTitle>
            <div class="flex gap-2 items-center text-lg">
              <ListPlusIcon class="w-6 mb-1" />
              {{ t('title') }}
            </div>
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div class="text-sm text-muted-foreground">
            {{ t('description') }}
          </div>
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
  description: 选择或自定义您的模型提供商
  manual: 自定义提供商
en-US:
  title: Other Providers
  description: Select a provider or customize your own
  manual: Custom Provider
</i18n>
