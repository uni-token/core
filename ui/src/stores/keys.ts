import type { Provider } from '@/lib/providers'
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { toast } from 'vue-sonner'
import { useI18n } from '@/lib/locals'
import { useKeysDb } from './db'

export interface APIKey {
  id: string
  name: string
  type: string
  protocol: 'openai'
  baseUrl: string
  token: string
}

export const useKeysStore = defineStore('keys', () => {
  const db = useKeysDb()
  const { t } = useI18n({
    'en-US': {
      addKeyFailed: 'Failed to add key',
      updateKeyFailed: 'Failed to update key',
      deleteKeyFailed: 'Failed to delete key',
    },
    'zh-CN': {
      addKeyFailed: '添加密钥失败',
      updateKeyFailed: '更新密钥失败',
      deleteKeyFailed: '删除密钥失败',
    },
  })

  // State
  const keys = ref<APIKey[]>([])
  const loading = ref(true)
  const loadingError = ref<string | null>(null)

  // Actions
  async function loadKeys() {
    loadingError.value = null

    try {
      keys.value = Object.values(await db.all())
    }
    catch (err) {
      loadingError.value = err instanceof Error ? err.message : 'Unknown error'
    }
    finally {
      loading.value = false
    }
  }

  async function createAndAddKey(provider: Provider) {
    const key = await provider.createKey()
    await addKey({
      name: provider.name,
      type: provider.id,
      protocol: 'openai',
      baseUrl: provider.baseURL,
      token: key,
    })
  }

  async function addKey(key: Omit<APIKey, 'id' | 'name'> & { name?: string }): Promise<APIKey> {
    const id = crypto.randomUUID()
    const data: APIKey = {
      ...key,
      id,
      name: key.name ?? '',
    }
    try {
      await db.put(id, data)
      await loadKeys()
      return data
    }
    catch (err) {
      toast.error(t('addKeyFailed'))
      throw err
    }
  }

  async function updateKey(keyId: string, key: APIKey) {
    try {
      await db.put(keyId, key)
      await loadKeys()
    }
    catch (err) {
      toast.error(t('updateKeyFailed'))
      throw err
    }
  }

  async function deleteKey(keyId: string) {
    try {
      await db.delete(keyId)
      await loadKeys()
    }
    catch (err) {
      toast.error(t('deleteKeyFailed'))
      throw err
    }
  }

  return {
    // State
    keys,
    loading,
    loadingError,
    // Actions
    loadKeys,
    createAndAddKey,
    addKey,
    updateKey,
    deleteKey,
  }
})
