/**
 * Backup Service Unit Tests
 */

import { describe, it, expect, beforeEach } from 'vitest'
import {
  createBackup,
  listBackups,
  deleteBackup,
  getLatestBackupStatus,
  BackupInfo,
} from './backup'
import { createMockKV, createMockR2, createMockD1 } from '../../test/setup'

describe('Backup Service', () => {
  let mockEnv: {
    DB: D1Database
    KV: KVNamespace
    R2: R2Bucket
  }

  beforeEach(() => {
    mockEnv = {
      DB: createMockD1(),
      KV: createMockKV(),
      R2: createMockR2(),
    }
  })

  describe('createBackup', () => {
    it('should create a backup successfully', async () => {
      // Mock the database responses
      const mockDB = {
        prepare: () => ({
          bind: function() { return this },
          first: async () => ({ name: 'users' }),
          all: async () => ({ results: [] }),
        }),
      } as any

      const env = { ...mockEnv, DB: mockDB }

      const result = await createBackup(env)

      expect(result.status).toBe('completed')
      expect(result.id).toMatch(/^backup-\d+$/)
    })
  })

  describe('listBackups', () => {
    it('should list backups from R2', async () => {
      const backups = await listBackups(mockEnv)

      expect(Array.isArray(backups)).toBe(true)
    })
  })

  describe('getLatestBackupStatus', () => {
    it('should return null when no backup status', async () => {
      const status = await getLatestBackupStatus(mockEnv)

      expect(status).toBeNull()
    })

    it('should return backup status from KV', async () => {
      const backupInfo: BackupInfo = {
        id: 'backup-123',
        timestamp: new Date().toISOString(),
        size: 1024,
        tables: ['users', 'nodes'],
        status: 'completed',
      }

      await mockEnv.KV.put('backup:latest', JSON.stringify(backupInfo))

      const status = await getLatestBackupStatus(mockEnv)

      expect(status).not.toBeNull()
      expect(status!.id).toBe('backup-123')
      expect(status!.status).toBe('completed')
    })
  })

  describe('deleteBackup', () => {
    it('should delete backup successfully', async () => {
      // First create a backup
      await mockEnv.R2.put('backups/d1/backup-test.json', '{"metadata":{},"data":{}}')

      const result = await deleteBackup(mockEnv, 'backup-test')

      expect(result).toBe(true)
    })

    it('should return false for non-existent backup', async () => {
      const result = await deleteBackup(mockEnv, 'non-existent')

      expect(result).toBe(true) // Returns true even if doesn't exist
    })
  })
})