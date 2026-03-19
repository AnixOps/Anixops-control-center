/**
 * Test setup file
 * Configures test environment with mocks for Cloudflare Workers
 */

import { beforeAll, afterAll, vi } from 'vitest'

// Mock KV Namespace
export function createMockKV(): KVNamespace {
  const store = new Map<string, { value: string; expiration?: number }>()

  return {
    get: vi.fn(async (key: string, options?: any) => {
      const item = store.get(key)
      if (!item) return null
      // Support both `kv.get(key, 'json')` and `kv.get(key, { type: 'json' })`
      const isJson = options === 'json' || options?.type === 'json'
      if (isJson) {
        try {
          return JSON.parse(item.value)
        } catch {
          return null
        }
      }
      return item.value
    }) as any,

    put: vi.fn(async (key: string, value: string, options?: any) => {
      store.set(key, { value, expiration: options?.expirationTtl })
    }) as any,

    delete: vi.fn(async (key: string) => {
      store.delete(key)
    }) as any,

    list: vi.fn(async () => ({
      keys: Array.from(store.keys()).map(name => ({ name })),
      list_complete: true,
    })) as any,

    getWithMetadata: vi.fn(async (key: string) => {
      const item = store.get(key)
      if (!item) return { value: null, metadata: null }
      return { value: item.value, metadata: null }
    }) as any,
  }
}

// Mock R2 Bucket
export function createMockR2(): R2Bucket {
  const store = new Map<string, { body: string; metadata?: any }>()

  return {
    get: vi.fn(async (key: string) => {
      const item = store.get(key)
      if (!item) return null
      return {
        key,
        body: item.body,
        size: item.body.length,
        text: async () => item.body,
        json: async () => JSON.parse(item.body),
        arrayBuffer: async () => new TextEncoder().encode(item.body).buffer,
      } as any
    }) as any,

    put: vi.fn(async (key: string, value: any, options?: any) => {
      const body = typeof value === 'string' ? value : JSON.stringify(value)
      store.set(key, { body, metadata: options?.customMetadata })
      return { key }
    }) as any,

    delete: vi.fn(async (key: string) => {
      store.delete(key)
    }) as any,

    list: vi.fn(async (options?: any) => ({
      objects: Array.from(store.entries())
        .filter(([key]) => !options?.prefix || key.startsWith(options.prefix))
        .map(([key, item]) => ({
          key,
          size: item.body.length,
          uploaded: new Date(),
        })),
      delimitedPrefixes: [],
    })) as any,

    head: vi.fn(async (key: string) => {
      const item = store.get(key)
      if (!item) return null
      return { key, size: item.body.length }
    }) as any,
  }
}

// Mock D1 Database
export function createMockD1(): D1Database {
  const tables = new Map<string, any[]>()

  return {
    prepare: vi.fn((sql: string) => {
      return {
        bind: vi.fn(function(this: any, ...args: any[]) {
          this._bindings = args
          return this
        }),
        first: vi.fn(async function(this: any) {
          // Simple mock implementation
          if (sql.includes('SELECT COUNT')) {
            return { count: 1 }
          }
          if (sql.includes('SELECT * FROM users WHERE email')) {
            return null
          }
          return null
        }),
        all: vi.fn(async function(this: any) {
          return { results: [] }
        }),
        run: vi.fn(async function(this: any) {
          return { success: true, meta: { changes: 1 } }
        }),
      }
    }) as any,

    batch: vi.fn(async (statements: any[]) => {
      return statements.map(() => ({ success: true }))
    }) as any,

    exec: vi.fn(async (sql: string) => {
      return { count: 0 }
    }) as any,

    withSession: vi.fn(() => createMockD1()) as any,
  }
}

// Global test environment
declare global {
  var testEnv: {
    KV: KVNamespace
    R2: R2Bucket
    DB: D1Database
    JWT_SECRET: string
    JWT_EXPIRE: string
    API_KEY_SALT: string
    ENVIRONMENT: string
  }
}

beforeAll(() => {
  globalThis.testEnv = {
    KV: createMockKV(),
    R2: createMockR2(),
    DB: createMockD1(),
    JWT_SECRET: 'test-secret-key-for-testing-min-32-characters',
    JWT_EXPIRE: '86400',
    API_KEY_SALT: 'test-salt',
    ENVIRONMENT: 'test',
  }
})