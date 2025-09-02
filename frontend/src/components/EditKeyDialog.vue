<script setup lang="ts">
import type { APIKey } from '@/stores'
import { Trash2 } from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger } from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { useKeysStore, usePresetsStore } from '@/stores'

interface EditConfig {
  name: string
  baseUrl: string
  token: string
}

const props = defineProps<{
  apiKey?: APIKey | null
}>()
const { t } = useI18n()
const open = defineModel<boolean>('open')
const keysStore = useKeysStore()
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
  if (!isConfigValid.value || !props.apiKey) {
    return
  }

  const success = await keysStore.updateKey(props.apiKey.id, {
    id: props.apiKey.id,
    name: config.value.name,
    type: props.apiKey.type,
    protocol: 'openai',
    baseUrl: config.value.baseUrl,
    token: config.value.token,
  })

  if (success) {
    open.value = false
  }
}

async function handleDelete() {
  if (!props.apiKey) {
    return
  }

  const success = await keysStore.deleteKey(props.apiKey.id)
  await presetsStore.loadPresets()

  if (success) {
    open.value = false
  }
}

watch(() => props.apiKey, (newKey) => {
  if (newKey) {
    config.value = {
      name: newKey.name,
      baseUrl: newKey.baseUrl,
      token: newKey.token,
    }
  }
}, { immediate: true, deep: true })
</script>

<template>
  <Dialog v-model:open="open">
    <DialogContent class="sm:max-w-lg">
      <DialogHeader>
        <DialogTitle>{{ t('editKeyDialog.title') }}</DialogTitle>
        <DialogDescription>
          {{ t('editKeyDialog.description') }}
        </DialogDescription>
      </DialogHeader>

      <div class="space-y-4">
        <div class="space-y-2">
          <label class="text-sm font-medium">{{ t('editKeyDialog.keyName') }}</label>
          <Input v-model="config.name" :placeholder="t('editKeyDialog.keyNamePlaceholder')" autocomplete="off" />
        </div>

        <div class="space-y-2">
          <label class="text-sm font-medium">{{ t('editKeyDialog.baseUrl') }}</label>
          <Input v-model="config.baseUrl" :placeholder="t('editKeyDialog.baseUrlPlaceholder')" autocomplete="off" />
        </div>

        <div class="space-y-2">
          <label class="text-sm font-medium">{{ t('editKeyDialog.apiKey') }}</label>
          <Input v-model="config.token" type="password" :placeholder="t('editKeyDialog.apiKeyPlaceholder')" autocomplete="new-password" />
        </div>

        <div class="flex gap-2 mt-6">
          <AlertDialog>
            <AlertDialogTrigger as-child>
              <Button variant="outline" class="text-red-600 hover:text-red-700">
                <Trash2 class="mr-2 h-4 w-4" />
                {{ t('editKeyDialog.delete') }}
              </Button>
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>{{ t('editKeyDialog.confirmDeleteTitle') }}</AlertDialogTitle>
                <AlertDialogDescription>
                  {{ t('editKeyDialog.confirmDeleteDescription') }}
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
            {{ t('editKeyDialog.saveChanges') }}
          </Button>
        </div>
      </div>
    </DialogContent>
  </Dialog>
</template>
