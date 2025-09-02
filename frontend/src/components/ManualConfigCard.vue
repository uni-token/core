<script setup lang="ts">
import type { ManualConfig } from '@/components/ManualConfigDialog.vue'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue-sonner'
import ManualConfigDialog from '@/components/ManualConfigDialog.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { useKeysStore } from '@/stores'

const emit = defineEmits<{
  configured: [key: string]
}>()

const { t } = useI18n()
const keysStore = useKeysStore()

const showEditDialog = ref(false)

async function handleManualConfigSave(config: ManualConfig) {
  const key = await keysStore.addKey({
    name: config.name,
    type: 'manual',
    protocol: config.protocol,
    baseUrl: config.baseUrl,
    token: config.token,
  })

  emit('configured', key.id)
  showEditDialog.value = false
  toast.success(t('success'))
}
</script>

<template>
  <Card>
    <CardHeader>
      <div class="flex items-center justify-between">
        <CardTitle class="text-lg">
          {{ t('title') }}
        </CardTitle>
      </div>
    </CardHeader>

    <CardContent class="flex-grow">
      <div class="text-sm text-muted-foreground">
        <p>{{ t('description') }}</p>
      </div>
    </CardContent>

    <CardFooter>
      <Button class="w-full" @click="showEditDialog = true">
        {{ t('add') }}
      </Button>
    </CardFooter>

    <ManualConfigDialog
      v-model:open="showEditDialog"
      @save="handleManualConfigSave"
    />
  </Card>
</template>

<i18n lang="yaml">
zh-CN:
  title: 手动配置
  description: 手动输入 Base URL 和 API Key
  add: 添加手动配置
  success: 配置成功
en-US:
  title: Manual Configuration
  description: Manually enter Base URL and API Key
  add: Add Manual Configuration
  success: Configuration Successful
</i18n>
