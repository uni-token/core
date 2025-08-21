<script setup>
import { Analytics } from '@vercel/analytics/vue'
import DefaultTheme from 'vitepress/theme'

const { Layout } = DefaultTheme
</script>

<template>
  <Layout>
    <template #aside-outline-before>
      <Analytics />
    </template>
  </Layout>
</template>
