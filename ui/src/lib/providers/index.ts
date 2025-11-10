import type { Component } from 'vue'
import { createSharedComposable } from '@vueuse/core'
import { defineDbStore } from '@/stores/db'

export interface ProviderUserInfo {
  name: string
  // Real-name verification status
  verified?: boolean
  phone?: string
  email?: string
  balance?: number
}

export interface ProviderVerificationInfo {
  name: string
  cardId: string
  time?: number
}

export interface Provider {
  readonly id: string
  readonly name: string
  readonly homepage: string

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
  }

  readonly baseURL: string
  readonly createKey: () => Promise<string>
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

export function defineProvider<P extends Provider>(provider: () => P): () => P {
  return createSharedComposable(provider)
}
