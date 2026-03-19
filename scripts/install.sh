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
REPO_OWNER="${REPO_OWNER:-Arlowen}"
REPO_NAME="${REPO_NAME:-cloudcanal-openapi-cli}"
RELEASE_VERSION="${RELEASE_VERSION:-latest}"
INSTALL_ROOT="${INSTALL_ROOT:-$HOME/.cloudcanal-cli}"
INSTALL_BIN_DIR="${INSTALL_BIN_DIR:-$INSTALL_ROOT/bin}"
INSTALL_PATH="$INSTALL_BIN_DIR/$APP_NAME"
INSTALL_BIN_PATH="$INSTALL_PATH"
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
DOWNLOAD_BASE_URL="${DOWNLOAD_BASE_URL:-}"
TMP_DIR="$(mktemp -d "${TMPDIR:-/tmp}/${REPO_NAME}.XXXXXX")"
ARCHIVE_PATH="$TMP_DIR/$APP_NAME.tar.gz"
CHECKSUMS_PATH="$TMP_DIR/checksums.txt"
EXTRACT_DIR="$TMP_DIR/extract"

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

detect_platform() {
  local os arch

  case "$(uname -s)" in
    Darwin) os="darwin" ;;
    Linux) os="linux" ;;
    *)
      log_error "Unsupported operating system: $(uname -s)"
      exit 1
      ;;
  esac

  case "$(uname -m)" in
    x86_64|amd64) arch="amd64" ;;
    arm64|aarch64) arch="arm64" ;;
    *)
      log_error "Unsupported CPU architecture: $(uname -m)"
      exit 1
      ;;
  esac

  PLATFORM_OS="$os"
  PLATFORM_ARCH="$arch"
  ARCHIVE_NAME="${APP_NAME}_${PLATFORM_OS}_${PLATFORM_ARCH}.tar.gz"
}

download_url() {
  if [[ -n "$DOWNLOAD_BASE_URL" ]]; then
    printf '%s/%s\n' "${DOWNLOAD_BASE_URL%/}" "$ARCHIVE_NAME"
    return 0
  fi

  if [[ "$RELEASE_VERSION" == "latest" ]]; then
    printf 'https://github.com/%s/%s/releases/latest/download/%s\n' "$REPO_OWNER" "$REPO_NAME" "$ARCHIVE_NAME"
    return 0
  fi

  printf 'https://github.com/%s/%s/releases/download/%s/%s\n' "$REPO_OWNER" "$REPO_NAME" "$RELEASE_VERSION" "$ARCHIVE_NAME"
}

checksums_url() {
  if [[ -n "$DOWNLOAD_BASE_URL" ]]; then
    printf '%s/checksums.txt\n' "${DOWNLOAD_BASE_URL%/}"
    return 0
  fi

  if [[ "$RELEASE_VERSION" == "latest" ]]; then
    printf 'https://github.com/%s/%s/releases/latest/download/checksums.txt\n' "$REPO_OWNER" "$REPO_NAME"
    return 0
  fi

  printf 'https://github.com/%s/%s/releases/download/%s/checksums.txt\n' "$REPO_OWNER" "$REPO_NAME" "$RELEASE_VERSION"
}

sha256_command() {
  if command -v sha256sum >/dev/null 2>&1; then
    printf 'sha256sum\n'
    return 0
  fi
  if command -v shasum >/dev/null 2>&1; then
    printf 'shasum\n'
    return 0
  fi
  if command -v openssl >/dev/null 2>&1; then
    printf 'openssl\n'
    return 0
  fi
  log_error "A SHA-256 tool is required (sha256sum, shasum, or openssl)"
  exit 1
}

compute_sha256() {
  local file_path="$1"

  case "$(sha256_command)" in
    sha256sum)
      sha256sum "$file_path" | awk '{print $1}'
      ;;
    shasum)
      shasum -a 256 "$file_path" | awk '{print $1}'
      ;;
    openssl)
      openssl dgst -sha256 -r "$file_path" | awk '{print $1}'
      ;;
  esac
}

verify_archive_checksum() {
  local url
  url="$(checksums_url)"
  log_info "Downloading checksum file from $url"
  curl -fsSL "$url" -o "$CHECKSUMS_PATH"

  local expected actual
  expected="$(
    awk -v target="$ARCHIVE_NAME" '
      {
        name = $2
        sub(/^\.\//, "", name)
        if (name == target) {
          print $1
          exit
        }
      }
    ' "$CHECKSUMS_PATH"
  )"
  if [[ -z "$expected" ]]; then
    log_error "Checksum entry not found for $ARCHIVE_NAME"
    exit 1
  fi

  actual="$(compute_sha256 "$ARCHIVE_PATH")"
  if [[ "$actual" != "$expected" ]]; then
    log_error "Checksum verification failed for $ARCHIVE_NAME"
    log_error "Expected: $expected"
    log_error "Actual:   $actual"
    exit 1
  fi

  log_success "Verified release checksum for $ARCHIVE_NAME"
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

install_binary() {
  mkdir -p "$INSTALL_BIN_DIR" "$EXTRACT_DIR"

  local url
  url="$(download_url)"
  log_info "Downloading release asset from $url"
  curl -fsSL "$url" -o "$ARCHIVE_PATH"
  verify_archive_checksum

  log_info "Extracting release asset"
  tar -xzf "$ARCHIVE_PATH" -C "$EXTRACT_DIR"

  local extracted_binary="$EXTRACT_DIR/$APP_NAME"
  if [[ ! -x "$extracted_binary" ]]; then
    log_error "Release asset does not contain executable $APP_NAME"
    exit 1
  fi

  install -m 755 "$extracted_binary" "$INSTALL_PATH"
  log_success "Installed binary to $INSTALL_PATH"
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

  "$INSTALL_BIN_PATH" completion zsh "$APP_NAME" > "$ZSH_COMPLETION_PATH"
  log_success "Installed zsh completion to $ZSH_COMPLETION_PATH"

  "$INSTALL_BIN_PATH" completion bash "$APP_NAME" > "$BASH_COMPLETION_PATH"
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

trap cleanup EXIT

log_info "CloudCanal OpenAPI CLI release install started"
require_command curl
require_command tar
detect_platform
install_binary
ensure_completion_files
ensure_path_block
ensure_completion_block

log_info "Open a new shell or source $INSTALL_SHELL_RC, then run: $APP_NAME jobs list"
log_success "Release install completed"
