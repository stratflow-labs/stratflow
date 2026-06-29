#!/usr/bin/env bash
set -euo pipefail

QUERY_GEN_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCRIPTS_DIR="$(cd "$QUERY_GEN_DIR/.." && pwd)"

source "$SCRIPTS_DIR/discover_services.sh"
source "$QUERY_GEN_DIR/schema_parser.sh"
source "$QUERY_GEN_DIR/query_generators.sh"

services_input="${1:-${SERVICES:-}}"
services=()

usage() {
  cat <<'USAGE'
Usage:
  ./scripts/query-gen/crud-query-gen.sh [services]

Examples:
  ./scripts/query-gen/crud-query-gen.sh identity
  ./scripts/query-gen/crud-query-gen.sh "identity,strategy-registry"
USAGE
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  usage
  exit 0
fi

if [[ -n "$services_input" ]]; then
  while IFS= read -r svc; do
    [[ -n "$svc" ]] && services+=("$svc")
  done < <(parse_services_input "$services_input")
else
  services=("${DEFAULT_SERVICES[@]}")
fi

if [[ "${#services[@]}" -eq 0 ]]; then
  echo "No services provided for CRUD query generation."
  exit 1
fi

generate_service_queries() {
  local svc="$1"
  local schema_file="services/$svc/db/gen/schema.sql"
  local sqlc_config="services/$svc/db/sqlc.yaml"
  local queries_dir="services/$svc/internal/adapters/postgres/sqlc/gen/queries"

  if [[ ! -f "$sqlc_config" ]] || ! grep -q '../internal/adapters/postgres/sqlc/gen/queries' "$sqlc_config"; then
    echo "[$svc] CRUD queries skipped: services/$svc/db/sqlc.yaml does not include ../internal/adapters/postgres/sqlc/gen/queries"
    return 0
  fi

  if [[ ! -f "$schema_file" ]]; then
    echo "[$svc] schema snapshot not found: $schema_file"
    return 1
  fi

  mkdir -p "$queries_dir"
  find "$queries_dir" -maxdepth 1 -type f -name '*.sql' -delete

  # Проверяем, нужны ли batch queries для этого сервиса
  local enable_batch=1
  if [[ "$svc" == "identity" ]]; then
    enable_batch=0
    echo "[$svc] Batch queries disabled for identity service"
  fi

  local current_table=""
  local specs=()
  local table column column_type

  while IFS='|' read -r table column column_type; do
    if [[ -z "$table" || -z "$column" || -z "$column_type" ]]; then
      continue
    fi

    if [[ -n "$current_table" && "$table" != "$current_table" ]]; then
      generate_table_queries "$queries_dir/$current_table.sql" "$current_table" "${specs[@]}"
      if [[ "$enable_batch" -eq 1 ]]; then
        generate_batch_queries "$queries_dir/$current_table.sql" "$current_table" "${specs[@]}"
      fi
      specs=()
    fi

    current_table="$table"
    if [[ "${#specs[@]}" -eq 0 ]] || ! contains_column "$column" "${specs[@]}"; then
      specs+=("$column|$column_type")
    fi
  done < <(extract_tables "$schema_file")

  if [[ -n "$current_table" ]]; then
    generate_table_queries "$queries_dir/$current_table.sql" "$current_table" "${specs[@]}"
    if [[ "$enable_batch" -eq 1 ]]; then
      generate_batch_queries "$queries_dir/$current_table.sql" "$current_table" "${specs[@]}"
    fi
  fi

  echo "[$svc] CRUD queries generated in $queries_dir"
}

for svc in "${services[@]}"; do
  generate_service_queries "$svc"
done
