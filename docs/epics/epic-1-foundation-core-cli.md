# Epic 1: Foundation & Core CLI

**Goal:** Establish robust project foundation with configuration system, error handling, and basic terminal UI framework. This epic sets up all infrastructure needed for rapid feature development in subsequent epics.

## Story 1.1: Project Initialization and Structure

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

## Story 1.2: Configuration File Parser

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

## Story 1.3: Kubernetes Context Detection

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

## Story 1.4: Basic TUI Framework

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
