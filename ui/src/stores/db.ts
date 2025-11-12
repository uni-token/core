import { defineStore } from 'pinia'
import { useServiceStore } from './service'

export function defineDbStore<T>(name: string) {
  return defineStore(`db:${name}`, () => {
    const serviceStore = useServiceStore()
    return {
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
