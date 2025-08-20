<script setup lang="ts">
import type { AuthState } from '@/composables/auth'
import Clarity from '@microsoft/clarity'
import { Zap } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import AppSidebar from '@/components/AppSidebar.vue'
import LoginForm from '@/components/LoginForm.vue'
import ThemeToggle from '@/components/ThemeToggle.vue'
import { SidebarInset, SidebarProvider, SidebarTrigger } from '@/components/ui/sidebar'
import { Toaster } from '@/components/ui/sonner'
import { useAuth } from '@/composables/auth'
import { useThemeStore } from '@/stores/theme'
import ActionHandler from '@/views/ActionHandler.vue'
import 'vue-sonner/style.css'

const { t } = useI18n()
const { checkAuth, isLoggedIn } = useAuth()
const { replace, currentRoute } = useRouter()
const themeStore = useThemeStore()

const authState = ref<AuthState | null>(null)
const isCheckingAuth = ref(true)

onMounted(async () => {
  // Initialize theme
  themeStore.initTheme()
  themeStore.setupSystemThemeListener()

  authState.value = await checkAuth()
  isCheckingAuth.value = false
  setTimeout(() => replace(currentRoute.value.path), 100)
  Clarity.init('sx60zbxtfz')
})
</script>

<template>
  <div class="min-h-screen bg-background">
    <Toaster />

    <div v-if="isCheckingAuth" class="min-h-screen flex items-center justify-center">
      <div class="text-center space-y-4">
        <div class="w-8 h-8 border-2 border-primary border-t-transparent rounded-full animate-spin mx-auto" />
        <p class="text-muted-foreground">
          {{ t('common.loading') }}
        </p>
      </div>
    </div>

    <LoginForm v-else-if="!isLoggedIn" :register="authState?.status === 'not_registered'" />

    <SidebarProvider v-else>
      <AppSidebar />
      <SidebarInset>
        <div class="md:hidden bg-background border-b border-border sticky top-0 z-50">
          <div class="flex items-center justify-between h-14 px-4">
            <div class="flex items-center gap-3">
              <SidebarTrigger class="h-8 w-8" />
              <div class="flex items-center gap-2">
                <div class="rounded-lg bg-primary p-1.5 text-primary-foreground">
                  <Zap class="h-3 w-3" />
                </div>
                <h1 class="text-base font-semibold">
                  {{ t('app.name') }}
                </h1>
              </div>
            </div>
            <ThemeToggle />
          </div>
        </div>
        <RouterView />
      </SidebarInset>
    </SidebarProvider>

    <ActionHandler />
  </div>
</template>

<style>
/* Global styles using UnoCSS */
</style>
