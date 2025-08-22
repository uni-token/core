const PORT_RANGE_START = 18760

export async function findServicePort(): Promise<number | null> {
  for (let port = PORT_RANGE_START; port < PORT_RANGE_START + 10; port++) {
    if (await checkPort(port)) {
      return port
    }
  }
  return null
}

async function checkPort(port: number) {
  try {
    const response = await fetch(`http://localhost:${port}/`, {
      cache: 'no-cache',
      method: 'GET',
    })
    if ((await response.json()).__uni_token) {
      return true
    }
  }
  catch {}
  return false
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

export type UniTokenOpenAIResult = {
  type: 'not-found'
} | {
  type: 'not-granted'
} | {
  type: 'success'
  /**
   * The base URL for the OpenAI API.
   */
  baseUrl: string
  /**
   * The API key granted by UniToken for accessing OpenAI.
   */
  apiKey: string
}

/**
 * Requests user for OpenAI token via UniToken service.
 * @param options - The options for the request.
 * @returns An object containing the baseUrl, and apiKey.
 *   If the user did not grant permission, apiKey will be null.
 * @throws Possible network issues or service errors.
 */
export async function requestUniTokenOpenAI(options: UniTokenOptions): Promise<UniTokenOpenAIResult> {
  const port = await findServicePort()
  if (!port) {
    return { type: 'not-found' }
  }

  const serverUrl = `http://localhost:${port}/`
  const baseUrl = `${serverUrl}openai/`

  const abortController = new AbortController()
  setTimeout(() => {
    abortController.abort()
  }, 10000) // 10 seconds timeout

  const response = await fetch(`${serverUrl}app/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      name: options.appName,
      description: options.description,
      uid: options.savedApiKey,
    }),
    signal: abortController.signal,
  })

  if (response.status === 403) {
    return { type: 'not-granted' }
  }

  if (!response.ok) {
    const errorText = await response.text()
    throw new Error(`Registration failed: HTTP ${response.status} - ${errorText}`)
  }

  const responseJson = await response.json()
  return {
    type: 'success',
    baseUrl,
    apiKey: responseJson.token,
  }
}

export function startService(): void {
  window.open('uni-token://start')
}

export const SUPPORTED_OS = ['Windows', 'macOS', 'Linux'] as const
export type SupportedOS = (typeof SUPPORTED_OS)[number]

export function downloadService(os: SupportedOS): void {
  const filename = {
    Linux: 'service-linux-amd64',
    macOS: 'service-darwin-amd64',
    Windows: 'service-windows-amd64.exe',
  }[os]
  const url = `https://uni-token.app/release/${filename}`
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  link.click()
  link.remove()
}
