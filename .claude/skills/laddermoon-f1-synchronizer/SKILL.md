---
name: laddermoon-f1-synchronizer
description: Synchronize repository changes to metadata. Use when code has changed since last sync to keep meta.md up to date.
license: MIT
compatibility: Requires initialized metadata (.laddermoon/meta.md)
metadata:
  author: laddermoon
  version: "1.0"
---

Synchronize repository changes to metadata.

**Input**: None required. The skill compares the current commit to the last synced commit.

**Steps**

1. **Check metadata exists**

   ```bash
   cat .laddermoon/meta-commit-id
   ```

   If `.laddermoon/meta-commit-id` doesn't exist, inform user to run `laddermoon-f1-initializer` first.

2. **Get sync state**

   ```bash
   # Last synced commit
   last_commit=$(cat .laddermoon/meta-commit-id)
   
   # Current commit
   current_commit=$(git rev-parse HEAD)
   
   echo "Last sync: $last_commit"
   echo "Current: $current_commit"
   ```

   If `last_commit == current_commit`, inform user metadata is up to date and STOP.

3. **Analyze changes since last sync**

   ```bash
   # Commit log
   git log --oneline ${last_commit}..${current_commit}
   
   # Changed files
   git diff --stat ${last_commit} ${current_commit}
   
   # Detailed changes for key files
   git diff ${last_commit} ${current_commit} -- README.md package.json go.mod
   ```

4. **Read current metadata**

   ```bash
   cat .laddermoon/meta.md
   ```

5. **Determine what needs updating**

   Analyze changes and categorize:

   | Change Type | META Section to Update |
   |-------------|------------------------|
   | README changed | Section 1 (Overview) |
   | Dependencies changed | Section 2.1 (Tech Stack) |
   | New directories/files | Section 2.2/2.3 (Architecture/Components) |
   | Files moved/renamed | Section 3 (Information Index) |
   | New docs added | Section 3 (Information Index) |

   **Decision rule**: Only update if the change affects the macro structure. Ignore:
   - Code implementation details
   - Minor file edits that don't change structure
   - Test file additions (unless test directory structure changed)

6. **Update meta.md if needed**

   For each section that needs updating:
   - Read the source file
   - Extract relevant information
   - Update the section with `[Source: <file>]` citation
   - Update the "Last updated" timestamp and commit reference

   **IMPORTANT**: Do NOT modify information that hasn't changed. Only add/update what's affected by the commits.

7. **Update meta-commit-id**

   ```bash
   git rev-parse HEAD > .laddermoon/meta-commit-id
   ```

**Output**

```
## Synchronization Complete

**Commits analyzed**: <count> (<last_commit>..<current_commit>)

### Changes Detected

- <change 1> [Source: <file>]
- <change 2> [Source: <file>]

### META Updates

| Section | Change | Source |
|---------|--------|--------|
| <section> | <what changed> | <file> |

### Sync State

- **Previous commit**: <last_commit>
- **Current commit**: <current_commit>
- **meta-commit-id**: Updated ✓
```

**Output When No Update Needed**

```
## Synchronization Complete

**Status**: Already up to date
**Current commit**: <commit_id>

No changes affect metadata structure. No updates made.
```

**Organization Principles (MUST FOLLOW)**

| Principle | Description |
|-----------|-------------|
| **Current state only** | Update to reflect current state, not change history |
| **Intuitive → Macro → Micro** | Maintain the three-level structure |
| **All sources cited** | Every update must have `[Source: path]` |
| **No fabrication** | Only update with verifiable information |
| **No inference** | Don't guess what changes mean for project direction |
| **Minimal updates** | Only touch sections affected by changes |

**Guardrails**

- Never run without existing metadata (use initializer first)
- Always compare commits before making changes
- Only update sections affected by actual changes
- Always cite sources for every modification
- Never infer user intent from code changes
- Update meta-commit-id only after successful sync
- If a change is unclear, leave meta.md unchanged (inspector will catch gaps)
