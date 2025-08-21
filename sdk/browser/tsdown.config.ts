import { defineConfig } from 'tsdown'

export default defineConfig({
  sourcemap: true,
  clean: true,
  dts: {
    oxc: true,
  },
  platform: 'browser',
})
