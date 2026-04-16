#!/usr/bin/env bash

set -euo pipefail

branch_name="$(git rev-parse --abbrev-ref HEAD)"

if [[ "$branch_name" == "main" ]]; then
  echo "Refusing to open a PR from main. Create a feat/*, fix/*, or ci/* branch first." >&2
  exit 1
fi

if [[ ! "$branch_name" =~ ^(feat|fix|ci)/ ]]; then
  echo "Branch must start with feat/, fix/, or ci/. Current branch: $branch_name" >&2
  exit 1
fi

if [[ -n "$(git status --porcelain)" ]]; then
  echo "Working tree is not clean. Commit or stash changes before opening a PR." >&2
  exit 1
fi

if ! command -v gh >/dev/null 2>&1; then
  echo "GitHub CLI (gh) is required to open a pull request." >&2
  exit 1
fi

if ! git remote get-url origin >/dev/null 2>&1; then
  echo "No origin remote found." >&2
  exit 1
fi

git push -u origin "$branch_name"

pr_title="${1:-}"
pr_body="${2:-}"

if [[ -z "$pr_title" ]]; then
  pr_title="$(git log -1 --pretty=%s)"
fi

if [[ -n "$pr_body" ]]; then
  gh pr create --base main --head "$branch_name" --title "$pr_title" --body "$pr_body"
else
  gh pr create --base main --head "$branch_name" --title "$pr_title" --fill
fi