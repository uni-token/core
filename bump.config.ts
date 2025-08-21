import { readFileSync, writeFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { parse, stringify } from '@std/toml'
import { defineConfig } from 'bumpp'
import { $ } from 'zx'

export default defineConfig({
  files: ['sdk/*/package.json'],
  all: true,
  async execute(config) {
    const projectToml = resolve(import.meta.dirname, './sdk/python/pyproject.toml')
    const projectConfig = parse(readFileSync(projectToml, 'utf-8')) as any
    const version = config.results.newVersion
    projectConfig.project.version = version
    writeFileSync(projectToml, stringify(projectConfig, {
      keyAlignment: true,
    }), 'utf-8')
    await $`pnpm lint --fix`
  },
})
