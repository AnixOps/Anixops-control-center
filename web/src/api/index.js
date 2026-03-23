import axios from 'axios'
import { useAuthStore } from '@/stores/auth'

// Use environment variable or default to Workers API
const baseURL = import.meta.env.VITE_API_URL || 'https://api.anixops.com/api/v1'

const api = axios.create({
  baseURL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor
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

// Response interceptor
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      const authStore = useAuthStore()
      authStore.logout()
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default api

// API methods
export const nodesApi = {
  list: () => api.get('/nodes'),
  get: (id) => api.get(`/nodes/${id}`),
  create: (data) => api.post('/nodes', data),
  update: (id, data) => api.put(`/nodes/${id}`, data),
  delete: (id) => api.delete(`/nodes/${id}`),
  stats: (id) => api.get(`/nodes/${id}/stats`)
}

export const agentsApi = {
  list: () => api.get('/agents'),
  connect: (data) => api.post('/agents/connect', data),
  disconnect: () => api.post('/agents/disconnect'),
  exec: (data) => api.post('/agents/exec', data),
  info: (id) => api.get(`/agents/${id}/info`)
}

export const usersApi = {
  list: (params) => api.get('/users', { params }),
  get: (id) => api.get(`/users/${id}`),
  create: (data) => api.post('/users', data),
  update: (id, data) => api.put(`/users/${id}`, data),
  delete: (id) => api.delete(`/users/${id}`),
  ban: (id) => api.post(`/users/${id}/ban`),
  unban: (id) => api.post(`/users/${id}/unban`)
}

export const playbooksApi = {
  list: () => api.get('/playbooks'),
  get: (name) => api.get(`/playbooks/${name}`),
  run: (data) => api.post('/playbooks/run', data),
  validate: (data) => api.post('/playbooks/validate', data)
}

export const dashboardApi = {
  get: () => api.get('/dashboard'),
  stats: () => api.get('/dashboard/stats')
}

export const pluginsApi = {
  list: () => api.get('/plugins'),
  get: (name) => api.get(`/plugins/${name}`),
  execute: (name, action, params) => api.post(`/plugins/${name}/execute`, { action, params }),
  status: (name) => api.get(`/plugins/${name}/status`),
  start: (name) => api.post(`/admin/plugins/${name}/start`),
  stop: (name) => api.post(`/admin/plugins/${name}/stop`)
}

export const logsApi = {
  list: (params) => api.get('/logs', { params })
}

export const settingsApi = {
  get: () => api.get('/settings'),
  update: (data) => api.put('/settings', data)
}

// AI Services
export const aiApi = {
  chat: (message, history = []) => api.post('/ai/chat', { message, history }),
  analyzeLog: (logContent) => api.post('/ai/analyze-log', { logContent }),
  opsAdvice: (context) => api.post('/ai/ops-advice', { context }),
  embedding: (text) => api.post('/ai/embedding', { text }),
  query: (query) => api.post('/ai/query', { query })
}

// Vectorize Services
export const vectorApi = {
  search: (embedding, options = {}) => api.post('/vectors/search', { embedding, ...options }),
  insert: (id, embedding, metadata) => api.post('/vectors', { id, embedding, metadata }),
  delete: (id) => api.delete(`/vectors/${id}`)
}

// IPFS Services
export const ipfsApi = {
  upload: (data) => api.post('/ipfs/upload', data, {
    headers: { 'Content-Type': 'multipart/form-data' }
  }),
  get: (cid) => api.get(`/ipfs/${cid}`)
}

// Web3 Services
export const web3Api = {
  challenge: (address) => api.post('/web3/challenge', { address }),
  verify: (address, signature, message) => api.post('/web3/verify', { address, signature, message }),
  audit: (auditData) => api.post('/web3/audit', auditData)
}