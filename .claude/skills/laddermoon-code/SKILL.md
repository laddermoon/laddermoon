# Skill: laddermoon-code

Version: 0.1.0

## Description

Write code to complete a development Task. This skill takes a Task as input and implements the required changes in a feature branch.

**Workflow**:
1. Create a feature branch from current HEAD
2. Implement the changes
3. Commit with clear messages
4. Output is ready for Review

## Input

- Task file path (e.g., `Tasks/task-001-implement-feature.md`)
- Or direct task description

## Steps

1. **Read the Task**

   If task file provided:
   ```bash
   cat Tasks/task-NNN-*.md
   ```

   Extract:
   - Task ID
   - Requirements
   - Acceptance criteria
   - Related Issue/Proposal (if any)

2. **Create feature branch**

   ```bash
   # Create branch from current HEAD
   git checkout -b lm-task-NNN
   ```

   Branch naming convention: `lm-task-<task-id>` or `lm-<issue/proposal-id>`

3. **Analyze requirements**

   Before coding:
   - Understand the scope
   - Identify files to modify
   - Plan the implementation approach

4. **Implement changes**

   Write the code following:
   - Project's existing code style
   - Best practices for the language/framework
   - Minimal changes to achieve the goal

   Use appropriate tools:
   ```bash
   # Read existing code
   cat path/to/file.go
   
   # Make edits
   # ... use edit tools ...
   
   # Run tests if available
   go test ./...
   ```

5. **Commit changes**

   ```bash
   git add <changed files>
   git commit -m "Task-NNN: <brief description>
   
   - <change 1>
   - <change 2>
   
   Related: <issue/proposal if any>"
   ```

   Commit message format:
   - First line: `Task-NNN: <brief summary>` (under 50 chars)
   - Body: bullet points of changes
   - Footer: related issues/proposals

6. **Verify implementation**

   ```bash
   # Build check
   go build ./...
   
   # Test check
   go test ./...
   
   # Lint check (if available)
   golangci-lint run
   ```

7. **Prepare for review**

   Output summary of changes:
   ```bash
   git diff main...lm-task-NNN --stat
   git log main..lm-task-NNN --oneline
   ```

## Output

```
## Code Summary

**Task**: task-NNN
**Branch**: lm-task-NNN
**Status**: Ready for Review

### Changes Made

| File | Change |
|------|--------|
| path/to/file.go | Added function X |
| path/to/other.go | Modified Y |

### Commits

- `abc1234` Task-NNN: Implement feature X
- `def5678` Task-NNN: Add tests for X

### Verification

- [ ] Build: ✓ Passed
- [ ] Tests: ✓ Passed (or N/A)
- [ ] Lint: ✓ Passed (or N/A)

### Next Steps

Run `lm review` to review this implementation.
```

## Guardrails

1. **Stay in scope** - Only implement what the Task requires
2. **Use feature branch** - Never commit directly to main
3. **Clear commits** - Each commit should be atomic and well-described
4. **Verify before done** - Always run available verification (build, test)
5. **No unrelated changes** - Don't fix other issues while implementing a task
6. **Follow style** - Match the project's existing code style
