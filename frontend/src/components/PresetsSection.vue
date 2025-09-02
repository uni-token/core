<script setup lang="ts">
import type { AppPreset } from '@/stores'
import { onClickOutside } from '@vueuse/core'
import { Edit, Plus, Trash2 } from 'lucide-vue-next'
import { computed, onMounted, reactive, ref } from 'vue'
import { VueDraggable } from 'vue-draggable-plus'
import { useI18n } from 'vue-i18n'
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger } from '@/components/ui/alert-dialog'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Skeleton } from '@/components/ui/skeleton'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { useKeysStore, usePresetsStore } from '@/stores'

const { t } = useI18n()
const presetsStore = usePresetsStore()
const keysStore = useKeysStore()
const showKeySelector = ref<string | null>(null)

const editContainerRefs = reactive<{ [key: string]: HTMLElement | null }>({})

onClickOutside(computed(() => presetsStore.activeEditId ? editContainerRefs[presetsStore.activeEditId] : null), () => {
  const preset = presetsStore.presets.find(p => p.id === presetsStore.activeEditId)
  if (preset) {
    presetsStore.autoSavePresetName(preset)
  }
})

// Add key to preset
async function addKeyToPreset(presetId: string, keyId: string) {
  await presetsStore.addKeyToPreset(presetId, keyId)
  showKeySelector.value = null
}

// Remove key from preset
async function removeKeyFromPreset(presetId: string, keyId: string) {
  await presetsStore.removeKeyFromPreset(presetId, keyId)
}

// Get unique keys for display (removes duplicates while preserving latest position)
function getUniqueKeys(keys: string[]) {
  if (!keys || keys.length === 0)
    return []

  const lastPositions = new Map<string, { key: string, position: number }>()

  for (let i = 0; i < keys.length; i++) {
    const key = keys[i]
    lastPositions.set(key, { key, position: i })
  }

  return Array.from(lastPositions.values())
    .sort((a, b) => a.position - b.position)
    .map(item => item.key)
}

// Get available keys that are not already in the preset
function getAvailableKeys(preset: AppPreset) {
  const available = []
  for (const key of keysStore.keys) {
    if (!preset.keys.includes(key.id)) {
      available.push(key)
    }
  }
  return available
}

// Handle drag and drop - update preset keys
async function updatePresetKeys(presetId: string, newKeys: (string | { newKey: string })[]) {
  const newKeyIds = new Set()
  for (const item of newKeys) {
    if (typeof item === 'object' && item.newKey) {
      newKeyIds.add(item.newKey)
    }
  }

  const uniqueKeys: string[] = []
  for (const item of newKeys) {
    if (typeof item === 'object' && item.newKey) {
      uniqueKeys.push(item.newKey)
    }
    if (typeof item === 'string' && !newKeyIds.has(item)) {
      uniqueKeys.push(item)
    }
  }

  await presetsStore.updatePreset(presetId, { keys: uniqueKeys })
}

onMounted(() => {
  presetsStore.loadPresets()
})
</script>

<template>
  <div class="flex flex-col flex-grow space-y-6 min-h-80">
    <div class="flex items-center justify-between">
      <h2 class="text-2xl font-bold">
        {{ t('presets.title') }}
      </h2>
      <Button @click="presetsStore.addPreset">
        <Plus class="mr-2 h-4 w-4" />
        {{ t('presets.addPreset') }}
      </Button>
    </div>

    <div v-if="presetsStore.loading" class="space-y-3">
      <div v-for="i in 3" :key="i" class="space-y-3">
        <Skeleton class="h-4 w-full" />
        <Skeleton class="h-4 w-3/4" />
      </div>
    </div>

    <div v-else-if="presetsStore.loadError && presetsStore.presets.length === 0" class="rounded-lg border border-red-200 bg-red-50 p-4">
      <p class="text-red-800">
        {{ t('presets.loadFailed') }}: {{ presetsStore.loadError }}
      </p>
    </div>

    <Card v-else class="h-0 flex flex-col flex-grow px-4 py-4">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead class="min-w-28">
              {{ t('presets.presetName') }}
            </TableHead>
            <TableHead class="w-full">
              {{ t('presets.keys') }}
            </TableHead>
            <TableHead>
              {{ t('presets.actions') }}
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="preset in presetsStore.presets" :key="preset.id">
            <TableCell>
              <div v-if="presetsStore.activeEditId === preset.id" class="flex items-center gap-2">
                <div :ref="(el: any) => (editContainerRefs[preset.id] = el)" class="flex flex-col gap-1">
                  <Input
                    v-model="presetsStore.editingPresetName[preset.id]"
                    class="h-8 w-24"
                    :class="{ 'border-red-500': presetsStore.nameValidationError }"
                    @keydown.enter="presetsStore.autoSavePresetName(preset)"
                    @keydown.escape="presetsStore.cancelEditPresetName(preset)"
                    @input="presetsStore.validatePresetName(presetsStore.editingPresetName[preset.id], preset.id)"
                  />
                  <span v-if="presetsStore.nameValidationError" class="text-sm text-red-500">
                    {{ presetsStore.nameValidationError }}
                  </span>
                </div>
              </div>
              <div v-else class="flex items-center gap-2">
                <span class="font-medium ml-1">{{ preset.name }}</span>
                <div class="flex-grow" />
                <Button
                  size="sm"
                  variant="ghost"
                  class="h-6 w-6 p-0"
                  @click="presetsStore.startEditPresetName(preset.id)"
                >
                  <Edit class="h-3 w-3 opacity-80" />
                </Button>
              </div>
            </TableCell>
            <TableCell>
              <AlertDialog>
                <AlertDialogTrigger as-child>
                  <div
                    class="preset-drop-area group relative flex items-center pl-2 pr-1 py-1 border-2 border-dashed border-gray-200 rounded-md hover:border-gray-300 hover:bg-gray-100 transition-colors"
                  >
                    <VueDraggable
                      :model-value="preset.keys"
                      :group="{ name: 'keys', pull: 'clone' }"
                      class="flex flex-wrap flex-grow gap-1 items-center min-h-8"
                      :class="
                        {
                          'group-hover:blur-[2px]': getAvailableKeys(preset).length > 0,
                        }"
                      :clone="key => ({ newKey: key })"
                      @update:model-value="(newKeys: (string | {newKey:string})[]) => updatePresetKeys(preset.id, newKeys)"
                      @click="(ev: MouseEvent) => getAvailableKeys(preset).length === 0 && ev.stopImmediatePropagation()"
                    >
                      <Badge
                        v-for="keyId in getUniqueKeys(preset.keys)"
                        :key="keyId"
                        variant="secondary"
                        class="text-sm flex items-center gap-1 cursor-move bg-gray-200 hover:bg-gray-300"
                        @click="(ev: MouseEvent) => ev.stopPropagation()"
                      >
                        {{ keysStore.keys.find(p => p.id === keyId)?.name || t('presets.unknownKey') }}
                        <button
                          class="ml-1 text-xs hover:text-red-600"
                          :title="t('presets.remove')"
                          @click="removeKeyFromPreset(preset.id, keyId)"
                        >
                          ×
                        </button>
                      </Badge>
                    </VueDraggable>

                    <div v-if="getAvailableKeys(preset).length > 0" class="absolute inset-0 flex pointer-events-none backdrop-blur-[0px]">
                      <div
                        class="w-10 flex flex-grow flex-wrap gap-x-1 gap-y-0 items-center justify-center text-sm text-center"
                        :class="getUniqueKeys(preset.keys).length === 0 ? 'opacity-90' : 'opacity-0 group-hover:opacity-80 bg-transparent group-hover:bg-gray-200 duration-400 transition-all'"
                        style="text-shadow: 0 0 8px rgba(255, 255, 255, 0.8), 0 0 20px rgba(255, 255, 255, 0.6);"
                      >
                        <Plus class="h-4 w-4" />
                        <span class="hidden sm:inline">
                          {{ t('presets.dragKeysHere') }}
                        </span>
                      </div>
                    </div>
                  </div>
                </AlertDialogTrigger>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>{{ t('presets.addKeyToPreset') }}</AlertDialogTitle>
                    <AlertDialogDescription>
                      {{ t('presets.selectKeyDescription') }} "{{ preset.name }}" {{ t('presets.selectKeyDescription2') }}
                    </AlertDialogDescription>
                  </AlertDialogHeader>
                  <div class="py-4">
                    <div v-if="getAvailableKeys(preset).length === 0" class="text-sm text-gray-500">
                      {{ t('presets.noAvailableKeys') }}
                    </div>
                    <div v-else class="space-y-2 flex flex-col max-h-[50vh] overflow-y-auto">
                      <Button
                        v-for="key in getAvailableKeys(preset)"
                        :key="key.id"
                        variant="outline"
                        class="flex items-center p-4 border rounded cursor-pointer hover:bg-gray-50"
                        @click="addKeyToPreset(preset.id, key.id)"
                      >
                        <div class="font-medium">
                          {{ key.name }}
                        </div>
                        <div class="text-sm text-gray-500">
                          {{ key.type }}
                        </div>
                        <div class="flex-grow" />
                      </Button>
                    </div>
                  </div>
                  <AlertDialogFooter>
                    <AlertDialogCancel>{{ t('presets.cancel') }}</AlertDialogCancel>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            </TableCell>
            <TableCell class="text-right">
              <div class="flex items-center justify-end gap-2">
                <AlertDialog>
                  <AlertDialogTrigger as-child>
                    <Button
                      variant="outline"
                      size="sm"
                      class="text-red-600 hover:text-red-700"
                      :disabled="preset.name === 'default'"
                    >
                      <Trash2 class="h-4 w-4" />
                    </Button>
                  </AlertDialogTrigger>
                  <AlertDialogContent>
                    <AlertDialogHeader>
                      <AlertDialogTitle>{{ t('presets.confirmDelete') }}</AlertDialogTitle>
                      <AlertDialogDescription>
                        {{ t('presets.confirmDeleteDescription') }} "{{ preset.name }}"?
                      </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                      <AlertDialogCancel>{{ t('presets.cancel') }}</AlertDialogCancel>
                      <AlertDialogAction @click="presetsStore.deletePreset(preset)">
                        {{ t('presets.delete') }}
                      </AlertDialogAction>
                    </AlertDialogFooter>
                  </AlertDialogContent>
                </AlertDialog>
              </div>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </Card>
  </div>
</template>
