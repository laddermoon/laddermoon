---
name: laddermoon-sync
description: Synchronize repository changes to META. Use when the codebase has changed and META needs to be updated to reflect the current state.
license: MIT
compatibility: Requires LadderMoon initialized (lm init)
metadata:
  author: laddermoon
  version: "0.1.0"
  role: "Repo Syncer"
---

Synchronize codebase changes into the META system.

**You are the Repo Syncer role.** Your job is to analyze code changes since the last sync and update META.md to reflect the current state of the project.

---

## Input

This skill is invoked when code has changed. It may be called:
- After commits to the main branch
- Periodically to catch up
- When the user explicitly requests a sync

---

## Steps

1. **Check sync state**

   ```bash
   git show laddermoon-meta:.sync_state 2>/dev/null || echo ""
   ```

   This contains the last synced commit ID.

2. **Get current commit**

   ```bash
   git rev-parse HEAD
   ```

3. **If already synced, report and exit**

   If `.sync_state` matches HEAD, nothing to do.

4. **Analyze changes since last sync**

   ```bash
   # Get commit log
   git log --oneline <last_sync>..HEAD
   
   # Get changed files
   git diff --stat <last_sync> HEAD
   
   # For significant changes, read the actual diffs
   git diff <last_sync> HEAD -- <important_files>
   ```

   If no previous sync, analyze the current state:
   ```bash
   git log --oneline -20
   git ls-files
   ```

5. **Read current META.md**

   ```bash
   git show laddermoon-meta:META.md
   ```

6. **Determine what META updates are needed**

   Based on the changes, consider:
   - New files/directories → update Architecture or Structure sections
   - Dependency changes → update Technical Stack
   - New features → update Current State
   - API changes → update relevant documentation
   - Refactoring → may not need META updates

7. **Update META.md if needed**

   Use git worktree:
   ```bash
   tmpdir=$(mktemp -d)
   git worktree add "$tmpdir" laddermoon-meta
   
   # Update META.md with relevant changes
   # Update .sync_state with current commit
   
   cd "$tmpdir"
   echo "<current_commit>" > .sync_state
   git add META.md .sync_state
   git commit -m "Sync: <brief summary of code changes>"
   
   cd -
   git worktree remove "$tmpdir"
   ```

8. **If no META updates needed, just update sync state**

   Still commit the new `.sync_state` to track the sync point.

---

## What to Track in META

Focus on **structural and architectural changes**:

| Change Type | META Action |
|-------------|-------------|
| New module/package | Update Architecture |
| New dependency | Update Technical Stack |
| New API endpoint | Update API section if exists |
| Database schema change | Note in Architecture |
| Configuration change | Update relevant section |
| Pure refactoring | Usually no META update |
| Bug fix | Usually no META update |
| New feature | Update Current State |

---

## Output

```
## Sync Complete

**Commits synced**: <count> (<from>..<to>)

**Code changes detected**:
- <change 1>
- <change 2>

**META updates**:
- <update 1> (or "No META updates needed")

**Sync state**: <new commit id>
```

---

## Guardrails

- **Don't over-document** - Only track significant structural changes
- **Keep META high-level** - Don't duplicate code documentation
- **Always update .sync_state** - Even if no META changes needed
- **Summarize, don't enumerate** - "Added 5 new API endpoints" not listing all 5
- **Focus on the "what" and "why"** - Not implementation details
- **Preserve manual META edits** - Don't overwrite user's careful documentation
