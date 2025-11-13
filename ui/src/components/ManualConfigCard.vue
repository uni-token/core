<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import ManualConfigDialog from '@/components/ManualConfigDialog.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'

const emit = defineEmits<{
  configured: [key: string]
}>()

const { t } = useI18n()

const showEditDialog = ref(false)

async function handleManualConfigSave(config: string) {
  emit('configured', config)
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
      <ManualConfigDialog
        v-model:open="showEditDialog"
        @save="handleManualConfigSave"
      >
        <template #trigger>
          <Button class="w-full" variant="secondary" @click="showEditDialog = true">
            {{ t('add') }}
          </Button>
        </template>
      </ManualConfigDialog>
    </CardFooter>
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
