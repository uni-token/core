import { defineStore } from 'pinia'
import { ref } from 'vue'
import { toast } from 'vue-sonner'
import { t } from '@/locales'
import { useServiceStore } from './service'

export interface AppPreset {
  id: string
  name: string
  keys: string[]
  createdAt: string
  updatedAt: string
}

export const usePresetsStore = defineStore('presets', () => {
  const { fetch } = useServiceStore()

  // State
  const presets = ref<AppPreset[]>([])
  const loading = ref(true)
  const loadError = ref<string | null>(null)
  const editingPresetName = ref<{ [key: string]: string }>({})
  const activeEditId = ref<string | null>(null)
  const nameValidationError = ref<string | null>(null)

  // Actions
  async function loadPresets() {
    loadError.value = null

    try {
      const response = await fetch('presets/list')
      if (response.ok) {
        const data = (await response.json()).data
        const newPresets: AppPreset[] = []
        for (const newPreset of data) {
          const existingPreset = presets.value.findIndex(p => p.id === newPreset.id)
          if (existingPreset !== -1) {
            presets.value[existingPreset] = newPreset
          }
          else {
            newPresets.push(newPreset)
          }
        }
        presets.value = [...newPresets, ...presets.value].filter(p => data.some((d: AppPreset) => d.id === p.id))

        editingPresetName.value = {}
        presets.value.forEach((preset) => {
          editingPresetName.value[preset.id] = preset.name
        })
      }
      else {
        loadError.value = `HTTP ${response.status}: ${response.statusText}`
      }
    }
    catch (err) {
      console.error('Error loading presets:', err)
      loadError.value = err instanceof Error ? err.message : 'Unknown error'
    }
    finally {
      loading.value = false
    }
  }

  async function addPreset() {
    const existingNames = presets.value.map(p => p.name)
    let counter = 1
    let defaultName = `${t('presets.preset')} ${counter}`

    while (existingNames.includes(defaultName)) {
      counter++
      defaultName = `${t('presets.preset')} ${counter}`
    }

    try {
      const response = await fetch('presets/add', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          name: defaultName,
          keys: [],
        }),
      })

      if (response.ok) {
        await loadPresets()
        return true
      }
      else {
        toast.error(t('stores.presets.addPresetFailed'))
        return false
      }
    }
    catch (err) {
      toast.error(err instanceof Error ? err.message : 'Unknown error')
      return false
    }
  }

  async function deletePreset(preset: AppPreset) {
    try {
      const response = await fetch(`presets/delete/${encodeURIComponent(preset.id)}`, {
        method: 'DELETE',
      })

      if (response.ok) {
        await loadPresets()
        return true
      }
      else {
        toast.error(`${t('stores.presets.deletePresetFailed')}, ${response.status} ${response.statusText}`)
        return false
      }
    }
    catch (err) {
      toast.error(err instanceof Error ? err.message : 'Unknown error')
      return false
    }
  }

  async function updatePreset(presetId: string, updates: Partial<Pick<AppPreset, 'name' | 'keys'>>) {
    const preset = presets.value.find(p => p.id === presetId)
    if (!preset) {
      toast.error(t('stores.presets.presetNotFound'))
      return false
    }

    if (updates.name) {
      preset.name = updates.name
      editingPresetName.value[presetId] = updates.name
    }
    if (updates.keys) {
      preset.keys = updates.keys
    }

    try {
      const response = await fetch(`presets/update/${encodeURIComponent(presetId)}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          name: updates.name || preset.name,
          keys: updates.keys || preset.keys,
        }),
      })

      if (response.ok) {
        await loadPresets()
        return true
      }
      else {
        const errorData = await response.json().catch(() => ({}))
        toast.error(errorData.error || t('stores.presets.updatePresetFailed'))
        return false
      }
    }
    catch (err) {
      toast.error(err instanceof Error ? err.message : 'Unknown error')
      return false
    }
  }

  // UI state management
  function startEditPresetName(presetId: string) {
    activeEditId.value = presetId
    nameValidationError.value = null
  }

  function validatePresetName(newName: string, currentPresetId: string): boolean {
    nameValidationError.value = null

    if (!newName.trim()) {
      nameValidationError.value = t('stores.presets.presetNameEmpty')
      return false
    }

    const existingPreset = presets.value.find(p => p.name === newName && p.id !== currentPresetId)
    if (existingPreset) {
      nameValidationError.value = t('stores.presets.presetNameExists')
      return false
    }

    return true
  }

  async function savePresetName(preset: AppPreset) {
    const newName = editingPresetName.value[preset.id]
    if (!newName || newName === preset.name) {
      activeEditId.value = null
      return true
    }

    // Validate name before saving
    if (!validatePresetName(newName, preset.id)) {
      return false
    }

    const success = await updatePreset(preset.id, { name: newName })
    if (success) {
      activeEditId.value = null
      nameValidationError.value = null
    }
    else {
      editingPresetName.value[preset.id] = preset.name
    }
    return success
  }

  function cancelEditPresetName(preset: AppPreset) {
    editingPresetName.value[preset.id] = preset.name
    activeEditId.value = null
    nameValidationError.value = null
  }

  async function autoSavePresetName(preset: AppPreset) {
    const newName = editingPresetName.value[preset.id]
    if (!newName || newName === preset.name) {
      activeEditId.value = null
      nameValidationError.value = null
      return true
    }

    // Validate name before saving
    if (!validatePresetName(newName, preset.id)) {
      // If validation fails, revert to original name
      editingPresetName.value[preset.id] = preset.name
      activeEditId.value = null
      return false
    }

    const success = await updatePreset(preset.id, { name: newName })
    if (success) {
      activeEditId.value = null
      nameValidationError.value = null
    }
    else {
      editingPresetName.value[preset.id] = preset.name
      activeEditId.value = null
    }
    return success
  }

  // Helper function to add key to preset
  async function addKeyToPreset(presetId: string, keyId: string) {
    const preset = presets.value.find(p => p.id === presetId)
    if (!preset) {
      toast.error(t('stores.presets.presetNotFound'))
      return false
    }

    // Check if key already exists
    const existingIndex = preset.keys.findIndex(p => p === keyId)
    let newKeys: string[]

    if (existingIndex !== -1) {
      // Key already exists, move it to the end
      newKeys = [...preset.keys]
      newKeys.splice(existingIndex, 1) // Remove from current position
      newKeys.push(keyId) // Add to the end
    }
    else {
      // Key doesn't exist, add it to the end
      newKeys = [...preset.keys, keyId]
    }

    return await updatePreset(presetId, { keys: newKeys })
  }

  // Helper function to remove key from preset
  async function removeKeyFromPreset(presetId: string, keyId: string) {
    const preset = presets.value.find(p => p.id === presetId)
    if (!preset) {
      toast.error(t('stores.presets.presetNotFound'))
      return false
    }

    const newKeys = preset.keys.filter(p => p !== keyId)
    return await updatePreset(presetId, { keys: newKeys })
  }

  return {
    // State
    presets,
    loading,
    loadError,
    editingPresetName,
    activeEditId,
    nameValidationError,
    // Actions
    getPresetName: (id: string) => editingPresetName.value[id] || '',
    loadPresets,
    addPreset,
    deletePreset,
    updatePreset,
    startEditPresetName,
    validatePresetName,
    savePresetName,
    autoSavePresetName,
    cancelEditPresetName,
    addKeyToPreset,
    removeKeyFromPreset,
  }
})
