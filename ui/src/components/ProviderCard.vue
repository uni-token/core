<script setup lang="ts">
import type { Provider } from '@/lib/providers'
import { ExternalLinkIcon } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue-sonner'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import ProviderPaymentDialog from './ProviderPaymentDialog.vue'
import ProviderRealNameDialog from './ProviderRealNameDialog.vue'

const props = defineProps<{
  provider: Provider
}>()

const emit = defineEmits<{
  configured: [key: string]
}>()

const { t } = useI18n()

const userInfo = computed(() => props.provider.user)
const openPaymentDialog = ref(false)
const openRealNameDialog = ref(false)

// Check login status on component mount
onMounted(async () => {
  await props.provider.refreshUser()
})

function handleRecharge() {
  if ('websiteURL' in props.provider.payment!) {
    window.open(props.provider.payment!.websiteURL, '_blank', 'width=600,height=600,noopener=yes,noreferrer=yes')
    return
  }

  if (!userInfo.value?.verified) {
    toast.error(t('realNameRequired'))
    openRealNameDialog.value = true
  }
  else {
    openPaymentDialog.value = true
  }
}
</script>

<template>
  <component :is="provider.Login" v-if="userInfo === null" />
  <Card v-else>
    <CardHeader>
      <CardTitle class="flex items-center gap-2">
        <div class="mr-4 text-lg">
          {{ provider.name }}
        </div>

        <div class="h-2 w-2 bg-green-500 dark:bg-green-400 rounded-full" />
        <div class="text-green-700 dark:text-green-500 text-base">
          {{ t('loggedIn', [provider.name]) }}
        </div>
      </CardTitle>
    </CardHeader>
    <CardContent class="flex-grow">
      <div class="space-y-2">
        <div v-if="userInfo?.name" class="flex justify-between items-center text-sm">
          <span class="text-muted-foreground">{{ t('userName') }}</span>
          <span class="font-medium text-ellipsis max-w-24 text-nowrap overflow-hidden" :title="userInfo.name">{{ userInfo.name }}</span>
        </div>
        <div v-if="userInfo?.verified != null" class="flex justify-between items-center text-sm">
          <span class="text-muted-foreground">{{ t('realname') }}</span>
          <button
            class="font-medium border-b border-dashed border-primary text-primary hover:bg-muted/50 transition-colors"
            @click="openRealNameDialog = true"
          >
            {{ userInfo.verified ? t('verified') : t('unverified') }}
          </button>
        </div>
        <div v-if="userInfo?.phone" class="flex justify-between items-center text-sm">
          <span class="text-muted-foreground">{{ t('phone') }}</span>
          <span class="font-medium">{{ userInfo.phone }}</span>
        </div>
        <div v-if="userInfo?.balance != null" class="flex justify-between items-center text-sm">
          <span class="text-muted-foreground">{{ t('balance') }}</span>
          <span class="font-medium">
            {{ userInfo.balance.currency === 'USD' ? '$' : '' }}{{ userInfo.balance.amount }}
            {{ userInfo.balance.currency === 'CNY' ? t('yuan') : userInfo.balance.currency === 'USD' ? '' : `(${userInfo.balance.currency})` }}
          </span>
        </div>
      </div>
    </CardContent>
    <CardFooter class="pt-2 flex gap-2">
      <Button variant="default" size="sm" class="h-8" @click="emit('configured', provider.id!)">
        {{ t('confirm') }}
      </Button>
      <div class="flex-grow" />
      <Button v-if="!!provider.payment" variant="secondary" size="sm" class="h-8" @click="handleRecharge">
        {{ t('recharge') }}
        <ExternalLinkIcon v-if="'websiteURL' in provider.payment" />
      </Button>
      <Button variant="secondary" size="sm" class="h-8" @click="provider.logout">
        {{ t('logout') }}
      </Button>
    </CardFooter>

    <ProviderRealNameDialog v-if="!!provider.verification" v-model="openRealNameDialog" :provider="provider" />
    <ProviderPaymentDialog v-if="!!provider.payment" v-model="openPaymentDialog" :provider="provider" />
  </Card>
</template>

<i18n lang="yaml">
zh-CN:
  loggedIn: 已登录
  realname: 实名认证
  verified: 已认证
  unverified: 未认证
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
  userName: 用户名
en-US:
  loggedIn: Configured
  realname: Legal Name
  verified: Verified
  unverified: Unverified
  phone: Phone
  balance: Balance
  yuan: CNY
  createKey: Create Key
  confirmCreateKey: Confirm Create Key
  confirmCreateKeyDescription: Are you sure you want to create a new API key? This will generate a new access key for your account.
  recharge: Recharge
  logout: Logout
  realNameRequired: Please complete real-name authentication first
  cancel: Cancel
  confirm: Confirm
  userName: User Name
</i18n>
