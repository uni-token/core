import process from 'node:process'
import { OpenAI } from 'openai'
import { requestUniTokenOpenAI } from './index.ts'

function loadApiKey(): string | null {
  // ...
  return null
}

function saveApiKey(apiKey: string | null): void {
  // ...
  // eslint-disable-next-line no-console
  console.log('API key: ', apiKey)
}

async function main() {
  const result = await requestUniTokenOpenAI({
    appName: 'MyApp',
    description: 'This is a test application',
    // If not provided, the user will always be prompted to grant permission to this app.
    savedApiKey: loadApiKey(),
  })
  saveApiKey(result.apiKey)
  if (!result.apiKey) {
    console.error('User did not grant permission for OpenAI token.')
    return
  }

  const openai = new OpenAI({
    baseURL: result.baseUrl,
    apiKey: result.apiKey,
  })

  const stream = await openai.chat.completions.create({
    model: 'gpt-4o-mini',
    stream: true,
    messages: [
      { role: 'system', content: 'You are a concise assistant.' },
      { role: 'user', content: 'Please write a one-sentence bedtime story.' },
    ],
  })

  for await (const chunk of stream) {
    const delta = chunk.choices[0].delta
    if (delta?.content) {
      process.stdout.write(delta.content)
    }
  }
  process.stdout.write('\n')
}

main()
