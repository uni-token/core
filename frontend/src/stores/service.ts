import { useIntervalFn } from '@vueuse/core'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const useServiceStore = defineStore('service', () => {
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

  const params = new URLSearchParams(window.location.search)
  const token = ref<string | null>(params.get('token') || null)

  const initialLoad = new Promise<void>((resolve) => {
    useIntervalFn(
      async () => {
        if (requireFindService) {
          const port = await findPort()
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
  }, 3000)

  async function findPort() {
    for (let port = 18000; port < 18010; port++) {
      if (await checkPort(port)) {
        return port
      }
    }
    return null
  }

  async function checkPort(port: number) {
    try {
      const response = await fetch(`http://localhost:${port}/`, {
        cache: 'no-cache',
        method: 'GET',
      })
      if ((await response.json()).__uni_token) {
        return true
      }
    }
    catch {}
    return false
  }

  async function tryStartService() {
    window.open('uni-token://start')
  }

  return {
    serverConnected,
    serviceHost,
    serviceUrl,
    servicePort,
    refreshService: () => {
      requireFindService = true
    },
    tryStartService,
    token,
    fetch: async (path: string, options?: RequestInit) => {
      await initialLoad
      if (!serviceUrl.value) {
        throw new Error('Service not available')
      }
      try {
        return await fetch(
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
      }
      catch (error) {
        requireFindService = true
        console.error(`Error fetching ${path}:`, error)
        throw error
      }
    },
  }
})
