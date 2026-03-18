import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface DashboardWidget {
  id: string
  type: WidgetType
  title: string
  size: WidgetSize
  position: { x: number; y: number }
  config: Record<string, any>
  enabled: boolean
  refreshInterval?: number
}

export type WidgetType =
  | 'node-status'
  | 'traffic-chart'
  | 'user-stats'
  | 'system-metrics'
  | 'alerts'
  | 'quick-actions'
  | 'recent-activity'
  | 'plugin-status'
  | 'traffic-map'
  | 'performance-gauge'

export type WidgetSize = 'small' | 'medium' | 'large' | 'xlarge'

export interface WidgetTemplate {
  type: WidgetType
  name: string
  description: string
  defaultSize: WidgetSize
  defaultConfig: Record<string, any>
  icon: string
}

const STORAGE_KEY = 'dashboard-widgets'

// Widget templates
export const widgetTemplates: WidgetTemplate[] = [
  {
    type: 'node-status',
    name: 'Node Status',
    description: 'Overview of all node statuses',
    defaultSize: 'medium',
    defaultConfig: { showOffline: true, refreshInterval: 30000 },
    icon: 'server'
  },
  {
    type: 'traffic-chart',
    name: 'Traffic Chart',
    description: 'Real-time traffic visualization',
    defaultSize: 'large',
    defaultConfig: { timeRange: '24h', showUpload: true, showDownload: true },
    icon: 'chart-line'
  },
  {
    type: 'user-stats',
    name: 'User Statistics',
    description: 'User activity and statistics',
    defaultSize: 'medium',
    defaultConfig: { showActive: true, showBanned: false },
    icon: 'users'
  },
  {
    type: 'system-metrics',
    name: 'System Metrics',
    description: 'CPU, Memory, Disk usage',
    defaultSize: 'medium',
    defaultConfig: { metrics: ['cpu', 'memory', 'disk'] },
    icon: 'gauge'
  },
  {
    type: 'alerts',
    name: 'Active Alerts',
    description: 'Current system alerts',
    defaultSize: 'small',
    defaultConfig: { maxAlerts: 5, severity: 'all' },
    icon: 'bell'
  },
  {
    type: 'quick-actions',
    name: 'Quick Actions',
    description: 'Common action shortcuts',
    defaultSize: 'small',
    defaultConfig: { actions: ['deploy-node', 'run-playbook', 'manage-users'] },
    icon: 'bolt'
  },
  {
    type: 'recent-activity',
    name: 'Recent Activity',
    description: 'Latest system events',
    defaultSize: 'medium',
    defaultConfig: { maxItems: 10, types: 'all' },
    icon: 'clock'
  },
  {
    type: 'plugin-status',
    name: 'Plugin Status',
    description: 'Status of installed plugins',
    defaultSize: 'small',
    defaultConfig: { showInactive: false },
    icon: 'puzzle'
  },
  {
    type: 'traffic-map',
    name: 'Traffic Map',
    description: 'Geographic traffic distribution',
    defaultSize: 'xlarge',
    defaultConfig: { showLabels: true, animate: true },
    icon: 'map'
  },
  {
    type: 'performance-gauge',
    name: 'Performance Gauge',
    description: 'Overall system performance score',
    defaultSize: 'small',
    defaultConfig: { metrics: ['latency', 'throughput', 'errors'] },
    icon: 'speedometer'
  }
]

// Default dashboard layout
const defaultWidgets: DashboardWidget[] = [
  { id: 'node-status-1', type: 'node-status', title: 'Node Status', size: 'medium', position: { x: 0, y: 0 }, config: {}, enabled: true },
  { id: 'traffic-chart-1', type: 'traffic-chart', title: 'Traffic Overview', size: 'large', position: { x: 2, y: 0 }, config: {}, enabled: true },
  { id: 'alerts-1', type: 'alerts', title: 'Active Alerts', size: 'small', position: { x: 0, y: 1 }, config: {}, enabled: true },
  { id: 'quick-actions-1', type: 'quick-actions', title: 'Quick Actions', size: 'small', position: { x: 1, y: 1 }, config: {}, enabled: true },
  { id: 'system-metrics-1', type: 'system-metrics', title: 'System Metrics', size: 'medium', position: { x: 2, y: 1 }, config: {}, enabled: true }
]

export const useDashboardWidgetsStore = defineStore('dashboardWidgets', () => {
  const widgets = ref<DashboardWidget[]>([])
  const isEditing = ref(false)
  const editMode = ref<'layout' | 'settings'>('layout')
  const selectedWidgetId = ref<string | null>(null)
  const isLoading = ref(false)

  // Computed
  const enabledWidgets = computed(() => {
    return widgets.value.filter(w => w.enabled).sort((a, b) => {
      if (a.position.y !== b.position.y) return a.position.y - b.position.y
      return a.position.x - b.position.x
    })
  })

  const selectedWidget = computed(() => {
    return widgets.value.find(w => w.id === selectedWidgetId.value) || null
  })

  const availableWidgets = computed(() => {
    const usedTypes = new Set(widgets.value.map(w => w.type))
    return widgetTemplates.filter(t => !usedTypes.has(t.type))
  })

  // Load widgets from storage
  function loadWidgets() {
    isLoading.value = true
    try {
      const stored = localStorage.getItem(STORAGE_KEY)
      if (stored) {
        widgets.value = JSON.parse(stored)
      } else {
        widgets.value = [...defaultWidgets]
      }
    } catch (e) {
      console.error('Failed to load widgets:', e)
      widgets.value = [...defaultWidgets]
    } finally {
      isLoading.value = false
    }
  }

  // Save widgets to storage
  function saveWidgets() {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(widgets.value))
  }

  // Add widget
  function addWidget(type: WidgetType, config?: Partial<DashboardWidget>) {
    const template = widgetTemplates.find(t => t.type === type)
    if (!template) return null

    const widget: DashboardWidget = {
      id: `${type}-${Date.now()}`,
      type,
      title: template.name,
      size: template.defaultSize,
      position: findEmptyPosition(template.defaultSize),
      config: { ...template.defaultConfig, ...config?.config },
      enabled: true,
      ...config
    }

    widgets.value.push(widget)
    saveWidgets()
    return widget
  }

  // Remove widget
  function removeWidget(id: string) {
    const index = widgets.value.findIndex(w => w.id === id)
    if (index !== -1) {
      widgets.value.splice(index, 1)
      saveWidgets()
    }
  }

  // Update widget
  function updateWidget(id: string, updates: Partial<DashboardWidget>) {
    const widget = widgets.value.find(w => w.id === id)
    if (widget) {
      Object.assign(widget, updates)
      saveWidgets()
    }
  }

  // Move widget
  function moveWidget(id: string, position: { x: number; y: number }) {
    const widget = widgets.value.find(w => w.id === id)
    if (widget) {
      widget.position = position
      saveWidgets()
    }
  }

  // Resize widget
  function resizeWidget(id: string, size: WidgetSize) {
    const widget = widgets.value.find(w => w.id === id)
    if (widget) {
      widget.size = size
      saveWidgets()
    }
  }

  // Toggle widget
  function toggleWidget(id: string) {
    const widget = widgets.value.find(w => w.id === id)
    if (widget) {
      widget.enabled = !widget.enabled
      saveWidgets()
    }
  }

  // Select widget
  function selectWidget(id: string | null) {
    selectedWidgetId.value = id
  }

  // Toggle edit mode
  function toggleEditMode(mode?: 'layout' | 'settings') {
    if (isEditing.value && editMode.value === mode) {
      isEditing.value = false
      selectedWidgetId.value = null
    } else {
      isEditing.value = true
      editMode.value = mode || 'layout'
    }
  }

  // Reset to default layout
  function resetToDefault() {
    widgets.value = [...defaultWidgets]
    saveWidgets()
  }

  // Find empty position
  function findEmptyPosition(size: WidgetSize): { x: number; y: number } {
    const gridSize = 4 // 4 columns
    const sizeWidth: Record<WidgetSize, number> = {
      small: 1,
      medium: 2,
      large: 3,
      xlarge: 4
    }

    const occupied = new Set<string>()
    widgets.value.forEach(w => {
      const width = sizeWidth[w.size]
      for (let dx = 0; dx < width; dx++) {
        occupied.add(`${w.position.x + dx},${w.position.y}`)
      }
    })

    const width = sizeWidth[size]
    for (let y = 0; y < 10; y++) {
      for (let x = 0; x <= gridSize - width; x++) {
        let canPlace = true
        for (let dx = 0; dx < width; dx++) {
          if (occupied.has(`${x + dx},${y}`)) {
            canPlace = false
            break
          }
        }
        if (canPlace) {
          return { x, y }
        }
      }
    }

    return { x: 0, y: 0 }
  }

  // Export layout
  function exportLayout(): string {
    return JSON.stringify({
      version: '1.0',
      widgets: widgets.value
    }, null, 2)
  }

  // Import layout
  function importLayout(json: string): boolean {
    try {
      const data = JSON.parse(json)
      if (data.widgets && Array.isArray(data.widgets)) {
        widgets.value = data.widgets
        saveWidgets()
        return true
      }
      return false
    } catch {
      return false
    }
  }

  // Initialize
  loadWidgets()

  return {
    widgets,
    enabledWidgets,
    availableWidgets,
    selectedWidget,
    isEditing,
    editMode,
    selectedWidgetId,
    isLoading,
    loadWidgets,
    saveWidgets,
    addWidget,
    removeWidget,
    updateWidget,
    moveWidget,
    resizeWidget,
    toggleWidget,
    selectWidget,
    toggleEditMode,
    resetToDefault,
    exportLayout,
    importLayout
  }
})