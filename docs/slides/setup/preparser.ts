import { existsSync } from 'node:fs'
import { resolve } from 'node:path'
import { definePreparserSetup } from '@slidev/types'

export default definePreparserSetup(({ filepath }) => {
  return [
    {
      transformRawLines(lines) {
        for (const i in lines) {
          const match = lines[i].match(/^@CONDITIONAL_IMPORT (.*)$/)
          if (match) {
            const importPath = match[1].trim()
            if (existsSync(resolve(filepath, '..', importPath))) {
              const importLines = [
                '---',
                `src: ${importPath}`,
                '---',
                '',
              ]
              lines.splice(Number(i), 1, ...importLines)
            }
            else {
              lines[i] = `<!-- No ${importPath} -->`
            }
          }
        }
      },
    },
  ]
})
