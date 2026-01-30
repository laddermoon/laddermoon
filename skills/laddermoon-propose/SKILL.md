---
name: laddermoon-propose
description: Propose Suggestions - improvements that help achieve user goals better. Propose ONLY creates Suggestions.
license: MIT
compatibility: Requires LadderMoon initialized and synced (lm init, lm sync)
metadata:
  author: laddermoon
  version: "0.3.0"
  role: "Suggester"
---

Propose improvements that help achieve user's stated goals better.

**You are the Suggester role.** Your job is to suggest improvements. This skill ONLY creates Suggestions - use `lm question` for Questions and `lm audit` for Issues.

---

## Suggestion Categories

Suggestions can be proposed from the following perspectives:

### Category 1: Unimplemented Intent
User intent that hasn't been implemented yet.
- Look at stated goals that don't have corresponding implementation
- Example: User wants "CLI help command [Feed #3]" but it doesn't exist yet
- Evidence: `[Feed #N]` intent + absence in code

### Category 2: Better Implementation
Current implementation works but could better achieve user intent.
- Improvements that more effectively serve stated goals
- Example: User wants "fast startup [Feed #2]", current is 3s, could be 1s
- Evidence: `[Serves: Feed #N]` + `[Current: path]` + proposed improvement

### Category 3: Implementation Improvement
Improvements to current implementation (not directly tied to specific intent).
- Better code quality, performance, maintainability
- Example: Refactor duplicated code in handlers
- Evidence: `[Current: path]` + what can be improved

---

## Core Principles (MUST FOLLOW)

### Principle 1: Suggestions are IMPROVEMENTS
- Suggestions make things better, not fix things that are wrong (that's Issues)
- Every suggestion should have clear benefit
- Cite what goal it serves or what it improves

### Principle 2: Propose ONLY creates Suggestions
- **DO NOT create Questions** - use `lm question` for that
- **DO NOT create Issues** - use `lm audit` for that
- **Respect Non-Goals** - never suggest things user explicitly rejected

### Principle 3: Source and Traceability
- **Every suggestion must cite**:
  - For unimplemented intent: `[Serves: Feed #N]`
  - For better implementation: `[Serves: Feed #N]` + `[Current: path]`
  - For implementation improvement: `[Improves: path]`
- **NO unsolicited suggestions against non-goals**

---

## Input

This skill is invoked via `lm propose [focus area]`. Focus areas:
- General suggestions (no argument)
- Specific: performance, DX, testing, etc.

---

## Steps

1. **Get current branch and verify sync**

   ```bash
   branch=$(git rev-parse --abbrev-ref HEAD)
   branch_dir=$(echo "$branch" | tr '/' '_')
   
   current=$(git rev-parse HEAD)
   synced=$(git show laddermoon-meta:${branch_dir}/.sync_state 2>/dev/null || echo "")
   ```

2. **Read META.md**

   ```bash
   git show laddermoon-meta:${branch_dir}/META.md
   ```

   Understand:
   - User goals (what to implement/improve)
   - Non-goals (what NOT to suggest)
   - Current reality state

3. **Check existing Suggestions**

   ```bash
   git ls-tree laddermoon-meta:${branch_dir}/Suggestions/
   ```

   Don't duplicate existing items.

4. **Identify suggestions for each category**

   | Category | What to look for | Evidence needed |
   |----------|------------------|-----------------|
   | Unimplemented Intent | Goals without implementation | `[Serves: Feed #N]` |
   | Better Implementation | Working code that could be better | `[Serves: Feed #N]` + `[Current: path]` |
   | Implementation Improvement | Code that can be improved | `[Improves: path]` |

5. **Filter against Non-Goals**

   Before proposing, check non-goals in META:
   - If user said "Don't care about test coverage" → don't suggest more tests
   - If user said "No external dependencies" → don't suggest new libraries

6. **Create Suggestion files**

   Create worktree **in project directory**:

   ```bash
   tmpdir=".lm-tmp-$(date +%s)"
   git worktree add "$tmpdir" laddermoon-meta
   
   cat > "$tmpdir/${branch_dir}/Suggestions/suggest-<NNN>-<slug>.md" << 'EOF'
   # Suggestion: <Title>
   
   **ID**: suggest-NNN
   **Category**: <unimplemented-intent | better-impl | impl-improvement>
   **Impact**: High / Medium / Low
   **Effort**: High / Medium / Low
   **Status**: Open
   **Proposed**: YYYY-MM-DD
   
   ## What
   
   <Clear description of the suggestion>
   
   ## Why
   
   <For unimplemented-intent>
   - Serves: [Feed #N] - <quote the intent>
   - Currently: Not implemented
   
   <For better-impl>
   - Serves: [Feed #N] - <quote the intent>
   - Current: [Source: path/to/file] - <how it works now>
   - Better: <how it could be improved>
   
   <For impl-improvement>
   - Improves: [Source: path/to/file]
   - Benefit: <what gets better>
   
   ## How
   
   <High-level implementation approach>
   
   ## Trade-offs
   
   <Any downsides or costs>
   EOF
   
   cd "$tmpdir"
   git add ${branch_dir}/Suggestions/
   git commit -m "Propose: Suggestion <NNN> - <brief title>"
   
   cd -
   git worktree remove "$tmpdir"
   ```

---

## Suggestion Prioritization

| Priority | Criteria |
|----------|----------|
| **P0 - Critical** | Directly enables core goal |
| **P1 - High** | Significantly improves goal achievement |
| **P2 - Medium** | Moderately helps stated goals |
| **P3 - Low** | Minor improvement |

Impact/Effort matrix:
- **Quick Win**: High impact, Low effort → Suggest first
- **Strategic**: High impact, High effort → Worth discussing
- **Fill-in**: Low impact, Low effort → When time permits
- **Avoid**: Low impact, High effort → Don't suggest

---

## Output Format

```
## Suggestions for Branch: <branch>

**Non-Goals Respected**:
- [Feed #2]: <non-goal> (no suggestions made against this)

### Suggestions Created

| ID | Category | Suggestion | Impact/Effort |
|----|----------|------------|---------------|
| suggest-001 | unimplemented-intent | <title> | High/Low |
| suggest-002 | better-impl | <title> | Medium/Medium |
| suggest-003 | impl-improvement | <title> | Low/Low |

---

**Files created**:
- Suggestions/suggest-001-add-help.md
- Suggestions/suggest-002-faster-startup.md

**Next steps**:
- Run `lm solve Suggestions/<file>` to implement
```

---

## Guardrails

- **ONLY create Suggestions** - not Questions or Issues
- **NEVER suggest against non-goals** - respect what user explicitly rejected
- **Cite evidence** - every suggestion needs references
- **Be specific** - concrete improvements with clear implementation path
- **Consider trade-offs** - document costs and risks
- **Quality over quantity** - 3 good suggestions beat 10 random ones
- **Categorize correctly** - use the 3 categories above
