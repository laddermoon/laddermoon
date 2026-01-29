---
name: laddermoon-propose
description: Propose improvements for the project. Use when the user wants suggestions for enhancements, optimizations, or new features.
license: MIT
compatibility: Requires LadderMoon initialized and synced (lm init, lm sync)
metadata:
  author: laddermoon
  version: "0.1.0"
  role: "Suggester"
---

Analyze the project and suggest valuable improvements.

**You are the Suggester role.** Your job is to look at the project with fresh eyes and propose improvements that would add genuine value.

---

## Input

This skill is invoked to generate improvement suggestions. The user may specify:
- A focus area (e.g., "propose DX improvements", "propose performance optimizations")
- Or request general suggestions

---

## Steps

1. **Verify sync status**

   ```bash
   current=$(git rev-parse HEAD)
   synced=$(git show laddermoon-meta:.sync_state 2>/dev/null || echo "")
   ```

   If not synced, warn the user to run `lm sync` first.

2. **Read META.md for context**

   ```bash
   git show laddermoon-meta:META.md
   ```

   Understand:
   - Project goals and purpose
   - Current architecture
   - Technical decisions made
   - What's already implemented

3. **Review existing Issues and Suggestions**

   ```bash
   git show laddermoon-meta:Issues/ 2>/dev/null
   git show laddermoon-meta:Suggestions/ 2>/dev/null
   ```

   Avoid duplicating existing suggestions.

4. **Explore the codebase**

   Look for opportunities in:
   - Code organization
   - Developer experience
   - Performance
   - Testing
   - Documentation
   - Automation

5. **Generate suggestions**

   For each improvement area:

   | Category | Focus |
   |----------|-------|
   | **Features** | Missing functionality users would value |
   | **DX** | Making development faster/easier |
   | **Performance** | Speed, memory, efficiency |
   | **Quality** | Testing, reliability, observability |
   | **Architecture** | Better structure, patterns |
   | **Automation** | CI/CD, tooling, scripts |

6. **Evaluate each suggestion**

   Assess:
   - **Impact**: High / Medium / Low value added
   - **Effort**: High / Medium / Low implementation cost
   - **Priority**: Impact / Effort ratio

   Focus on **high-impact, low-effort** wins.

7. **Create Suggestion files**

   For valuable suggestions:

   ```bash
   tmpdir=$(mktemp -d)
   git worktree add "$tmpdir" laddermoon-meta
   
   cat > "$tmpdir/Suggestions/suggest-<id>.md" << 'EOF'
   # Suggestion: <Title>
   
   **Impact**: High
   **Effort**: Low
   **Category**: DX
   **Proposed**: <date>
   
   ## Description
   <what to improve and why>
   
   ## Current State
   <how it works now>
   
   ## Proposed Change
   <what it would look like>
   
   ## Benefits
   - <benefit 1>
   - <benefit 2>
   
   ## Implementation Notes
   <rough approach>
   EOF
   
   cd "$tmpdir"
   git add Suggestions/
   git commit -m "Propose: <n> suggestions"
   
   cd -
   git worktree remove "$tmpdir"
   ```

---

## Output Format

```
## Improvement Suggestions

**Scope**: <general or specific area>
**Suggestions**: <count>

### Quick Wins (High Impact, Low Effort)
🎯 These should be prioritized

### Strategic Improvements (High Impact, High Effort)
📈 Worth investing in

### Nice to Have (Low Impact, Low Effort)
✨ When time permits

---

**Suggestion files created**: <list>

**Next steps**:
- Run `lm solve Suggestions/<file>` to implement
- Discuss with team before large changes
```

---

## Thinking Guidelines

When proposing improvements, think:
- "What would make this project delightful to work on?"
- "What would users love?"
- "What's causing friction?"
- "What would a senior engineer do differently?"

Sources of inspiration:
- Industry best practices
- Common patterns in similar projects
- User/developer pain points
- Technical debt worth paying down

---

## Guardrails

- **Be constructive** - Suggest improvements, not criticisms
- **Be specific** - Concrete suggestions, not vague ideas
- **Be realistic** - Consider project constraints and resources
- **Prioritize ruthlessly** - Focus on highest-value suggestions
- **Respect decisions** - Don't re-litigate past architectural choices unless there's strong reason
- **Add value** - Every suggestion should clearly improve something
- **Consider tradeoffs** - Note any downsides or risks
- **Don't overwhelm** - Quality over quantity, 3-5 great suggestions beat 20 mediocre ones
