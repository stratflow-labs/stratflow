#!/usr/bin/env bash
set -euo pipefail

compose_string="${COMPOSE:-docker compose}"
read -r -a compose_cmd <<< "$compose_string"
compose() {
  "${compose_cmd[@]}" "$@"
}

go_bin="${GO_BIN:-}"
if [[ -z "$go_bin" ]]; then
  go_bin="$(command -v go || true)"
fi
if [[ -z "$go_bin" && -x "/opt/homebrew/bin/go" ]]; then
  go_bin="/opt/homebrew/bin/go"
fi
if [[ -z "$go_bin" ]]; then
  echo "Go binary not found. Set GO_BIN or add go to PATH."
  exit 1
fi

compose up -d postgres >/dev/null 2>&1

if [[ $# -gt 0 ]]; then
  LOG_LEVEL="${LOG_LEVEL:-error}" "$go_bin" run ./services/identity/cmd/api/main.go admin "$@"
  exit $?
fi

admin_email="${ADMIN_EMAIL:-}"
admin_login="${ADMIN_LOGIN:-admin}"
admin_name="${ADMIN_NAME:-Admin}"
admin_last_name="${ADMIN_LAST_NAME:-Admin}"
admin_password="${ADMIN_PASSWORD:-}"

read_value() {
  local prompt="$1"
  local default_value="$2"
  local value=""
  if [[ -n "$default_value" ]]; then
    read -r -p "$prompt [$default_value]: " value
    printf '%s' "${value:-$default_value}"
  else
    read -r -p "$prompt: " value
    printf '%s' "$value"
  fi
}

admin_email="$(read_value "Admin email" "$admin_email")"
admin_login="$(read_value "Admin login" "$admin_login")"
admin_name="$(read_value "Admin name" "$admin_name")"
admin_last_name="$(read_value "Admin last name" "$admin_last_name")"

read -r -s -p "Admin password [optional]: " typed_password
echo
if [[ -n "$typed_password" ]]; then
  admin_password="$typed_password"
fi

cmd=(
  "$go_bin" run ./services/identity/cmd/api/main.go admin
  --email "$admin_email"
  --login "$admin_login"
  --name "$admin_name"
  --last-name "$admin_last_name"
  --role admin
  --gender-code 1
)

if [[ -n "$admin_password" ]]; then
  cmd+=(--password "$admin_password")
fi

LOG_LEVEL="${LOG_LEVEL:-error}" "${cmd[@]}"
