# Sprint Change Proposal: Epic 6 Overhaul

**Date:** 2025-10-07
**Author:** Sarah (Product Owner)
**Status:** Pending Approval

---

## Executive Summary

Epic 6 requires complete overhaul to address 9 critical UI/UX bugs discovered during testing. The existing "Extended Features & Polish" epic (Stories 6.1-6.4) is being replaced with "UI/UX Bug Fixes & Refinement" containing 3 focused stories addressing namespace navigation, pod selection/actions, and error handling issues.

**Impact:** Critical for MVP release quality
**Effort:** Medium (estimated 13-21 SP across 3 stories)
**Timeline:** No delay to MVP (these fixes are blockers for release)

---

## 1. Change Context & Trigger

### Triggering Issue

Story 6.1 (Favorite Namespaces Display) was partially implemented but contains multiple bugs. During testing, accumulated 9 critical UI/UX issues across namespace navigation, pod selection, actions system, and error handling that block MVP release.

### Issue Classification

- **Type:** Accumulated technical bugs + outdated stories
- **Scope:** UI/UX layer (TUI components)
- **Severity:** Critical - all 9 bugs must be fixed for MVP

### Root Cause

Epic 1-5 focused on feature delivery. Epic 6 was originally planned for "nice-to-have" features (keyboard shortcuts help, config reload, performance optimization) but actual need is bugfix-focused epic to ensure production quality.

---

## 2. Detailed Bug List (9 Issues)

### Namespace List Issues (4 bugs):
1. **Favorites visualization** - Star icon and separator present, should be color-only; not respecting config order
2. **Cursor disappears** - Cursor becomes invisible after several scroll operations
3. **No cursor centering** - Long list scrolling doesn't center cursor in viewport
4. **Cursor not persisting** - Selected namespace loses visual highlight when focus shifts to pod panel

### Pod Selection & Actions Issues (3 bugs):
5. **"Special pods" concept** - Error "multiple pods match pattern ^backend.*" shown when should silently select first
6. **Actions panel focusable** - Actions can receive keyboard focus, should be display-only; no multi-column layout for long lists
7. **Template substitution broken** - Actions not correctly substituting `{{.namespace}}` and `{{.pod}}` variables

### Error Handling & Loading Issues (2 bugs):
8. **Log panel at bottom** - Error logs displayed at bottom of screen, should use modal dialog with retry
9. **No loading indicators** - Missing spinners during async operations (namespace fetch, pod fetch, action execution)

---

## 3. Epic Impact Analysis

### Current Epic 6 Status

**Epic 6: Extended Features & Polish**
- ❌ Story 6.1: Favorite Namespaces Display - Partially implemented, **4 critical bugs**
- ❌ Story 6.2: Keyboard Shortcuts Help - Not started, **not needed for MVP**
- ❌ Story 6.3: Configuration Reload - Not started, **not needed for MVP**
- ❌ Story 6.4: Performance Optimization - Not started, **not needed for MVP**

### Recommended Action

**Delete** existing Epic 6 entirely, **replace** with new Epic 6 containing 3 bugfix-focused stories.

### Impact on Other Epics

- **Epic 1-5:** ✅ Completed, no changes required
- **Future Epics:** ✅ No other epics planned currently

---

## 4. Artifact Conflicts & Required Updates

### PRD Updates Required

| Section | Change Required | Rationale |
|---------|----------------|-----------|
| **FR9** | Update favorites description: "color highlighting only (no icons/separators), preserving config file order" | Current text mentions separator/sorting behavior that's incorrect |
| **FR10** | Expand to include: "loading spinners, error modal dialogs with retry capability" | Current text only mentions action execution feedback |
| **FR11** | Clarify: "silently selecting first pod alphabetically (no special pod concept)" | Current text doesn't address "special pods" removal |
| **NFR5** | Add: "with loading spinners displayed during async operations" | Missing specification for loading states |
| **UI Design Goals** | Add "Error Handling View" section describing modal dialog | Missing from current UI specifications |
| **Epic 6** | Complete replacement with new 3-story epic | Existing epic no longer relevant |

### Architecture Updates Required

| Section | Change Required | Rationale |
|---------|----------------|-----------|
| **Components** | Add: Error Modal Component, Spinner Component descriptions | New UI components needed for bugfixes |
| **Error Handling Strategy** | Add: Error Modal Pattern with implementation example | Current strategy uses inline errors, need modal-based approach |
| **Architectural Patterns** | Add: UI Pattern Updates (Epic 6) section | Document the pattern changes from bugfixes |
| **Source Tree** | Add: `internal/tui/components/` directory with modal.go and spinner.go | New source files required |

---

## 5. Proposed New Epic 6

### Epic 6: UI/UX Bug Fixes & Refinement

**Goal:** Fix critical UI/UX bugs in namespace list navigation, pod selection, actions system, and error handling to ensure production-ready user experience.

---

### Story 6.1: Namespace List Navigation Fixes

**Acceptance Criteria:** (4 ACs addressing bugs 1-4)

1. Favorites displayed with color highlighting only, preserving config file order
2. Cursor remains visible at all times during scroll operations
3. Cursor centers in viewport during long list scrolling
4. Cursor persists on selected namespace when focus changes to pod panel

**Estimate:** 3-5 SP
**Components:** `internal/tui/panels/namespaces.go`, `internal/config/config.go`

---

### Story 6.2: Pod Selection & Actions System Fixes

**Acceptance Criteria:** (3 ACs addressing bugs 5-7)

1. Remove "special pods" logic, silently select first matching pod alphabetically
2. Actions panel non-focusable, multi-column layout for long action lists
3. Template variable substitution correctly replaces `{{.context}}`, `{{.namespace}}`, `{{.pod}}`

**Estimate:** 5-8 SP
**Components:** `internal/tui/panels/pods.go`, `internal/tui/panels/actions.go`, `internal/executor/executor.go`

---

### Story 6.3: Error Handling & Loading States

**Acceptance Criteria:** (3 ACs addressing bugs 8-9)

1. Remove bottom log panel, logs only to file
2. Error modal dialog with red border, retry on Enter, dismiss on ESC
3. Loading spinners for all async operations (namespace fetch, pod fetch, action execution)

**Estimate:** 5-8 SP
**Components:** `internal/tui/app.go`, `internal/tui/components/modal.go`, `internal/tui/components/spinner.go`

---

## 6. Path Forward Evaluation

### Option 1: Direct Adjustment ✅ **RECOMMENDED**

**Action:** Replace Epic 6 with bugfix-focused epic, update PRD & Architecture

**Effort:**
- Documentation updates: 2-3 hours
- Implementation: 13-21 SP (3 stories)
- Testing: Manual testing checklist for all 9 bugs

**Work Discarded:** None (Stories 6.2-6.4 never started)

**Risks:** Minimal (isolated to UI layer, no architectural changes)

**Timeline Impact:** None (bugfixes are MVP blockers anyway)

**Long-term Sustainability:** Excellent (eliminates technical debt before release)

### Option 2: Rollback ❌ NOT APPLICABLE

Story 6.1 needs refinement, not rollback. Epic 1-5 unaffected.

### Option 3: MVP Re-scoping ❌ NOT NEEDED

MVP scope unchanged. This is bugfix work, not feature changes.

---

## 7. Specific Proposed Edits

### ✅ Completed Edits

**New Epic File:**
- ✅ Created: `docs/epics/epic-6-ui-ux-bugfixes.md` with 3 stories

**PRD Changes:**
- ✅ FR9: Updated favorites description
- ✅ FR10: Expanded visual feedback requirements
- ✅ FR11: Clarified multiple pods handling
- ✅ NFR5: Added loading spinner requirement
- ✅ UI Design Goals: Added Error Handling View section
- ✅ Epic 6: Completely replaced with new bugfix epic
- ✅ Epic List: Updated Epic 6 summary line

**Architecture Changes:**
- ✅ Components: Added Error Modal Component and Spinner Component
- ✅ Error Handling Strategy: Added Error Modal Pattern section
- ✅ Architectural Patterns: Added UI Pattern Updates (Epic 6)
- ✅ Source Tree: Added `internal/tui/components/` directory

**Old Epic File:**
- ⏳ TODO: Archive or delete `docs/epics/epic-6-extended-features-polish.md`

---

## 8. Next Steps

### Immediate Actions

1. **✅ Get user approval** for this Sprint Change Proposal
2. **⏳ Archive old epic:** Move `docs/epics/epic-6-extended-features-polish.md` to archive or delete
3. **⏳ Update backlog:** Add 3 new stories (6.1, 6.2, 6.3) to sprint backlog
4. **⏳ Prioritize:** Story 6.1 highest priority (most visible UX issues)

### Development Sequence

**Recommended order:**
1. Story 6.1 (Namespace Navigation) - Most visible bugs, independent
2. Story 6.3 (Error Handling) - Foundation for better error UX
3. Story 6.2 (Pod Selection/Actions) - Depends on Story 6.3 for error modal

**Parallel work possible:**
- Story 6.1 and 6.3 can be developed in parallel (different components)

### Testing & Validation

**Manual Testing Checklist:** Create checklist covering all 9 bug scenarios

**Acceptance Gates:**
- All 9 bugs verified fixed
- No regressions in Epic 1-5 functionality
- Manual testing checklist 100% pass

---

## 9. Risk Assessment

### Low Risk ✅

- **Scope:** Isolated to UI layer (TUI components)
- **Architecture:** No fundamental changes, additive components only
- **Dependencies:** No external dependencies affected
- **Rollback:** Easy to revert if needed (Git)

### Mitigation Strategies

- **Regression Prevention:** Manual test checklist for Epic 1-5 functionality
- **Incremental Delivery:** Complete stories in priority order, validate before next
- **Code Review:** Careful review of template substitution changes (Story 6.2)

---

## 10. Success Criteria

✅ **Epic 6 considered successful when:**

1. All 9 bugs fixed and verified
2. PRD and Architecture documents updated and consistent
3. Manual testing checklist passes 100%
4. No regressions in Epic 1-5 functionality
5. User acceptance of new error handling UX (modal dialogs)
6. User acceptance of loading indicators (spinners)

---

## Approval

**Proposed by:** Sarah (Product Owner)
**Date:** 2025-10-07

**Awaiting approval from:** Project Stakeholder

**Approval Status:** ⏳ Pending

---

**Once approved, proceed with:**
1. Archive old Epic 6 file
2. Create backlog tickets for Stories 6.1, 6.2, 6.3
3. Begin development starting with Story 6.1
