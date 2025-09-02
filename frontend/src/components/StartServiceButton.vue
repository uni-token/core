<script setup lang="ts">
import { downloadService, startService, SUPPORTED_OS } from '@uni-token/browser-sdk'
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

    toast.success(t('startSuccess'))
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
          {{ t('tryStart') }}
        </Button>
      </slot>
    </AlertDialogTrigger>
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle>
          {{ t('startLocalService') }}
        </AlertDialogTitle>
        <AlertDialogDescription>
          <p>
            {{ t('browserPrompt') }}
          </p>
          <p>
            {{ t('adminPermission') }}
          </p>
          <p>
            {{ t('troubleshootingPrefix') }}
            <Dialog v-model:open="showDownloadDialog">
              <DialogTrigger class="inline">
                <Button
                  size="sm"
                  variant="link"
                  class="inline! text-sm underline px-0 opacity-80 hover:opacity-100"
                >
                  {{ t('manualDownload') }}
                </Button>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>
                    {{ t('downloadService') }}
                  </DialogTitle>
                  <DialogDescription>
                    {{ t('selectOS') }}

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

                    {{ t('afterDownload') }}
                  </DialogDescription>
                </DialogHeader>
              </DialogContent>
            </Dialog>
          </p>
        </AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <AlertDialogCancel>{{ t('cancel') }}</AlertDialogCancel>
        <AlertDialogAction
          @click="startService()"
        >
          {{ t('confirm') }}
        </AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>

<i18n lang="yaml">
en-US:
  startSuccess: Service started successfully!
  tryStart: Start Service
  startLocalService: Start Local Service
  browserPrompt: In the browser prompt that appears, please select to allow opening the UniToken application.
  adminPermission: When starting for the first time, UniToken will request administrator permissions. Please allow.
  troubleshootingPrefix: 'If the service fails to start properly, please try restarting the related Agent software, or '
  manualDownload: download manually
  downloadService: Download UniToken Service
  selectOS: Please select your operating system.
  afterDownload: 'After downloading, please run it to start the service.'
  cancel: Cancel
  confirm: Confirm

zh-CN:
  startSuccess: 服务已成功启动！
  tryStart: 启动服务
  startLocalService: 启动本地服务
  browserPrompt: 在浏览器弹出的提示框中，请选择允许打开 UniToken 应用。
  adminPermission: 首次启动时，UniToken 会请求管理员权限，请允许。
  troubleshootingPrefix: 若服务未能正常启动，请尝试重新启动相关 Agent 软件，或
  manualDownload: 手动下载
  downloadService: 下载 UniToken 服务
  selectOS: 请选择你的操作系统。
  afterDownload: 下载后，请运行它以启动服务。
  cancel: 取消
  confirm: 确认
</i18n>
