# Epic 5: Extended Actions & Polish

**Goal:** Implement URL and local action types, add favorite namespace support, and refine UI/UX based on MVP usage. This epic completes the MVP feature set.

## Story 5.1: URL Action Type

**As a** user,
**I want** to open URLs with namespace/pod variable substitution,
**so that** I can quickly access web interfaces related to my namespace.

**Acceptance Criteria:**

1. URL action type supported in configuration
2. Variables {{namespace}} and {{pod}} substituted in URL template
3. URL opened in default browser automatically
4. URL also printed to console for reference
5. Error shown if browser launch fails
6. Multiple URL actions can be configured per context
7. URL validation performed during config load

## Story 5.2: Local Action Type

**As a** user,
**I want** to execute commands on my local machine,
**so that** I can run kubectl commands or scripts with context.

**Acceptance Criteria:**

1. Local action type supported in configuration
2. Commands executed in local shell with kubectl context set
3. Variables {{namespace}} and {{pod}} available in command template
4. Command output shown in terminal (TUI minimized)
5. Exit code captured and displayed
6. Error handling for command not found or execution failure
7. Environment variables preserved from parent shell

## Story 5.3: Favorite Namespaces Display

**As a** user,
**I want** favorite namespaces displayed at top of the list,
**so that** I can quickly access frequently used namespaces.

**Acceptance Criteria:**

1. Favorite namespaces from config shown before regular namespaces
2. Visual separator between favorites and regular list
3. Favorite namespaces sorted alphabetically
4. Regular namespaces sorted alphabetically
5. Favorites persist across context switches (per-context)
6. Empty favorites list handled gracefully (no separator shown)
7. Favorites update when configuration reloaded

## Story 5.4: Keyboard Shortcuts Help

**As a** user,
**I want** to view all keyboard shortcuts,
**so that** I can learn and remember available commands.

**Acceptance Criteria:**

1. Pressing '?' displays help overlay
2. Help shows: navigation keys, search activation, action shortcuts, quit
3. Help overlay dismissable with ESC or '?'
4. Help overlay does not block main view entirely (semi-transparent or bordered)
5. Global shortcuts listed separately from action shortcuts
6. Help accessible from any screen

## Story 5.5: Configuration Reload

**As a** user,
**I want** to reload configuration without restarting,
**so that** I can test configuration changes quickly.

**Acceptance Criteria:**

1. Pressing 'r' triggers configuration reload
2. Configuration file re-read and re-validated
3. Current view refreshed with new configuration
4. Success message shown if reload successful
5. Error message shown with details if reload fails
6. Current context and namespace selection preserved if still valid
7. Actions list updated to reflect new configuration

## Story 5.6: Performance Optimization

**As a** developer,
**I want** kubertino to meet NFR performance targets,
**so that** it provides the fast experience promised.

**Acceptance Criteria:**

1. Startup time measured and optimized to < 1 second
2. Namespace fetch parallelized for multiple contexts (if needed)
3. Pod fetch optimized with kubectl options (--no-headers, specific fields)
4. Configuration parsing optimized (avoid redundant parsing)
5. Memory profiling done to ensure < 50MB footprint
6. Binary size measured (should be < 20MB)
7. Performance benchmarks documented in README