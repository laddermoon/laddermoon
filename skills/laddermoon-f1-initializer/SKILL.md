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

3. **Explore the repository with curiosity**

   Approach the repository like a curious person discovering a new project. Follow the organization principles:

   **Intuitive Level** - What is this project:
   - Look for project description files (README, package manifests, etc.)
   - Understand: What is this? What problem does it solve? Who uses it?
   - **Must cite source for every claim**

   **Macro Level** - Complete structural understanding:
   
   **Critical: Directory Structure & Purpose**
   - List all major directories at the root level
   - For each directory, determine its purpose:
     - What kind of files are in it?
     - What role does it play in the project?
     - Is it source code, documentation, configuration, tests, build artifacts?
   - If a directory's purpose is unclear from its contents, note it for deeper exploration
   - Explore subdirectories when needed to understand the structure
   
   **Technical Stack & Architecture**:
   - Identify languages, frameworks, dependencies
   - Understand the logical architecture: How is the code organized?
   - What are the main modules/packages/components?
   - How do they relate to each other?
   
   **Core Implementation Details**:
   - What are the key entry points? (main files, CLI commands, API endpoints)
   - Where are the most important features implemented?
   - What are the critical code paths?
   
   **Usage & Operations**:
   - How do you run this project?
   - How do you build/test/deploy it?
   - What configuration is needed?
   
   **Code Standards & Conventions**:
   - Are there coding standards documented?
   - What patterns or conventions are followed?
   - Are there linting/formatting rules?

   **Micro Level** - Information Index:
   - Document WHERE to find specific information
   - Point to files/directories, don't copy their contents
   - **Must cite source for each entry**

   **Exploration Approach**: Start broad, then go deeper as needed. If something is unclear at one level, explore the next level down until you understand it.

4. **Create meta.md following the organization principles**

   Structure the metadata to answer the questions a curious person would have:

   ```markdown
   # Project Metadata

   > Last updated: YYYY-MM-DD
   > Commit: <current-commit-id>

   ## 1. Overview (Intuitive)

   <One paragraph: What is this project? What problem does it solve?>

   [Source: <file>]

   ## 2. Structure (Macro)

   ### 2.1 Directory Structure

   **REQUIRED**: List all major directories and their purposes:

   | Directory | Purpose | Key Contents |
   |-----------|---------|-------------|
   | <path> | <what it's for> | <what's in it> |

   [Source: directory exploration]

   ### 2.2 Technical Stack

   - **Languages**: <languages used> [Source: <file>]
   - **Key Dependencies**: <main dependencies> [Source: <file>]
   - **Build/Runtime**: <how it runs> [Source: <file>]

   ### 2.3 Architecture

   <How is the code logically organized? What are the main modules/components?>

   [Source: <files analyzed>]

   ### 2.4 Core Implementation

   - **Entry Points**: <main files, CLI commands, etc.> [Source: <file>]
   - **Key Features**: <where critical functionality is implemented> [Source: <code exploration>]

   ### 2.5 Usage

   - **How to Run**: <commands/steps> [Source: <file>]
   - **How to Build**: <build process> [Source: <file>]
   - **How to Test**: <testing approach> [Source: <file>]

   ### 2.6 Standards & Conventions

   <Coding standards, patterns, conventions followed> [Source: <file>]

   ## 3. Information Index (Micro)

   | Topic | Location | Format |
   |-------|----------|--------|
   | <topic> | <path> | <format> |

   [Source: directory analysis]

   ## 4. Notes

   <Any observations, uncertainties, or areas needing clarification>
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
