# Contribution Guide

Keep changes small, reviewable, and focused.

## Branches

- `feat/*` for new features
- `fix/*` for bug fixes
- `ci/*` for CI and release workflow changes
- Use short, clear branch names

## Flow

1. Create a branch from `main`
2. Make one logical change per branch
3. Validate locally when possible
4. Open a pull request to `main`
5. Ask for review from another team member
6. Merge only after approval

Use `scripts/create-pr.sh` to push the branch and open the PR.

Hotfixes also follow the same flow with a `fix/*` branch.

## Rules

- Do not push directly to `main`
- Keep PRs small and easy to review
- Prefer docs and tests with code changes
- Use the release workflow only after version bumps on `main`
- Hotfixes merged into `main` will be tagged automatically by the release workflow

## For Agents

- Follow the branch naming rules above
- Make the smallest safe change
- Do not skip the PR step
- Keep explanations simple and concise