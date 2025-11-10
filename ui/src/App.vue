<script setup lang="ts">
import type { AuthState } from './stores'
import Clarity from '@microsoft/clarity'
import { Zap } from 'lucide-vue-next'
import { onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import AppSidebar from '@/components/AppSidebar.vue'
import LoginForm from '@/components/LoginForm.vue'
import { SidebarInset, SidebarProvider, SidebarTrigger } from '@/components/ui/sidebar'
import { Toaster } from '@/components/ui/sonner'
import { useThemeStore } from '@/stores/theme'
import { useAuthStore, useServiceStore } from './stores'
import 'vue-sonner/style.css'

const { t } = useI18n()
const authStore = useAuthStore()
const { currentRoute } = useRouter()
const themeStore = useThemeStore()
const serviceStore = useServiceStore()

const authState = ref<AuthState | null>(null)
const isCheckingAuth = ref(true)

onMounted(async () => {
  // Initialize theme
  themeStore.initTheme()
  themeStore.setupSystemThemeListener()

  authState.value = await authStore.checkAuth()
  isCheckingAuth.value = false
  setTimeout(() => currentRoute.value.query = {}, 100)
  Clarity.init('sx60zbxtfz')
})

watch(() => serviceStore.serverConnected, async (connected) => {
  if (connected) {
    authState.value = await authStore.checkAuth()
  }
})

async function onUIOpened() {
  const params = new URLSearchParams(window.location.search)
  const session = params.get('session')
  if (!session)
    return
  pingActive()

  async function pingActive() {
    try {
      const resp = await serviceStore.api('ui/active', {
        method: 'POST',
        body: JSON.stringify({ session }),
      })
      const data = await resp.json()
      if (data.continue !== false) {
        setTimeout(pingActive, 1000)
      }
    }
    catch {
      setTimeout(pingActive, 1000)
    }
  }
}
onUIOpened()
</script>

<template>
  <div class="min-h-screen bg-background">
    <Toaster />

    <div v-if="isCheckingAuth" class="min-h-screen flex items-center justify-center">
      <div class="text-center space-y-4">
        <div class="w-8 h-8 border-2 border-primary border-t-transparent rounded-full animate-spin mx-auto" />
        <p class="text-muted-foreground">
          {{ t('loading') }}
        </p>
      </div>
    </div>

    <LoginForm v-else-if="!authStore.isLoggedIn" :register="authState?.status === 'not_registered'" />

    <SidebarProvider v-else>
      <AppSidebar />
      <SidebarInset>
        <div class="md:hidden bg-background border-b border-border sticky top-0 z-50">
          <div class="flex items-center gap-3 h-14 px-4">
            <SidebarTrigger class="h-8 w-8" />
            <div class="flex items-center gap-2">
              <div class="rounded-lg bg-primary p-1.5 text-primary-foreground">
                <Zap class="h-3 w-3" />
              </div>
              <h1 class="text-base font-semibold">
                {{ t('appName') }}
              </h1>
            </div>
          </div>
        </div>
        <RouterView />
      </SidebarInset>
    </SidebarProvider>
  </div>
</template>

<i18n lang="yaml">
zh-CN:
  loading: 加载中...
  appName: UniToken
en-US:
  loading: Loading...
  appName: UniToken
</i18n>

<style>
/* Global styles using UnoCSS */
</style>
