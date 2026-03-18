/**
 * Form validation utilities
 */

/**
 * Validation rules
 */
export const rules = {
  required: (message = 'This field is required') => (value) => {
    if (value === null || value === undefined || value === '') return message
    if (Array.isArray(value) && value.length === 0) return message
    return true
  },

  email: (message = 'Invalid email address') => (value) => {
    if (!value) return true
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    return emailRegex.test(value) || message
  },

  minLength: (min, message) => (value) => {
    if (!value) return true
    return value.length >= min || message || `Minimum ${min} characters required`
  },

  maxLength: (max, message) => (value) => {
    if (!value) return true
    return value.length <= max || message || `Maximum ${max} characters allowed`
  },

  min: (minVal, message) => (value) => {
    if (value === null || value === undefined || value === '') return true
    return Number(value) >= minVal || message || `Minimum value is ${minVal}`
  },

  max: (maxVal, message) => (value) => {
    if (value === null || value === undefined || value === '') return true
    return Number(value) <= maxVal || message || `Maximum value is ${maxVal}`
  },

  pattern: (regex, message = 'Invalid format') => (value) => {
    if (!value) return true
    return regex.test(value) || message
  },

  url: (message = 'Invalid URL') => (value) => {
    if (!value) return true
    try {
      new URL(value)
      return true
    } catch {
      return message
    }
  },

  numeric: (message = 'Must be a number') => (value) => {
    if (!value) return true
    return !isNaN(Number(value)) || message
  },

  integer: (message = 'Must be an integer') => (value) => {
    if (!value) return true
    return Number.isInteger(Number(value)) || message
  },

  positive: (message = 'Must be positive') => (value) => {
    if (!value) return true
    return Number(value) > 0 || message
  },

  match: (fieldName, getValue, message) => (value) => {
    const otherValue = getValue()
    return value === otherValue || message || `Must match ${fieldName}`
  },

  custom: (validator, message) => (value) => {
    return validator(value) || message
  }
}

/**
 * Validate a single value against rules
 */
export function validate(value, ruleList) {
  for (const rule of ruleList) {
    const result = rule(value)
    if (result !== true) {
      return result
    }
  }
  return true
}

/**
 * Validate form data against schema
 */
export function validateForm(data, schema) {
  const errors = {}
  let isValid = true

  for (const [field, rules] of Object.entries(schema)) {
    const value = data[field]
    const error = validate(value, rules)
    if (error !== true) {
      errors[field] = error
      isValid = false
    }
  }

  return { isValid, errors }
}

/**
 * Common validation schemas
 */
export const schemas = {
  login: {
    email: [rules.required(), rules.email()],
    password: [rules.required(), rules.minLength(6)]
  },

  register: {
    email: [rules.required(), rules.email()],
    password: [rules.required(), rules.minLength(8)],
    confirmPassword: [rules.required()]
  },

  node: {
    name: [rules.required(), rules.minLength(2), rules.maxLength(50)],
    host: [rules.required()],
    port: [rules.required(), rules.numeric(), rules.min(1), rules.max(65535)]
  },

  user: {
    email: [rules.required(), rules.email()],
    name: [rules.maxLength(100)],
    password: [rules.minLength(8)]
  }
}