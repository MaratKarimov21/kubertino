# Sprint Change Proposal
## Action System Simplification - Architectural Refactoring

**Date:** 2025-10-05
**Trigger:** Story 4.2 completion + real-world usage feedback
**Prepared by:** Sarah (Product Owner)
**Status:** ✅ APPROVED

---

## 1. Analysis Summary

### 1.1 Identified Issue

After real-world usage of Kubertino, the current multi-typed action system (pod_exec, url, local) creates unnecessary complexity in understanding and maintaining the codebase. The core insight: users want simple command aliasing with variable substitution, not different action "types."

**Root Cause:** Over-engineering of action system with type-based Strategy pattern when a simpler template-based approach suffices.

### 1.2 Epic Impact

- **Epic 4 (Action System Core):** ✅ Complete - no changes needed
  - Story 4.1: Actions Display Panel - Done
  - Story 4.2: Pod Exec Action Execution - Done (QA passed 90/100)
  - Story 4.3-4.4: Merged into 4.2

- **Epic 5 (Extended Actions & Polish):** 🔄 Restructured → Epic 6
  - Stories 5.1-5.2 (URL/Local action types) - **REMOVED** (replaced by refactoring)
  - Stories 5.3-5.6 - **MOVED** to Epic 6 (renumbered 6.1-6.4)

- **Epic 5 (NEW): Action System Simplification** - Refactoring epic created

### 1.3 Artifact Adjustment Needs

| Artifact | Impact Level | Changes Required |
|----------|-------------|------------------|
| **PRD** | Medium | Update FR4-FR6 (merge into universal template requirement), update config examples |
| **Architecture** | High | Update Config/Action/Context models, remove Strategy pattern, add Template pattern |
| **Epic 4** | Low | Update status of stories 4.3-4.4 to "Merged into 4.2" |
| **Epic 5** | High | **CREATE NEW** - Action System Simplification refactoring epic |
| **Epic 6** | Medium | **RENAME** from Epic 5, remove stories 5.1-5.2, renumber 5.3-5.6 → 6.1-6.4 |

### 1.4 Recommended Path Forward

**Selected: Option 1 - Direct Adjustment / Integration**

Complete Epic 4 as-is, create new Epic 5 for refactoring action system to universal templates, restructure current Epic 5 as Epic 6 with URL/Local stories removed.

**Rationale:**
- Preserves Epic 4 validated work (QA 90/100, 74.5% test coverage)
- Lower risk incremental refactoring vs rewrite
- Clear documentation of architectural evolution

### 1.5 PRD MVP Impact

**MVP Scope Changes:**
- ❌ Remove: FR5 (URL action type), FR6 (Local action type) as separate implementations
- ✅ Add: Universal template-based action system (simpler, more flexible)
- ✅ Keep: All other MVP features unchanged

**Original MVP Goals:** Still achievable with simplified architecture

---

## 2. Specific Proposed Edits

### 2.1 PRD Updates

**File:** `docs/prd.md`

**Edit 1: Update Functional Requirements (Lines 34-46)**

**From:**
```markdown
**FR4:** Support pod_exec action type for running commands inside pods matching configured regex patterns

**FR5:** Support url action type for opening web links with variable substitution ({{namespace}}, {{pod}})

**FR6:** Support local action type for executing commands on local machine with kubectl context
```

**To:**
```markdown
**FR4:** Support universal action template system with variable substitution ({{context}}, {{namespace}}, {{pod}}) for flexible command aliasing

**FR5:** [REMOVED - replaced by FR4 universal templates]

**FR6:** [REMOVED - replaced by FR4 universal templates]
```

**Edit 2: Add Configuration Example Section (New section after Line 149)**

**Add:**
```markdown
### Example Configuration

**Basic kubertino.yml structure:**

```yaml
version: "1.0"

# Optional: Override default kubeconfig path
kubeconfig: ~/.kube/config

# Global actions (available for all contexts)
actions:
  - name: "Logs"
    shortcut: "l"
    command: "kubectl logs -n {{namespace}} {{pod}} -f"
  - name: "Open Dashboard"
    shortcut: "d"
    command: "open https://dashboard.example.com/{{context}}/{{namespace}}"

# Favorites - supports two formats:

# Format A: Per-context favorites
favorites:
  production:
    - critical-namespace
    - monitoring
  staging:
    - staging-app

# Format B: Global favorites (alternative)
# favorites:
#   - namespace1
#   - namespace2

contexts:
  - name: production
    default_pod_pattern: ".*"
    # Per-context actions (optional, extend/override global)
    actions:
      - name: "Rails Console"
        shortcut: "c"
        command: "kubectl exec -n {{namespace}} {{pod}} -it -- bundle exec rails console"
        pod_pattern: "rails-web-.*"

  - name: staging
    default_pod_pattern: ".*"
```
```

---

### 2.2 Architecture Document Updates

**File:** `docs/architecture.md`

**Edit 1: Update Configuration Model (Lines 126-148)**

**From:**
```go
type Config struct {
    Version  string    `yaml:"version"`
    Contexts []Context `yaml:"contexts"`
}

type Context struct {
    Name               string   `yaml:"name"`
    Kubeconfig         string   `yaml:"kubeconfig"`
    ClusterURL         string   `yaml:"cluster_url,omitempty"`
    DefaultPodPattern  string   `yaml:"default_pod_pattern"`
    FavoriteNamespaces []string `yaml:"favorite_namespaces"`
    Actions            []Action `yaml:"actions"`
}

type Action struct {
    Name       string `yaml:"name"`
    Shortcut   string `yaml:"shortcut"`
    Type       string `yaml:"type"` // pod_exec, url, local
    Command    string `yaml:"command,omitempty"`
    URL        string `yaml:"url,omitempty"`
    PodPattern string `yaml:"pod_pattern,omitempty"`
}
```

**To:**
```go
type Config struct {
    Version    string              `yaml:"version"`
    Kubeconfig string              `yaml:"kubeconfig,omitempty"` // Optional kubeconfig path override
    Actions    []Action            `yaml:"actions"`              // Global actions for all contexts
    Favorites  interface{}         `yaml:"favorites,omitempty"`  // map[string][]string OR []string
    Contexts   []Context           `yaml:"contexts"`
}

type Context struct {
    Name              string   `yaml:"name"`
    DefaultPodPattern string   `yaml:"default_pod_pattern,omitempty"`
    Actions           []Action `yaml:"actions,omitempty"` // Per-context actions (extend/override global)
}

type Action struct {
    Name        string `yaml:"name"`
    Shortcut    string `yaml:"shortcut"`
    Command     string `yaml:"command"`           // Template with {{context}}, {{namespace}}, {{pod}}
    PodPattern  string `yaml:"pod_pattern,omitempty"` // Optional pod regex override
    Container   string `yaml:"container,omitempty"`   // Optional container name for multi-container pods
    Destructive bool   `yaml:"destructive,omitempty"` // Requires confirmation (optional)
}
```

**Edit 2: Update Architectural Patterns (Lines 87-91)**

**From:**
```markdown
**Command Pattern:** Actions encapsulated as executable commands with templating
- Rationale: Flexible action system, easy to extend with new action types

**Strategy Pattern:** Different action types (pod_exec, url, local) implement common interface
- Rationale: Clean abstraction for action execution, extensible design
```

**To:**
```markdown
**Command Pattern:** Actions encapsulated as executable commands with templating
- Rationale: Flexible action system, simple universal interface

**Template Pattern:** Action commands use template syntax for variable substitution ({{context}}, {{namespace}}, {{pod}})
- Rationale: Maximum flexibility, users define any command structure, eliminates need for action type distinction
```

**Edit 3: Update Components - Command Executor (Section ~line 300+)**

**From:**
```markdown
### Command Executor

**Responsibility:** Execute different action types (pod_exec, url, local) with appropriate handling

**Key Interfaces:**
- `ExecutePodExec(action, context, namespace, pods) error`
- `ExecuteURL(action, context, namespace, pod) error`
- `ExecuteLocal(action, context, namespace, pod) error`
```

**To:**
```markdown
### Command Executor

**Responsibility:** Execute action commands with template variable substitution

**Key Interfaces:**
- `Execute(action, context, namespace, pod) error` - Parse template, substitute variables, execute command

**Template Variables:**
- `{{context}}` - Current Kubernetes context name
- `{{namespace}}` - Selected namespace
- `{{pod}}` - Matched pod name (using pod_pattern regex)

**Implementation:**
- Parse action.Command as Go template
- Substitute variables from current state
- Execute via tea.ExecProcess (for interactive) or exec.Command (for background)
```

---

### 2.3 Epic Updates

**File:** `docs/epics/epic-4-action-system-core.md`

**Edit: Add Story Status Updates (After line 68)**

**Add:**
```markdown

---

## Story 4.3: Context Box Display

**Status:** ✅ Merged into Story 4.2

**Note:** Context box functionality was implemented as part of Story 4.2 (executor package includes context_box.go). No separate story needed.

---

## Story 4.4: Error Handling for Actions

**Status:** ✅ Merged into Story 4.2

**Note:** Comprehensive error handling was implemented in Story 4.2 (executor/errors.go, TUI error display). All acceptance criteria met within Story 4.2.
```

---

**File:** `docs/epics/epic-5-extended-actions-polish.md`

**Action: Rename file to** `docs/epics/epic-6-extended-features-polish.md`

**Content Updates:**

**Change Title (Line 1):**
```markdown
# Epic 6: Extended Features & Polish
```

**Update Goal (Line 3):**
```markdown
**Goal:** Add favorite namespace support, keyboard shortcuts help, configuration reload, and performance optimization to complete the MVP feature set.
```

**Remove Stories 5.1-5.2 (Lines 5-36):**
```markdown
[DELETE Story 5.1: URL Action Type - replaced by Epic 5 universal templates]
[DELETE Story 5.2: Local Action Type - replaced by Epic 5 universal templates]
```

**Renumber Remaining Stories:**
- Story 5.3 → Story 6.1 (Favorite Namespaces Display)
- Story 5.4 → Story 6.2 (Keyboard Shortcuts Help)
- Story 5.5 → Story 6.3 (Configuration Reload)
- Story 5.6 → Story 6.4 (Performance Optimization)

---

**File:** `docs/epics/epic-5-action-system-simplification.md` (NEW FILE)

**Create New Epic File** - See full content in Section 3 below.

---

## 3. New Epic 5: Action System Simplification

**Full Epic Definition:**

```markdown
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
**I want** action commands to support {{context}}, {{namespace}}, {{pod}} variables,
**so that** I can create flexible command aliases.

**Acceptance Criteria:**

1. Action.Command parsed as Go template
2. Variables {{context}}, {{namespace}}, {{pod}} substituted from current state
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

**As a** user upgrading from Epic 4 config,
**I want** clear migration instructions,
**so that** I can update my configuration without breaking my workflow.

**Acceptance Criteria:**

1. Migration guide document created (docs/migration/epic4-to-epic5.md)
2. Before/after config examples provided
3. Breaking changes clearly listed
4. Step-by-step migration instructions
5. Common pitfalls and troubleshooting section
6. Optional: Validation script to check old config format
7. README updated with link to migration guide
8. Examples cover both global and per-context action configurations
```

---

## 4. High-Level Action Plan

**Phase 1: Documentation Updates** (Story 5.4)
1. Update PRD (FR4-FR6, add config examples)
2. Update Architecture (models, patterns, components)
3. Update Epic 4 (mark 4.3-4.4 merged)
4. Create Epic 5, rename Epic 5→6, clean up stories

**Phase 2: Config Model Refactoring** (Story 5.1)
1. Update Config, Context, Action structs
2. Implement favorites dual-format parsing
3. Add global/per-context action merging logic
4. Update validation for new structure
5. Write unit tests

**Phase 3: Executor Refactoring** (Story 5.2)
1. Remove Type-based executor logic
2. Implement template parsing
3. Add variable substitution
4. Update Story 4.2 tests
5. Add template-specific tests

**Phase 4: Favorites Implementation** (Story 5.3)
1. Implement dual-format favorites support
2. Update namespace panel logic
3. Test both formats
4. Verify visual display

**Phase 5: Migration Support** (Story 5.5)
1. Write migration guide
2. Create before/after examples
3. Optional validation script
4. Update README

---

## 5. Agent Handoff Plan

**Immediate Next Steps:**

1. **User Approval:** ✅ APPROVED
2. **Documentation Agent (Sarah - PO):** Execute Story 5.4 (update PRD, Architecture, Epics)
3. **Development Agent (Dave):** Execute Stories 5.1-5.3 (config model, executor, favorites)
4. **Documentation Agent (Sarah - PO):** Execute Story 5.5 (migration guide)
5. **QA Agent (Quinn):** Review and test Epic 5 stories as they complete

**Role Assignments:**
- **PO (Sarah):** Stories 5.4, 5.5 (documentation)
- **Dev (Dave):** Stories 5.1, 5.2, 5.3 (implementation)
- **QA (Quinn):** Review all stories, run tests
- **PM:** Review final Epic 5 completion, approve for production

---

## 6. Success Criteria

**Epic 5 Complete When:**
- ✅ All 5 stories pass QA gates
- ✅ Documentation (PRD, Architecture) updated and reviewed
- ✅ Config model supports global/per-context actions and dual-format favorites
- ✅ Executor uses template-based approach (no Type field)
- ✅ Migration guide published
- ✅ All tests passing with >70% coverage maintained
- ✅ No regression in Epic 4 functionality

**Quality Gates:**
- Code coverage: ≥70% (maintain Epic 4 level)
- Security: Input validation for templates (prevent injection)
- Performance: Template parsing <5ms per action
- Documentation: All artifacts consistent and approved

---

**Prepared by:** Sarah (Product Owner)
**Date:** 2025-10-05
**Approved:** 2025-10-05
**Status:** ✅ APPROVED - Ready for Implementation
