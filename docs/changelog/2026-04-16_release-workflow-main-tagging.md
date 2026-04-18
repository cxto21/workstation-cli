# 2026-04-16 Release Workflow Main Tagging

## Summary

- Changed the GitHub release workflow to trigger from pushes to `main` instead of tag-only pushes.
- The release job now reads the version from `Cargo.toml`, creates the matching `vX.Y.Z` tag, pushes it, and publishes the GitHub release from that tag.
- Updated the release guide so the documented process matches the automated `main`-driven flow.

## Notes

- This keeps the release source of truth in `Cargo.toml`.
- A version bump is still required before releasing again.
- The workflow is idempotent for reruns on the same commit: an existing matching tag is reused, while a tag pointing at another commit still fails fast.

## Follow-up

- Added a short contribution guide at `CONTRIB.md` so small-team branch and PR flow stays explicit for humans and agents.
- Added `scripts/create-pr.sh` and aligned `AGENTS.md` with the branch, push, and PR flow used by the team.
- Clarified hotfix handling: hotfixes should still use `fix/*` and PRs to `main`, while any emergency change landed on `main` will tag automatically from the version in `Cargo.toml`.
- Updated release workflow macOS runner from `macos-13` to `macos-14` after CI failure with unsupported runner configuration.
- Fixed `install.sh` asset selection loop to avoid silent failure when the first Linux candidate asset is missing.
- Updated README install/clone/repo URLs to use `reflecterlabs/workstation-cli` consistently.
- Opted GitHub JavaScript actions into Node 24 in `release.yml` to remove Node 20 deprecation warnings before enforcement.
- Bumped project version to `0.9.9` for next automatic release tag from `main`.
- Switched project license metadata and README badge from MIT to BSD-3-Clause and added a `LICENSE` file.
- Added `TRADEMARK.md` to keep the code open under BSD-3-Clause while protecting the project name, logo, and branding.
- Aligned `Cargo.lock` with the `0.9.9` package version so the release workflow can build with `--locked`.
- Switched the x86_64 Linux release build from `ubuntu-20.04` to `ubuntu-24.04` to avoid the run cancellation that prevented tag publishing.