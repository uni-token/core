<script setup lang="ts">
import { BarChart3, Brain, Grid3X3, InfoIcon, Settings } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'
import ServiceStatus from '@/components/ServiceStatus.vue'
import ThemeToggle from '@/components/ThemeToggle.vue'
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
import LogoSvg from '/logo.svg?raw'

const { t } = useI18n()
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
      <div class="px-2">
        <ServiceStatus variant="full">
          <template #actions>
            <ThemeToggle />
          </template>
        </ServiceStatus>
      </div>
    </SidebarFooter>
  </Sidebar>
</template>
