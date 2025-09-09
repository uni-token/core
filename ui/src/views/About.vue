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
  return import.meta.env.VITE_BUILD_VERSION || t('unknown')
})

const releaseDate = computed(() => {
  const timestamp = import.meta.env.VITE_RELEASE_TIMESTAMP
  return timestamp ? new Date(Number(timestamp)).toLocaleDateString() : t('unknown')
})
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
          <CardTitle class="flex items-center gap-2">
            <Zap class="h-5 w-5 text-primary" />
            {{ t('introduction') }}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            <p class="text-muted-foreground">
              {{ t('description') }}
            </p>
            <p class="text-muted-foreground">
              {{ t('features') }}
            </p>
            <ul class="list-disc list-inside space-y-2 text-muted-foreground ml-4">
              <li>{{ t('feature1') }}</li>
              <li>{{ t('feature2') }}</li>
              <li>{{ t('feature3') }}</li>
            </ul>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle class="flex items-center gap-2">
            <InfoIcon class="h-5 w-5 text-primary" />
            {{ t('versionInfo') }}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-12">
            <div class="space-y-2">
              <div class="flex justify-between">
                <span class="font-medium">{{ t('version') }}: </span>
                <span class="text-muted-foreground">{{ buildVersion }}</span>
              </div>
            </div>
            <div class="space-y-2">
              <div class="flex justify-between">
                <span class="font-medium">{{ t('releaseDate') }}: </span>
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
            {{ t('contact') }}
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
              <span class="text-sm">{{ t('website') }}: </span>
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
            {{ t('license') }}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <p class="text-sm text-muted-foreground">
            {{ t('licenseText') }}
            <a
              href="https://github.com/uni-token/core/blob/main/LICENSE"
              class="text-primary hover:underline text-sm"
            >LICENSE</a>
            {{ t('licenseFile') }}
          </p>
        </CardContent>
      </Card>
    </div>
  </div>
</template>

<i18n lang="yaml">
zh-CN:
  title: 关于
  introduction: 关于 UniToken
  description: UniToken 是一个高效的 AI Token 聚合服务，为开发者提供统一的 API 接口来访问各种 AI 模型。
  features: 主要特性：
  feature1: 统一的 API 接口访问多个 AI 提供商
  feature2: 智能负载均衡和故障转移
  feature3: 实时用量统计和成本分析
  versionInfo: 版本信息
  version: 版本
  releaseDate: 发布日期
  unknown: 未知
  contact: 联系方式
  website: 网站
  license: 许可证
  licenseText: 本项目使用 MIT 许可证开源。
  licenseFile: 查看完整的许可证文件
en-US:
  title: About
  introduction: About UniToken
  description: UniToken is an efficient AI Token aggregation service that provides developers with a unified API interface to access various AI models.
  features: Key Features：
  feature1: Unified API interface to access multiple AI providers
  feature2: Intelligent load balancing and failover
  feature3: Real-time usage statistics and cost analysis
  versionInfo: Version Information
  version: Version
  releaseDate: Release Date
  unknown: Unknown
  contact: Contact
  website: Website
  license: License
  licenseText: This project is open sourced under the MIT license.
  licenseFile: View complete license file
</i18n>
