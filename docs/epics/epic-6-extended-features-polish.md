# Epic 6: Extended Features & Polish

**Goal:** Add favorite namespace support, keyboard shortcuts help, configuration reload, and performance optimization to complete the MVP feature set.

**Note:** Stories 5.1 (URL Action Type) and 5.2 (Local Action Type) were removed and replaced by Epic 5's universal template-based action system. URL and local command execution is now handled through template commands (e.g., `command: "open https://..."` or `command: "kubectl get pods -n {{.namespace}}"`).

---

## Story 6.1: Favorite Namespaces Display

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

## Story 6.2: Keyboard Shortcuts Help

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

## Story 6.3: Configuration Reload

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

## Story 6.4: Performance Optimization

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