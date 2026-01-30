# Skill: laddermoon-apply

Version: 0.1.0

## Description

Merge a feature branch into main. If merge conflicts occur, attempt to resolve them automatically or with Claude's help.

**Workflow**:
1. Attempt to merge the feature branch
2. If successful, done
3. If conflicts, analyze and resolve them
4. If resolution fails, request human intervention

## Input

- Feature branch name (e.g., `lm-task-001`)
- Target branch (default: `main`)

## Steps

1. **Check branch status**

   ```bash
   # Verify feature branch exists
   git branch --list lm-task-*
   
   # Check current branch
   git branch --show-current
   
   # Ensure working directory is clean
   git status --porcelain
   ```

2. **Switch to target branch**

   ```bash
   git checkout main
   git pull origin main  # If remote exists
   ```

3. **Attempt merge**

   ```bash
   git merge lm-task-NNN --no-ff -m "Merge lm-task-NNN: <description>"
   ```

4. **Handle merge result**

   **If successful**:
   - Report success
   - Optionally delete feature branch

   **If conflicts**:
   - List conflicting files
   - Attempt automatic resolution

5. **Resolve conflicts (if any)**

   ```bash
   # List conflicts
   git diff --name-only --diff-filter=U
   
   # For each conflicting file:
   # 1. Read the file with conflict markers
   # 2. Understand both versions
   # 3. Determine correct resolution
   # 4. Edit file to resolve
   # 5. Stage the resolved file
   ```

   Resolution strategies:
   - **Code conflicts**: Analyze intent of both changes, merge logically
   - **Config conflicts**: Usually take the feature branch version
   - **Documentation conflicts**: Combine both additions

6. **Complete merge after resolution**

   ```bash
   git add <resolved files>
   git commit -m "Merge lm-task-NNN: <description>

   Resolved conflicts in:
   - file1.go
   - file2.go"
   ```

7. **Cleanup**

   ```bash
   # Delete feature branch after successful merge
   git branch -d lm-task-NNN
   ```

## Output

### Success case:
```
## Apply Summary

**Branch**: lm-task-NNN → main
**Status**: ✓ Merged successfully

### Merge Details

- Commits merged: N
- Files changed: M
- Conflicts: None

### Cleanup

- Feature branch deleted: Yes/No

### Next Steps

Run `lm sync` to update META with the new changes.
```

### Conflict resolved case:
```
## Apply Summary

**Branch**: lm-task-NNN → main
**Status**: ✓ Merged with conflict resolution

### Conflicts Resolved

| File | Resolution |
|------|------------|
| path/to/file.go | Combined both changes |
| path/to/config.yaml | Used feature branch version |

### Next Steps

1. Verify the resolution is correct
2. Run `lm sync` to update META
```

### Conflict unresolved case:
```
## Apply Summary

**Branch**: lm-task-NNN → main
**Status**: ✗ Needs manual intervention

### Unresolved Conflicts

| File | Issue |
|------|-------|
| path/to/complex.go | Complex logic conflict, needs human decision |

### How to resolve manually

```bash
# 1. View conflicts
git diff

# 2. Edit files to resolve
# 3. Stage resolved files
git add <files>

# 4. Complete merge
git commit
```

### Or abort the merge

```bash
git merge --abort
```
```

## Guardrails

1. **Clean state required** - Working directory must be clean before merge
2. **No force push** - Never use `--force` options
3. **Preserve history** - Use `--no-ff` to create merge commit
4. **Verify resolution** - After resolving conflicts, ensure code compiles/tests pass
5. **Clear communication** - If can't resolve, clearly explain what needs human attention
6. **Backup branch** - Don't delete feature branch until merge is confirmed successful
