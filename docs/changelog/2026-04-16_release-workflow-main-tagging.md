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