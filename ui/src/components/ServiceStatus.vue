<script setup lang="ts">
import { AlertTriangle } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'
import {
  HoverCard,
  HoverCardContent,
  HoverCardTrigger,
} from '@/components/ui/hover-card'
import { useServiceStore } from '@/stores'
import StartServiceButton from './StartServiceButton.vue'

defineProps<{
  variant: 'full' | 'compact'
}>()

const { t } = useI18n()
const serviceStore = useServiceStore()
</script>

<template>
  <!-- Full variant with warning message (used in sidebar) -->
  <div v-if="variant === 'full'" class="space-y-2">
    <div v-if="!serviceStore.serverConnected">
      <div class="rounded-lg border border-orange-200 bg-orange-50 p-2 dark:border-orange-900 dark:bg-orange-950">
        <div class="flex items-start gap-2">
          <AlertTriangle class="h-4 w-4 text-orange-600 dark:text-orange-400 mt-0.5 flex-shrink-0" />
          <div class="space-y-1">
            <p class="text-xs font-medium text-orange-800 dark:text-orange-200">
              {{ t('connectionFailed') }}
            </p>
            <p class="text-xs text-orange-700 dark:text-orange-300">
              {{ t('reconnecting') }}
            </p>
            <p class="text-xs text-orange-700 dark:text-orange-300">
              {{ t('restartAgent') }}<br>{{ t('or') }}
              <StartServiceButton class="text-xs!" />
            </p>
          </div>
        </div>
      </div>
    </div>

    <div class="flex items-center justify-between">
      <HoverCard :open-delay="400">
        <HoverCardTrigger>
          <div class="h-8 flex items-center gap-1 text-xs text-sidebar-foreground/70">
            <div
              class="mx-1 h-2 w-2 rounded-full mb-[2px]"
              :class="serviceStore.serverConnected ? 'bg-green-500' : 'bg-red-500'"
            />
            <span class="text-sm select-none">{{ serviceStore.serverConnected ? t('connected') : t('disconnected') }}</span>
          </div>
        </HoverCardTrigger>
        <HoverCardContent v-if="serviceStore.serverConnected" :side-offset="0" class="pt-2 py-1">
          <div class="font-mono text-sm text-center mt-1">
            {{ serviceStore.serviceHost }}
          </div>
        </HoverCardContent>
      </HoverCard>
      <slot name="actions" />
    </div>
  </div>

  <!-- Compact variant (used in login form) -->
  <div v-else-if="variant === 'compact'" class="space-y-2">
    <div class="flex items-center justify-center gap-2 text-sm">
      <div
        class="h-2 w-2 rounded-full"
        :class="serviceStore.serverConnected ? 'bg-green-500' : 'bg-red-500'"
      />
      <span :class="serviceStore.serverConnected ? 'text-green-700 dark:text-green-400' : 'text-red-700 dark:text-red-400'">
        {{ serviceStore.serverConnected ? t('connected') : t('disconnected') }}
      </span>
    </div>

    <!-- Permanent troubleshooting message when disconnected -->
    <div v-if="!serviceStore.serverConnected" class="text-center space-y-2">
      <p class="text-sm text-orange-800 dark:text-orange-200">
        {{ t('troubleshooting') }}
        <StartServiceButton />
      </p>
    </div>
  </div>
</template>

<i18n lang="yaml">
zh-CN:
  connected: 本地服务已连接
  disconnected: 本地服务未连接
  connectionFailed: 服务连接失败
  reconnecting: 正在重新连接...
  restartAgent: 可能是本地服务未能正常启动，请重新启动相关 Agent 软件，
  troubleshooting: 请重新启动 AI Agent，或
  tryStart: 启动服务
  startSuccess: 服务已成功启动！
  startLocalService: 启动本地服务
  browserPrompt: 在浏览器弹出的提示框中，请选择允许打开 UniToken 应用。
  adminPermission: 首次启动时，UniToken 会请求管理员权限，请允许。
  troubleshootingPrefix: 若服务未能正常启动，请尝试重新启动相关 Agent 软件，或
  manualDownload: 手动下载
  downloadService: 下载 UniToken 服务
  selectOS: 请选择你的操作系统。
  afterDownload: 下载后，请运行它以启动服务。

en-US:
  connected: Local Service Connected
  disconnected: Service Disconnected
  connectionFailed: Service Connection Failed
  reconnecting: Reconnecting...
  restartAgent: |
    Local service may not have started properly. Please restart the related
    Agent software.
  troubleshooting: 'Please restart the AI agent, or'
  tryStart: Start Service
  startSuccess: Service started successfully!
  startLocalService: Start Local Service
  browserPrompt: |
    In the browser prompt that appears, please select to allow opening the
    UniToken application.
  adminPermission: |
    When starting for the first time, UniToken will request administrator
    permissions. Please allow.
  troubleshootingPrefix: |
    If the service fails to start properly, please try restarting the related
    Agent software, or
  manualDownload: download manually
  downloadService: Download UniToken Service
  selectOS: Please select your operating system.
  afterDownload: 'After downloading, please run it to start the service.'
</i18n>
