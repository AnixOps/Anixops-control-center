import type { Context } from 'hono'
import type { Env, Playbook } from '../types'

/**
 * 获取 Playbook 列表
 */
export async function listPlaybooksHandler(c: Context<{ Bindings: Env }>) {
  const result = await c.env.DB
    .prepare('SELECT id, name, description, created_at, updated_at FROM playbooks ORDER BY name')
    .all<Omit<Playbook, 'storage_key'>>()

  return c.json({
    success: true,
    data: result.results,
  })
}

/**
 * 获取单个 Playbook
 */
export async function getPlaybookHandler(c: Context<{ Bindings: Env }>) {
  const name = c.req.param('name')

  const meta = await c.env.DB
    .prepare('SELECT * FROM playbooks WHERE name = ?')
    .bind(name)
    .first<Playbook>()

  if (!meta) {
    return c.json({ success: false, error: 'Playbook not found' }, 404)
  }

  // 从 R2 获取内容
  const object = await c.env.R2.get(meta.storage_key)
  if (!object) {
    return c.json({ success: false, error: 'Playbook content not found' }, 404)
  }

  const content = await object.text()

  return c.json({
    success: true,
    data: {
      ...meta,
      content,
    },
  })
}

/**
 * 上传 Playbook
 */
export async function uploadPlaybookHandler(c: Context<{ Bindings: Env }>) {
  const body = await c.req.json<{ name: string; content: string; description?: string }>()
  const user = c.get('user')

  if (!body.name || !body.content) {
    return c.json({ success: false, error: 'Missing name or content' }, 400)
  }

  const storageKey = `playbooks/${body.name}.yml`

  // 存储到 R2
  await c.env.R2.put(storageKey, body.content, {
    httpMetadata: {
      contentType: 'text/yaml',
    },
    customMetadata: {
      uploaded_by: String(user.sub),
    },
  })

  // 更新元数据
  await c.env.DB
    .prepare(`
      INSERT INTO playbooks (name, storage_key, description, created_at, updated_at)
      VALUES (?, ?, ?, datetime('now'), datetime('now'))
      ON CONFLICT(name) DO UPDATE SET
        storage_key = excluded.storage_key,
        description = COALESCE(excluded.description, description),
        updated_at = datetime('now')
    `)
    .bind(body.name, storageKey, body.description || null)
    .run()

  await logAudit(c, user.sub, 'upload_playbook', 'playbook', { name: body.name })

  return c.json({
    success: true,
    data: { name: body.name, storage_key: storageKey },
  }, 201)
}

/**
 * 运行 Playbook
 */
export async function runPlaybookHandler(c: Context<{ Bindings: Env }>) {
  const name = c.req.param('name')
  const body = await c.req.json<{ targets?: string[]; extra_vars?: Record<string, unknown> }>()
  const user = c.get('user')

  // 验证 playbook 存在
  const meta = await c.env.DB
    .prepare('SELECT id FROM playbooks WHERE name = ?')
    .bind(name)
    .first()

  if (!meta) {
    return c.json({ success: false, error: 'Playbook not found' }, 404)
  }

  // 这里应该触发实际的 playbook 执行
  // 可以通过 Durable Objects 或发送到消息队列

  await logAudit(c, user.sub, 'run_playbook', 'playbook', {
    name,
    targets: body.targets,
  })

  return c.json({
    success: true,
    data: {
      job_id: crypto.randomUUID(),
      status: 'pending',
      message: 'Playbook execution started',
    },
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