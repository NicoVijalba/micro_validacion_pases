#!/usr/bin/env bash
set -euo pipefail

host="${DB_HOST:-db}"
port="${DB_PORT:-3306}"

until nc -z "$host" "$port"; do
  echo "waiting for mysql ${host}:${port}"
  sleep 2
done

echo "mysql is reachable"
