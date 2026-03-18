import type { Context } from 'hono'
import type { Env, AuditLog } from '../types'

/**
 * 获取审计日志列表
 */
export async function listAuditLogsHandler(c: Context<{ Bindings: Env }>) {
  const page = parseInt(c.req.query('page') || '1', 10)
  const perPage = parseInt(c.req.query('per_page') || '50', 10)
  const userId = c.req.query('user_id')
  const action = c.req.query('action')
  const startDate = c.req.query('start_date')
  const endDate = c.req.query('end_date')

  let query = 'SELECT * FROM audit_logs WHERE 1=1'
  const params: (string | number)[] = []

  if (userId) {
    query += ' AND user_id = ?'
    params.push(userId)
  }

  if (action) {
    query += ' AND action LIKE ?'
    params.push(`%${action}%`)
  }

  if (startDate) {
    query += ' AND created_at >= ?'
    params.push(startDate)
  }

  if (endDate) {
    query += ' AND created_at <= ?'
    params.push(endDate)
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
    .all<AuditLog>()

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