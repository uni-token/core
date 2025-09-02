<script setup lang="ts">
import type { App } from '@/stores/app'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue-sonner'
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger } from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { useAppStore } from '@/stores/app'

const props = defineProps<{
  app: App
}>()
const { t } = useI18n()
const open = defineModel<boolean>('open', { default: false })

const appStore = useAppStore()
const loading = ref(false)

async function handleDeleteApp() {
  loading.value = true
  try {
    await appStore.deleteApp(props.app.id)
    toast.success(t('deleteSuccess'))
    open.value = false
  }
  catch (error) {
    console.error('Delete Failed:', error)
    toast.error(t('deleteFailed'))
  }
  finally {
    loading.value = false
  }
}

function formatDate(dateString: string) {
  if (!dateString)
    return t('unknown')
  try {
    return new Date(dateString).toLocaleString()
  }
  catch {
    return t('invalidDate')
  }
}
</script>

<template>
  <Dialog v-model:open="open">
    <DialogContent class="max-w-md">
      <DialogHeader>
        <DialogTitle>{{ t('title') }}</DialogTitle>
      </DialogHeader>

      <div class="space-y-4">
        <div class="space-y-3">
          <div class="flex items-center justify-between">
            <h3 class="font-semibold text-lg">
              {{ app.name }}
            </h3>
          </div>

          <p class="text-muted-foreground text-sm">
            {{ app.description || t('noDescription') }}
          </p>
        </div>

        <div class="space-y-3 border-t pt-4">
          <div class="text-sm">
            <div>
              <span class="text-gray-500 dark:text-gray-200">{{ t('appId') }}:</span>
              <p class="font-mono text-xs bg-gray-100 dark:bg-white/10 p-1 rounded mt-1 break-all">
                {{ app.id }}
              </p>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4 text-sm">
            <div>
              <span class="text-gray-500">{{ t('createdAt') }}:</span>
              <p class="mt-1">
                {{ formatDate(app.createdAt) }}
              </p>
            </div>
            <div v-if="app.lastActiveAt">
              <span class="text-gray-500">{{ t('lastActiveAt') }}:</span>
              <p class="mt-1">
                {{ formatDate(app.lastActiveAt) }}
              </p>
            </div>
          </div>
        </div>

        <div class="border-t border-red-200 dark:border-red-800 pt-4">
          <h4 class="text-sm font-medium text-red-800 dark:text-red-500 mb-2">
            {{ t('dangerZone') }}
          </h4>
          <AlertDialog>
            <AlertDialogTrigger as-child>
              <Button variant="destructive" size="sm" class="w-full bg-red-600/70! hover:bg-red-600/60!" :disabled="loading">
                {{ t('deleteApp') }}
              </Button>
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>{{ t('confirmDeleteTitle') }}</AlertDialogTitle>
                <AlertDialogDescription>
                  {{ t('confirmDeleteDescription', { appName: app.name }) }}
                </AlertDialogDescription>
              </AlertDialogHeader>
              <AlertDialogFooter>
                <AlertDialogCancel>{{ t('cancel') }}</AlertDialogCancel>
                <AlertDialogAction @click="handleDeleteApp">
                  {{ t('delete') }}
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        </div>
      </div>
    </DialogContent>
  </Dialog>
</template>

<i18n lang="yaml">
en-US:
  title: Application Details
  noDescription: No description
  appId: Application ID
  createdAt: Created At
  lastActiveAt: Last Active
  unknown: Unknown
  invalidDate: Invalid Date
  dangerZone: Danger Zone
  deleteApp: Delete Application
  confirmDeleteTitle: Confirm Delete Application
  confirmDeleteDescription: |
    Are you sure you want to delete application "{appName}"? This action cannot
    be undone and all application data will be permanently deleted.
  deleteSuccess: Application deleted successfully
  deleteFailed: Failed to delete application
  delete: Delete
  cancel: Cancel
zh-CN:
  title: 应用详细信息
  noDescription: 暂无描述
  appId: 应用ID
  createdAt: 创建时间
  lastActiveAt: 最后活跃
  unknown: 未知
  invalidDate: 无效日期
  dangerZone: 危险操作
  deleteApp: 删除应用
  confirmDeleteTitle: 确认删除应用
  confirmDeleteDescription: '您确定要删除应用 "{appName}" 吗？此操作不可撤销，应用的所有数据将被永久删除。'
  deleteSuccess: 应用删除成功
  deleteFailed: 删除应用失败
  delete: 删除
  cancel: 取消
</i18n>
