#!/usr/bin/env bash

indent_list() {
  local values=("$@")
  local idx
  for idx in "${!values[@]}"; do
    local suffix=","
    if [[ "$idx" -eq $((${#values[@]} - 1)) ]]; then
      suffix=""
    fi
    printf '    %s%s\n' "${values[$idx]}" "$suffix"
  done
}

placeholder_list() {
  local count="$1"
  local idx
  for ((idx = 1; idx <= count; idx++)); do
    local suffix=","
    if [[ "$idx" -eq "$count" ]]; then
      suffix=""
    fi
    printf '    $%d%s\n' "$idx" "$suffix"
  done
}

set_list() {
  local arg_offset="$1"
  shift
  local values=("$@")
  local idx
  for idx in "${!values[@]}"; do
    local suffix=","
    if [[ "$idx" -eq $((${#values[@]} - 1)) ]]; then
      suffix=""
    fi
    printf '    %s = $%d%s\n' "${values[$idx]}" "$((idx + arg_offset))" "$suffix"
  done
}

where_columns() {
  local idx=1
  local column
  for column in "$@"; do
    if [[ "$idx" -eq 1 ]]; then
      printf 'WHERE %s = $%d\n' "$column" "$idx"
    else
      printf '  AND %s = $%d\n' "$column" "$idx"
    fi
    idx=$((idx + 1))
  done
}
