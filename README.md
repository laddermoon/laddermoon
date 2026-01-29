🧬 LadderMoon (lm)
Climb the code, reach the intent.
LadderMoon 是一个 AI 原生的元架构引擎。它通过在 Git 影子分支（Shadow Branch）中构建项目的数字孪生，实现“意图”与“代码”的深度对齐。

它不仅仅是一个开发助手，它是 AI AS ME —— 你的数字化架构分身，在代码库中投射你的意志。

🌙 核心哲学：AI AS ME
在传统的开发中，代码是肉身，意图是灵魂。随着项目膨胀，灵魂往往会迷失。

LadderMoon 通过三个维度重塑开发体验：

不仅仅是代码生成： 它是对你决策逻辑的捕捉与复刻。

不仅仅是文档： 它是存储在 .ai-shadow 分支中的动态 DNA。

不仅仅是助手： 它是你的 Digital Twin，以你的视角审计代码，以你的偏好驱动进化。

🚀 核心功能
🧬 META 系统：通过影子分支维护项目的元数据，不污染主分支代码。

🔄 Smart Sync：自动分析 Git Diff，将代码变更实时映射到架构意图（DNA.md）。

🛡️ Architectural Audit：基于项目 META 自动探测潜在的架构性风险与偏差。

🧠 Self-Evolution：通过 Decision Log 学习用户的否定与反馈，实现 AI 的个性化成长。

🛠️ 安装与起步
安装
目前支持通过 Go 环境直接编译安装：

Bash

go install github.com/LadderMoon/LadderMoon/cmd/lm@latest
快速开始
初始化项目基因：

Bash

lm init
这会创建影子分支 ai-shadow 并初始化 META 结构。

同步当前进度：

Bash

lm sync
将本地代码变动同步至 META 库，确保 AI “理解”最新的修改。

探测风险与建议：

Bash

lm audit    # 发现问题 (Issuer)
lm propose  # 获取改进建议 (Suggester)
处理任务：

Bash

lm solve [ISSUE_ID]
📂 角色定义 (The 9 Skills)
LadderMoon 内部集成了 9 个专业化角色，共同维护项目的生命周期：

Syncer: 仓库同步专家。

Questioner / Solver: 消除认知模糊的闭环。

Issuer / Suggester: 架构审计与优化提案。

Coder / Reviewer: 意图到代码的高质量实现。

Processor: 用户信息的输入接口。

Self-Improver: 系统的自我进化中枢。

🎨 界面预览
当你执行 lm init 时，你将看到 LadderMoon 的苏醒：

Plaintext

╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║  █    ██▀▀▀█ █▀▀▀█ █▀▀▀▄ ██▀▀▀ ██▀▀▀█ █▀▀▀█▄ █▀▀▀█ █▀▀▀█ █▄  █  ║
║  █    █▄▄▄█ █▄▄▄█ █   █ █▄▄▄  █▄▄▄█ █   █ █   █ █   █ █ █ █  ║
║  █    █   █ █   █ █   █ █     █   █ █   █ █   █ █   █ █  ▀█  ║
║  █▄▄▄ █   █ █   █ █▄▄▄▀ ██▄▄▄ █   █ █▄▄▄█▀ █▄▄▄█ █▄▄▄█ █   █  ║
║                                                              ║
║                    LadderMoon CLI v1.0.0                     ║
║              AI AS ME: Your Architectural Twin               ║
╚══════════════════════════════════════════════════════════════╝
⚖️ 开源协议
基于 MIT License 开源。

"This tool is designed to liberate human intent from mechanical coding. Climb the code, reach the intent."