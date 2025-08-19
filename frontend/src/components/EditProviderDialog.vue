<script setup lang="ts">
import type { LLMProvider } from '@/stores'
import { Trash2 } from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger } from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { usePresetsStore, useProvidersStore } from '@/stores'

interface EditConfig {
  name: string
  baseUrl: string
  token: string
}

const props = defineProps<{
  provider?: LLMProvider | null
}>()
const { t } = useI18n()
const open = defineModel<boolean>('open')
const providersStore = useProvidersStore()
const presetsStore = usePresetsStore()

const config = ref<EditConfig>({
  name: '',
  baseUrl: '',
  token: '',
})

const isConfigValid = computed(() => {
  return config.value.name && config.value.baseUrl && config.value.token
})

async function handleSave() {
  if (!isConfigValid.value || !props.provider) {
    return
  }

  const success = await providersStore.updateProvider(props.provider.id, {
    id: props.provider.id,
    name: config.value.name,
    type: props.provider.type,
    protocol: 'openai',
    baseUrl: config.value.baseUrl,
    token: config.value.token,
  })

  if (success) {
    open.value = false
  }
}

async function handleDelete() {
  if (!props.provider) {
    return
  }

  const success = await providersStore.deleteProvider(props.provider.id)
  await presetsStore.loadPresets()

  if (success) {
    open.value = false
  }
}

watch(() => props.provider, (newProvider) => {
  if (newProvider) {
    config.value = {
      name: newProvider.name,
      baseUrl: newProvider.baseUrl,
      token: newProvider.token,
    }
  }
}, { immediate: true, deep: true })
</script>

<template>
  <Dialog v-model:open="open">
    <DialogContent class="sm:max-w-lg">
      <DialogHeader>
        <DialogTitle>{{ t('editProviderDialog.title') }}</DialogTitle>
        <DialogDescription>
          {{ t('editProviderDialog.description') }}
        </DialogDescription>
      </DialogHeader>

      <div class="space-y-4">
        <div class="space-y-2">
          <label class="text-sm font-medium">{{ t('editProviderDialog.providerName') }}</label>
          <Input v-model="config.name" :placeholder="t('editProviderDialog.providerNamePlaceholder')" autocomplete="off" />
        </div>

        <div class="space-y-2">
          <label class="text-sm font-medium">{{ t('editProviderDialog.baseUrl') }}</label>
          <Input v-model="config.baseUrl" :placeholder="t('editProviderDialog.baseUrlPlaceholder')" autocomplete="off" />
        </div>

        <div class="space-y-2">
          <label class="text-sm font-medium">{{ t('editProviderDialog.apiKey') }}</label>
          <Input v-model="config.token" type="password" :placeholder="t('editProviderDialog.apiKeyPlaceholder')" autocomplete="new-password" />
        </div>

        <div class="flex gap-2 mt-6">
          <AlertDialog>
            <AlertDialogTrigger as-child>
              <Button variant="outline" class="text-red-600 hover:text-red-700">
                <Trash2 class="mr-2 h-4 w-4" />
                {{ t('editProviderDialog.delete') }}
              </Button>
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>{{ t('editProviderDialog.confirmDeleteTitle') }}</AlertDialogTitle>
                <AlertDialogDescription>
                  {{ t('editProviderDialog.confirmDeleteDescription') }}
                </AlertDialogDescription>
              </AlertDialogHeader>
              <AlertDialogFooter>
                <AlertDialogCancel tabindex="1">
                  {{ t('common.cancel') }}
                </AlertDialogCancel>
                <AlertDialogAction tabindex="2" @click="handleDelete">
                  {{ t('common.delete') }}
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
          <Button class="flex-1" :disabled="!isConfigValid" @click="handleSave">
            {{ t('editProviderDialog.saveChanges') }}
          </Button>
        </div>
      </div>
    </DialogContent>
  </Dialog>
</template>
