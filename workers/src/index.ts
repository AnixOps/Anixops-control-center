import { Hono } from 'hono'
import { cors } from 'hono/cors'
import { logger } from 'hono/logger'
import { prettyJSON } from 'hono/pretty-json'
import type { Env } from './types'

// Handlers
import { healthHandler, readinessHandler } from './handlers/health'
import { loginHandler, registerHandler, refreshHandler, logoutHandler, meHandler } from './handlers/auth'
import { listNodesHandler, getNodeHandler, createNodeHandler, updateNodeHandler, deleteNodeHandler } from './handlers/nodes'
import { listPlaybooksHandler, getPlaybookHandler, uploadPlaybookHandler, runPlaybookHandler } from './handlers/playbooks'
import { listPluginsHandler, getPluginHandler, executePluginHandler } from './handlers/plugins'
import { dashboardHandler, statsHandler } from './handlers/dashboard'
import { listAuditLogsHandler } from './handlers/audit'
import { listUsersHandler, getUserHandler, createUserHandler, updateUserHandler, deleteUserHandler } from './handlers/users'

// Middleware
import { authMiddleware, rbacMiddleware } from './middleware/auth'
import { rateLimitMiddleware } from './middleware/rate-limit'

// 创建应用
const app = new Hono<{ Bindings: Env }>()

// 全局中间件
app.use('*', logger())
app.use('*', prettyJSON())
app.use('*', cors({
  origin: (origin) => {
    // 允许的域名
    const allowed = [
      'http://localhost:3000',
      'http://localhost:5173',
      'https://anixops.pages.dev',
      'https://anixops.dev',
    ]
    if (allowed.includes(origin)) return origin
    return allowed[0]
  },
  allowMethods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS', 'PATCH'],
  allowHeaders: ['Content-Type', 'Authorization', 'X-API-Key'],
  exposeHeaders: ['X-Total-Count'],
  credentials: true,
  maxAge: 86400,
}))

// ==================== 公开路由 ====================

// 健康检查
app.get('/health', healthHandler)
app.get('/readiness', readinessHandler)

// 认证 (公开)
app.post('/api/v1/auth/login', rateLimitMiddleware({ windowMs: 60000, max: 5 }), loginHandler)
app.post('/api/v1/auth/register', rateLimitMiddleware({ windowMs: 60000, max: 3 }), registerHandler)
app.post('/api/v1/auth/refresh', refreshHandler)
app.post('/api/v1/auth/logout', logoutHandler)

// ==================== 受保护路由 ====================

// 用户信息
app.get('/api/v1/users/me', authMiddleware, meHandler)

// 用户管理 (需要管理员权限)
app.get('/api/v1/users', authMiddleware, rbacMiddleware(['admin']), listUsersHandler)
app.get('/api/v1/users/:id', authMiddleware, rbacMiddleware(['admin']), getUserHandler)
app.post('/api/v1/users', authMiddleware, rbacMiddleware(['admin']), createUserHandler)
app.put('/api/v1/users/:id', authMiddleware, rbacMiddleware(['admin']), updateUserHandler)
app.delete('/api/v1/users/:id', authMiddleware, rbacMiddleware(['admin']), deleteUserHandler)

// 节点管理
app.get('/api/v1/nodes', authMiddleware, listNodesHandler)
app.get('/api/v1/nodes/:id', authMiddleware, getNodeHandler)
app.post('/api/v1/nodes', authMiddleware, rbacMiddleware(['admin', 'operator']), createNodeHandler)
app.put('/api/v1/nodes/:id', authMiddleware, rbacMiddleware(['admin', 'operator']), updateNodeHandler)
app.delete('/api/v1/nodes/:id', authMiddleware, rbacMiddleware(['admin']), deleteNodeHandler)

// Playbook 管理
app.get('/api/v1/playbooks', authMiddleware, listPlaybooksHandler)
app.get('/api/v1/playbooks/:name', authMiddleware, getPlaybookHandler)
app.post('/api/v1/playbooks', authMiddleware, rbacMiddleware(['admin', 'operator']), uploadPlaybookHandler)
app.post('/api/v1/playbooks/:name/run', authMiddleware, rbacMiddleware(['admin', 'operator']), runPlaybookHandler)

// 插件管理
app.get('/api/v1/plugins', authMiddleware, listPluginsHandler)
app.get('/api/v1/plugins/:name', authMiddleware, getPluginHandler)
app.post('/api/v1/plugins/:name/execute', authMiddleware, rbacMiddleware(['admin', 'operator']), executePluginHandler)

// Dashboard
app.get('/api/v1/dashboard', authMiddleware, dashboardHandler)
app.get('/api/v1/dashboard/stats', authMiddleware, statsHandler)

// 审计日志
app.get('/api/v1/audit-logs', authMiddleware, rbacMiddleware(['admin']), listAuditLogsHandler)

// ==================== WebSocket ====================

app.get('/api/v1/ws', async (c) => {
  const id = c.env.WEBSOCKET_SERVER.idFromName('global')
  const stub = c.env.WEBSOCKET_SERVER.get(id)
  return stub.fetch(c.req.raw)
})

// ==================== 错误处理 ====================

app.notFound((c) => {
  return c.json({ success: false, error: 'Not Found' }, 404)
})

app.onError((err, c) => {
  console.error('Error:', err)

  // 开发环境返回详细错误
  if (c.env.ENVIRONMENT === 'development') {
    return c.json({
      success: false,
      error: err.message,
      stack: err.stack,
    }, 500)
  }

  return c.json({ success: false, error: 'Internal Server Error' }, 500)
})

// ==================== Durable Object ====================

export class WebSocketServer {
  private state: DurableObjectState
  private sessions: Map<WebSocket, { userId?: number }>

  constructor(state: DurableObjectState) {
    this.state = state
    this.sessions = new Map()
  }

  async fetch(request: Request): Promise<Response> {
    const url = new URL(request.url)

    if (url.pathname === '/api/v1/ws') {
      const { 0: client, 1: server } = new WebSocketPair()

      this.handleSession(server)

      return new Response(null, { status: 101, webSocket: client })
    }

    return new Response('Not Found', { status: 404 })
  }

  private handleSession(ws: WebSocket) {
    ws.accept()
    this.sessions.set(ws, {})

    ws.addEventListener('message', async (event) => {
      try {
        const data = JSON.parse(event.data as string)
        await this.handleMessage(ws, data)
      } catch (err) {
        ws.send(JSON.stringify({ error: 'Invalid message format' }))
      }
    })

    ws.addEventListener('close', () => {
      this.sessions.delete(ws)
    })

    ws.addEventListener('error', () => {
      this.sessions.delete(ws)
    })
  }

  private async handleMessage(ws: WebSocket, data: { type: string; payload?: unknown }) {
    switch (data.type) {
      case 'ping':
        ws.send(JSON.stringify({ type: 'pong' }))
        break

      case 'subscribe':
        // 订阅节点状态更新
        ws.send(JSON.stringify({ type: 'subscribed', channel: data.payload }))
        break

      case 'unsubscribe':
        // 取消订阅
        ws.send(JSON.stringify({ type: 'unsubscribed', channel: data.payload }))
        break

      default:
        ws.send(JSON.stringify({ error: 'Unknown message type' }))
    }
  }

  // 广播消息给所有连接
  broadcast(message: unknown) {
    for (const [ws] of this.sessions) {
      try {
        ws.send(JSON.stringify(message))
      } catch {
        this.sessions.delete(ws)
      }
    }
  }
}

// 导出
export default app