#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
source "$SCRIPT_DIR/lib/log.sh"

BIN_DIR="$ROOT_DIR/bin"
BIN_PATH="$BIN_DIR/cloudcanal"
CLI_HOME_DIR="${HOME}/.cloudcanal-cli"
LOG_ROOT_DIR="${CLI_HOME_DIR}/logs"

mkdir -p "$LOG_ROOT_DIR"
LOG_DIR="$(mktemp -d "${LOG_ROOT_DIR}/cloudcanal-openapi-cli-build.XXXXXX")"

cleanup() {
  rm -rf "$LOG_DIR"
}

trap cleanup EXIT

run_step() {
  local title="$1"
  local log_name="$2"
  shift 2

  local log_path="$LOG_DIR/$log_name.log"
  log_info "[$STEP_NO] $title"

  if "$@" >"$log_path" 2>&1; then
    log_success "$title completed"
    return 0
  fi

  log_error "$title failed"
  log_error "Command output:"
  sed 's/^/    /' "$log_path" >&2
  exit 1
}

cd "$ROOT_DIR"

log_info "CloudCanal OpenAPI CLI build started"

STEP_NO="1/3"
log_info "[$STEP_NO] Clean build artifacts"
if [[ -d "$BIN_DIR" ]]; then
  rm -rf "$BIN_DIR"
  log_success "Removed $BIN_DIR"
else
  log_success "No existing build artifacts"
fi

STEP_NO="2/3"
run_step "Run tests" "test" make test

STEP_NO="3/3"
run_step "Build CLI" "build" make build

log_success "Binary ready at $BIN_PATH"
print_run_summary "Build completed"
