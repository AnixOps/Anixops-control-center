import type { Context, Next } from 'hono'
import type { Env } from '../types'

interface RateLimitOptions {
  windowMs: number  // 时间窗口 (毫秒)
  max: number       // 最大请求数
  keyGenerator?: (c: Context<{ Bindings: Env }>) => string
}

/**
 * 速率限制中间件
 * 使用 KV 存储实现
 */
export function rateLimitMiddleware(options: RateLimitOptions) {
  return async (c: Context<{ Bindings: Env }>, next: Next) => {
    const key = options.keyGenerator
      ? options.keyGenerator(c)
      : `ratelimit:${c.req.header('CF-Connecting-IP') || 'unknown'}:${c.req.path}`

    try {
      // 获取当前计数
      const stored = await c.env.KV.get(key)
      const count = stored ? parseInt(stored, 10) : 0

      if (count >= options.max) {
        return c.json({
          success: false,
          error: 'Too Many Requests',
          retry_after: Math.ceil(options.windowMs / 1000),
        }, 429, {
          'Retry-After': String(Math.ceil(options.windowMs / 1000)),
          'X-RateLimit-Limit': String(options.max),
          'X-RateLimit-Remaining': '0',
          'X-RateLimit-Reset': String(Math.ceil(Date.now() / 1000 + options.windowMs / 1000)),
        })
      }

      // 增加计数
      await c.env.KV.put(key, String(count + 1), {
        expirationTtl: Math.ceil(options.windowMs / 1000),
      })

      // 设置响应头
      c.header('X-RateLimit-Limit', String(options.max))
      c.header('X-RateLimit-Remaining', String(options.max - count - 1))

      await next()
    } catch (err) {
      // KV 错误不应该阻止请求
      console.error('Rate limit error:', err)
      await next()
    }
  }
}

/**
 * IP 白名单中间件
 */
export function ipWhitelistMiddleware(allowedIPs: string[]) {
  return async (c: Context<{ Bindings: Env }>, next: Next) => {
    const clientIP = c.req.header('CF-Connecting-IP')

    if (!clientIP || !allowedIPs.includes(clientIP)) {
      return c.json({ success: false, error: 'Forbidden: IP not allowed' }, 403)
    }

    await next()
  }
}