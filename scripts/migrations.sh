#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=./discover_services.sh
source "$SCRIPT_DIR/discover_services.sh"

mode="up"
services_input="${SERVICES:-}"
parallel="${PARALLEL:-0}"
services_from_flags=()
positional=()

usage() {
  cat <<'USAGE'
Usage:
  ./scripts/migrations.sh [up] [services] [parallel]
  ./scripts/migrations.sh reset
  ./scripts/migrations.sh up -s identity
  ./scripts/migrations.sh --services "identity"

Options:
  -s, --service NAME       Add one service (repeatable)
      --services LIST      Comma or space separated services
  -p, --parallel           Enable parallel mode for migrations
      --parallel=0|1       Set parallel mode explicitly
      --no-parallel        Disable parallel mode
  -h, --help               Show this help
USAGE
}

append_services() {
  local raw="$1"
  local svc
  while IFS= read -r svc; do
    [[ -n "$svc" ]] && services_from_flags+=("$svc")
  done < <(parse_services_input "$raw")
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    -s|--service)
      [[ $# -ge 2 ]] || { echo "Option $1 requires a value."; exit 1; }
      services_from_flags+=("$2")
      shift 2
      ;;
    --services)
      [[ $# -ge 2 ]] || { echo "Option $1 requires a value."; exit 1; }
      append_services "$2"
      shift 2
      ;;
    --services=*)
      append_services "${1#*=}"
      shift
      ;;
    -p|--parallel)
      parallel="1"
      shift
      ;;
    --parallel=*)
      parallel="${1#*=}"
      shift
      ;;
    --no-parallel)
      parallel="0"
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    --)
      shift
      while [[ $# -gt 0 ]]; do
        positional+=("$1")
        shift
      done
      ;;
    -*)
      echo "Unknown option: $1"
      usage
      exit 1
      ;;
    *)
      positional+=("$1")
      shift
      ;;
  esac
done

if [[ "${#positional[@]}" -gt 0 ]]; then
  case "${positional[0]}" in
    up|reset)
      mode="${positional[0]}"
      positional=("${positional[@]:1}")
      ;;
  esac
fi

if [[ "${#services_from_flags[@]}" -gt 0 ]]; then
  services_input="${services_from_flags[*]}"
elif [[ "${#positional[@]}" -gt 0 && -z "$services_input" ]]; then
  services_input="${positional[0]}"
fi

if [[ "${#positional[@]}" -gt 1 ]]; then
  parallel="${positional[1]}"
fi

compose_string="${COMPOSE:-docker compose}"
read -r -a compose_cmd <<< "$compose_string"
compose() {
  "${compose_cmd[@]}" "$@"
}

if [[ "$mode" == "reset" ]]; then
  compose up -d postgres >/dev/null 2>&1
  echo "Resetting database..."
  compose exec -T postgres sh -ec '
    psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d postgres -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '\''$POSTGRES_DB'\'' AND pid <> pg_backend_pid();" >/dev/null
    psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d postgres -c "DROP DATABASE IF EXISTS \"$POSTGRES_DB\";" >/dev/null
    psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d postgres -c "CREATE DATABASE \"$POSTGRES_DB\";" >/dev/null
  '
  echo "Database reset complete."
  exit 0
fi

if [[ "$mode" != "up" ]]; then
  echo "Unknown mode: $mode"
  usage
  exit 1
fi

if [[ "$parallel" != "0" && "$parallel" != "1" ]]; then
  echo "Invalid parallel value: $parallel (expected 0 or 1)"
  exit 1
fi

log_color="$(printf '%s' "${LOG_COLOR:-}" | tr '[:upper:]' '[:lower:]')"
disable_color="0"
if [[ -n "$log_color" ]]; then
  case "$log_color" in
    1|true)
      disable_color="0"
      ;;
    0|false)
      disable_color="1"
      ;;
    *)
      disable_color="0"
      ;;
  esac
elif [[ -n "${NO_COLOR:-}" || "${TERM:-}" == "dumb" ]]; then
  # Backward compatibility for existing environments using NO_COLOR.
  disable_color="1"
fi

if [[ "$disable_color" == "1" ]]; then
  c_reset=""
  c_cyan=""
  c_green=""
  c_yellow=""
  c_red=""
  c_err_badge=""
  c_badge_reset=""
else
  c_reset="$(printf '\033[0m')"
  c_cyan="$(printf '\033[36m')"
  c_green="$(printf '\033[32m')"
  c_yellow="$(printf '\033[33m')"
  c_red="$(printf '\033[31m')"
  c_err_badge="$(printf '\033[37;41;1m')"
  c_badge_reset="$(printf '\033[0m')"
fi

tmpdir="$(mktemp -d)"
cleanup() {
  rm -rf "$tmpdir"
}
trap cleanup EXIT

services=()
if [[ -n "$services_input" ]]; then
  while IFS= read -r svc; do
    services+=("$svc")
  done < <(parse_services_input "$services_input")
else
  services=("${DEFAULT_SERVICES[@]}")
fi

if [[ "${#services[@]}" -eq 0 ]]; then
  echo "No services provided for migrations."
  exit 1
fi

for svc in "${services[@]}"; do
  if [[ ! -f "services/$svc/cmd/ctl/main.go" ]]; then
    echo "Unknown service: $svc"
    echo "Expected file: services/$svc/cmd/ctl/main.go"
    exit 1
  fi
done

display_services="$(join_services_csv "${services[@]}")"

compose up -d postgres >/dev/null 2>&1
echo "Running migrations for: $display_services"

migration_cmd_for_service() {
  local svc="$1"
  if [[ -f "./services/$svc/cmd/ctl/main.go" ]]; then
    echo "./services/$svc/cmd/ctl/main.go"
    return 0
  fi
  echo "Unknown service: $svc"
  echo "Expected file: services/$svc/cmd/ctl/main.go"
  return 1
}

run_migration() {
  local svc="$1"
  local output_file="$2"
  local cmd
  local cmd_status=0
  if ! cmd="$(migration_cmd_for_service "$svc")"; then
    {
      printf "[%s] migrate\n" "$svc"
      migration_cmd_for_service "$svc"
      printf "\n"
    } >"$output_file" 2>&1
    return 1
  fi

  {
    printf "[%s] migrate\n" "$svc"
    LOG_LEVEL=error go run "$cmd" migrate || cmd_status=$?
    printf "\n"
  } >"$output_file" 2>&1
  return "$cmd_status"
}

render_log() {
  local file="$1"
  while IFS= read -r line || [[ -n "$line" ]]; do
    case "$line" in
      *"goose: successfully migrated database to version:"*)
        continue
        ;;
      \[*\]\ migrate)
        printf "%s%s%s\n" "$c_cyan" "$line" "$c_reset"
        ;;
      OK\ *)
        printf "%s\n" "$line"
        ;;
      *"No migrations to run"*)
        printf "%s%s%s\n" "$c_yellow" "$line" "$c_reset"
        ;;
      *"ERROR"*|*"error"*)
        printf "%s%s%s\n" "$c_red" "$line" "$c_reset"
        ;;
      *)
        printf "%s\n" "$line"
        ;;
    esac
  done <"$file"
}

status=0
if [[ "$parallel" == "1" ]]; then
  pids=()
  for svc in "${services[@]}"; do
    run_migration "$svc" "$tmpdir/$svc.log" &
    pids+=("$!")
  done

  for pid in "${pids[@]}"; do
    if ! wait "$pid"; then
      status=1
    fi
  done

  for svc in "${services[@]}"; do
    render_log "$tmpdir/$svc.log"
  done
else
  for svc in "${services[@]}"; do
    log_file="$tmpdir/$svc.log"
    if ! run_migration "$svc" "$log_file"; then
      status=1
      render_log "$log_file"
      break
    fi
    render_log "$log_file"
  done
fi

if [[ "$status" -eq 0 ]]; then
  printf "%sMigrations completed successfully.%s\n" "$c_green" "$c_reset"
  exit 0
fi

printf "%s[ERROR]%s Migrations failed.\n" "$c_err_badge" "$c_badge_reset"
exit "$status"
