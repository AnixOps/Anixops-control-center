import type { Context } from 'hono'
import { z } from 'zod'
import type { Env, Node } from '../types'

const createNodeSchema = z.object({
  name: z.string().min(1).max(100),
  host: z.string().min(1),
  port: z.number().int().min(1).max(65535).default(22),
  config: z.record(z.unknown()).optional(),
})

const updateNodeSchema = z.object({
  name: z.string().min(1).max(100).optional(),
  host: z.string().min(1).optional(),
  port: z.number().int().min(1).max(65535).optional(),
  status: z.enum(['online', 'offline', 'maintenance']).optional(),
  config: z.record(z.unknown()).optional(),
})

/**
 * 获取节点列表
 */
export async function listNodesHandler(c: Context<{ Bindings: Env }>) {
  const page = parseInt(c.req.query('page') || '1', 10)
  const perPage = parseInt(c.req.query('per_page') || '20', 10)
  const status = c.req.query('status')

  let query = 'SELECT * FROM nodes'
  const params: (string | number)[] = []

  if (status) {
    query += ' WHERE status = ?'
    params.push(status)
  }

  // 获取总数
  const countResult = await c.env.DB
    .prepare(`SELECT COUNT(*) as total FROM (${query})`)
    .bind(...params)
    .first<{ total: number }>()

  // 获取分页数据
  query += ' ORDER BY created_at DESC LIMIT ? OFFSET ?'
  params.push(perPage, (page - 1) * perPage)

  const result = await c.env.DB
    .prepare(query)
    .bind(...params)
    .all<Node>()

  return c.json({
    success: true,
    data: {
      items: result.results,
      total: countResult?.total || 0,
      page,
      per_page: perPage,
      total_pages: Math.ceil((countResult?.total || 0) / perPage),
    },
  })
}

/**
 * 获取单个节点
 */
export async function getNodeHandler(c: Context<{ Bindings: Env }>) {
  const id = c.req.param('id')

  const node = await c.env.DB
    .prepare('SELECT * FROM nodes WHERE id = ?')
    .bind(id)
    .first<Node>()

  if (!node) {
    return c.json({ success: false, error: 'Node not found' }, 404)
  }

  return c.json({
    success: true,
    data: node,
  })
}

/**
 * 创建节点
 */
export async function createNodeHandler(c: Context<{ Bindings: Env }>) {
  try {
    const body = await c.req.json()
    const data = createNodeSchema.parse(body)
    const user = c.get('user')

    // 检查名称是否已存在
    const existing = await c.env.DB
      .prepare('SELECT id FROM nodes WHERE name = ?')
      .bind(data.name)
      .first()

    if (existing) {
      return c.json({ success: false, error: 'Node name already exists' }, 409)
    }

    const result = await c.env.DB
      .prepare(`
        INSERT INTO nodes (name, host, port, status, config)
        VALUES (?, ?, ?, 'offline', ?)
        RETURNING *
      `)
      .bind(
        data.name,
        data.host,
        data.port,
        data.config ? JSON.stringify(data.config) : null
      )
      .first<Node>()

    // 记录审计日志
    await logAudit(c, user.sub, 'create_node', 'node', { node_id: result?.id, name: data.name })

    return c.json({
      success: true,
      data: result,
    }, 201)
  } catch (err) {
    if (err instanceof z.ZodError) {
      return c.json({ success: false, error: 'Validation error', details: err.errors }, 400)
    }
    throw err
  }
}

/**
 * 更新节点
 */
export async function updateNodeHandler(c: Context<{ Bindings: Env }>) {
  const id = c.req.param('id')
  const user = c.get('user')

  try {
    const body = await c.req.json()
    const data = updateNodeSchema.parse(body)

    // 检查节点是否存在
    const existing = await c.env.DB
      .prepare('SELECT id FROM nodes WHERE id = ?')
      .bind(id)
      .first()

    if (!existing) {
      return c.json({ success: false, error: 'Node not found' }, 404)
    }

    // 构建更新语句
    const updates: string[] = []
    const values: (string | number | null)[] = []

    if (data.name) {
      updates.push('name = ?')
      values.push(data.name)
    }
    if (data.host) {
      updates.push('host = ?')
      values.push(data.host)
    }
    if (data.port) {
      updates.push('port = ?')
      values.push(data.port)
    }
    if (data.status) {
      updates.push('status = ?')
      values.push(data.status)
    }
    if (data.config !== undefined) {
      updates.push('config = ?')
      values.push(JSON.stringify(data.config))
    }

    if (updates.length === 0) {
      return c.json({ success: false, error: 'No fields to update' }, 400)
    }

    updates.push('updated_at = datetime(\'now\')')
    values.push(id)

    const result = await c.env.DB
      .prepare(`UPDATE nodes SET ${updates.join(', ')} WHERE id = ? RETURNING *`)
      .bind(...values)
      .first<Node>()

    await logAudit(c, user.sub, 'update_node', 'node', { node_id: id })

    return c.json({
      success: true,
      data: result,
    })
  } catch (err) {
    if (err instanceof z.ZodError) {
      return c.json({ success: false, error: 'Validation error', details: err.errors }, 400)
    }
    throw err
  }
}

/**
 * 删除节点
 */
export async function deleteNodeHandler(c: Context<{ Bindings: Env }>) {
  const id = c.req.param('id')
  const user = c.get('user')

  const result = await c.env.DB
    .prepare('DELETE FROM nodes WHERE id = ? RETURNING id')
    .bind(id)
    .first()

  if (!result) {
    return c.json({ success: false, error: 'Node not found' }, 404)
  }

  await logAudit(c, user.sub, 'delete_node', 'node', { node_id: id })

  return c.json({
    success: true,
    message: 'Node deleted successfully',
  })
}

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