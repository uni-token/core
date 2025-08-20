import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue-sonner'
import { useServiceStore } from './service'

export interface LLMProvider {
  id: string
  name: string
  type: string
  protocol: 'openai'
  baseUrl: string
  token: string
}

export const useProvidersStore = defineStore('providers', () => {
  const { fetch } = useServiceStore()
  const { t } = useI18n()

  // State
  const providers = ref<LLMProvider[]>([])
  const loading = ref(true)
  const loadingError = ref<string | null>(null)

  // Actions
  async function loadProviders() {
    loadingError.value = null

    try {
      const response = await fetch('providers/list')
      if (response.ok) {
        const data = await response.json()
        providers.value = data.data
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

  async function addProvider(provider: Omit<LLMProvider, 'id' | 'name'> & { name?: string }): Promise<LLMProvider> {
    try {
      const response = await fetch('providers/add', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(provider),
      })

      if (response.ok) {
        await loadProviders()
        return (await response.json()).data
      }
      else {
        toast.error(t('stores.providers.addProviderFailed'))
        throw new Error(`HTTP ${response.status}: ${response.statusText}`)
      }
    }
    catch (err) {
      toast.error(err instanceof Error ? err.message : 'Unknown error')
      throw err
    }
  }

  async function updateProvider(providerId: string, provider: LLMProvider) {
    try {
      const response = await fetch(`providers/update/${encodeURIComponent(providerId)}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(provider),
      })

      if (response.ok) {
        await loadProviders()
        return true
      }
      else {
        toast.error(t('stores.providers.updateProviderFailed'))
        return false
      }
    }
    catch (err) {
      toast.error(err instanceof Error ? err.message : 'Unknown error')
      return false
    }
  }

  async function deleteProvider(providerId: string) {
    try {
      const response = await fetch(`providers/delete/${encodeURIComponent(providerId)}`, {
        method: 'DELETE',
      })

      if (response.ok) {
        await loadProviders()
        return true
      }
      else {
        toast.error(t('stores.providers.deleteProviderFailed'))
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
    providers,
    loading,
    loadingError,
    // Actions
    loadProviders,
    addProvider,
    updateProvider,
    deleteProvider,
  }
})
