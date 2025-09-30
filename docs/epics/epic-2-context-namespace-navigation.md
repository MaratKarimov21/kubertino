# Epic 2: Context & Namespace Navigation

**Goal:** Implement core navigation functionality allowing users to select Kubernetes contexts and filter namespaces with fuzzy search. This epic delivers the primary workflow foundation.

## Story 2.1: Context Selection Screen

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

## Story 2.2: Namespace List Display

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

## Story 2.3: Fuzzy Search for Namespaces

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

## Story 2.4: Split-Pane Layout Implementation

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
