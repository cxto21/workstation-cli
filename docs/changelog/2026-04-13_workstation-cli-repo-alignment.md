# 2026-04-13 Workstation CLI Repo Alignment

## Session Goal

Retarget public install and update references to `cxto21/workstation-cli` so the codebase matches the requested repository namespace.

## Changes Made

- Updated the installer defaults and usage example in `install.sh`.
- Updated the self-update command in `src/main.rs`.
- Updated the daemon release-check URL in `src/daemon/service.rs`.
- Updated public repository links and badges in `README.md`.

## Session Update: Prebuilt Release Flow

- Added a GitHub Actions release workflow that builds precompiled Linux and macOS tarballs and publishes checksums to GitHub Releases.
- Updated the release guide to describe the prebuilt distribution path and release artifact expectations.
- Clarified in the README that the quick install path uses release binaries instead of local compilation.
- Verified the installed CLI command with `workstation-cli --version` in the local environment after placing the binary on the PATH.

## Verification

`git diff --check` passed. `cargo test --locked` passed for the project except for one pre-existing failure in `tests/input_tests.rs` (`content_delete_page_and_function_keys_are_encoded`), which is unrelated to the release workflow/docs edits.

## Session Update: Onboarding Palette and Template Cleanup

- Changed the onboarding UI to a strict black-and-white palette.
- Removed the `Mato Creator Office` template from both onboarding entry points.
- Bumped the crate version to `0.9.7` so the next release can carry these UI changes.