# Инструкция по деплою

Этот документ описывает процесс настройки и деплоя сервера на удалённую виртуальную машину.

## Предварительные требования

- Удалённая VM с Linux (Ubuntu/Debian/CentOS)
- Доступ по SSH с ключами
- Права sudo на VM
- Go 1.25.1 или выше (для локальной сборки)

## Настройка VM (выполняется один раз)

### 1. Проверка платформы

Проверьте архитектуру вашей VM:

```bash
uname -m
```

Результат:
- `x86_64` - используйте бинарник `perf-server-linux-amd64`
- `aarch64` - используйте бинарник `perf-server-linux-arm64`

Или используйте скрипт:

```bash
make check-platform DEPLOY_HOST=your-vm.com DEPLOY_USER=ubuntu
```

### 2. Создание пользователя для сервиса

Создайте непривилегированного пользователя для запуска сервера (это **отдельный** пользователь от вашего SSH-пользователя):

```bash
sudo useradd -r -s /bin/false -d /local/perf-server perf-server
sudo mkdir -p /local/perf-server
sudo chown perf-server:perf-server /local/perf-server
```

**Важно:** Пользователь `perf-server` используется только для запуска самого сервера Go. Для SSH-подключения и деплоя используется ваш обычный пользователь на VM (который будет указан в `DEPLOY_USER`).

### 3. Настройка SSH доступа

#### Генерация SSH ключа (если ещё нет)

На вашей локальной машине:

```bash
ssh-keygen -t ed25519 -C "github-actions-deploy" -f ~/.ssh/deploy_key
```

#### Копирование публичного ключа на VM

Замените `user` на вашего SSH-пользователя (этот же пользователь будет использоваться в `DEPLOY_USER`):

```bash
ssh-copy-id -i ~/.ssh/deploy_key.pub user@your-vm.com
```

Или вручную:

```bash
cat ~/.ssh/deploy_key.pub | ssh user@your-vm.com "mkdir -p ~/.ssh && cat >> ~/.ssh/authorized_keys"
```

#### Настройка sudo без пароля (для пользователя деплоя)

На VM выполните:

```bash
sudo visudo
```

Добавьте строку (замените `user` на вашего SSH-пользователя - того же, что будет в `DEPLOY_USER`):

```
user ALL=(ALL) NOPASSWD: /bin/systemctl daemon-reload, /bin/systemctl restart perf-server, /bin/systemctl restart perf-server@*, /bin/systemctl status perf-server, /bin/systemctl status perf-server@*
```

**Примечание:** `DEPLOY_USER` - это тот же пользователь, с которым вы подключаетесь по SSH. Это не отдельный пользователь, а ваш обычный пользователь на VM (например, `ubuntu`, `deploy` или ваш личный пользователь).

### 4. Установка systemd service

Скопируйте service файл на VM:

```bash
scp deploy/perf-server.service user@your-vm.com:/tmp/
```

На VM установите service:

```bash
sudo cp /tmp/perf-server.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable perf-server
```

### 5. Настройка переменных окружения

Отредактируйте service файл при необходимости:

```bash
sudo systemctl edit perf-server
```

Добавьте или измените переменные окружения:

```ini
[Service]
Environment="RUN_ADDRESS=:8080"
```

Или отредактируйте напрямую:

```bash
sudo nano /etc/systemd/system/perf-server.service
```

После изменений:

```bash
sudo systemctl daemon-reload
sudo systemctl restart perf-server
```

## Настройка GitHub Actions

### 1. Добавление секретов в GitHub

Перейдите в настройки репозитория: `Settings → Secrets and variables → Actions`

Добавьте следующие секреты:

- **DEPLOY_HOST** - IP адрес или hostname вашей VM (например, `192.168.1.100` или `vm.example.com`)
- **DEPLOY_USER** - **тот же пользователь, с которым вы подключаетесь по SSH** (например, `ubuntu`, `deploy` или ваш обычный пользователь на VM). Это не отдельный пользователь - используйте того же, с кем обычно подключаетесь по SSH.
- **DEPLOY_SSH_KEY** - приватный SSH ключ (содержимое файла `~/.ssh/deploy_key`)
- **DEPLOY_PATH** (опционально) - путь для размещения бинарника (по умолчанию `/local/perf-server`)
- **DEPLOY_PORT** (опционально) - SSH порт (по умолчанию `22`)

#### Как получить приватный ключ:

```bash
cat ~/.ssh/deploy_key
```

Скопируйте весь вывод, включая строки `-----BEGIN ... KEY-----` и `-----END ... KEY-----`.

### 2. Проверка workflow

После добавления секретов, при push в ветку `main` автоматически запустится workflow, который:

1. Запустит тесты
2. Соберёт бинарники для Linux (amd64 и arm64)
3. Скопирует бинарники на VM
4. Определит архитектуру и выберет нужный бинарник
5. Заменит бинарник и перезапустит сервис
6. Проверит статус сервиса

## Ручной деплой

### Использование Makefile

```bash
make deploy DEPLOY_HOST=your-vm.com DEPLOY_USER=ubuntu
```

Или через переменные окружения:

```bash
export DEPLOY_HOST=your-vm.com
export DEPLOY_USER=ubuntu
export DEPLOY_PATH=/local/perf-server  # опционально
make deploy
```

### Использование скрипта напрямую

```bash
./scripts/deploy.sh your-vm.com ubuntu /local/perf-server
```

## Управление сервисом

### Проверка статуса

```bash
sudo systemctl status perf-server
```

### Просмотр логов

```bash
sudo journalctl -u perf-server -f
```

Или последние 100 строк:

```bash
sudo journalctl -u perf-server -n 100
```

### Перезапуск

```bash
sudo systemctl restart perf-server
```

### Остановка

```bash
sudo systemctl stop perf-server
```

### Запуск

```bash
sudo systemctl start perf-server
```

## Устранение проблем

### Сервис не запускается

1. Проверьте логи:
   ```bash
   sudo journalctl -u perf-server -n 50
   ```

2. Проверьте права на файл:
   ```bash
   ls -l /local/perf-server/perf-server
   sudo chmod +x /local/perf-server/perf-server
   ```

3. Проверьте, что бинарник для правильной архитектуры:
   ```bash
   file /local/perf-server/perf-server
   ```

### SSH подключение не работает

1. Проверьте, что ключ добавлен в GitHub Secrets правильно (включая переносы строк)
2. Проверьте права на ключ на VM:
   ```bash
   chmod 600 ~/.ssh/authorized_keys
   ```
3. Проверьте SSH подключение вручную:
   ```bash
   ssh -i ~/.ssh/deploy_key user@your-vm.com
   ```

### Permission denied при перезапуске сервиса

Убедитесь, что настроен sudo без пароля (см. раздел "Настройка SSH доступа").

## Проверка работы сервера

После деплоя проверьте, что сервер отвечает:

```bash
curl http://your-vm.com:8080/hello
```

Должен вернуться ответ: `Hello, World!`

## Масштабирование воркеров (multi-instance)

Для увеличения пропускной способности можно запустить несколько инстансов сервера за nginx с балансировкой нагрузки.

### 1. Установка template unit

Скопируйте template unit на VM:

```bash
scp deploy/perf-server@.service user@your-vm.com:/tmp/
```

На VM:

```bash
sudo cp /tmp/perf-server@.service /etc/systemd/system/
sudo systemctl daemon-reload
```

### 2. Запуск нескольких воркеров

Остановите одиночный инстанс (если был):

```bash
sudo systemctl stop perf-server
sudo systemctl disable perf-server
```

Включите и запустите 4 воркера (порты 8080–8083):

```bash
sudo systemctl enable perf-server@0 perf-server@1 perf-server@2 perf-server@3
sudo systemctl start perf-server@0 perf-server@1 perf-server@2 perf-server@3
```

Проверка статуса:

```bash
sudo systemctl status perf-server@0 perf-server@1 perf-server@2 perf-server@3
```

### 3. Управление инстансами

```bash
# Перезапуск всех воркеров
sudo systemctl restart 'perf-server@*'

# Логи конкретного инстанса
sudo journalctl -u perf-server@0 -f

# Добавить ещё один воркер (порт 8084)
sudo systemctl enable perf-server@4
sudo systemctl start perf-server@4
```

## Интеграция с nginx (опционально)

Если вы используете nginx как reverse proxy, скопируйте конфиг:

```bash
scp deploy/nginx-perf-server.conf user@your-vm.com:/tmp/
```

На VM:

```bash
sudo cp /tmp/nginx-perf-server.conf /etc/nginx/sites-available/perf-server
sudo ln -sf /etc/nginx/sites-available/perf-server /etc/nginx/sites-enabled/
# Отредактируйте server_name при необходимости
sudo nginx -t && sudo systemctl reload nginx
```

Конфиг `deploy/nginx-perf-server.conf` содержит upstream с балансировкой между портами 8080–8083. Для одиночного инстанса используйте упрощённый вариант:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```
