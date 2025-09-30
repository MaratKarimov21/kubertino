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

**As a** user,
**I want** to see a decorative box showing context and namespace,
**so that** I know where I'm executing commands.

**Acceptance Criteria:**

1. Box displayed at top of terminal when TUI minimizes
2. Box includes: context name, namespace name, pod name, action name
3. Box styled with borders using ASCII art or Unicode box drawing
4. Box color-coded or highlighted for visibility
5. Box remains visible throughout command execution
6. Box automatically cleared when returning to TUI

## Story 4.4: Error Handling for Actions

**As a** user,
**I want** clear error messages when actions fail,
**so that** I understand what went wrong and how to fix it.

**Acceptance Criteria:**

1. Pod not found: show error with actual pod pattern and found pods
2. Kubectl exec failed: show kubectl error message
3. Permission denied: show RBAC-related error with context
4. Network timeout: show timeout error with retry option
5. Invalid configuration: prevent action execution, show config error
6. All errors return user to TUI (don't exit application)
7. Errors logged to file (~/.kubertino/logs/errors.log) for debugging
