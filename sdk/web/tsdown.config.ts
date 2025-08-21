import { readFileSync, writeFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { defineConfig } from 'tsdown'

export default defineConfig({
  entry: {
    index: 'src/index.ts',
  },
  sourcemap: true,
  clean: true,
  dts: {
    oxc: true,
  },
  platform: 'browser',
  onSuccess() {
    const readme = readFileSync(resolve(import.meta.dirname, '../../README.md'), 'utf8')
    const trimmed = readme.replace(/\n<!--DEV-->([\s\S]*?)<!--\/DEV-->\n/g, '')
    writeFileSync(resolve(import.meta.dirname, 'README.md'), trimmed)
  },
})
