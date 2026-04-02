# Changelog

All notable changes to this project will be documented in this file.

## [0.1.1] - 2026-04-02

### Changed

- Migrated the canonical GitHub repository owner to `ClouGence`.
- Updated the Go module path and internal imports to `github.com/ClouGence/cloudcanal-openapi-cli`.
- Updated install, uninstall, and help documentation links to use the ClouGence repository URLs.

## [0.1.0] - 2026-03-19

### Added

- Added GitHub Releases based installation and uninstall scripts with checksum verification.
- Added `--output json` support for machine-readable CLI output.
- Added configurable HTTP timeout and read-retry settings.
- Added zsh and bash TAB completion support.

### Changed

- Moved installed binaries, completions, config, and build logs under `~/.cloudcanal-cli`.
- Simplified the README quick start flow and installation documentation.
- Improved CLI initialization resilience and network behavior.

### Removed

- Removed old directory compatibility cleanup logic from install and uninstall scripts.
- Removed automatic migration from `~/.cloudcanal` to `~/.cloudcanal-cli`.
