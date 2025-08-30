<script setup lang="ts">
import { Plus } from 'lucide-vue-next'
import { onMounted, ref, watchEffect } from 'vue'
import { useI18n } from 'vue-i18n'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { useProvidersStore } from '@/stores'
import AddProviderDialog from './AddProviderDialog.vue'
import ManualConfigCard from './ManualConfigCard.vue'
import SiliconFlowConfigDialog from './SiliconFlowConfigDialog.vue'
import Button from './ui/button/Button.vue'

const props = defineProps<{
  compact?: boolean
}>()
const selectedProviderId = defineModel<string>()
const { t } = useI18n()
const providersStore = useProvidersStore()

const openAddProviderDialog = ref(false)
const showSiliconFlowConfig = ref(false)

onMounted(() => providersStore.loadProviders())

watchEffect(() => {
  if (!selectedProviderId.value) {
    const firstProvider = providersStore.providers[0]
    if (firstProvider) {
      selectedProviderId.value = firstProvider.id
    }
  }
})
</script>

<template>
  <div class="space-y-2">
    <div v-if="providersStore.loading" class="text-sm text-muted-foreground">
      {{ t('providers.loadingProviders') }}
    </div>
    <div v-else-if="providersStore.providers.length === 0" class="text-sm text-muted-foreground">
      <template v-if="props.compact">
        <span class="text-red-500 dark:text-red-600">
          {{ t('providers.noProvidersAvailable') }}
        </span>

        <Button variant="link" class="text-blue-600 dark:text-blue-400 underline hover:opacity-80" @click="openAddProviderDialog = true">
          {{ t('providers.addNewProvider') }}
        </Button>
      </template>
      <template v-else>
        <div class="grid gap-4 grid-cols-2 mb-4">
          <!-- SiliconFlow configuration card (if not configured) -->
          <Card class="relative gap-2">
            <CardHeader>
              <div class="flex items-center justify-between">
                <CardTitle class="text-lg">
                  {{ t('providers.siliconFlow') }}
                </CardTitle>
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

          <ManualConfigCard class="relative gap-2" />
        </div>
      </template>
    </div>
    <div v-else class="flex gap-4">
      <div :class="compact ? 'text-sm font-medium mb-1' : 'text-base font-bold mb-1'">
        {{ t('providers.selectProvider') }}
      </div>
      <div class="w-0 flex-grow flex flex-wrap gap-2 h-fit">
        <Badge
          v-for="provider in providersStore.providers"
          :key="provider.id"
          :variant="selectedProviderId === provider.id ? 'default' : 'secondary'"
          class="cursor-pointer transition-colors"
          :class="{
            'bg-blue-500 hover:bg-blue-600 dark:bg-blue-800 dark:hover:bg-blue-600 text-white': selectedProviderId === provider.id,
            'bg-muted hover:bg-gray-300 dark:hover:bg-gray-700': selectedProviderId !== provider.id,
          }"
          @click="selectedProviderId = provider.id"
        >
          {{ provider.name }}
        </Badge>

        <!-- Add new provider button -->
        <Badge
          variant="secondary"
          class="ml-auto text-xs cursor-pointer transition-colors bg-muted hover:bg-gray-300 dark:hover:bg-gray-700 h-6"
          @click="openAddProviderDialog = true"
        >
          <Plus class="inline h-4 w-4" />
          {{ t('providers.addNewProvider') }}
        </Badge>
      </div>
    </div>

    <SiliconFlowConfigDialog
      v-model:open="showSiliconFlowConfig"
      @configured="selectedProviderId = $event"
    />
    <AddProviderDialog v-model:open="openAddProviderDialog" @configured="selectedProviderId = $event" />
  </div>
</template>
