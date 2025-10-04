<!-- Powered by BMAD™ Core -->

# dev-debug-fix

Fix bug documented in debug log with structured approach tracking.

## Purpose

Implement fix for bugs found during user testing, documenting approach and changes in the debug log for iterative refinement.

## When to Use

- QA has created debug report in `docs/debug/{epic}.{story}.md`
- User has confirmed QA's analysis
- Ready to implement fix

## Inputs

```yaml
required:
  - debug_log_path: 'docs/debug/{epic}.{story}.md'
  - story_id: '{epic}.{story}' # e.g., "3.2"
```

## Prerequisites

- Debug log exists with QA's analysis and instructions
- Original story file is accessible
- Have read permissions to modify code

## Process

### 1. Read and Understand

**Read these files in order**:
1. Debug log: `docs/debug/{epic}.{story}.md`
2. Original story: `docs/stories/{epic}.{story}.*.md`
3. Files mentioned in QA's "Affected Components"

### 2. Analyze QA's Hypothesis

- Review QA's root cause hypothesis
- Verify the reproduction steps yourself (if possible)
- Check QA's suggested approach
- Form your own implementation plan

### 3. Document Your Approach

Find the latest "## Session N - Dev Response" section in debug log and fill it:

```markdown
## Session N - Dev Response

**Date**: [ISO-8601 timestamp]
**Developer**: James (Dev)

### Approach

**Root Cause Confirmed**: [Yes/No/Partially]
[Explain what you found]

**Implementation Plan**:
1. [Step 1 of your fix]
2. [Step 2]
3. [Step 3]

**Hypothesis**: [Your understanding of why this will fix the issue]

### Changes Made

**Files Modified**:
- `[filename]`
  - **Change**: [What you changed]
  - **Why**: [Why this change fixes the issue]
  - **Line(s)**: [Line numbers or function name]

**Tests Added/Modified**:
- `[test filename]`
  - **Test**: [What the test validates]
  - **Coverage**: [What scenario this covers]

### Validation Performed

- [ ] Reproduced original bug (confirmed it exists)
- [ ] Applied fix
- [ ] Verified fix resolves the issue per test scenario
- [ ] Ran existing tests - all pass
- [ ] Tested edge cases
- [ ] No regression in related functionality

**Verification Result**: [Pass/Fail/Partial]

### Status

[Ready for User Testing | Need More Information | Alternative Approach Needed]

**Notes**:
[Any additional context, concerns, or follow-up needed]
```

### 4. Implement the Fix

- Make minimal, focused changes addressing the root cause
- Follow the story's coding standards
- Add or update tests to prevent regression
- Run full test suite to ensure no breakage

### 5. Test Against QA's Test Scenario

Execute the exact test scenario provided by QA in the debug log:
- Follow "Given/When/Then" scenario
- Execute "Verification Steps"
- Document results in your Dev Response section

### 6. Update Debug Log Status

If fix is successful and verified:

```markdown
**Current Status**: Resolved - Ready for User Verification

**Resolution**: Fixed in Session N by James (Dev) - [Date]
```

If fix is partial or needs more work:

```markdown
**Current Status**: In Progress - Iteration N

**Next Steps**: [What needs to happen next]
```

### 7. Clean Up Unnecessary Changes

**CRITICAL**: After fixing the bug, review ALL your changes:

- Remove any debug code (console.logs, commented code)
- Remove experimental code that didn't help
- Remove alternative approaches you tried but abandoned
- Keep ONLY the minimal changes that fix the issue
- Ensure code is clean and production-ready

Document removed changes:

```markdown
### Changes Reverted

- `[filename]` - Removed [what] because [why it wasn't needed]
```

## Completion Criteria

Fix is complete when:
- [ ] Root cause identified and documented
- [ ] Minimal fix implemented (no unnecessary changes)
- [ ] Tests added/updated and passing
- [ ] QA's test scenario verified successfully
- [ ] No regression in existing functionality
- [ ] Debug log updated with approach and results
- [ ] All debug/experimental code removed
- [ ] Ready for user verification

## Key Principles

- Document your thinking process for learning and collaboration
- Make minimal changes - resist over-engineering
- Test thoroughly before claiming fix is complete
- Be honest about partial fixes or uncertainties
- Clean up after yourself - remove failed attempts
- Focus on root cause, not symptoms

## Blocking Conditions

Stop and request help if:
- Cannot reproduce the bug from QA's steps
- QA's analysis seems incorrect after investigation
- Fix requires architectural changes beyond bug scope
- Need clarification on expected behavior
- After 3 different approaches, still failing

**CRITICAL**: If blocked after 3 attempts:

```markdown
### Debug Session N - BLOCKED

**Attempts Made**: 3
**Status**: Need Assistance

**What I've Tried**:
1. [Approach 1] - [Result]
2. [Approach 2] - [Result]
3. [Approach 3] - [Result]

**Requesting**: [QA re-analysis | User clarification | Architecture review]
```

Update debug log status to: `**Current Status**: Blocked - Need Assistance`

## After User Testing

If user confirms fix works, QA agent should run final review:

### QA Final Review (New Session)

QA agent adds to debug log:

```markdown
## Session N+1 - QA Final Review

**Date**: [ISO-8601 timestamp]
**Reviewer**: Quinn (QA)

### Verification

**User Testing**: [Pass/Fail]
**Code Review**: [Clean/Has Issues]

### Code Quality Assessment

**Unnecessary Changes Found**:
- [ ] No unnecessary changes - code is clean ✓
- [ ] Found extra changes that should be removed:
  - `[filename]` - [What should be removed and why]

**Recommended Actions**:
- [Keep as-is | Clean up extra code | Refactor approach]

### Conclusion

[Bug resolved successfully | Needs minor cleanup | Requires rework]

**Final Status**: [Closed | Needs Cleanup]
```

If cleanup needed, Dev agent makes final changes and QA reviews again.

## Next Steps After This Task

1. Dev agent implements fix and updates debug log
2. User tests the fix manually
3. If fix works:
   - QA agent performs final review of changes
   - QA ensures no unnecessary code remains
   - QA updates debug log with final conclusion
   - Debug log status → "Closed"
4. If fix doesn't work:
   - QA agent creates new session in debug log
   - Describes what changed (better/worse/same)
   - Provides updated analysis and instructions
   - Cycle repeats

## Story File Updates

**CRITICAL**: This task operates OUTSIDE normal story workflow:
- Do NOT update story status
- Do NOT update story File List
- Do NOT modify story Tasks/Acceptance Criteria

**ONLY update** story's QA Results → Debug Sessions section:
```markdown
### Debug Sessions

- Session N: [Date] - [Status] - See docs/debug/{epic}.{story}.md
```
