import type { Context } from 'hono'
import { z } from 'zod'
import { hash, compare } from 'bcryptjs'
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

const changePasswordSchema = z.object({
  current_password: z.string().min(1),
  new_password: z.string().min(8),
})

const createTokenSchema = z.object({
  name: z.string().min(1).max(100),
  expires_in_days: z.number().int().min(1).max(365).optional(),
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

/**
 * 修改密码
 */
export async function changePasswordHandler(c: Context<{ Bindings: Env }>) {
  const currentUser = c.get('user')

  try {
    const body = await c.req.json()
    const data = changePasswordSchema.parse(body)

    // 获取用户当前密码
    const user = await c.env.DB
      .prepare('SELECT id, password_hash FROM users WHERE id = ?')
      .bind(currentUser.sub)
      .first<{ id: number; password_hash: string }>()

    if (!user) {
      return c.json({ success: false, error: 'User not found' }, 404)
    }

    // 验证当前密码
    const isValid = await compare(data.current_password, user.password_hash)
    if (!isValid) {
      return c.json({ success: false, error: 'Current password is incorrect' }, 400)
    }

    // 更新密码
    const newPasswordHash = await hash(data.new_password, 10)
    await c.env.DB
      .prepare('UPDATE users SET password_hash = ?, updated_at = datetime(\'now\') WHERE id = ?')
      .bind(newPasswordHash, currentUser.sub)
      .run()

    await logAudit(c, currentUser.sub, 'change_password', 'user', { user_id: currentUser.sub })

    return c.json({
      success: true,
      message: 'Password changed successfully',
    })
  } catch (err) {
    if (err instanceof z.ZodError) {
      return c.json({ success: false, error: 'Validation error', details: err.errors }, 400)
    }
    throw err
  }
}

/**
 * 获取当前用户信息
 */
export async function getCurrentUserHandler(c: Context<{ Bindings: Env }>) {
  const currentUser = c.get('user')

  const user = await c.env.DB
    .prepare('SELECT id, email, role, auth_provider, enabled, last_login_at, created_at FROM users WHERE id = ?')
    .bind(currentUser.sub)
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
 * 更新当前用户信息
 */
export async function updateCurrentUserHandler(c: Context<{ Bindings: Env }>) {
  const currentUser = c.get('user')

  try {
    const body = await c.req.json()
    const { name, email } = body as { name?: string; email?: string }

    // 更新用户信息
    if (email) {
      const existing = await c.env.DB
        .prepare('SELECT id FROM users WHERE email = ? AND id != ?')
        .bind(email, currentUser.sub)
        .first()

      if (existing) {
        return c.json({ success: false, error: 'Email already in use' }, 409)
      }

      await c.env.DB
        .prepare('UPDATE users SET email = ?, updated_at = datetime(\'now\') WHERE id = ?')
        .bind(email, currentUser.sub)
        .run()
    }

    const user = await c.env.DB
      .prepare('SELECT id, email, role, auth_provider, enabled, last_login_at, created_at FROM users WHERE id = ?')
      .bind(currentUser.sub)
      .first<Omit<User, 'password_hash'>>()

    await logAudit(c, currentUser.sub, 'update_profile', 'user', { user_id: currentUser.sub })

    return c.json({
      success: true,
      data: user,
    })
  } catch (err) {
    throw err
  }
}

// ==================== API Tokens ====================

/**
 * 生成随机token
 */
function generateToken(): string {
  const chars = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'
  let token = 'sk_live_'
  for (let i = 0; i < 32; i++) {
    token += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  return token
}

/**
 * 列出API tokens
 */
export async function listApiTokensHandler(c: Context<{ Bindings: Env }>) {
  const currentUser = c.get('user')

  const tokens = await c.env.DB
    .prepare(`
      SELECT id, name, created_at, last_used, expires_at
      FROM api_tokens
      WHERE user_id = ?
      ORDER BY created_at DESC
    `)
    .bind(currentUser.sub)
    .all<{ id: number; name: string; created_at: string; last_used: string | null; expires_at: string | null }>()

  return c.json({
    success: true,
    data: tokens.results,
  })
}

/**
 * 创建API token
 */
export async function createApiTokenHandler(c: Context<{ Bindings: Env }>) {
  const currentUser = c.get('user')

  try {
    const body = await c.req.json()
    const data = createTokenSchema.parse(body)

    const token = generateToken()
    const tokenHash = await hash(token, 10)

    let expiresAt: string | null = null
    if (data.expires_in_days) {
      const expires = new Date()
      expires.setDate(expires.getDate() + data.expires_in_days)
      expiresAt = expires.toISOString()
    }

    const result = await c.env.DB
      .prepare(`
        INSERT INTO api_tokens (user_id, name, token, created_at, expires_at)
        VALUES (?, ?, ?, datetime('now'), ?)
        RETURNING id, name, created_at, expires_at
      `)
      .bind(currentUser.sub, data.name, tokenHash, expiresAt)
      .first<{ id: number; name: string; created_at: string; expires_at: string | null }>()

    await logAudit(c, currentUser.sub, 'create_api_token', 'api_token', { token_id: result?.id, name: data.name })

    // 返回token（仅此一次可见）
    return c.json({
      success: true,
      data: {
        ...result,
        token, // 明文token，只返回一次
      },
    }, 201)
  } catch (err) {
    if (err instanceof z.ZodError) {
      return c.json({ success: false, error: 'Validation error', details: err.errors }, 400)
    }
    throw err
  }
}

/**
 * 删除API token
 */
export async function deleteApiTokenHandler(c: Context<{ Bindings: Env }>) {
  const currentUser = c.get('user')
  const tokenId = c.req.param('id')

  const result = await c.env.DB
    .prepare('DELETE FROM api_tokens WHERE id = ? AND user_id = ? RETURNING id, name')
    .bind(tokenId, currentUser.sub)
    .first<{ id: number; name: string }>()

  if (!result) {
    return c.json({ success: false, error: 'Token not found' }, 404)
  }

  await logAudit(c, currentUser.sub, 'delete_api_token', 'api_token', { token_id: tokenId, name: result.name })

  return c.json({
    success: true,
    message: 'API token deleted successfully',
  })
}

// ==================== Sessions ====================

/**
 * 获取活跃会话
 */
export async function listSessionsHandler(c: Context<{ Bindings: Env }>) {
  const currentUser = c.get('user')

  // 从KV获取活跃会话
  const sessionsKey = `sessions:${currentUser.sub}`
  const sessions = await c.env.KV.get(sessionsKey, 'json') as Array<{
    id: string;
    device: string;
    ip: string;
    last_active: string;
    created_at: string;
  }> | null

  // 获取当前会话ID（从请求头或生成）
  const currentSessionId = c.req.header('X-Session-Id') || 'current'

  return c.json({
    success: true,
    data: {
      sessions: sessions || [],
      current_session_id: currentSessionId,
    },
  })
}

/**
 * 删除其他会话
 */
export async function deleteOtherSessionsHandler(c: Context<{ Bindings: Env }>) {
  const currentUser = c.get('user')
  const sessionToKeep = c.req.header('X-Session-Id') || 'current'

  // 清除其他会话
  await c.env.KV.put(`sessions:${currentUser.sub}`, JSON.stringify([{
    id: sessionToKeep,
    device: 'Current Device',
    ip: c.req.header('CF-Connecting-IP') || 'Unknown',
    last_active: new Date().toISOString(),
    created_at: new Date().toISOString(),
  }]))

  await logAudit(c, currentUser.sub, 'delete_sessions', 'session', { kept_session: sessionToKeep })

  return c.json({
    success: true,
    message: 'All other sessions have been signed out',
  })
}