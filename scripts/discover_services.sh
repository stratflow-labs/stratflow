#!/usr/bin/env bash

DEFAULT_SERVICES=(identity strategy-registry)

# parse_services_input prints services one per line.
# Supports input formats: "a b c" and "a,b,c".
parse_services_input() {
  local raw="${1:-}"
  local normalized="${raw//,/ }"
  local parsed=()
  read -r -a parsed <<< "$normalized"
  local svc
  for svc in "${parsed[@]}"; do
    [[ -n "$svc" ]] && printf '%s\n' "$svc"
  done
}

# join_services_csv joins services by comma+space for display.
join_services_csv() {
  local out=""
  local svc
  for svc in "$@"; do
    if [[ -n "$out" ]]; then
      out+=", "
    fi
    out+="$svc"
  done
  printf '%s' "$out"
}
