<script setup lang="ts">
import { RefreshCw } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
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
    error.value = err instanceof Error ? err.message : t('common.error')
  }
  finally {
    loading.value = false
  }
}

function refreshStats() {
  loadStats()
}

function formatNumber(num: number) {
  if (num >= 1000000) {
    return `${(num / 1000000).toFixed(1)}M`
  }
  else if (num >= 1000) {
    return `${(num / 1000).toFixed(1)}K`
  }
  else {
    return num.toString()
  }
}

function formatTime(timestamp: string) {
  try {
    return new Date(timestamp).toLocaleString('zh-CN')
  }
  catch {
    return t('common.invalidTime')
  }
}

onMounted(() => {
  loadStats()
})
</script>

<template>
  <div class="p-6">
    <div class="space-y-6">
      <div class="flex items-center justify-between">
        <h2 class="text-2xl font-bold">
          {{ t('usage.title') }}
        </h2>
        <div class="flex items-center gap-2">
          <Select v-model="selectedDays" @update:model-value="loadStats">
            <SelectTrigger class="w-[180px]">
              <SelectValue>
                {{ selectedDays === '7' ? t('usage.last7Days') : selectedDays === '30' ? t('usage.last30Days') : t('usage.last90Days') }}
              </SelectValue>
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="7">
                {{ t('usage.last7Days') }}
              </SelectItem>
              <SelectItem value="30">
                {{ t('usage.last30Days') }}
              </SelectItem>
              <SelectItem value="90">
                {{ t('usage.last90Days') }}
              </SelectItem>
            </SelectContent>
          </Select>
          <Button variant="outline" :disabled="loading" @click="refreshStats">
            <RefreshCw class="mr-2 h-4 w-4" :class="{ 'animate-spin': loading }" />
            {{ t('usage.refresh') }}
          </Button>
        </div>
      </div>

      <!-- Overview Cards -->
      <div v-if="stats" class="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader class="pb-2">
            <CardTitle class="text-sm font-medium text-muted-foreground">
              {{ t('usage.totalTokens') }}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">
              {{ formatNumber(stats.totalTokens) }}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader class="pb-2">
            <CardTitle class="text-sm font-medium text-muted-foreground">
              {{ t('usage.totalCost') }}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">
              ¥{{ stats.totalCost.toFixed(2) }}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader class="pb-2">
            <CardTitle class="text-sm font-medium text-muted-foreground">
              {{ t('usage.totalRequests') }}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">
              {{ formatNumber(stats.totalRequests) }}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader class="pb-2">
            <CardTitle class="text-sm font-medium text-muted-foreground">
              {{ t('usage.averageTokensPerRequest') }}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">
              {{ stats.totalRequests > 0 ? Math.round(stats.totalTokens / stats.totalRequests) : 0 }}
            </div>
          </CardContent>
        </Card>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="space-y-4">
        <div class="grid gap-4 md:grid-cols-4">
          <Skeleton v-for="i in 4" :key="i" class="h-24" />
        </div>
        <div class="grid gap-4 md:grid-cols-3">
          <Skeleton v-for="i in 3" :key="i" class="h-64" />
        </div>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="rounded-lg border border-red-200 bg-red-50 p-4">
        <p class="text-red-800">
          {{ t('usage.loadFailed') }}: {{ error }}
        </p>
      </div>

      <!-- Statistics Tables -->
      <div v-else-if="stats" class="grid gap-6 md:grid-cols-3">
        <!-- By Application -->
        <Card>
          <CardHeader>
            <CardTitle class="text-lg">
              {{ t('usage.byApp') }}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div class="space-y-2">
              <div v-if="Object.keys(stats.byApp).length === 0" class="text-sm text-muted-foreground text-center py-4">
                {{ t('usage.noData') }}
              </div>
              <div
                v-for="(app, appId) in stats.byApp" v-else :key="appId"
                class="flex justify-between items-center p-2 rounded border"
              >
                <div>
                  <div class="font-medium text-sm">
                    {{ app.appName || t('usage.unknownApp') }}
                  </div>
                  <div class="text-xs text-muted-foreground">
                    {{ app.requestCount }} {{ t('usage.requests') }}
                  </div>
                </div>
                <div class="text-right">
                  <div class="font-mono text-sm">
                    {{ formatNumber(app.totalTokens) }}
                  </div>
                  <div class="text-xs text-muted-foreground">
                    ¥{{ app.totalCost.toFixed(2) }}
                  </div>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        <!-- By Provider -->
        <Card>
          <CardHeader>
            <CardTitle class="text-lg">
              {{ t('usage.byProvider') }}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div class="space-y-2">
              <div v-if="Object.keys(stats.byKey).length === 0" class="text-sm text-muted-foreground text-center py-4">
                {{ t('usage.noData') }}
              </div>
              <div
                v-for="(provider, name) in stats.byKey" v-else :key="name"
                class="flex justify-between items-center p-2 rounded border"
              >
                <div>
                  <div class="font-medium text-sm">
                    {{ name }}
                  </div>
                  <div class="text-xs text-muted-foreground">
                    {{ provider.requestCount }} {{ t('usage.requests') }}
                  </div>
                </div>
                <div class="text-right">
                  <div class="font-mono text-sm">
                    {{ formatNumber(provider.totalTokens) }}
                  </div>
                  <div class="text-xs text-muted-foreground">
                    ¥{{ provider.totalCost.toFixed(2) }}
                  </div>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        <!-- By Model -->
        <Card>
          <CardHeader>
            <CardTitle class="text-lg">
              {{ t('usage.byModel') }}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div class="space-y-2">
              <div v-if="Object.keys(stats.byModel).length === 0" class="text-sm text-muted-foreground text-center py-4">
                {{ t('usage.noData') }}
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
                    {{ model.key }} · {{ model.requestCount }} {{ t('usage.times') }}
                  </div>
                </div>
                <div class="text-right">
                  <div class="font-mono text-sm">
                    {{ formatNumber(model.totalTokens) }}
                  </div>
                  <div class="text-xs text-muted-foreground">
                    ¥{{ model.totalCost.toFixed(2) }}
                  </div>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      <!-- Recent Usage Records -->
      <Card v-if="stats">
        <CardHeader>
          <CardTitle class="text-lg">
            {{ t('usage.recentUsageRecords') }}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div v-if="stats.recentUsages.length === 0" class="text-center py-8 text-muted-foreground">
            {{ t('usage.noUsageRecords') }}
          </div>
          <div v-else class="space-y-2 max-h-96 overflow-y-auto">
            <div
              v-for="usage in stats.recentUsages" :key="usage.id"
              class="flex items-center justify-between p-3 rounded border hover:bg-muted/50"
            >
              <div class="flex-1">
                <div class="flex items-center gap-2">
                  <span class="font-medium text-sm">{{ usage.appName || t('usage.unknownApp') }}</span>
                  <span class="text-xs px-2 py-1 bg-muted rounded">{{ usage.key }}</span>
                  <span class="text-xs text-muted-foreground">{{ usage.model }}</span>
                </div>
                <div class="text-xs text-muted-foreground mt-1">
                  {{ formatTime(usage.timestamp) }} · {{ usage.endpoint }}
                </div>
              </div>
              <div class="text-right">
                <div class="font-mono text-sm">
                  {{ formatNumber(usage.totalTokens) }} {{ t('usage.tokensUnit') }}
                </div>
                <div class="text-xs text-muted-foreground">
                  ¥{{ usage.cost.toFixed(4) }}
                  <span :class="usage.status === 'success' ? 'text-green-600' : 'text-red-600'" class="ml-1">
                    {{ usage.status === 'success' ? '✓' : '✗' }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
