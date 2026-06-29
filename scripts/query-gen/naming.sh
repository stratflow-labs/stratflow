#!/usr/bin/env bash
# Утилиты для преобразования имен таблиц и колонок

camel_word() {
  local word="$1"
  printf '%s%s' "$(tr '[:lower:]' '[:upper:]' <<< "${word:0:1}")" "${word:1}"
}

pascalize_parts() {
  local raw="$1"
  local out=""
  local part

  IFS='_' read -r -a parts <<< "$raw"
  for part in "${parts[@]}"; do
    out+="$(camel_word "$part")"
  done

  printf '%s' "$out"
}

pascal_name() {
  local raw="$1"
  local singular="$raw"
  
  if [[ "$singular" == *_status || "$singular" == "status" ]]; then
    singular="$singular"
  elif [[ "$singular" == *_sessions ]]; then
    singular="${singular%_sessions}_session"
  elif [[ "$singular" == *ies ]]; then
    singular="${singular%ies}y"
  elif [[ "$singular" == *s && "$singular" != "session" ]]; then
    singular="${singular%s}"
  fi

  pascalize_parts "$singular"
}

pascal_plural_name() {
  local raw="$1"

  local plural="$raw"
  if [[ "$plural" == *y ]]; then
    plural="${plural%y}ies"
  elif [[ "$plural" != *s ]]; then
    plural="${plural}s"
  fi

  pascalize_parts "$plural"
}

column_param_name() {
  local column="$1"
  
  if [[ "$column" == "id" ]]; then
    printf 'ID'
    return
  fi
  
  if [[ "$column" == *_id ]]; then
    printf '%sID' "$(pascal_name "${column%_id}")"
    return
  fi
  
  pascal_name "$column"
}
