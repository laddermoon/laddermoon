---
name: laddermoon-improve
description: Self-improvement of the LadderMoon META system. Analyze META quality and suggest improvements to how project information is organized.
license: MIT
compatibility: Requires LadderMoon initialized (lm init)
metadata:
  author: laddermoon
  version: "0.1.0"
  role: "Self-Improver"
---

Analyze and improve the quality of the META system itself.

**You are the Self-Improver role.** Your job is to analyze how well the META system is working and suggest improvements to META organization, information quality, and LadderMoon workflow effectiveness.

---

## What Improve Does

1. **Analyzes META.md quality** - Is information well-organized?
2. **Checks workflow effectiveness** - Are Issues/Suggestions/Questions being resolved?
3. **Identifies META gaps** - What's missing or could be better?
4. **Suggests META improvements** - How to organize information better

---

## Improvement Areas

### Area 1: META.md Structure
Is META.md well-organized following the principles?
- **Intuitive level**: Is it clear what the project is about?
- **Macro level**: Is the information structure logical?
- **Micro level**: Are there pointers to details, not duplicated content?

### Area 2: Traceability
Is information properly cited?
- Do all statements have `[Feed #N]` or `[Source: path]`?
- Are conflict markers properly used and resolved?
- Can you trace any statement to its origin?

### Area 3: Intent vs Reality Balance
Is there a good balance between intent and reality documentation?
- Are user goals clearly documented?
- Is the current state accurately reflected?
- Are the gaps between intent and reality visible?

### Area 4: Workflow Health
Are Issues/Suggestions/Questions being managed?
- How many open Items exist?
- Are items being resolved?
- Are there stale items?

### Area 5: Information Freshness
Is META.md up to date?
- When was last sync?
- Have there been code changes since?
- Is the documented reality still accurate?

---

## Core Principles (MUST FOLLOW)

### Principle 1: Improve META, Not Project
- Focus on META system quality
- Don't create Issues/Suggestions for code
- Improvements are about how information is organized

### Principle 2: Follow META Principles
- Intuitive → Macro → Micro
- Intent vs Reality
- Source and Traceability

### Principle 3: Actionable Recommendations
- Be specific about what to improve
- Explain how to improve it
- Prioritize by impact

---

## Input

This skill is invoked via `lm improve`.

---

## Steps

1. **Get current branch and read META**

   ```bash
   branch=$(git rev-parse --abbrev-ref HEAD)
   branch_dir=$(echo "$branch" | tr '/' '_')
   
   git show laddermoon-meta:${branch_dir}/META.md
   ```

2. **Analyze structure quality**

   Check:
   - Is there a clear project overview?
   - Are sections logically organized?
   - Is information at the right level of detail?

3. **Check traceability**

   Scan for:
   - Statements without source citations
   - Unresolved conflict markers
   - Broken references

4. **Examine workflow items**

   ```bash
   git ls-tree laddermoon-meta:${branch_dir}/Issues/
   git ls-tree laddermoon-meta:${branch_dir}/Suggestions/
   git ls-tree laddermoon-meta:${branch_dir}/Questions/
   ```

   Assess:
   - Number of open items
   - Age of items
   - Resolution rate

5. **Check freshness**

   ```bash
   git show laddermoon-meta:${branch_dir}/.sync_state
   git rev-parse HEAD
   ```

   Compare last sync with current state.

6. **Generate improvement report**

   Create a report with findings and recommendations.

---

## Output Format

```
## META System Health Report

**Branch**: <branch>
**Last Sync**: <date/commit>
**Current Commit**: <commit>

### Structure Quality

| Aspect | Status | Notes |
|--------|--------|-------|
| Intuitive Overview | ✓ Good / ⚠ Needs Work | <notes> |
| Macro Structure | ✓ Good / ⚠ Needs Work | <notes> |
| Micro Detail Level | ✓ Good / ⚠ Needs Work | <notes> |

### Traceability

| Check | Result |
|-------|--------|
| Statements with sources | X of Y (Z%) |
| Unresolved conflicts | N |
| Broken references | N |

### Workflow Health

| Item Type | Open | Resolved | Stale (>7d) |
|-----------|------|----------|-------------|
| Issues | N | M | K |
| Suggestions | N | M | K |
| Questions | N | M | K |

### Freshness

- Last sync: <commit> (<N commits behind>)
- Recommendation: <run sync / up to date>

---

## Recommendations

### High Priority
1. <recommendation> - <why it matters>

### Medium Priority
1. <recommendation> - <why it matters>

### Low Priority
1. <recommendation> - <why it matters>

---

**Next steps**:
- Address high priority items first
- Run `lm sync` if behind
- Answer open Questions via `lm answer`
```

---

## Example Recommendations

| Issue | Recommendation |
|-------|----------------|
| No project overview | Add intuitive description at top of META |
| Missing source citations | Add `[Feed #N]` or `[Source: path]` to uncited statements |
| Stale Questions | Answer or close Questions older than 7 days |
| Out of date | Run `lm sync` to catch up with code changes |
| Too much detail | Move detailed info to files, add pointers in META |
| Duplicate info | Remove duplicates, keep single source of truth |

---

## Guardrails

- **Focus on META system** - not project code
- **Be constructive** - provide actionable recommendations
- **Prioritize by impact** - most important first
- **Don't create Issues/Suggestions** - this is a report, not action
- **Be specific** - point to exact problems and solutions
- **Follow META principles** - recommendations should align with principles
