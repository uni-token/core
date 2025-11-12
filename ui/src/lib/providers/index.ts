import type { Component } from 'vue'
import { createSharedComposable } from '@vueuse/core'
import { markRaw } from 'vue'
import { defineDbStore } from '@/stores/db'
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

export interface Provider<A = unknown> {
  readonly id: string
  readonly name: string
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

  readonly payment?: {
    createWeChatPay: (options: {
      amount: string
    }) => Promise<{
      orderId: string
      qrcUrl: string
      interval?: number
      timeout?: number
    }>

    checkWeChatPay: (options: {
      orderId: string
    }) => Promise<'success' | 'wait' | 'canceled'>
  } | {
    readonly websiteURL: string
  }

  readonly baseURL: string
  readonly createKey: () => Promise<string>

  readonly apis: A
}

const useProviderSessionsDb = defineDbStore<unknown>('provider_sessions')

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
    useSiliconFlowProvider(),
    useDeepSeekProvider(),
    useOpenRouterProvider(),
  ]
  const map = Object.fromEntries(list.map(p => [p.id, markRaw(p)])) as Record<string, Provider>

  for (const provider of list) {
    provider.refreshUser()
  }

  return map
})
