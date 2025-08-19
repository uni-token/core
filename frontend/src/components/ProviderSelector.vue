<script setup lang="ts">
import { Plus } from 'lucide-vue-next'
import { onMounted, ref, watchEffect } from 'vue'
import { useI18n } from 'vue-i18n'
import { Badge } from '@/components/ui/badge'
import { useProvidersStore } from '@/stores'
import AddProviderDialog from './AddProviderDialog.vue'
import Button from './ui/button/Button.vue'

const selectedProviderId = defineModel<string>()

const { t } = useI18n()
const providersStore = useProvidersStore()

const openAddProviderDialog = ref(false)

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
  <div class="space-y-1">
    <div class="text-sm font-medium text-gray-700">
      {{ t('providers.selectProvider') }}:
    </div>
    <div v-if="providersStore.loading" class="text-sm text-gray-500">
      {{ t('providers.loadingProviders') }}
    </div>
    <div v-else-if="providersStore.providers.length === 0" class="text-sm text-gray-500">
      <span class="text-red-500">
        {{ t('providers.noProvidersAvailable') }}
      </span>

      <Button variant="link" class="text-blue-600 hover:underline" @click="openAddProviderDialog = true">
        {{ t('providers.addNewProvider') }}
      </Button>
    </div>
    <div v-else class="flex flex-wrap gap-2">
      <Badge
        v-for="provider in providersStore.providers"
        :key="provider.id"
        :variant="selectedProviderId === provider.id ? 'default' : 'secondary'"
        class="cursor-pointer transition-colors hover:bg-gray-300"
        :class="{
          'bg-blue-500 hover:bg-blue-600 text-white': selectedProviderId === provider.id,
          'bg-gray-200 hover:bg-gray-300': selectedProviderId !== provider.id,
        }"
        @click="selectedProviderId = provider.id"
      >
        {{ provider.name }}
      </Badge>

      <!-- Add new provider button -->
      <Badge
        key="new"
        variant="secondary"
        class="text-xs cursor-pointer transition-colors hover:bg-gray-300 h-6"
        @click="openAddProviderDialog = true"
      >
        <Plus class="inline h-4 w-4" />
        {{ t('providers.addNewProvider') }}
      </Badge>
    </div>

    <AddProviderDialog v-model:open="openAddProviderDialog" @configured="selectedProviderId = $event" />
  </div>
</template>
