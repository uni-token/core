import { createSharedComposable, useIntervalFn } from '@vueuse/core'
import { computed, ref } from 'vue'

export const useService = createSharedComposable(() => {
  const serverConnected = ref(true)
  const servicePort = ref<number | null>(null)
  const serviceUrl = computed(() => {
    if (servicePort.value) {
      return `http://localhost:${servicePort.value}/`
    }
    return null
  })
  let requireFindService = true
  let failureSince = 0

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
            failureSince = 0
          }
          else {
            failureSince ||= Date.now()
            if (Date.now() - failureSince > 3 * 1000) {
              serverConnected.value = false
            }
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
      if ((await response.json()).__united_token) {
        return true
      }
    }
    catch {}
    return false
  }

  return {
    serverConnected,
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
