#!/usr/bin/env bash

contains_exact_column() {
  local needle="$1"
  shift
  local column
  for column in "$@"; do
    if [[ "$column" == "$needle" ]]; then
      return 0
    fi
  done
  return 1
}

is_searchable_column() {
  local column="$1"
  case "$column" in
    slug|name|description|value|title|code|email|login|label)
      return 0
      ;;
    *)
      return 1
      ;;
  esac
}

is_mutable_column() {
  local column="$1"

  case "$column" in
    id|created_at)
      return 1
      ;;
    *_id)
      return 1
      ;;
    *)
      return 0
      ;;
  esac
}

column_name_from_spec() {
  local spec="$1"
  printf '%s' "${spec%%|*}"
}

column_type_from_spec() {
  local spec="$1"
  local type="${spec#*|}"
  type=$(echo "$type" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
  printf '%s' "$type"
}

collect_table_metadata() {
  local table="$1"
  shift
  local specs=("$@")

  TABLE_NAME="$table"
  TABLE_SINGULAR_NAME="$(pascal_name "$table")"
  TABLE_PLURAL_NAME="$(pascal_plural_name "$table")"
  TABLE_COLUMNS=()
  TABLE_UPDATE_COLUMNS=()
  TABLE_FK_COLUMNS=()
  TABLE_SELECT_COLUMNS=()
  TABLE_TEXT_COLUMNS=()
  TABLE_HAS_ID=0
  TABLE_HAS_IS_REMOVED=0
  TABLE_HAS_CREATED_AT=0
  TABLE_HAS_UPDATED_AT=0

  local spec column
  for spec in "${specs[@]}"; do
    column="$(column_name_from_spec "$spec")"
    TABLE_COLUMNS+=("$column")
    TABLE_SELECT_COLUMNS+=("$column")

    if [[ "$column" == "id" ]]; then
      TABLE_HAS_ID=1
    fi
    if [[ "$column" == "is_removed" ]]; then
      TABLE_HAS_IS_REMOVED=1
    fi
    if [[ "$column" == "created_at" ]]; then
      TABLE_HAS_CREATED_AT=1
    fi
    if [[ "$column" == "updated_at" ]]; then
      TABLE_HAS_UPDATED_AT=1
    fi
    if is_searchable_column "$column"; then
      TABLE_TEXT_COLUMNS+=("$column")
    fi
    if is_mutable_column "$column"; then
      TABLE_UPDATE_COLUMNS+=("$column")
    fi
    if [[ "$column" == *_id && "$column" != "id" ]]; then
      TABLE_FK_COLUMNS+=("$column")
    fi
  done
}
