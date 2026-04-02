#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
source "$SCRIPT_DIR/lib/log.sh"

DIST_DIR="${DIST_DIR:-$ROOT_DIR/dist}"
VERSION="${VERSION:-}"
COMMIT="${COMMIT:-$(git -C "$ROOT_DIR" rev-parse HEAD 2>/dev/null || printf 'unknown\n')}"
BUILD_TIME="${BUILD_TIME:-$(date -u +"%Y-%m-%dT%H:%M:%SZ")}"

if [[ -z "$VERSION" ]]; then
  if [[ "${GITHUB_REF_NAME:-}" == v* ]]; then
    VERSION="${GITHUB_REF_NAME#v}"
  else
    VERSION="dev"
  fi
fi

checksum_file() {
  local file_path="$1"
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$file_path"
    return 0
  fi
  if command -v shasum >/dev/null 2>&1; then
    shasum -a 256 "$file_path"
    return 0
  fi

  log_error "A SHA-256 tool is required (sha256sum or shasum)"
  exit 1
}

generate_release_notes() {
  local output="$DIST_DIR/release-notes.md"

  awk -v version="$VERSION" '
    $0 ~ "^## \\[" version "\\]" { in_section = 1 }
    /^## \[/ && in_section && $0 !~ "^## \\[" version "\\]" { exit }
    in_section { print }
  ' "$ROOT_DIR/CHANGELOG.md" > "$output"

  if [[ -s "$output" ]]; then
    return 0
  fi

  if [[ "$VERSION" == "dev" ]]; then
    cat > "$output" <<EOF
# Release Notes

- Development build
- commit: $COMMIT
- buildTime: $BUILD_TIME
EOF
    return 0
  fi

  log_error "Release notes for version ${VERSION} not found in CHANGELOG.md"
  exit 1
}

build_asset() {
  local goos="$1"
  local goarch="$2"
  local asset_name="cloudcanal_${goos}_${goarch}.tar.gz"
  local workdir
  workdir="$(mktemp -d)"

  log_info "Building ${goos}/${goarch}"
  (
    cd "$ROOT_DIR"
    CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" \
      make build \
      BIN="$workdir/cloudcanal" \
      GO_BUILD_FLAGS="-trimpath" \
      EXTRA_LDFLAGS="-s -w" \
      VERSION="$VERSION" \
      COMMIT="$COMMIT" \
      BUILD_TIME="$BUILD_TIME"
  )
  tar -C "$workdir" -czf "$DIST_DIR/$asset_name" cloudcanal
  rm -rf "$workdir"
}

log_info "CloudCanal OpenAPI CLI release asset build started"
rm -rf "$DIST_DIR"
mkdir -p "$DIST_DIR"

for target in \
  "darwin amd64" \
  "darwin arm64" \
  "linux amd64" \
  "linux arm64"
do
  read -r goos goarch <<<"$target"
  build_asset "$goos" "$goarch"
done

cp "$ROOT_DIR/scripts/install.sh" "$DIST_DIR/install.sh"
cp "$ROOT_DIR/scripts/uninstall.sh" "$DIST_DIR/uninstall.sh"
chmod +x "$DIST_DIR/install.sh" "$DIST_DIR/uninstall.sh"

generate_release_notes

(
  cd "$DIST_DIR"
  : > checksums.txt
  for asset in cloudcanal_*.tar.gz install.sh uninstall.sh release-notes.md; do
    checksum_file "$asset" >> checksums.txt
  done
)

log_success "Release assets ready in $DIST_DIR"
print_run_summary "Release asset build completed"
