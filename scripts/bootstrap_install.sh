#!/usr/bin/env bash

set -euo pipefail

_log() {
  local level="$1"
  local color="$2"
  shift 2
  printf '\033[%sm[%s]\033[0m %s %s\n' "$color" "$level" "$(date '+%Y-%m-%d %H:%M:%S')" "$*"
}

log_info()    { _log "INFO" "34" "$@"; }
log_success() { _log "SUCCESS" "32" "$@"; }
log_error()   { _log "ERROR" "31" "$@" >&2; }

APP_NAME="${APP_NAME:-cloudcanal}"
REPO_OWNER="${REPO_OWNER:-Arlowen}"
REPO_NAME="${REPO_NAME:-cloudcanal-openapi-cli}"
REPO_REF="${REPO_REF:-main}"
ARCHIVE_URL="${ARCHIVE_URL:-https://github.com/$REPO_OWNER/$REPO_NAME/archive/refs/heads/$REPO_REF.tar.gz}"
INSTALL_ROOT="${INSTALL_ROOT:-$HOME/.local/share/$REPO_NAME}"
REPOSITORY_DIR="$INSTALL_ROOT/repository"
TMP_DIR="$(mktemp -d "${TMPDIR:-/tmp}/${REPO_NAME}.XXXXXX")"
ARCHIVE_PATH="$TMP_DIR/source.tar.gz"
EXTRACT_DIR="$TMP_DIR/extract"
STAGED_REPOSITORY_DIR="$INSTALL_ROOT/repository.new"

cleanup() {
  rm -rf "$TMP_DIR"
}

require_command() {
  local command_name="$1"
  if command -v "$command_name" >/dev/null 2>&1; then
    return 0
  fi
  log_error "Required command not found: $command_name"
  exit 1
}

trap cleanup EXIT

log_info "CloudCanal OpenAPI CLI bootstrap install started"
require_command curl
require_command tar
require_command go

mkdir -p "$INSTALL_ROOT" "$EXTRACT_DIR"

log_info "Downloading source archive from $ARCHIVE_URL"
curl -fsSL "$ARCHIVE_URL" -o "$ARCHIVE_PATH"

log_info "Extracting source archive"
tar -xzf "$ARCHIVE_PATH" -C "$EXTRACT_DIR"

SOURCE_DIR="$(find "$EXTRACT_DIR" -mindepth 1 -maxdepth 1 -type d | head -n 1)"
if [[ -z "${SOURCE_DIR:-}" ]]; then
  log_error "Failed to locate extracted source directory"
  exit 1
fi

log_info "Building project from source"
"$SOURCE_DIR/scripts/all_build.sh"

rm -rf "$STAGED_REPOSITORY_DIR"
mv "$SOURCE_DIR" "$STAGED_REPOSITORY_DIR"
rm -rf "$REPOSITORY_DIR"
mv "$STAGED_REPOSITORY_DIR" "$REPOSITORY_DIR"
log_success "Installed source repository at $REPOSITORY_DIR"

log_info "Running project install script"
"$REPOSITORY_DIR/scripts/install.sh"

log_success "Bootstrap install completed"
log_info "One-line install command:"
log_info "curl -fsSL https://raw.githubusercontent.com/$REPO_OWNER/$REPO_NAME/$REPO_REF/scripts/bootstrap_install.sh | bash"
