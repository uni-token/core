<script setup lang="ts">
import type { ManualConfig } from '@/components/ManualConfigDialog.vue'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue-sonner'
import ManualConfigDialog from '@/components/ManualConfigDialog.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { useProvidersStore } from '@/stores'

const emit = defineEmits<{
  configured: [provider: string]
}>()

const { t } = useI18n()
const providersStore = useProvidersStore()

const showEditDialog = ref(false)

async function handleManualConfigSave(config: ManualConfig) {
  const provider = await providersStore.addProvider({
    name: config.name,
    type: 'manual',
    protocol: config.protocol,
    baseUrl: config.baseUrl,
    token: config.token,
  })

  emit('configured', provider.id)
  showEditDialog.value = false
  toast.success(t('providers.manualConfigSuccess'))
}
</script>

<template>
  <Card>
    <CardHeader>
      <div class="flex items-center justify-between">
        <CardTitle class="text-lg">
          {{ t('providers.manualConfig') }}
        </CardTitle>
      </div>
    </CardHeader>

    <CardContent class="flex-grow">
      <div class="text-sm text-gray-600">
        <p>{{ t('providers.manualConfigDescription') }}</p>
      </div>
    </CardContent>

    <CardFooter>
      <Button class="w-full" @click="showEditDialog = true">
        {{ t('providers.addManualConfig') }}
      </Button>
    </CardFooter>

    <ManualConfigDialog
      v-model:open="showEditDialog"
      @save="handleManualConfigSave"
    />
  </Card>
</template>
