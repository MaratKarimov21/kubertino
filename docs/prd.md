# Kubertino Product Requirements Document (PRD)

## Goals and Background Context

### Goals

- Create a fast, responsive Kubernetes TUI tool that loads instantly unlike k9s
- Enable rapid namespace selection and navigation with fuzzy search
- Provide configurable actions system for team-specific workflows
- Support multiple Kubernetes contexts seamlessly
- Deliver kubectl-level performance with TUI convenience

### Background Context

Development teams working with large Kubernetes clusters (100+ namespaces on dev, 20+ on prod) face significant productivity loss due to k9s performance issues. k9s uses eager loading - attempting to fetch data for all namespaces simultaneously, resulting in:

- 20-30 second initial load times
- Multiple RBAC permission checks (many failing and retrying)
- Unresponsive UI during data fetching
- Blocking workflow for simple operations like accessing Rails console

This creates friction in daily development workflows, especially for teams deploying per-feature-branch environments. The core issue is architectural: k9s prioritizes "show everything" over "show what's needed now."

Kubertino takes a lazy loading approach - only fetching data for the selected namespace, similar to kubectl's focused operations. This architectural choice enables kubectl-level performance while maintaining TUI navigation benefits.

### Change Log

| Date | Version | Description | Author |
|------|---------|-------------|---------|
| 2025-09-29 | 1.0 | Initial PRD creation | PM |

## Requirements

### Functional

**FR1:** Support multiple Kubernetes contexts with seamless switching via configuration file

**FR2:** Display namespace list with fuzzy search capability (activated via `/` key)

**FR3:** Execute configurable actions on selected namespace using keyboard shortcuts

**FR4:** Support universal action template system with variable substitution ({{.context}}, {{.namespace}}, {{.pod}}) for flexible command aliasing, enabling users to execute any command (kubectl exec, open URLs, local scripts) with dynamic context-aware parameters

**FR5:** Support global actions (available for all contexts) and per-context actions (override/extend global actions)

**FR6:** Display pods matching default pod pattern for selected namespace

**FR7:** Validate configuration file (~/.kubertino.yml) on startup with clear error messages

**FR8:** Show context and namespace information in decorative box when executing pod commands

**FR9:** Support favorite namespaces configuration in two formats: per-context map or global list. Favorites displayed at top of namespace list with color highlighting only (no icons or separators), preserving config file order.

**FR10:** Provide visual feedback for all operations: loading spinners during data fetching, error modal dialogs with retry capability for failures, and command execution status

**FR11:** Handle multiple pods matching regex pattern by silently selecting first pod alphabetically (no special pod concept, no error messages for multiple matches)

### Non Functional

**NFR1:** Initial startup time must be under 1 second (comparable to kubectl)

**NFR2:** Namespace list loading must complete within 2 seconds maximum

**NFR3:** Action execution latency must match kubectl performance (no added overhead)

**NFR4:** Configuration file parsing must fail fast with actionable error messages

**NFR5:** UI must remain responsive during all operations (no blocking renders), with loading spinners displayed during async operations (namespace fetch, pod fetch, action execution)

**NFR6:** Memory footprint must be minimal (< 50MB for typical operation)

**NFR7:** Must work on Linux and macOS (primary developer platforms)

**NFR8:** Binary size should be reasonable for distribution (< 20MB)

### Compatibility Requirements

**CR1:** Compatible with kubectl configuration files (~/.kube/config)

**CR2:** Supports all kubectl contexts without modification

**CR3:** Works with RBAC-restricted clusters (graceful degradation)

**CR4:** Compatible with standard terminal emulators (iTerm2, Terminal.app, Alacritty, etc.)

## User Interface Design Goals

### Overall UX Vision

Kubertino prioritizes speed and keyboard-driven navigation. The interface focuses on information density while maintaining clarity. Split-pane layout provides simultaneous visibility of namespaces, pods, and available actions, reducing context switching and cognitive load.

### Key Interaction Paradigms

- **Keyboard-first navigation**: All operations accessible via keyboard shortcuts
- **Fuzzy search**: Fast filtering without leaving current context
- **Immediate feedback**: Visual confirmation for all actions
- **Collapsible TUI**: Minimizes to decorative header when executing interactive commands

### Core Screens and Views

**Main View (Split-pane Layout):**
- Left panel (50% width): Namespace list with search bar
- Right top panel (25% height): Pod list for selected namespace
- Right bottom panel (25% height): Available actions with shortcuts

**Context Selection View:**
- Full-screen list of configured contexts
- Displayed on initial startup if multiple contexts configured

**Command Execution View:**
- Minimized TUI showing context/namespace box
- Full terminal control handed to executed command
- Returns to TUI on command exit

**Error Handling View:**
- Modal dialog overlay with red border for all errors
- Error message with context and retry option
- "[Press Enter to retry]" / "[Press ESC to dismiss]" instructions
- Modal blocks interaction with main view until dismissed

### Accessibility

- WCAG AA compliance not required (developer tool)
- Clear visual hierarchy with standard terminal colors
- Keyboard navigation for all functions

### Target Platforms

- Terminal-based application (TUI)
- Linux and macOS primary targets
- Windows support not required for MVP

## Technical Assumptions

### Repository Structure

**Monorepo:** Single repository containing all components

### Service Architecture

**Monolith:** Single binary CLI application

### Testing Requirements

**Unit + Integration:**
- Unit tests for configuration parsing and business logic
- Integration tests for kubectl interaction
- Manual testing for UI/UX validation
- No automated E2E tests required for MVP

### Additional Technical Assumptions and Requests

**Language & Runtime:**
- Go 1.21+ (latest stable)
- No CGO dependencies for easy cross-compilation

**UI Framework:**
- Bubble Tea (latest stable version)
- Lip Gloss for styling components

**Kubernetes Interaction:**
- Use client-go library for Kubernetes API interaction
- Alternatively: shell out to kubectl for simplicity in MVP
- Config parsing via kubernetes client-go config package

**Configuration:**
- YAML parsing via gopkg.in/yaml.v3
- Config location: ~/.kubertino.yml
- Support for config validation on startup

**Terminal Management:**
- PTY handling for interactive commands
- Terminal size detection and responsive layout
- Clean terminal restoration on exit

**Build & Distribution:**
- Goreleaser for multi-platform builds
- GitHub releases for distribution
- Homebrew tap for macOS installation

**Development Tools:**
- Standard Go toolchain
- golangci-lint for code quality
- No IDE-specific requirements

### Configuration File Example

**Example ~/.kubertino.yml structure:**

```yaml
version: "1.0"

# Optional: Override default kubeconfig path
kubeconfig: ~/.kube/config

# Global actions (available for all contexts)
actions:
  - name: "View Logs"
    shortcut: "l"
    command: "kubectl logs -n {{.namespace}} {{.pod}} -f --tail=100"

  - name: "Open Dashboard"
    shortcut: "d"
    command: "open https://dashboard.example.com/{{.context}}/{{.namespace}}"

  - name: "Port Forward"
    shortcut: "p"
    command: "kubectl port-forward -n {{.namespace}} {{.pod}} 3000:3000"

# Favorites - Format A: Per-context map
favorites:
  production:
    - critical-app
    - monitoring
    - payment-service
  staging:
    - staging-app
    - test-namespace

# Favorites - Format B: Global list (alternative)
# favorites:
#   - my-namespace
#   - another-namespace

contexts:
  - name: production
    default_pod_pattern: ".*"
    # Per-context actions (optional, extend/override global)
    actions:
      - name: "Rails Console"
        shortcut: "c"
        command: "kubectl exec -n {{.namespace}} {{.pod}} -it -- bundle exec rails console"
        pod_pattern: "rails-web-.*"

      - name: "DB Console"
        shortcut: "b"
        command: "kubectl exec -n {{.namespace}} {{.pod}} -it -- bundle exec rails dbconsole"
        pod_pattern: "rails-web-.*"

  - name: staging
    default_pod_pattern: ".*"
    actions:
      - name: "Bash Shell"
        shortcut: "s"
        command: "kubectl exec -n {{.namespace}} {{.pod}} -it -- /bin/bash"
```

**Key Configuration Features:**

- **Global Actions:** Defined at top-level, available across all contexts
- **Per-Context Actions:** Defined within context, extend or override global actions
- **Template Variables:** `{{.context}}`, `{{.namespace}}`, `{{.pod}}` substituted at runtime
- **Dual Favorites Format:** Supports both per-context map and global list
- **Pod Pattern Matching:** Optional regex for selecting specific pods
- **Flexible Commands:** Any shell command (kubectl, open, scripts, etc.)

## Epic List

The following epics deliver Kubertino functionality in logical, sequential increments:

**Epic 1: Foundation & Core CLI** - Establish project structure, configuration parsing, and basic CLI framework

**Epic 2: Context & Namespace Navigation** - Implement context selection and namespace list with fuzzy search

**Epic 3: Pod Discovery & Display** - Add pod listing and pattern matching for selected namespace

**Epic 4: Action System Core** - Implement configurable action execution framework with pod_exec type

**Epic 5: Action System Simplification** - Refactor actions to universal template-based system with global/per-context configuration

**Epic 6: UI/UX Bug Fixes & Refinement** - Fix critical navigation, pod selection, actions, and error handling bugs

## Epic 1: Foundation & Core CLI

**Goal:** Establish robust project foundation with configuration system, error handling, and basic terminal UI framework. This epic sets up all infrastructure needed for rapid feature development in subsequent epics.

### Story 1.1: Project Initialization and Structure

**As a** developer,
**I want** a well-organized Go project with proper module structure,
**so that** the codebase is maintainable and follows Go best practices.

**Acceptance Criteria:**

1. Go module initialized with appropriate module path
2. Project structure includes: cmd/, internal/, pkg/ directories
3. README.md with project description and build instructions
4. Makefile with build, test, lint targets
5. .gitignore configured for Go projects
6. GitHub repository initialized with MIT license
7. golangci-lint configuration present

### Story 1.2: Configuration File Parser

**As a** user,
**I want** the tool to read and parse ~/.kubertino.yml configuration,
**so that** I can customize contexts, namespaces, and actions.

**Acceptance Criteria:**

1. YAML configuration file loaded from ~/.kubertino.yml
2. Configuration structure includes: contexts, default_pod_pattern, favorite_namespaces, actions
3. Each context includes: name, kubeconfig path, namespace favorites, actions array
4. Actions include: name, shortcut, type (pod_exec/url/local), command/url template
5. Configuration validation on load with specific error messages
6. Example configuration file provided in repository
7. Unit tests cover valid and invalid configuration scenarios
8. Shortcut conflicts detected and reported

### Story 1.3: Kubernetes Context Detection

**As a** user,
**I want** kubertino to detect available kubectl contexts,
**so that** I can work with my existing Kubernetes configurations.

**Acceptance Criteria:**

1. Read kubectl config from standard location (~/.kube/config)
2. Parse available contexts from kubeconfig
3. Match configured contexts in ~/.kubertino.yml with available kubectl contexts
4. Warning displayed for configured contexts not found in kubeconfig
5. Error displayed if no valid contexts available
6. Unit tests mock kubeconfig file reading

### Story 1.4: Basic TUI Framework

**As a** developer,
**I want** a Bubble Tea TUI framework initialized,
**so that** subsequent epics can build UI components efficiently.

**Acceptance Criteria:**

1. Bubble Tea application scaffold created
2. Basic model/update/view pattern implemented
3. Keyboard input handling framework established
4. Terminal size detection and responsive layout foundation
5. Clean exit handling (Ctrl+C, ESC, 'q')
6. Error display component for showing validation errors
7. Application launches without crashing

## Epic 2: Context & Namespace Navigation

**Goal:** Implement core navigation functionality allowing users to select Kubernetes contexts and filter namespaces with fuzzy search. This epic delivers the primary workflow foundation.

### Story 2.1: Context Selection Screen

**As a** user,
**I want** to select a Kubernetes context on startup,
**so that** I can work with the correct cluster.

**Acceptance Criteria:**

1. Full-screen context list displayed on startup
2. Contexts from configuration file shown with names
3. Arrow keys navigate context list
4. Enter key selects context and transitions to namespace view
5. Currently selected context highlighted visually
6. Number of namespaces shown per context (if available in config)
7. ESC or 'q' exits application

### Story 2.2: Namespace List Display

**As a** user,
**I want** to see a list of namespaces for selected context,
**so that** I can choose which namespace to work with.

**Acceptance Criteria:**

1. Left panel (50% width) displays namespace list
2. Namespaces fetched via kubectl for selected context
3. Favorite namespaces (from config) displayed at top of list
4. Arrow keys navigate namespace list
5. Currently selected namespace highlighted
6. Namespace count displayed in header
7. Loading indicator shown during namespace fetch
8. Error message displayed if namespace fetch fails

### Story 2.3: Fuzzy Search for Namespaces

**As a** user,
**I want** to filter namespaces using fuzzy search,
**so that** I can quickly find specific namespaces in large lists.

**Acceptance Criteria:**

1. Pressing '/' activates search mode
2. Search input box appears at bottom of namespace panel
3. Typing filters namespace list in real-time
4. Fuzzy matching algorithm matches non-contiguous characters
5. Matching characters highlighted in results
6. ESC clears search and returns to full list
7. Enter on filtered result selects that namespace
8. Empty search shows all namespaces
9. "No matches" message when search returns empty

### Story 2.4: Split-Pane Layout Implementation

**As a** user,
**I want** to see namespaces, pods, and actions in split-pane layout,
**so that** I have full context visibility.

**Acceptance Criteria:**

1. Layout divided: left 50% (namespaces), right 50% split horizontally
2. Right top 50% (pods section) - initially empty with placeholder text
3. Right bottom 50% (actions section) - initially empty with placeholder text
4. Header bar shows selected context name
5. Panel borders clearly defined
6. Layout responsive to terminal resize
7. Minimum terminal size enforced (80x24 characters)
8. Warning message if terminal too small

## Epic 3: Pod Discovery & Display

**Goal:** Implement pod listing functionality for selected namespace with pattern matching, enabling users to see available pods before executing actions.

### Story 3.1: Pod List Retrieval

**As a** user,
**I want** to see pods in the selected namespace,
**so that** I know what resources are available.

**Acceptance Criteria:**

1. When namespace selected, fetch pods via kubectl
2. Display pod names in right-top panel (25% of screen)
3. Show pod status (Running, Pending, Failed) with color coding
4. Loading indicator while fetching pods
5. Error message if pod fetch fails
6. Empty state message if no pods in namespace
7. Pod list updates when switching namespaces

### Story 3.2: Default Pod Pattern Matching

**As a** user,
**I want** pods matching default_pod_pattern to be highlighted,
**so that** I understand which pod will be used for actions.

**Acceptance Criteria:**

1. Apply default_pod_pattern regex to pod list
2. First matching pod marked as "default" with visual indicator
3. If no pods match pattern, show warning message
4. If multiple pods match, first one selected automatically
5. Pattern compiled once per context (cached)
6. Invalid regex in config shows error on context load
7. Tooltip or help text explains default pod concept

### Story 3.3: Pod List Navigation

**As a** user,
**I want** to navigate the pod list with keyboard,
**so that** I can inspect different pods (future: execute pod-specific actions).

**Acceptance Criteria:**

1. Tab key switches focus between namespace and pod panels
2. Arrow keys navigate pod list when focused
3. Visual indicator shows which panel has focus
4. Selected pod highlighted differently than default pod marker
5. Scrolling works for long pod lists
6. Shift+Tab moves focus backwards

## Epic 4: Action System Core

**Goal:** Implement the core action execution system with pod_exec type, enabling users to run commands inside pods via keyboard shortcuts. This delivers the primary use case.

### Story 4.1: Actions Display Panel

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

### Story 4.2: Pod Exec Action Execution

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

### Story 4.3: Context Box Display

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

### Story 4.4: Error Handling for Actions

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

## Epic 5: Action System Simplification

**Goal:** Refactor action system from type-based (pod_exec, url, local) to universal template-based approach, simplifying architecture and improving flexibility.

**Note:** This epic was added after Epic 4 to address unnecessary complexity in the multi-typed action system. Stories 5.1-5.3 refactor the config model and executor to use Go templates with variable substitution, replacing the type-based approach.

### Story 5.1: Refactor Config Model

**Status:** ✅ Complete

Configuration model refactored to support global/per-context actions and dual-format favorites (per-context map or global list).

### Story 5.2: Implement Template-Based Action Execution

**Status:** ✅ Complete

Action commands now use Go template syntax with {{.context}}, {{.namespace}}, {{.pod}} variables, eliminating the need for separate action types.

### Story 5.3: Implement Favorites Dual-Format Support

**Status:** ✅ Complete

Favorites configuration supports both per-context map and global list formats with proper TUI integration.

## Epic 6: UI/UX Bug Fixes & Refinement

**Goal:** Fix critical UI/UX bugs in namespace list navigation, pod selection, actions system, and error handling to ensure production-ready user experience.

**Context:** Epic 1-5 delivered core functionality, but during testing accumulated 9 critical UI/UX issues that must be resolved before MVP release.

**Note:** This epic replaces the previous "Extended Features & Polish" epic. Stories related to keyboard shortcuts help, configuration reload, and performance optimization were deemed non-critical for MVP and removed from scope.

### Story 6.1: Namespace List Navigation Fixes

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

### Story 6.2: Pod Selection & Actions System Fixes

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

### Story 6.3: Error Handling & Loading States

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

## Next Steps

### UX Expert Prompt

**Not Required:** This is a terminal-based tool for developers. UI/UX specification is already defined in the PRD (split-pane layout, keyboard navigation). Proceed directly to architecture.

### Architect Prompt

Create a technical architecture document for Kubertino that addresses:

1. **Go Project Structure:** Detailed package organization following best practices
2. **Bubble Tea Architecture:** Model/update/view pattern for each screen component
3. **Kubernetes Client Integration:** Whether to use client-go or shell out to kubectl
4. **Configuration Management:** YAML parsing, validation, and runtime reloading
5. **Terminal Management:** PTY handling for interactive commands, terminal state management
6. **Error Handling Strategy:** Error types, logging approach, user-facing error messages
7. **Testing Strategy:** Unit test coverage, integration test approach, manual test scenarios
8. **Build and Distribution:** Cross-compilation, releases, installation methods

Reference this PRD for all requirements and epic details. Focus on architectural decisions that enable fast, maintainable development while meeting performance NFRs.