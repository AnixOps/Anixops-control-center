import type { Context } from 'hono'
import { z } from 'zod'
import type { Env } from '../types'

const createNotificationSchema = z.object({
  title: z.string().min(1).max(200),
  message: z.string().min(1),
  type: z.enum(['info', 'warning', 'error', 'success']).default('info'),
  user_id: z.number().int().optional(), // 如果指定则发送给特定用户，否则广播
  data: z.record(z.unknown()).optional(),
})

/**
 * 获取通知列表
 */
export async function listNotificationsHandler(c: Context<{ Bindings: Env }>) {
  const user = c.get('user')
  const page = parseInt(c.req.query('page') || '1', 10)
  const perPage = parseInt(c.req.query('per_page') || '20', 10)
  const unreadOnly = c.req.query('unread_only') === 'true'

  // 从KV获取用户通知
  const notificationsKey = `notifications:${user.sub}`
  let notifications = await c.env.KV.get(notificationsKey, 'json') as Array<{
    id: string
    title: string
    message: string
    type: string
    read: boolean
    created_at: string
    data?: Record<string, unknown>
  }> | null

  if (!notifications) {
    // 生成一些模拟通知
    notifications = generateMockNotifications()
    await c.env.KV.put(notificationsKey, JSON.stringify(notifications), { expirationTtl: 86400 * 7 })
  }

  // 过滤未读
  if (unreadOnly) {
    notifications = notifications.filter(n => !n.read)
  }

  // 分页
  const total = notifications.length
  const startIndex = (page - 1) * perPage
  const paginatedNotifications = notifications.slice(startIndex, startIndex + perPage)

  return c.json({
    success: true,
    data: {
      items: paginatedNotifications,
      total,
      page,
      per_page: perPage,
      total_pages: Math.ceil(total / perPage),
      unread_count: notifications.filter(n => !n.read).length,
    },
  })
}

/**
 * 标记通知为已读
 */
export async function markNotificationReadHandler(c: Context<{ Bindings: Env }>) {
  const user = c.get('user')
  const notificationId = c.req.param('id')

  const notificationsKey = `notifications:${user.sub}`
  const notifications = await c.env.KV.get(notificationsKey, 'json') as Array<{
    id: string
    title: string
    message: string
    type: string
    read: boolean
    created_at: string
  }> | null

  if (notifications) {
    const index = notifications.findIndex(n => n.id === notificationId)
    if (index !== -1) {
      notifications[index].read = true
      await c.env.KV.put(notificationsKey, JSON.stringify(notifications), { expirationTtl: 86400 * 7 })
    }
  }

  return c.json({
    success: true,
    message: 'Notification marked as read',
  })
}

/**
 * 标记所有通知为已读
 */
export async function markAllNotificationsReadHandler(c: Context<{ Bindings: Env }>) {
  const user = c.get('user')

  const notificationsKey = `notifications:${user.sub}`
  const notifications = await c.env.KV.get(notificationsKey, 'json') as Array<{
    id: string
    read: boolean
  }> | null

  if (notifications) {
    for (const notification of notifications) {
      notification.read = true
    }
    await c.env.KV.put(notificationsKey, JSON.stringify(notifications), { expirationTtl: 86400 * 7 })
  }

  return c.json({
    success: true,
    message: 'All notifications marked as read',
  })
}

/**
 * 删除通知
 */
export async function deleteNotificationHandler(c: Context<{ Bindings: Env }>) {
  const user = c.get('user')
  const notificationId = c.req.param('id')

  const notificationsKey = `notifications:${user.sub}`
  const notifications = await c.env.KV.get(notificationsKey, 'json') as Array<{
    id: string
  }> | null

  if (notifications) {
    const filtered = notifications.filter(n => n.id !== notificationId)
    await c.env.KV.put(notificationsKey, JSON.stringify(filtered), { expirationTtl: 86400 * 7 })
  }

  return c.json({
    success: true,
    message: 'Notification deleted',
  })
}

/**
 * 创建通知（内部API）
 */
export async function createNotificationHandler(c: Context<{ Bindings: Env }>) {
  const user = c.get('user')

  try {
    const body = await c.req.json()
    const data = createNotificationSchema.parse(body)

    const notification = {
      id: crypto.randomUUID(),
      title: data.title,
      message: data.message,
      type: data.type,
      read: false,
      created_at: new Date().toISOString(),
      data: data.data,
    }

    if (data.user_id) {
      // 发送给特定用户
      const notificationsKey = `notifications:${data.user_id}`
      const notifications = await c.env.KV.get(notificationsKey, 'json') as Array<typeof notification> | null
      await c.env.KV.put(notificationsKey, JSON.stringify([notification, ...(notifications || [])]), { expirationTtl: 86400 * 7 })
    } else {
      // 广播给所有用户 - 这里简化处理，只记录到系统通知
      const broadcastKey = 'notifications:broadcast'
      const broadcasts = await c.env.KV.get(broadcastKey, 'json') as Array<typeof notification> | null
      await c.env.KV.put(broadcastKey, JSON.stringify([notification, ...(broadcasts || [])]), { expirationTtl: 86400 * 7 })
    }

    await logAudit(c, user.sub, 'create_notification', 'notification', { title: data.title })

    return c.json({
      success: true,
      data: notification,
    }, 201)
  } catch (err) {
    if (err instanceof z.ZodError) {
      return c.json({ success: false, error: 'Validation error', details: err.errors }, 400)
    }
    throw err
  }
}

/**
 * 获取未读通知数量
 */
export async function getUnreadCountHandler(c: Context<{ Bindings: Env }>) {
  const user = c.get('user')

  const notificationsKey = `notifications:${user.sub}`
  const notifications = await c.env.KV.get(notificationsKey, 'json') as Array<{
    read: boolean
  }> | null

  const unreadCount = notifications ? notifications.filter(n => !n.read).length : 0

  return c.json({
    success: true,
    data: {
      unread_count: unreadCount,
    },
  })
}

// 辅助函数：生成模拟通知
function generateMockNotifications(): Array<{
  id: string
  title: string
  message: string
  type: string
  read: boolean
  created_at: string
  data?: Record<string, unknown>
}> {
  const types = ['info', 'warning', 'error', 'success']
  const templates = [
    { title: 'Node Offline', message: 'Node "US-East-1" has gone offline', type: 'error' },
    { title: 'High CPU Usage', message: 'Node "EU-West-2" CPU usage exceeded 90%', type: 'warning' },
    { title: 'New User Registered', message: 'A new user has registered on the platform', type: 'info' },
    { title: 'Backup Completed', message: 'Daily backup completed successfully', type: 'success' },
    { title: 'SSL Certificate Expiring', message: 'SSL certificate for api.anixops.com expires in 7 days', type: 'warning' },
    { title: 'Playbook Executed', message: 'Ansible playbook "deploy-app" completed successfully', type: 'success' },
  ]

  return templates.map((template, index) => ({
    id: `notification-${index + 1}`,
    title: template.title,
    message: template.message,
    type: template.type,
    read: index > 2, // 前3条未读
    created_at: new Date(Date.now() - index * 3600000).toISOString(),
  }))
}

// 辅助函数：记录审计日志
async function logAudit(
  c: Context<{ Bindings: Env }>,
  userId: number,
  action: string,
  resource: string,
  details?: Record<string, unknown>
) {
  try {
    await c.env.DB
      .prepare(`
        INSERT INTO audit_logs (user_id, action, resource, ip, user_agent, details)
        VALUES (?, ?, ?, ?, ?, ?)
      `)
      .bind(
        userId,
        action,
        resource,
        c.req.header('CF-Connecting-IP') || null,
        c.req.header('User-Agent') || null,
        details ? JSON.stringify(details) : null
      )
      .run()
  } catch (err) {
    console.error('Failed to log audit:', err)
  }
}