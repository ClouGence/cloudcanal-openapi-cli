#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
source "$SCRIPT_DIR/lib/log.sh"

default_shell_rc() {
  case "$(basename "${SHELL:-}")" in
    bash) printf '%s\n' "$HOME/.bashrc" ;;
    *) printf '%s\n' "$HOME/.zshrc" ;;
  esac
}

APP_NAME="${APP_NAME:-cloudcanal}"
BIN_PATH="$ROOT_DIR/bin/$APP_NAME"
CLOUDCANAL_AUTO_RELOAD_SHELL="${CLOUDCANAL_AUTO_RELOAD_SHELL:-1}"
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

ensure_binary() {
  if [[ -x "$BIN_PATH" ]]; then
    log_success "Found binary at $BIN_PATH"
    return 0
  fi

  log_info "Binary not found, running all_build.sh first"
  "$SCRIPT_DIR/all_build.sh"
}

strip_rc_block() {
  local start_mark="$1"
  local end_mark="$2"
  local tmp_file
  tmp_file="$(mktemp)"

  if [[ -f "$INSTALL_SHELL_RC" ]]; then
    awk -v start="$start_mark" -v end="$end_mark" '
      $0 == start {skip = 1; next}
      $0 == end {skip = 0; next}
      !skip {print}
    ' "$INSTALL_SHELL_RC" > "$tmp_file"
  fi

  mv "$tmp_file" "$INSTALL_SHELL_RC"
}

ensure_path_block() {
  mkdir -p "$(dirname "$INSTALL_SHELL_RC")"
  touch "$INSTALL_SHELL_RC"
  strip_rc_block "$PATH_MARK_START" "$PATH_MARK_END"

  {
    printf '\n%s\n' "$PATH_MARK_START"
    printf 'export PATH="%s:$PATH"\n' "$INSTALL_BIN_DIR"
    printf '%s\n' "$PATH_MARK_END"
  } >> "$INSTALL_SHELL_RC"

  log_success "Updated $INSTALL_SHELL_RC"
}

ensure_completion_files() {
  mkdir -p "$INSTALL_ZSH_COMPLETION_DIR" "$INSTALL_BASH_COMPLETION_DIR"

  "$INSTALL_PATH" completion zsh "$APP_NAME" > "$ZSH_COMPLETION_PATH"
  log_success "Installed zsh completion to $ZSH_COMPLETION_PATH"

  "$INSTALL_PATH" completion bash "$APP_NAME" > "$BASH_COMPLETION_PATH"
  log_success "Installed bash completion to $BASH_COMPLETION_PATH"
}

ensure_completion_block() {
  mkdir -p "$(dirname "$INSTALL_SHELL_RC")"
  touch "$INSTALL_SHELL_RC"
  strip_rc_block "$COMPLETION_MARK_START" "$COMPLETION_MARK_END"

  {
    printf '\n%s\n' "$COMPLETION_MARK_START"
    printf 'if [[ -n "${ZSH_VERSION:-}" ]] && [[ -d "%s" ]]; then\n' "$INSTALL_ZSH_COMPLETION_DIR"
    printf '  fpath=("%s" $fpath)\n' "$INSTALL_ZSH_COMPLETION_DIR"
    printf '  autoload -Uz compinit\n'
    printf '  compinit\n'
    printf 'fi\n'
    printf 'if [[ -n "${BASH_VERSION:-}" ]] && [[ -f "%s" ]]; then\n' "$BASH_COMPLETION_PATH"
    printf '  source "%s"\n' "$BASH_COMPLETION_PATH"
    printf 'fi\n'
    printf '%s\n' "$COMPLETION_MARK_END"
  } >> "$INSTALL_SHELL_RC"

  log_success "Updated $INSTALL_SHELL_RC"
}

reload_shell_session() {
  if [[ "$CLOUDCANAL_AUTO_RELOAD_SHELL" != "1" ]] || [[ -n "${CI:-}" ]] || [[ ! -t 0 ]] || [[ ! -t 1 ]]; then
    log_info "Open a new shell or source $INSTALL_SHELL_RC, then run: $APP_NAME jobs list"
    return 0
  fi

  if [[ -z "${SHELL:-}" ]] || [[ ! -x "$SHELL" ]]; then
    log_info "Open a new shell or source $INSTALL_SHELL_RC, then run: $APP_NAME jobs list"
    return 0
  fi

  log_info "Launching a new login shell so $APP_NAME is available immediately"
  exec "$SHELL" -l
}

log_info "Install $APP_NAME command"
ensure_binary
mkdir -p "$INSTALL_BIN_DIR"
ln -sfn "$BIN_PATH" "$INSTALL_PATH"
log_success "Installed $INSTALL_PATH"
ensure_completion_files
ensure_path_block
ensure_completion_block

print_run_summary "Install completed"
reload_shell_session
