<script setup lang="ts">
import { Layers, RefreshCw } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Skeleton } from '@/components/ui/skeleton'
import { useServiceStore } from '@/stores'

interface UsageStats {
  totalTokens: number
  totalCost: number
  totalRequests: number
  byApp: Record<string, {
    appName: string
    totalTokens: number
    totalCost: number
    requestCount: number
  }>
  byKey: Record<string, {
    totalTokens: number
    totalCost: number
    requestCount: number
  }>
  byModel: Record<string, {
    key: string
    totalTokens: number
    totalCost: number
    requestCount: number
  }>
  recentUsages: Array<{
    id: string
    appName: string
    key: string
    model: string
    totalTokens: number
    cost: number
    timestamp: string
    endpoint: string
    status: string
  }>
}

const { t } = useI18n()
const { fetch } = useServiceStore()
const stats = ref<UsageStats | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)
const selectedDays = ref('30')

async function loadStats() {
  loading.value = true
  error.value = null

  try {
    const response = await fetch(`usage/stats?days=${selectedDays.value}`)
    if (response.ok) {
      const data = await response.json()
      stats.value = data.data
    }
    else {
      error.value = `HTTP ${response.status}: ${response.statusText}`
    }
  }
  catch (err) {
    error.value = err instanceof Error ? err.message : t('error')
  }
  finally {
    loading.value = false
  }
}

function refreshStats() {
  loadStats()
}

function formatNumber(num: number) {
  if (num >= 1000000)
    return `${(num / 1000000).toFixed(1)}M`
  if (num >= 1000)
    return `${(num / 1000).toFixed(1)}K`
  return num.toString()
}

function formatTime(timestamp: string) {
  try {
    return new Date(timestamp).toLocaleString('zh-CN')
  }
  catch {
    return t('invalidTime')
  }
}

function formatCurrency(num: number) {
  return `¥${num.toFixed(2)}`
}

function getPercentage(part: number, total: number) {
  if (!total)
    return 0
  return Math.round((part / total) * 100)
}

const totalTokens = computed(() => stats.value?.totalTokens ?? 0)
const totalCost = computed(() => stats.value?.totalCost ?? 0)
const totalRequests = computed(() => stats.value?.totalRequests ?? 0)
const avgTokens = computed(() => {
  if (!stats.value?.totalRequests)
    return 0
  return Math.round(stats.value.totalTokens / stats.value.totalRequests)
})

const appStats = computed(() => {
  if (!stats.value)
    return []
  return Object.entries(stats.value.byApp)
    .map(([id, app]) => ({
      id,
      name: app.appName || t('unknownApp'),
      tokens: app.totalTokens,
      cost: app.totalCost,
      requests: app.requestCount,
    }))
    .sort((a, b) => b.tokens - a.tokens)
})

const providerStats = computed(() => {
  if (!stats.value)
    return []
  return Object.entries(stats.value.byKey)
    .map(([name, provider]) => ({
      name,
      tokens: provider.totalTokens,
      cost: provider.totalCost,
      requests: provider.requestCount,
    }))
    .sort((a, b) => b.tokens - a.tokens)
})

onMounted(() => {
  loadStats()
})
</script>

<template>
  <div class="p-6">
    <div class="space-y-6">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <p class="text-xs uppercase tracking-wide text-muted-foreground">
            {{ t('title') }}
          </p>
          <h2 class="text-2xl font-bold">
            {{ t('overview') }}
          </h2>
        </div>
        <div class="flex items-center gap-2">
          <Select v-model="selectedDays" @update:model-value="loadStats">
            <SelectTrigger class="w-[150px]">
              <SelectValue>
                {{ selectedDays === '7' ? t('last7Days') : selectedDays === '30' ? t('last30Days') : t('last90Days') }}
              </SelectValue>
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="7">
                {{ t('last7Days') }}
              </SelectItem>
              <SelectItem value="30">
                {{ t('last30Days') }}
              </SelectItem>
              <SelectItem value="90">
                {{ t('last90Days') }}
              </SelectItem>
            </SelectContent>
          </Select>
          <Button variant="outline" size="icon" :disabled="loading" @click="refreshStats">
            <RefreshCw :class="['h-4 w-4', loading ? 'animate-spin' : '']" />
          </Button>
        </div>
      </div>

      <div v-if="loading && !stats" class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <Skeleton v-for="i in 3" :key="i" class="h-32 rounded-2xl" />
      </div>

      <div v-else-if="error" class="p-6 bg-destructive/5 border border-destructive/50 rounded-2xl">
        <p class="text-destructive font-medium">{{ t('loadFailed') }}</p>
        <p class="text-sm text-destructive/80">{{ error }}</p>
      </div>

      <div v-else-if="stats" class="space-y-6">
        <Card class="flex flex-col gap-6 rounded-xl border py-6">
          <CardContent class="p-6 space-y-6 text-foreground">
            <div class="flex flex-wrap items-center justify-between gap-4">
              <div>
                <p class="text-sm text-muted-foreground">{{ t('totalTokens') }}</p>
                <h3 class="text-4xl font-semibold text-foreground">
                  {{ formatNumber(totalTokens) }}
                  <span class="text-xl font-normal text-muted-foreground">{{ t('tokensUnit') }}</span>
                </h3>
              </div>
              <span class="px-4 py-1 rounded-full bg-muted text-sm font-medium text-foreground">
                {{ selectedDays === '7' ? t('last7Days') : selectedDays === '30' ? t('last30Days') : t('last90Days') }}
              </span>
            </div>
            <div class="grid gap-4 md:grid-cols-3">
              <div class="rounded-2xl border bg-muted/30 p-4">
                <p class="text-sm text-muted-foreground">{{ t('totalCost') }}</p>
                <p class="text-2xl font-semibold mt-1 text-foreground">{{ formatCurrency(totalCost) }}</p>
              </div>
              <div class="rounded-2xl border bg-muted/30 p-4">
                <p class="text-sm text-muted-foreground">{{ t('totalRequests') }}</p>
                <p class="text-2xl font-semibold mt-1 text-foreground">{{ formatNumber(totalRequests) }}</p>
              </div>
              <div class="rounded-2xl border bg-muted/30 p-4">
                <p class="text-sm text-muted-foreground">{{ t('averageTokensPerRequest') }}</p>
                <p class="text-2xl font-semibold mt-1 text-foreground">{{ formatNumber(avgTokens) }}</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <div class="grid gap-6 lg:grid-cols-2">
          <Card>
            <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
              <div>
                <p class="text-sm text-muted-foreground uppercase tracking-wide">
                  {{ t('byApp') }}
                </p>
                <CardTitle class="text-xl font-bold text-foreground">
                  {{ t('appUsage') }}
                </CardTitle>
              </div>
              <div class="size-9 rounded-full bg-primary/10 text-primary flex items-center justify-center font-semibold">
                {{ appStats.length }}
              </div>
            </CardHeader>
            <CardContent>
              <div v-if="appStats.length === 0" class="text-center py-8 text-muted-foreground">
                {{ t('noData') }}
              </div>
              <div v-else class="space-y-4">
                <div
                  v-for="app in appStats"
                  :key="app.id"
                  class="rounded-xl border bg-muted/30 p-3"
                >
                  <div class="flex items-center justify-between">
                    <div>
                      <p class="font-medium text-sm">{{ app.name }}</p>
                      <p class="text-xs text-muted-foreground">
                        {{ app.requests }} {{ t('requests') }} · {{ formatCurrency(app.cost) }}
                      </p>
                    </div>
                    <p class="font-mono text-sm">{{ formatNumber(app.tokens) }}</p>
                  </div>
                  <div class="mt-3 h-2 rounded-full bg-muted overflow-hidden">
                    <div
                      class="h-full rounded-full bg-primary"
                      :style="{ width: `${getPercentage(app.tokens, totalTokens)}%` }"
                    />
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
              <div>
                <p class="text-sm text-muted-foreground uppercase tracking-wide">
                  {{ t('byProvider') }}
                </p>
                <CardTitle class="text-xl font-bold text-foreground">
                  {{ t('providerUsage') }}
                </CardTitle>
              </div>
              <div class="size-9 rounded-full bg-muted flex items-center justify-center">
                <Layers class="h-4 w-4 text-muted-foreground" />
              </div>
            </CardHeader>
            <CardContent>
              <div v-if="providerStats.length === 0" class="text-center py-8 text-muted-foreground">
                {{ t('noData') }}
              </div>
              <div v-else class="space-y-4">
                <div
                  v-for="provider in providerStats"
                  :key="provider.name"
                  class="rounded-xl border p-3 hover:border-primary/40 transition"
                >
                  <div class="flex items-center justify-between">
                    <div>
                      <p class="font-medium text-sm">{{ provider.name }}</p>
                      <p class="text-xs text-muted-foreground">
                        {{ provider.requests }} {{ t('requests') }} · {{ formatCurrency(provider.cost) }}
                      </p>
                    </div>
                    <p class="font-mono text-sm">{{ formatNumber(provider.tokens) }}</p>
                  </div>
                  <div class="mt-3 h-2 rounded-full bg-muted overflow-hidden">
                    <div
                      class="h-full rounded-full bg-slate-500"
                      :style="{ width: `${getPercentage(provider.tokens, totalTokens)}%` }"
                    />
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        <Card>
          <CardHeader>
            <CardTitle class="text-lg">
              {{ t('byModel') }}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div class="space-y-2">
              <div v-if="Object.keys(stats.byModel).length === 0" class="text-sm text-muted-foreground text-center py-4">
                {{ t('noData') }}
              </div>
              <div
                v-for="(model, name) in stats.byModel" v-else :key="name"
                class="flex justify-between items-center p-2 rounded border"
              >
                <div>
                  <div class="font-medium text-sm">
                    {{ name.split('/')[1] || name }}
                  </div>
                  <div class="text-xs text-muted-foreground">
                    {{ model.key }} · {{ model.requestCount }} {{ t('times') }}
                  </div>
                </div>
                <div class="text-right">
                  <div class="font-mono text-sm">
                    {{ formatNumber(model.totalTokens) }}
                  </div>
                  <div class="text-xs text-muted-foreground">
                    {{ formatCurrency(model.totalCost) }}
                  </div>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle class="text-lg">
              {{ t('recentUsageRecords') }}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div v-if="stats.recentUsages.length === 0" class="text-center py-8 text-muted-foreground">
              {{ t('noUsageRecords') }}
            </div>
            <div v-else class="space-y-2 max-h-96 overflow-y-auto pr-1">
              <div
                v-for="usage in stats.recentUsages"
                :key="usage.id"
                class="flex items-center justify-between gap-3 p-4 rounded-xl border bg-muted/30 hover:bg-muted/60 transition"
              >
                <div class="flex-1">
                  <div class="flex flex-wrap items-center gap-2">
                    <span class="font-medium text-sm">{{ usage.appName || t('unknownApp') }}</span>
                    <span class="text-xs px-2 py-1 bg-secondary rounded">{{ usage.key }}</span>
                    <span class="text-xs text-muted-foreground">{{ usage.model }}</span>
                  </div>
                  <div class="text-xs text-muted-foreground mt-1">
                    {{ formatTime(usage.timestamp) }} · {{ usage.endpoint }}
                  </div>
                </div>
                <div class="text-right">
                  <div class="font-mono text-sm">
                    {{ formatNumber(usage.totalTokens) }} {{ t('tokensUnit') }}
                  </div>
                  <div class="text-xs text-muted-foreground flex items-center justify-end gap-2">
                    {{ formatCurrency(usage.cost) }}
                    <span
                      :class="usage.status === 'success' ? 'text-green-600' : 'text-red-500'"
                      class="text-xs font-medium"
                    >
                      {{ usage.status === 'success' ? t('statusSuccess') : t('statusFailed') }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  </div>
</template>

<i18n lang="yaml">
zh-CN:
  title: Token 用量统计
  overview: 用量总览
  last7Days: 最近7天
  last30Days: 最近30天
  last90Days: 最近90天
  refresh: 刷新
  totalTokens: 总Token消耗
  totalCost: 总费用
  totalRequests: 总请求数
  averageTokensPerRequest: 平均每请求Token
  loadFailed: 加载用量统计失败
  byApp: 按应用统计
  byProvider: 按提供商统计
  appUsage: 应用用量
  providerUsage: 提供商用量
  byModel: 按模型统计
  noData: 暂无数据
  unknownApp: 未知应用
  requests: 次请求
  times: 次
  recentUsageRecords: 最近使用记录
  noUsageRecords: 暂无使用记录
  tokensUnit: tokens
  error: 未知错误
  invalidTime: 无效时间
  statusSuccess: 成功
  statusFailed: 失败
en-US:
  title: Token Usage Statistics
  overview: Overview
  last7Days: Last 7 Days
  last30Days: Last 30 Days
  last90Days: Last 90 Days
  refresh: Refresh
  totalTokens: Total Token Consumption
  totalCost: Total Cost
  totalRequests: Total Requests
  averageTokensPerRequest: Average Tokens per Request
  loadFailed: Failed to load usage statistics
  byApp: By Application
  byProvider: By Provider
  appUsage: App Usage
  providerUsage: Provider Usage
  byModel: By Model
  noData: No data available
  unknownApp: Unknown App
  requests: requests
  times: times
  recentUsageRecords: Recent Usage Records
  noUsageRecords: No usage records
  tokensUnit: tokens
  error: Unknown error
  invalidTime: Invalid time
  statusSuccess: Success
  statusFailed: Failed
</i18n>
