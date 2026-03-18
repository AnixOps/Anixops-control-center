import type { Context } from 'hono'
import { sign } from 'jose'
import { hash, compare } from 'bcryptjs'
import { z } from 'zod'
import type { Env, User } from '../types'

// 验证 schema
const loginSchema = z.object({
  email: z.string().email(),
  password: z.string().min(6),
})

const registerSchema = z.object({
  email: z.string().email(),
  password: z.string().min(8),
  role: z.enum(['admin', 'operator', 'viewer']).optional(),
})

/**
 * 登录
 */
export async function loginHandler(c: Context<{ Bindings: Env }>) {
  try {
    const body = await c.req.json()
    const { email, password } = loginSchema.parse(body)

    // 查找用户
    const user = await c.env.DB
      .prepare('SELECT * FROM users WHERE email = ? AND enabled = 1')
      .bind(email)
      .first<User>()

    if (!user || !user.password_hash) {
      return c.json({ success: false, error: 'Invalid credentials' }, 401)
    }

    // 验证密码
    const valid = await compare(password, user.password_hash)
    if (!valid) {
      return c.json({ success: false, error: 'Invalid credentials' }, 401)
    }

    // 生成 JWT
    const secret = new TextEncoder().encode(c.env.JWT_SECRET)
    const expire = parseInt(c.env.JWT_EXPIRE, 10) || 86400

    const accessToken = await sign(
      {
        sub: user.id,
        email: user.email,
        role: user.role,
      },
      secret,
      { expiresIn: `${expire}s` }
    )

    const refreshToken = await sign(
      {
        sub: user.id,
        type: 'refresh',
      },
      secret,
      { expiresIn: '7d' }
    )

    // 更新最后登录时间
    await c.env.DB
      .prepare('UPDATE users SET last_login_at = datetime(\'now\') WHERE id = ?')
      .bind(user.id)
      .run()

    // 记录审计日志
    await logAudit(c, user.id, 'login', 'auth', { ip: c.req.header('CF-Connecting-IP') })

    return c.json({
      success: true,
      data: {
        access_token: accessToken,
        refresh_token: refreshToken,
        token_type: 'Bearer',
        expires_in: expire,
        user: {
          id: user.id,
          email: user.email,
          role: user.role,
        },
      },
    })
  } catch (err) {
    if (err instanceof z.ZodError) {
      return c.json({ success: false, error: 'Validation error', details: err.errors }, 400)
    }
    throw err
  }
}

/**
 * 注册
 */
export async function registerHandler(c: Context<{ Bindings: Env }>) {
  try {
    const body = await c.req.json()
    const { email, password, role } = registerSchema.parse(body)

    // 检查用户是否已存在
    const existing = await c.env.DB
      .prepare('SELECT id FROM users WHERE email = ?')
      .bind(email)
      .first()

    if (existing) {
      return c.json({ success: false, error: 'Email already registered' }, 409)
    }

    // 哈希密码
    const passwordHash = await hash(password, 10)

    // 创建用户
    const result = await c.env.DB
      .prepare(`
        INSERT INTO users (email, password_hash, role, auth_provider, enabled)
        VALUES (?, ?, ?, 'local', 1)
        RETURNING id, email, role, created_at
      `)
      .bind(email, passwordHash, role || 'viewer')
      .first<{ id: number; email: string; role: string; created_at: string }>()

    // 记录审计日志
    await logAudit(c, result?.id, 'register', 'user', { email })

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
 * 刷新 Token
 */
export async function refreshHandler(c: Context<{ Bindings: Env }>) {
  const body = await c.req.json<{ refresh_token?: string }>()

  if (!body.refresh_token) {
    return c.json({ success: false, error: 'Missing refresh token' }, 400)
  }

  try {
    const secret = new TextEncoder().encode(c.env.JWT_SECRET)
    const { payload } = await import('jose').then(jose =>
      jose.verify(body.refresh_token!, secret, { algorithms: ['HS256'] })
    )

    if (payload.type !== 'refresh') {
      return c.json({ success: false, error: 'Invalid token type' }, 401)
    }

    // 获取用户信息
    const user = await c.env.DB
      .prepare('SELECT id, email, role FROM users WHERE id = ? AND enabled = 1')
      .bind(payload.sub)
      .first<User>()

    if (!user) {
      return c.json({ success: false, error: 'User not found' }, 404)
    }

    // 生成新的 access token
    const expire = parseInt(c.env.JWT_EXPIRE, 10) || 86400
    const accessToken = await sign(
      {
        sub: user.id,
        email: user.email,
        role: user.role,
      },
      secret,
      { expiresIn: `${expire}s` }
    )

    return c.json({
      success: true,
      data: {
        access_token: accessToken,
        token_type: 'Bearer',
        expires_in: expire,
      },
    })
  } catch {
    return c.json({ success: false, error: 'Invalid refresh token' }, 401)
  }
}

/**
 * 登出
 */
export async function logoutHandler(c: Context<{ Bindings: Env }>) {
  // 可以在这里实现 token 黑名单 (使用 KV)
  // await c.env.KV.put(`blacklist:${token}`, '1', { expirationTtl: 86400 })

  return c.json({
    success: true,
    message: 'Logged out successfully',
  })
}

/**
 * 获取当前用户信息
 */
export async function meHandler(c: Context<{ Bindings: Env }>) {
  const user = c.get('user')

  const userInfo = await c.env.DB
    .prepare('SELECT id, email, role, auth_provider, last_login_at, created_at FROM users WHERE id = ?')
    .bind(user.sub)
    .first<User>()

  if (!userInfo) {
    return c.json({ success: false, error: 'User not found' }, 404)
  }

  return c.json({
    success: true,
    data: userInfo,
  })
}

/**
 * 记录审计日志
 */
async function logAudit(
  c: Context<{ Bindings: Env }>,
  userId: number | undefined,
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
        userId || null,
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