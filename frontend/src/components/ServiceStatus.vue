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
              {{ t('service.connectionFailed') }}
            </p>
            <p class="text-xs text-orange-700 dark:text-orange-300">
              {{ t('service.reconnecting') }}
            </p>
            <p class="text-xs text-orange-700 dark:text-orange-300">
              {{ t('service.restartAgent') }}<br>{{ t('common.or') }}
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
            <span class="text-sm select-none">{{ serviceStore.serverConnected ? t('service.connected') : t('service.disconnected') }}</span>
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
        {{ serviceStore.serverConnected ? t('service.connected') : t('service.disconnected') }}
      </span>
    </div>

    <!-- Permanent troubleshooting message when disconnected -->
    <div v-if="!serviceStore.serverConnected" class="text-center space-y-2">
      <p class="text-sm text-orange-800 dark:text-orange-200">
        {{ t('service.troubleshooting') }}
        <StartServiceButton />
      </p>
    </div>
  </div>
</template>
