# Cloudflare 资源 ID 获取指南

## 前提条件

1. 拥有 Cloudflare 账号 (免费注册: https://dash.cloudflare.com/sign-up)
2. 安装 Wrangler CLI: `npm install -g wrangler`
3. 登录: `wrangler login`

---

## 方法一：通过 Wrangler CLI (推荐)

### 1. D1 数据库

```bash
# 创建数据库
wrangler d1 create anixops-db

# 输出示例:
# ✅ Successfully created DB 'anixops-db'!
# database_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
#              ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
#              这就是你的 database_id
```

### 2. KV 命名空间

```bash
# 创建 KV 命名空间
wrangler kv:namespace create KV

# 输出示例:
# ✅ Created namespace: KV
# ID = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
#      ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
#      这就是你的 KV id
```

### 3. R2 存储桶

```bash
# 创建 R2 存储桶
wrangler r2 bucket create anixops-files

# 输出示例:
# ✅ Created bucket 'anixops-files'
# R2 不需要 ID，直接使用 bucket_name
```

---

## 方法二：通过 Cloudflare Dashboard

### 1. D1 数据库 ID

1. 打开: https://dash.cloudflare.com/
2. 左侧菜单 → **Workers & Pages** → **D1**
3. 点击 **Create database**
4. 输入名称: `anixops-db`
5. 创建后，点击数据库名称
6. 在 **Settings** 标签页找到 **Database ID**

```
Database ID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

### 2. KV 命名空间 ID

1. 打开: https://dash.cloudflare.com/
2. 左侧菜单 → **Workers & Pages** → **KV**
3. 点击 **Create a namespace**
4. 输入名称: `anixops-kv` (或直接用 `KV`)
5. 创建后，列表中显示的 **ID** 列就是你的 KV ID

```
ID: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

### 3. R2 存储桶

1. 打开: https://dash.cloudflare.com/
2. 左侧菜单 → **R2**
3. 点击 **Create bucket**
4. 输入名称: `anixops-files`
5. 创建后，bucket 名称就是标识符

---

## 更新 wrangler.toml

获取到 ID 后，更新 `wrangler.toml`:

```toml
# D1 数据库绑定
[[d1_databases]]
binding = "DB"
database_name = "anixops-db"
database_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"  # 替换这里

# KV 命名空间绑定
[[kv_namespaces]]
binding = "KV"
id = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"  # 替换这里

# R2 存储桶绑定
[[r2_buckets]]
binding = "R2"
bucket_name = "anixops-files"  # 不需要 ID，用名称
```

---

## 设置密钥 (Secrets)

```bash
# 设置 JWT 密钥 (必须)
wrangler secret put JWT_SECRET
# 输入一个强随机字符串，例如: openssl rand -base64 32

# 设置 API Key 盐值 (必须)
wrangler secret put API_KEY_SALT
# 输入一个随机字符串
```

---

## 快速命令汇总

```bash
# 1. 安装 wrangler
npm install -g wrangler

# 2. 登录 Cloudflare
wrangler login

# 3. 创建资源
wrangler d1 create anixops-db
wrangler kv:namespace create KV
wrangler r2 bucket create anixops-files

# 4. 设置密钥
wrangler secret put JWT_SECRET
wrangler secret put API_KEY_SALT

# 5. 更新 wrangler.toml 中的 ID

# 6. 运行迁移
cd workers
npm install
wrangler d1 migrations apply anixops-db

# 7. 本地测试
npm run dev

# 8. 部署
npm run deploy
```

---

## 查看 Dashboard URL

| 资源 | Dashboard URL |
|------|---------------|
| D1 | https://dash.cloudflare.com/ → Workers & Pages → D1 |
| KV | https://dash.cloudflare.com/ → Workers & Pages → KV |
| R2 | https://dash.cloudflare.com/ → R2 |
| Workers | https://dash.cloudflare.com/ → Workers & Pages |

---

## 验证部署

部署后，API 将在以下 URL 可用:

```
https://anixops-api.<你的子域>.workers.dev
```

测试健康检查:

```bash
curl https://anixops-api.<子域>.workers.dev/health
```

预期响应:

```json
{
  "status": "healthy",
  "version": "1.0.0",
  "timestamp": "2026-03-18T15:00:00.000Z",
  "environment": "production"
}
```