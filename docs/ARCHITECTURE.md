# LadderMoon 架构设计文档

## 系统架构概览

LadderMoon 采用 **"分而治之，外部调度"** 的架构模式，将 AI 视为**无状态的计算函数**：给它上下文 → 输出结果并修改文件 → 关闭。

核心原则：**"Session 是临时的，数据（DNA）是永恒的。"**

---

## 角色系统

系统包含 9 个角色，每个角色是一个独立的工具/应用，具体实现为 Claude Code 中的 Skill。

| 角色 | 职责 |
|------|------|
| **User Input Processor** | 处理用户主动提供的项目信息，更新 META 文件 |
| **Repo Syncer** | 将 Repository 更新同步到 META 文件 |
| **Questioner** | 提出疑问 |
| **Question Solver** | 解决疑问 |
| **Issuer** | 提出问题 |
| **Suggester** | 提出改进建议 |
| **Coder** | 编写代码解决 Issue/Suggestion |
| **Reviewer** | 验收代码修改 |
| **Self-Improver** | 根据反馈优化 Issuer 和 Suggester |

---

## Session 规划策略

### 为何不使用单一长 Session

长 Session 在复杂项目管理中有三个致命伤：

1. **上下文偏移 (Context Drift)**：对话变长后，AI 会遗漏开头的关键约束
2. **Token 膨胀与成本**：每次对话都要重复发送整个历史
3. **状态不可控**：一个环节出错会污染后续环节

### 推荐方案：原子化 Session + 状态持久化

每个操作都一个原子Session，使用 claude "...." 这样的方式。
流程由lm驱动。

---

## 数据存储设计

### 影子分支 (Shadow Branch)

使用独立的 Git 分支 `laddermoon-meta` 存储 META 信息，不与主分支合并。

**优点：**
- 利用 Git 自身功能管理文档版本
- META 更新与代码提交解耦
- 支持 CommitID 同步机制

### CommitID 同步机制

META 文件记录最后同步的 CommitID：
- 若与 Repo 的 CommitID 一致 → META 有效
- 若不一致 → 触发增量同步，仅处理未同步的 commits

### 核心文件结构

```
laddermoon-meta 分支
├── META.md              # 项目元信息主文件
├── .sync_state          # 记录最后同步的 CommitID
├── UserFeed.log         # 用户 Feed 的原始记录
├── Questions/           # 待解决的疑问
│   └── .gitkeep
├── Issues/              # 已识别的问题
│   └── .gitkeep
└── Suggestions/         # 改进建议
    └── .gitkeep
```
---

## 代码结构设计

```
LadderMoon/
├── cmd/
│   └── lm/                   # CLI 入口
│       ├── main.go
│       └── cmd/
│           ├── root.go       # 根命令定义
│           ├── init.go       # lm init
│           ├── status.go     # lm status
│           ├── feed.go       # lm feed
│           ├── sync.go       # lm sync
│           ├── audit.go      # lm audit (Issuer)
│           ├── propose.go    # lm propose (Suggester)
│           ├── solve.go      # lm solve (Coder)
│           └── version.go    # lm version
├── pkg/
│   └── meta/                 # META 文件操作库
│       └── meta.go           # Git 操作、文件读写、同步状态管理
├── docs/
│   ├── PRODUCT.md            # 产品文档
│   ├── ARCHITECTURE.md       # 架构设计文档
│   └── IMPLEMENTATION.md     # 实现计划文档
└── README.md                 # 项目文档
```

---

## 架构优势

1. **容错性**：单个环节失败可单独重启，不影响其他环节
2. **可测试性**：可单独测试每个角色，使用伪造数据验证表现
3. **自迭代空间**：Self-Improver 可异步运行，扫描历史 Session Log 优化 Prompt
4. **成本可控**：避免 Token 膨胀，按需加载上下文
