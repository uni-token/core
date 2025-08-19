<script setup lang="ts">
import {
  FileText,
  Github,
  Globe,
  InfoIcon,
  Mail,
  Zap,
} from 'lucide-vue-next'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'

const { t } = useI18n()

const buildVersion = computed(() => {
  return import.meta.env.VITE_BUILD_VERSION || t('about.unknown')
})

const releaseDate = computed(() => {
  const timestamp = import.meta.env.VITE_RELEASE_TIMESTAMP
  return timestamp ? new Date(Number(timestamp)).toLocaleDateString() : t('about.unknown')
})
</script>

<template>
  <div class="p-6">
    <div class="space-y-6 max-w-3xl">
      <div class="flex items-center justify-between">
        <h2 class="text-2xl font-bold">
          {{ t('about.title') }}
        </h2>
      </div>

      <Card>
        <CardHeader>
          <CardTitle class="flex items-center gap-2">
            <Zap class="h-5 w-5 text-primary" />
            {{ t('about.introduction') }}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            <p class="text-muted-foreground">
              {{ t('about.description') }}
            </p>
            <p class="text-muted-foreground">
              {{ t('about.features') }}
            </p>
            <ul class="list-disc list-inside space-y-2 text-muted-foreground ml-4">
              <li>{{ t('about.feature1') }}</li>
              <li>{{ t('about.feature2') }}</li>
              <li>{{ t('about.feature3') }}</li>
            </ul>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle class="flex items-center gap-2">
            <InfoIcon class="h-5 w-5 text-primary" />
            {{ t('about.versionInfo') }}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-12">
            <div class="space-y-2">
              <div class="flex justify-between">
                <span class="font-medium">{{ t('about.version') }}: </span>
                <span class="text-muted-foreground">{{ buildVersion }}</span>
              </div>
            </div>
            <div class="space-y-2">
              <div class="flex justify-between">
                <span class="font-medium">{{ t('about.releaseDate') }}: </span>
                <span class="text-muted-foreground">{{ releaseDate }}</span>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle class="flex items-center gap-2">
            <Mail class="h-5 w-5 text-primary" />
            {{ t('about.contact') }}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            <div class="flex items-center gap-3">
              <Github class="h-4 w-4 text-muted-foreground" />
              <span class="text-sm">GitHub: </span>
              <a href="https://github.com/uni-token" target="_blank" class="text-primary hover:underline text-sm">
                github.com/uni-token
              </a>
            </div>
            <div class="flex items-center gap-3">
              <Globe class="h-4 w-4 text-muted-foreground" />
              <span class="text-sm">{{ t('about.website') }}: </span>
              <a href="https://uni-token.app" class="text-primary hover:underline text-sm">
                uni-token.app
              </a>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle class="flex items-center gap-2">
            <FileText class="h-5 w-5 text-primary" />
            {{ t('about.license') }}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <p class="text-sm text-muted-foreground">
            {{ t('about.licenseText') }}
            <a
              href="https://github.com/uni-token/core/blob/main/LICENSE"
              class="text-primary hover:underline text-sm"
            >LICENSE</a>
            {{ t('about.licenseFile') }}
          </p>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
