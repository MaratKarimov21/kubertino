# Epic 3: Pod Discovery & Display

**Goal:** Implement pod listing functionality for selected namespace with pattern matching, enabling users to see available pods before executing actions.

## Story 3.1: Pod List Retrieval

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

## Story 3.2: Default Pod Pattern Matching

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

## Story 3.3: Pod List Navigation

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
