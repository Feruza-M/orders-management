#!/usr/bin/env bash
set -euo pipefail
APP_DIR=/opt/order-management-app
BIN_NAME=orders-app
SERVICE_NAME=orders-app
APP_USER=orders
DB_NAME=orders_db
DB_USER=orders_app
if [[ $EUID -ne 0 ]]; then echo "Run as root"; exit 1; fi
apt update
apt install -y postgresql postgresql-contrib nginx curl ca-certificates
id -u "$APP_USER" >/dev/null 2>&1 || useradd --system --create-home --shell /usr/sbin/nologin "$APP_USER"
mkdir -p "$APP_DIR"
cp -r ./* "$APP_DIR"/
cd "$APP_DIR"
if ! command -v go >/dev/null 2>&1; then
  curl -fsSL https://go.dev/dl/go1.22.12.linux-amd64.tar.gz -o /tmp/go.tar.gz
  rm -rf /usr/local/go
  tar -C /usr/local -xzf /tmp/go.tar.gz
  export PATH=$PATH:/usr/local/go/bin
fi
export PATH=$PATH:/usr/local/go/bin
go mod tidy
go build -o "$BIN_NAME" ./cmd/server
cp .env.example .env || true
chown -R "$APP_USER":"$APP_USER" "$APP_DIR"
runuser -u postgres -- psql -c "DO $$ BEGIN IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = '${DB_USER}') THEN CREATE ROLE ${DB_USER} LOGIN PASSWORD 'change_me'; END IF; END $$;"
runuser -u postgres -- psql -tc "SELECT 1 FROM pg_database WHERE datname='${DB_NAME}'" | grep -q 1 || runuser -u postgres -- createdb "$DB_NAME"
runuser -u postgres -- psql -d "$DB_NAME" -c "GRANT ALL PRIVILEGES ON DATABASE ${DB_NAME} TO ${DB_USER};"
runuser -u postgres -- psql -d "$DB_NAME" -c "GRANT ALL ON SCHEMA public TO ${DB_USER};"
PGPASSWORD=change_me psql -h 127.0.0.1 -U "$DB_USER" -d "$DB_NAME" -f migrations/001_init.sql
cp deploy/systemd/orders-app.service /etc/systemd/system/orders-app.service
systemctl daemon-reload
systemctl enable --now "$SERVICE_NAME"
cp deploy/nginx/orders.conf /etc/nginx/sites-available/orders.conf
ln -sf /etc/nginx/sites-available/orders.conf /etc/nginx/sites-enabled/orders.conf
rm -f /etc/nginx/sites-enabled/default
nginx -t
systemctl restart nginx
echo "Done. Open http://SERVER_IP/"
