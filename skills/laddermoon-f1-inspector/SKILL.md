---
name: laddermoon-f1-inspector
description: Inspect metadata clarity and create Questions for gaps. Use to check if meta.md is clear enough for users and AI to understand the project.
license: MIT
compatibility: Requires initialized metadata (.laddermoon/meta.md)
metadata:
  author: laddermoon
  version: "1.0"
---

Inspect metadata clarity and create Questions for identified gaps.

**Input**: Question filename to write to (e.g., `question-3.md`). If not provided, auto-generate based on `last-question-id`.

**Steps**

1. **Read current metadata**

   ```bash
   cat .laddermoon/meta.md
   ```

   If metadata doesn't exist, inform user to run `laddermoon-f1-initializer` first.

2. **Check against organization principles**

   **Principle 0: Current state only**
   - Does metadata describe current state, or does it include historical/future info?
   - Are there outdated references?

   **Principle 1-4: Intuitive → Macro → Micro structure**

   | Level | Check | Question if failing |
   |-------|-------|---------------------|
   | Intuitive | Is there a clear one-paragraph overview? | "What is this project's primary purpose?" |
   | Intuitive | Can a newcomer understand what this is? | "What problem does this project solve?" |
   | Macro | **CRITICAL: Are all major directories listed with purposes?** | "What are the main directories and what does each do?" |
   | Macro | Is the directory structure clear? | "What is in each directory?" |
   | Macro | Is the architecture described? | "How is the code logically organized?" |
   | Macro | Are main components/modules listed? | "What are the key components?" |
   | Macro | Is tech stack documented? | "What technologies are used?" |
   | Macro | Are entry points documented? | "How do you run/use this project?" |
   | Macro | Are usage patterns documented? | "How do you build/test/deploy?" |
   | Macro | Are standards/conventions documented? | "What coding standards are followed?" |
   | Micro | Is there an information index? | "Where can I find documentation about X?" |
   | Micro | Are key file locations documented? | "Where is the configuration located?" |

   **Principle 5: Sources cited**
   - Does every claim have a `[Source: path]` citation?
   - If not, which claims lack sources?

   **Principle 6: No fabrication**
   - Spot-check claims against the actual repo
   - Are there claims that cannot be verified?

3. **Evaluate from curious newcomer perspective**

   Someone discovering this project would want to know:
   - **What is this?** (purpose, scope)
   - **What are the main directories?** (critical - must be documented)
   - **What does each directory do?** (purpose of each major directory)
   - **How is the code organized?** (architecture, modules, components)
   - **What are the key parts?** (entry points, core features)
   - **How do I use it?** (run, build, test, deploy)
   - **What technologies does it use?** (languages, frameworks, dependencies)
   - **Where do I find things?** (docs, configs, tests)
   - **What standards does it follow?** (coding conventions, patterns)

   For each question that meta.md doesn't clearly answer, create a Question.

4. **Get next question ID**

   ```bash
   last_id=$(cat .laddermoon/questions/last-question-id)
   next_id=$((last_id + 1))
   ```

5. **Create Question files for identified gaps**

   For each gap, create a question file:

   **Filename**: `.laddermoon/questions/question-<N>.md`

   ```markdown
   # Question <N>

   ## Status
   open

   ## Category
   <organization-principle | project-manager-concern>

   ## Question
   <The specific question that needs answering>

   ## Context
   <What in meta.md triggered this question, or what is missing>

   ## Why It Matters
   <How answering this improves metadata clarity>
   ```

6. **Update question tracking**

   ```bash
   # Update last-question-id
   echo "<highest_new_id>" > .laddermoon/questions/last-question-id
   
   # Append to question-meta.jsonl
   echo '{"id": <N>, "status": "open"}' >> .laddermoon/questions/question-meta.jsonl
   ```

**Output**

```
## Inspection Complete

### Clarity Assessment

| Level | Status | Notes |
|-------|--------|-------|
| Intuitive | ✓/✗ | <brief note> |
| Macro | ✓/✗ | <brief note> |
| Micro | ✓/✗ | <brief note> |
| Sources | ✓/✗ | <brief note> |

### Questions Created: <N>

| ID | Category | Question |
|----|----------|----------|
| question-<N> | <category> | <question summary> |

### Files Written

- .laddermoon/questions/question-<N>.md
- .laddermoon/questions/last-question-id (updated)
- .laddermoon/questions/question-meta.jsonl (appended)

### Next Steps

1. Run `laddermoon-f1-clarifier` on each question to find answers
2. Re-run inspector after clarification to check for remaining gaps
```

**Output When No Questions Needed**

```
## Inspection Complete

### Clarity Assessment

| Level | Status |
|-------|--------|
| Intuitive | ✓ |
| Macro | ✓ |
| Micro | ✓ |
| Sources | ✓ |

**Metadata is clear.** No questions created.

The metadata satisfies all organization principles and answers key project manager concerns.
```

**Question Categories**

| Category | When to use |
|----------|-------------|
| `missing-overview` | Intuitive level is unclear or missing |
| `incomplete-structure` | Macro level has gaps |
| `missing-index` | Micro level (information index) incomplete |
| `missing-source` | Claim lacks `[Source: path]` citation |
| `unverifiable` | Claim cannot be verified in repo |
| `outdated` | Information appears stale |
| `project-concern` | Project manager would need this info |

**Guardrails**

- **Only create Questions** - Never modify meta.md directly
- **Be specific** - Each question must be clearly answerable
- **Cite context** - Reference what in meta.md (or its absence) triggered the question
- **No fabrication** - Only question actual gaps, not imagined ones
- **Prioritize** - Focus on questions that most improve clarity
- **Don't duplicate** - Check existing questions before creating new ones
