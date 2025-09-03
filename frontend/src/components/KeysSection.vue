<script setup lang="ts">
import type { Provider } from '@/lib/providers'
import type { APIKey } from '@/stores'
import { Edit } from 'lucide-vue-next'
import { onMounted, ref, shallowRef } from 'vue'
import { VueDraggable } from 'vue-draggable-plus'
import { useI18n } from 'vue-i18n'
import EditKeyDialog from '@/components/EditKeyDialog.vue'
import ManualConfigCard from '@/components/ManualConfigCard.vue'
import ProviderConfigDialog from '@/components/ProviderConfigDialog.vue'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { useKeysStore, useProvidersStore } from '@/stores'

const { t } = useI18n()
const keysStore = useKeysStore()
const providersStore = useProvidersStore()

const showProviderDialog = shallowRef<Provider | null>(null)
const showEditDialog = ref(false)
const editingKey = ref<APIKey | null>(null)

function editKey(key: APIKey) {
  if (providersStore.map[key.type]) {
    showProviderDialog.value = providersStore.map[key.type]
    return
  }
  editingKey.value = key
  showEditDialog.value = true
}

function formatUrl(url: string): string {
  // Remove protocol if present
  return url.replace(/^(https?:\/\/)?/, '')
}

onMounted(() => {
  keysStore.loadKeys()
})
</script>

<template>
  <div class="space-y-6 flex-grow flex flex-col min-h-100">
    <div class="flex items-center justify-between">
      <h2 class="text-2xl font-bold">
        {{ t('title') }}
      </h2>
    </div>

    <div v-if="keysStore.loading" class="space-y-3">
      <div v-for="i in 3" :key="i" class="space-y-3">
        <Skeleton class="h-4 w-full" />
        <Skeleton class="h-4 w-3/4" />
      </div>
    </div>

    <div v-else-if="keysStore.loadingError" class="rounded-lg border border-red-200 bg-red-50 p-4">
      <p class="text-red-800">
        {{ t('loadFailed') }}: {{ keysStore.loadingError }}
      </p>
    </div>

    <div v-else class="flex-grow flex flex-col">
      <!-- Static cards for configuration (not draggable) -->
      <div class="grid gap-4 grid-cols-2 mb-4">
        <!-- SiliconFlow configuration card (if not configured) -->
        <Card v-for="provider in providersStore.list" :key="provider.id" class="relative gap-2">
          <CardHeader>
            <div class="flex items-center justify-between">
              <CardTitle class="text-lg">
                {{ t('siliconFlow') }}
              </CardTitle>
              <div v-show="provider.user !== undefined" class="flex items-center gap-2">
                <span
                  class="inline-flex items-center rounded-full px-2 py-1 text-xs font-medium text-accent-foreground"
                  :class="{
                    'bg-emerald-200 dark:bg-green-500/60': !!provider.user,
                    'bg-accent': !provider.user,
                  }"
                >
                  {{ provider.user ? t('loggedIn') : t('loggedOut') }}
                </span>
              </div>
            </div>
          </CardHeader>

          <CardContent>
            <div class="text-sm text-muted-foreground">
              <p>
                {{ t('siliconFlowDescription1') }}
                <a href="https://www.siliconflow.cn" target="_blank" class="text-blue-900 dark:text-blue-200 hover:underline">
                  {{ t('siliconFlow') }}
                </a>
                {{ t('siliconFlowDescription2') }}
              </p>
            </div>
          </CardContent>

          <CardFooter>
            <Button class="w-full" @click="showProviderDialog = provider">
              {{ t('configureSiliconFlow') }}
            </Button>
          </CardFooter>
        </Card>

        <!-- Manual Configuration Provider Card -->
        <ManualConfigCard class="relative gap-2" />
      </div>
      <VueDraggable
        v-model="keysStore.keys"
        :group="{ name: 'keys', pull: false /* 'clone' */, put: false }"
        :sort="false"
        :disabled="true"
        :clone="key => ({ newKey: key.id })"
        class="max-h-fit grid gap-4 overflow-y-auto pb-1 flex-grow h-0"
        :style="{
          'grid-template-columns': keysStore.keys.length > 2 ? 'repeat(auto-fit, minmax(250px, 1fr))' : 'repeat(3, 1fr)',
        }"
      >
        <div v-for="key in keysStore.keys" :key="key.id" class="draggable-key">
          <Card class="draggable-key-card relative cursor-select hover:shadow-md transition-shadow gap-0! pt-4 pb-2">
            <CardHeader>
              <div class="flex items-center justify-between">
                <CardTitle class="text-lg">
                  {{ key.name }}
                </CardTitle>
                <div class="flex items-center gap-2">
                  <Button
                    variant="outline"
                    size="sm"
                    @click="editKey(key)"
                  >
                    <Edit class="h-4 w-4" />
                  </Button>
                </div>
              </div>
            </CardHeader>

            <CardContent class="flex mb-2 mt-2">
              <div class="font-mono text-xs truncate" :title="key.baseUrl">
                {{ formatUrl(key.baseUrl) }}
              </div>
            </CardContent>
          </Card>
          <Badge
            variant="default"
            class="draggable-key-badge text-sm flex items-center gap-1 cursor-move"
          >
            {{ key.name }}
            <button
              class="ml-1 text-xs hover:text-red-600"
            >
              ×
            </button>
          </Badge>
        </div>
      </VueDraggable>
    </div>

    <ProviderConfigDialog
      v-if="showProviderDialog != null"
      :provider="showProviderDialog"
      @close="showProviderDialog = null"
    />

    <EditKeyDialog
      v-model:open="showEditDialog"
      :api-key="editingKey"
    />
  </div>
</template>

<style>
.preset-drop-area .draggable-key-card {
  display: none;
}
.draggable-key-badge {
  display: none;
}
.preset-drop-area .draggable-key-badge {
  display: block;
}
</style>

<i18n lang="yaml">
zh-CN:
  title: 模型供应商
  loadFailed: 加载失败
  siliconFlow: 硅基流动
  siliconFlowDescription1: 通过
  siliconFlowDescription2: 购买和配置 API
  loggedIn: 已登录
  loggedOut: 未登录
  configureSiliconFlow: 配置 硅基流动

en-US:
  title: Model Providers
  loadFailed: Failed to load
  siliconFlow: SiliconFlow
  siliconFlowDescription1: Purchase and configure API through
  siliconFlowDescription2: ''
  loggedIn: Logged In
  loggedOut: Logged Out
  configureSiliconFlow: Configure SiliconFlow
</i18n>
