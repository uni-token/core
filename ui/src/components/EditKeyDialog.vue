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
        <DialogTitle>{{ t('title') }}</DialogTitle>
        <DialogDescription>
          {{ t('description') }}
        </DialogDescription>
      </DialogHeader>

      <div class="space-y-4">
        <div class="space-y-2">
          <label class="text-sm font-medium">{{ t('name') }}</label>
          <Input v-model="config.name" :placeholder="t('namePlaceholder')" autocomplete="off" />
        </div>

        <div class="space-y-2">
          <label class="text-sm font-medium">{{ t('baseUrl') }}</label>
          <Input v-model="config.baseUrl" :placeholder="t('baseUrlPlaceholder')" autocomplete="off" />
        </div>

        <div class="space-y-2">
          <label class="text-sm font-medium">{{ t('apiKey') }}</label>
          <Input v-model="config.token" type="password" :placeholder="t('apiKeyPlaceholder')" autocomplete="new-password" />
        </div>

        <div class="flex gap-2 mt-6">
          <AlertDialog>
            <AlertDialogTrigger as-child>
              <Button variant="outline" class="text-red-600 hover:text-red-700">
                <Trash2 class="mr-2 h-4 w-4" />
                {{ t('delete') }}
              </Button>
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>{{ t('confirmDeleteTitle') }}</AlertDialogTitle>
                <AlertDialogDescription>
                  {{ t('confirmDeleteDescription') }}
                </AlertDialogDescription>
              </AlertDialogHeader>
              <AlertDialogFooter>
                <AlertDialogCancel tabindex="1">
                  {{ t('cancel') }}
                </AlertDialogCancel>
                <AlertDialogAction tabindex="2" @click="handleDelete">
                  {{ t('delete') }}
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
          <div class="flex-1" />
          <Button :disabled="!isConfigValid" @click="handleSave">
            {{ t('save') }}
          </Button>
        </div>
      </div>
    </DialogContent>
  </Dialog>
</template>

<i18n lang="yaml">
en-US:
  title: Edit API Key
  description: Modify API Key configuration information
  name: Name
  namePlaceholder: 'e.g.: OpenAI'
  baseUrl: Base URL
  baseUrlPlaceholder: 'https://api.openai.com/v1'
  apiKey: API Key
  apiKeyPlaceholder: sk-...
  delete: Delete
  confirmDeleteTitle: Confirm Delete
  confirmDeleteDescription: Are you sure you want to delete this Provider? This action cannot be undone.
  cancel: Cancel
  save: Save

zh-CN:
  title: 编辑 API Key
  description: 修改 API Key 的配置信息
  name: 名称
  namePlaceholder: '例如: OpenAI'
  baseUrl: Base URL
  baseUrlPlaceholder: 'https://api.openai.com/v1'
  apiKey: API Key
  apiKeyPlaceholder: sk-...
  delete: 删除
  confirmDeleteTitle: 确认删除
  confirmDeleteDescription: 确定要删除 Provider 吗？此操作无法撤销。
  cancel: 取消
  save: 保存
</i18n>
