<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'

export interface ManualConfig {
  name: string
  protocol: 'openai'
  baseUrl: string
  token: string
}

const emit = defineEmits<{
  save: [config: ManualConfig]
}>()
const { t } = useI18n()
const open = defineModel<boolean>('open')

const config = ref<ManualConfig>({
  name: '',
  protocol: 'openai',
  baseUrl: '',
  token: '',
})

const isConfigValid = computed(() => {
  return config.value.name && config.value.baseUrl && config.value.token
})

function handleSave() {
  if (isConfigValid.value) {
    emit('save', { ...config.value })
    resetConfig()
  }
}

function resetConfig() {
  config.value = {
    name: '',
    protocol: 'openai',
    baseUrl: '',
    token: '',
  }
}

watch(open, (newValue) => {
  if (!newValue) {
    resetConfig()
  }
})
</script>

<template>
  <Dialog v-model:open="open">
    <DialogContent class="sm:max-w-lg">
      <DialogHeader>
        <DialogTitle>{{ t('title') }}</DialogTitle>
        <DialogDescription>
          {{ t('description') }}
        </DialogDescription>
      </DialogHeader>

      <div class="mt-6 space-y-4">
        <div class="space-y-2">
          <label class="text-sm font-medium">{{ t('providerName') }}</label>
          <Input
            v-model="config.name"
            :placeholder="t('providerNamePlaceholder')"
            autocomplete="off"
            data-lpignore="true"
            data-form-type="other"
          />
        </div>

        <div class="space-y-2">
          <label class="text-sm font-medium">{{ t('baseUrl') }}</label>
          <Input
            v-model="config.baseUrl"
            :placeholder="t('baseUrlPlaceholder')"
            autocomplete="off"
            data-lpignore="true"
            data-form-type="other"
          />
        </div>

        <div class="space-y-2">
          <label class="text-sm font-medium">{{ t('apiKey') }}</label>
          <Input
            v-model="config.token"
            type="password"
            :placeholder="t('apiKeyPlaceholder')"
            autocomplete="off"
            data-form-type="other"
            data-lpignore="true"
          />
        </div>

        <Button class="w-full" :disabled="!isConfigValid" @click="handleSave">
          {{ t('save') }}
        </Button>
      </div>
    </DialogContent>
  </Dialog>
</template>

<i18n lang="yaml">
zh-CN:
  title: 手动配置 API
  description: 请输入您的 API 配置信息
  providerName: 提供商名称
  providerNamePlaceholder: 例如：OpenAI、Claude 等
  baseUrl: Base URL
  baseUrlPlaceholder: 例如：https://api.openai.com/v1
  apiKey: API Key
  apiKeyPlaceholder: 输入您的 API Key
  save: 保存配置
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
