import UnoCss from 'unocss/vite'
import { defineConfig } from 'vitepress'

export default defineConfig({
  lang: 'en-US',
  title: 'UniToken',
  description: 'UniToken documentation.',
  themeConfig: {
    sidebar: [
      {
        text: 'Guide',
        items: [
          {
            text: 'Introduction',
            link: '/guide/',
          },
        ],
      },
      {
        text: 'SDK',
        items: [
          {
            text: 'Python SDK',
            link: '/sdk/python',
          },
          {
            text: 'Node.js SDK',
            link: '/sdk/nodejs',
          },
          {
            text: 'Go SDK',
            link: '/sdk/go',
          },
          {
            text: 'Browser SDK',
            link: '/sdk/browser',
          },
        ],
      },
    ],
    socialLinks: [
      { icon: 'github', link: 'https://github.com/uni-token/core' },
      { icon: 'discord', link: 'https://discord.gg/UCEu5gTEHg' },
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
