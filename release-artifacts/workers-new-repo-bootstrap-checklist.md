# Workers 新仓初始化清单

## 推荐仓库名
- `anixops-control-center-workers`

## 推荐目录结构

```text
anixops-control-center-workers/
├── .github/
│   └── workflows/
│       ├── ci.yml
│       └── release.yml
├── docs/
│   ├── architecture.md
│   ├── deploy.md
│   ├── api-contract.md
│   └── cloudflare-integration.md
├── migrations/
├── src/
│   ├── handlers/
│   ├── middleware/
│   ├── services/
│   ├── utils/
│   ├── index.ts
│   └── types.ts
├── test/
├── .gitignore
├── CHANGELOG.md
├── VERSION
├── README.md
├── SETUP.md
├── package.json
├── package-lock.json
├── tsconfig.json
├── vitest.config.ts
└── wrangler.toml
```

## 首批迁移文件

### 从主仓迁移
- `workers/src/**`
- `workers/test/**`
- `workers/migrations/**`
- `workers/package.json`
- `workers/package-lock.json`
- `workers/tsconfig.json`
- `workers/vitest.config.ts`
- `workers/wrangler.toml`
- `workers/README.md`
- `workers/SETUP.md`

### 从主仓文档提取
- `docs/cloudflare-integration-plan.md` 中 workers / Cloudflare 相关章节

## 新仓初始化时要新建
- `.gitignore`
- `CHANGELOG.md`
- `VERSION`
- `docs/architecture.md`
- `docs/deploy.md`
- `docs/api-contract.md`
- `.github/workflows/ci.yml`
- `.github/workflows/release.yml`

## 不应迁移的内容
- `workers/node_modules/**`
- `workers/.wrangler/**`
- `workers/coverage/**`
- `workers/node_modules/.vite/**`

## 首批 CI 能力
- install dependencies
- typecheck / lint
- vitest
- coverage upload

## 首批 release 能力
- staging deploy
- production deploy
- changelog / release notes
- Cloudflare bindings 校验
