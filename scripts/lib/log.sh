#!/usr/bin/env bash

SCRIPT_START_TS="$(date +%s)"

_log() {
  local level="$1"
  local color="$2"
  shift 2

  local now
  now="$(date '+%Y-%m-%d %H:%M:%S')"

  if [[ -t 1 ]]; then
    printf '\033[%sm[%s]\033[0m %s %s\n' "$color" "$level" "$now" "$*"
  else
    printf '[%s] %s %s\n' "$level" "$now" "$*"
  fi
}

log_info()    { _log "INFO"    "34" "$@"; }
log_success() { _log "SUCCESS" "32" "$@"; }
log_error()   { _log "ERROR"   "31" "$@" >&2; }

print_run_summary() {
  local message="$1"
  local elapsed end_at
  elapsed="$(( $(date +%s) - SCRIPT_START_TS ))"
  end_at="$(date '+%Y-%m-%d %H:%M:%S %Z')"

  log_success "$message"
  log_info "Elapsed: ${elapsed}s"
  log_info "Completed at: ${end_at}"
}
