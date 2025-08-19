import { Buffer } from 'node:buffer'
import { spawnSync } from 'node:child_process'
import * as fs from 'node:fs'
import * as os from 'node:os'
import * as path from 'node:path'

function setupServiceRootPath(): string {
  const userHome = os.homedir()

  let root: string
  if (os.platform() === 'win32') {
    root = path.join(userHome, 'AppData', 'Local', 'UnitedToken')
  }
  else {
    root = path.join(userHome, '.local', 'share', 'uni-token')
  }

  if (!fs.existsSync(root)) {
    fs.mkdirSync(root, { recursive: true })
  }

  return root
}

async function detectRunningUrlFromFile(rootPath: string): Promise<string | null> {
  const filePath = path.join(rootPath, 'service.json')

  if (!fs.existsSync(filePath)) {
    return null
  }

  try {
    const fileContent = fs.readFileSync(filePath, 'utf8')
    const serviceInfo = JSON.parse(fileContent)
    const serverUrl = serviceInfo.url

    if (serverUrl) {
      const response = await fetch(serverUrl)
      const data = await response.json()
      if (data && typeof data === 'object' && '__united_token' in data) {
        return serverUrl
      }
    }
  }
  catch { }
  return null
}

async function startService(rootPath: string): Promise<string> {
  const execPath = os.platform() === 'win32'
    ? path.join(rootPath, 'service.exe')
    : path.join(rootPath, 'service')

  if (!fs.existsSync(execPath)) {
    await downloadService(execPath)
  }

  const result = spawnSync(execPath, ['start'], {
    stdio: 'inherit',
    timeout: 10000,
    encoding: 'utf8',
  })

  if (result.error) {
    throw new Error(`Execution error: ${result.error.message}`)
  }
  if (result.status !== 0) {
    const output = (result.stdout || result.stderr || '').toString()
    throw new Error(`Service exited with code ${result.status}. Output: ${output}`)
  }

  const serverUrl = await detectRunningUrlFromFile(rootPath)
  if (!serverUrl) {
    throw new Error('Service started but URL not detected.')
  }
  return serverUrl
}

async function downloadService(execPath: string): Promise<void> {
  const platformMap: { [key: string]: string } = {
    linux: 'service-linux-amd64',
    darwin: 'service-darwin-amd64',
    win32: 'service-windows-amd64.exe',
  }
  const filename = platformMap[os.platform()]

  if (!filename) {
    throw new Error(`Unsupported platform: ${os.platform()}`)
  }

  const url = `https://uni-token.app/release/${filename}`
  const response = await fetch(url)

  if (!response.ok) {
    throw new Error(`Download failed: HTTP ${response.status} - ${response.statusText}`)
  }

  const arrayBuffer = await response.arrayBuffer()
  fs.writeFileSync(execPath, Buffer.from(arrayBuffer))
  fs.chmodSync(execPath, 0o755)
}

export interface UniTokenOptions {
  /**
   * The name of the application requesting the OpenAI token.
   */
  appName: string
  /**
   * A brief description of the application.
   */
  description: string
  /**
   * Optional saved API key, if the user has previously granted permission.
   */
  savedApiKey?: string | null | undefined
}

/**
 * Requests user for OpenAI token via UniToken service.
 * @param options - The options for the request.
 * @returns An object containing the appID, baseUrl, and apiKey, or null if the user does not grant permission.
 * @throws Possible network issues or service errors.
 */
export async function requestUniTokenOpenAI(options: UniTokenOptions): Promise<{
  baseUrl: string
  apiKey: string
} | null> {
  const rootPath = setupServiceRootPath()
  const serverUrl = await detectRunningUrlFromFile(rootPath) || await startService(rootPath)

  const response = await fetch(`${serverUrl}app/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      name: options.appName,
      description: options.description,
      uid: options.savedApiKey,
    }),
  })

  if (response.status === 403) {
    return null
  }

  if (!response.ok) {
    const errorText = await response.text()
    throw new Error(`Registration failed: HTTP ${response.status} - ${errorText}`)
  }

  const responseJson = await response.json()
  return {
    baseUrl: `${serverUrl}openai/`,
    apiKey: responseJson.token,
  }
}
