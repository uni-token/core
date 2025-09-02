<script setup lang="ts">
import {
  Home,
  RotateCcw,
  Search,
  Zap,
} from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'

const router = useRouter()
const { t } = useI18n()

function goHome() {
  router.push('/')
}

function goBack() {
  router.back()
}
</script>

<template>
  <div class="p-6">
    <div class="min-h-[calc(100vh-8rem)] flex items-center justify-center">
      <div class="max-w-2xl w-full space-y-8 text-center">
        <div class="space-y-4">
          <div class="flex justify-center">
            <div class="relative">
              <div class="text-8xl font-bold text-muted-foreground/20 select-none">
                404
              </div>
              <div class="absolute inset-0 flex items-center justify-center">
                <div class="rounded-full bg-primary/10 p-6">
                  <Search class="h-12 w-12 text-primary" />
                </div>
              </div>
            </div>
          </div>

          <div class="space-y-2">
            <h1 class="text-3xl font-bold tracking-tight">
              {{ t('title') }}
            </h1>
            <p class="text-lg text-muted-foreground">
              {{ t('description') }}
            </p>
          </div>
        </div>

        <Card class="mx-auto max-w-md">
          <CardHeader>
            <CardTitle class="flex items-center justify-center gap-2 text-lg">
              <Zap class="h-5 w-5 text-primary" />
              {{ t('possibleReasons') }}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <ul class="text-sm text-muted-foreground space-y-2 text-left">
              <li class="flex items-start gap-2">
                <span class="text-primary mt-1">•</span>
                <span>{{ t('reason1') }}</span>
              </li>
              <li class="flex items-start gap-2">
                <span class="text-primary mt-1">•</span>
                <span>{{ t('reason2') }}</span>
              </li>
              <li class="flex items-start gap-2">
                <span class="text-primary mt-1">•</span>
                <span>{{ t('reason3') }}</span>
              </li>
              <li class="flex items-start gap-2">
                <span class="text-primary mt-1">•</span>
                <span>{{ t('reason4') }}</span>
              </li>
            </ul>
          </CardContent>
        </Card>

        <div class="flex flex-col sm:flex-row gap-4 justify-center">
          <Button class="flex items-center gap-2" @click="goHome">
            <Home class="h-4 w-4" />
            {{ t('goHome') }}
          </Button>
          <Button variant="outline" class="flex items-center gap-2" @click="goBack">
            <RotateCcw class="h-4 w-4" />
            {{ t('goBack') }}
          </Button>
        </div>

        <div class="text-sm text-muted-foreground">
          <p>
            {{ t('helpText') }}
            <button
              class="text-primary hover:underline font-medium"
              @click="goHome"
            >
              {{ t('goHome') }}
            </button>
            {{ t('reNavigate') }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<i18n lang="yaml">
en-US:
  title: Page Not Found
  description: Sorry, the page you are looking for does not exist or has been moved
  possibleReasons: Possible Reasons
  reason1: Incorrect URL address
  reason2: Page has been deleted or moved
  reason3: You don't have permission to access this page
  reason4: Website is under maintenance
  goHome: Go Home
  goBack: Go Back
  helpText: If the problem persists, please check if the URL is correct, or
  reNavigate: navigate again

zh-CN:
  title: 页面未找到
  description: 抱歉，您访问的页面不存在或已被移动
  possibleReasons: 可能的原因
  reason1: URL 地址输入错误
  reason2: 页面已被删除或移动
  reason3: 您没有访问此页面的权限
  reason4: 网站正在维护中
  goHome: 返回首页
  goBack: 返回上页
  helpText: 如果问题持续存在，请检查 URL 是否正确，或
  reNavigate: 重新导航
</i18n>
