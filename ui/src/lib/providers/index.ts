import type { Component } from 'vue'
import { createSharedComposable } from '@vueuse/core'
import { markRaw } from 'vue'
import { useProviderSessionsDb } from '@/stores/db'
import { useAIHubMixProvider } from './aihubmix'
import { useDeepSeekProvider } from './deepseek'
import { useOpenRouterProvider } from './openrouter'
import { useSiliconFlowProvider } from './siliconflow'

export interface ProviderUserInfo {
  name: string
  // Real-name verification status
  verified?: boolean
  phone?: string
  email?: string
  balance?: {
    amount: number
    currency?: 'USD' | 'CNY' | string
  }
}

export interface ProviderVerificationInfo {
  name: string
  cardId: string
  time?: number
}

export interface ProviderPaymentQRC {
  readonly type: 'qrc'
  readonly currency: 'CNY' | 'USD'
  readonly platform: string
  readonly create: (options: {
    amount: string
  }) => Promise<{
    orderId: string
    qrcUrl: string
    interval?: number
    timeout?: number
  }>
  readonly check: (options: {
    orderId: string
  }) => Promise<'success' | 'wait' | 'canceled'>
}

export interface ProviderPaymentWebsite {
  readonly type: 'website'
  readonly websiteURL: string
}

export type ProviderPayment = ProviderPaymentQRC | ProviderPaymentWebsite

export interface Provider<A = unknown> {
  readonly id: string
  readonly name: string
  readonly description: string
  /**
   * Reference: https://github.com/CherryHQ/cherry-studio/tree/main/src/renderer/src/assets/images/providers
   */
  readonly logo: string
  readonly homepage: string

  /**
   * - `undefined`: loading
   * - `null`: not logged in
   * - `ProviderUserInfo`: logged in
   */
  readonly user: undefined | null | ProviderUserInfo
  readonly refreshUser: () => Promise<void>

  readonly Login: Component
  readonly logout: () => Promise<void>

  readonly verification?: {
    check: () => Promise<null | ProviderVerificationInfo>

    readonly cardTypes: Array<{ value: number, label: string }>

    submit: (data: {
      name: string
      cardType: number
      cardId: string
    }) => Promise<'success' | 'failed' | {
      qrcUrl: string
    }>
  }

  readonly payment?: ProviderPayment

  readonly baseURL: string
  readonly createKey: () => Promise<string>

  readonly apis: A
}

export function useProviderSession<T>(providerId: string) {
  const db = useProviderSessionsDb()

  return {
    get() {
      return db.get(providerId) as Promise<T | null>
    },
    put(session: T) {
      return db.put(providerId, session)
    },
    delete() {
      return db.delete(providerId)
    },
  }
}

export function defineProvider<A>(provider: () => Provider<A>): () => Provider<A> {
  return createSharedComposable(provider)
}

export const useProviders = createSharedComposable(() => {
  const list = [
    useOpenRouterProvider(),
    useSiliconFlowProvider(),
    useDeepSeekProvider(),
    useAIHubMixProvider(),
  ]
  const map = Object.fromEntries(list.map(p => [p.id, markRaw(p)])) as Record<string, Provider>

  for (const provider of list) {
    provider.refreshUser()
  }

  return map
})
