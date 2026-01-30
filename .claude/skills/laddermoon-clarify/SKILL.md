# Skill: laddermoon-clarify

Version: 0.1.0

## Description

Resolve a Question by analyzing the codebase first, then asking the user if needed.

This skill takes a Question file as input and attempts to find the answer through:
1. Analyzing the codebase and existing documentation
2. If analysis is insufficient, asking the user for clarification

**Key principle**: Try to find the answer in the code first, only ask the user when necessary.

## Input

- Question file path (e.g., `Questions/question-001-overview.md`)
- Access to the project repository

## Steps

1. **Read the Question file**

   ```bash
   # In worktree on laddermoon-meta branch
   cat Questions/question-NNN-*.md
   ```

   Extract:
   - Question ID
   - Category
   - The actual question
   - Context

2. **Analyze codebase for answers**

   Based on the question category, search the appropriate sources:

   | Category | Where to Look |
   |----------|---------------|
   | `missing-overview` | README, package.json, go.mod, main entry files |
   | `incomplete-structure` | Directory structure, imports, main modules |
   | `unclear-detail` | Code comments, related source files |
   | `unverifiable-claim` | Actual implementation files |
   | `missing-source` | Git history, commit messages |

   Use tools like:
   ```bash
   # Find relevant files
   find . -name "*.go" -o -name "*.md" | head -20
   
   # Search for keywords
   grep -r "keyword" --include="*.go"
   
   # Check git history
   git log --oneline -10
   ```

3. **Evaluate findings**

   - **If answer found in code**: Proceed to update META
   - **If answer NOT found**: Ask user for clarification

4. **If asking user**

   Present the question clearly:
   ```
   ## Question Needs Your Input
   
   **Question**: <the question>
   
   **Context**: <why this matters>
   
   **What I found**: <any partial information from code analysis>
   
   **Please provide**: <specific information needed>
   ```

   Wait for user response.

5. **Update META with the answer**

   ```bash
   # Create worktree in project directory
   tmpdir=$(mktemp -d .lm-tmp-XXXXXX)
   git worktree add "$tmpdir" laddermoon-meta

   cd "$tmpdir"
   
   # Update META.md with the clarified information
   # Add proper source citation: [Source: Clarify from code] or [Source: User clarification]
   
   # Update Question status to Resolved
   # In Questions/question-NNN-*.md, change Status: Open -> Resolved
   # Add Resolution section
   
   git add META.md Questions/
   git commit -m "Clarify: resolved question-NNN - <brief summary>"

   cd ..
   git worktree remove "$tmpdir"
   ```

6. **Update Question file with resolution**

   Add to the Question file:
   ```markdown
   ## Resolution

   **Resolved**: YYYY-MM-DD
   **Method**: <Code analysis | User clarification>
   **Answer**: <the answer>
   **META Updated**: Yes/No
   ```

## Output

```
## Clarify Summary

**Question**: question-NNN
**Status**: Resolved

### Resolution

**Method**: <Code analysis | User clarification>
**Answer**: <brief summary of the answer>

### META Updates

- Updated section: <section name>
- Added information: <what was added>
- Source citation: [Source: <method>]

### Next Steps

Run `lm criticize` to check if META now has better clarity.
```

## Guardrails

1. **Code first, ask second** - Always try to find the answer in the codebase before asking the user
2. **Cite sources** - Always indicate whether the answer came from code analysis or user input
3. **Update META** - The answer should be integrated into META.md, not just recorded
4. **Mark resolved** - Update the Question file status to Resolved
5. **Be specific** - When asking the user, be very specific about what information is needed
6. **No fabrication** - Never invent answers; if unsure, ask the user
