import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_URL || '/api/v1'

// Create axios instance
const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor - add auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// Response interceptor - handle errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response) {
      // Handle 401 Unauthorized
      if (error.response.status === 401) {
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        window.location.href = '/login'
      }
      // Return structured error
      return Promise.reject({
        status: error.response.status,
        message: error.response.data?.message || error.response.data?.error || 'An error occurred',
        errors: error.response.data?.errors || {}
      })
    }
    return Promise.reject({
      status: 0,
      message: error.message || 'Network error',
      errors: {}
    })
  }
)

export default api

// ============================================
// Auth API
// ============================================
export const authApi = {
  login: (email, password) => api.post('/auth/login', { email, password }),
  logout: () => api.post('/auth/logout'),
  refresh: () => api.post('/auth/refresh'),
  register: (data) => api.post('/auth/register', data),
  me: () => api.get('/auth/me'),
  updatePassword: (data) => api.put('/auth/password', data),
  forgotPassword: (email) => api.post('/auth/forgot-password', { email }),
  resetPassword: (token, password) => api.post('/auth/reset-password', { token, password })
}

// ============================================
// Nodes API
// ============================================
export const nodesApi = {
  list: (params) => api.get('/nodes', { params }),
  get: (id) => api.get(`/nodes/${id}`),
  create: (data) => api.post('/nodes', data),
  update: (id, data) => api.put(`/nodes/${id}`, data),
  delete: (id) => api.delete(`/nodes/${id}`),
  batchDelete: (ids) => api.post('/nodes/batch-delete', { ids }),
  start: (id) => api.post(`/nodes/${id}/start`),
  stop: (id) => api.post(`/nodes/${id}/stop`),
  restart: (id) => api.post(`/nodes/${id}/restart`),
  stats: (id) => api.get(`/nodes/${id}/stats`),
  logs: (id, params) => api.get(`/nodes/${id}/logs`, { params }),
  testConnection: (id) => api.post(`/nodes/${id}/test`),
  sync: (id) => api.post(`/nodes/${id}/sync`)
}

// ============================================
// Plugins API
// ============================================
export const pluginsApi = {
  list: () => api.get('/plugins'),
  get: (name) => api.get(`/plugins/${name}`),
  execute: (name, action, params) => api.post(`/plugins/${name}/execute`, { action, params }),
  status: (name) => api.get(`/plugins/${name}/status`),
  config: (name) => api.get(`/plugins/${name}/config`),
  updateConfig: (name, config) => api.put(`/plugins/${name}/config`, config),
  enable: (name) => api.post(`/plugins/${name}/enable`),
  disable: (name) => api.post(`/plugins/${name}/disable`),
  start: (name) => api.post(`/plugins/${name}/start`),
  stop: (name) => api.post(`/plugins/${name}/stop`),
  restart: (name) => api.post(`/plugins/${name}/restart`),
  logs: (name, params) => api.get(`/plugins/${name}/logs`, { params })
}

// ============================================
// Playbooks API
// ============================================
export const playbooksApi = {
  list: (params) => api.get('/playbooks', { params }),
  get: (id) => api.get(`/playbooks/${id}`),
  create: (data) => api.post('/playbooks', data),
  update: (id, data) => api.put(`/playbooks/${id}`, data),
  delete: (id) => api.delete(`/playbooks/${id}`),
  duplicate: (id) => api.post(`/playbooks/${id}/duplicate`),
  run: (id, params) => api.post(`/playbooks/${id}/run`, params),
  runDirect: (data) => api.post('/playbooks/run', data),
  logs: (id) => api.get(`/playbooks/${id}/logs`),
  stop: (id) => api.post(`/playbooks/${id}/stop`),
  schedule: (id, cron) => api.post(`/playbooks/${id}/schedule`, { cron }),
  cancelSchedule: (id) => api.delete(`/playbooks/${id}/schedule`),
  validate: (data) => api.post('/playbooks/validate', data),
  templates: () => api.get('/playbooks/templates')
}

// ============================================
// Users API
// ============================================
export const usersApi = {
  list: (params) => api.get('/users', { params }),
  get: (id) => api.get(`/users/${id}`),
  create: (data) => api.post('/users', data),
  update: (id, data) => api.put(`/users/${id}`, data),
  delete: (id) => api.delete(`/users/${id}`),
  ban: (id) => api.post(`/users/${id}/ban`),
  unban: (id) => api.post(`/users/${id}/unban`),
  resetPassword: (id) => api.post(`/users/${id}/reset-password`),
  updateRole: (id, role) => api.put(`/users/${id}/role`, { role }),
  export: (params) => api.get('/users/export', { params, responseType: 'blob' }),
  import: (data) => api.post('/users/import', data)
}

// ============================================
// Agents API
// ============================================
export const agentsApi = {
  list: (params) => api.get('/agents', { params }),
  get: (id) => api.get(`/agents/${id}`),
  create: (data) => api.post('/agents', data),
  update: (id, data) => api.put(`/agents/${id}`, data),
  delete: (id) => api.delete(`/agents/${id}`),
  connect: (id) => api.post(`/agents/${id}/connect`),
  disconnect: (id) => api.post(`/agents/${id}/disconnect`),
  execute: (id, command) => api.post(`/agents/${id}/execute`, { command }),
  stats: (id) => api.get(`/agents/${id}/stats`),
  logs: (id, params) => api.get(`/agents/${id}/logs`, { params }),
  restart: (id) => api.post(`/agents/${id}/restart`),
  update: (id) => api.post(`/agents/${id}/update`)
}

// ============================================
// System API
// ============================================
export const systemApi = {
  info: () => api.get('/system/info'),
  health: () => api.get('/system/health'),
  metrics: () => api.get('/system/metrics'),
  logs: (params) => api.get('/system/logs', { params }),
  config: () => api.get('/system/config'),
  updateConfig: (config) => api.put('/system/config', config),
  backup: () => api.post('/system/backup'),
  restore: (data) => api.post('/system/restore', data),
  restart: () => api.post('/system/restart'),
  update: () => api.post('/system/update'),
  checkUpdate: () => api.get('/system/update/check')
}

// ============================================
// Logs API
// ============================================
export const logsApi = {
  list: (params) => api.get('/logs', { params }),
  stream: (params) => {
    const token = localStorage.getItem('token')
    const query = new URLSearchParams(params).toString()
    const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    return `${wsProtocol}//${window.location.host}${API_BASE_URL}/logs/stream?token=${token}&${query}`
  },
  export: (params) => api.get('/logs/export', { params, responseType: 'blob' }),
  clear: (before) => api.delete('/logs', { params: { before } }),
  level: () => api.get('/logs/level'),
  setLevel: (level) => api.put('/logs/level', { level })
}

// ============================================
// Settings API
// ============================================
export const settingsApi = {
  get: () => api.get('/settings'),
  update: (settings) => api.put('/settings', settings),
  testEmail: (email) => api.post('/settings/test-email', { email }),
  testConnection: () => api.post('/settings/test-connection'),
  backup: () => api.get('/settings/backup'),
  restore: (data) => api.post('/settings/restore', data),
  reset: () => api.post('/settings/reset')
}

// ============================================
// Dashboard API
// ============================================
export const dashboardApi = {
  stats: () => api.get('/dashboard/stats'),
  activities: (limit = 10) => api.get('/dashboard/activities', { params: { limit } }),
  alerts: () => api.get('/dashboard/alerts'),
  traffic: (period = '24h') => api.get('/dashboard/traffic', { params: { period } }),
  nodes: () => api.get('/dashboard/nodes'),
  charts: (type, period) => api.get('/dashboard/charts', { params: { type, period } })
}

// ============================================
// Roles API
// ============================================
export const rolesApi = {
  list: () => api.get('/roles'),
  get: (id) => api.get(`/roles/${id}`),
  create: (data) => api.post('/roles', data),
  update: (id, data) => api.put(`/roles/${id}`, data),
  delete: (id) => api.delete(`/roles/${id}`),
  permissions: (id) => api.get(`/roles/${id}/permissions`),
  updatePermissions: (id, permissions) => api.put(`/roles/${id}/permissions`, { permissions })
}

// ============================================
// Audit API
// ============================================
export const auditApi = {
  list: (params) => api.get('/audit', { params }),
  get: (id) => api.get(`/audit/${id}`),
  export: (params) => api.get('/audit/export', { params, responseType: 'blob' }),
  clear: (before) => api.delete('/audit', { params: { before } })
}