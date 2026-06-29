#!/usr/bin/env bash

QUERY_GEN_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

source "$QUERY_GEN_DIR/naming.sh"
source "$QUERY_GEN_DIR/rendering.sh"
source "$QUERY_GEN_DIR/column_specs.sh"
source "$QUERY_GEN_DIR/crud_queries.sh"
