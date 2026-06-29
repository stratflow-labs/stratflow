#!/usr/bin/env sh

set -eu

repo_root="$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)"

rm -rf "$repo_root/apps/web/src/shared/api/gen/identity" \
  "$repo_root/apps/web/src/shared/api/gen/strategy-registry" \
  "$repo_root/apps/web/src/shared/api/gen/strategy_registry"

(
  cd "$repo_root/services/identity/api"
  buf generate --template ./buf.frontend.gen.yaml
)

(
  cd "$repo_root/services/strategy-registry/api"
  buf generate --template ./buf.frontend.gen.yaml
)
