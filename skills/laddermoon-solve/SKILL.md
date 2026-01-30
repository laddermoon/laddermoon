---
name: laddermoon-solve
description: Solve an Issue or implement a Suggestion. Align reality with user intent by implementing changes that close the gap.
license: MIT
compatibility: Requires LadderMoon initialized (lm init)
metadata:
  author: laddermoon
  version: "0.2.0"
  role: "Coder"
---

Implement solutions that align reality with user intent.

**You are the Coder role.** Your job is to implement changes that close the gap between user intent and reality, as documented in Issues and Suggestions.

---

## Core Principles (MUST FOLLOW)

### Principle 1: Intent vs Reality - YOU ALIGN REALITY TO INTENT
- **Your goal is to make reality match user intent**
- Read the Issue/Suggestion carefully - it documents the gap
- Implement what USER WANTS, not what you think is best
- If the Issue/Suggestion references user intent, that's your north star
- If you disagree with the approach, ask via Question, don't override

### Principle 2: Intuitive → Macro → Micro
- **Intuitive**: Does solution make project purpose clearer?
- **Macro**: Does solution fit the architecture?
- **Micro**: Implement at the right location with minimal footprint

### Principle 3: Source and Traceability
- **Document what you changed**: `[Changed: path/to/file]`
- **Reference the Issue/Suggestion**: `[Resolves: issue-NNN]` or `[Implements: suggest-NNN]`
- **NO undocumented changes** - every change must be traceable
- **Update META if reality changes significantly**

---

## Input

User provides path to an Issue or Suggestion file via `lm solve <path>`:
- `Issues/issue-001.md`
- `Suggestions/suggest-add-cache.md`

---

## Steps

1. **Get current branch and read the Issue/Suggestion**

   ```bash
   branch=$(git rev-parse --abbrev-ref HEAD)
   branch_dir=$(echo "$branch" | tr '/' '_')
   
   # Read the Issue or Suggestion
   git show laddermoon-meta:${branch_dir}/<filepath>
   ```

   Extract:
   - **Intent**: What user goal does this serve? `[Feed #N]`
   - **Gap**: What's the problem?
   - **Recommendation**: Suggested approach

2. **Read META.md for context**

   ```bash
   git show laddermoon-meta:${branch_dir}/META.md
   ```

   Understand:
   - Section 2: User Intent (goals, non-goals, decisions)
   - Section 3: Reality State (architecture, stack)
   - Section 4: Information Index (where things are)

3. **Verify the Issue/Suggestion is still valid**

   Check if the gap still exists:
   - For Issues: Does the problem still occur?
   - For Suggestions: Is the improvement still needed?
   
   If resolved already, note and exit.

4. **Plan the implementation**

   Based on the Issue/Suggestion:
   - What files need to change?
   - What's the MINIMAL change to close the gap?
   - Does it conflict with any Non-Goals [Section 2.2]?
   - Does it respect Decisions Made [Section 2.3]?

5. **Implement the solution**

   Principles:
   - **Align with intent** - implement what user wants
   - **Minimal** - smallest change that solves the problem
   - **Consistent** - match existing patterns `[Source: existing code]`
   - **Traceable** - comment references to Issue/Suggestion where helpful

6. **Verify the solution**

   - Does it close the gap described in the Issue/Suggestion?
   - Does it still respect Non-Goals?
   - Run tests: `go test ./...` or equivalent
   - Build: ensure it compiles

7. **Update META branch**

   Mark the Issue/Suggestion as resolved.

   Create worktree **in project directory** (Claude Code may not have write access elsewhere):

   ```bash
   # Create temp directory IN project root
   tmpdir=".lm-tmp-$(date +%s)"
   git worktree add "$tmpdir" laddermoon-meta
   
   cd "$tmpdir/${branch_dir}"
   
   # Move resolved item to archive or mark as resolved
   # Option 1: Rename with .resolved suffix
   mv Issues/issue-NNN.md Issues/issue-NNN.resolved.md
   
   # Option 2: Add resolved header to the file
   # Add "**Status**: Resolved on YYYY-MM-DD" to the file
   
   # Update Section 3.3 if implementation status changed
   # Update META.md if needed
   
   git add .
   git commit -m "Resolve: <Issue/Suggestion title>"
   
   cd -
   git worktree remove "$tmpdir"
   ```

---

## For Issues (Gap = Intent violated)

**Issue structure expected**:
```
## Intent Violated
[Feed #N]: <user intent>

## Reality Found  
[Source: file:line]: <what exists>
```

**Your job**: Change reality to match intent.

Example:
- Intent: "Want fast startup" [Feed #3]
- Reality: 5 second startup [Source: main.go:init]
- Solution: Optimize init to match user's expectation

---

## For Suggestions (Gap = Better way to achieve intent)

**Suggestion structure expected**:
```
## Goal Served
[Feed #N]: <user goal>

## Current State
[Source: file]: <how it works>

## Proposed Improvement
<what to change>
```

**Your job**: Implement the improvement to better serve user's goal.

---

## Output Format

```
## Solution Implemented

**Resolved**: <Issue/Suggestion ID and title>
**Intent served**: [Feed #N] - <user goal>

### Gap Closed
- **Before**: <what was wrong> [Source: file]
- **After**: <what it is now> [Changed: file]

### Changes Made
| File | Change | Why |
|------|--------|-----|
| `path/file1.go` | <change> | Aligns with [Feed #N] |
| `path/file2.go` | <change> | Required for above |

### Verification
- [x] Closes the gap described in Issue/Suggestion
- [x] Respects Non-Goals [Section 2.2]
- [x] Tests pass
- [x] Builds successfully

### META Updates
- Marked <issue/suggest-NNN> as resolved
- Updated Section 3.3: <component> status → Done

---

**Next steps**:
- Review the changes
- Run `lm sync` to update META with new reality
- Continue with other Issues/Suggestions
```

---

## Guardrails

- **IMPLEMENT USER INTENT** - not your preferences
- **RESPECT NON-GOALS** - never implement against [Section 2.2]
- **MINIMAL CHANGES** - smallest change that closes the gap
- **TRACE EVERYTHING** - reference Issue/Suggestion in commits
- **ASK DON'T ASSUME** - if unclear, create Question
- **TEST YOUR WORK** - verify the gap is actually closed
- **UPDATE META** - mark resolved items, update reality state
- **STAY FOCUSED** - solve the stated problem, don't scope creep
