# Workers 独立仓库迁移实施指南

## 目标
将 `workers/` 从当前聚合仓库中拆分为独立 git 仓库，形成独立的 Cloudflare Workers API 后端产品线，同时保证：
- 主仓 Web / Mobile / Go 客户端不受影响
- Cloudflare production / staging 配置可迁移
- CI/CD、版本与文档职责清晰分离
- 整个迁移具备可回滚能力

---

## 当前前置条件

### 发布与环境状态
- production release workflow 已完成
- dev release workflow 已完成
- production PR 已合并
- `api.anixops.com` 当前 health / readiness 可用

### 当前 Cloudflare Workers 关键配置
配置来源：`workers/wrangler.toml`
- Worker service: `anixops-api-v2`
- Route: `api.anixops.com`
- Production D1: `anixops-db`
- Production KV: `35a39ccc4cea47e89aa8b1459f85e1cf`
- Production R2: `anixops-files`
- AI binding: `AI`
- Vectorize index: `anixops-embeddings`
- Staging env 已存在独立 name / D1 / KV / R2 配置框架

### 当前测试基线
- Workers full suite: `565 passed`
- Web full unit suite: `358 passed`
- Mobile full unit suite: `233 passed`
- Mobile Android integration: `5 passed`

---

## 推荐迁移策略

采用 **双仓过渡迁移**：
1. 先建立新仓并跑通 staging
2. 再修改主仓 CI / release / 文档
3. 进入双仓并行观察期
4. 最后从主仓移除 `workers/`

不建议一次性硬切，以避免 CI、文档和 Cloudflare 配置同时断裂。

---

## Phase 0：冻结与基线

### 冻结范围
仅允许阻断型修复：
- `workers/src/**`
- `workers/test/**`
- `workers/migrations/**`
- `workers/wrangler.toml`

### 记录快照
- `workers/wrangler.toml`
- `.github/workflows/ci.yml`
- `.github/workflows/release.yml`
- `web/src/api/index.js`
- `mobile/lib/core/services/api_client.dart`
- `mobile/lib/core/providers/api_providers.dart`
- `README.md`
- `CHANGELOG.md`

---

## Phase 1：新仓初始化

### 推荐仓库名
- `anixops-control-center-workers`（推荐）
- 备选：`anixops-workers`

### 首批必须迁移文件
#### 代码
- `workers/src/index.ts`
- `workers/src/types.ts`
- `workers/src/handlers/**`
- `workers/src/middleware/**`
- `workers/src/services/**`
- `workers/src/utils/**`

#### 测试
- `workers/test/**`

#### 数据库与配置
- `workers/migrations/**`
- `workers/package.json`
- `workers/package-lock.json`
- `workers/tsconfig.json`
- `workers/vitest.config.ts`
- `workers/wrangler.toml`

#### 文档
- `workers/README.md`
- `workers/SETUP.md`
- `docs/cloudflare-integration-plan.md` 中 workers 相关内容（建议迁入新仓 docs）

### 新仓首批建议新增文件
- `.gitignore`
- `CHANGELOG.md`
- `VERSION`
- `.github/workflows/ci.yml`
- `.github/workflows/release.yml`
- `docs/architecture.md`
- `docs/deploy.md`
- `docs/api-contract.md`

### 不应迁移的产物
- `workers/node_modules/**`
- `workers/.wrangler/**`
- `workers/coverage/**`
- `workers/node_modules/.vite/**`

---

## Phase 2：新仓 CI/CD 与 Cloudflare 落地

### 新仓 CI 最小集合
- install dependencies
- lint / typecheck
- workers tests
- coverage upload / threshold

### 新仓 Release / Deploy 最小集合
- staging deploy
- production deploy
- wrangler 配置校验
- release notes / changelog 发布

### Cloudflare 配置迁移顺序
1. 复制 `wrangler.toml`
2. 在新仓绑定 staging / production secrets
3. 校验 D1 / KV / R2 / AI / Vectorize
4. staging deploy + 冒烟
5. production deploy 切换

---

## Phase 3：主仓拆耦

### 必改文件
- `.github/workflows/ci.yml`
- `.github/workflows/release.yml`
- `scripts/version.sh`
- `README.md`
- `CHANGELOG.md`

### 必改内容
#### CI
删除：
- `working-directory: workers`
- `cache-dependency-path: workers/package-lock.json`
- `./workers/coverage/lcov.info`
- `workers-test` 对 quality gate 的直接依赖

#### Release
主仓 release 仅保留：
- Web
- Mobile
- CLI / TUI
- 聚合产品级说明

#### 版本
- workers 不再跟随主仓统一版本脚本推进
- 新仓独立维护 Workers 版本与 changelog

#### 文档
- 主仓 README 仅保留“API 由独立 workers 仓提供”
- 修复 workers README 中对 `../docs/...` 的相对引用

---

## Phase 4：双仓并行期

### 目标
- 新仓负责 Workers staging / prod 发布
- 主仓仅消费 API
- 保持关键客户端路径稳定：`https://api.anixops.com/api/v1`

### 必做验证
- Web API 调用正常
- Mobile API / SSE 正常
- health / readiness 正常
- auth / realtime / AI / Vectorize / Web3 / IPFS 正常

### 建议观察期
- 7 天

---

## Phase 5：最终切换

### 主仓保留
- `web/`
- `mobile/`
- Go 主程序
- 聚合文档
- 指向 workers 新仓的集成说明

### 主仓移除
- `workers/` 源码实现目录
- workers 专属 CI / release 逻辑

---

## 风险

1. 主仓 CI 断裂
2. 文档断链
3. Cloudflare 配置漂移
4. 版本叙事混乱
5. API 契约漂移导致客户端隐性故障

---

## 回滚策略

- 双仓并行期保留主仓 `workers/` 只读快照
- 若新仓 staging/prod 部署失败：保持旧 route / 旧发布路径
- 若切换后 production 出现异常：回切 Cloudflare route 或 worker version

---

## 验证清单

### 新仓验证
- CI 全绿
- staging deploy 成功
- 关键 API 冒烟通过

### 主仓验证
- 不再引用 `workers/*`
- Web / Mobile 通过新仓 API 正常工作
- 文档链接无断裂

### 切换验收
- production / staging 配置一致
- `api.anixops.com` 指向新仓发布链
- 观察期内无关键路径回归
