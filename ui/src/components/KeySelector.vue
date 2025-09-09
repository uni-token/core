<script setup lang="ts">
import type { Provider } from '@/lib/providers'
import { Plus } from 'lucide-vue-next'
import { onMounted, ref, shallowRef, watchEffect } from 'vue'
import { useI18n } from 'vue-i18n'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { useKeysStore, useProvidersStore } from '@/stores'
import AddKeyDialog from './AddKeyDialog.vue'
import ManualConfigCard from './ManualConfigCard.vue'
import ProviderConfigDialog from './ProviderConfigDialog.vue'
import Button from './ui/button/Button.vue'

const props = defineProps<{
  compact?: boolean
}>()
const selectedKeyId = defineModel<string>()

const providersStore = useProvidersStore()
const keysStore = useKeysStore()
const { t } = useI18n()

const openAddKeyDialog = ref(false)
const showProviderDialog = shallowRef<Provider | null>(null)

onMounted(() => keysStore.loadKeys())

watchEffect(() => {
  if (!selectedKeyId.value) {
    const firstKey = keysStore.keys[0]
    if (firstKey) {
      selectedKeyId.value = firstKey.id
    }
  }
})
</script>

<template>
  <div class="space-y-2">
    <div v-if="keysStore.loading" class="text-sm text-muted-foreground">
      {{ t('loadingKeys') }}
    </div>
    <div v-else-if="keysStore.keys.length === 0" class="text-sm text-muted-foreground">
      <template v-if="props.compact">
        <span class="text-red-500 dark:text-red-600">
          {{ t('noKeysAvailable') }}
        </span>

        <Button variant="link" class="text-blue-600 dark:text-blue-400 underline hover:opacity-80" @click="openAddKeyDialog = true">
          {{ t('addNewKey') }}
        </Button>
      </template>
      <template v-else>
        <div class="grid gap-4 grid-cols-2 mb-4">
          <Card v-for="provider in providersStore.list" :key="provider.id" class="relative gap-2">
            <CardHeader>
              <div class="flex items-center justify-between">
                <CardTitle class="text-lg">
                  {{ provider.name }}
                </CardTitle>
              </div>
            </CardHeader>

            <CardContent>
              <div class="text-sm text-muted-foreground">
                <p>
                  {{ t('description1') }}
                  <a :href="provider.homepage" target="_blank" class="text-blue-900 dark:text-blue-200 hover:underline">
                    {{ provider.name }}
                  </a>
                  {{ t('description2') }}
                </p>
              </div>
            </CardContent>

            <CardFooter>
              <Button class="w-full" @click="showProviderDialog = provider">
                {{ t('configure', [provider.name]) }}
              </Button>
            </CardFooter>
          </Card>

          <ManualConfigCard class="relative gap-2" />
        </div>
      </template>
    </div>
    <div v-else class="flex gap-4">
      <div :class="compact ? 'text-sm font-medium mb-1' : 'text-base font-bold mb-1'">
        {{ t('selectKey') }}
      </div>
      <div class="w-0 flex-grow flex flex-wrap gap-2 h-fit">
        <Badge
          v-for="key in keysStore.keys"
          :key="key.id"
          :variant="selectedKeyId === key.id ? 'default' : 'secondary'"
          class="cursor-pointer transition-colors"
          :class="{
            'bg-blue-500 hover:bg-blue-600 dark:bg-blue-800 dark:hover:bg-blue-600 text-white': selectedKeyId === key.id,
            'bg-muted hover:bg-gray-300 dark:hover:bg-gray-700': selectedKeyId !== key.id,
          }"
          @click="selectedKeyId = key.id"
        >
          {{ key.name }}
        </Badge>

        <!-- Add new key button -->
        <Badge
          variant="secondary"
          class="ml-auto text-xs cursor-pointer transition-colors bg-muted hover:bg-gray-300 dark:hover:bg-gray-700 h-6"
          @click="openAddKeyDialog = true"
        >
          <Plus class="inline h-4 w-4" />
          {{ t('addNewKey') }}
        </Badge>
      </div>
    </div>

    <ProviderConfigDialog
      v-if="showProviderDialog != null"
      :provider="showProviderDialog"
      @close="showProviderDialog = null"
    />
    <AddKeyDialog v-model:open="openAddKeyDialog" @configured="selectedKeyId = $event" />
  </div>
</template>

<i18n lang="yaml">
en-US:
  loadingKeys: Loading...
  noKeysAvailable: No providers configured
  addNewKey: Add Provider
  description1: Purchase and configure API through
  description2: ' '
  configure: Configure {0}
  selectKey: Provider

zh-CN:
  loadingKeys: 加载中
  noKeysAvailable: 未配置提供商
  addNewKey: 添加提供商
  description1: 通过
  description2: 购买和配置 API
  configure: 配置 {0}
  selectKey: 提供商
</i18n>
