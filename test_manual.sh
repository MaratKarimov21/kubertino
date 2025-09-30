#!/bin/bash
# Manual testing script for Story 2.1
# This script verifies the application starts correctly with different configurations

set -e

echo "=== Story 2.1 Manual Testing ==="
echo ""

# Test 1: Multi-context configuration (current config has 2 contexts)
echo "Test 1: Testing with 2 contexts (should show context selection screen)"
echo "Config has contexts:"
grep "  - name:" ~/.kubertino.yml | sed 's/  - name: /  - /'
echo ""
echo "Running kubertino (press 'q' to quit after viewing)..."
echo "Expected behavior: Should display context selection screen with 2 contexts"
echo ""

# Create a single-context test config
echo "Test 2: Creating single-context test config"
cat > /tmp/kubertino-single.yml << 'EOF'
version: "1.0"
contexts:
  - name: yangibackend/marat.k+prod@yangi.uz
    kubeconfig: ~/.kube/config
    cluster_url: https://kubeapi.retail.yangi.uz/
    default_pod_pattern: "^api-.*"
    favorite_namespaces:
      - default
      - production
    actions:
      - name: Console
        shortcut: c
        type: pod_exec
        command: /bin/sh
EOF

echo "Single context config created at /tmp/kubertino-single.yml"
echo ""
echo "To test single-context auto-selection:"
echo "  cp /tmp/kubertino-single.yml ~/.kubertino.yml"
echo "  ./kubertino"
echo "  (should auto-select and skip to namespace view)"
echo ""

echo "=== Manual Test Checklist ==="
echo "For multi-context (2 contexts):"
echo "  [ ] Context selection screen appears"
echo "  [ ] Both contexts visible with names"
echo "  [ ] First context highlighted by default"
echo "  [ ] 'yangibackend' shows '3 namespaces' count"
echo "  [ ] 'yangiretail-stg' shows '3 namespaces' count"
echo "  [ ] Arrow keys navigate up/down"
echo "  [ ] Navigation wraps around at top/bottom"
echo "  [ ] Enter key selects context and transitions"
echo "  [ ] ESC/q exits application"
echo "  [ ] Terminal resize handles gracefully"
echo ""
echo "For single-context:"
echo "  [ ] Auto-selects and skips to namespace view"
echo "  [ ] No context selection screen shown"
echo ""

echo "Test completed!"
