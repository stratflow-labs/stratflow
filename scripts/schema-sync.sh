#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=./discover_services.sh
source "$SCRIPT_DIR/discover_services.sh"

services_input="${1:-${SERVICES:-}}"
services=()

usage() {
  cat <<'USAGE'
Usage:
  ./scripts/schema-sync.sh [services]

Examples:
  ./scripts/schema-sync.sh
  ./scripts/schema-sync.sh "identity"
  ./scripts/schema-sync.sh "identity"
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
  echo "No services provided for schema sync."
  exit 1
fi

extract_up_sql() {
  local file="$1"
  awk '
    BEGIN { in_up = 1; in_statement_block = 0 }
    {
      if ($0 ~ /^--[[:space:]]*\+goose[[:space:]]+Down/) {
        in_up = 0
      }
      if (!in_up) {
        next
      }
      if ($0 ~ /^--[[:space:]]*\+goose[[:space:]]+StatementBegin/) {
        in_statement_block = 1
        next
      }
      if ($0 ~ /^--[[:space:]]*\+goose[[:space:]]+StatementEnd/) {
        in_statement_block = 0
        next
      }
      if (in_statement_block) {
        next
      }
      if ($0 ~ /^--[[:space:]]*\+goose/) {
        next
      }
      print
    }
  ' "$file"
}

filter_schema_statements() {
  awk '
    function trim(s) {
      gsub(/^[[:space:]\n]+/, "", s)
      gsub(/[[:space:]\n]+$/, "", s)
      return s
    }
    function flush( statement, upper ) {
      statement = trim(buf)
      if (statement == "") {
        buf = ""
        return
      }
      sub(/;[[:space:]]*$/, "", statement)
      upper = toupper(statement)
      if (upper ~ /^(CREATE[[:space:]]+TABLE|ALTER[[:space:]]+TABLE|CREATE[[:space:]]+UNIQUE[[:space:]]+INDEX|CREATE[[:space:]]+INDEX|CREATE[[:space:]]+TYPE|ALTER[[:space:]]+TYPE|CREATE[[:space:]]+DOMAIN|CREATE[[:space:]]+EXTENSION|CREATE[[:space:]]+VIEW|CREATE[[:space:]]+MATERIALIZED[[:space:]]+VIEW)/) {
        gsub(/[Cc][Oo][Nn][Cc][Uu][Rr][Rr][Ee][Nn][Tt][Ll][Yy][[:space:]]+/, "", statement)
        print statement ";"
        print ""
      }
      buf = ""
    }
    {
      line = $0
      if (line ~ /^[[:space:]]*--/) {
        next
      }
      buf = buf line "\n"
      if (line ~ /;[[:space:]]*$/) {
        flush()
      }
    }
    END {
      if (trim(buf) != "") {
        flush()
      }
    }
  '
}

sync_service_schema() {
  local svc="$1"
  local migrations_dir="services/$svc/db/migrations"
  local schema_dir="services/$svc/db/gen"
  local schema_file="$schema_dir/schema.sql"
  local tmp_up
  local tmp_schema
  tmp_up="$(mktemp)"
  tmp_schema="$(mktemp)"

  if [[ ! -d "$migrations_dir" ]]; then
    echo "[$svc] migrations dir not found: $migrations_dir"
    rm -f "$tmp_up" "$tmp_schema"
    return 1
  fi

  local migration_files=()
  while IFS= read -r file; do
    migration_files+=("$file")
  done < <(find "$migrations_dir" -maxdepth 1 -type f -name '*.sql' | sort)

  if [[ "${#migration_files[@]}" -eq 0 ]]; then
    echo "[$svc] no migration files found in $migrations_dir"
    rm -f "$tmp_up" "$tmp_schema"
    return 1
  fi

  for file in "${migration_files[@]}"; do
    extract_up_sql "$file" >> "$tmp_up"
    printf '\n' >> "$tmp_up"
  done

  filter_schema_statements < "$tmp_up" > "$tmp_schema"

  if ! grep -q '[^[:space:]]' "$tmp_schema"; then
    echo "[$svc] generated schema is empty, refusing to overwrite $schema_file"
    rm -f "$tmp_up" "$tmp_schema"
    return 1
  fi

  mkdir -p "$schema_dir"
  {
    echo "-- Code generated from db/migrations. DO NOT EDIT."
    echo
    cat "$tmp_schema"
  } > "$schema_file"

  rm -f "$tmp_up" "$tmp_schema"
  echo "[$svc] schema snapshot updated: $schema_file"
}

echo "Syncing schema snapshots from migrations for: $(join_services_csv "${services[@]}")"
for svc in "${services[@]}"; do
  sync_service_schema "$svc"
done
echo "Schema sync completed successfully."
