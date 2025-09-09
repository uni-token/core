import type { Provider } from '@/lib/providers'
import { defineStore } from 'pinia'
import { useDeepSeekProvider } from '@/lib/providers/deepseek'
import { useOpenRouterProvider } from '@/lib/providers/openrouter'
import { useSiliconFlowProvider } from '@/lib/providers/siliconflow'

export const useProvidersStore = defineStore('providers', () => {
  const map = {
    siliconFlow: useSiliconFlowProvider(),
    deepSeek: useDeepSeekProvider(),
    openRouter: useOpenRouterProvider(),
  } as Record<string, Provider>
  const list = Object.values(map)

  for (const provider of list) {
    provider.refreshUser()
  }

  return {
    map,
    list,
  }
})
