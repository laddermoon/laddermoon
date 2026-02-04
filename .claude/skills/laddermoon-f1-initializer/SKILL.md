---
name: laddermoon-f1-initializer
description: Initialize project metadata by extracting information from the repository and creating meta.md. Use when starting metadata tracking for a project.
license: MIT
compatibility: Requires Git repository
metadata:
  author: laddermoon
  version: "1.0"
---

Initialize project metadata from the repository.

**Input**: None required. The skill operates on the current repository.

**Steps**

1. **Check if metadata already exists**

   ```bash
   ls -la .laddermoon/
   ```

   If `.laddermoon/meta.md` exists, inform user and suggest using `laddermoon-f1-synchronizer` instead.

2. **Create the metadata directory structure**

   ```bash
   mkdir -p .laddermoon/questions
   mkdir -p .laddermoon/issues
   ```

3. **Analyze the repository**

   Gather information following the organization principles:

   **Intuitive Level** - What is this project:
   - Read `README.md`, `package.json`, `go.mod`, `Cargo.toml`, or similar project files
   - Extract project name, description, purpose
   - **Must cite source**: `[Source: README.md]`

   **Macro Level** - Complete information structure:
   ```bash
   # Get project structure
   find . -type f -name "*.md" | head -20
   ls -la
   git ls-files | head -50
   
   # Get tech stack
   cat package.json 2>/dev/null || cat go.mod 2>/dev/null || cat Cargo.toml 2>/dev/null
   
   # Get recent history
   git log --oneline -10
   ```

   **Micro Level** - Where to find information:
   - Document locations of key files (docs, configs, tests)
   - Point to information sources, not copy their contents
   - **Must cite source for each entry**

4. **Create meta.md following the organization principles**

   ```markdown
   # Project Metadata

   > Last updated: YYYY-MM-DD
   > Commit: <current-commit-id>

   ## 1. Overview (Intuitive)

   <One paragraph summary of what this project is>

   [Source: README.md]

   ## 2. Structure (Macro)

   ### 2.1 Technical Stack

   - **Language**: <language> [Source: <file>]
   - **Dependencies**: <key dependencies> [Source: <file>]

   ### 2.2 Architecture

   <High-level structure description>

   [Source: directory structure analysis]

   ### 2.3 Key Components

   | Component | Location | Purpose |
   |-----------|----------|--------|
   | <name> | <path> | <brief description> |

   [Source: <files analyzed>]

   ## 3. Information Index (Micro)

   | Topic | Location | Format |
   |-------|----------|--------|
   | Documentation | <path> | Markdown |
   | Configuration | <path> | <format> |
   | Tests | <path> | <format> |

   [Source: directory analysis]

   ## 4. Notes

   <Any observations that don't fit above categories>
   ```

5. **Write metadata files**

   Create `.laddermoon/meta.md` with the gathered information.

6. **Record current commit**

   ```bash
   git rev-parse HEAD > .laddermoon/meta-commit-id
   ```

7. **Initialize question/issue tracking files**

   ```bash
   echo "0" > .laddermoon/questions/last-question-id
   touch .laddermoon/questions/question-meta.jsonl
   echo "0" > .laddermoon/issues/last-issue-id
   touch .laddermoon/issues/issue-meta.jsonl
   ```

**Output**

```
## Initialization Complete

**Metadata created at**: .laddermoon/meta.md
**Synced to commit**: <commit-id>

### Information Extracted

- **Overview**: <brief summary>
- **Tech Stack**: <languages/frameworks>
- **Key Components**: <count> identified
- **Information Index**: <count> entries

### Sources Cited

- <list of files used as sources>

### Next Steps

1. Run `laddermoon-f1-inspector` to check metadata clarity
2. Run `laddermoon-f1-synchronizer` after code changes
```

**Organization Principles (MUST FOLLOW)**

| Principle | Description |
|-----------|-------------|
| **Current state only** | Metadata reflects what IS, not what was or will be |
| **Intuitive → Macro → Micro** | Structure information from overview to details |
| **All sources cited** | Every piece of information must have `[Source: path]` |
| **No fabrication** | Only document what can be verified in the repo |
| **No inference** | Don't guess user intent or project direction |
| **Point, don't copy** | For micro level, indicate WHERE to find info, not the details |

**Guardrails**

- Never create metadata if it already exists
- Always cite the source file for every claim
- Never infer or assume information not in the repo
- Keep the overview genuinely intuitive (one paragraph)
- Information Index should point to files, not duplicate content
- If something is unclear, leave it out (inspector will catch it later)
