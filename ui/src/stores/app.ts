import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { toast } from 'vue-sonner'
import { useI18n } from '@/lib/locals'
import { useServiceStore } from './service'

export interface App {
  id: string
  name: string
  description?: string
  key: string
  granted: boolean
  createdAt: string
  lastActiveAt: string
}

export const useAppStore = defineStore('app', () => {
  const { api: fetch } = useServiceStore()
  const { t } = useI18n({
    'zh-CN': {
      appDeleted: '应用已删除',
      allAppsDeleted: '所有应用已删除',
      appAuthorized: '应用已授权',
      appAuthorizationRevoked: '应用授权已撤销',
    },
    'en-US': {
      appDeleted: 'Application deleted',
      allAppsDeleted: 'All applications deleted',
      appAuthorized: 'Application authorized',
      appAuthorizationRevoked: 'Application authorization revoked',
    },
  })

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

  const toggleAppAuthorization = async (id: string, granted: boolean, key?: string) => {
    try {
      const body: any = { id, granted }
      if (granted && key) {
        body.key = key
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

        toast.success(granted ? t('appAuthorized') : t('appAuthorizationRevoked'))
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

        toast.success(t('appDeleted'))
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
        toast.success(t('allAppsDeleted'))
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
