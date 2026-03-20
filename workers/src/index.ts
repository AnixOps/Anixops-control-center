import { Hono } from 'hono'
import { cors } from 'hono/cors'
import { logger } from 'hono/logger'
import { prettyJSON } from 'hono/pretty-json'
import type { Env } from './types'

// Handlers
import { healthHandler, readinessHandler } from './handlers/health'
import { loginHandler, registerHandler, refreshHandler, logoutHandler, meHandler } from './handlers/auth'
import { listNodesHandler, getNodeHandler, createNodeHandler, updateNodeHandler, deleteNodeHandler, startNodeHandler, stopNodeHandler, restartNodeHandler, getNodeStatsHandler, getNodeLogsHandler, testNodeConnectionHandler, syncNodeHandler, bulkActionHandler } from './handlers/nodes'
import { listPlaybooksHandler, getPlaybookHandler, uploadPlaybookHandler, deletePlaybookHandler, listBuiltInPlaybooksHandler, getPlaybookCategoriesHandler, syncBuiltInPlaybooksHandler } from './handlers/playbooks'
import { listPluginsHandler, getPluginHandler, executePluginHandler } from './handlers/plugins'
import { dashboardHandler, statsHandler } from './handlers/dashboard'
import { listAuditLogsHandler } from './handlers/audit'
import { listUsersHandler, getUserHandler, createUserHandler, updateUserHandler, deleteUserHandler, changePasswordHandler, getCurrentUserHandler, updateCurrentUserHandler, listApiTokensHandler, createApiTokenHandler, deleteApiTokenHandler, listSessionsHandler, deleteOtherSessionsHandler, getUserLockoutHandler, unlockUserHandler } from './handlers/users'
import { testConnectionHandler, importServerHandler, detectServerTypeHandler } from './handlers/ssh'
import { listNotificationsHandler, markNotificationReadHandler, markAllNotificationsReadHandler, deleteNotificationHandler, createNotificationHandler, getUnreadCountHandler } from './handlers/notifications'
import { listTasksHandler, getTaskHandler, createTaskHandler, cancelTaskHandler, retryTaskHandler, getTaskLogsHandler } from './handlers/tasks'
import { listSchedulesHandler, getScheduleHandler, createScheduleHandler, updateScheduleHandler, deleteScheduleHandler, toggleScheduleHandler, runScheduleNowHandler } from './handlers/schedules'
import { listNodeGroupsHandler, getNodeGroupHandler, createNodeGroupHandler, updateNodeGroupHandler, deleteNodeGroupHandler, addNodesToGroupHandler, removeNodesFromGroupHandler } from './handlers/node-groups'
import { sseHandler, sseSubscribeHandler, sseUnsubscribeHandler, sseStatusHandler } from './handlers/sse'
import { createBackupHandler, listBackupsHandler, getBackupHandler, deleteBackupHandler, downloadBackupHandler, restoreBackupHandler, cleanupBackupsHandler, backupStatusHandler } from './handlers/backup'
import {
  registerAgentHandler,
  agentHeartbeatHandler,
  agentMetricsHandler,
  agentCommandResultHandler,
  sendAgentCommandHandler,
  getAgentMetricsHandler,
  generateInstallScriptHandler,
} from './handlers/agents'
import {
  getMFAStatusHandler,
  setupMFAHandler,
  enableMFAHandler,
  disableMFAHandler,
  verifyMFAHandler,
  regenerateRecoveryCodesHandler,
  adminDisableMFAHandler,
} from './handlers/mfa'
import { batchOperationsHandler, bulkNodeStatusHandler } from './handlers/batch'
import { prometheusMetricsHandler, detailedHealthHandler, readinessHandler as k8sReadinessHandler, livenessHandler } from './handlers/metrics'
import { cacheMiddleware } from './middleware/cache'

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
app.get('/health/detailed', detailedHealthHandler)
app.get('/readiness', readinessHandler)
app.get('/liveness', livenessHandler)
app.get('/metrics', prometheusMetricsHandler)

// 认证 (公开)
app.post('/api/v1/auth/login', rateLimitMiddleware({ windowMs: 60000, max: 5 }), loginHandler)
app.post('/api/v1/auth/register', rateLimitMiddleware({ windowMs: 60000, max: 3 }), registerHandler)
app.post('/api/v1/auth/refresh', refreshHandler)
app.post('/api/v1/auth/logout', logoutHandler)

// ==================== 受保护路由 ====================

// 用户信息
app.get('/api/v1/users/me', authMiddleware, getCurrentUserHandler)
app.put('/api/v1/users/me', authMiddleware, updateCurrentUserHandler)
app.put('/api/v1/auth/password', authMiddleware, changePasswordHandler)

// API Tokens
app.get('/api/v1/users/me/tokens', authMiddleware, listApiTokensHandler)
app.post('/api/v1/users/me/tokens', authMiddleware, createApiTokenHandler)
app.delete('/api/v1/users/me/tokens/:id', authMiddleware, deleteApiTokenHandler)

// Sessions
app.get('/api/v1/users/me/sessions', authMiddleware, listSessionsHandler)
app.delete('/api/v1/users/me/sessions/others', authMiddleware, deleteOtherSessionsHandler)

// 用户管理 (需要管理员权限)
app.get('/api/v1/users', authMiddleware, rbacMiddleware(['admin']), listUsersHandler)
app.get('/api/v1/users/:id', authMiddleware, rbacMiddleware(['admin']), getUserHandler)
app.post('/api/v1/users', authMiddleware, rbacMiddleware(['admin']), createUserHandler)
app.put('/api/v1/users/:id', authMiddleware, rbacMiddleware(['admin']), updateUserHandler)
app.delete('/api/v1/users/:id', authMiddleware, rbacMiddleware(['admin']), deleteUserHandler)

// 账户锁定管理 (需要管理员权限)
app.get('/api/v1/users/:id/lockout', authMiddleware, rbacMiddleware(['admin']), getUserLockoutHandler)
app.post('/api/v1/users/:id/unlock', authMiddleware, rbacMiddleware(['admin']), unlockUserHandler)

// SSH导入
app.post('/api/v1/ssh/test', authMiddleware, rbacMiddleware(['admin', 'operator']), testConnectionHandler)
app.post('/api/v1/ssh/import', authMiddleware, rbacMiddleware(['admin', 'operator']), importServerHandler)
app.post('/api/v1/ssh/detect', authMiddleware, rbacMiddleware(['admin', 'operator']), detectServerTypeHandler)

// 节点管理
app.get('/api/v1/nodes', authMiddleware, cacheMiddleware({ ttl: 30, private: true }), listNodesHandler)
app.get('/api/v1/nodes/:id', authMiddleware, getNodeHandler)
app.get('/api/v1/nodes/:id/stats', authMiddleware, getNodeStatsHandler)
app.get('/api/v1/nodes/:id/logs', authMiddleware, getNodeLogsHandler)
app.post('/api/v1/nodes', authMiddleware, rbacMiddleware(['admin', 'operator']), createNodeHandler)
app.post('/api/v1/nodes/bulk', authMiddleware, rbacMiddleware(['admin', 'operator']), bulkActionHandler)
app.post('/api/v1/nodes/bulk-status', authMiddleware, rbacMiddleware(['admin', 'operator']), bulkNodeStatusHandler)
app.post('/api/v1/nodes/:id/start', authMiddleware, rbacMiddleware(['admin', 'operator']), startNodeHandler)
app.post('/api/v1/nodes/:id/stop', authMiddleware, rbacMiddleware(['admin', 'operator']), stopNodeHandler)
app.post('/api/v1/nodes/:id/restart', authMiddleware, rbacMiddleware(['admin', 'operator']), restartNodeHandler)
app.post('/api/v1/nodes/:id/test', authMiddleware, testNodeConnectionHandler)
app.post('/api/v1/nodes/:id/sync', authMiddleware, rbacMiddleware(['admin', 'operator']), syncNodeHandler)
app.put('/api/v1/nodes/:id', authMiddleware, rbacMiddleware(['admin', 'operator']), updateNodeHandler)
app.delete('/api/v1/nodes/:id', authMiddleware, rbacMiddleware(['admin']), deleteNodeHandler)

// Playbook 管理
app.get('/api/v1/playbooks', authMiddleware, listPlaybooksHandler)
app.get('/api/v1/playbooks/built-in', authMiddleware, listBuiltInPlaybooksHandler)
app.get('/api/v1/playbooks/categories', authMiddleware, getPlaybookCategoriesHandler)
app.post('/api/v1/playbooks/sync-builtin', authMiddleware, rbacMiddleware(['admin']), syncBuiltInPlaybooksHandler)
app.get('/api/v1/playbooks/:name', authMiddleware, getPlaybookHandler)
app.post('/api/v1/playbooks', authMiddleware, rbacMiddleware(['admin', 'operator']), uploadPlaybookHandler)
app.delete('/api/v1/playbooks/:name', authMiddleware, rbacMiddleware(['admin']), deletePlaybookHandler)

// 任务管理
app.get('/api/v1/tasks', authMiddleware, listTasksHandler)
app.post('/api/v1/tasks', authMiddleware, rbacMiddleware(['admin', 'operator']), createTaskHandler)
app.get('/api/v1/tasks/:id', authMiddleware, getTaskHandler)
app.get('/api/v1/tasks/:id/logs', authMiddleware, getTaskLogsHandler)
app.post('/api/v1/tasks/:id/cancel', authMiddleware, rbacMiddleware(['admin', 'operator']), cancelTaskHandler)
app.post('/api/v1/tasks/:id/retry', authMiddleware, rbacMiddleware(['admin', 'operator']), retryTaskHandler)

// 调度管理
app.get('/api/v1/schedules', authMiddleware, listSchedulesHandler)
app.post('/api/v1/schedules', authMiddleware, rbacMiddleware(['admin', 'operator']), createScheduleHandler)
app.get('/api/v1/schedules/:id', authMiddleware, getScheduleHandler)
app.put('/api/v1/schedules/:id', authMiddleware, rbacMiddleware(['admin', 'operator']), updateScheduleHandler)
app.delete('/api/v1/schedules/:id', authMiddleware, rbacMiddleware(['admin']), deleteScheduleHandler)
app.post('/api/v1/schedules/:id/toggle', authMiddleware, rbacMiddleware(['admin', 'operator']), toggleScheduleHandler)
app.post('/api/v1/schedules/:id/run', authMiddleware, rbacMiddleware(['admin', 'operator']), runScheduleNowHandler)

// 节点组管理
app.get('/api/v1/node-groups', authMiddleware, listNodeGroupsHandler)
app.post('/api/v1/node-groups', authMiddleware, rbacMiddleware(['admin', 'operator']), createNodeGroupHandler)
app.get('/api/v1/node-groups/:id', authMiddleware, getNodeGroupHandler)
app.put('/api/v1/node-groups/:id', authMiddleware, rbacMiddleware(['admin', 'operator']), updateNodeGroupHandler)
app.delete('/api/v1/node-groups/:id', authMiddleware, rbacMiddleware(['admin']), deleteNodeGroupHandler)
app.post('/api/v1/node-groups/:id/nodes', authMiddleware, rbacMiddleware(['admin', 'operator']), addNodesToGroupHandler)
app.delete('/api/v1/node-groups/:id/nodes', authMiddleware, rbacMiddleware(['admin', 'operator']), removeNodesFromGroupHandler)

// 插件管理
app.get('/api/v1/plugins', authMiddleware, listPluginsHandler)
app.get('/api/v1/plugins/:name', authMiddleware, getPluginHandler)
app.post('/api/v1/plugins/:name/execute', authMiddleware, rbacMiddleware(['admin', 'operator']), executePluginHandler)

// Dashboard
app.get('/api/v1/dashboard', authMiddleware, dashboardHandler)
app.get('/api/v1/dashboard/stats', authMiddleware, statsHandler)

// 审计日志
app.get('/api/v1/audit-logs', authMiddleware, rbacMiddleware(['admin']), listAuditLogsHandler)

// 通知管理
app.get('/api/v1/notifications', authMiddleware, listNotificationsHandler)
app.get('/api/v1/notifications/unread-count', authMiddleware, getUnreadCountHandler)
app.post('/api/v1/notifications', authMiddleware, rbacMiddleware(['admin', 'operator']), createNotificationHandler)
app.put('/api/v1/notifications/:id/read', authMiddleware, markNotificationReadHandler)
app.put('/api/v1/notifications/read-all', authMiddleware, markAllNotificationsReadHandler)
app.delete('/api/v1/notifications/:id', authMiddleware, deleteNotificationHandler)

// ==================== SSE (实时通信) ====================
app.get('/api/v1/sse', authMiddleware, sseHandler)
app.post('/api/v1/sse/subscribe', authMiddleware, sseSubscribeHandler)
app.post('/api/v1/sse/unsubscribe', authMiddleware, sseUnsubscribeHandler)
app.get('/api/v1/sse/status', authMiddleware, sseStatusHandler)

// ==================== 备份管理 ====================
app.get('/api/v1/backups', authMiddleware, rbacMiddleware(['admin']), listBackupsHandler)
app.get('/api/v1/backups/status', authMiddleware, rbacMiddleware(['admin']), backupStatusHandler)
app.post('/api/v1/backups', authMiddleware, rbacMiddleware(['admin']), createBackupHandler)
app.get('/api/v1/backups/:id', authMiddleware, rbacMiddleware(['admin']), getBackupHandler)
app.get('/api/v1/backups/:id/download', authMiddleware, rbacMiddleware(['admin']), downloadBackupHandler)
app.post('/api/v1/backups/:id/restore', authMiddleware, rbacMiddleware(['admin']), restoreBackupHandler)
app.delete('/api/v1/backups/:id', authMiddleware, rbacMiddleware(['admin']), deleteBackupHandler)
app.post('/api/v1/backups/cleanup', authMiddleware, rbacMiddleware(['admin']), cleanupBackupsHandler)

// ==================== Agent 管理 ====================
// Agent API (认证通过Header)
app.post('/api/v1/agents/register', registerAgentHandler)
app.post('/api/v1/agents/heartbeat', agentHeartbeatHandler)
app.post('/api/v1/agents/metrics', agentMetricsHandler)
app.post('/api/v1/agents/command-result', agentCommandResultHandler)

// Agent 管理API (需要用户认证)
app.get('/api/v1/agents/:agentId/metrics', authMiddleware, getAgentMetricsHandler)
app.post('/api/v1/agents/:agentId/command', authMiddleware, rbacMiddleware(['admin', 'operator']), sendAgentCommandHandler)
app.get('/api/v1/nodes/:nodeId/install-script', authMiddleware, rbacMiddleware(['admin', 'operator']), generateInstallScriptHandler)

// ==================== MFA 双因素认证 ====================
app.get('/api/v1/mfa/status', authMiddleware, getMFAStatusHandler)
app.post('/api/v1/mfa/setup', authMiddleware, setupMFAHandler)
app.post('/api/v1/mfa/enable', authMiddleware, enableMFAHandler)
app.post('/api/v1/mfa/disable', authMiddleware, disableMFAHandler)
app.post('/api/v1/mfa/verify', authMiddleware, verifyMFAHandler)
app.post('/api/v1/mfa/recovery-codes', authMiddleware, regenerateRecoveryCodesHandler)
app.post('/api/v1/admin/users/:id/mfa/disable', authMiddleware, rbacMiddleware(['admin']), adminDisableMFAHandler)

// ==================== 批量操作 ====================
app.post('/api/v1/batch', authMiddleware, rbacMiddleware(['admin', 'operator']), batchOperationsHandler)

// ==================== WebSocket ====================
// WebSocket 暂时禁用 - Durable Object 有问题
/*
app.get('/api/v1/ws', async (c) => {
  const id = c.env.WEBSOCKET_SERVER.idFromName('global')
  const stub = c.env.WEBSOCKET_SERVER.get(id)
  return stub.fetch(c.req.raw)
})
*/

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
// 暂时禁用以调试 Cloudflare 错误 1101
/*
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
*/

// Export
export default app