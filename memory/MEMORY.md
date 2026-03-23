# AnixOps Control Center - 记忆索引

## 项目概述
AnixOps Control Center 是一个统一的 TUI/GUI 控制中心，用于管理所有 AnixOps 产品。

## 当前版本: v2.2.0-rc.21

## 记忆文件

### 项目规划
- [cloudflare_advanced_features.md](cloudflare_advanced_features.md) - Cloudflare 炫技与 Web3.0 功能规划
- [project_progress.md](project_progress.md) - 项目进度跟踪
- [development_summary.md](development_summary.md) - 开发总结

## 关键技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Cloudflare Workers + Hono |
| 数据库 | D1 (SQLite) + KV + R2 |
| Web 前端 | Vue 3 + Pinia + Tailwind |
| Mobile/Desktop | Flutter + Riverpod |
| 实时通信 | SSE + WebSocket |
| AI | Workers AI + Vectorize |
| Web3 | IPFS Gateway + Ethereum Gateway |

## 测试统计
- Workers: 502 tests
- Web: 350 tests
- Mobile: 228 tests
- E2E: 7 tests
- **Total: 1087 tests**

## 架构

```
┌─────────────────┐          ┌─────────────────┐
│   Mobile App    │          │    Web App      │
│   (Flutter)     │          │    (Vue 3)      │
├─────────────────┤          ├─────────────────┤
│ • SSE Service   │◄────────►│ • useSSE.js     │
│ • WebSocket     │          │ • Pinia Stores  │
│ • API Client    │          │ • API Client    │
└────────┬────────┘          └────────┬────────┘
         │                            │
         │  REST + SSE + WebSocket    │
         ▼                            ▼
┌─────────────────────────────────────────────┐
│              Cloudflare Workers             │
├─────────────────────────────────────────────┤
│ • REST API (25 handlers)                    │
│ • SSE /api/v1/sse                           │
│ • WebSocket /api/v1/ws                      │
│ • Workers AI (LLM 推理)                     │
│ • Vectorize (向量数据库)                    │
│ • IPFS/Ethereum Gateway                     │
└─────────────────────────────────────────────┘
```