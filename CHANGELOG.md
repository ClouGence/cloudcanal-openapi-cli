# Changelog

All notable changes to this project will be documented in this file.

## [0.1.3] - 2026-04-02

### Added

- Added startup-time release checks that compare the local build version with the latest GitHub release and suggest a one-line `curl` upgrade command when an update is available.
- Added a dedicated `internal/updatecheck` module with semver-aware comparison and tests for release redirect parsing.

### Changed

- Updated one-shot text commands to print the active profile and API endpoint before command output so users can see which environment is in use.
- Improved interactive initialization so pressing `Ctrl+C` is treated as a clean cancellation instead of surfacing the raw `prompt aborted` error.

## [0.1.2] - 2026-04-02

### Added

- Added `version` and `--version` commands with `version`, `commit`, and `buildTime` output.
- Added profile-aware configuration management with `config profiles list|use|add|remove`.
- Added build metadata injection and release asset packaging via a shared `make release-assets` flow.

### Changed

- Switched CLI configuration storage to `language + currentProfile + profiles` schema under `~/.cloudcanal-cli/config.json`.
- Updated `config show`, REPL prompt, help text, completion, and docs to expose the active profile context.
- Enhanced release delivery to publish installer assets and print installed build metadata after installation.

### Removed

- Removed support for silently reusing the legacy single-profile config format; users are now prompted to reinitialize into the profile-based schema.

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
