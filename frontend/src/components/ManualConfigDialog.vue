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
        <DialogTitle>{{ t('manualConfigDialog.title') }}</DialogTitle>
        <DialogDescription>
          {{ t('manualConfigDialog.description') }}
        </DialogDescription>
      </DialogHeader>

      <div class="mt-6 space-y-4">
        <div class="space-y-2">
          <label class="text-sm font-medium">{{ t('manualConfigDialog.providerName') }}</label>
          <Input
            v-model="config.name"
            :placeholder="t('manualConfigDialog.providerNamePlaceholder')"
            autocomplete="off"
            data-lpignore="true"
            data-form-type="other"
          />
        </div>

        <div class="space-y-2">
          <label class="text-sm font-medium">{{ t('manualConfigDialog.baseUrl') }}</label>
          <Input
            v-model="config.baseUrl"
            :placeholder="t('manualConfigDialog.baseUrlPlaceholder')"
            autocomplete="off"
            data-lpignore="true"
            data-form-type="other"
          />
        </div>

        <div class="space-y-2">
          <label class="text-sm font-medium">{{ t('manualConfigDialog.apiKey') }}</label>
          <Input
            v-model="config.token"
            type="password"
            :placeholder="t('manualConfigDialog.apiKeyPlaceholder')"
            autocomplete="off"
            data-form-type="other"
            data-lpignore="true"
          />
        </div>

        <Button class="w-full" :disabled="!isConfigValid" @click="handleSave">
          {{ t('manualConfigDialog.saveProvider') }}
        </Button>
      </div>
    </DialogContent>
  </Dialog>
</template>
