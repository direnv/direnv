#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null && pwd)"
cd "$SCRIPT_DIR/.."

version="${1:-}"
if [[ -z $version ]]; then
  echo "USAGE: $0 version" >&2
  echo "Example: $0 v2.37.0" >&2
  exit 1
fi

# Remove 'v' prefix if present for consistency
version_number="${version#v}"
version_tag="v${version_number}"

if [[ "$(git symbolic-ref --short HEAD)" != "master" ]]; then
  echo "must be on master branch" >&2
  exit 1
fi

waitForPr() {
  local pr_branch=$1
  while true; do
    if gh pr view "$pr_branch" --json state --jq '.state' | grep -q 'MERGED'; then
      break
    fi
    echo "Waiting for PR to be merged..."
    sleep 5
  done
}

# Ensure we are up-to-date
uncommitted_changes=$(git diff --compact-summary)
if [[ -n $uncommitted_changes ]]; then
  echo -e "There are uncommitted changes, exiting:\n${uncommitted_changes}" >&2
  exit 1
fi

git fetch origin
unpushed_commits=$(git log --format=oneline origin/master..master)
if [[ $unpushed_commits != "" ]]; then
  echo -e "\nThere are unpushed changes, exiting:\n$unpushed_commits" >&2
  exit 1
fi

# Make sure we're up to date with remote
git pull origin master

# Make sure tag does not exist
if git tag -l | grep -q "^${version_tag}\$"; then
  echo "Tag ${version_tag} already exists, exiting" >&2
  exit 1
fi

echo "Generating changelog..."
git changelog

echo ""
echo "Please review and edit the CHANGELOG.md file."
echo "Opening editor..."
"${EDITOR:-nano}" CHANGELOG.md

echo ""
read -p "Continue with release ${version_tag}? (y/N): " -r
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
  echo "Release aborted. Resetting changes..."
  git checkout -- CHANGELOG.md
  exit 0
fi

echo "Updating version.txt to ${version_number}..."
echo "${version_number}" > version.txt

echo "Creating release branch..."
git branch -D "release-${version_tag}" 2>/dev/null || true
git checkout -b "release-${version_tag}"

echo "Committing changes..."
git add version.txt CHANGELOG.md
git commit -m "Release ${version_tag}"

echo "Pushing release branch..."
git push origin "release-${version_tag}"

echo "Creating pull request..."
pr_url=$(gh pr create \
  --base master \
  --head "release-${version_tag}" \
  --title "Release ${version_tag}" \
  --body "Release ${version_tag}")

# Extract PR number from URL
pr_number=$(echo "$pr_url" | grep -oE '[0-9]+$')

echo "Enabling auto-merge..."
gh pr merge "$pr_number" --auto --merge --delete-branch

echo "Switching back to master..."
git checkout master

echo "Waiting for PR to be merged..."
waitForPr "release-${version_tag}"

echo "Fetching latest master..."
git pull origin master

echo "Creating and pushing tag..."
git tag "${version_tag}"
git push origin "${version_tag}"

echo ""
echo "✅ Release ${version_tag} has been created!"
echo "🚀 CI will now build and publish the binaries automatically."
echo "📦 Check the release at: https://github.com/direnv/direnv/releases/tag/${version_tag}"