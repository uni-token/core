<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import Captcha from '@/components/Captcha.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { useSiliconFlowStore } from '@/stores'

const { t } = useI18n()
const siliconFlowStore = useSiliconFlowStore()
</script>

<template>
  <Card>
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
</template>
