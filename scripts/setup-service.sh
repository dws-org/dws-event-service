#!/usr/bin/env bash

set -euo pipefail

usage() {
  cat <<'USAGE'
Usage: scripts/setup-service.sh <module-path> [service-name] [service-port]

Apply project-wide renames so this template is ready to be used as a dedicated
microservice repository.

Arguments:
  module-path   Fully qualified Go module path (e.g. github.com/acme/user-service)
  service-name  Human-readable service name (defaults to the module name)
  service-port  HTTP port without the leading colon (defaults to 8080)
USAGE
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  usage
  exit 0
fi

if [[ $# -lt 1 ]]; then
  echo "error: module-path is required"
  echo
  usage
  exit 1
fi

MODULE_PATH="$1"
SERVICE_NAME="${2:-$(basename "$MODULE_PATH")}"
SERVICE_PORT="${3:-8080}"

# Normalise the service slug
to_slug() {
  local input="$1"
  local lowercase
  lowercase=$(printf '%s' "$input" | tr '[:upper:]' '[:lower:]')
  local slug
  slug=$(printf '%s' "$lowercase" | sed -E 's/[^a-z0-9]+/-/g; s/^-+|-+$//g; s/-+/-/g')
  if [[ -z "$slug" ]]; then
    slug="service"
  fi
  printf '%s' "$slug"
}

SERVICE_SLUG="$(to_slug "$SERVICE_NAME")"
SERVICE_DESCRIPTION="${SERVICE_NAME} microservice API"

if [[ "$SERVICE_PORT" =~ ^: ]]; then
  PORT_VALUE="$SERVICE_PORT"
else
  PORT_VALUE=":$SERVICE_PORT"
fi

command -v go >/dev/null 2>&1 || {
  echo "error: go is required to run this script" >&2
  exit 1
}

command -v python3 >/dev/null 2>&1 || {
  echo "error: python3 is required to run this script" >&2
  exit 1
}

echo "Updating Go module path to ${MODULE_PATH}"
go mod edit -module "$MODULE_PATH"

echo "Replacing import paths"
python3 <<PY
from pathlib import Path

old_path = "boiler-go-api"
new_path = "${MODULE_PATH}"

for path in Path(".").rglob("*.go"):
    if any(part.startswith(".git") for part in path.parts):
        continue
    text = path.read_text()
    if old_path in text:
        path.write_text(text.replace(old_path, new_path))
PY

echo "Updating service metadata in configs/config.yaml"
python3 <<PY
from pathlib import Path
import re

config_path = Path("configs/config.yaml")
raw = config_path.read_text()

patterns = {
    r'name:\s*".*?"': f'name: "${SERVICE_NAME}"',
    r'slug:\s*".*?"': f'slug: "${SERVICE_SLUG}"',
    r'description:\s*".*?"': f'description: "${SERVICE_DESCRIPTION}"',
    r'version:\s*".*?"': 'version: "0.1.0"',
    r'port:\s*".*?"': f'port: "${PORT_VALUE}"',
}

for pattern, replacement in patterns.items():
    raw, count = re.subn(pattern, replacement, raw, count=1, flags=re.MULTILINE)
    if count == 0:
        print(f"warning: could not apply pattern '{pattern}' in {config_path}")

config_path.write_text(raw)
PY

echo "Updating CLI examples"
python3 <<PY
from pathlib import Path

target = Path("cmd/version.go")
text = target.read_text()
text = text.replace("service-template version", "${SERVICE_SLUG} version")
target.write_text(text)
PY

echo "Running go mod tidy"
go mod tidy

cat <<DONE

Setup complete!
- Module path: ${MODULE_PATH}
- Service name: ${SERVICE_NAME}
- Service slug: ${SERVICE_SLUG}
- HTTP port: ${PORT_VALUE}

Next steps:
  1. Review configs/config.yaml for any additional customisations.
  2. Update docs/readme to reflect this service's responsibilities.
  3. Commit the changes and push to a new repository.
DONE

