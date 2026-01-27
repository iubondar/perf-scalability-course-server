# Performance Scalability Course Server

Простой HTTP сервер на Go для обучающего курса по производительности и масштабируемости. Используется для нагрузочного тестирования и демонстрации различных техник оптимизации.

## Описание

Сервер предоставляет простой HTTP API и предназначен для:
- Нагрузочного тестирования
- Демонстрации работы с nginx как reverse proxy
- Интеграции с Redis и PostgreSQL (планируется)
- Изучения техник оптимизации производительности

## Структура проекта

```
.
├── cmd/
│   └── server/          # Точка входа приложения
├── internal/
│   ├── config/          # Конфигурация (флаги и переменные окружения)
│   ├── handler/         # HTTP handlers
│   ├── router/          # Маршрутизация (chi router)
│   └── server/          # HTTP server с graceful shutdown
├── deploy/              # Файлы для деплоя (systemd service)
├── scripts/             # Скрипты для деплоя и утилиты
└── .github/
    └── workflows/       # GitHub Actions CI/CD
```

## Быстрый старт

### Локальная разработка

1. Клонируйте репозиторий:
```bash
git clone https://github.com/iubondar/perf-scalability-cource-server.git
cd perf-scalability-cource-server
```

2. Установите зависимости:
```bash
go mod download
```

3. Запустите сервер:
```bash
make run
# или
go run cmd/server/main.go
```

Сервер запустится на `localhost:8000` (по умолчанию).

### Использование Makefile

```bash
make help              # Показать все доступные команды
make build             # Собрать бинарник
make test              # Запустить тесты
make run               # Запустить сервер локально
make build-linux       # Собрать бинарники для Linux
```

## Конфигурация

Сервер можно настроить через:
- **Флаги командной строки**: `-a localhost:8080`
- **Переменные окружения**: `RUN_ADDRESS=:8080`

Приоритет: переменные окружения > флаги > значения по умолчанию

### Примеры

```bash
# Через флаг
go run cmd/server/main.go -a :8080

# Через переменную окружения
RUN_ADDRESS=:8080 go run cmd/server/main.go

# В systemd (см. deploy/perf-server.service)
Environment="RUN_ADDRESS=:8080"
```

## API Endpoints

- `GET /hello` - Hello World handler (возвращает "Hello, World!")

## Тестирование

Запуск всех тестов:
```bash
make test
```

Запуск тестов с покрытием:
```bash
make test-coverage
```

## Деплой

Полная инструкция по деплою находится в [DEPLOY.md](DEPLOY.md).

### Краткая инструкция

1. Настройте VM (см. [DEPLOY.md](DEPLOY.md#настройка-vm-выполняется-один-раз))
2. Добавьте секреты в GitHub (Settings → Secrets and variables → Actions):
   - `DEPLOY_HOST` - IP или hostname VM
   - `DEPLOY_USER` - пользователь для SSH
   - `DEPLOY_SSH_KEY` - приватный SSH ключ
   - `DEPLOY_PATH` (опционально) - путь для деплоя
3. Push в ветку `main` - автоматический деплой запустится

### Ручной деплой

```bash
make deploy DEPLOY_HOST=your-vm.com DEPLOY_USER=ubuntu
```

## CI/CD

Проект использует GitHub Actions для автоматического деплоя:
- При push в `main` автоматически запускаются тесты
- Собираются бинарники для Linux (amd64 и arm64)
- Выполняется деплой на удалённую VM
- Автоматически перезапускается systemd сервис

Workflow файл: [.github/workflows/deploy.yml](.github/workflows/deploy.yml)

## Управление сервисом на VM

```bash
# Статус
sudo systemctl status perf-server

# Логи
sudo journalctl -u perf-server -f

# Перезапуск
sudo systemctl restart perf-server
```

## Технологии

- **Go 1.25.1** - язык программирования
- **chi** - HTTP router
- **zap** - структурированное логирование
- **systemd** - управление сервисом
- **GitHub Actions** - CI/CD

## Планы развития

- [ ] Интеграция с Redis
- [ ] Интеграция с PostgreSQL
- [ ] Health check endpoint (`/health`)
- [ ] Метрики для мониторинга
- [ ] Дополнительные handlers для нагрузочного тестирования

## Лицензия

Этот проект создан для обучающих целей.
