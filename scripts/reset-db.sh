#!/usr/bin/env bash
set -euo pipefail

# Reset PostgreSQL database for Study-UPC.
# Usage: sudo ./scripts/reset-db.sh <db_name> <db_user> <db_password>

DB_NAME="${1:-}"
DB_USER="${2:-}"
DB_PASSWORD="${3:-}"

if [[ -z "$DB_NAME" || -z "$DB_USER" || -z "$DB_PASSWORD" ]]; then
  echo "Usage: sudo $0 <db_name> <db_user> <db_password>"
  exit 1
fi

echo "[WARN] This will DROP and recreate database: ${DB_NAME}"
read -r -p "Type 'yes' to continue: " CONFIRM
if [[ "$CONFIRM" != "yes" ]]; then
  echo "[INFO] Aborted."
  exit 0
fi

echo "[INFO] Terminating active connections..."
sudo -u postgres -H psql -v ON_ERROR_STOP=1 -c \
  "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname='${DB_NAME}' AND pid <> pg_backend_pid();"

echo "[INFO] Dropping database if exists..."
sudo -u postgres -H psql -v ON_ERROR_STOP=1 -c "DROP DATABASE IF EXISTS ${DB_NAME};"

echo "[INFO] Ensuring role exists..."
sudo -u postgres -H psql -v ON_ERROR_STOP=1 -tc "SELECT 1 FROM pg_roles WHERE rolname='${DB_USER}'" | grep -q 1 || \
  sudo -u postgres -H psql -v ON_ERROR_STOP=1 -c "CREATE ROLE ${DB_USER} LOGIN PASSWORD '${DB_PASSWORD}';"

echo "[INFO] Creating database..."
sudo -u postgres -H psql -v ON_ERROR_STOP=1 -c "CREATE DATABASE ${DB_NAME} OWNER ${DB_USER};"

echo "[INFO] Granting privileges..."
sudo -u postgres -H psql -v ON_ERROR_STOP=1 -c "GRANT ALL PRIVILEGES ON DATABASE ${DB_NAME} TO ${DB_USER};"

echo "[INFO] Database reset complete."
echo "Next: start backend to run migrations."
