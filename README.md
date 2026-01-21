# Kubertino

A fast, lazy-loading Kubernetes TUI (Terminal User Interface) tool built with Go and Bubble Tea. Kubertino provides a performant alternative to k9s, optimized for developers working with large Kubernetes clusters.

<img width="1272" height="784" alt="Снимок экрана 2026-01-21 в 13 04 46" src="https://github.com/user-attachments/assets/e4a72906-0a7d-4e7f-b621-348d384c8df2" />
<img width="932" height="632" alt="Снимок экрана 2026-01-21 в 13 21 17" src="https://github.com/user-attachments/assets/93d72be6-32f9-4f9e-95a8-1630248cb167" />
<img width="1277" height="783" alt="Снимок экрана 2025-10-12 в 23 17 27" src="https://github.com/user-attachments/assets/422d780d-8b9a-442f-9849-692578dee2b2" />


## Overview

Kubertino offers:
- **Fast startup**: <1s initialization time
- **Lazy loading**: Only loads data for the selected namespace, avoiding the performance bottleneck of loading all namespaces simultaneously
- **kubectl-level performance**: Shells out to kubectl for operations, maintaining native compatibility
- **Intuitive split-pane UI**: Three-panel interface showing namespaces, pods, and configurable actions
- **Fuzzy search**: Quick namespace filtering
- **Configurable actions**: Execute pod commands, open URLs, or run local commands via keyboard shortcuts
- **Interactive pod access**: Execute commands directly in pods (shell, logs, port-forward, etc.)

## Requirements

- **Go**: Version 1.21 or higher
- **kubectl**: Must be installed and configured with access to your Kubernetes clusters

## Installation

### Build from Source

```bash
git clone https://github.com/maratkarimov/kubertino.git
cd kubertino
make build
```

The binary will be created as `kubertino` in the project root.

### Install Binary

```bash
sudo mv kubertino /usr/local/bin/
```

## Usage

```bash
# Launch kubertino
kubertino

# The TUI will start and allow you to:
# 1. Select a Kubernetes context (if multiple are configured)
# 2. Search and select a namespace
# 3. View pods in the selected namespace
# 4. Execute configured actions via keyboard shortcuts
```

## Configuration

Kubertino uses a YAML configuration file located at `~/.kubertino.yml`.

An example configuration will be available in the `examples/` directory.

Configuration supports:
- Multiple Kubernetes contexts
- Custom kubeconfig file paths
- Per-context action definitions
- Favorite namespaces
- Default pod patterns for action targeting
- Custom keyboard shortcuts

## Logs

Kubertino writes application logs to `~/.kubertino/kubertino.log` to avoid interfering with the TUI display.

You can view logs in real-time using:

```bash
tail -f ~/.kubertino/kubertino.log
```

## Development

### Build

```bash
make build
```

### Test

```bash
make test
```

### Lint

```bash
make lint
```

### Clean Build Artifacts

```bash
make clean
```

## Project Structure

```
kubertino/
├── cmd/kubertino/          # Application entry point
├── internal/               # Private application code
│   ├── config/            # Configuration parsing and validation
│   ├── k8s/               # Kubernetes adapter (kubectl integration)
│   ├── tui/               # Bubble Tea TUI components
│   ├── executor/          # Action execution (pod exec, URLs, local commands)
│   └── search/            # Fuzzy search implementation
├── pkg/                   # Public libraries (future use)
├── examples/              # Example configuration files
└── scripts/               # Utility scripts
```

## License

MIT License - See LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.
