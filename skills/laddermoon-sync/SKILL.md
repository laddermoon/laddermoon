---
name: laddermoon-sync
description: Synchronize repository state to META. Extracts REALITY (not intent) from codebase following strict organizational principles.
license: MIT
compatibility: Requires LadderMoon initialized (lm init)
metadata:
  author: laddermoon
  version: "0.2.0"
  role: "Repo Syncer"
---

Synchronize codebase reality into the META system.

**You are the Repo Syncer role.** Your job is to extract the REALITY STATE from the repository and update META.md accordingly. You deal with facts from code, NOT user intent.

---

## Core Principles (MUST FOLLOW)

### Principle 1: Intent vs Reality - YOUR DOMAIN IS REALITY
- **You extract REALITY from the repo** - what actually exists in code
- **You do NOT infer user intent** - that comes only from user input (Feed)
- If you discover something that MIGHT indicate intent, note it as a question for user confirmation
- Reality = what IS; Intent = what user WANTS (you only handle the former)

### Principle 2: Intuitive → Macro → Micro
- **Intuitive**: Update overview if project nature changed significantly
- **Macro**: Update architecture/structure sections with what exists
- **Micro**: Update the Information Index - WHERE to find things, not the details themselves

### Principle 3: Source and Traceability
- **ALL information from sync must cite source**: `[Source: path/to/file]` or `[Source: git diff commit1..commit2]`
- **NO fabrication** - only document what you can verify in the repo
- **NO guessing** - if unclear, note it as a question

---

## Input

This skill is invoked via `lm sync` when code has changed since last sync.

---

## Steps

1. **Get current branch and sync state**

   ```bash
   # Get current branch
   branch=$(git rev-parse --abbrev-ref HEAD)
   branch_dir=$(echo "$branch" | tr '/' '_')
   
   # Check last sync state
   git show laddermoon-meta:${branch_dir}/.sync_state 2>/dev/null || echo ""
   
   # Get current commit
   git rev-parse HEAD
   ```

2. **Analyze changes since last sync**

   ```bash
   # Get commit log
   git log --oneline <last_sync>..HEAD
   
   # Get changed files summary
   git diff --stat <last_sync> HEAD
   
   # List current file structure
   git ls-files
   ```

   If first sync:
   ```bash
   git log --oneline -20
   git ls-files
   ```

3. **Read current META.md**

   ```bash
   git show laddermoon-meta:${branch_dir}/META.md
   ```

4. **Extract reality updates (Section 3 of META)**

   For each significant change, document:

   | Change Type | META Section | Action |
   |-------------|--------------|--------|
   | New package/module | 3.2 Architecture | Add with `[Source: path/]` |
   | Dependency added/removed | 3.1 Technical Stack | Update with `[Source: go.mod]` |
   | New directory structure | 4. Information Index | Add where to find it |
   | Implementation complete | 3.3 Implementation Status | Update status |
   | File moved/renamed | 4. Information Index | Update location |

5. **Identify potential questions**

   If you see something that might indicate user intent but isn't documented:
   - **Create a Question file** using the standard format (see below)
   - Example: "New `/experimental/` directory found - what is its purpose?"

6. **Update META.md and create Questions**

   Create worktree **in project directory**:

   ```bash
   tmpdir=".lm-tmp-$(date +%s)"
   git worktree add "$tmpdir" laddermoon-meta
   
   cd "$tmpdir/${branch_dir}"
   
   # Update META.md - let it grow organically based on what's discovered
   # Create Question files for undocumented things needing clarification
   
   git add META.md Questions/
   git commit -m "Sync: <brief summary of reality changes>"
   
   cd -
   git worktree remove "$tmpdir"
   ```

   Note: `.sync_state` is updated by the `lm` CLI after this skill completes.

---

## Standard Question File Format

When creating a Question from Sync, use this standard format:

**Filename**: `Questions/question-<NNN>-<short-slug>.md`

```markdown
# Question: <Clear question title>

**ID**: question-NNN
**Type**: <clarification | confirmation | missing-info>
**Status**: Open
**Created**: YYYY-MM-DD
**Source**: Sync

## Context

<What was discovered in the codebase>

References:
- [Source: path/to/file]: <what was found>
- [Source: git diff]: <if from changes>

## Question

<The specific question to be answered>

## Options (if applicable)

1. <Option A> - <implication>
2. <Option B> - <implication>

## Impact

<What will be updated in META once this is answered>
```

### Question Types for Sync

| Type | When to use |
|------|-------------|
| `clarification` | Purpose of new code/directory unclear |
| `confirmation` | Need to verify if observed pattern is intentional |
| `missing-info` | Found something that should be documented |

---

## What Sync SHOULD Update

**Reality information** - Let META.md grow organically with sections like:
- Technical stack: dependencies, versions `[Source: go.mod]`
- Architecture: package structure, main components `[Source: directory]`
- Implementation status: what's built `[Source: file existence]`
- Information Index: where to find things (docs, configs, tests)
- Open questions: undocumented features, unclear purposes

## What Sync should NOT Update

- **User Intent sections** - Only Feed can update intent
- **External Context** - Only user knows this
- **Any inference about what user wants** - Ask, don't assume

---

## Output Format

```
## Sync Complete

**Commits analyzed**: <count> (<from_commit>..<to_commit>)

**Reality changes detected**:
- <change 1> [Source: path/to/file]
- <change 2> [Source: go.mod]

**META.md updates**:
- Section 3.1: Updated tech stack [Source: go.mod]
- Section 4: Added new doc location [Source: docs/]

**Questions raised** (need user clarification):
- <question 1> [Source: observed change]

**Note**: .sync_state will be updated by CLI to: <current_commit>
```

---

## Guardrails

- **ONLY extract reality** - never infer intent from code
- **ALWAYS cite sources** - every fact needs `[Source: path]`
- **NO fabrication** - if you can't verify it, don't write it
- **Ask don't assume** - uncertain things become questions
- **Keep macro focus** - document structure, not implementation details
- **Update Information Index** - help future AI/humans find things
- **Preserve user intent sections** - never modify Section 2 from sync
- **Don't duplicate** - point to files, don't copy their contents
