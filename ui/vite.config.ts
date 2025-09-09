import path from 'node:path'
import { fileURLToPath, URL } from 'node:url'
import VueI18n from '@intlify/unplugin-vue-i18n/vite'
import Tailwindcss from '@tailwindcss/vite'
import Vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'
import { version } from './package.json'

export default defineConfig({
  plugins: [
    Vue(),
    Tailwindcss(),
    VueI18n({
      runtimeOnly: true,
      compositionOnly: true,
      fullInstall: false,
      include: [
        path.resolve(import.meta.dirname, 'src/locales/**'),
        path.resolve(import.meta.dirname, 'src/**/*.vue'),
      ],
    }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  define: {
    'import.meta.env.VITE_BUILD_VERSION': JSON.stringify(version),
    'import.meta.env.VITE_RELEASE_TIMESTAMP': JSON.stringify(Date.now()),
  },
})
