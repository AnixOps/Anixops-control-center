/**
 * Utility functions for the AnixOps Web GUI
 */

/**
 * Format bytes to human readable string
 */
export function formatBytes(bytes, decimals = 2) {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(decimals)) + ' ' + sizes[i]
}

/**
 * Format number with thousands separator
 */
export function formatNumber(num, locale = 'en-US') {
  return new Intl.NumberFormat(locale).format(num)
}

/**
 * Format currency
 */
export function formatCurrency(amount, currency = 'USD', locale = 'en-US') {
  return new Intl.NumberFormat(locale, {
    style: 'currency',
    currency
  }).format(amount)
}

/**
 * Format date
 */
export function formatDate(date, options = {}) {
  const defaultOptions = {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    ...options
  }
  return new Intl.DateTimeFormat('en-US', defaultOptions).format(new Date(date))
}

/**
 * Format relative time (e.g., "2 hours ago")
 */
export function formatRelativeTime(date) {
  const now = new Date()
  const past = new Date(date)
  const diffMs = now - past
  const diffSec = Math.floor(diffMs / 1000)
  const diffMin = Math.floor(diffSec / 60)
  const diffHour = Math.floor(diffMin / 60)
  const diffDay = Math.floor(diffHour / 24)
  const diffWeek = Math.floor(diffDay / 7)
  const diffMonth = Math.floor(diffDay / 30)
  const diffYear = Math.floor(diffDay / 365)

  if (diffSec < 60) return 'just now'
  if (diffMin < 60) return `${diffMin} minute${diffMin > 1 ? 's' : ''} ago`
  if (diffHour < 24) return `${diffHour} hour${diffHour > 1 ? 's' : ''} ago`
  if (diffDay < 7) return `${diffDay} day${diffDay > 1 ? 's' : ''} ago`
  if (diffWeek < 4) return `${diffWeek} week${diffWeek > 1 ? 's' : ''} ago`
  if (diffMonth < 12) return `${diffMonth} month${diffMonth > 1 ? 's' : ''} ago`
  return `${diffYear} year${diffYear > 1 ? 's' : ''} ago`
}

/**
 * Format duration in seconds to human readable
 */
export function formatDuration(seconds) {
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  const secs = seconds % 60

  if (hours > 0) {
    return `${hours}h ${minutes}m ${secs}s`
  }
  if (minutes > 0) {
    return `${minutes}m ${secs}s`
  }
  return `${secs}s`
}

/**
 * Truncate string with ellipsis
 */
export function truncate(str, length = 50) {
  if (!str || str.length <= length) return str
  return str.slice(0, length) + '...'
}

/**
 * Generate a random ID
 */
export function generateId(prefix = '') {
  const random = Math.random().toString(36).substring(2, 9)
  const timestamp = Date.now().toString(36)
  return prefix ? `${prefix}_${timestamp}${random}` : `${timestamp}${random}`
}

/**
 * Deep clone an object
 */
export function deepClone(obj) {
  return JSON.parse(JSON.stringify(obj))
}

/**
 * Check if value is empty (null, undefined, empty string, empty array, empty object)
 */
export function isEmpty(value) {
  if (value === null || value === undefined) return true
  if (typeof value === 'string') return value.trim() === ''
  if (Array.isArray(value)) return value.length === 0
  if (typeof value === 'object') return Object.keys(value).length === 0
  return false
}

/**
 * Sleep for a given number of milliseconds
 */
export function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms))
}

/**
 * Debounce function
 */
export function debounce(func, wait) {
  let timeout
  return function executedFunction(...args) {
    const later = () => {
      clearTimeout(timeout)
      func(...args)
    }
    clearTimeout(timeout)
    timeout = setTimeout(later, wait)
  }
}

/**
 * Throttle function
 */
export function throttle(func, limit) {
  let inThrottle
  return function executedFunction(...args) {
    if (!inThrottle) {
      func(...args)
      inThrottle = true
      setTimeout(() => (inThrottle = false), limit)
    }
  }
}

/**
 * Parse query string to object
 */
export function parseQueryString(queryString) {
  const params = new URLSearchParams(queryString)
  const result = {}
  for (const [key, value] of params) {
    result[key] = value
  }
  return result
}

/**
 * Convert object to query string
 */
export function toQueryString(obj) {
  const params = new URLSearchParams()
  for (const [key, value] of Object.entries(obj)) {
    if (value !== null && value !== undefined && value !== '') {
      params.append(key, value)
    }
  }
  return params.toString()
}

/**
 * Copy text to clipboard
 */
export async function copyToClipboard(text) {
  try {
    await navigator.clipboard.writeText(text)
    return true
  } catch {
    // Fallback for older browsers
    const textarea = document.createElement('textarea')
    textarea.value = text
    textarea.style.position = 'fixed'
    textarea.style.opacity = '0'
    document.body.appendChild(textarea)
    textarea.select()
    const success = document.execCommand('copy')
    document.body.removeChild(textarea)
    return success
  }
}

/**
 * Download a file from a blob
 */
export function downloadBlob(blob, filename) {
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}

/**
 * Download data as JSON file
 */
export function downloadAsJson(data, filename) {
  const json = JSON.stringify(data, null, 2)
  const blob = new Blob([json], { type: 'application/json' })
  downloadBlob(blob, filename)
}

/**
 * Download data as CSV file
 */
export function downloadAsCsv(data, filename) {
  if (!Array.isArray(data) || data.length === 0) return

  const headers = Object.keys(data[0])
  const csvRows = [
    headers.join(','),
    ...data.map(row =>
      headers.map(header => {
        const cell = row[header]
        const escaped = String(cell).replace(/"/g, '""')
        return `"${escaped}"`
      }).join(',')
    )
  ]

  const csv = csvRows.join('\n')
  const blob = new Blob([csv], { type: 'text/csv' })
  downloadBlob(blob, filename)
}

/**
 * Get file extension from filename
 */
export function getFileExtension(filename) {
  const parts = filename.split('.')
  return parts.length > 1 ? parts.pop().toLowerCase() : ''
}

/**
 * Format file size for display
 */
export function formatFileSize(bytes) {
  return formatBytes(bytes, 1)
}

/**
 * Check if a color is light or dark
 */
export function isLightColor(color) {
  const hex = color.replace('#', '')
  const r = parseInt(hex.substr(0, 2), 16)
  const g = parseInt(hex.substr(2, 2), 16)
  const b = parseInt(hex.substr(4, 2), 16)
  const brightness = (r * 299 + g * 587 + b * 114) / 1000
  return brightness > 128
}

/**
 * Generate a random color
 */
export function randomColor() {
  return '#' + Math.floor(Math.random() * 16777215).toString(16).padStart(6, '0')
}

/**
 * Get initials from a name
 */
export function getInitials(name, maxLength = 2) {
  if (!name) return ''
  return name
    .split(' ')
    .map(word => word.charAt(0).toUpperCase())
    .slice(0, maxLength)
    .join('')
}

/**
 * Validate email address
 */
export function isValidEmail(email) {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

/**
 * Validate URL
 */
export function isValidUrl(url) {
  try {
    new URL(url)
    return true
  } catch {
    return false
  }
}

/**
 * Parse JSON safely
 */
export function safeJsonParse(str, defaultValue = null) {
  try {
    return JSON.parse(str)
  } catch {
    return defaultValue
  }
}

/**
 * Get browser info
 */
export function getBrowserInfo() {
  const ua = navigator.userAgent
  let browser = 'Unknown'
  let version = ''

  if (ua.includes('Firefox')) {
    browser = 'Firefox'
    version = ua.match(/Firefox\/(\d+)/)?.[1] || ''
  } else if (ua.includes('Chrome')) {
    browser = 'Chrome'
    version = ua.match(/Chrome\/(\d+)/)?.[1] || ''
  } else if (ua.includes('Safari')) {
    browser = 'Safari'
    version = ua.match(/Version\/(\d+)/)?.[1] || ''
  } else if (ua.includes('Edge')) {
    browser = 'Edge'
    version = ua.match(/Edge\/(\d+)/)?.[1] || ''
  }

  return { browser, version, userAgent: ua }
}

/**
 * Check if running in electron
 */
export function isElectron() {
  return typeof process !== 'undefined' && process.versions && process.versions.electron
}

/**
 * Check if running in mobile browser
 */
export function isMobile() {
  return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent)
}

/**
 * Get contrast color (black or white) for a background color
 */
export function getContrastColor(hexColor) {
  const hex = hexColor.replace('#', '')
  const r = parseInt(hex.substr(0, 2), 16)
  const g = parseInt(hex.substr(2, 2), 16)
  const b = parseInt(hex.substr(4, 2), 16)
  const luminance = (0.299 * r + 0.587 * g + 0.114 * b) / 255
  return luminance > 0.5 ? '#000000' : '#ffffff'
}