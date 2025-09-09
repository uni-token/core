import { findServicePort } from '@uni-token/browser-sdk'
import { useIntervalFn } from '@vueuse/core'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { useAuthStore } from '@/stores/auth'

export const useServiceStore = defineStore('service', () => {
  const authStore = useAuthStore()

  const params = new URLSearchParams(window.location.search)
  const token = ref<string | null>(params.get('token') || null)

  const serverConnected = ref(true)
  const servicePort = ref<number | null>(null)
  const serviceHost = computed(() => {
    if (servicePort.value) {
      return `localhost:${servicePort.value}`
    }
    return null
  })
  const serviceUrl = computed(() => {
    if (serviceHost.value) {
      return `http://${serviceHost.value}/`
    }
    return null
  })
  let requireFindService = true

  const initialLoad = new Promise<void>((resolve) => {
    const givenServerPort = Number.parseInt(params.get('port') || '')
    if (givenServerPort > 0 && givenServerPort < 65535) {
      servicePort.value = givenServerPort
      requireFindService = false
      resolve()
    }

    useIntervalFn(
      async () => {
        if (requireFindService) {
          const port = await findServicePort()
          resolve()
          if (port) {
            serverConnected.value = true
            servicePort.value = port
            requireFindService = false
          }
          else {
            serverConnected.value = false
          }
        }
      },
      1000,
      {
        immediate: true,
        immediateCallback: true,
      },
    )
  })

  useIntervalFn(() => {
    requireFindService = true
  }, 5000)

  return {
    serverConnected,
    serviceHost,
    serviceUrl,
    servicePort,
    refreshService: () => {
      requireFindService = true
    },
    token,
    fetch: async (path: string, options?: RequestInit) => {
      await initialLoad
      if (!serviceUrl.value) {
        throw new Error('Service not available')
      }
      try {
        const resp = await fetch(
          `${serviceUrl.value}${path}`,
          token.value
            ? {
                ...options,
                headers: {
                  Authorization: `Bearer ${token.value}`,
                  ...options?.headers,
                },
              }
            : options,
        )

        if (resp.status === 401) {
          if (resp.headers.get('X-Uni-Token-Error')) {
            console.error('Unauthorized access, please check your token')
            authStore.currentUser = null
            token.value = null
          }
        }

        return resp
      }
      catch (error) {
        requireFindService = true
        console.error(`Error fetching ${path}:`, error)
        throw error
      }
    },
  }
})
