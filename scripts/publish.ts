import { resolve } from 'node:path'
import { $ } from 'zx'

async function publishPythonSdk() {
  console.log('Building package...')

  $.cwd = resolve(import.meta.dirname, '../sdk/python')

  await $`rm -rf dist/ build/ *.egg-info/`

  await $`uvx --from build pyproject-build --installer uv`

  await $`uvx twine check dist/*`

  console.log(`Publishing to pypi...`)

  await $`uvx twine upload dist/*`
  console.log('Package published to PyPI successfully!')
}

publishPythonSdk()
