# Skill: laddermoon-criticize

Version: 0.1.0

## Description

Analyze META for clarity and completeness, then file Questions for areas that need clarification.

This skill examines META.md against the organizational principles (Intuitive → Macro → Micro) and identifies:
- Missing information that prevents understanding the project
- Unclear or ambiguous descriptions
- Incomplete structural coverage
- Gaps between what META claims and what can be verified

**This skill ONLY creates Questions. It does NOT answer them.**

## Principles

### META Clarity Standards

A well-formed META should allow someone to:
1. **Intuitive**: Understand what this project is in one paragraph
2. **Macro**: See the complete high-level structure and architecture
3. **Micro**: Know where to find important details

### Question Categories

| Category | Description | Example |
|----------|-------------|---------|
| `missing-overview` | No intuitive description of the project | "What is this project for?" |
| `incomplete-structure` | Missing macro-level information | "What are the main components?" |
| `unclear-detail` | Ambiguous or vague descriptions | "What does 'flexible' mean here?" |
| `unverifiable-claim` | Claims that cannot be verified | "Where is the 'modular design' implemented?" |
| `missing-source` | Information without traceability | "Where did this decision come from?" |

## Input

- META.md content
- Project repository for verification

## Steps

1. **Read META.md**

   ```bash
   # In worktree on laddermoon-meta branch
   cat META.md
   ```

2. **Check Intuitive Level**

   Verify META has:
   - Clear project description (what it is, what it does)
   - Target audience or use case
   - Key value proposition

   If missing or unclear → File Question with category `missing-overview`

3. **Check Macro Level**

   Verify META has:
   - Architecture overview
   - Main components and their relationships
   - Technology stack with purposes
   - Current implementation status

   If incomplete → File Question with category `incomplete-structure`

4. **Check Micro Level**

   Verify META:
   - Points to where details can be found
   - Has information index for key topics
   - Documents important decisions with rationale

   If missing → File Question with category `unclear-detail`

5. **Verify Claims**

   For each claim in META:
   - Can it be verified in the codebase?
   - Is the source cited?

   If unverifiable → File Question with category `unverifiable-claim`
   If no source → File Question with category `missing-source`

6. **Create Question files**

   For each identified gap, create a Question file in `Questions/` directory.

   **Question file format**:
   ```markdown
   # Question: <clear title>

   **ID**: question-NNN
   **Category**: <missing-overview | incomplete-structure | unclear-detail | unverifiable-claim | missing-source>
   **Status**: Open
   **Created**: YYYY-MM-DD
   **Source**: Criticize

   ## Context

   <What triggered this question - quote from META if applicable>

   ## Question

   <The specific question that needs answering>

   ## Why This Matters

   <Why answering this question improves META clarity>

   ## Suggested Resolution

   - [ ] Option A: <description>
   - [ ] Option B: <description>
   - [ ] Provide information directly
   ```

7. **Commit changes**

   ```bash
   # Create worktree in project directory
   tmpdir=$(mktemp -d .lm-tmp-XXXXXX)
   git worktree add "$tmpdir" laddermoon-meta

   cd "$tmpdir"
   # Create Question files
   mkdir -p Questions
   # Write question files...

   git add Questions/
   git commit -m "Criticize: filed N questions for META clarity"

   cd ..
   git worktree remove "$tmpdir"
   ```

## Output

Summary of filed Questions:

```
## Criticize Summary

### Questions Filed: N

| ID | Category | Question |
|----|----------|----------|
| question-001 | missing-overview | What is the primary purpose? |
| question-002 | incomplete-structure | What are the main modules? |

### META Clarity Score

- Intuitive: ✓/✗
- Macro: ✓/✗  
- Micro: ✓/✗

### Next Steps

1. Run `lm verify question` to review and approve questions
2. Run `lm clarify` to resolve approved questions
```

## Guardrails

1. **Only file Questions** - Never modify META.md directly
2. **Be specific** - Each question should have a clear, answerable form
3. **Cite context** - Always quote the relevant part of META that triggered the question
4. **No fabrication** - Only question what is actually missing or unclear
5. **Actionable** - Each question should lead to a concrete improvement
