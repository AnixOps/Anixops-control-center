import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export type MFAStatus = 'disabled' | 'enabled' | 'pending'
export type MFAMethod = 'totp' | 'sms' | 'email' | 'backup_codes'

export interface MFAConfig {
  enabled: boolean
  method: MFAMethod
  backupCodes: string[]
  phoneNumber?: string
  email?: string
  totpSecret?: string
  createdAt?: string
  lastUsedAt?: string
}

export interface MFASetupResponse {
  secret: string
  qrCodeUrl: string
  backupCodes: string[]
  manualEntryKey: string
}

export interface MFAVerification {
  code: string
  method: MFAMethod
  rememberDevice: boolean
}

const STORAGE_KEY = 'mfa-config'

export const useMFAStore = defineStore('mfa', () => {
  const status = ref<MFAStatus>('disabled')
  const config = ref<MFAConfig | null>(null)
  const isSetupInProgress = ref(false)
  const setupStep = ref(1)
  const setupData = ref<MFASetupResponse | null>(null)
  const verificationCode = ref('')
  const rememberDevice = ref(false)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // Computed
  const isEnabled = computed(() => status.value === 'enabled')
  const isPending = computed(() => status.value === 'pending')
  const hasBackupCodes = computed(() => config.value?.backupCodes?.length > 0)
  const availableMethods = computed(() => {
    return [
      { id: 'totp', name: 'Authenticator App', description: 'Use Google Authenticator, Authy, or similar app', icon: 'smartphone' },
      { id: 'sms', name: 'SMS', description: 'Receive codes via text message', icon: 'message' },
      { id: 'email', name: 'Email', description: 'Receive codes via email', icon: 'mail' },
      { id: 'backup_codes', name: 'Backup Codes', description: 'Use one-time backup codes', icon: 'key' }
    ]
  })

  // Load MFA config
  async function loadConfig() {
    isLoading.value = true
    error.value = null

    try {
      // Simulate API call
      const stored = localStorage.getItem(STORAGE_KEY)
      if (stored) {
        config.value = JSON.parse(stored)
        status.value = config.value?.enabled ? 'enabled' : 'disabled'
      }
    } catch (e) {
      error.value = 'Failed to load MFA configuration'
    } finally {
      isLoading.value = false
    }
  }

  // Start MFA setup
  async function startSetup(method: MFAMethod): Promise<MFASetupResponse | null> {
    isLoading.value = true
    error.value = null
    isSetupInProgress.value = true
    setupStep.value = 1

    try {
      // Generate setup data (in real app, this comes from server)
      const secret = generateSecret()
      const backupCodes = generateBackupCodes()

      const setup: MFASetupResponse = {
        secret,
        qrCodeUrl: `otpauth://totp/AnixOps:user@example.com?secret=${secret}&issuer=AnixOps`,
        backupCodes,
        manualEntryKey: formatSecret(secret)
      }

      setupData.value = setup
      return setup
    } catch (e) {
      error.value = 'Failed to start MFA setup'
      return null
    } finally {
      isLoading.value = false
    }
  }

  // Verify setup code
  async function verifySetupCode(code: string): Promise<boolean> {
    isLoading.value = true
    error.value = null

    try {
      // Verify the code (in real app, this validates against server)
      if (code.length === 6 && /^\d+$/.test(code)) {
        setupStep.value = 2
        return true
      }
      error.value = 'Invalid verification code'
      return false
    } finally {
      isLoading.value = false
    }
  }

  // Complete setup
  async function completeSetup(): Promise<boolean> {
    if (!setupData.value) return false

    isLoading.value = true
    error.value = null

    try {
      config.value = {
        enabled: true,
        method: 'totp',
        backupCodes: setupData.value.backupCodes,
        createdAt: new Date().toISOString()
      }

      localStorage.setItem(STORAGE_KEY, JSON.stringify(config.value))
      status.value = 'enabled'
      isSetupInProgress.value = false
      setupStep.value = 0
      return true
    } catch (e) {
      error.value = 'Failed to complete MFA setup'
      return false
    } finally {
      isLoading.value = false
    }
  }

  // Cancel setup
  function cancelSetup() {
    isSetupInProgress.value = false
    setupStep.value = 0
    setupData.value = null
    verificationCode.value = ''
    error.value = null
  }

  // Verify code for login
  async function verifyCode(verification: MFAVerification): Promise<boolean> {
    isLoading.value = true
    error.value = null

    try {
      // Verify the code
      if (verification.code.length === 6 && /^\d+$/.test(verification.code)) {
        if (config.value) {
          config.value.lastUsedAt = new Date().toISOString()
          localStorage.setItem(STORAGE_KEY, JSON.stringify(config.value))
        }
        return true
      }

      // Check backup codes
      if (config.value?.backupCodes.includes(verification.code)) {
        // Remove used backup code
        config.value.backupCodes = config.value.backupCodes.filter(c => c !== verification.code)
        localStorage.setItem(STORAGE_KEY, JSON.stringify(config.value))
        return true
      }

      error.value = 'Invalid verification code'
      return false
    } finally {
      isLoading.value = false
    }
  }

  // Regenerate backup codes
  async function regenerateBackupCodes(): Promise<string[]> {
    isLoading.value = true

    try {
      const codes = generateBackupCodes()
      if (config.value) {
        config.value.backupCodes = codes
        localStorage.setItem(STORAGE_KEY, JSON.stringify(config.value))
      }
      return codes
    } finally {
      isLoading.value = false
    }
  }

  // Disable MFA
  async function disableMFA(password: string): Promise<boolean> {
    isLoading.value = true
    error.value = null

    try {
      // Verify password (in real app, server validates)
      config.value = null
      status.value = 'disabled'
      localStorage.removeItem(STORAGE_KEY)
      return true
    } catch (e) {
      error.value = 'Failed to disable MFA'
      return false
    } finally {
      isLoading.value = false
    }
  }

  // Helper functions
  function generateSecret(): string {
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ234567'
    let secret = ''
    for (let i = 0; i < 16; i++) {
      secret += chars.charAt(Math.floor(Math.random() * chars.length))
    }
    return secret
  }

  function formatSecret(secret: string): string {
    return secret.match(/.{1,4}/g)?.join(' ') || secret
  }

  function generateBackupCodes(): string[] {
    const codes: string[] = []
    for (let i = 0; i < 8; i++) {
      const code = Math.random().toString(36).substring(2, 8).toUpperCase()
      codes.push(code)
    }
    return codes
  }

  // Initialize
  loadConfig()

  return {
    status,
    config,
    isSetupInProgress,
    setupStep,
    setupData,
    verificationCode,
    rememberDevice,
    isLoading,
    error,
    isEnabled,
    isPending,
    hasBackupCodes,
    availableMethods,
    loadConfig,
    startSetup,
    verifySetupCode,
    completeSetup,
    cancelSetup,
    verifyCode,
    regenerateBackupCodes,
    disableMFA
  }
})