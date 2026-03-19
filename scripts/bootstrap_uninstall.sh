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
REPO_NAME="${REPO_NAME:-cloudcanal-openapi-cli}"
INSTALL_ROOT="${INSTALL_ROOT:-$HOME/.local/share/$REPO_NAME}"
INSTALL_BIN_DIR="${INSTALL_BIN_DIR:-$HOME/bin}"
INSTALL_PATH="$INSTALL_BIN_DIR/$APP_NAME"
INSTALL_BIN_PATH="$INSTALL_ROOT/bin/$APP_NAME"
INSTALL_SHELL_RC="${INSTALL_SHELL_RC:-$HOME/.zshrc}"
INSTALL_ZSH_COMPLETION_DIR="${INSTALL_ZSH_COMPLETION_DIR:-$HOME/.zsh/completions}"
INSTALL_BASH_COMPLETION_DIR="${INSTALL_BASH_COMPLETION_DIR:-$HOME/.local/share/bash-completion/completions}"
ZSH_COMPLETION_PATH="$INSTALL_ZSH_COMPLETION_DIR/_$APP_NAME"
BASH_COMPLETION_PATH="$INSTALL_BASH_COMPLETION_DIR/$APP_NAME"
PATH_MARK_START="# >>> cloudcanal-openapi-cli >>>"
PATH_MARK_END="# <<< cloudcanal-openapi-cli <<<"
COMPLETION_MARK_START="# >>> cloudcanal-openapi-cli completion >>>"
COMPLETION_MARK_END="# <<< cloudcanal-openapi-cli completion <<<"

remove_link() {
  if [[ -L "$INSTALL_PATH" ]]; then
    local target
    target="$(readlink "$INSTALL_PATH")"
    if [[ "$target" == "$INSTALL_BIN_PATH" ]]; then
      rm -f "$INSTALL_PATH"
      log_success "Removed $INSTALL_PATH"
      return 0
    fi
    log_info "Skipped $INSTALL_PATH because it is not managed by the release installer"
    return 0
  fi

  if [[ -e "$INSTALL_PATH" ]]; then
    log_info "Skipped $INSTALL_PATH because it is not a managed symlink"
    return 0
  fi

  log_info "No installed command found at $INSTALL_PATH"
}

remove_rc_block() {
  local start_mark="$1"
  local end_mark="$2"
  local description="$3"

  if [[ ! -f "$INSTALL_SHELL_RC" ]] || ! grep -Fq "$start_mark" "$INSTALL_SHELL_RC"; then
    log_info "No $description to remove from $INSTALL_SHELL_RC"
    return 0
  fi

  local tmp_file
  tmp_file="$(mktemp)"

  awk -v start="$start_mark" -v end="$end_mark" '
    $0 == start {skip = 1; next}
    $0 == end {skip = 0; next}
    !skip {print}
  ' "$INSTALL_SHELL_RC" > "$tmp_file"

  mv "$tmp_file" "$INSTALL_SHELL_RC"
  log_success "Updated $INSTALL_SHELL_RC"
}

remove_completion_files() {
  if [[ -f "$ZSH_COMPLETION_PATH" ]]; then
    rm -f "$ZSH_COMPLETION_PATH"
    log_success "Removed $ZSH_COMPLETION_PATH"
  else
    log_info "No zsh completion file found at $ZSH_COMPLETION_PATH"
  fi

  if [[ -f "$BASH_COMPLETION_PATH" ]]; then
    rm -f "$BASH_COMPLETION_PATH"
    log_success "Removed $BASH_COMPLETION_PATH"
  else
    log_info "No bash completion file found at $BASH_COMPLETION_PATH"
  fi
}

remove_install_root() {
  if [[ -d "$INSTALL_ROOT" ]]; then
    rm -rf "$INSTALL_ROOT"
    log_success "Removed $INSTALL_ROOT"
  else
    log_info "No install root found at $INSTALL_ROOT"
  fi
}

log_info "CloudCanal OpenAPI CLI release uninstall started"
remove_link
remove_rc_block "$PATH_MARK_START" "$PATH_MARK_END" "PATH configuration"
remove_completion_files
remove_rc_block "$COMPLETION_MARK_START" "$COMPLETION_MARK_END" "shell completion configuration"
remove_install_root
log_info "Open a new shell or source $INSTALL_SHELL_RC to refresh PATH"
log_success "Release uninstall completed"
