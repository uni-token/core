import { readFileSync, writeFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { defineConfig } from 'tsdown'

export default defineConfig({
  entry: 'src/index.ts',
  outDir: 'dist',
  format: 'esm',
  sourcemap: true,
  clean: true,
  onSuccess() {
    const readme = readFileSync(resolve(import.meta.dirname, '../../README.md'), 'utf8')
    const trimmed = readme.replace(/\n<!--DEV-->([\s\S]*?)<!--\/DEV-->\n/g, '')
    writeFileSync(resolve(import.meta.dirname, 'README.md'), trimmed)
  },
})
