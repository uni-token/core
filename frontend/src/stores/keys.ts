import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue-sonner'
import { useServiceStore } from './service'

export interface APIKey {
  id: string
  name: string
  type: string
  protocol: 'openai'
  baseUrl: string
  token: string
}

export const useKeysStore = defineStore('keys', () => {
  const { fetch } = useServiceStore()
  const { t } = useI18n()

  // State
  const keys = ref<APIKey[]>([])
  const loading = ref(true)
  const loadingError = ref<string | null>(null)

  // Actions
  async function loadKeys() {
    loadingError.value = null

    try {
      const response = await fetch('keys/list')
      if (response.ok) {
        const data = await response.json()
        keys.value = data.data
      }
      else {
        loadingError.value = `HTTP ${response.status}: ${response.statusText}`
      }
    }
    catch (err) {
      loadingError.value = err instanceof Error ? err.message : 'Unknown error'
    }
    finally {
      loading.value = false
    }
  }

  async function addKey(key: Omit<APIKey, 'id' | 'name'> & { name?: string }): Promise<APIKey> {
    try {
      const response = await fetch('keys/add', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(key),
      })

      if (response.ok) {
        await loadKeys()
        return (await response.json()).data
      }
      else {
        toast.error(t('stores.keys.addKeyFailed'))
        throw new Error(`HTTP ${response.status}: ${response.statusText}`)
      }
    }
    catch (err) {
      toast.error(err instanceof Error ? err.message : 'Unknown error')
      throw err
    }
  }

  async function updateKey(keyId: string, key: APIKey) {
    try {
      const response = await fetch(`keys/update/${encodeURIComponent(keyId)}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(key),
      })

      if (response.ok) {
        await loadKeys()
        return true
      }
      else {
        toast.error(t('stores.keys.updateKeyFailed'))
        return false
      }
    }
    catch (err) {
      toast.error(err instanceof Error ? err.message : 'Unknown error')
      return false
    }
  }

  async function deleteKey(keyId: string) {
    try {
      const response = await fetch(`keys/delete/${encodeURIComponent(keyId)}`, {
        method: 'DELETE',
      })

      if (response.ok) {
        await loadKeys()
        return true
      }
      else {
        toast.error(t('stores.keys.deleteKeyFailed'))
        return false
      }
    }
    catch (err) {
      toast.error(err instanceof Error ? err.message : 'Unknown error')
      return false
    }
  }

  return {
    // State
    keys,
    loading,
    loadingError,
    // Actions
    loadKeys,
    addKey,
    updateKey,
    deleteKey,
  }
})
