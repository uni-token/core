import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue-sonner'
import { useService } from '@/composables/service'

export interface App {
  id: string
  name: string
  description?: string
  provider: string
  granted: boolean
  createdAt: string
  lastActiveAt: string
}

export const useAppStore = defineStore('app', () => {
  const { fetch } = useService()
  const { t } = useI18n()

  // State
  const apps = ref<App[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Getters
  // const activeApps = computed(() => apps.value.filter(app => app.status === 'active'))
  const grantedApps = computed(() => apps.value.filter(app => app.granted))
  const appCount = computed(() => apps.value.length)

  const getAppById = computed(() => (id: string) => apps.value.find(app => app.id === id))

  // Actions
  const loadApps = async () => {
    loading.value = true
    error.value = null

    try {
      const response = await fetch('app/list')
      if (response.ok) {
        const data = await response.json()
        apps.value = data.data || data || []
      }
      else {
        error.value = `HTTP ${response.status}: ${response.statusText}`
        toast.error(error.value)
      }
    }
    catch (err) {
      error.value = err instanceof Error ? err.message : 'Unknown error'
      toast.error(error.value)
    }
    finally {
      loading.value = false
    }
  }

  const refreshApps = async () => {
    await loadApps()
  }

  const toggleAppAuthorization = async (id: string, granted: boolean, provider?: string) => {
    try {
      const body: any = { id, granted }
      if (granted && provider) {
        body.provider = provider
      }

      const response = await fetch('app/toggle', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(body),
      })

      if (response.ok) {
        const appIndex = apps.value.findIndex(app => app.id === id)
        if (appIndex !== -1) {
          apps.value[appIndex].granted = granted
        }

        toast.success(granted ? t('stores.app.appAuthorized') : t('stores.app.appAuthorizationRevoked'))
      }
      else {
        const errorData = await response.json()
        toast.error(errorData.error || `Operation failed: HTTP ${response.status}`)
        throw new Error(errorData.error || `Operation failed: HTTP ${response.status}`)
      }
    }
    catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Operation failed'
      toast.error(errorMessage)
      throw err
    }
  }

  const deleteApp = async (id: string) => {
    try {
      const response = await fetch(`app/delete/${id}`, {
        method: 'DELETE',
      })

      if (response.ok) {
        const appIndex = apps.value.findIndex(app => app.id === id)
        if (appIndex !== -1) {
          apps.value.splice(appIndex, 1)
        }

        toast.success(t('stores.app.appDeleted'))
      }
      else {
        const errorData = await response.json()
        toast.error(errorData.error || `Delete failed: HTTP ${response.status}`)
        throw new Error(errorData.error || `Delete failed: HTTP ${response.status}`)
      }
    }
    catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Delete failed'
      toast.error(errorMessage)
      throw err
    }
  }

  const deleteAllApps = async () => {
    try {
      const response = await fetch('app/clear', {
        method: 'DELETE',
      })

      if (response.ok) {
        apps.value = []
        toast.success(t('stores.app.allAppsDeleted'))
      }
      else {
        const errorData = await response.json()
        toast.error(errorData.error || `Delete failed: HTTP ${response.status}`)
        throw new Error(errorData.error || `Delete failed: HTTP ${response.status}`)
      }
    }
    catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Delete failed'
      toast.error(errorMessage)
      throw err
    }
  }

  const clearError = () => {
    error.value = null
  }

  const addApp = (app: App) => {
    apps.value.push(app)
  }

  const updateApp = (id: string, updates: Partial<App>) => {
    const appIndex = apps.value.findIndex(app => app.id === id)
    if (appIndex !== -1) {
      apps.value[appIndex] = { ...apps.value[appIndex], ...updates }
    }
  }

  return {
    // State
    apps,
    loading,
    error,

    // Getters
    // activeApps,
    grantedApps,
    appCount,
    getAppById,

    // Actions
    loadApps,
    refreshApps,
    toggleAppAuthorization,
    deleteApp,
    deleteAllApps,
    clearError,
    addApp,
    updateApp,
  }
})
