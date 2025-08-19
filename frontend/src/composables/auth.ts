import { createSharedComposable, useLocalStorage } from '@vueuse/core'
import { computed, ref } from 'vue'
import { useService } from './service'

export interface AuthState {
  status: 'success' | 'error' | 'not_registered'
  message?: string
  username?: string
  token?: string
}

export const useAuth = createSharedComposable(() => {
  const { fetch: serviceFetch, token } = useService()

  const params = new URLSearchParams(window.location.search)
  const currentUser = ref<string | null>(params.get('username') || null)
  const isLoggedIn = computed(() => !!currentUser.value)
  const isLoading = ref(false)

  const savedUsername = useLocalStorage('united_token_username', '')
  const savedPassword = useLocalStorage('united_token_password', '')

  function checkAuth() {
    return login(savedUsername.value, savedPassword.value)
  }

  async function login(username: string, password: string): Promise<AuthState> {
    isLoading.value = true
    try {
      const response = await serviceFetch('auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
      })

      const result: AuthState = await response.json()

      if (result.status === 'success' && result.username) {
        currentUser.value = result.username
        token.value = result.token!
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
      const response = await serviceFetch('auth/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
      })

      const result: AuthState = await response.json()

      if (result.status === 'success' && result.username) {
        currentUser.value = result.username
        token.value = result.token!
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
    token.value = null
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
