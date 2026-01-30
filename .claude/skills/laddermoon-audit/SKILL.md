---
name: laddermoon-audit
description: Detect and document Issues - problems that need to be fixed. Audit ONLY creates Issues, not Questions or Suggestions.
license: MIT
compatibility: Requires LadderMoon initialized and synced (lm init, lm sync)
metadata:
  author: laddermoon
  version: "0.3.0"
  role: "Issuer"
---

Analyze the project and identify genuine Issues that need to be fixed.

**You are the Issuer role.** Your job is to find and document real problems. This skill ONLY creates Issues - use `lm question` for Questions and `lm propose` for Suggestions.

---

## Issue Categories

Issues can be identified from the following perspectives:

### Category 1: Intent Self-Conflict
User's stated intentions conflict with each other.
- Example: "Want minimal dependencies [Feed #2]" vs "Use React framework [Feed #5]"
- Evidence: Both intents in META with `[Feed #N]` citations

### Category 2: Implementation vs Intent Conflict
What's implemented contradicts what user wants.
- **Only check implemented parts** - ignore unimplemented intentions (those are Suggestions)
- Example: User wants "REST API [Feed #3]" but code uses GraphQL [Source: api/]
- Evidence: Intent `[Feed #N]` + Reality `[Source: path]`

### Category 3: Project Reality Problems
Problems in current implementation itself (structural or detail).
- Code quality issues, architectural problems, bugs
- Example: Circular dependency between packages [Source: pkg/a imports pkg/b, pkg/b imports pkg/a]
- Evidence: `[Source: path]` showing the problem

### Category 4: META Information Gap
Important information missing that prevents understanding the project.
- Example: No overview of what the project does
- Example: Architecture described but key component locations missing
- Evidence: What's missing and why it matters

### Category 5: Verification Failure
Tests or validations fail when conditions are met.
- Run tests/validations as documented in META's Information Index
- Example: `go test ./...` fails [Source: test output]
- Evidence: Expected vs Actual results

---

## Core Principles (MUST FOLLOW)

### Principle 1: Issues are PROBLEMS to FIX
- Issues represent things that are **wrong** and need fixing
- Every issue must have clear evidence
- Cite sources for all claims

### Principle 2: Audit ONLY creates Issues
- **DO NOT create Questions** - use `lm question` for that
- **DO NOT create Suggestions** - use `lm propose` for that
- If unsure whether something is an issue, it's probably not - skip it

### Principle 3: Source and Traceability
- **Every issue must cite evidence**:
  - For intent conflicts: `[Feed #N]` vs `[Feed #M]`
  - For implementation issues: `[Source: path/to/file:line]`
  - For verification failures: `[Expected: X] [Actual: Y]`
- **NO fabricated issues** - only real problems you can prove

---

## Input

This skill is invoked via `lm audit [focus area]`. Focus areas:
- General audit (no argument)
- Specific: security, performance, architecture, etc.

---

## Steps

1. **Get current branch and verify sync**

   ```bash
   branch=$(git rev-parse --abbrev-ref HEAD)
   branch_dir=$(echo "$branch" | tr '/' '_')
   
   current=$(git rev-parse HEAD)
   synced=$(git show laddermoon-meta:${branch_dir}/.sync_state 2>/dev/null || echo "")
   ```

   If not synced, warn user to run `lm sync` first.

2. **Read META.md**

   ```bash
   git show laddermoon-meta:${branch_dir}/META.md
   ```

   Understand:
   - User goals and non-goals (intent)
   - Current reality state
   - Information index (how to verify)

3. **Audit each category**

   | Category | What to check | Evidence needed |
   |----------|---------------|-----------------|
   | Intent Self-Conflict | Compare all stated intents | `[Feed #N]` vs `[Feed #M]` |
   | Implementation vs Intent | Compare code with stated intent | `[Feed #N]` + `[Source: path]` |
   | Reality Problems | Code quality, bugs, structure | `[Source: path]` |
   | META Information Gap | What's missing in META | What's needed and why |
   | Verification Failure | Run tests per Information Index | `[Expected: X] [Actual: Y]` |

4. **Create Issue files**

   Create worktree **in project directory**:

   ```bash
   tmpdir=".lm-tmp-$(date +%s)"
   git worktree add "$tmpdir" laddermoon-meta
   
   cat > "$tmpdir/${branch_dir}/Issues/issue-<NNN>-<slug>.md" << 'EOF'
   # Issue: <Title>
   
   **ID**: issue-NNN
   **Category**: <intent-conflict | impl-vs-intent | reality-problem | meta-gap | verification-failure>
   **Severity**: Critical / High / Medium / Low
   **Status**: Open
   **Detected**: YYYY-MM-DD
   
   ## Problem
   
   <Clear description of what's wrong>
   
   ## Evidence
   
   <For intent-conflict>
   - [Feed #N]: <quote intent A>
   - [Feed #M]: <quote intent B>
   
   <For impl-vs-intent>
   - Intent: [Feed #N] says <what user wants>
   - Reality: [Source: path/to/file] shows <what exists>
   
   <For reality-problem>
   - [Source: path/to/file:line]: <the problem>
   
   <For meta-gap>
   - Missing: <what information is needed>
   - Impact: <why it matters>
   
   <For verification-failure>
   - Test: <what was run>
   - Expected: <what should happen>
   - Actual: <what happened>
   - [Source: test output or log]
   
   ## Recommendation
   
   <How to fix this issue>
   EOF
   
   cd "$tmpdir"
   git add ${branch_dir}/Issues/
   git commit -m "Audit: Issue <NNN> - <brief title>"
   
   cd -
   git worktree remove "$tmpdir"
   ```

---

## Issue Severity Guidelines

| Severity | Definition |
|----------|------------|
| **Critical** | Breaks core functionality or contradicts primary goal |
| **High** | Significant problem affecting major features |
| **Medium** | Notable issue that should be fixed |
| **Low** | Minor problem, can be fixed later |

---

## Output Format

```
## Audit Results

**Branch**: <branch>
**Scope**: <general or focus area>

### Issues Found

| ID | Category | Severity | Problem |
|----|----------|----------|---------|
| issue-001 | intent-conflict | High | Goals X and Y conflict |
| issue-002 | impl-vs-intent | Medium | Wants X, has Y |
| issue-003 | reality-problem | Low | Circular dependency |
| issue-004 | meta-gap | Medium | Missing architecture info |
| issue-005 | verification-failure | High | Tests fail |

### Verified OK
- Intent [Feed #1]: ✓ Implemented correctly
- Tests: ✓ All pass

---

**Files created**:
- Issues/issue-001-intent-conflict.md
- Issues/issue-002-wrong-api.md

**Next steps**:
- Fix issues via `lm solve Issues/issue-001-intent-conflict.md`
```

---

## Guardrails

- **ONLY create Issues** - not Questions or Suggestions
- **ALWAYS cite evidence** - every issue needs proof
- **NO opinion-based issues** - only real problems with evidence
- **NO fabrication** - only document verifiable issues
- **Be specific** - exact file paths and line numbers
- **Categorize correctly** - use the 5 categories above
- **Skip uncertain items** - if not sure it's an issue, don't create it
