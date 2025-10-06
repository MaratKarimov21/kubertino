# Epic 5: Action System Simplification

**Goal:** Refactor action system from type-based (pod_exec, url, local) to universal template-based approach, simplifying architecture and improving flexibility.

**Background:** After completing Epic 4 and real-world usage, the multi-typed action system proved unnecessarily complex. Users want simple command aliasing with variable substitution, not different action "types." This epic refactors the system to a simpler, more flexible template-based approach.

**Dependencies:** Epic 4 complete

---

## Story 5.1: Refactor Config Model

**As a** developer,
**I want** updated config structures supporting global/per-context actions and dual-format favorites,
**so that** the config model matches the simplified architecture.

**Acceptance Criteria:**

1. Config struct includes top-level Kubeconfig (optional), Actions (global), Favorites (dual-format)
2. Context struct removes FavoriteNamespaces and Actions fields (moved to top-level)
3. Action struct removes Type and URL fields, keeps Command with template syntax
4. Favorites supports two formats: map[string][]string (per-context) OR []string (global)
5. Config validation handles both favorites formats
6. Per-context actions can extend/override global actions
7. Backward compatibility checks warn users of breaking changes
8. Unit tests cover both favorites formats and action merging logic

---

## Story 5.2: Implement Template-Based Action Execution

**As a** user,
**I want** action commands to support {{.context}}, {{.namespace}}, {{.pod}} variables,
**so that** I can create flexible command aliases.

**Acceptance Criteria:**

1. Action.Command parsed as Go template
2. Variables {{.context}}, {{.namespace}}, {{.pod}} substituted from current state
3. Invalid templates detected during config validation
4. Template parsing errors shown with helpful messages
5. Executor refactored to remove Type-based logic
6. tea.ExecProcess pattern maintained for interactive commands
7. All existing Story 4.2 tests updated to template-based approach
8. New tests cover template parsing and variable substitution

---

## Story 5.3: Implement Favorites Dual-Format Support

**As a** user,
**I want** to configure favorites as per-context map OR global list,
**so that** I have flexibility in organizing my namespaces.

**Acceptance Criteria:**

1. Config parser detects favorites format (map vs list)
2. Per-context format: favorites[context_name] returns []string for selected context
3. Global format: favorites returns []string for all contexts
4. Namespace panel displays favorites based on current context
5. Visual separator between favorites and regular namespaces
6. Favorites sorted alphabetically within their section
7. Empty favorites handled gracefully (no separator shown)
8. Unit tests cover both formats and edge cases

---

## Story 5.4: Update Documentation Artifacts

**As a** developer and user,
**I want** PRD and Architecture documents updated,
**so that** documentation reflects the new simplified system.

**Acceptance Criteria:**

1. PRD updated: FR4-FR6 merged into universal template requirement
2. PRD includes configuration example with new structure
3. Architecture updated: Config, Context, Action models reflect changes
4. Architecture updated: Strategy pattern removed, Template pattern added
5. Architecture updated: Executor component description simplified
6. Epic 4 updated: Stories 4.3-4.4 marked as "Merged into 4.2"
7. Epic 5 (old) renamed to Epic 6, stories 5.1-5.2 removed, 5.3-5.6 renumbered
8. All documentation changes reviewed and approved

---

## Story 5.5: Create Migration Guide

**Status:** ⏸️ Deferred (Out of Scope)

**Note:** Migration guide creation has been deferred as users can reference the updated documentation (PRD configuration example, architecture.md) and examples directory for migration guidance. The breaking changes are documented in architecture.md "Key Changes from Previous Version" section.
