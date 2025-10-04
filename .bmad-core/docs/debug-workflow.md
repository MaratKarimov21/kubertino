<!-- Powered by BMAD™ Core -->

# Debug Workflow

Structured iterative debugging process for bugs found during user testing after story completion.

## Overview

This workflow handles the scenario where:
- Story passed initial QA review (status: "Done")
- User finds bugs during manual/real-world testing
- Multiple debugging iterations may be needed
- Need to track approaches and avoid accumulating unnecessary code changes

## Workflow Diagram

```
Story "Done" → User finds bug
    ↓
QA: *debug-report {story}
    ↓
Creates: docs/debug/{epic}.{story}.md
    ↓
User confirms QA analysis
    ↓
Dev: *debug-fix docs/debug/{epic}.{story}.md
    ↓
Dev documents approach & implements fix
    ↓
User tests fix
    ↓
    ├─ Works? → QA: *review-debug {story} → Close debug log
    └─ Doesn't work? → QA updates debug log → Back to Dev
```

## Detailed Process

### Step 1: User Reports Bug

User finds issue during testing and provides:
- Description of what went wrong
- Steps to reproduce
- Expected vs actual behavior

### Step 2: QA Creates Debug Report

```bash
@qa
*debug-report 3.2

[Provide bug details when prompted]
```

**QA Agent will**:
- Create `docs/debug/3.2.md`
- Analyze the bug and hypothesize root cause
- Provide specific instructions for Dev
- Include test scenario for verification
- Update story file with debug session reference

### Step 3: Dev Implements Fix

```bash
@dev
*debug-fix docs/debug/3.2.md
```

**Dev Agent will**:
1. Read debug log and QA's analysis
2. Document implementation approach
3. Make minimal, focused changes
4. Add/update tests
5. Verify fix against QA's test scenario
6. Remove any debug/experimental code
7. Update debug log with results

### Step 4: User Tests Fix

User manually tests using same reproduction steps to verify:
- Bug is resolved
- No regression in other functionality

### Step 5A: Fix Works - QA Final Review

```bash
@qa
*review-debug 3.2
```

**QA Agent will**:
- Review code changes for cleanliness
- Check for unnecessary changes left behind
- Verify only minimal required changes remain
- Update debug log with final conclusion
- Close debug log if all clean

### Step 5B: Fix Doesn't Work - Iterate

```bash
@qa
*debug-report 3.2

[Describe what changed: better/worse/same]
```

**QA Agent will**:
- Add new session to existing debug log
- Update analysis based on new information
- Provide revised instructions
- Cycle continues until resolved

## Debug Log Structure

Each debug log (`docs/debug/{epic}.{story}.md`) contains:

### Session N (Iterative)

**QA Analysis**:
- Bug description and reproduction steps
- Root cause hypothesis
- Affected components
- Test scenario for verification
- Specific instructions for Dev

**Dev Response**:
- Implementation approach and hypothesis
- Changes made (files, lines, why)
- Tests added/modified
- Validation results
- Status and notes

**QA Final Review** (when resolved):
- Verification of user testing
- Code quality assessment
- Unnecessary changes check
- Final conclusion

## Key Principles

### For QA Agent
- Be specific and actionable in analysis
- Provide clear hypothesis, not just description
- Give Dev specific starting points
- Include verification steps for fix confirmation
- Track iteration progress

### For Dev Agent
- Document thinking process for learning
- Make minimal changes - resist over-engineering
- Test thoroughly before claiming completion
- Be honest about uncertainties
- **Clean up after yourself** - remove failed attempts
- HALT and request help after 3 failed approaches

### For User
- Test thoroughly with original reproduction steps
- Report clearly whether fix improved situation
- Confirm when bug is fully resolved

## File Locations

- **Debug Logs**: `docs/debug/{epic}.{story}.md`
- **Original Story**: `docs/stories/{epic}.{story}.*.md`
- **Story Reference**: QA Results section → Debug Sessions

## Benefits

1. **Structured Communication**: Clear handoffs between QA/Dev/User
2. **Knowledge Capture**: Documents debugging thought process
3. **Clean Code**: Forces removal of unnecessary changes
4. **Iteration Tracking**: Easy to see what's been tried
5. **Learning Tool**: Helps team understand problem-solving approaches

## Example Usage

```bash
# User finds bug in story 3.2
@qa
*debug-report 3.2

# QA prompts for details, creates docs/debug/3.2.md

# User reviews and approves QA analysis
@dev
*debug-fix docs/debug/3.2.md

# Dev implements, updates debug log

# User tests - works!
@qa
*review-debug 3.2

# QA verifies code is clean, closes debug log
```

## Commands Summary

### QA Agent Commands

- `*debug-report {story}` - Create initial debug report
- `*review-debug {story}` - Final review after fix works

### Dev Agent Commands

- `*debug-fix {debug_log}` - Implement fix from debug log

## Integration with Main Workflow

Debug workflow is **separate from** normal story development:

- Does NOT update story status
- Does NOT update story File List
- Does NOT modify story Tasks/Acceptance Criteria
- ONLY adds reference in story's QA Results → Debug Sessions

Story remains "Done" - debug sessions are post-completion fixes.
