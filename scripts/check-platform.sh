#!/bin/bash

# Скрипт для проверки платформы удалённой VM
# Использование: ./scripts/check-platform.sh [host] [user]

set -e

HOST="${1:-${DEPLOY_HOST}}"
USER="${2:-${DEPLOY_USER}}"

if [ -z "$HOST" ] || [ -z "$USER" ]; then
    echo "Usage: $0 <host> <user>"
    echo "Or set environment variables: DEPLOY_HOST, DEPLOY_USER"
    exit 1
fi

echo "Checking platform on ${USER}@${HOST}..."

ssh ${USER}@${HOST} << EOF
    echo "=== System Information ==="
    echo "Architecture: \$(uname -m)"
    echo "Kernel: \$(uname -r)"
    echo "OS: \$(cat /etc/os-release | grep PRETTY_NAME | cut -d'"' -f2)"
    echo ""
    echo "=== Go Binary Compatibility ==="
    ARCH=\$(uname -m)
    if [ "\$ARCH" = "x86_64" ]; then
        echo "✓ Compatible with: linux/amd64"
        echo "Binary name: perf-server-linux-amd64"
    elif [ "\$ARCH" = "aarch64" ]; then
        echo "✓ Compatible with: linux/arm64"
        echo "Binary name: perf-server-linux-arm64"
    else
        echo "⚠ Unsupported architecture: \$ARCH"
        echo "You may need to build for a different target"
    fi
    echo ""
    echo "=== Systemd Status ==="
    if systemctl list-unit-files | grep -q perf-server.service; then
        echo "✓ perf-server.service is installed"
        systemctl status perf-server --no-pager || true
    else
        echo "⚠ perf-server.service is not installed"
        echo "Install it with: sudo cp deploy/perf-server.service /etc/systemd/system/"
    fi
EOF
