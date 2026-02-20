#!/bin/bash

# Скрипт для ручного деплоя на удалённую VM
# Использование: ./scripts/deploy.sh [host] [user] [deploy_path]

set -e

HOST="${1:-${DEPLOY_HOST}}"
USER="${2:-${DEPLOY_USER}}"
DEPLOY_PATH="${3:-${DEPLOY_PATH:-/local/perf-server}}"

if [ -z "$HOST" ] || [ -z "$USER" ]; then
    echo "Usage: $0 <host> <user> [deploy_path]"
    echo "Or set environment variables: DEPLOY_HOST, DEPLOY_USER, DEPLOY_PATH"
    exit 1
fi

echo "Building binaries..."
GOOS=linux GOARCH=amd64 go build -o perf-server-linux-amd64 ./cmd/server
GOOS=linux GOARCH=arm64 go build -o perf-server-linux-arm64 ./cmd/server

echo "Copying binaries to server..."
scp perf-server-linux-* ${USER}@${HOST}:${DEPLOY_PATH}/

echo "Deploying on server..."
ssh ${USER}@${HOST} << EOF
    set -e
    
    # Определяем архитектуру
    ARCH=\$(uname -m)
    if [ "\$ARCH" = "x86_64" ]; then
        BINARY_NAME="perf-server-linux-amd64"
    elif [ "\$ARCH" = "aarch64" ]; then
        BINARY_NAME="perf-server-linux-arm64"
    else
        echo "Unsupported architecture: \$ARCH"
        exit 1
    fi
    
    # Создаём директорию для деплоя, если её нет
    mkdir -p ${DEPLOY_PATH}
    
    # Копируем новый бинарник
    mv ${DEPLOY_PATH}/\$BINARY_NAME ${DEPLOY_PATH}/perf-server
    chmod +x ${DEPLOY_PATH}/perf-server
    
    # Перезапускаем сервис (все инстансы при multi-instance setup)
    sudo systemctl daemon-reload
    sudo systemctl restart 'perf-server@*' 2>/dev/null || sudo systemctl restart perf-server
    
    # Проверяем статус
    sleep 2
    sudo systemctl status 'perf-server@*' --no-pager 2>/dev/null || sudo systemctl status perf-server --no-pager || exit 1
    
    echo "Deployment completed successfully!"
EOF

echo "Cleaning up local binaries..."
rm -f perf-server-linux-*

echo "Deployment finished!"
