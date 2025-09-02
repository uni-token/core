<script setup lang="ts">
import { Plus } from 'lucide-vue-next'
import { onMounted, ref, watchEffect } from 'vue'
import { useI18n } from 'vue-i18n'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { useKeysStore } from '@/stores'
import AddKeyDialog from './AddKeyDialog.vue'
import ManualConfigCard from './ManualConfigCard.vue'
import SiliconFlowConfigDialog from './SiliconFlowConfigDialog.vue'
import Button from './ui/button/Button.vue'

const props = defineProps<{
  compact?: boolean
}>()
const selectedKeyId = defineModel<string>()
const { t } = useI18n()
const keysStore = useKeysStore()

const openAddKeyDialog = ref(false)
const showSiliconFlowConfig = ref(false)

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
          <!-- SiliconFlow configuration card (if not configured) -->
          <Card class="relative gap-2">
            <CardHeader>
              <div class="flex items-center justify-between">
                <CardTitle class="text-lg">
                  {{ t('siliconFlow') }}
                </CardTitle>
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
              <Button class="w-full" @click="showSiliconFlowConfig = true">
                {{ t('configureSiliconFlow') }}
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

    <SiliconFlowConfigDialog
      v-model:open="showSiliconFlowConfig"
      @configured="selectedKeyId = $event"
    />
    <AddKeyDialog v-model:open="openAddKeyDialog" @configured="selectedKeyId = $event" />
  </div>
</template>

<i18n lang="yaml">
en-US:
  loadingKeys: Loading...
  noKeysAvailable: No providers configured
  addNewKey: Add Provider
  siliconFlow: SiliconFlow
  siliconFlowDescription1: Purchase and configure API through
  siliconFlowDescription2: ' '
  configureSiliconFlow: Configure SiliconFlow
  selectKey: Provider

zh-CN:
  loadingKeys: 加载中
  noKeysAvailable: 未配置提供商
  addNewKey: 添加提供商
  siliconFlow: 硅基流动
  siliconFlowDescription1: 通过
  siliconFlowDescription2: 购买和配置 API
  configureSiliconFlow: 配置 硅基流动
  selectKey: 提供商
</i18n>
