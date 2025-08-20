<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue-sonner'
import Captcha from '@/components/Captcha.vue'
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
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Skeleton } from '@/components/ui/skeleton'
import { useSiliconFlowStore } from '@/stores'
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

async function handleCreateKey() {
  const result = await siliconFlowStore.createApiKeyAndApply()
  if (result) {
    emit('configured', result.id)
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
      <AlertDialog>
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
      </AlertDialog>
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
  <Card v-else-if="siliconFlowStore.userInfo">
    <CardHeader class="pb-4">
      <CardTitle class="text-lg">
        {{ t('siliconFlow.loginTitle') }}
      </CardTitle>
      <CardDescription class="text-sm">
        {{ t('siliconFlow.loginDescription') }}
      </CardDescription>
    </CardHeader>
    <CardContent class="space-y-4">
      <div class="flex rounded-md border border-input bg-background">
        <div class="flex items-center px-3 border-r border-input bg-muted/50 rounded-l-md">
          <span class="text-sm font-medium text-muted-foreground">+86</span>
        </div>
        <Input
          id="phone" v-model="siliconFlowStore.phoneNumber" :placeholder="t('siliconFlow.phoneNumber')" type="tel"
          class="border-0 rounded-l-none focus-visible:ring-0 focus-visible:ring-offset-0 h-10"
        />
      </div>
      <div class="flex rounded-md border border-input bg-background">
        <Input
          id="sms" v-model="siliconFlowStore.smsCode" :placeholder="t('siliconFlow.smsCode')" type="text" maxlength="6"
          class="w-fit flex-grow border-0 rounded-r-none focus-visible:ring-0 focus-visible:ring-offset-0 h-10"
        />
        <div class="border-l border-input">
          <Captcha
            :enabled="siliconFlowStore.phoneNumber.length > 0" :config="siliconFlowStore.captchaConfig"
            class="h-10 px-4 bg-muted/50 rounded-r-md border-0 text-xs text-primary hover:bg-muted/70 transition-colors disabled:opacity-50"
            @next="siliconFlowStore.sendSMS"
          />
        </div>
      </div>
    </CardContent>
    <CardFooter class="flex flex-col space-y-3 pt-3">
      <div class="flex items-center space-x-2 text-xs text-muted-foreground">
        <input id="agree" v-model="siliconFlowStore.agreed" type="checkbox" class="h-3 w-3 rounded border border-input">
        <label for="agree" class="flex items-center gap-1 cursor-pointer">
          <span>{{ t('siliconFlow.agreeToTerms') }}</span>
          <a
            href="https://docs.siliconflow.cn/docs/user-agreement" target="_blank"
            class="text-primary hover:underline"
          >{{ t('siliconFlow.userAgreement') }}</a>
          <span>{{ t('siliconFlow.and') }}</span>
          <a
            href="https://docs.siliconflow.cn/docs/privacy-policy" target="_blank"
            class="text-primary hover:underline"
          >{{ t('siliconFlow.privacyPolicy') }}</a>
        </label>
      </div>
      <Button class="w-full h-10" :disabled="!siliconFlowStore.canLogin" @click="siliconFlowStore.login">
        <span v-if="siliconFlowStore.isLoading">{{ t('siliconFlow.loggingIn') }}</span>
        <span v-else>{{ t('siliconFlow.registerLogin') }}</span>
      </Button>
      <div class="flex items-center space-x-2 text-xs text-muted-foreground">
        <input id="keep" v-model="siliconFlowStore.keepLogin" type="checkbox" class="h-3 w-3 rounded border border-input" checked>
        <label for="keep" class="cursor-pointer">{{ t('siliconFlow.keepLoggedIn30Days') }}</label>
      </div>
    </CardFooter>
  </Card>

  <!-- Loading state -->
  <Skeleton v-else class="h-48 w-full max-w-md mx-auto" />
</template>
