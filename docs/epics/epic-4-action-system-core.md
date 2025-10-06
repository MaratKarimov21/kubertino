# Epic 4: Action System Core

**Goal:** Implement the core action execution system with pod_exec type, enabling users to run commands inside pods via keyboard shortcuts. This delivers the primary use case.

## Story 4.1: Actions Display Panel

**As a** user,
**I want** to see available actions with shortcuts,
**so that** I know what operations I can perform.

**Acceptance Criteria:**

1. Right-bottom panel displays actions from configuration
2. Each action shows: shortcut key + action name
3. Actions filtered by current context
4. Actions displayed in configured order
5. Visual grouping by action type (if multiple types present)
6. Scrolling supported for many actions
7. Empty state if no actions configured for context

## Story 4.2: Pod Exec Action Execution

**As a** user,
**I want** to execute pod_exec actions via shortcuts,
**so that** I can quickly run commands in pods.

**Acceptance Criteria:**

1. Pressing configured shortcut key triggers action
2. TUI minimizes to show context/namespace decorative box
3. kubectl exec command constructed with: namespace, pod name, command
4. Command executed with full terminal control (interactive mode)
5. User interacts with command normally (e.g., Rails console)
6. On command exit (Ctrl+D or exit), TUI restores
7. Error message shown if pod not found or exec fails
8. Confirmation prompt for destructive actions (if configured)

## Story 4.3: Context Box Display

**Status:** ✅ Merged into Story 4.2

**Note:** Context box functionality was implemented as part of Story 4.2 (executor package includes `internal/executor/context_box.go`). The context box displays before command execution showing context, namespace, pod, and action name. All acceptance criteria for Story 4.3 were fulfilled within Story 4.2 implementation.

---

## Story 4.4: Error Handling for Actions

**Status:** ✅ Merged into Story 4.2

**Note:** Comprehensive error handling was implemented in Story 4.2 (`internal/executor/errors.go`, TUI error display in `internal/tui/app.go`). All acceptance criteria met within Story 4.2:
- Pod pattern matching errors with helpful context
- Kubectl exec failure handling
- TUI error panel display
- Non-fatal error handling (always return to TUI)

Error logging to file (AC 7) was deferred as optional feature.
