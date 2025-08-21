<script setup lang="ts">
import { downloadService, startService, SUPPORTED_OS } from '@uni-token/web-sdk'
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue-sonner'
import { useServiceStore } from '@/stores'
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger } from './ui/alert-dialog'
import { Button } from './ui/button'
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from './ui/dialog'

const props = defineProps<{
  class?: string
}>()

const { t } = useI18n()
const serviceStore = useServiceStore()

const showStartServiceDialog = ref(false)
const showDownloadDialog = ref(false)

watch(() => serviceStore.serverConnected, (connected) => {
  if (connected) {
    showStartServiceDialog.value = false
    showDownloadDialog.value = false

    toast.success(t('service.startSuccess'))
  }
})
</script>

<template>
  <AlertDialog v-model:open="showStartServiceDialog">
    <AlertDialogTrigger as-child>
      <slot>
        <Button
          size="sm"
          variant="link"
          class="inline underline px-0 opacity-80 hover:opacity-100 h-6"
          :class="props.class"
        >
          {{ t('service.tryStart') }}
        </Button>
      </slot>
    </AlertDialogTrigger>
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle>
          {{ t('service.startLocalService') }}
        </AlertDialogTitle>
        <AlertDialogDescription>
          <p>
            {{ t('service.browserPrompt') }}
          </p>
          <p>
            {{ t('service.adminPermission') }}
          </p>
          <p>
            {{ t('service.troubleshootingPrefix') }}
            <Dialog v-model:open="showDownloadDialog">
              <DialogTrigger class="inline">
                <Button
                  size="sm"
                  variant="link"
                  class="inline! text-sm underline px-0 opacity-80 hover:opacity-100"
                >
                  {{ t('service.manualDownload') }}
                </Button>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>
                    {{ t('service.downloadService') }}
                  </DialogTitle>
                  <DialogDescription>
                    {{ t('service.selectOS') }}

                    <div class="flex flex-col gap-2 my-4">
                      <Button
                        v-for="os in SUPPORTED_OS"
                        :key="os"
                        size="sm"
                        variant="secondary"
                        class="w-full"
                        @click="downloadService(os)"
                      >
                        {{ os }}
                      </Button>
                    </div>

                    {{ t('service.afterDownload') }}
                  </DialogDescription>
                </DialogHeader>
              </DialogContent>
            </Dialog>
          </p>
        </AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <AlertDialogCancel>{{ t('common.cancel') }}</AlertDialogCancel>
        <AlertDialogAction
          @click="startService()"
        >
          {{ t('common.confirm') }}
        </AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>
