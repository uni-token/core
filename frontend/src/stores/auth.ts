import { useLocalStorage } from '@vueuse/core'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { useServiceStore } from './service'

export interface AuthState {
  status: 'success' | 'error' | 'not_registered'
  message?: string
  username?: string
  token?: string
}

export const useAuthStore = defineStore('auth', () => {
  const serviceStore = useServiceStore()

  const params = new URLSearchParams(window.location.search)
  const currentUser = ref<string | null>(params.get('username') || null)
  const isLoggedIn = computed(() => !!currentUser.value)
  const isLoading = ref(false)

  const savedUsername = useLocalStorage('uni_token_username', '')
  const savedPassword = useLocalStorage('uni_token_password', '')

  function checkAuth() {
    return login(savedUsername.value, savedPassword.value)
  }

  async function login(username: string, password: string): Promise<AuthState> {
    isLoading.value = true
    try {
      const response = await serviceStore.fetch('auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
      })

      const result: AuthState = await response.json()

      if (result.status === 'success' && result.username) {
        currentUser.value = result.username
        serviceStore.token = result.token!
        savedUsername.value = username
        savedPassword.value = password
      }

      return result
    }
    catch (error) {
      console.error('Login failed:', error)
      return {
        status: 'error',
        message: 'Connection error',
      }
    }
    finally {
      isLoading.value = false
    }
  }

  async function register(
    username: string,
    password: string,
  ): Promise<AuthState> {
    isLoading.value = true
    try {
      const response = await serviceStore.fetch('auth/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
      })

      const result: AuthState = await response.json()

      if (result.status === 'success' && result.username) {
        currentUser.value = result.username
        serviceStore.token = result.token!
        savedUsername.value = username
        savedPassword.value = password
      }

      return result
    }
    catch (error) {
      console.error('Register failed:', error)
      return {
        status: 'error',
        message: 'Connection error',
      }
    }
    finally {
      isLoading.value = false
    }
  }

  async function logout() {
    currentUser.value = null
    serviceStore.token = null
    savedUsername.value = ''
    savedPassword.value = ''
  }

  return {
    currentUser,
    isLoggedIn,
    isLoading,
    checkAuth,
    login,
    register,
    logout,
  }
})
