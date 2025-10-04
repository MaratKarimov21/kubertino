<!-- Powered by BMAD™ Core -->

# qa-debug-report

Create structured debug report for bugs found during user testing after story completion.

## Purpose

Document bugs found during real-world testing in a structured format that enables iterative debugging between QA and Dev agents.

## When to Use

- Story status is "Done" but user found bugs during manual testing
- Issue reproduces reliably with specific steps
- Need to track multiple debugging iterations

## Inputs

```yaml
required:
  - story_id: '{epic}.{story}' # e.g., "3.2"
  - bug_description: 'User-provided description of what went wrong'
  - reproduction_steps: 'Steps to reproduce the issue'
optional:
  - expected_behavior: 'What should happen'
  - actual_behavior: 'What actually happens'
  - severity: 'critical|high|medium|low'
```

## Prerequisites

- Story has status "Done" (passed initial QA review)
- Bug is reproducible by user
- Clear steps to reproduce are available

## Process

### 1. Initialize Debug Log

Create directory if missing:
```bash
mkdir -p docs/debug
```

Create debug log file: `docs/debug/{epic}.{story}.md`

### 2. Analyze the Bug

- Review the original story file to understand what was implemented
- Review the code changes from the story's File List
- Analyze the bug description and reproduction steps
- Identify potential root cause categories:
  - Logic error in implementation
  - Missing edge case handling
  - Integration issue between components
  - Configuration or environment issue
  - Misunderstood requirement

### 3. Create Debug Report Structure

Write to `docs/debug/{epic}.{story}.md`:

```markdown
# Debug Log: Story {epic}.{story}

**Story Title**: [Story title from story file]
**Story File**: docs/stories/{epic}.{story}.*.md
**Debug Started**: [ISO-8601 timestamp]
**Current Status**: Open

---

## Session 1 - QA Initial Analysis

**Date**: [ISO-8601 timestamp]
**Reported By**: User
**Analyzed By**: Quinn (QA)

### Bug Description

[User's description of the bug]

### Reproduction Steps

1. [Step 1]
2. [Step 2]
3. [Step 3]

### Expected Behavior

[What should happen]

### Actual Behavior

[What actually happens]

### Severity

[critical|high|medium|low]

### QA Analysis

**Root Cause Hypothesis**:
[Your analysis of what might be causing this issue]

**Affected Components**:
- [Component/file 1]
- [Component/file 2]

**Related Code**:
- File: [filename:line]
  - Relevant code section or function name

### Test Scenario for Verification

**Given**: [Initial state]
**When**: [Action performed]
**Then**: [Expected result]

**Verification Steps**:
1. [How to verify the fix works]
2. [What to check to confirm no regression]

### Instructions for Dev

**Priority**: [P0|P1|P2]

**Specific Tasks**:
1. [Task 1 - be specific about what to investigate/change]
2. [Task 2]
3. [Task 3]

**Important Considerations**:
- [Any constraints or things to be careful about]
- [Areas that might be affected by changes]

---

## Session 1 - Dev Response

[Dev agent will fill this section]
```

### 4. Update Story File Reference

Append to the story file's QA Results section (if it exists) or create new section:

```markdown
### Debug Sessions

- Session 1: Started [Date] - See docs/debug/{epic}.{story}.md
```

### 5. Notify User

Inform user that debug report has been created at `docs/debug/{epic}.{story}.md` and Dev agent should now be invoked to address the issue.

## Output

Creates structured debug log file with:
- Clear bug description and reproduction steps
- QA's analysis and hypothesis
- Specific instructions for Dev agent
- Test scenario for verification
- Placeholder for Dev agent's response

## Key Principles

- Be specific and actionable in your analysis
- Provide clear hypothesis, not just description
- Give Dev agent specific starting points for investigation
- Include verification steps so fix can be confirmed
- Think like a senior engineer debugging a production issue
- Focus on root cause, not just symptoms

## Blocking Conditions

Stop and request clarification if:
- Bug description is too vague to analyze
- Cannot reproduce issue from provided steps
- Missing critical information about expected behavior
- Story file is missing or incomplete
- Cannot access code files mentioned in story

## Next Steps After This Task

1. QA agent creates this report
2. User reviews and confirms analysis is correct
3. User invokes Dev agent with: "Fix bug in docs/debug/{epic}.{story}.md"
4. Dev agent executes `dev-debug-fix.md` task
5. Cycle repeats until bug is resolved
