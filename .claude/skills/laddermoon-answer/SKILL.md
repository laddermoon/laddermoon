---
name: laddermoon-answer
description: Process user's answer to a Question, update META accordingly, and resolve conflicts. Answer resolves Questions and updates META.
license: MIT
compatibility: Requires LadderMoon initialized (lm init)
metadata:
  author: laddermoon
  version: "0.1.0"
  role: "Question Solver"
---

Process user's answer to a Question and update META to reflect the resolution.

**You are the Question Solver role.** Your job is to take user's answer to a filed Question, update META.md accordingly, and mark the Question as resolved.

---

## What Answer Does

1. **Reads the Question file** to understand what was asked
2. **Processes user's answer** (provided in the prompt)
3. **Updates META.md** to reflect the resolution
4. **Resolves conflicts** if the Question was about intent conflicts
5. **Marks Question as resolved**

---

## Core Principles (MUST FOLLOW)

### Principle 1: User Answer is Authoritative
- The user's answer is the final word
- Update META to accurately reflect user's decision
- Remove conflict markers once resolved

### Principle 2: Traceability
- Record the answer with source: `[Feed #N]` (the answer is a form of user input)
- Reference which Question was resolved: `[Resolves: question-NNN]`
- Remove `[CONFLICT: see question-NNN]` markers when resolved

### Principle 3: Complete the Loop
- Update META.md with the resolved information
- Mark Question status as Resolved
- Clean up any conflict markers

---

## Input

This skill is invoked via `lm answer <question-file> "<user's answer>"`.

Example:
```
lm answer Questions/question-001-priority.md "Minimal dependencies is more important, so no React"
```

The prompt includes:
- Path to the Question file
- User's answer text

---

## Steps

1. **Get current branch and read the Question**

   ```bash
   branch=$(git rev-parse --abbrev-ref HEAD)
   branch_dir=$(echo "$branch" | tr '/' '_')
   
   git show laddermoon-meta:${branch_dir}/<question-file>
   ```

   Understand:
   - What was the question?
   - What type was it? (intent-conflict, clarification, etc.)
   - What was the expected impact?

2. **Read current META.md**

   ```bash
   git show laddermoon-meta:${branch_dir}/META.md
   ```

   Identify:
   - Where the conflict/unclear parts are
   - What needs to be updated

3. **Process the user's answer**

   Based on question type:
   
   | Question Type | Action on Answer |
   |---------------|------------------|
   | `intent-conflict` | Keep winning intent, remove or note losing intent |
   | `clarification` | Add clarified detail to META |
   | `confirmation` | Confirm or deny, update accordingly |
   | `missing-info` | Add the new information to META |

4. **Update META.md**

   - Remove `[CONFLICT: see question-NNN]` markers
   - Update the relevant sections with resolved information
   - Add source citation for the answer: `[Feed #N]` or `[Answer: question-NNN]`

5. **Mark Question as Resolved**

   Update the Question file:
   - Change `**Status**: Open` to `**Status**: Resolved`
   - Add resolution section

6. **Commit changes**

   Create worktree **in project directory**:

   ```bash
   tmpdir=".lm-tmp-$(date +%s)"
   git worktree add "$tmpdir" laddermoon-meta
   
   cd "$tmpdir/${branch_dir}"
   
   # Update META.md with resolved information
   # Update Question file with resolution
   
   git add META.md Questions/
   git commit -m "Answer: Resolve question-NNN - <brief summary>"
   
   cd -
   git worktree remove "$tmpdir"
   ```

---

## Question Resolution Format

Update the Question file to mark resolution:

```markdown
# Question: <Title>

**ID**: question-NNN
**Type**: <type>
**Status**: Resolved
**Created**: YYYY-MM-DD
**Resolved**: YYYY-MM-DD

## Context
<original context>

## Question
<original question>

## Resolution

**Answer**: <user's answer>
**Action taken**: <what was updated in META>
**Source**: [Feed #N] or direct answer
```

---

## Output Format

```
## Question Resolved

**Question**: question-NNN - <title>
**Type**: <type>
**Answer**: <user's answer>

### META Updates
- <section>: <what changed>
- Removed conflict marker: [CONFLICT: see question-NNN]

### Question File Updated
- Status: Open → Resolved
- Added resolution section

---

**Next steps**:
- Run `lm sync` if code changes are needed
- Continue with `lm audit` or `lm propose`
```

---

## Guardrails

- **Honor user's answer** - don't second-guess or modify
- **Update META completely** - remove all conflict markers
- **Cite the resolution** - trace back to user's answer
- **Mark Question resolved** - don't leave it Open
- **Be precise** - update exactly what the Question impacted
- **Create worktree in project directory** - use `.lm-tmp-*` prefix
