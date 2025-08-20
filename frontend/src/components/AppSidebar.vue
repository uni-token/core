<script setup lang="ts">
import { AlertTriangle, BarChart3, Brain, Grid3X3, InfoIcon, Settings } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'
import ThemeToggle from '@/components/ThemeToggle.vue'
import {
  HoverCard,
  HoverCardContent,
  HoverCardTrigger,
} from '@/components/ui/hover-card'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '@/components/ui/sidebar'
import { useServiceStore } from '@/stores'
import LogoSvg from '/logo.svg?raw'

const { t } = useI18n()
const serviceStore = useServiceStore()
</script>

<template>
  <Sidebar>
    <SidebarHeader>
      <div class="flex items-center gap-2 px-2 py-2">
        <div class="rounded-lg bg-primary p-1 text-primary-foreground">
          <div class="h-6 w-6" v-html="LogoSvg" />
        </div>
        <div>
          <h1 class="text-lg font-semibold">
            {{ t('app.name') }}
          </h1>
          <p class="text-xs text-sidebar-foreground/70">
            {{ t('app.description') }}
          </p>
        </div>
      </div>
    </SidebarHeader>

    <SidebarContent>
      <SidebarGroup>
        <SidebarGroupLabel>{{ t('navigation.title') }}</SidebarGroupLabel>
        <SidebarGroupContent>
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton as-child :is-active="$route.path === '/apps'">
                <router-link to="/apps" class="flex items-center gap-2">
                  <Grid3X3 class="h-4 w-4" />
                  <span>{{ t('navigation.appManagement') }}</span>
                </router-link>
              </SidebarMenuButton>
            </SidebarMenuItem>

            <SidebarMenuItem>
              <SidebarMenuButton as-child :is-active="$route.path === '/models'">
                <router-link to="/models" class="flex items-center gap-2">
                  <Brain class="h-4 w-4" />
                  <span>{{ t('navigation.models') }}</span>
                </router-link>
              </SidebarMenuButton>
            </SidebarMenuItem>

            <SidebarMenuItem>
              <SidebarMenuButton as-child :is-active="$route.path === '/usage'">
                <router-link to="/usage" class="flex items-center gap-2">
                  <BarChart3 class="h-4 w-4" />
                  <span>{{ t('navigation.usage') }}</span>
                </router-link>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>

      <SidebarGroup>
        <SidebarGroupLabel>{{ t('navigation.tools') }}</SidebarGroupLabel>
        <SidebarGroupContent>
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton as-child :is-active="$route.path === '/settings'">
                <router-link to="/settings" class="flex items-center gap-2">
                  <Settings class="h-4 w-4" />
                  <span>{{ t('navigation.settings') }}</span>
                </router-link>
              </SidebarMenuButton>
            </SidebarMenuItem>

            <SidebarMenuItem>
              <SidebarMenuButton as-child :is-active="$route.path === '/about'">
                <router-link to="/about" class="flex items-center gap-2">
                  <InfoIcon class="h-4 w-4" />
                  <span>{{ t('navigation.about') }}</span>
                </router-link>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>
    </SidebarContent>

    <SidebarFooter>
      <div class="px-2 space-y-2">
        <div v-if="!serviceStore.serverConnected">
          <div class="rounded-lg border border-orange-200 bg-orange-50 p-2 dark:border-orange-900 dark:bg-orange-950">
            <div class="flex items-start gap-2">
              <AlertTriangle class="h-4 w-4 text-orange-600 dark:text-orange-400 mt-0.5 flex-shrink-0" />
              <div class="space-y-1">
                <p class="text-xs font-medium text-orange-800 dark:text-orange-200">
                  {{ t('service.connectionFailed') }}
                </p>
                <p class="text-xs text-orange-700 dark:text-orange-300">
                  {{ t('service.reconnecting') }}
                </p>
                <p class="text-xs text-orange-700 dark:text-orange-300">
                  {{ t('service.restartAgent') }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <div class="flex items-center justify-between">
          <HoverCard :open-delay="400">
            <HoverCardTrigger>
              <div class="h-8 flex items-center gap-1 text-xs text-sidebar-foreground/70">
                <div
                  class="mx-1 h-2 w-2 rounded-full mb-[2px]"
                  :class="serviceStore.serverConnected ? 'bg-green-500' : 'bg-red-500'"
                />
                <span class="text-sm select-none">{{ serviceStore.serverConnected ? t('service.connected') : t('service.disconnected') }}</span>
              </div>
            </HoverCardTrigger>
            <HoverCardContent v-if="serviceStore.serverConnected" :side-offset="0" class="pt-2 py-1">
              <div class="font-mono text-sm text-center mt-1 ">
                {{ serviceStore.serviceUrl }}
              </div>
            </HoverCardContent>
          </HoverCard>
          <ThemeToggle />
        </div>
      </div>
    </SidebarFooter>
  </Sidebar>
</template>
