#!/usr/bin/env bash
# Seeds the repo with issues (for Issue-to-Code) and flawed pull requests
# (for Automated Code Review). Requires the GitHub CLI, authenticated, and a
# remote named origin that points at your repo.
set -euo pipefail

cd "$(dirname "$0")/.."
git rev-parse --abbrev-ref HEAD | grep -qx main || { echo "Run from the main branch"; exit 1; }

echo "Creating issues..."
for f in seed/issues/*.md; do
  title=$(head -1 "$f" | sed 's/^# //')
  body=$(tail -n +3 "$f")
  gh issue create --title "$title" --body "$body"
done

echo "Creating flawed pull requests..."
title_for() {
  case "$1" in
    01-python-click-log) echo "Log each click with its user agent" ;;
    02-node-admin-delete) echo "Add admin endpoint to delete a link" ;;
    03-go-top-links) echo "Add endpoint for top links by clicks" ;;
    *) echo "Seeded change $1" ;;
  esac
}
for p in seed/patches/*.patch; do
  key=$(basename "$p" .patch)
  branch="seed/$key"
  git checkout -q main
  git checkout -q -B "$branch"
  git apply "$p"
  git commit -qam "$(title_for "$key")"
  git push -q -u origin "$branch"
  gh pr create --base main --head "$branch" --title "$(title_for "$key")" \
    --body "Seeded pull request for the Dispatch trial. Intentionally flawed; see README for what to look for."
done
git checkout -q main
echo "Done."
