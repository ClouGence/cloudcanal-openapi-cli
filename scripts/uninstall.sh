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

default_shell_rc() {
  case "$(basename "${SHELL:-}")" in
    bash) printf '%s\n' "$HOME/.bashrc" ;;
    *) printf '%s\n' "$HOME/.zshrc" ;;
  esac
}

APP_NAME="${APP_NAME:-cloudcanal}"
REPO_NAME="${REPO_NAME:-cloudcanal-openapi-cli}"
INSTALL_ROOT="${INSTALL_ROOT:-$HOME/.cloudcanal-cli}"
INSTALL_BIN_DIR="${INSTALL_BIN_DIR:-$INSTALL_ROOT/bin}"
INSTALL_PATH="$INSTALL_BIN_DIR/$APP_NAME"
INSTALL_SHELL_RC="${INSTALL_SHELL_RC:-$(default_shell_rc)}"
INSTALL_COMPLETION_DIR="${INSTALL_COMPLETION_DIR:-$INSTALL_ROOT/completions}"
INSTALL_ZSH_COMPLETION_DIR="${INSTALL_ZSH_COMPLETION_DIR:-$INSTALL_COMPLETION_DIR/zsh}"
INSTALL_BASH_COMPLETION_DIR="${INSTALL_BASH_COMPLETION_DIR:-$INSTALL_COMPLETION_DIR/bash}"
ZSH_COMPLETION_PATH="$INSTALL_ZSH_COMPLETION_DIR/_$APP_NAME"
BASH_COMPLETION_PATH="$INSTALL_BASH_COMPLETION_DIR/$APP_NAME"
PATH_MARK_START="# >>> cloudcanal-openapi-cli >>>"
PATH_MARK_END="# <<< cloudcanal-openapi-cli <<<"
COMPLETION_MARK_START="# >>> cloudcanal-openapi-cli completion >>>"
COMPLETION_MARK_END="# <<< cloudcanal-openapi-cli completion <<<"

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

remove_if_empty_dir() {
  local dir="$1"
  if [[ -d "$dir" ]] && [[ -z "$(ls -A "$dir")" ]]; then
    rmdir "$dir"
  fi
}

prune_install_dirs() {
  remove_if_empty_dir "$INSTALL_ZSH_COMPLETION_DIR"
  remove_if_empty_dir "$INSTALL_BASH_COMPLETION_DIR"
  remove_if_empty_dir "$INSTALL_COMPLETION_DIR"
  remove_if_empty_dir "$INSTALL_BIN_DIR"
  remove_if_empty_dir "$INSTALL_ROOT"
}

remove_binary() {
  if [[ -e "$INSTALL_PATH" ]]; then
    rm -f "$INSTALL_PATH"
    log_success "Removed $INSTALL_PATH"
  else
    log_info "No installed command found at $INSTALL_PATH"
  fi
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

log_info "CloudCanal OpenAPI CLI release uninstall started"
remove_binary
remove_rc_block "$PATH_MARK_START" "$PATH_MARK_END" "PATH configuration"
remove_completion_files
remove_rc_block "$COMPLETION_MARK_START" "$COMPLETION_MARK_END" "shell completion configuration"
prune_install_dirs
log_info "Open a new shell or source $INSTALL_SHELL_RC to refresh PATH"
log_success "Release uninstall completed"
