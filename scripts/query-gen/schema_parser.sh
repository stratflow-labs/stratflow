#!/usr/bin/env bash


extract_tables() {
  local schema_file="$1"
  
  awk '
    function trim(s) {
      gsub(/^[[:space:]]+/, "", s)
      gsub(/[[:space:]]+$/, "", s)
      return s
    }
    
    function is_type_token(token) {
      return token ~ /^[A-Z][A-Z0-9_]*(\([0-9,]+\))?(\[\])?$/
    }

    function emit_col(line, parts, col, type) {
      line = trim(line)
      sub(/,$/, "", line)
      if (line == "") return
      if (line ~ /^(CONSTRAINT|PRIMARY|FOREIGN|UNIQUE|CHECK|EXCLUDE)[[:space:]]/) return
      split(line, parts, /[[:space:]]+/)
      col = parts[1]
      type = parts[2]
      gsub(/"/, "", col)
      if (col == "" || col == ")" || col == "AND" || col == "OR") return
      if (!is_type_token(type)) return
      if (col != "" && type != "") print table "|" col "|" type
    }
    
    BEGIN { in_table = 0; pending_alter_table = "" }
    
    /^[[:space:]]*CREATE[[:space:]]+TABLE/ {
      line = $0
      sub(/.*CREATE[[:space:]]+TABLE[[:space:]]+(IF[[:space:]]+NOT[[:space:]]+EXISTS[[:space:]]+)?/, "", line)
      sub(/[[:space:]]*\(.*/, "", line)
      gsub(/"/, "", line)
      table = trim(line)
      in_table = 1
      next
    }
    
    in_table && /^[[:space:]]*\);/ {
      in_table = 0
      next
    }
    
    in_table {
      emit_col($0)
    }
    
    /^[[:space:]]*ALTER[[:space:]]+TABLE[[:space:]]+/ && /[[:space:]]+ADD[[:space:]]+COLUMN[[:space:]]+/ {
      line = $0
      sub(/.*ALTER[[:space:]]+TABLE[[:space:]]+/, "", line)
      alter_table = line
      sub(/[[:space:]]+ADD[[:space:]]+COLUMN.*/, "", alter_table)
      gsub(/"/, "", alter_table)
      alter_table = trim(alter_table)

      alter_col = line
      sub(/.*ADD[[:space:]]+COLUMN[[:space:]]+(IF[[:space:]]+NOT[[:space:]]+EXISTS[[:space:]]+)?/, "", alter_col)
      split(trim(alter_col), parts, /[[:space:]]+/)
      col = parts[1]
      gsub(/"/, "", col)
      type = parts[2]
      if (col != "" && type != "" && is_type_token(type)) print alter_table "|" col "|" type
      next
    }
    
    /^[[:space:]]*ALTER[[:space:]]+TABLE[[:space:]]+/ {
      line = $0
      sub(/.*ALTER[[:space:]]+TABLE[[:space:]]+/, "", line)
      gsub(/"/, "", line)
      pending_alter_table = trim(line)
      next
    }
    
    pending_alter_table != "" && /^[[:space:]]*ADD[[:space:]]+COLUMN[[:space:]]+/ {
      line = $0
      sub(/^[[:space:]]*ADD[[:space:]]+COLUMN[[:space:]]+(IF[[:space:]]+NOT[[:space:]]+EXISTS[[:space:]]+)?/, "", line)
      split(trim(line), parts, /[[:space:]]+/)
      col = parts[1]
      gsub(/"/, "", col)
      type = parts[2]
      if (col != "" && type != "" && is_type_token(type)) print pending_alter_table "|" col "|" type
      pending_alter_table = ""
      next
    }
    
    pending_alter_table != "" && /;/ {
      pending_alter_table = ""
    }
  ' "$schema_file"
}

contains_column() {
  local needle="$1"
  shift
  local item
  for item in "$@"; do
    if [[ "$item" == "$needle" ]]; then
      return 0
    fi
  done
  return 1
}
