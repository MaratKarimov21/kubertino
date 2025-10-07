# Epic 6: UI/UX Bug Fixes & Refinement

**Goal:** Fix critical UI/UX bugs in namespace list navigation, pod selection, actions system, and error handling to ensure production-ready user experience.

**Context:** Epic 1-5 delivered core functionality, but during testing accumulated 9 critical UI/UX issues that must be resolved before MVP release. This epic consolidates all bugfixes into logical stories.

---

## Story 6.1: Namespace List Navigation Fixes

**As a** user,
**I want** namespace list navigation to behave correctly with proper cursor management and favorites display,
**so that** I can efficiently navigate large namespace lists without confusion.

**Acceptance Criteria:**

1. **Favorites Visual Display:**
   - Favorite namespaces displayed with color highlighting only (no star icon, no visual separator)
   - Favorites appear in the same order as defined in config file (not alphabetically sorted)
   - Regular namespaces continue after favorites in alphabetical order

2. **Cursor Visibility During Scroll:**
   - Cursor remains visible at all times when navigating namespace list
   - No disappearing cursor after several key presses during scroll

3. **Cursor Centering Behavior:**
   - When scrolling down a long list:
     - Cursor starts at top of visible list
     - Cursor moves down until reaching middle of visible viewport
     - List scrolls while keeping cursor centered in viewport
     - When bottom of list reached, cursor continues to bottom
   - Same behavior applies in reverse when scrolling up

4. **Cursor Persistence on Focus Change:**
   - When pressing Enter on namespace, cursor remains on selected namespace
   - Focus shifts to pod panel, but namespace selection visually persists
   - User can see which namespace is currently active

**Technical Notes:**
- Component: `internal/tui/panels/namespaces.go`
- Related: Favorites parsing in `internal/config/config.go`

---

## Story 6.2: Pod Selection & Actions System Fixes

**As a** user,
**I want** pod selection and actions to work correctly without confusing "special pods" concept,
**so that** I can execute actions reliably on selected pods.

**Acceptance Criteria:**

1. **Remove "Special Pods" Concept:**
   - Eliminate all logic related to distinguishing "special" pods
   - Always use first pod matching the action's `pod_pattern` (or default pattern)
   - Remove error "multiple pods match pattern ^backend.*" - this is expected behavior
   - If multiple pods match, silently select first one (alphabetically)

2. **Actions Panel Non-Focusable:**
   - Actions panel cannot receive keyboard focus (remove from tab order)
   - Actions always visible but not interactive via navigation
   - Action shortcuts work globally regardless of focused panel
   - If actions list is long, display in multiple columns (2-3 columns) instead of scrolling

3. **Action Template Variable Substitution:**
   - When action executed, correctly substitute `{{.namespace}}` with selected namespace
   - Correctly substitute `{{.pod}}` with matched pod name
   - Correctly substitute `{{.context}}` with current context name
   - Verify substitution works for all action command templates

**Technical Notes:**
- Components: `internal/tui/panels/pods.go`, `internal/tui/panels/actions.go`
- Executor: `internal/executor/executor.go`
- Template engine: Go `text/template`

---

## Story 6.3: Error Handling & Loading States

**As a** user,
**I want** clear error feedback via modal dialogs and loading indicators,
**so that** I understand application state and can recover from errors easily.

**Acceptance Criteria:**

1. **Remove Bottom Log Panel:**
   - Completely remove any log/status display at bottom of screen
   - Logs should only go to log file (`~/.kubertino/logs/kubertino.log`)
   - Screen real estate reclaimed for main panels

2. **Error Modal Dialog:**
   - When any error occurs, display modal dialog overlay on top of all panels
   - Modal has red border and red-themed styling
   - Modal shows: error message, affected operation, suggestion (if any)
   - Modal displays "[Press Enter to retry]" message
   - Pressing Enter re-executes the failed operation
   - Pressing ESC dismisses modal and returns to previous state
   - Modal blocks all other input while displayed

3. **Loading Spinners:**
   - Display spinner indicator when fetching namespaces (on context switch)
   - Display spinner indicator when fetching pods (on namespace selection)
   - Display spinner indicator during action execution (before command takes over terminal)
   - Spinner shown in relevant panel (namespace panel during NS fetch, pod panel during pod fetch)
   - Spinner styling consistent across all panels

**Technical Notes:**
- Components: `internal/tui/app.go` (global error state), all panels (loading states)
- New component: `internal/tui/components/modal.go` (error modal)
- New component: `internal/tui/components/spinner.go` (loading indicator)
- Error handling flow refactor required

---

**Epic 6 Completion Criteria:**
- All 3 stories implemented and tested
- Manual testing checklist passed for all scenarios
- No regressions in Epic 1-5 functionality
