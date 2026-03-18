import { ref, onMounted, onUnmounted, computed, watch } from 'vue'

/**
 * Composable for keyboard navigation
 */
export function useKeyboardNavigation(options = {}) {
  const {
    onEnter,
    onEscape,
    onArrowUp,
    onArrowDown,
    onArrowLeft,
    onArrowRight,
    onTab,
    onSpace,
    onHome,
    onEnd,
    onPageUp,
    onPageDown,
    preventDefault = true
  } = options

  function handleKeyDown(event) {
    const handlers = {
      'Enter': onEnter,
      'Escape': onEscape,
      'ArrowUp': onArrowUp,
      'ArrowDown': onArrowDown,
      'ArrowLeft': onArrowLeft,
      'ArrowRight': onArrowRight,
      'Tab': onTab,
      ' ': onSpace,
      'Home': onHome,
      'End': onEnd,
      'PageUp': onPageUp,
      'PageDown': onPageDown
    }

    const handler = handlers[event.key]
    if (handler) {
      if (preventDefault) {
        event.preventDefault()
      }
      handler(event)
    }
  }

  onMounted(() => {
    document.addEventListener('keydown', handleKeyDown)
  })

  onUnmounted(() => {
    document.removeEventListener('keydown', handleKeyDown)
  })

  return {
    handleKeyDown
  }
}

/**
 * Composable for focus trap (modal dialogs)
 */
export function useFocusTrap(containerRef) {
  const focusableSelectors = [
    'button:not([disabled])',
    'input:not([disabled])',
    'select:not([disabled])',
    'textarea:not([disabled])',
    'a[href]',
    '[tabindex]:not([tabindex="-1"])'
  ].join(', ')

  let lastFocusedElement = null

  function getFocusableElements() {
    if (!containerRef.value) return []
    return Array.from(containerRef.value.querySelectorAll(focusableSelectors))
  }

  function trapFocus(event) {
    const focusableElements = getFocusableElements()
    if (focusableElements.length === 0) return

    const firstElement = focusableElements[0]
    const lastElement = focusableElements[focusableElements.length - 1]

    if (event.key === 'Tab') {
      if (event.shiftKey) {
        if (document.activeElement === firstElement) {
          event.preventDefault()
          lastElement.focus()
        }
      } else {
        if (document.activeElement === lastElement) {
          event.preventDefault()
          firstElement.focus()
        }
      }
    }
  }

  function activate() {
    lastFocusedElement = document.activeElement
    document.addEventListener('keydown', trapFocus)

    // Focus first focusable element
    const focusableElements = getFocusableElements()
    if (focusableElements.length > 0) {
      focusableElements[0].focus()
    }
  }

  function deactivate() {
    document.removeEventListener('keydown', trapFocus)
    if (lastFocusedElement) {
      lastFocusedElement.focus()
    }
  }

  return {
    activate,
    deactivate,
    getFocusableElements
  }
}

/**
 * Composable for skip links
 */
export function useSkipLink(targetId) {
  function skipToContent() {
    const target = document.getElementById(targetId)
    if (target) {
      target.setAttribute('tabindex', '-1')
      target.focus()
      target.removeAttribute('tabindex')
    }
  }

  return {
    skipToContent
  }
}

/**
 * Composable for high contrast mode detection
 */
export function useHighContrast() {
  const isHighContrast = ref(false)

  function checkHighContrast() {
    // Check for Windows high contrast mode
    const isWindowsHighContrast = window.matchMedia('(forced-colors: active)').matches
    // Check for macOS high contrast
    const isMacHighContrast = window.matchMedia('(prefers-contrast: more)').matches

    isHighContrast.value = isWindowsHighContrast || isMacHighContrast

    // Apply high contrast class to body
    if (isHighContrast.value) {
      document.documentElement.classList.add('high-contrast')
    } else {
      document.documentElement.classList.remove('high-contrast')
    }
  }

  onMounted(() => {
    checkHighContrast()
    window.matchMedia('(forced-colors: active)').addEventListener('change', checkHighContrast)
    window.matchMedia('(prefers-contrast: more)').addEventListener('change', checkHighContrast)
  })

  onUnmounted(() => {
    window.matchMedia('(forced-colors: active)').removeEventListener('change', checkHighContrast)
    window.matchMedia('(prefers-contrast: more)').removeEventListener('change', checkHighContrast)
  })

  return {
    isHighContrast
  }
}

/**
 * Composable for reduced motion preferences
 */
export function useReducedMotion() {
  const prefersReducedMotion = ref(false)

  function checkReducedMotion() {
    prefersReducedMotion.value = window.matchMedia('(prefers-reduced-motion: reduce)').matches

    if (prefersReducedMotion.value) {
      document.documentElement.classList.add('reduced-motion')
    } else {
      document.documentElement.classList.remove('reduced-motion')
    }
  }

  onMounted(() => {
    checkReducedMotion()
    window.matchMedia('(prefers-reduced-motion: reduce)').addEventListener('change', checkReducedMotion)
  })

  onUnmounted(() => {
    window.matchMedia('(prefers-reduced-motion: reduce)').removeEventListener('change', checkReducedMotion)
  })

  return {
    prefersReducedMotion
  }
}

/**
 * Composable for screen reader announcements
 */
export function useAnnouncer() {
  function announce(message, priority = 'polite') {
    const announcer = document.createElement('div')
    announcer.setAttribute('aria-live', priority)
    announcer.setAttribute('aria-atomic', 'true')
    announcer.setAttribute('class', 'sr-only')
    announcer.style.cssText = 'position: absolute; width: 1px; height: 1px; padding: 0; margin: -1px; overflow: hidden; clip: rect(0, 0, 0, 0); white-space: nowrap; border: 0;'

    document.body.appendChild(announcer)

    // Small delay to ensure the element is in the DOM
    setTimeout(() => {
      announcer.textContent = message
    }, 100)

    // Remove after announcement
    setTimeout(() => {
      document.body.removeChild(announcer)
    }, 1000)
  }

  return {
    announce
  }
}

/**
 * Composable for WebSocket real-time updates
 */
export function useWebSocket(url) {
  const ws = ref(null)
  const isConnected = ref(false)
  const error = ref(null)
  const lastMessage = ref(null)
  const reconnectAttempts = ref(0)

  let reconnectTimer = null
  let heartbeatTimer = null
  const maxReconnectAttempts = 5
  const reconnectDelay = 3000
  const heartbeatInterval = 30000

  const eventHandlers = new Map()

  function connect() {
    if (ws.value?.readyState === WebSocket.OPEN) return

    try {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
      const wsUrl = url || `${protocol}//${window.location.host}/api/v1/ws`
      const token = localStorage.getItem('token')

      ws.value = new WebSocket(`${wsUrl}?token=${token}`)

      ws.value.onopen = () => {
        isConnected.value = true
        error.value = null
        reconnectAttempts.value = 0
        startHeartbeat()
      }

      ws.value.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          lastMessage.value = data

          // Dispatch to event handlers
          if (data.event && eventHandlers.has(data.event)) {
            eventHandlers.get(data.event).forEach(handler => handler(data.data))
          }
        } catch (e) {
          console.error('WebSocket message parse error:', e)
        }
      }

      ws.value.onerror = (e) => {
        error.value = 'WebSocket error'
        console.error('WebSocket error:', e)
      }

      ws.value.onclose = () => {
        isConnected.value = false
        scheduleReconnect()
      }
    } catch (e) {
      error.value = e.message
      scheduleReconnect()
    }
  }

  function disconnect() {
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
    if (heartbeatTimer) {
      clearInterval(heartbeatTimer)
      heartbeatTimer = null
    }
    if (ws.value) {
      ws.value.close()
      ws.value = null
    }
    isConnected.value = false
  }

  function scheduleReconnect() {
    if (reconnectAttempts.value >= maxReconnectAttempts) return

    reconnectTimer = setTimeout(() => {
      reconnectAttempts.value++
      connect()
    }, reconnectDelay)
  }

  function startHeartbeat() {
    heartbeatTimer = setInterval(() => {
      if (ws.value?.readyState === WebSocket.OPEN) {
        ws.value.send(JSON.stringify({ event: 'ping', data: {} }))
      }
    }, heartbeatInterval)
  }

  function subscribe(event, handler) {
    if (!eventHandlers.has(event)) {
      eventHandlers.set(event, [])
    }
    eventHandlers.get(event).push(handler)

    // Return unsubscribe function
    return () => {
      const handlers = eventHandlers.get(event)
      if (handlers) {
        const index = handlers.indexOf(handler)
        if (index !== -1) {
          handlers.splice(index, 1)
        }
      }
    }
  }

  function emit(event, data) {
    if (ws.value?.readyState === WebSocket.OPEN) {
      ws.value.send(JSON.stringify({ event, data }))
    }
  }

  onMounted(() => {
    connect()
  })

  onUnmounted(() => {
    disconnect()
  })

  return {
    ws,
    isConnected,
    error,
    lastMessage,
    connect,
    disconnect,
    subscribe,
    emit
  }
}

/**
 * Composable for API requests with loading state
 */
export function useApi() {
  const loading = ref(false)
  const error = ref(null)
  const data = ref(null)

  async function execute(apiCall) {
    loading.value = true
    error.value = null

    try {
      const response = await apiCall()
      data.value = response.data
      return response.data
    } catch (e) {
      error.value = e.response?.data?.message || e.message || 'An error occurred'
      throw e
    } finally {
      loading.value = false
    }
  }

  return {
    loading,
    error,
    data,
    execute
  }
}

/**
 * Composable for pagination
 */
export function usePagination(initialPage = 1, initialLimit = 20) {
  const currentPage = ref(initialPage)
  const pageSize = ref(initialLimit)
  const total = ref(0)

  const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

  function setPage(page) {
    currentPage.value = Math.max(1, Math.min(page, totalPages.value))
  }

  function setPageSize(size) {
    pageSize.value = size
    currentPage.value = 1
  }

  function setTotal(newTotal) {
    total.value = newTotal
  }

  function reset() {
    currentPage.value = 1
    total.value = 0
  }

  return {
    currentPage,
    pageSize,
    total,
    totalPages,
    setPage,
    setPageSize,
    setTotal,
    reset
  }
}

/**
 * Composable for debounced search
 */
export function useDebounce(initialValue = '', delay = 300) {
  const value = ref(initialValue)
  const debouncedValue = ref(initialValue)
  let timer = null

  watch(value, (newValue) => {
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => {
      debouncedValue.value = newValue
    }, delay)
  })

  return {
    value,
    debouncedValue
  }
}

/**
 * Composable for toast notifications
 */
export function useToast() {
  const toasts = ref([])

  function show(message, type = 'success', duration = 5000) {
    const id = Date.now()
    toasts.value.push({ id, message, type })

    if (duration > 0) {
      setTimeout(() => remove(id), duration)
    }

    return id
  }

  function remove(id) {
    const index = toasts.value.findIndex(t => t.id === id)
    if (index !== -1) {
      toasts.value.splice(index, 1)
    }
  }

  function success(message, duration) {
    return show(message, 'success', duration)
  }

  function error(message, duration) {
    return show(message, 'error', duration)
  }

  function warning(message, duration) {
    return show(message, 'warning', duration)
  }

  function info(message, duration) {
    return show(message, 'info', duration)
  }

  return {
    toasts,
    show,
    remove,
    success,
    error,
    warning,
    info
  }
}

/**
 * Composable for local storage with reactivity
 */
export function useLocalStorage(key, defaultValue = null) {
  const stored = localStorage.getItem(key)
  const value = ref(stored ? JSON.parse(stored) : defaultValue)

  watch(value, (newValue) => {
    if (newValue === null) {
      localStorage.removeItem(key)
    } else {
      localStorage.setItem(key, JSON.stringify(newValue))
    }
  }, { deep: true })

  return value
}

/**
 * Composable for form validation
 */
export function useForm(initialValues = {}, validationRules = {}) {
  const values = ref({ ...initialValues })
  const errors = ref({})
  const touched = ref({})

  function validate() {
    errors.value = {}

    for (const [field, rules] of Object.entries(validationRules)) {
      const value = values.value[field]

      for (const rule of rules) {
        const error = rule(value, values.value)
        if (error) {
          errors.value[field] = error
          break
        }
      }
    }

    return Object.keys(errors.value).length === 0
  }

  function setFieldError(field, error) {
    errors.value[field] = error
  }

  function clearFieldError(field) {
    delete errors.value[field]
  }

  function setTouched(field) {
    touched.value[field] = true
  }

  function reset() {
    values.value = { ...initialValues }
    errors.value = {}
    touched.value = {}
  }

  return {
    values,
    errors,
    touched,
    validate,
    setFieldError,
    clearFieldError,
    setTouched,
    reset
  }
}

// Validation rules
export const validators = {
  required: (message = 'This field is required') => (value) => {
    return value === null || value === undefined || value === '' ? message : null
  },

  email: (message = 'Invalid email address') => (value) => {
    if (!value) return null
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    return emailRegex.test(value) ? null : message
  },

  minLength: (min, message) => (value) => {
    if (!value) return null
    return value.length >= min ? null : message || `Minimum ${min} characters required`
  },

  maxLength: (max, message) => (value) => {
    if (!value) return null
    return value.length <= max ? null : message || `Maximum ${max} characters allowed`
  },

  pattern: (regex, message = 'Invalid format') => (value) => {
    if (!value) return null
    return regex.test(value) ? null : message
  },

  min: (minVal, message) => (value) => {
    if (value === null || value === undefined) return null
    return value >= minVal ? null : message || `Minimum value is ${minVal}`
  },

  max: (maxVal, message) => (value) => {
    if (value === null || value === undefined) return null
    return value <= maxVal ? null : message || `Maximum value is ${maxVal}`
  },

  match: (otherField, message = 'Fields do not match') => (value, values) => {
    return value === values[otherField] ? null : message
  }
}