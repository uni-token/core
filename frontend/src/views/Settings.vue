<script setup lang="ts">
import { LogOut, Trash2 } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import LanguageSelector from '@/components/LanguageSelector.vue'
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
import { useAuth } from '@/composables/auth'
import { useService } from '@/composables/service'
import { useAppStore } from '@/stores/app'

const { t } = useI18n()
const router = useRouter()
const { logout, currentUser } = useAuth()
const { fetch, serverConnected } = useService()
const { deleteAllApps } = useAppStore()

const clearingRecords = ref(false)
const deletingApps = ref(false)

const serviceStatus = computed(() => {
  return serverConnected.value ? t('service.connected') : t('service.disconnected')
})

function handleLogout() {
  logout()
  toast.success(t('settings.logoutSuccess'))
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
      toast.success(t('settings.clearRecordsSuccess'))
    }
    else {
      const errorData = await response.json()
      toast.error(errorData.error || t('settings.clearRecordsFailed'))
    }
  }
  catch (error) {
    console.error('Clear records failed:', error)
    toast.error(t('settings.clearRecordsFailedRetry'))
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
          {{ t('settings.title') }}
        </h2>
      </div>

      <Card>
        <CardHeader>
          <CardTitle class="text-lg">
            {{ t('settings.general') }}
          </CardTitle>
          <CardDescription>{{ t('settings.generalDescription') }}</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            <LanguageSelector />
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle class="text-lg">
            {{ t('settings.account') }}
          </CardTitle>
          <CardDescription>{{ t('settings.accountDescription') }}</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <div>
                <h4 class="text-sm font-medium">
                  {{ t('settings.currentUser') }}
                </h4>
                <p class="text-sm text-muted-foreground">
                  {{ currentUser || t('settings.unknownUser') }}
                </p>
              </div>
            </div>

            <div class="flex items-center justify-between pt-4 border-t">
              <div>
                <h4 class="text-sm font-medium">
                  {{ t('settings.logout') }}
                </h4>
                <p class="text-sm text-muted-foreground">
                  {{ t('settings.logoutDescription') }}
                </p>
              </div>
              <AlertDialog>
                <AlertDialogTrigger as-child>
                  <Button variant="outline" class="text-red-600 hover:text-red-700">
                    <LogOut class="mr-2 h-4 w-4" />
                    {{ t('settings.logout') }}
                  </Button>
                </AlertDialogTrigger>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>{{ t('settings.confirmLogout') }}</AlertDialogTitle>
                    <AlertDialogDescription>
                      {{ t('settings.confirmLogoutDescription') }}
                    </AlertDialogDescription>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>{{ t('settings.cancel') }}</AlertDialogCancel>
                    <AlertDialogAction class="bg-red-600 hover:bg-red-700" @click="handleLogout">
                      {{ t('settings.confirm') }}
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
            {{ t('settings.dataManagement') }}
          </CardTitle>
          <CardDescription>{{ t('settings.dataManagementDescription') }}</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <div>
                <h4 class="text-sm font-medium">
                  {{ t('settings.clearRecords') }}
                </h4>
                <p class="text-sm text-muted-foreground">
                  {{ t('settings.clearRecordsDescription') }}
                </p>
              </div>
              <AlertDialog>
                <AlertDialogTrigger as-child>
                  <Button variant="outline" class="text-red-600 hover:text-red-700" :disabled="clearingRecords">
                    <Trash2 class="mr-2 h-4 w-4" />
                    {{ clearingRecords ? t('settings.clearing') : t('settings.clearRecords') }}
                  </Button>
                </AlertDialogTrigger>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>{{ t('settings.confirmClearRecords') }}</AlertDialogTitle>
                    <AlertDialogDescription>
                      {{ t('settings.confirmClearRecordsDescription') }}
                    </AlertDialogDescription>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>{{ t('settings.cancel') }}</AlertDialogCancel>
                    <AlertDialogAction class="bg-red-600 hover:bg-red-700" @click="handleClearRecords">
                      {{ t('settings.confirm') }}
                    </AlertDialogAction>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            </div>

            <div class="flex items-center justify-between pt-4 border-t">
              <div>
                <h4 class="text-sm font-medium">
                  {{ t('settings.deleteAllApps') }}
                </h4>
                <p class="text-sm text-muted-foreground">
                  {{ t('settings.deleteAllAppsDescription') }}
                </p>
              </div>
              <AlertDialog>
                <AlertDialogTrigger as-child>
                  <Button variant="outline" class="text-red-600 hover:text-red-700" :disabled="deletingApps">
                    <Trash2 class="mr-2 h-4 w-4" />
                    {{ deletingApps ? t('settings.deleting') : t('settings.deleteAllApps') }}
                  </Button>
                </AlertDialogTrigger>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>{{ t('settings.confirmDeleteAllApps') }}</AlertDialogTitle>
                    <AlertDialogDescription>
                      {{ t('settings.confirmDeleteAllAppsDescription') }}
                    </AlertDialogDescription>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>{{ t('settings.cancel') }}</AlertDialogCancel>
                    <AlertDialogAction class="bg-red-600 hover:bg-red-700" @click="handleDeleteAllApps">
                      {{ t('settings.confirmDelete') }}
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
            {{ t('settings.appInfo') }}
          </CardTitle>
          <CardDescription>{{ t('settings.appInfoDescription') }}</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="space-y-3 text-sm">
            <div class="flex justify-between">
              <span class="text-muted-foreground">{{ t('settings.appNameLabel') }}:</span>
              <span>{{ t('app.name') }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-muted-foreground">{{ t('settings.versionLabel') }}:</span>
              <span>v1.0.0</span>
            </div>
            <div class="flex justify-between">
              <span class="text-muted-foreground">{{ t('settings.serviceStatus') }}:</span>
              <div class="flex items-center gap-2">
                <div class="h-2 w-2 rounded-full" :class="serverConnected ? 'bg-green-500' : 'bg-red-500'" />
                <span>{{ serviceStatus }}</span>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
