import { defineStore } from 'pinia'
import { ref, computed, onMounted, onUnmounted } from 'vue'

export interface KeyboardShortcut {
  id: string
  name: string
  description: string
  keys: string[]
  category: 'navigation' | 'actions' | 'editing' | 'view' | 'system'
  action: () => void
  enabled: boolean
  scope?: string // e.g., 'global', 'nodes', 'users'
}

export interface ShortcutCategory {
  id: string
  name: string
  shortcuts: KeyboardShortcut[]
}

const STORAGE_KEY = 'keyboard-shortcuts'

// Default shortcuts
const defaultShortcuts: KeyboardShortcut[] = [
  // Navigation
  { id: 'nav-dashboard', name: 'Go to Dashboard', description: 'Navigate to dashboard', keys: ['g', 'd'], category: 'navigation', action: () => {}, enabled: true, scope: 'global' },
  { id: 'nav-nodes', name: 'Go to Nodes', description: 'Navigate to nodes', keys: ['g', 'n'], category: 'navigation', action: () => {}, enabled: true, scope: 'global' },
  { id: 'nav-users', name: 'Go to Users', description: 'Navigate to users', keys: ['g', 'u'], category: 'navigation', action: () => {}, enabled: true, scope: 'global' },
  { id: 'nav-plugins', name: 'Go to Plugins', description: 'Navigate to plugins', keys: ['g', 'p'], category: 'navigation', action: () => {}, enabled: true, scope: 'global' },
  { id: 'nav-settings', name: 'Go to Settings', description: 'Navigate to settings', keys: ['g', 's'], category: 'navigation', action: () => {}, enabled: true, scope: 'global' },

  // Actions
  { id: 'action-search', name: 'Search', description: 'Open search', keys: ['/'], category: 'actions', action: () => {}, enabled: true, scope: 'global' },
  { id: 'action-new', name: 'New Item', description: 'Create new item', keys: ['n'], category: 'actions', action: () => {}, enabled: true, scope: 'global' },
  { id: 'action-delete', name: 'Delete', description: 'Delete selected item', keys: ['Delete'], category: 'actions', action: () => {}, enabled: true, scope: 'global' },
  { id: 'action-refresh', name: 'Refresh', description: 'Refresh current view', keys: ['r'], category: 'actions', action: () => {}, enabled: true, scope: 'global' },
  { id: 'action-save', name: 'Save', description: 'Save current changes', keys: ['Ctrl', 's'], category: 'actions', action: () => {}, enabled: true, scope: 'global' },

  // View
  { id: 'view-fullscreen', name: 'Toggle Fullscreen', description: 'Toggle fullscreen mode', keys: ['F11'], category: 'view', action: () => {}, enabled: true, scope: 'global' },
  { id: 'view-sidebar', name: 'Toggle Sidebar', description: 'Show/hide sidebar', keys: ['b'], category: 'view', action: () => {}, enabled: true, scope: 'global' },
  { id: 'view-theme', name: 'Toggle Theme', description: 'Switch between dark/light', keys: ['t'], category: 'view', action: () => {}, enabled: true, scope: 'global' },

  // Editing
  { id: 'edit-copy', name: 'Copy', description: 'Copy selected', keys: ['Ctrl', 'c'], category: 'editing', action: () => {}, enabled: true, scope: 'global' },
  { id: 'edit-paste', name: 'Paste', description: 'Paste from clipboard', keys: ['Ctrl', 'v'], category: 'editing', action: () => {}, enabled: true, scope: 'global' },
  { id: 'edit-cut', name: 'Cut', description: 'Cut selected', keys: ['Ctrl', 'x'], category: 'editing', action: () => {}, enabled: true, scope: 'global' },
  { id: 'edit-undo', name: 'Undo', description: 'Undo last action', keys: ['Ctrl', 'z'], category: 'editing', action: () => {}, enabled: true, scope: 'global' },
  { id: 'edit-redo', name: 'Redo', description: 'Redo last action', keys: ['Ctrl', 'Shift', 'z'], category: 'editing', action: () => {}, enabled: true, scope: 'global' },

  // System
  { id: 'help', name: 'Help', description: 'Show keyboard shortcuts', keys: ['?'], category: 'system', action: () => {}, enabled: true, scope: 'global' },
  { id: 'escape', name: 'Close/Cancel', description: 'Close dialog or cancel', keys: ['Escape'], category: 'system', action: () => {}, enabled: true, scope: 'global' },
]

export const useKeyboardShortcutsStore = defineStore('keyboardShortcuts', () => {
  const shortcuts = ref<KeyboardShortcut[]>([...defaultShortcuts])
  const isListening = ref(true)
  const showHelp = ref(false)
  const currentScope = ref('global')
  const keySequence = ref<string[]>([])
  const sequenceTimeout = ref<number | null>(null)

  // Computed
  const categories = computed<ShortcutCategory[]>(() => {
    const cats: Record<string, ShortcutCategory> = {
      navigation: { id: 'navigation', name: 'Navigation', shortcuts: [] },
      actions: { id: 'actions', name: 'Actions', shortcuts: [] },
      editing: { id: 'editing', name: 'Editing', shortcuts: [] },
      view: { id: 'view', name: 'View', shortcuts: [] },
      system: { id: 'system', name: 'System', shortcuts: [] }
    }

    shortcuts.value
      .filter(s => s.enabled)
      .forEach(s => {
        if (cats[s.category]) {
          cats[s.category].shortcuts.push(s)
        }
      })

    return Object.values(cats).filter(c => c.shortcuts.length > 0)
  })

  const enabledShortcuts = computed(() => {
    return shortcuts.value.filter(s => s.enabled && (s.scope === 'global' || s.scope === currentScope.value))
  })

  // Load custom shortcuts from storage
  function loadCustomShortcuts() {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored) {
      try {
        const custom = JSON.parse(stored)
        // Merge custom with defaults
        custom.forEach((cs: Partial<KeyboardShortcut>) => {
          const existing = shortcuts.value.find(s => s.id === cs.id)
          if (existing && cs.keys) {
            existing.keys = cs.keys
          }
        })
      } catch (e) {
        console.error('Failed to load custom shortcuts:', e)
      }
    }
  }

  // Save custom shortcuts
  function saveCustomShortcuts() {
    const custom = shortcuts.value.map(s => ({
      id: s.id,
      keys: s.keys
    }))
    localStorage.setItem(STORAGE_KEY, JSON.stringify(custom))
  }

  // Update shortcut
  function updateShortcut(id: string, keys: string[]) {
    const shortcut = shortcuts.value.find(s => s.id === id)
    if (shortcut) {
      shortcut.keys = keys
      saveCustomShortcuts()
    }
  }

  // Reset to defaults
  function resetToDefaults() {
    shortcuts.value = [...defaultShortcuts]
    localStorage.removeItem(STORAGE_KEY)
  }

  // Enable/disable shortcut
  function setShortcutEnabled(id: string, enabled: boolean) {
    const shortcut = shortcuts.value.find(s => s.id === id)
    if (shortcut) {
      shortcut.enabled = enabled
    }
  }

  // Register a new shortcut
  function registerShortcut(shortcut: KeyboardShortcut) {
    const existing = shortcuts.value.find(s => s.id === shortcut.id)
    if (existing) {
      Object.assign(existing, shortcut)
    } else {
      shortcuts.value.push(shortcut)
    }
  }

  // Unregister a shortcut
  function unregisterShortcut(id: string) {
    const index = shortcuts.value.findIndex(s => s.id === id)
    if (index !== -1) {
      shortcuts.value.splice(index, 1)
    }
  }

  // Set current scope
  function setScope(scope: string) {
    currentScope.value = scope
  }

  // Handle keydown event
  function handleKeyDown(event: KeyboardEvent) {
    if (!isListening.value) return

    // Ignore if input is focused
    const target = event.target as HTMLElement
    if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || target.isContentEditable) {
      // Allow some global shortcuts even in inputs
      if (!event.ctrlKey && !event.metaKey) {
        return
      }
    }

    const key = getNormalizedKey(event)

    // Add to sequence
    keySequence.value.push(key)

    // Clear sequence after timeout
    if (sequenceTimeout.value) {
      clearTimeout(sequenceTimeout.value)
    }
    sequenceTimeout.value = window.setTimeout(() => {
      keySequence.value = []
    }, 500)

    // Check for matching shortcuts
    for (const shortcut of enabledShortcuts.value) {
      if (matchesSequence(shortcut.keys, keySequence.value)) {
        event.preventDefault()
        shortcut.action()
        keySequence.value = []
        return
      }
    }

    // If sequence doesn't match any shortcut, check partial matches
    const hasPartialMatch = enabledShortcuts.value.some(s => {
      return s.keys.length > keySequence.value.length &&
        s.keys.slice(0, keySequence.value.length).every((k, i) => k === keySequence.value[i])
    })

    if (!hasPartialMatch) {
      keySequence.value = []
    }
  }

  // Get normalized key name
  function getNormalizedKey(event: KeyboardEvent): string {
    if (event.key === ' ') return 'Space'
    if (event.key === '/') return '/'
    if (event.key === '?') return '?'

    const key = event.key.length === 1 ? event.key.toUpperCase() : event.key

    // Normalize modifier keys
    if (event.ctrlKey && key !== 'Control') return `Ctrl+${key}`
    if (event.metaKey && key !== 'Meta') return `Cmd+${key}`
    if (event.altKey && key !== 'Alt') return `Alt+${key}`
    if (event.shiftKey && key.length === 1) return key

    return key
  }

  // Check if key sequence matches
  function matchesSequence(shortcutKeys: string[], pressedKeys: string[]): boolean {
    if (shortcutKeys.length !== pressedKeys.length) return false
    return shortcutKeys.every((k, i) => k.toLowerCase() === pressedKeys[i].toLowerCase())
  }

  // Toggle help
  function toggleHelp() {
    showHelp.value = !showHelp.value
  }

  // Format keys for display
  function formatKeys(keys: string[]): string {
    return keys.map(k => {
      const keyMap: Record<string, string> = {
        'Ctrl': '⌃',
        'Cmd': '⌘',
        'Alt': '⌥',
        'Shift': '⇧',
        'Enter': '↵',
        'Escape': 'Esc',
        'ArrowUp': '↑',
        'ArrowDown': '↓',
        'ArrowLeft': '←',
        'ArrowRight': '→',
        'Delete': 'Del',
        'Backspace': '⌫',
        'Space': '␣',
      }
      return keyMap[k] || k
    }).join(' + ')
  }

  // Lifecycle
  onMounted(() => {
    loadCustomShortcuts()
    document.addEventListener('keydown', handleKeyDown)
  })

  onUnmounted(() => {
    document.removeEventListener('keydown', handleKeyDown)
    if (sequenceTimeout.value) {
      clearTimeout(sequenceTimeout.value)
    }
  })

  return {
    shortcuts,
    categories,
    enabledShortcuts,
    isListening,
    showHelp,
    currentScope,
    loadCustomShortcuts,
    saveCustomShortcuts,
    updateShortcut,
    resetToDefaults,
    setShortcutEnabled,
    registerShortcut,
    unregisterShortcut,
    setScope,
    toggleHelp,
    formatKeys
  }
})