---
name: laddermoon-audit
description: Detect potential issues in the project. Use when the user wants to find problems, bugs, risks, or architectural concerns.
license: MIT
compatibility: Requires LadderMoon initialized and synced (lm init, lm sync)
metadata:
  author: laddermoon
  version: "0.1.0"
  role: "Issuer"
---

Analyze the project and identify potential issues.

**You are the Issuer role.** Your job is to critically examine the project through the lens of META and identify genuine problems that need attention.

---

## Input

This skill is invoked to perform a project audit. The user may specify:
- A focus area (e.g., "audit security", "audit performance")
- Or request a general audit

---

## Steps

1. **Verify sync status**

   ```bash
   # Check if synced
   current=$(git rev-parse HEAD)
   synced=$(git show laddermoon-meta:.sync_state 2>/dev/null || echo "")
   ```

   If not synced, warn the user to run `lm sync` first.

2. **Read META.md for context**

   ```bash
   git show laddermoon-meta:META.md
   ```

   Understand:
   - Project architecture
   - Technical stack
   - Design decisions
   - Current state

3. **Explore the codebase**

   Based on META context, examine:
   - Code structure and organization
   - Dependencies and versions
   - Configuration and security
   - Error handling patterns
   - Test coverage
   - Documentation state

4. **Identify issues**

   Look for problems in these categories:

   | Category | Examples |
   |----------|----------|
   | **Security** | Hardcoded secrets, injection vulnerabilities, auth issues |
   | **Architecture** | Tight coupling, circular dependencies, god classes |
   | **Performance** | N+1 queries, memory leaks, inefficient algorithms |
   | **Reliability** | Missing error handling, no retries, race conditions |
   | **Maintainability** | Dead code, inconsistent patterns, missing tests |
   | **Documentation** | Outdated docs, missing README, unclear APIs |

5. **Prioritize and document**

   For each issue found, assess:
   - **Severity**: Critical / High / Medium / Low
   - **Effort**: How hard to fix
   - **Impact**: What happens if not fixed

6. **Create Issue files**

   For significant issues, create files in Issues/:

   ```bash
   tmpdir=$(mktemp -d)
   git worktree add "$tmpdir" laddermoon-meta
   
   # Create issue file
   cat > "$tmpdir/Issues/issue-<id>.md" << 'EOF'
   # Issue: <Title>
   
   **Severity**: High
   **Category**: Security
   **Detected**: <date>
   
   ## Description
   <detailed description>
   
   ## Location
   <file paths and line numbers>
   
   ## Impact
   <what could go wrong>
   
   ## Recommendation
   <how to fix>
   EOF
   
   cd "$tmpdir"
   git add Issues/
   git commit -m "Audit: Found <n> issues"
   
   cd -
   git worktree remove "$tmpdir"
   ```

---

## Output Format

```
## Audit Results

**Scope**: <general or specific area>
**Issues found**: <count>

### Critical Issues
(Issues that need immediate attention)

### High Priority
(Significant problems that should be addressed soon)

### Medium Priority
(Issues to address when convenient)

### Low Priority
(Minor improvements)

---

**Issue files created**: <list>

**Next steps**:
- Run `lm solve Issues/<file>` to address specific issues
- Run `lm propose` for improvement suggestions
```

---

## Thinking Guidelines

When auditing, think like:
- A security researcher looking for vulnerabilities
- A performance engineer looking for bottlenecks
- A maintainer who will own this code for years
- A new developer trying to understand the codebase

Ask yourself:
- "What could go wrong here?"
- "Will this scale?"
- "Is this secure?"
- "Can someone understand this in 6 months?"

---

## Guardrails

- **Be specific** - Point to exact files and lines, not vague concerns
- **Be actionable** - Every issue should have a clear path to resolution
- **Be honest** - Don't manufacture issues, only report real problems
- **Be prioritized** - Not all issues are equal, rank them
- **Don't nitpick** - Focus on meaningful issues, not style preferences
- **Consider context** - An MVP has different standards than production code
- **Verify sync** - Don't audit stale META, ensure it's current
