<script setup lang="ts">
import { Plus } from 'lucide-vue-next'
import { onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { Badge } from '@/components/ui/badge'
import { usePresetsStore } from '@/stores'

const props = defineProps<Props>()

const emit = defineEmits<Emits>()

const { t } = useI18n()

interface Props {
  modelValue?: string
}

interface Emits {
  (e: 'update:modelValue', value: string): void
}

const presetsStore = usePresetsStore()
const selectedPresetId = ref<string>('')

onMounted(() => presetsStore.loadPresets())

// Auto-select the first preset when presets are loaded
watch(() => presetsStore.presets, (presets) => {
  if (presets.length > 0 && !selectedPresetId.value) {
    selectedPresetId.value = presets[0].id
    emit('update:modelValue', presets[0].id)
  }
}, { immediate: true })

// Listen to external value changes
watch(() => props.modelValue, (value) => {
  if (value) {
    selectedPresetId.value = value
  }
}, { immediate: true })

function selectPreset(presetId: string) {
  selectedPresetId.value = presetId
  emit('update:modelValue', presetId)
}

// Load presets (if not already loaded)
if (presetsStore.presets.length === 0 && !presetsStore.loading) {
  presetsStore.loadPresets()
}
</script>

<template>
  <div class="space-y-1 max-h-20">
    <div class="text-sm font-medium text-gray-700">
      {{ t('presets.selectPreset') }}:
    </div>
    <div v-if="presetsStore.loading" class="text-sm text-gray-500">
      {{ t('presets.loadingPresets') }}
    </div>
    <div v-else-if="presetsStore.presets.length === 0" class="text-sm text-gray-500">
      {{ t('presets.noPresetsAvailable') }}
    </div>
    <div v-else class="flex flex-wrap gap-2">
      <Badge
        v-for="preset in presetsStore.presets"
        :key="preset.id"
        :variant="selectedPresetId === preset.id ? 'default' : 'secondary'"
        class="cursor-pointer transition-colors hover:bg-gray-300"
        :class="{
          'bg-blue-500 hover:bg-blue-600 text-white': selectedPresetId === preset.id,
          'bg-gray-200 hover:bg-gray-300': selectedPresetId !== preset.id,
        }"
        @click="selectPreset(preset.id)"
      >
        {{ preset.name }}
      </Badge>

      <!-- Add new preset button -->
      <RouterLink
        to="/models"
        target="_blank"
      >
        <Badge
          key="new"
          variant="secondary"
          class="text-xs cursor-pointer transition-colors hover:bg-gray-300 h-6"
        >
          <Plus class="inline h-4 w-4" />
          {{ t('presets.addNewPreset') }}
        </Badge>
      </RouterLink>
    </div>
  </div>
</template>
