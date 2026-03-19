import type { Context } from 'hono'
import type { Env } from '../types'

/**
 * Get required route parameter, throws if not present
 */
export function getRequiredParam(c: Context<{ Bindings: Env }>, name: string): string {
  const value = c.req.param(name)
  if (!value) {
    throw new Error(`Missing required parameter: ${name}`)
  }
  return value
}

/**
 * Get client IP address from request
 */
export function getClientIP(c: Context<{ Bindings: Env }>): string | null {
  return c.req.header('CF-Connecting-IP') || null
}

/**
 * Get user agent from request
 */
export function getUserAgent(c: Context<{ Bindings: Env }>): string | null {
  return c.req.header('User-Agent') || null
}

/**
 * Log audit entry to database
 */
export async function logAudit(
  c: Context<{ Bindings: Env }>,
  userId: number | undefined,
  action: string,
  resource: string,
  details?: Record<string, unknown>
): Promise<void> {
  try {
    await c.env.DB
      .prepare(`
        INSERT INTO audit_logs (user_id, action, resource, ip, user_agent, details)
        VALUES (?, ?, ?, ?, ?, ?)
      `)
      .bind(
        userId ?? null,
        action,
        resource,
        getClientIP(c),
        getUserAgent(c),
        details ? JSON.stringify(details) : null
      )
      .run()
  } catch (err) {
    console.error('Failed to log audit:', err)
  }
}