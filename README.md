# Order Management App (Go + PostgreSQL + Nginx + systemd)

Минимальная система управления заказами с backend на Go, PostgreSQL, простым frontend и reverse proxy через Nginx.

## Что внутри
- REST API для заказов
- PostgreSQL как БД
- Встроенный frontend на HTML/CSS/JS
- Запуск приложения как systemd unit
- Готовый Nginx-конфиг
- SQL-миграция и install-скрипт

## Структура API
- `GET /api/health`
- `GET /api/orders`
- `GET /api/orders/:id`
- `POST /api/orders`
- `PUT /api/orders/:id`
- `DELETE /api/orders/:id`

Пример JSON:
```json
{
  "customer": "Ivan Petrov",
  "email": "ivan@example.com",
  "amount": 1599.90,
  "status": "new",
  "description": "Первый заказ"
}
```

## Быстрый запуск
```bash
cd order-management-app
sudo bash scripts/install.sh
```

После установки:
1. Открой `/opt/order-management-app/.env`
2. Поменяй `DB_PASSWORD`
3. Задай такой же пароль в PostgreSQL:
```bash
sudo -u postgres psql
ALTER ROLE orders_app WITH PASSWORD 'your_strong_password';
\q
```
4. Перезапусти:
```bash
sudo systemctl restart orders-app
```

## Ручной запуск
### Пакеты
```bash
sudo apt update
sudo apt install -y postgresql postgresql-contrib nginx curl ca-certificates
```
### Go
```bash
curl -fsSL https://go.dev/dl/go1.22.12.linux-amd64.tar.gz -o /tmp/go.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf /tmp/go.tar.gz
export PATH=$PATH:/usr/local/go/bin
```
### Пользователь и каталог
```bash
sudo useradd --system --create-home --shell /usr/sbin/nologin orders || true
sudo mkdir -p /opt/order-management-app
sudo cp -r ./* /opt/order-management-app/
cd /opt/order-management-app
sudo chown -R orders:orders /opt/order-management-app
```
### PostgreSQL
```bash
sudo -u postgres psql
CREATE ROLE orders_app WITH LOGIN PASSWORD 'change_me';
CREATE DATABASE orders_db OWNER orders_app;
\q
PGPASSWORD=change_me psql -h 127.0.0.1 -U orders_app -d orders_db -f migrations/001_init.sql
```
### Env
```bash
cp .env.example .env
nano .env
```
### Сборка
```bash
go mod tidy
go build -o orders-app ./cmd/server
```
### systemd
```bash
sudo cp deploy/systemd/orders-app.service /etc/systemd/system/orders-app.service
sudo systemctl daemon-reload
sudo systemctl enable --now orders-app
```
### Nginx
```bash
sudo cp deploy/nginx/orders.conf /etc/nginx/sites-available/orders.conf
sudo ln -sf /etc/nginx/sites-available/orders.conf /etc/nginx/sites-enabled/orders.conf
sudo rm -f /etc/nginx/sites-enabled/default
sudo nginx -t
sudo systemctl restart nginx
```

Проверка:
```bash
curl http://127.0.0.1:8080/api/health
curl http://127.0.0.1/api/health
systemctl status orders-app
systemctl status nginx
```
