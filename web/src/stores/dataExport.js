import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface ExportData {
  version: string
  timestamp: string
  platform: string
  data: {
    nodes?: any[]
    users?: any[]
    plugins?: any[]
    settings?: any
    playbooks?: any[]
  }
  metadata: {
    exportedBy: string
    exportedAt: string
    includeSensitive: boolean
  }
}

export type ExportFormat = 'json' | 'csv' | 'yaml'
export type ImportMode = 'merge' | 'replace' | 'append'

export const useDataExportStore = defineStore('dataExport', () => {
  const isExporting = ref(false)
  const isImporting = ref(false)
  const exportHistory = ref<ExportData[]>([])
  const lastExport = ref<string | null>(null)
  const lastImport = ref<string | null>(null)

  // Computed
  const hasHistory = computed(() => exportHistory.value.length > 0)

  // Export data
  async function exportData(
    options: {
      includeNodes?: boolean
      includeUsers?: boolean
      includePlugins?: boolean
      includeSettings?: boolean
      includePlaybooks?: boolean
      includeSensitive?: boolean
      format?: ExportFormat
    } = {}
  ): Promise<Blob> {
    isExporting.value = true

    try {
      const data: ExportData = {
        version: '1.1.0',
        timestamp: new Date().toISOString(),
        platform: 'web',
        data: {},
        metadata: {
          exportedBy: 'AnixOps Control Center',
          exportedAt: new Date().toISOString(),
          includeSensitive: options.includeSensitive || false
        }
      }

      // Collect data based on options
      if (options.includeNodes) {
        // Get nodes from API or store
        data.data.nodes = []
      }

      if (options.includeUsers) {
        data.data.users = []
      }

      if (options.includePlugins) {
        data.data.plugins = []
      }

      if (options.includeSettings) {
        data.data.settings = {}
      }

      if (options.includePlaybooks) {
        data.data.playbooks = []
      }

      // Convert to requested format
      const format = options.format || 'json'
      let content: string
      let mimeType: string

      switch (format) {
        case 'json':
          content = JSON.stringify(data, null, 2)
          mimeType = 'application/json'
          break
        case 'csv':
          content = convertToCSV(data)
          mimeType = 'text/csv'
          break
        case 'yaml':
          content = convertToYAML(data)
          mimeType = 'text/yaml'
          break
        default:
          content = JSON.stringify(data, null, 2)
          mimeType = 'application/json'
      }

      // Update history
      exportHistory.value.unshift(data)
      lastExport.value = data.timestamp

      // Keep only last 10 exports in history
      if (exportHistory.value.length > 10) {
        exportHistory.value = exportHistory.value.slice(0, 10)
      }

      return new Blob([content], { type: mimeType })
    } finally {
      isExporting.value = false
    }
  }

  // Import data
  async function importData(
    file: File,
    options: {
      mode?: ImportMode
      validateOnly?: boolean
    } = {}
  ): Promise<{ success: boolean; message: string; stats?: any }> {
    isImporting.value = true

    try {
      const content = await readFileContent(file)
      let data: ExportData

      // Parse based on file type
      if (file.name.endsWith('.json')) {
        data = JSON.parse(content)
      } else if (file.name.endsWith('.yaml') || file.name.endsWith('.yml')) {
        data = parseYAML(content)
      } else {
        return { success: false, message: 'Unsupported file format' }
      }

      // Validate data
      const validation = validateImportData(data)
      if (!validation.valid) {
        return { success: false, message: validation.error || 'Invalid data format' }
      }

      if (options.validateOnly) {
        return { success: true, message: 'Data is valid', stats: getImportStats(data) }
      }

      // Import based on mode
      const mode = options.mode || 'merge'
      const stats = {
        nodes: 0,
        users: 0,
        plugins: 0,
        settings: false,
        playbooks: 0
      }

      if (data.data.nodes) {
        stats.nodes = await importNodes(data.data.nodes, mode)
      }

      if (data.data.users) {
        stats.users = await importUsers(data.data.users, mode)
      }

      if (data.data.plugins) {
        stats.plugins = await importPlugins(data.data.plugins, mode)
      }

      if (data.data.settings) {
        stats.settings = await importSettings(data.data.settings, mode)
      }

      if (data.data.playbooks) {
        stats.playbooks = await importPlaybooks(data.data.playbooks, mode)
      }

      lastImport.value = new Date().toISOString()

      return { success: true, message: 'Import completed', stats }
    } catch (error) {
      return { success: false, message: `Import failed: ${error}` }
    } finally {
      isImporting.value = false
    }
  }

  // Download exported file
  function downloadExport(blob: Blob, filename?: string) {
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = filename || `anixops-export-${new Date().toISOString().split('T')[0]}.json`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
  }

  // Helper functions
  async function readFileContent(file: File): Promise<string> {
    return new Promise((resolve, reject) => {
      const reader = new FileReader()
      reader.onload = () => resolve(reader.result as string)
      reader.onerror = () => reject(reader.error)
      reader.readAsText(file)
    })
  }

  function convertToCSV(data: ExportData): string {
    // Simple CSV conversion for nodes/users
    const lines: string[] = []

    if (data.data.nodes && data.data.nodes.length > 0) {
      lines.push('# Nodes')
      lines.push(Object.keys(data.data.nodes[0]).join(','))
      data.data.nodes.forEach(node => {
        lines.push(Object.values(node).map(v => `"${v}"`).join(','))
      })
      lines.push('')
    }

    if (data.data.users && data.data.users.length > 0) {
      lines.push('# Users')
      lines.push(Object.keys(data.data.users[0]).join(','))
      data.data.users.forEach(user => {
        lines.push(Object.values(user).map(v => `"${v}"`).join(','))
      })
    }

    return lines.join('\n')
  }

  function convertToYAML(data: ExportData): string {
    // Simple YAML conversion
    let yaml = `version: "${data.version}"
timestamp: "${data.timestamp}"
platform: "${data.platform}"
metadata:
  exportedBy: "${data.metadata.exportedBy}"
  exportedAt: "${data.metadata.exportedAt}"
  includeSensitive: ${data.metadata.includeSensitive}
data:
`

    if (data.data.nodes) {
      yaml += `  nodes:\n`
      data.data.nodes.forEach(node => {
        yaml += `    - ${JSON.stringify(node)}\n`
      })
    }

    if (data.data.users) {
      yaml += `  users:\n`
      data.data.users.forEach(user => {
        yaml += `    - ${JSON.stringify(user)}\n`
      })
    }

    return yaml
  }

  function parseYAML(content: string): ExportData {
    // Basic YAML parsing (in production, use a library)
    // For now, assume JSON-like structure
    return JSON.parse(content)
  }

  function validateImportData(data: any): { valid: boolean; error?: string } {
    if (!data) return { valid: false, error: 'No data provided' }
    if (!data.version) return { valid: false, error: 'Missing version' }
    if (!data.data) return { valid: false, error: 'Missing data section' }

    return { valid: true }
  }

  function getImportStats(data: ExportData): any {
    return {
      nodes: data.data.nodes?.length || 0,
      users: data.data.users?.length || 0,
      plugins: data.data.plugins?.length || 0,
      playbooks: data.data.playbooks?.length || 0,
      hasSettings: !!data.data.settings
    }
  }

  async function importNodes(nodes: any[], mode: ImportMode): Promise<number> {
    // Implementation would call API
    return nodes.length
  }

  async function importUsers(users: any[], mode: ImportMode): Promise<number> {
    return users.length
  }

  async function importPlugins(plugins: any[], mode: ImportMode): Promise<number> {
    return plugins.length
  }

  async function importSettings(settings: any, mode: ImportMode): Promise<boolean> {
    return true
  }

  async function importPlaybooks(playbooks: any[], mode: ImportMode): Promise<number> {
    return playbooks.length
  }

  return {
    isExporting,
    isImporting,
    exportHistory,
    lastExport,
    lastImport,
    hasHistory,
    exportData,
    importData,
    downloadExport
  }
})