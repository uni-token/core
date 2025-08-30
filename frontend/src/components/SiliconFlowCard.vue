<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue-sonner'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { useSiliconFlowStore } from '@/stores'
import SiliconFlowLoginCard from './SiliconFlowLoginCard.vue'
import SiliconFlowPaymentDialog from './SiliconFlowPaymentDialog.vue'
import SiliconFlowRealName from './SiliconFlowRealName.vue'

const props = defineProps<{
  showTitle?: boolean
}>()

const emit = defineEmits<{
  configured: [provider: string]
}>()

const { t } = useI18n()
const siliconFlowStore = useSiliconFlowStore()

const openRealNameDialog = ref(false)

// Check login status on component mount
onMounted(async () => {
  await siliconFlowStore.init()
})

function handleRecharge() {
  if (!siliconFlowStore.authed) {
    toast.error(t('siliconFlow.realNameRequired'))
    openRealNameDialog.value = true
  }
  else {
    siliconFlowStore.openPaymentDialog()
  }
}
</script>

<template>
  <!-- User logged in state -->
  <Card v-if="siliconFlowStore.userInfo?.isLoggedIn">
    <CardHeader>
      <CardTitle class="flex items-center gap-2">
        <div v-if="props.showTitle" class="mr-4 text-lg">
          {{ t('providers.siliconFlow') }}
        </div>

        <div class="h-2 w-2 bg-green-500 dark:bg-green-400 rounded-full" />
        <div class="text-green-700 dark:text-green-500 text-base">
          {{ t('siliconFlow.loggedIn') }}
        </div>

        <div class="flex-grow" />

        <span class="font-medium text-black">{{ siliconFlowStore.userInfo.data?.name }}</span>
      </CardTitle>
    </CardHeader>
    <CardContent>
      <div v-if="siliconFlowStore.userInfo.data" class="space-y-2">
        <div class="flex justify-between items-center text-sm">
          <span class="text-muted-foreground">{{ t('siliconFlow.realname') }}</span>
          <button
            class="font-medium border-b border-dashed border-primary text-primary hover:bg-muted/50 transition-colors"
            @click="openRealNameDialog = true"
          >
            {{ siliconFlowStore.authed ? t('siliconFlow.auth-success') : t('siliconFlow.not-authenticated') }}
          </button>
        </div>
        <div v-if="siliconFlowStore.userInfo.data.phone" class="flex justify-between items-center text-sm">
          <span class="text-muted-foreground">{{ t('siliconFlow.phone') }}</span>
          <span class="font-medium">{{ siliconFlowStore.userInfo.data.phone }}</span>
        </div>
        <div v-if="siliconFlowStore.userInfo.data.totalBalance" class="flex justify-between items-center text-sm">
          <span class="text-muted-foreground">{{ t('siliconFlow.balance') }}</span>
          <span class="font-medium">{{ siliconFlowStore.userInfo.data.totalBalance }} {{ t('siliconFlow.yuan') }}</span>
        </div>
      </div>
    </CardContent>
    <CardFooter class="pt-2 flex gap-2">
      <!-- <AlertDialog>
        <AlertDialogTrigger as-child>
          <Button variant="default" size="sm" class="h-8">
            {{ t('siliconFlow.createKey') }}
          </Button>
        </AlertDialogTrigger>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{{ t('siliconFlow.confirmCreateKey') }}</AlertDialogTitle>
            <AlertDialogDescription>
              {{ t('siliconFlow.confirmCreateKeyDescription') }}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>{{ t('common.cancel') }}</AlertDialogCancel>
            <AlertDialogAction @click="handleCreateKey">
              {{ t('common.confirm') }}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog> -->
      <Button variant="default" size="sm" class="h-8" @click="emit('configured', siliconFlowStore.providerId!)">
        {{ t('common.confirm') }}
      </Button>
      <div class="flex-grow" />
      <Button variant="secondary" size="sm" class="h-8" @click="handleRecharge">
        {{ t('siliconFlow.recharge') }}
      </Button>
      <Button variant="secondary" size="sm" class="h-8" @click="siliconFlowStore.logout">
        {{ t('siliconFlow.logout') }}
      </Button>
    </CardFooter>

    <!-- Real Name Dialog -->
    <SiliconFlowRealName v-model="openRealNameDialog" />

    <!-- Payment Dialog -->
    <SiliconFlowPaymentDialog v-model="siliconFlowStore.showPaymentDialog" />
  </Card>

  <!-- User not logged in state -->
  <SiliconFlowLoginCard v-else-if="siliconFlowStore.userInfo" />

  <!-- Loading state -->
  <Skeleton v-else class="h-48 w-full max-w-md mx-auto" />
</template>
