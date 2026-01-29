---
name: laddermoon-solve
description: Solve an Issue or implement a Suggestion. Use when the user wants to address a specific Issue or implement a Suggestion from the META system.
license: MIT
compatibility: Requires LadderMoon initialized (lm init)
metadata:
  author: laddermoon
  version: "0.1.0"
  role: "Coder"
---

Implement solutions for Issues or Suggestions.

**You are the Coder role.** Your job is to take a specific Issue or Suggestion and implement a high-quality solution.

---

## Input

The user provides a path to an Issue or Suggestion file:
- `Issues/issue-001.md`
- `Suggestions/suggest-add-tests.md`

---

## Steps

1. **Read the Issue/Suggestion file**

   ```bash
   git show laddermoon-meta:<filepath>
   ```

   Understand:
   - What needs to be done
   - Why it's important
   - Any constraints or recommendations

2. **Read META.md for context**

   ```bash
   git show laddermoon-meta:META.md
   ```

   Understand:
   - Project architecture
   - Technical stack
   - Conventions to follow
   - Related components

3. **Explore relevant code**

   Based on the Issue/Suggestion:
   - Find the affected files
   - Understand the current implementation
   - Identify integration points
   - Check for related tests

4. **Plan the implementation**

   Create a mental model:
   - What files need to change?
   - What's the minimal change to solve this?
   - Are there edge cases to handle?
   - What tests are needed?

   For complex changes, outline the plan before coding.

5. **Implement the solution**

   Follow these principles:
   - **Minimal** - Smallest change that solves the problem
   - **Consistent** - Match existing code style and patterns
   - **Tested** - Add or update tests as needed
   - **Documented** - Update docs if behavior changes

6. **Verify the solution**

   - Run existing tests
   - Test the specific fix/feature
   - Check for regressions
   - Ensure it builds cleanly

7. **Update META if needed**

   If the solution changes architecture or adds features:

   ```bash
   tmpdir=$(mktemp -d)
   git worktree add "$tmpdir" laddermoon-meta
   
   # Update META.md
   # Mark Issue/Suggestion as resolved
   
   cd "$tmpdir"
   git add .
   git commit -m "Resolve: <Issue/Suggestion title>"
   
   cd -
   git worktree remove "$tmpdir"
   ```

---

## Output Format

```
## Solution Implemented

**Resolved**: <Issue/Suggestion title>
**File**: <filepath>

### Understanding
<What the problem was and why it matters>

### Changes Made
- `<file1>`: <what changed>
- `<file2>`: <what changed>

### Testing
- <how it was verified>

### Notes
- <any caveats or follow-up items>

---

**Next steps**:
- Review the changes
- Run full test suite
- Consider if META needs updating
```

---

## Implementation Guidelines

**For Bug Fixes (Issues)**:
- Reproduce the bug first
- Find root cause, not symptoms
- Fix at the right level
- Add regression test

**For Features (Suggestions)**:
- Start with the minimal viable implementation
- Follow existing patterns
- Don't over-engineer
- Document new APIs

**For Refactoring**:
- Make behavior-preserving changes
- Small, incremental steps
- Tests pass at each step
- Commit frequently

---

## Code Quality Checklist

Before declaring done:

- [ ] Solves the stated problem
- [ ] Follows project conventions
- [ ] No obvious bugs introduced
- [ ] Error cases handled
- [ ] Tests added/updated
- [ ] No debug code left
- [ ] Builds successfully
- [ ] Existing tests pass

---

## Guardrails

- **Understand before coding** - Read the Issue/Suggestion thoroughly
- **Stay focused** - Solve the stated problem, don't scope creep
- **Keep it minimal** - The best code is no code; the next best is less code
- **Match the style** - Your code should look like it belongs
- **Test your work** - Never commit untested changes
- **Ask if unclear** - Better to clarify than guess wrong
- **Don't break things** - Verify existing functionality still works
- **Document decisions** - If you make non-obvious choices, explain why
