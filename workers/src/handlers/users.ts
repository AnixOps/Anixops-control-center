import type { Context } from 'hono'
import { z } from 'zod'
import { hash } from 'bcryptjs'
import type { Env, User } from '../types'

const createUserSchema = z.object({
  email: z.string().email(),
  password: z.string().min(8),
  role: z.enum(['admin', 'operator', 'viewer']).default('viewer'),
})

const updateUserSchema = z.object({
  email: z.string().email().optional(),
  password: z.string().min(8).optional(),
  role: z.enum(['admin', 'operator', 'viewer']).optional(),
  enabled: z.boolean().optional(),
})

/**
 * 获取用户列表
 */
export async function listUsersHandler(c: Context<{ Bindings: Env }>) {
  const page = parseInt(c.req.query('page') || '1', 10)
  const perPage = parseInt(c.req.query('per_page') || '50', 10)
  const search = c.req.query('search')

  let query = 'SELECT id, email, role, auth_provider, enabled, last_login_at, created_at FROM users WHERE 1=1'
  const params: (string | number)[] = []

  if (search) {
    query += ' AND email LIKE ?'
    params.push(`%${search}%`)
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
    .all<Omit<User, 'password_hash'>>()

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
 * 获取单个用户
 */
export async function getUserHandler(c: Context<{ Bindings: Env }>) {
  const id = c.req.param('id')

  const user = await c.env.DB
    .prepare('SELECT id, email, role, auth_provider, enabled, last_login_at, created_at FROM users WHERE id = ?')
    .bind(id)
    .first<Omit<User, 'password_hash'>>()

  if (!user) {
    return c.json({ success: false, error: 'User not found' }, 404)
  }

  return c.json({
    success: true,
    data: user,
  })
}

/**
 * 创建用户
 */
export async function createUserHandler(c: Context<{ Bindings: Env }>) {
  const currentUser = c.get('user')

  try {
    const body = await c.req.json()
    const data = createUserSchema.parse(body)

    // 检查邮箱是否已存在
    const existing = await c.env.DB
      .prepare('SELECT id FROM users WHERE email = ?')
      .bind(data.email)
      .first()

    if (existing) {
      return c.json({ success: false, error: 'Email already exists' }, 409)
    }

    // 哈希密码
    const passwordHash = await hash(data.password, 10)

    const result = await c.env.DB
      .prepare(`
        INSERT INTO users (email, password_hash, role, auth_provider, enabled)
        VALUES (?, ?, ?, 'local', 1)
        RETURNING id, email, role, created_at
      `)
      .bind(data.email, passwordHash, data.role)
      .first<{ id: number; email: string; role: string; created_at: string }>()

    await logAudit(c, currentUser.sub, 'create_user', 'user', { user_id: result?.id, email: data.email })

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
 * 更新用户
 */
export async function updateUserHandler(c: Context<{ Bindings: Env }>) {
  const id = c.req.param('id')
  const currentUser = c.get('user')

  try {
    const body = await c.req.json()
    const data = updateUserSchema.parse(body)

    // 检查用户是否存在
    const existing = await c.env.DB
      .prepare('SELECT id FROM users WHERE id = ?')
      .bind(id)
      .first()

    if (!existing) {
      return c.json({ success: false, error: 'User not found' }, 404)
    }

    // 构建更新语句
    const updates: string[] = []
    const values: (string | number | null)[] = []

    if (data.email) {
      updates.push('email = ?')
      values.push(data.email)
    }
    if (data.password) {
      updates.push('password_hash = ?')
      values.push(await hash(data.password, 10))
    }
    if (data.role) {
      updates.push('role = ?')
      values.push(data.role)
    }
    if (data.enabled !== undefined) {
      updates.push('enabled = ?')
      values.push(data.enabled ? 1 : 0)
    }

    if (updates.length === 0) {
      return c.json({ success: false, error: 'No fields to update' }, 400)
    }

    updates.push('updated_at = datetime(\'now\')')
    values.push(id)

    const result = await c.env.DB
      .prepare(`UPDATE users SET ${updates.join(', ')} WHERE id = ? RETURNING id, email, role, enabled, updated_at`)
      .bind(...values)
      .first()

    await logAudit(c, currentUser.sub, 'update_user', 'user', { user_id: id })

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
 * 删除用户
 */
export async function deleteUserHandler(c: Context<{ Bindings: Env }>) {
  const id = c.req.param('id')
  const currentUser = c.get('user')

  // 不能删除自己
  if (parseInt(id, 10) === currentUser.sub) {
    return c.json({ success: false, error: 'Cannot delete yourself' }, 400)
  }

  const result = await c.env.DB
    .prepare('DELETE FROM users WHERE id = ? RETURNING id')
    .bind(id)
    .first()

  if (!result) {
    return c.json({ success: false, error: 'User not found' }, 404)
  }

  await logAudit(c, currentUser.sub, 'delete_user', 'user', { user_id: id })

  return c.json({
    success: true,
    message: 'User deleted successfully',
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