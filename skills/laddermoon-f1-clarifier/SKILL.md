---
name: laddermoon-f1-clarifier
description: Answer Inspector's Questions by searching the repo. If answer cannot be found, create an Issue. Use to resolve Questions and improve metadata.
license: MIT
compatibility: Requires initialized metadata and existing Questions
metadata:
  author: laddermoon
  version: "1.0"
---

Answer a Question by searching the repository. If the answer cannot be found, create an Issue.

**Input**: 
- Question filename (e.g., `question-3.md`)
- Issue filename to create if needed (e.g., `issue-1.md`). If not provided, auto-generate based on `last-issue-id`.

**Steps**

1. **Read the Question file**

   ```bash
   cat .laddermoon/questions/<question-filename>
   ```

   Extract:
   - Question ID
   - Category
   - The actual question
   - Context

2. **Search the repository for answers**

   Explore the repository to find the answer. Be thorough and curious:

   **For questions about project purpose/overview**:
   - Look for README files, package manifests, documentation
   - Check project root files that describe what this is

   **For questions about directory structure/purpose**:
   - List directory contents to see what's inside
   - Look for README files within directories
   - Examine file types and naming patterns
   - Check for package/module definitions
   - If still unclear, explore subdirectories

   **For questions about architecture/components**:
   - Examine directory organization
   - Look at import/dependency patterns
   - Check for architectural documentation
   - Review main entry points

   **For questions about usage/operations**:
   - Look for build files (Makefile, package.json scripts, etc.)
   - Check for documentation on running/building
   - Examine CI/CD configurations

   **For questions about standards/conventions**:
   - Look for linter configs, style guides
   - Check for CONTRIBUTING.md or similar
   - Examine code patterns in main files

   **Exploration principle**: Use git, file system tools, grep, and file reading to thoroughly investigate. Don't stop at the first file - explore until you find a definitive answer or confirm the information doesn't exist.

3. **Evaluate findings**

   **If answer found in repo**:
   - Proceed to Step 4 (Update metadata)
   - The answer MUST cite the source: `[Source: <file-path>]`

   **If answer NOT found in repo**:
   - Proceed to Step 5 (Create Issue)
   - This indicates the project itself has a gap that needs development work

4. **If answer found: Update metadata**

   a. **Update meta.md** with the new information:
      - Add to appropriate section (Intuitive/Macro/Micro)
      - Include source citation: `[Source: <file-path>]`

   b. **Update Question status**:
      - Change status to `solved` in `question-meta.jsonl`
      - Update the question file with resolution

   Question file update:
   ```markdown
   ## Resolution

   **Status**: solved
   **Resolved**: YYYY-MM-DD
   **Method**: Found in repository
   **Source**: <file-path>
   **Answer**: <the answer>
   ```

   ```bash
   # Update question-meta.jsonl - change status from "open" to "solved"
   ```

5. **If answer NOT found: Create Issue**

   This means the repository itself lacks the information needed. Create an Issue to track this gap.

   a. **Get next issue ID**:
   ```bash
   last_id=$(cat .laddermoon/issues/last-issue-id)
   next_id=$((last_id + 1))
   ```

   b. **Create Issue file**: `.laddermoon/issues/issue-<N>.md`

   ```markdown
   # Issue <N>

   ## Status
   open

   ## Type
   missing-information

   ## Related Question
   question-<M>

   ## Description
   <What information is missing from the repository>

   ## Why It Matters
   <Why this information should exist in the project>

   ## Suggested Resolution
   <What needs to be added to the repository - e.g., add section to README, create documentation file, etc.>
   ```

   c. **Update tracking files**:
   ```bash
   # Update last-issue-id
   echo "<next_id>" > .laddermoon/issues/last-issue-id
   
   # Append to issue-meta.jsonl
   echo '{"id": <N>, "status": "open"}' >> .laddermoon/issues/issue-meta.jsonl
   ```

   d. **Update Question status** to `issued`:
   ```bash
   # Update question-meta.jsonl - change status to "issued"
   ```

   Question file update:
   ```markdown
   ## Resolution

   **Status**: issued
   **Resolved**: YYYY-MM-DD
   **Method**: Created issue (answer not found in repo)
   **Issue**: issue-<N>
   ```

**Output When Answer Found**

```
## Clarification Complete

**Question**: question-<M>
**Status**: solved

### Answer Found

**Source**: <file-path>
**Answer**: <summary of the answer>

### META Updated

- **Section**: <which section was updated>
- **Information added**: <what was added>
- **Citation**: [Source: <file-path>]

### Files Modified

- .laddermoon/meta.md (updated)
- .laddermoon/questions/question-<M>.md (resolved)
- .laddermoon/questions/question-meta.jsonl (status → solved)
```

**Output When Issue Created**

```
## Clarification Complete

**Question**: question-<M>
**Status**: issued

### Answer Not Found

Searched:
- <file-1>
- <file-2>
- ...

The repository does not contain the information needed to answer this question.

### Issue Created

**Issue**: issue-<N>
**Type**: missing-information
**Description**: <what's missing>

### Files Created/Modified

- .laddermoon/issues/issue-<N>.md (created)
- .laddermoon/issues/last-issue-id (updated)
- .laddermoon/issues/issue-meta.jsonl (appended)
- .laddermoon/questions/question-<M>.md (status → issued)
- .laddermoon/questions/question-meta.jsonl (status → issued)

### Next Steps

The Issue indicates the project needs development work to add the missing information.
After the Issue is resolved, re-run the metadata workflow.
```

**Issue Types**

| Type | When to use |
|------|-------------|
| `missing-information` | Repo lacks documentation or explanation |
| `incomplete-structure` | Project structure is unclear/missing parts |
| `missing-documentation` | Specific docs should exist but don't |
| `unclear-architecture` | Architecture is not documented |

**Guardrails**

- **Search thoroughly** before creating an Issue
- **Never fabricate answers** - if not in repo, create Issue
- **Always cite sources** - every answer needs `[Source: path]`
- **Update both systems** - metadata AND question/issue tracking
- **Be specific in Issues** - clearly describe what's missing and how to fix it
- **Don't modify repo code** - only metadata files, Issues describe what developers should do
