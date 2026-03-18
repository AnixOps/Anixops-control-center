import type { Context, Next } from 'hono'
import { jwtVerify } from 'jose'
import type { Env, JWTPayload } from '../types'

/**
 * JWT 认证中间件
 */
export async function authMiddleware(c: Context<{ Bindings: Env }>, next: Next) {
  const authHeader = c.req.header('Authorization')

  if (!authHeader || !authHeader.startsWith('Bearer ')) {
    return c.json({ success: false, error: 'Unauthorized: Missing or invalid Authorization header' }, 401)
  }

  const token = authHeader.substring(7)

  try {
    const secret = new TextEncoder().encode(c.env.JWT_SECRET)
    const { payload } = await jwtVerify(token, secret)

    c.set('user', payload as JWTPayload)
    await next()
  } catch (err) {
    if (err instanceof Error) {
      if (err.message.includes('expired')) {
        return c.json({ success: false, error: 'Token expired' }, 401)
      }
    }
    return c.json({ success: false, error: 'Invalid token' }, 401)
  }
}

/**
 * 基于角色的访问控制中间件
 */
export function rbacMiddleware(allowedRoles: string[]) {
  return async (c: Context<{ Bindings: Env }>, next: Next) => {
    const user = c.get('user')

    if (!user) {
      return c.json({ success: false, error: 'Unauthorized' }, 401)
    }

    if (!allowedRoles.includes(user.role)) {
      return c.json({
        success: false,
        error: 'Forbidden: Insufficient permissions',
        required_roles: allowedRoles,
        your_role: user.role,
      }, 403)
    }

    await next()
  }
}

/**
 * API Key 认证中间件 (用于代理/CLI)
 */
export async function apiKeyMiddleware(c: Context<{ Bindings: Env }>, next: Next) {
  const apiKey = c.req.header('X-API-Key')

  if (!apiKey) {
    return c.json({ success: false, error: 'Missing API Key' }, 401)
  }

  // 验证 API Key
  const result = await c.env.DB
    .prepare(`
      SELECT u.id, u.email, u.role
      FROM api_tokens t
      JOIN users u ON t.user_id = u.id
      WHERE t.token = ? AND (t.expires_at IS NULL OR t.expires_at > datetime('now'))
    `)
    .bind(apiKey)
    .first<{ id: number; email: string; role: string }>()

  if (!result) {
    return c.json({ success: false, error: 'Invalid or expired API Key' }, 401)
  }

  // 更新最后使用时间
  await c.env.DB
    .prepare('UPDATE api_tokens SET last_used = datetime(\'now\') WHERE token = ?')
    .bind(apiKey)
    .run()

  c.set('user', {
    sub: result.id,
    email: result.email,
    role: result.role,
    iat: Date.now() / 1000,
    exp: 0,
  })

  await next()
}