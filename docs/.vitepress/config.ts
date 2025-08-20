import UnoCss from 'unocss/vite'
import { defineConfig } from 'vitepress'

export default defineConfig({
  lang: 'en-US',
  title: 'UniToken',
  description: 'UniToken documentation.',
  themeConfig: {
    sidebar: [
      {
        text: 'Introduction',
        link: '/guide/',
      },
      {
        text: 'SDK',
        items: [
          {
            text: 'Node.js SDK',
            link: '/sdk/nodejs',
          },
          {
            text: 'Python SDK',
            link: '/sdk/python',
          },
        ],
      },
    ],
  },
  head: [
    ['link', { rel: 'icon', href: '/logo-light.png' }],
    ['link', { rel: 'icon', href: '/logo-dark.png', media: '(prefers-color-scheme: dark)' }],
  ],
  vite: {
    plugins: [
      UnoCss(),
    ],
  },
})
