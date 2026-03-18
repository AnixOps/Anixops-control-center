# AnixOps API - Cloudflare Workers

基于 Cloudflare Workers 的边缘 API 服务，使用 Hono 框架构建。

## 技术栈

- **Runtime**: Cloudflare Workers
- **Framework**: [Hono](https://hono.dev/)
- **Database**: Cloudflare D1 (SQLite)
- **Cache**: Cloudflare KV
- **Storage**: Cloudflare R2
- **Auth**: JWT + bcrypt

## 快速开始

### 1. 安装依赖

```bash
cd workers
npm install
```

### 2. 创建 D1 数据库

```bash
# 创建数据库
wrangler d1 create anixops-db

# 记录返回的 database_id，更新 wrangler.toml 中的 database_id
```

### 3. 创建 KV 命名空间

```bash
# 创建 KV 命名空间
wrangler kv:namespace create KV

# 记录返回的 id，更新 wrangler.toml 中的 id
```

### 4. 创建 R2 存储桶

```bash
# 创建 R2 存储桶
wrangler r2 bucket create anixops-files
```

### 5. 设置密钥

```bash
# 设置 JWT 密钥
wrangler secret put JWT_SECRET

# 设置 API Key 盐值
wrangler secret put API_KEY_SALT
```

### 6. 运行迁移

```bash
# 应用数据库迁移
wrangler d1 migrations apply anixops-db
```

### 7. 本地开发

```bash
# 启动开发服务器
npm run dev

# API 将在 http://localhost:8787 可用
```

### 8. 部署

```bash
# 部署到 Cloudflare
npm run deploy
```

## API 端点

### 公开端点

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/health` | 健康检查 |
| GET | `/readiness` | 就绪检查 |
| POST | `/api/v1/auth/login` | 用户登录 |
| POST | `/api/v1/auth/register` | 用户注册 |
| POST | `/api/v1/auth/refresh` | 刷新 Token |

### 认证端点 (需要 JWT)

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/v1/users/me` | 当前用户信息 |
| GET | `/api/v1/nodes` | 节点列表 |
| GET | `/api/v1/playbooks` | Playbook 列表 |
| GET | `/api/v1/plugins` | 插件列表 |
| GET | `/api/v1/dashboard` | Dashboard 概览 |

### 管理端点 (需要 admin 角色)

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/v1/users` | 用户列表 |
| POST | `/api/v1/users` | 创建用户 |
| GET | `/api/v1/audit-logs` | 审计日志 |

## 环境变量

| 变量名 | 描述 | 必需 |
|--------|------|------|
| `JWT_SECRET` | JWT 签名密钥 | ✅ |
| `JWT_EXPIRE` | Token 过期时间(秒) | ❌ (默认: 86400) |
| `API_KEY_SALT` | API Key 生成盐值 | ✅ |

## 目录结构

```
workers/
├── src/
│   ├── index.ts          # 主入口
│   ├── types.ts          # TypeScript 类型定义
│   ├── handlers/         # API 处理器
│   │   ├── auth.ts
│   │   ├── nodes.ts
│   │   ├── playbooks.ts
│   │   ├── plugins.ts
│   │   ├── dashboard.ts
│   │   ├── audit.ts
│   │   ├── users.ts
│   │   └── health.ts
│   └── middleware/       # 中间件
│       ├── auth.ts
│       └── rate-limit.ts
├── migrations/           # D1 迁移文件
│   └── 0001_initial.sql
├── package.json
├── tsconfig.json
└── wrangler.toml         # Cloudflare 配置
```

## 测试

```bash
# 运行测试
npm test

# 类型检查
npm run typecheck
```

## 监控

```bash
# 实时查看日志
npm run tail
```

## 相关文档

- [Cloudflare Workers 文档](https://developers.cloudflare.com/workers/)
- [Cloudflare D1 文档](https://developers.cloudflare.com/d1/)
- [Hono 文档](https://hono.dev/)
- [详细集成方案](../docs/cloudflare-integration-plan.md)