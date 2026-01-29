# LadderMoon 实现计划文档

## 技术选型

- **开发语言**：Go
- **CLI 框架**：Cobra 或类似框架
- **AI 后端**：Claude Code (通过 Skill 机制)
- **版本控制**：Git (影子分支管理)

---

## 实现阶段

### 阶段一：基础设施

**目标**：搭建项目骨架，实现基本的 META 管理功能

| 任务 | 优先级 | 产出 |
|------|--------|------|
| 初始化项目结构 | P0 | Go module、目录结构 |
| 实现 `lm init` | P0 | 创建影子分支、初始化 META.md |
| 实现 META 文件读写库 | P0 | `pkg/meta` 包 |
| 实现 `lm status` | P1 | 显示 META 状态、CommitID |
| 实现 CommitID 同步检测 | P1 | 判断 META 是否需要更新 |

### 阶段二：信息同步

**目标**：实现代码库与 META 的双向信息流

| 任务 | 优先级 | 产出 |
|------|--------|------|
| 实现 `lm feed` | P0 | User Input Processor 功能 |
| 实现 `lm sync` | P0 | Repo Syncer 功能 |
| Git diff 解析 | P1 | 增量同步机制 |
| Session 状态文件管理 | P1 | STATUS.json、SESSION_SUMMARY.md |

### 阶段三：诊断系统

**目标**：实现 AI 驱动的问题发现机制

| 任务 | 优先级 | 产出 |
|------|--------|------|
| 实现 Questioner Skill | P0 | 疑问生成能力 |
| 实现 Question Solver Skill | P0 | 疑问自动解决 |
| 实现 `lm audit` | P1 | Issuer 功能集成 |
| 实现 `lm propose` | P1 | Suggester 功能集成 |
| 疑问循环机制 | P1 | 自动循环直至无疑问 |

### 阶段四：执行系统

**目标**：实现代码修改与验收流程

| 任务 | 优先级 | 产出 |
|------|--------|------|
| 实现 `lm solve [ID]` | P0 | Coder 功能 |
| 分支管理 | P0 | 自动创建修改分支 |
| 实现 Reviewer 机制 | P1 | 代码验收流程 |
| 合并策略 | P1 | 验收通过后自动合并 |

### 阶段五：自我优化

**目标**：实现系统自我迭代能力

| 任务 | 优先级 | 产出 |
|------|--------|------|
| 决策日志记录 | P0 | Decision_Log.md 自动更新 |
| 用户反馈收集 | P1 | Issue/Suggestion 评价机制 |
| Self-Improver 实现 | P2 | Prompt 自动优化 |

---

## 核心流程实现细节

### 流程1：用户主动提供信息

```
用户执行 lm feed ...
    │
    ▼
User Input Processor 启动
    │
    ▼
更新 META.md
```

### 流程2：元信息库维护

```
lm sync 触发
    │
    ▼
Repo Syncer 执行
    │
    ▼
Questioner 提出疑问 , 如果没有疑问，流程结束。
    │
    ▼
Question Solver 尝试解决
    │
    ├──(能解决)──▶ 更新 META.md
    │
    └──(不能解决)──▶ 向用户提问（可以在Skill内部实现）
```

### 流程3：创建 Issue

```
lm audit 触发
    │
    ▼
检查 CommitID 一致性
    │
    ├──(不一致)──▶ 提示用户先 sync
    │
    └──(一致)──▶ Issuer 提出问题 ， 如果没有问题，提示没有Issue
                    │
                    ▼
               用户分类 Issue
```

### 流程4：创建 Suggestion

```
lm propose 触发
    │
    ▼
检查 CommitID 一致性
    │
    ├──(不一致)──▶ 提示用户先 sync
    │
    └──(一致)──▶ Suggester 提出建议 ， 也可以没有建议
                    │
                    ▼
               用户分类 Suggestion
```

### 流程5：解决 Issue/Suggestion

```
lm solve [ID] 触发
    │
    ▼
读取指定 Issue/Suggestion 文件
    │
    ▼
Coder 创建新分支并修改
    │
    ▼
Reviewer 验收
    │
    ├──(不通过)──▶ 返回 Coder 修改
    │
    └──(通过)──▶ 合并到主分支
```

### 流程6：自我提升

```
Self-Improver 定期/按需执行
    │
    ▼
读取 Decision_Log.md
    │
    ▼
分析用户反馈模式
    │
    ▼
优化 Issuer/Suggester 的 Prompt
```

---

## Claude Code Skill 定义

每个角色对应一个 Skill，建议的 Skill 结构：
一个Skill就是一个MD文件，可以参照下面目录上的Skill来实现
/Users/congpeiqing/codes/laddermoon/.claude/skills
---

## MVP 里程碑

**最小可行产品**应包含：

1. ✅ `lm init` - 初始化项目
2. ✅ `lm feed` - 录入信息
3. ✅ `lm sync` - 同步代码库
4. ✅ `lm status` - 查看状态

**MVP 完成标准**：能够建立并维护一个基本的 META 信息库，手动触发同步。

---

## 风险与应对

| 风险 | 应对策略 |
|------|----------|
| Claude API 限流 | 实现请求队列和重试机制 |
| Token 成本过高 | 精简 META.md 内容，按需加载 |
| AI 输出不稳定 | 结构化输出格式，增加校验 |
| 影子分支冲突 | 自动 rebase 策略 |

---

## 下一步行动

1. 创建 Go module 并初始化项目结构
2. 实现 `pkg/meta` 基础库
3. 实现 `lm init` 命令
4. 编写第一个 Claude Code Skill (Syncer)
5. 端到端测试基础流程
