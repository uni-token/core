<script setup lang="ts">
import type { LLMProvider } from '@/stores'
import { Edit } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import { VueDraggable } from 'vue-draggable-plus'
import { useI18n } from 'vue-i18n'
import EditProviderDialog from '@/components/EditProviderDialog.vue'
import ManualConfigCard from '@/components/ManualConfigCard.vue'
import SiliconFlowConfigDialog from '@/components/SiliconFlowConfigDialog.vue'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { useProvidersStore, useSiliconFlowStore } from '@/stores'

const { t } = useI18n()
const providersStore = useProvidersStore()
const siliconFlowStore = useSiliconFlowStore()
const showSiliconFlowConfig = ref(false)
const showEditDialog = ref(false)
const editingProvider = ref<LLMProvider | null>(null)

function editProvider(provider: LLMProvider) {
  editingProvider.value = provider
  showEditDialog.value = true
}

function formatUrl(url: string): string {
  // Remove protocol if present
  return url.replace(/^(https?:\/\/)?/, '')
}

onMounted(() => {
  providersStore.loadProviders()
  siliconFlowStore.checkLoginStatus()
})
</script>

<template>
  <div class="space-y-6 flex-grow flex flex-col min-h-100">
    <div class="flex items-center justify-between">
      <h2 class="text-2xl font-bold">
        {{ t('providers.title') }}
      </h2>
    </div>

    <div v-if="providersStore.loading" class="space-y-3">
      <div v-for="i in 3" :key="i" class="space-y-3">
        <Skeleton class="h-4 w-full" />
        <Skeleton class="h-4 w-3/4" />
      </div>
    </div>

    <div v-else-if="providersStore.loadingError" class="rounded-lg border border-red-200 bg-red-50 p-4">
      <p class="text-red-800">
        {{ t('providers.loadFailed') }}: {{ providersStore.loadingError }}
      </p>
    </div>

    <div v-else class="flex-grow flex flex-col">
      <!-- Static cards for configuration (not draggable) -->
      <div class="grid gap-4 grid-cols-2 mb-4">
        <!-- SiliconFlow configuration card (if not configured) -->
        <Card class="relative gap-2">
          <CardHeader>
            <div class="flex items-center justify-between">
              <CardTitle class="text-lg">
                {{ t('providers.siliconFlow') }}
              </CardTitle>
              <div v-show="!siliconFlowStore.isLoading" class="flex items-center gap-2">
                <span
                  class="inline-flex items-center rounded-full px-2 py-1 text-xs font-medium text-accent-foreground"
                  :class="{
                    'bg-emerald-200 dark:bg-green-500/60': siliconFlowStore.isLoggedIn,
                    'bg-accent': !siliconFlowStore.isLoggedIn,
                  }"
                >
                  {{ siliconFlowStore.isLoggedIn ? t('providers.loggedIn') : t('providers.loggedOut') }}
                </span>
              </div>
            </div>
          </CardHeader>

          <CardContent>
            <div class="text-sm text-muted-foreground">
              <p>
                {{ t('providers.siliconFlowDescription1') }}
                <a href="https://www.siliconflow.cn" target="_blank" class="text-blue-900 dark:text-blue-200 hover:underline">
                  {{ t('providers.siliconFlow') }}
                </a>
                {{ t('providers.siliconFlowDescription2') }}
              </p>
            </div>
          </CardContent>

          <CardFooter>
            <Button class="w-full" @click="showSiliconFlowConfig = true">
              {{ t('providers.configureSiliconFlow') }}
            </Button>
          </CardFooter>
        </Card>

        <!-- Manual Configuration Provider Card -->
        <ManualConfigCard class="relative gap-2" />
      </div>
      <VueDraggable
        v-model="providersStore.providers"
        :group="{ name: 'providers', pull: false /* 'clone' */, put: false }"
        :sort="false"
        :disabled="true"
        :clone="provider => ({ newProvider: provider.id })"
        class="max-h-fit grid gap-4 overflow-y-auto pb-1 flex-grow h-0"
        :style="{
          'grid-template-columns': providersStore.providers.length > 2 ? 'repeat(auto-fit, minmax(250px, 1fr))' : 'repeat(3, 1fr)',
        }"
      >
        <div v-for="provider in providersStore.providers" :key="provider.id" class="draggable-provider">
          <Card class="draggable-provider-card relative cursor-select hover:shadow-md transition-shadow gap-0! pt-4 pb-2">
            <CardHeader>
              <div class="flex items-center justify-between">
                <CardTitle class="text-lg">
                  {{ provider.name }}
                </CardTitle>
                <div class="flex items-center gap-2">
                  <Button
                    variant="outline"
                    size="sm"
                    @click="editProvider(provider)"
                  >
                    <Edit class="h-4 w-4" />
                  </Button>
                </div>
              </div>
            </CardHeader>

            <CardContent class="flex mb-2 mt-2">
              <div class="font-mono text-xs truncate" :title="provider.baseUrl">
                {{ formatUrl(provider.baseUrl) }}
              </div>
            </CardContent>
          </Card>
          <Badge
            variant="default"
            class="draggable-provider-badge text-sm flex items-center gap-1 cursor-move"
          >
            {{ provider.name }}
            <button
              class="ml-1 text-xs hover:text-red-600"
            >
              ×
            </button>
          </Badge>
        </div>
      </VueDraggable>
    </div>

    <!-- Dialog Components -->
    <SiliconFlowConfigDialog
      v-model:open="showSiliconFlowConfig"
    />

    <EditProviderDialog
      v-model:open="showEditDialog"
      :provider="editingProvider"
    />
  </div>
</template>

<style>
.preset-drop-area .draggable-provider-card {
  display: none;
}
.draggable-provider-badge {
  display: none;
}
.preset-drop-area .draggable-provider-badge {
  display: block;
}
</style>
