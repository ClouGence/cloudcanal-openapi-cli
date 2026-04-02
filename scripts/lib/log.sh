#!/usr/bin/env bash

SCRIPT_START_TS="$(date +%s)"

_log() {
  local level="$1"
  local color="$2"
  local stream="${3:-stdout}"
  shift 3

  local now
  now="$(date '+%Y-%m-%d %H:%M:%S')"
  local line
  line="${now} [${level}] $*"
  local label
  label="[${level}]"

  if [[ "$stream" == "stderr" ]]; then
    if [[ -t 2 ]]; then
      printf '%s \033[%sm%s\033[0m %s\n' "$now" "$color" "$label" "$*" >&2
    else
      printf '%s\n' "$line" >&2
    fi
    return
  fi

  if [[ -t 1 ]]; then
    printf '%s \033[%sm%s\033[0m %s\n' "$now" "$color" "$label" "$*"
    return
  fi

  printf '%s\n' "$line"
}

log_info()    { _log "INFO" "32" "stdout" "$@"; }
log_warn()    { _log "WARN" "33" "stdout" "$@"; }
log_success() { log_info "$@"; }
log_error()   { _log "ERROR" "31" "stderr" "$@"; }

print_run_summary() {
  local message="$1"
  local elapsed end_at
  elapsed="$(( $(date +%s) - SCRIPT_START_TS ))"
  end_at="$(date '+%Y-%m-%d %H:%M:%S %Z')"

  log_info "$message"
  log_info "Elapsed: ${elapsed}s"
  log_info "Completed at: ${end_at}"
}
