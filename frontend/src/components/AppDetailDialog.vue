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
    toast.success(t('appDetailDialog.deleteSuccess'))
    open.value = false
  }
  catch (error) {
    console.error('Delete Failed:', error)
    toast.error(t('appDetailDialog.deleteFailed'))
  }
  finally {
    loading.value = false
  }
}

function formatDate(dateString: string) {
  if (!dateString)
    return t('appDetailDialog.unknown')
  try {
    return new Date(dateString).toLocaleString()
  }
  catch {
    return t('appDetailDialog.invalidDate')
  }
}
</script>

<template>
  <Dialog v-model:open="open">
    <DialogContent class="max-w-md">
      <DialogHeader>
        <DialogTitle>{{ t('appDetailDialog.title') }}</DialogTitle>
      </DialogHeader>

      <div class="space-y-4">
        <div class="space-y-3">
          <div class="flex items-center justify-between">
            <h3 class="font-semibold text-lg">
              {{ app.name }}
            </h3>
          </div>

          <p class="text-muted-foreground text-sm">
            {{ app.description || t('appDetailDialog.noDescription') }}
          </p>
        </div>

        <div class="space-y-3 border-t pt-4">
          <div class="text-sm">
            <div>
              <span class="text-gray-500 dark:text-gray-200">{{ t('appDetailDialog.appId') }}:</span>
              <p class="font-mono text-xs bg-gray-100 dark:bg-white/10 p-1 rounded mt-1 break-all">
                {{ app.id }}
              </p>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4 text-sm">
            <div>
              <span class="text-gray-500">{{ t('appDetailDialog.createdAt') }}:</span>
              <p class="mt-1">
                {{ formatDate(app.createdAt) }}
              </p>
            </div>
            <div v-if="app.lastActiveAt">
              <span class="text-gray-500">{{ t('appDetailDialog.lastActiveAt') }}:</span>
              <p class="mt-1">
                {{ formatDate(app.lastActiveAt) }}
              </p>
            </div>
          </div>
        </div>

        <div class="border-t border-red-200 dark:border-red-800 pt-4">
          <h4 class="text-sm font-medium text-red-800 dark:text-red-500 mb-2">
            {{ t('appDetailDialog.dangerZone') }}
          </h4>
          <AlertDialog>
            <AlertDialogTrigger as-child>
              <Button variant="destructive" size="sm" class="w-full bg-red-600/70! hover:bg-red-600/60!" :disabled="loading">
                {{ t('appDetailDialog.deleteApp') }}
              </Button>
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>{{ t('appDetailDialog.confirmDeleteTitle') }}</AlertDialogTitle>
                <AlertDialogDescription>
                  {{ t('appDetailDialog.confirmDeleteDescription', { appName: app.name }) }}
                </AlertDialogDescription>
              </AlertDialogHeader>
              <AlertDialogFooter>
                <AlertDialogCancel>{{ t('common.cancel') }}</AlertDialogCancel>
                <AlertDialogAction @click="handleDeleteApp">
                  {{ t('common.delete') }}
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        </div>
      </div>
    </DialogContent>
  </Dialog>
</template>
