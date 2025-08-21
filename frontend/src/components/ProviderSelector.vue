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
  <div class="space-y-2">
    <div class="text-sm font-medium">
      {{ t('providers.selectProvider') }}:
    </div>
    <div v-if="providersStore.loading" class="text-sm text-muted-foreground">
      {{ t('providers.loadingProviders') }}
    </div>
    <div v-else-if="providersStore.providers.length === 0" class="text-sm text-muted-foreground">
      <span class="text-red-500 dark:text-red-600">
        {{ t('providers.noProvidersAvailable') }}
      </span>

      <Button variant="link" class="text-blue-600 dark:text-blue-400 underline hover:opacity-80" @click="openAddProviderDialog = true">
        {{ t('providers.addNewProvider') }}
      </Button>
    </div>
    <div v-else class="flex flex-wrap gap-2">
      <Badge
        v-for="provider in providersStore.providers"
        :key="provider.id"
        :variant="selectedProviderId === provider.id ? 'default' : 'secondary'"
        class="cursor-pointer transition-colors"
        :class="{
          'bg-blue-500 hover:bg-blue-600 dark:bg-blue-800 dark:hover:bg-blue-600 text-white': selectedProviderId === provider.id,
          'bg-muted hover:bg-muted-foreground': selectedProviderId !== provider.id,
        }"
        @click="selectedProviderId = provider.id"
      >
        {{ provider.name }}
      </Badge>

      <!-- Add new provider button -->
      <Badge
        variant="secondary"
        class="text-xs cursor-pointer transition-colors bg-muted hover:bg-muted-foreground h-6"
        @click="openAddProviderDialog = true"
      >
        <Plus class="inline h-4 w-4" />
        {{ t('providers.addNewProvider') }}
      </Badge>
    </div>

    <AddProviderDialog v-model:open="openAddProviderDialog" @configured="selectedProviderId = $event" />
  </div>
</template>
