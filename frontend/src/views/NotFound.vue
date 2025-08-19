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
              {{ t('notFound.title') }}
            </h1>
            <p class="text-lg text-muted-foreground">
              {{ t('notFound.description') }}
            </p>
          </div>
        </div>

        <Card class="mx-auto max-w-md">
          <CardHeader>
            <CardTitle class="flex items-center justify-center gap-2 text-lg">
              <Zap class="h-5 w-5 text-primary" />
              {{ t('notFound.possibleReasons') }}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <ul class="text-sm text-muted-foreground space-y-2 text-left">
              <li class="flex items-start gap-2">
                <span class="text-primary mt-1">•</span>
                <span>{{ t('notFound.reason1') }}</span>
              </li>
              <li class="flex items-start gap-2">
                <span class="text-primary mt-1">•</span>
                <span>{{ t('notFound.reason2') }}</span>
              </li>
              <li class="flex items-start gap-2">
                <span class="text-primary mt-1">•</span>
                <span>{{ t('notFound.reason3') }}</span>
              </li>
              <li class="flex items-start gap-2">
                <span class="text-primary mt-1">•</span>
                <span>{{ t('notFound.reason4') }}</span>
              </li>
            </ul>
          </CardContent>
        </Card>

        <div class="flex flex-col sm:flex-row gap-4 justify-center">
          <Button class="flex items-center gap-2" @click="goHome">
            <Home class="h-4 w-4" />
            {{ t('notFound.goHome') }}
          </Button>
          <Button variant="outline" class="flex items-center gap-2" @click="goBack">
            <RotateCcw class="h-4 w-4" />
            {{ t('notFound.goBack') }}
          </Button>
        </div>

        <div class="text-sm text-muted-foreground">
          <p>
            {{ t('notFound.helpText') }}
            <button
              class="text-primary hover:underline font-medium"
              @click="goHome"
            >
              {{ t('notFound.goHome') }}
            </button>
            {{ t('notFound.reNavigate') }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
