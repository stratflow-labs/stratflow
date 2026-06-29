#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=./discover_services.sh
source "$SCRIPT_DIR/discover_services.sh"

mode="${1:-gen}"
services_input="${2:-${SERVICES:-}}"
services=()

action_for_mode() {
  case "$1" in
    gen|generate)
      echo "generate"
      ;;
    vet)
      echo "vet"
      ;;
    verify)
      echo "verify"
      ;;
    *)
      echo ""
      ;;
  esac
}

usage() {
  cat <<'USAGE'
Usage:
  ./scripts/sqlc.sh <gen|vet|verify> [services]

Examples:
  ./scripts/sqlc.sh gen
  ./scripts/sqlc.sh gen "identity"
  ./scripts/sqlc.sh vet "identity"
  ./scripts/sqlc.sh verify "identity"
USAGE
}

action="$(action_for_mode "$mode")"
if [[ -z "$action" ]]; then
  echo "Unknown mode: $mode"
  usage
  exit 1
fi

if [[ -n "$services_input" ]]; then
  while IFS= read -r svc; do
    services+=("$svc")
  done < <(parse_services_input "$services_input")
else
  services=("${DEFAULT_SERVICES[@]}")
fi

if [[ "${#services[@]}" -eq 0 ]]; then
  echo "No services provided for sqlc $action."
  exit 1
fi

if ! command -v sqlc >/dev/null 2>&1; then
  echo "sqlc is not installed."
  echo "Install: go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest"
  exit 1
fi

for svc in "${services[@]}"; do
  if [[ ! -f "services/$svc/cmd/ctl/main.go" ]]; then
    echo "Unknown service: $svc"
    echo "Expected file: services/$svc/cmd/ctl/main.go"
    exit 1
  fi
  if [[ ! -f "services/$svc/db/sqlc.yaml" ]]; then
    echo "Service $svc does not have sqlc config."
    echo "Expected file: services/$svc/db/sqlc.yaml"
    exit 1
  fi
done

echo "Running sqlc $action for: $(join_services_csv "${services[@]}")"

if [[ "$action" == "generate" ]]; then
  echo "Syncing schema snapshots from migrations..."
  "$SCRIPT_DIR/schema-sync.sh" "${services[*]}"
  echo "Generating CRUD query snapshots from schemas..."
  "$SCRIPT_DIR/query-gen/crud-query-gen.sh" "${services[*]}"
fi

for svc in "${services[@]}"; do
  echo "[$svc] sqlc $action"
  if [[ "$action" == "generate" ]]; then
    find "services/$svc/internal/adapters/postgres/sqlc/gen" -maxdepth 1 -type f -name '*.sql.go' -delete
  fi
  (cd "services/$svc/db" && sqlc "$action")
done

echo "SQLC $action completed successfully."
