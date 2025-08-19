<script setup lang="ts">
import { RefreshCw } from 'lucide-vue-next'
import { storeToRefs } from 'pinia'
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import AppDetailDialog from '@/components/AppDetailDialog.vue'
import ProviderSelector from '@/components/ProviderSelector.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { Switch } from '@/components/ui/switch'
import { useAppStore, usePresetsStore } from '@/stores'

const { t } = useI18n()
const appStore = useAppStore()
const { apps, loading, error } = storeToRefs(appStore)
const { loadApps, refreshApps, toggleAppAuthorization } = appStore
const presetsStore = usePresetsStore()

const showDetailDialog = ref(false)
const selectedApp = ref<any | null>(null)

function openAppDetail(app: any) {
  selectedApp.value = app
  showDetailDialog.value = true
}

onMounted(() => {
  loadApps()
  presetsStore.loadPresets()
})
</script>

<template>
  <div class="p-6">
    <div class="space-y-6">
      <div class="flex items-center justify-between">
        <h2 class="text-2xl font-bold">
          {{ t('appManagement.title') }}
        </h2>
        <Button :disabled="loading" @click="refreshApps">
          <RefreshCw class="mr-2 h-4 w-4" :class="{ 'animate-spin': loading }" />
          {{ t('appManagement.refresh') }}
        </Button>
      </div>

      <div v-if="loading" class="space-y-3">
        <div v-for="i in 3" :key="i" class="space-y-3">
          <Skeleton class="h-4 w-full" />
          <Skeleton class="h-4 w-3/4" />
        </div>
      </div>

      <div v-else-if="error" class="rounded-lg border border-red-200 bg-red-50 p-4">
        <p class="text-red-800">
          {{ t('appManagement.loadFailed') }}: {{ error }}
        </p>
      </div>

      <div v-else-if="apps.length === 0" class="rounded-lg border border-gray-200 bg-gray-50 p-8 text-center">
        <p class="text-gray-600">
          {{ t('appManagement.noApps') }}
        </p>
      </div>

      <div v-else class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <Card v-for="app in apps" :key="app.id" class="hover:shadow-md transition-shadow">
          <CardHeader>
            <div class="flex items-center justify-between">
              <CardTitle class="text-lg">
                {{ app.name }}
              </CardTitle>
              <div class="flex-grow" />
              <Button variant="outline" size="sm" :disabled="loading" @click="openAppDetail(app)">
                {{ t('appManagement.details') }}
              </Button>
            </div>
            <CardDescription>{{ app.description || t('appManagement.noDescription') }}</CardDescription>
          </CardHeader>

          <CardContent class="flex flex-col">
            <div>
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium">{{ t('appManagement.authorizationStatus') }}:</span>
                <Switch
                  v-model="app.granted"
                  :disabled="loading"
                  @update:model-value="(checked: boolean) => toggleAppAuthorization(app.id, checked)"
                />
              </div>
              <p class="text-xs text-gray-500 mt-1">
                {{ app.granted ? t('appManagement.hasAccess') : t('appManagement.noAccess') }}
              </p>
            </div>
            <div :class="{ 'opacity-40 pointer-events-none select-none': !app.granted }" class="mt-4">
              <ProviderSelector v-model="app.provider" @update:model-value="appStore.toggleAppAuthorization(app.id, true, app.provider)" />
            </div>
          </CardContent>
        </Card>
      </div>
    </div>

    <AppDetailDialog
      v-if="selectedApp" v-model:open="showDetailDialog" :app="selectedApp"
    />
  </div>
</template>
