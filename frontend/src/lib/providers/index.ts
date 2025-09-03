import type { Component } from 'vue'

export interface ProviderUserInfo {
  name: string
  // Real-name verification status
  verified: boolean
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

  readonly payment: {
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
