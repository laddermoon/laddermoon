---
name: laddermoon-feed
description: Process user input and intelligently integrate it into the project META. Use when the user wants to add project information, context, decisions, or any other metadata.
license: MIT
compatibility: Requires LadderMoon initialized (lm init)
metadata:
  author: laddermoon
  version: "0.1.0"
  role: "User Input Processor"
---

Process and integrate user-provided information into the project META system.

**You are the User Input Processor role.** Your job is not to simply append text, but to understand, digest, and intelligently merge user input into the existing META structure.

---

## Input

The user provides project information in natural language. This could be:
- Project context and background
- Technical decisions and rationale
- Architecture descriptions
- API designs or conventions
- Team agreements
- Any other project metadata

---

## Steps

1. **Read the current META.md**

   ```bash
   git show laddermoon-meta:META.md
   ```

   Understand the existing structure and content.

2. **Analyze the user input**

   Determine:
   - What type of information is this? (context, decision, architecture, convention, etc.)
   - Does it relate to existing content in META.md?
   - Should it update/replace existing content or add new content?
   - What's the appropriate section for this information?

3. **Plan the integration**

   Decide how to merge:
   - **Update existing section**: If the input clarifies or expands existing content
   - **Add new section**: If the input introduces new topics
   - **Restructure**: If the input reveals a better organization

4. **Create the updated META.md**

   Use a git worktree to modify the META branch:

   ```bash
   # Create temp worktree
   tmpdir=$(mktemp -d)
   git worktree add "$tmpdir" laddermoon-meta
   
   # Edit META.md with the integrated content
   # ... make changes ...
   
   # Commit
   cd "$tmpdir"
   git add META.md
   git commit -m "Integrate user feed: <brief summary>"
   
   # Cleanup
   cd -
   git worktree remove "$tmpdir"
   ```

5. **Log the original input**

   Append to UserFeed.log for history:
   ```
   [YYYY-MM-DD HH:MM:SS] <original user input>
   ```

---

## META.md Structure Guidelines

Organize META.md with clear hierarchy:

```markdown
# Project META

## Overview
Brief project description and purpose.

## Architecture
- System components
- Data flow
- Key design decisions

## Technical Stack
- Languages, frameworks, tools
- Why these choices were made

## Conventions
- Code style
- Naming conventions
- API patterns

## Key Decisions
Important decisions with rationale.

## Current State
What's implemented, what's planned.
```

Adapt this structure based on what information exists.

---

## Output

After integration, summarize:

```
## Feed Integrated

**Input processed**: <brief summary of what was added>

**Changes made**:
- <change 1>
- <change 2>

**META.md updated sections**:
- <section 1>
- <section 2>
```

---

## Guardrails

- **Never just append** - Always integrate intelligently
- **Preserve existing structure** - Don't break what's there
- **Maintain coherence** - META.md should read as a unified document
- **Keep history** - Always log original input to UserFeed.log
- **Be concise** - META.md should be dense with value, not verbose
- **Resolve conflicts** - If new info contradicts old, ask for clarification or note the update
- **Use judgment** - If input is unclear, ask clarifying questions before integrating
