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
  configured: [key: string]
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
    toast.error(t('realNameRequired'))
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
          {{ t('siliconFlow') }}
        </div>

        <div class="h-2 w-2 bg-green-500 dark:bg-green-400 rounded-full" />
        <div class="text-green-700 dark:text-green-500 text-base">
          {{ t('loggedIn') }}
        </div>

        <div class="flex-grow" />

        <span class="font-medium text-black">{{ siliconFlowStore.userInfo.data?.name }}</span>
      </CardTitle>
    </CardHeader>
    <CardContent>
      <div v-if="siliconFlowStore.userInfo.data" class="space-y-2">
        <div class="flex justify-between items-center text-sm">
          <span class="text-muted-foreground">{{ t('realname') }}</span>
          <button
            class="font-medium border-b border-dashed border-primary text-primary hover:bg-muted/50 transition-colors"
            @click="openRealNameDialog = true"
          >
            {{ siliconFlowStore.authed ? t('auth-success') : t('not-authenticated') }}
          </button>
        </div>
        <div v-if="siliconFlowStore.userInfo.data.phone" class="flex justify-between items-center text-sm">
          <span class="text-muted-foreground">{{ t('phone') }}</span>
          <span class="font-medium">{{ siliconFlowStore.userInfo.data.phone }}</span>
        </div>
        <div v-if="siliconFlowStore.userInfo.data.totalBalance" class="flex justify-between items-center text-sm">
          <span class="text-muted-foreground">{{ t('balance') }}</span>
          <span class="font-medium">{{ siliconFlowStore.userInfo.data.totalBalance }} {{ t('yuan') }}</span>
        </div>
      </div>
    </CardContent>
    <CardFooter class="pt-2 flex gap-2">
      <!-- <AlertDialog>
        <AlertDialogTrigger as-child>
          <Button variant="default" size="sm" class="h-8">
            {{ t('createKey') }}
          </Button>
        </AlertDialogTrigger>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{{ t('confirmCreateKey') }}</AlertDialogTitle>
            <AlertDialogDescription>
              {{ t('confirmCreateKeyDescription') }}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>{{ t('cancel') }}</AlertDialogCancel>
            <AlertDialogAction @click="handleCreateKey">
              {{ t('confirm') }}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog> -->
      <Button variant="default" size="sm" class="h-8" @click="emit('configured', siliconFlowStore.keyId!)">
        {{ t('confirm') }}
      </Button>
      <div class="flex-grow" />
      <Button variant="secondary" size="sm" class="h-8" @click="handleRecharge">
        {{ t('recharge') }}
      </Button>
      <Button variant="secondary" size="sm" class="h-8" @click="siliconFlowStore.logout">
        {{ t('logout') }}
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

<i18n lang="yaml">
zh-CN:
  siliconFlow: 硅基流动
  loggedIn: 已登录硅基流动
  realname: 实名
  authSuccess: 已认证
  notAuthenticated: 未认证
  phone: 手机
  balance: 余额
  yuan: 元
  createKey: 创建密钥
  confirmCreateKey: 确认创建密钥
  confirmCreateKeyDescription: 您确定要创建一个新的 API 密钥吗？这将为您的账户生成一个新的访问密钥。
  recharge: 充值
  logout: 退出登录
  realNameRequired: 请先完成实名认证
  cancel: 取消
  confirm: 确认
en-US:
  siliconFlow: SiliconFlow
  loggedIn: Logged in to SiliconFlow
  realname: Real Name
  authSuccess: Authenticated
  notAuthenticated: Not Authenticated
  phone: Phone
  balance: Balance
  yuan: Yuan
  createKey: Create Key
  confirmCreateKey: Confirm Create Key
  confirmCreateKeyDescription: Are you sure you want to create a new API key? This will generate a new access key for your account.
  recharge: Recharge
  logout: Logout
  realNameRequired: Please complete real-name authentication first
  cancel: Cancel
  confirm: Confirm
</i18n>
