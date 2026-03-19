/**
 * Audit Utilities Unit Tests
 */

import { describe, it, expect, vi, beforeEach } from 'vitest'
import {
  getRequiredParam,
  getClientIP,
  getUserAgent,
  logAudit,
} from './audit'
import { createMockKV, createMockD1 } from '../../test/setup'

// Mock Hono Context
function createMockContext(overrides: any = {}): any {
  return {
    req: {
      param: vi.fn((name: string) => overrides.param?.[name]),
      header: vi.fn((name: string) => overrides.headers?.[name]),
    },
    env: {
      DB: createMockD1(),
      KV: createMockKV(),
      ...overrides.env,
    },
    ...overrides,
  }
}

describe('Audit Utilities', () => {
  describe('getRequiredParam', () => {
    it('should return param value when present', () => {
      const c = createMockContext({ param: { id: '123' } })
      const result = getRequiredParam(c, 'id')
      expect(result).toBe('123')
    })

    it('should throw error when param is missing', () => {
      const c = createMockContext()
      expect(() => getRequiredParam(c, 'missing')).toThrow('Missing required parameter: missing')
    })
  })

  describe('getClientIP', () => {
    it('should return client IP from CF-Connecting-IP header', () => {
      const c = createMockContext({ headers: { 'CF-Connecting-IP': '192.168.1.1' } })
      const result = getClientIP(c)
      expect(result).toBe('192.168.1.1')
    })

    it('should return null when header is missing', () => {
      const c = createMockContext({ headers: {} })
      const result = getClientIP(c)
      expect(result).toBeNull()
    })
  })

  describe('getUserAgent', () => {
    it('should return user agent from header', () => {
      const c = createMockContext({ headers: { 'User-Agent': 'Mozilla/5.0' } })
      const result = getUserAgent(c)
      expect(result).toBe('Mozilla/5.0')
    })

    it('should return null when header is missing', () => {
      const c = createMockContext({ headers: {} })
      const result = getUserAgent(c)
      expect(result).toBeNull()
    })
  })

  describe('logAudit', () => {
    it('should insert audit log to database', async () => {
      const mockDB = createMockD1()
      const c = createMockContext({
        env: { DB: mockDB },
        headers: {
          'CF-Connecting-IP': '10.0.0.1',
          'User-Agent': 'TestAgent',
        },
      })

      await logAudit(c, 1, 'login', 'auth', { email: 'test@example.com' })

      // Verify prepare was called
      expect(mockDB.prepare).toHaveBeenCalled()
    })

    it('should handle undefined user ID', async () => {
      const c = createMockContext()
      await logAudit(c, undefined, 'test', 'test')
      expect(c.env.DB.prepare).toHaveBeenCalled()
    })

    it('should handle errors gracefully', async () => {
      const c = createMockContext({
        env: {
          DB: {
            prepare: vi.fn(() => ({
              bind: vi.fn(() => ({
                run: vi.fn(() => Promise.reject(new Error('DB error'))),
              })),
            })),
          },
        },
      })

      // Should not throw
      await expect(logAudit(c, 1, 'test', 'test')).resolves.not.toThrow()
    })
  })
})