import type { App } from './app'
import type { APIKey } from './keys'
import { defineStore } from 'pinia'
import { useServiceStore } from './service'

function defineDbStore<T>(name: string) {
  return defineStore(`db:${name}`, () => {
    const serviceStore = useServiceStore()
    return {
      async all(): Promise<Record<string, T>> {
        const resp = await serviceStore.api(`store/${name}`)
        if (resp.status === 200) {
          const raw = await resp.json() as Record<string, string>
          const result: Record<string, T> = {}
          for (const [key, value] of Object.entries(raw)) {
            result[key] = JSON.parse(value) as T
          }
          return result
        }
        else {
          throw new Error(`Error fetching all keys from store ${name}: ${resp.statusText}`)
        }
      },
      async clear(): Promise<void> {
        const resp = await serviceStore.api(`store/${name}`, {
          method: 'DELETE',
        })
        if (resp.status !== 200) {
          throw new Error(`Error clearing store ${name}: ${resp.statusText}`)
        }
      },
      async get(key: string): Promise<T | null> {
        const resp = await serviceStore.api(`store/${name}/${key}`)
        if (resp.status === 200) {
          return await resp.json() as T
        }
        else if (resp.status === 404) {
          return null
        }
        else {
          throw new Error(`Error fetching key ${key} from store ${name}: ${resp.statusText}`)
        }
      },
      async put(key: string, value: T): Promise<void> {
        const resp = await serviceStore.api(`store/${name}/${key}`, {
          method: 'PUT',
          body: JSON.stringify(value),
        })
        if (resp.status !== 200) {
          throw new Error(`Error putting key ${key} to store ${name}: ${resp.statusText}`)
        }
      },
      async delete(key: string): Promise<void> {
        const resp = await serviceStore.api(`store/${name}/${key}`, {
          method: 'DELETE',
        })
        if (resp.status !== 200) {
          throw new Error(`Error deleting key ${key} from store ${name}: ${resp.statusText}`)
        }
      },
    }
  })
}

export const useAppsDb = defineDbStore<App>('apps')
export const useProviderSessionsDb = defineDbStore<unknown>('provider_sessions')
export const useKeysDb = defineDbStore<APIKey>('llm_keys')
