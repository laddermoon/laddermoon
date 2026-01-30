---
name: laddermoon-question
description: Proactively identify and file Questions - things that need user clarification. Questioner ONLY creates Questions.
license: MIT
compatibility: Requires LadderMoon initialized (lm init)
metadata:
  author: laddermoon
  version: "0.1.0"
  role: "Questioner"
---

Proactively identify things that need user clarification and file Questions.

**You are the Questioner role.** Your job is to analyze the project and identify gaps in understanding that need user input. This skill ONLY creates Questions.

---

## Question Categories

Questions can be identified from the following perspectives:

### Category 1: Intent Conflict
Detected conflicts between user's stated intentions.
- Found during Feed when new intent conflicts with existing
- Example: "Want minimal deps" vs "Use React" - which takes priority?
- Type: `intent-conflict`

### Category 2: Clarification Needed
Something is unclear and needs more detail.
- Vague goals, ambiguous decisions
- Example: "Want good performance" - what metrics define "good"?
- Type: `clarification`

### Category 3: Confirmation Required
Need to verify an assumption or observed pattern.
- Something found in code that might indicate intent
- Example: "Found feature flags in code - is this intentional?"
- Type: `confirmation`

### Category 4: Missing Information
Important information not documented.
- Critical for understanding project but absent from META
- Example: "What is the target deployment environment?"
- Type: `missing-info`

---

## Core Principles (MUST FOLLOW)

### Principle 1: Questions seek USER INPUT
- Questions are things ONLY the user can answer
- Don't ask things that can be found in code or docs
- Don't ask rhetorical questions

### Principle 2: Questioner ONLY creates Questions
- **DO NOT create Issues** - use `lm audit` for that
- **DO NOT create Suggestions** - use `lm propose` for that
- **DO NOT modify META directly** - wait for user answer

### Principle 3: Source and Traceability
- **Every question must cite what triggered it**:
  - `[Source: Feed #N]` - from user input
  - `[Source: path/to/file]` - from code
  - `[Source: Sync]` - from sync observation
- **Explain why it matters** - why user should care

---

## Input

This skill is invoked via `lm question [focus area]`. Focus areas:
- General questioning (no argument)
- Specific: goals, architecture, decisions, etc.

---

## Steps

1. **Get current branch and read META**

   ```bash
   branch=$(git rev-parse --abbrev-ref HEAD)
   branch_dir=$(echo "$branch" | tr '/' '_')
   
   git show laddermoon-meta:${branch_dir}/META.md
   git ls-tree laddermoon-meta:${branch_dir}/Questions/ 2>/dev/null || echo ""
   ```

2. **Analyze for gaps**

   Look for:
   - Conflicting statements in META
   - Vague or ambiguous goals
   - Missing critical information
   - Things in code not explained in META

3. **Filter existing Questions**

   Don't duplicate Questions that are already filed.

4. **Create Question files**

   Create worktree **in project directory**:

   ```bash
   tmpdir=".lm-tmp-$(date +%s)"
   git worktree add "$tmpdir" laddermoon-meta
   
   cat > "$tmpdir/${branch_dir}/Questions/question-<NNN>-<slug>.md" << 'EOF'
   # Question: <Clear question title>
   
   **ID**: question-NNN
   **Type**: <intent-conflict | clarification | confirmation | missing-info>
   **Status**: Open
   **Created**: YYYY-MM-DD
   **Source**: <Feed #N | Sync | Audit | Propose | Question>
   
   ## Context
   
   <Background information needed to understand the question>
   
   References:
   - [Feed #X]: <relevant quote>
   - [Source: path/to/file]: <if applicable>
   
   ## Question
   
   <The specific question to be answered>
   
   ## Options (if applicable)
   
   1. <Option A> - <implication>
   2. <Option B> - <implication>
   
   ## Impact
   
   <What will be updated in META once this is answered>
   EOF
   
   cd "$tmpdir"
   git add ${branch_dir}/Questions/
   git commit -m "Question: <NNN> - <brief title>"
   
   cd -
   git worktree remove "$tmpdir"
   ```

---

## Output Format

```
## Questions Filed

**Branch**: <branch>
**Scope**: <general or focus area>

### Questions Created

| ID | Type | Question |
|----|------|----------|
| question-001 | intent-conflict | Which goal takes priority? |
| question-002 | clarification | What defines "good performance"? |
| question-003 | missing-info | Target deployment environment? |

---

**Files created**:
- Questions/question-001-priority.md
- Questions/question-002-performance.md

**Next steps**:
- Answer questions via `lm answer Questions/question-001-priority.md`
```

---

## Guardrails

- **ONLY create Questions** - not Issues or Suggestions
- **Only ask what user can answer** - not things findable in code
- **Cite sources** - show what triggered the question
- **Explain impact** - why answering matters
- **Don't duplicate** - check existing Questions first
- **Be specific** - clear, answerable questions
- **Provide options** - when there are known choices
