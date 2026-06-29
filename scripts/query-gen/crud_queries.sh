#!/usr/bin/env bash

emit_get_by_id_query() {
  echo "-- name: Get${TABLE_SINGULAR_NAME}ByID :one"
  echo "SELECT"
  indent_list "${TABLE_SELECT_COLUMNS[@]}"
  echo "FROM ${TABLE_NAME}"
  if [[ "$TABLE_HAS_IS_REMOVED" -eq 1 ]]; then
    echo "WHERE id = \$1"
    echo "  AND is_removed = FALSE;"
  else
    echo "WHERE id = \$1;"
  fi
  echo
}

emit_list_query() {
  echo "-- name: List${TABLE_PLURAL_NAME} :many"
  echo "SELECT"
  indent_list "${TABLE_SELECT_COLUMNS[@]}"
  echo "FROM ${TABLE_NAME}"
  if [[ "$TABLE_HAS_IS_REMOVED" -eq 1 ]]; then
    echo "WHERE is_removed = FALSE"
  fi
  if [[ "$TABLE_HAS_CREATED_AT" -eq 1 ]]; then
    echo "ORDER BY created_at DESC"
  else
    echo "ORDER BY id ASC"
  fi
  echo "LIMIT \$1 OFFSET \$2;"
  echo
}

emit_count_query() {
  echo "-- name: Count${TABLE_PLURAL_NAME} :one"
  echo "SELECT COUNT(*) AS total"
  echo "FROM ${TABLE_NAME}"
  if [[ "$TABLE_HAS_IS_REMOVED" -eq 1 ]]; then
    echo "WHERE is_removed = FALSE;"
  else
    echo ";"
  fi
  echo
}

emit_search_queries() {
  if [[ "${#TABLE_TEXT_COLUMNS[@]}" -eq 0 ]]; then
    return
  fi

  echo "-- name: List${TABLE_PLURAL_NAME}WithSearch :many"
  echo "SELECT"
  indent_list "${TABLE_SELECT_COLUMNS[@]}"
  echo "FROM ${TABLE_NAME}"
  
  # Если есть FK, добавляем фильтр
  if [[ "${#TABLE_FK_COLUMNS[@]}" -gt 0 ]]; then
    local fk="${TABLE_FK_COLUMNS[0]}"
    echo "WHERE ${fk} = sqlc.arg(${fk})"
    echo "  AND ("
  else
    echo "WHERE ("
  fi
  
  echo "    sqlc.narg(search)::text = '' OR"

  local idx=0
  local column
  for column in "${TABLE_TEXT_COLUMNS[@]}"; do
    if [[ $idx -gt 0 ]]; then
      echo "    OR"
    fi
    echo "    ${column} ILIKE '%' || sqlc.narg(search)::text || '%'"
    idx=$((idx + 1))
  done

  echo ")"
  if [[ "$TABLE_HAS_IS_REMOVED" -eq 1 ]]; then
    echo "  AND is_removed = FALSE"
  fi
  if [[ "$TABLE_HAS_CREATED_AT" -eq 1 ]]; then
    echo "ORDER BY"
    echo "    CASE WHEN sqlc.arg(sort_asc)::bool THEN created_at END ASC,"
    echo "    CASE WHEN NOT sqlc.arg(sort_asc)::bool THEN created_at END DESC"
  fi
  echo "LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);"
  echo

  # Count with search
  echo "-- name: Count${TABLE_PLURAL_NAME}WithSearch :one"
  echo "SELECT COUNT(*) AS total"
  echo "FROM ${TABLE_NAME}"
  
  if [[ "${#TABLE_FK_COLUMNS[@]}" -gt 0 ]]; then
    local fk="${TABLE_FK_COLUMNS[0]}"
    echo "WHERE ${fk} = sqlc.arg(${fk})"
    echo "  AND ("
  else
    echo "WHERE ("
  fi
  
  echo "    sqlc.narg(search)::text = '' OR"

  for column in "${TABLE_TEXT_COLUMNS[@]}"; do
    echo "    ${column} ILIKE '%' || sqlc.narg(search)::text || '%' OR"
  done
  echo "    FALSE"
  echo ")"
  if [[ "$TABLE_HAS_IS_REMOVED" -eq 1 ]]; then
    echo "  AND is_removed = FALSE;"
  else
    echo ";"
  fi
  echo
}

emit_get_by_slug_query() {
  if ! printf '%s\n' "${TABLE_COLUMNS[@]}" | grep -qx 'slug'; then
    return
  fi

  if [[ "${#TABLE_FK_COLUMNS[@]}" -gt 0 ]]; then
    local fk="${TABLE_FK_COLUMNS[0]}"
    local fk_name
    fk_name="$(column_param_name "$fk")"
    echo "-- name: Get${TABLE_SINGULAR_NAME}BySlug :one"
    echo "SELECT"
    indent_list "${TABLE_SELECT_COLUMNS[@]}"
    echo "FROM ${TABLE_NAME}"
    where_columns "$fk" "slug"
    if [[ "$TABLE_HAS_IS_REMOVED" -eq 1 ]]; then
      echo "  AND is_removed = FALSE"
    fi
    echo ";"
    echo
    return
  fi

  echo "-- name: Get${TABLE_SINGULAR_NAME}BySlug :one"
  echo "SELECT"
  indent_list "${TABLE_SELECT_COLUMNS[@]}"
  echo "FROM ${TABLE_NAME}"
  echo "WHERE slug = \$1"
  if [[ "$TABLE_HAS_IS_REMOVED" -eq 1 ]]; then
    echo "  AND is_removed = FALSE"
  fi
  echo ";"
  echo
}

emit_get_by_name_query() {
  if ! printf '%s\n' "${TABLE_COLUMNS[@]}" | grep -qx 'name'; then
    return
  fi

  echo "-- name: Get${TABLE_SINGULAR_NAME}ByName :one"
  echo "SELECT"
  indent_list "${TABLE_SELECT_COLUMNS[@]}"
  echo "FROM ${TABLE_NAME}"
  
  # Если есть FK, добавляем фильтр
  if [[ "${#TABLE_FK_COLUMNS[@]}" -gt 0 ]]; then
    local fk="${TABLE_FK_COLUMNS[0]}"
    where_columns "$fk" "name"
  else
    echo "WHERE name = \$1"
  fi
  
  if [[ "$TABLE_HAS_IS_REMOVED" -eq 1 ]]; then
    echo "  AND is_removed = FALSE"
  fi
  echo ";"
  echo
}

emit_get_by_public_key_query() {
  if ! printf '%s\n' "${TABLE_COLUMNS[@]}" | grep -qx 'public_key'; then
    return
  fi

  echo "-- name: Get${TABLE_SINGULAR_NAME}ByPublicKey :one"
  echo "SELECT"
  indent_list "${TABLE_SELECT_COLUMNS[@]}"
  echo "FROM ${TABLE_NAME}"
  echo "WHERE public_key = \$1"
  if [[ "$TABLE_HAS_IS_REMOVED" -eq 1 ]]; then
    echo "  AND is_removed = FALSE"
  fi
  echo ";"
  echo
}

emit_fk_queries() {
  if [[ "${#TABLE_FK_COLUMNS[@]}" -eq 0 ]]; then
    return
  fi

  local fk
  for fk in "${TABLE_FK_COLUMNS[@]}"; do
    local fk_name
    fk_name="$(column_param_name "$fk")"

    echo "-- name: List${TABLE_PLURAL_NAME}By${fk_name} :many"
    echo "SELECT"
    indent_list "${TABLE_SELECT_COLUMNS[@]}"
    echo "FROM ${TABLE_NAME}"
    echo "WHERE ${fk} = \$1"
    if [[ "$TABLE_HAS_IS_REMOVED" -eq 1 ]]; then
      echo "  AND is_removed = FALSE"
    fi
    if [[ "$TABLE_HAS_CREATED_AT" -eq 1 ]]; then
      echo "ORDER BY created_at DESC"
    else
      echo "ORDER BY id ASC"
    fi
    echo "LIMIT \$2 OFFSET \$3;"
    echo

    echo "-- name: Count${TABLE_PLURAL_NAME}By${fk_name} :one"
    echo "SELECT COUNT(*) AS total"
    echo "FROM ${TABLE_NAME}"
    echo "WHERE ${fk} = \$1"
    if [[ "$TABLE_HAS_IS_REMOVED" -eq 1 ]]; then
      echo "  AND is_removed = FALSE"
    fi
    echo ";"
    echo

    if [[ "$TABLE_HAS_CREATED_AT" -eq 1 ]]; then
      echo "-- name: List${TABLE_PLURAL_NAME}By${fk_name}Asc :many"
      echo "SELECT"
      indent_list "${TABLE_SELECT_COLUMNS[@]}"
      echo "FROM ${TABLE_NAME}"
      echo "WHERE ${fk} = \$1"
      if [[ "$TABLE_HAS_IS_REMOVED" -eq 1 ]]; then
        echo "  AND is_removed = FALSE"
      fi
      echo "ORDER BY created_at ASC"
      echo "LIMIT \$2 OFFSET \$3;"
      echo
    fi

    echo "-- name: List${TABLE_PLURAL_NAME}By${fk_name}Array :many"
    echo "SELECT"
    indent_list "${TABLE_SELECT_COLUMNS[@]}"
    echo "FROM ${TABLE_NAME}"
    echo "WHERE ${fk} = ANY(\$1::uuid[])"
    if [[ "$TABLE_HAS_IS_REMOVED" -eq 1 ]]; then
      echo "  AND is_removed = FALSE"
    fi
    if [[ "$TABLE_HAS_CREATED_AT" -eq 1 ]]; then
      echo "ORDER BY created_at ASC"
    else
      echo "ORDER BY id ASC"
    fi
    echo "LIMIT \$2 OFFSET \$3;"
    echo

    echo "-- name: Count${TABLE_PLURAL_NAME}By${fk_name}Array :one"
    echo "SELECT COUNT(*) AS total"
    echo "FROM ${TABLE_NAME}"
    echo "WHERE ${fk} = ANY(\$1::uuid[])"
    if [[ "$TABLE_HAS_IS_REMOVED" -eq 1 ]]; then
      echo "  AND is_removed = FALSE"
    fi
    echo ";"
    echo
  done
}

emit_update_query() {
  if [[ "${#TABLE_UPDATE_COLUMNS[@]}" -eq 0 ]]; then
    return
  fi

  echo "-- name: Update${TABLE_SINGULAR_NAME} :one"
  echo "UPDATE ${TABLE_NAME}"
  echo "SET"
  set_list 2 "${TABLE_UPDATE_COLUMNS[@]}"
  echo "WHERE id = \$1"
  if [[ "$TABLE_HAS_IS_REMOVED" -eq 1 ]]; then
    echo "  AND is_removed = FALSE"
  fi
  echo "RETURNING"
  indent_list "${TABLE_SELECT_COLUMNS[@]}"
  echo ";"
  echo
}

emit_delete_query() {
  echo "-- name: Delete${TABLE_SINGULAR_NAME} :execrows"
  if [[ "$TABLE_HAS_IS_REMOVED" -eq 1 ]]; then
    echo "UPDATE ${TABLE_NAME}"
    echo "SET"
    if [[ "$TABLE_HAS_UPDATED_AT" -eq 1 ]]; then
      echo "    is_removed = TRUE,"
      echo "    updated_at = NOW()"
    else
      echo "    is_removed = TRUE"
    fi
    echo "WHERE id = \$1"
    echo "  AND is_removed = FALSE;"
    return
  fi

  echo "DELETE FROM ${TABLE_NAME} WHERE id = \$1;"
}

generate_table_queries() {
  local out_file="$1"
  local table="$2"
  shift 2
  local specs=("$@")

  collect_table_metadata "$table" "${specs[@]}"

  {
    echo "-- Code generated by sqlc. DO NOT EDIT."
    echo
    echo "-- name: Create${TABLE_SINGULAR_NAME} :one"
    echo "INSERT INTO ${TABLE_NAME} ("
    indent_list "${TABLE_COLUMNS[@]}"
    echo ") VALUES ("
    placeholder_list "${#TABLE_COLUMNS[@]}"
    echo ") RETURNING"
    indent_list "${TABLE_SELECT_COLUMNS[@]}"
    echo ";"
    echo

    if [[ "$TABLE_HAS_ID" -eq 1 ]]; then
      emit_get_by_id_query
      emit_list_query
      emit_count_query
      emit_search_queries
      emit_get_by_slug_query
      emit_get_by_name_query
      emit_get_by_public_key_query
      emit_fk_queries
      emit_update_query
      emit_delete_query
    fi
  } > "$out_file"
}

generate_batch_queries() {
  :
}
