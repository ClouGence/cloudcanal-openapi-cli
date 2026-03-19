#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
source "$SCRIPT_DIR/lib/log.sh"

APP_NAME="${APP_NAME:-cloudcanal}"
BIN_PATH="$ROOT_DIR/bin/$APP_NAME"
INSTALL_BIN_DIR="${INSTALL_BIN_DIR:-$HOME/bin}"
INSTALL_PATH="$INSTALL_BIN_DIR/$APP_NAME"
INSTALL_SHELL_RC="${INSTALL_SHELL_RC:-$HOME/.zshrc}"
PATH_MARK_START="# >>> cloudcanal-openapi-cli >>>"
PATH_MARK_END="# <<< cloudcanal-openapi-cli <<<"

remove_link() {
  if [[ -L "$INSTALL_PATH" ]]; then
    local target
    target="$(readlink "$INSTALL_PATH")"
    if [[ "$target" == "$BIN_PATH" ]]; then
      rm -f "$INSTALL_PATH"
      log_success "Removed $INSTALL_PATH"
      return 0
    fi
    log_info "Skipped $INSTALL_PATH because it is not managed by this project"
    return 0
  fi

  if [[ -e "$INSTALL_PATH" ]]; then
    log_info "Skipped $INSTALL_PATH because it is not a symlink created by this project"
    return 0
  fi

  log_info "No installed command found at $INSTALL_PATH"
}

remove_path_block() {
  if [[ ! -f "$INSTALL_SHELL_RC" ]] || ! grep -Fq "$PATH_MARK_START" "$INSTALL_SHELL_RC"; then
    log_info "No PATH configuration to remove from $INSTALL_SHELL_RC"
    return 0
  fi

  local tmp_file
  tmp_file="$(mktemp)"

  awk -v start="$PATH_MARK_START" -v end="$PATH_MARK_END" '
    $0 == start {skip = 1; next}
    $0 == end {skip = 0; next}
    !skip {print}
  ' "$INSTALL_SHELL_RC" > "$tmp_file"

  mv "$tmp_file" "$INSTALL_SHELL_RC"
  log_success "Updated $INSTALL_SHELL_RC"
}

log_info "Uninstall $APP_NAME command"
remove_link
remove_path_block

log_info "Open a new shell or source $INSTALL_SHELL_RC to refresh PATH"
print_run_summary "Uninstall completed"
