<script setup lang="ts">
import { LogOut, Trash2 } from 'lucide-vue-next'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import LanguageSelector from '@/components/LanguageSelector.vue'
import ThemeSelector from '@/components/ThemeSelector.vue'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { useAuthStore, useServiceStore } from '@/stores'
import { useAppStore } from '@/stores/app'

const { t } = useI18n()
const router = useRouter()
const authStore = useAuthStore()
const serviceStore = useServiceStore()
const { deleteAllApps } = useAppStore()

const clearingRecords = ref(false)
const deletingApps = ref(false)

function handleLogout() {
  authStore.logout()
  toast.success(t('logoutSuccess'))
  router.push('/')
}

async function handleClearRecords() {
  if (clearingRecords.value)
    return

  clearingRecords.value = true

  try {
    const response = await fetch('usage/clear', {
      method: 'POST',
    })

    if (response.ok) {
      toast.success(t('clearRecordsSuccess'))
    }
    else {
      const errorData = await response.json()
      toast.error(errorData.error || t('clearRecordsFailed'))
    }
  }
  catch (error) {
    console.error('Clear records failed:', error)
    toast.error(t('clearRecordsFailedRetry'))
  }
  finally {
    clearingRecords.value = false
  }
}

async function handleDeleteAllApps() {
  if (deletingApps.value)
    return

  deletingApps.value = true
  try {
    await deleteAllApps()
  }
  catch (error) {
    console.error('Delete all apps failed:', error)
  }
  finally {
    deletingApps.value = false
  }
}
</script>

<template>
  <div class="p-6">
    <div class="space-y-6">
      <div class="flex items-center justify-between">
        <h2 class="text-2xl font-bold">
          {{ t('title') }}
        </h2>
      </div>

      <Card>
        <CardHeader>
          <CardTitle class="text-lg">
            {{ t('general') }}
          </CardTitle>
          <CardDescription>{{ t('generalDescription') }}</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            <LanguageSelector />
            <ThemeSelector />
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle class="text-lg">
            {{ t('account') }}
          </CardTitle>
          <CardDescription>{{ t('accountDescription') }}</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <div>
                <h4 class="text-sm font-medium">
                  {{ t('currentUser') }}
                </h4>
                <p class="text-sm text-muted-foreground">
                  {{ authStore.currentUser || t('unknownUser') }}
                </p>
              </div>
            </div>

            <div class="flex items-center justify-between pt-4 border-t">
              <div>
                <h4 class="text-sm font-medium">
                  {{ t('logout') }}
                </h4>
                <p class="text-sm text-muted-foreground">
                  {{ t('logoutDescription') }}
                </p>
              </div>
              <AlertDialog>
                <AlertDialogTrigger as-child>
                  <Button variant="outline" class="text-red-600 hover:text-red-700">
                    <LogOut class="mr-2 h-4 w-4" />
                    {{ t('logout') }}
                  </Button>
                </AlertDialogTrigger>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>{{ t('confirmLogout') }}</AlertDialogTitle>
                    <AlertDialogDescription>
                      {{ t('confirmLogoutDescription') }}
                    </AlertDialogDescription>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>{{ t('cancel') }}</AlertDialogCancel>
                    <AlertDialogAction class="bg-red-600 hover:bg-red-700 text-white" @click="handleLogout">
                      {{ t('confirm') }}
                    </AlertDialogAction>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle class="text-lg">
            {{ t('dataManagement') }}
          </CardTitle>
          <CardDescription>{{ t('dataManagementDescription') }}</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <div>
                <h4 class="text-sm font-medium">
                  {{ t('clearRecords') }}
                </h4>
                <p class="text-sm text-muted-foreground">
                  {{ t('clearRecordsDescription') }}
                </p>
              </div>
              <AlertDialog>
                <AlertDialogTrigger as-child>
                  <Button variant="outline" class="text-red-600 hover:text-red-700" :disabled="clearingRecords">
                    <Trash2 class="mr-2 h-4 w-4" />
                    {{ clearingRecords ? t('clearing') : t('clearRecords') }}
                  </Button>
                </AlertDialogTrigger>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>{{ t('confirmClearRecords') }}</AlertDialogTitle>
                    <AlertDialogDescription>
                      {{ t('confirmClearRecordsDescription') }}
                    </AlertDialogDescription>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>{{ t('cancel') }}</AlertDialogCancel>
                    <AlertDialogAction class="bg-red-600 hover:bg-red-700 text-white" @click="handleClearRecords">
                      {{ t('confirm') }}
                    </AlertDialogAction>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            </div>

            <div class="flex items-center justify-between pt-4 border-t">
              <div>
                <h4 class="text-sm font-medium">
                  {{ t('deleteAllApps') }}
                </h4>
                <p class="text-sm text-muted-foreground">
                  {{ t('deleteAllAppsDescription') }}
                </p>
              </div>
              <AlertDialog>
                <AlertDialogTrigger as-child>
                  <Button variant="outline" class="text-red-600 hover:text-red-700" :disabled="deletingApps">
                    <Trash2 class="mr-2 h-4 w-4" />
                    {{ deletingApps ? t('deleting') : t('deleteAllApps') }}
                  </Button>
                </AlertDialogTrigger>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>{{ t('confirmDeleteAllApps') }}</AlertDialogTitle>
                    <AlertDialogDescription>
                      {{ t('confirmDeleteAllAppsDescription') }}
                    </AlertDialogDescription>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>{{ t('cancel') }}</AlertDialogCancel>
                    <AlertDialogAction class="bg-red-600 hover:bg-red-700 text-white" @click="handleDeleteAllApps">
                      {{ t('confirmDelete') }}
                    </AlertDialogAction>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle class="text-lg">
            {{ t('appInfo') }}
          </CardTitle>
          <CardDescription>{{ t('appInfoDescription') }}</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="space-y-3 text-sm">
            <div class="flex justify-between">
              <span class="text-muted-foreground">{{ t('appNameLabel') }}:</span>
              <span> UniToken </span>
            </div>
            <div class="flex justify-between">
              <span class="text-muted-foreground">{{ t('versionLabel') }}:</span>
              <span>v1.0.0</span>
            </div>
            <div class="flex justify-between">
              <span class="text-muted-foreground">{{ t('serviceStatus') }}:</span>
              <div class="flex items-center gap-2">
                <div class="h-2 w-2 rounded-full" :class="serviceStore.serverConnected ? 'bg-green-500' : 'bg-red-500'" />
                <span>{{ serviceStore.serviceHost || t('disconnected') }}</span>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>

<i18n lang="yaml">
zh-CN:
  title: 设置
  general: 通用设置
  language: 语言
  theme: 主题
  themeDescription: 选择您偏好的颜色主题
  darkMode: 深色模式
  lightMode: 浅色模式
  systemMode: 跟随系统
  notifications: 通知设置
  emailNotifications: 邮件通知
  pushNotifications: 推送通知
  account: 账户设置
  currentUser: 当前用户
  logout: 退出登录
  logoutSuccess: 已成功退出登录
  dangerZone: 危险操作
  clearRecords: 清空使用记录
  clearRecordsDescription: 清空所有 Token 使用记录，此操作不可撤销
  clearRecordsSuccess: 使用记录已清空
  clearRecordsFailed: 清空记录失败
  clearRecordsFailedRetry: 清空记录失败，请稍后重试
  deleteAllApps: 删除所有应用
  deleteAllAppsDescription: 删除所有注册的应用，此操作不可撤销
  systemInfo: 系统信息
  serviceStatus: 服务状态
  dataManagement: 数据管理
  dataManagementDescription: 管理您的应用数据和使用记录
  accountDescription: 管理您的账户信息和安全设置
  generalDescription: 配置应用的基本设置和偏好
  unknownUser: 未知用户
  logoutDescription: 退出当前账户，返回登录页面
  confirmLogout: 确认退出登录
  confirmLogoutDescription: 退出登录后，您需要重新输入账户信息才能登录。
  confirmClearRecords: 确认清空记录
  confirmClearRecordsDescription: 此操作将永久删除所有使用记录和统计数据，无法恢复。请谨慎操作。
  clearing: 清空中...
  confirmDelete: 确认删除
  confirm: 确认
  cancel: 取消
  confirmDeleteAllApps: 确认删除所有应用
  confirmDeleteAllAppsDescription: 此操作将永久删除所有已连接的应用和相关数据，无法恢复。请谨慎操作。
  deleting: 删除中...
  appInfo: 应用信息
  appInfoDescription: 关于 UniToken 的版本和服务信息
  appNameLabel: 应用名称
  versionLabel: 版本

en-US:
  title: Settings
  general: General Settings
  language: Language
  theme: Theme
  themeDescription: Choose your preferred color theme
  darkMode: Dark Mode
  lightMode: Light Mode
  systemMode: Follow System
  notifications: Notification Settings
  emailNotifications: Email Notifications
  pushNotifications: Push Notifications
  account: Account Settings
  currentUser: Current User
  logout: Logout
  logoutSuccess: Successfully logged out
  dangerZone: Danger Zone
  clearRecords: Clear Usage Records
  clearRecordsDescription: Clear all token usage records. This action cannot be undone
  clearRecordsSuccess: Usage records cleared
  clearRecordsFailed: Failed to clear records
  clearRecordsFailedRetry: 'Failed to clear records, please try again later'
  deleteAllApps: Delete All Apps
  deleteAllAppsDescription: Delete all registered applications. This action cannot be undone
  systemInfo: System Information
  serviceStatus: Service Status
  dataManagement: Data Management
  dataManagementDescription: Manage your application data and usage records
  accountDescription: Manage your account information and security settings
  generalDescription: Configure basic application settings and preferences
  unknownUser: Unknown User
  logoutDescription: Sign out of current account and return to login page
  confirmLogout: Confirm Logout
  confirmLogoutDescription: >-
    After logging out, you will need to re-enter your account information to log
    in again.
  confirmClearRecords: Confirm Clear Records
  confirmClearRecordsDescription: >-
    This action will permanently delete all usage records and statistics data
    and cannot be recovered. Please proceed with caution.
  clearing: Clearing...
  confirmDelete: Confirm Delete
  confirm: Confirm
  cancel: Cancel
  confirmDeleteAllApps: Confirm Delete All Apps
  confirmDeleteAllAppsDescription: >-
    This action will permanently delete all connected applications and related
    data and cannot be recovered. Please proceed with caution.
  deleting: Deleting...
  appInfo: Application Information
  appInfoDescription: Version and service information about UniToken
  appNameLabel: Application Name
  versionLabel: Version
</i18n>
