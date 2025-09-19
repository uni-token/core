import type { Provider } from '@/lib/providers'
import { defineStore } from 'pinia'
import { useDeepSeekProvider } from '@/lib/providers/deepseek'
import { useOpenRouterProvider } from '@/lib/providers/openrouter'
import { useSiliconFlowProvider } from '@/lib/providers/siliconflow'

export const useProvidersStore = defineStore('providers', () => {
  const list = [
    useSiliconFlowProvider(),
    useDeepSeekProvider(),
    useOpenRouterProvider(),
  ]
  const map = Object.fromEntries(list.map(p => [p.id, p])) as Record<string, Provider>

  for (const provider of list) {
    provider.refreshUser()
  }

  return {
    list,
    map,
  }
})
