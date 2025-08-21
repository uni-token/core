<script setup lang="ts">
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

function download(os: 'windows' | 'macos' | 'linux') {
  const filename = {
    linux: 'service-linux-amd64',
    macos: 'service-darwin-amd64',
    windows: 'service-windows-amd64.exe',
  }[os]
  const url = `https://uni-token.app/release/${filename}`
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  link.click()
  link.remove()
}

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
      <Button
        size="sm"
        variant="link"
        class="inline underline px-0 opacity-80 hover:opacity-100 h-6"
        :class="props.class"
      >
        {{ t('service.tryStart') }}
      </Button>
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
                        size="sm"
                        variant="secondary"
                        class="w-full"
                        @click="download('windows')"
                      >
                        Windows
                      </Button>
                      <Button
                        size="sm"
                        variant="secondary"
                        class="w-full"
                        @click="download('macos')"
                      >
                        macOS
                      </Button>
                      <Button
                        size="sm"
                        variant="secondary"
                        class="w-full"
                        @click="download('linux')"
                      >
                        Linux
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
          @click="serviceStore.tryStartService"
        >
          {{ t('common.confirm') }}
        </AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>
