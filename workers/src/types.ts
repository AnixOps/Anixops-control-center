// Cloudflare Workers 环境类型定义

export interface Env {
  // 环境变量
  ENVIRONMENT: 'development' | 'production'
  JWT_SECRET: string
  JWT_EXPIRE: string
  API_KEY_SALT: string

  // D1 数据库
  DB: D1Database

  // KV 命名空间
  KV: KVNamespace

  // R2 存储桶
  R2: R2Bucket

  // Durable Objects
  WEBSOCKET_SERVER: DurableObjectNamespace
}

// 用户类型
export interface User {
  id: number
  email: string
  password_hash?: string
  role: 'admin' | 'operator' | 'viewer'
  auth_provider: 'local' | 'github' | 'google' | 'cloudflare'
  enabled: boolean
  last_login_at?: string
  created_at: string
  updated_at: string
}

// JWT Payload
export interface JWTPayload {
  sub: number
  email: string
  role: string
  iat: number
  exp: number
}

// 节点类型
export interface Node {
  id: number
  name: string
  host: string
  port: number
  status: 'online' | 'offline' | 'maintenance'
  last_seen?: string
  config?: string
  created_at: string
  updated_at: string
}

// Playbook 类型
export interface Playbook {
  id: number
  name: string
  storage_key: string
  description?: string
  created_at: string
  updated_at: string
}

// 审计日志类型
export interface AuditLog {
  id: number
  user_id?: number
  action: string
  resource?: string
  ip?: string
  user_agent?: string
  details?: string
  created_at: string
}

// API 响应类型
export interface ApiResponse<T = unknown> {
  success: boolean
  data?: T
  error?: string
  message?: string
}

// 分页响应
export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  per_page: number
  total_pages: number
}

// Hono 上下文变量
declare module 'hono' {
  interface ContextVariableMap {
    user: JWTPayload
  }
}