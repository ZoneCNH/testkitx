#!/usr/bin/env bash
set -euo pipefail

version="${1:-${VERSION:-}}"

if [[ -z "$version" ]]; then
  echo "ERROR: set VERSION=vX.Y.Z when running release preflight"
  exit 1
fi

if [[ ! "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+([-+][0-9A-Za-z.-]+)?$ ]]; then
  echo "ERROR: release version must look like vX.Y.Z: $version"
  exit 1
fi

if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  echo "ERROR: release preflight must run inside a git worktree"
  exit 1
fi

branch="$(git rev-parse --abbrev-ref HEAD)"
if [[ "$branch" != "main" ]]; then
  echo "ERROR: release preflight must run on main; current branch is $branch"
  exit 1
fi

if [[ -n "$(git status --porcelain)" ]]; then
  echo "ERROR: release preflight requires a clean git worktree"
  git status --short
  exit 1
fi

git fetch --quiet origin main --tags

head_sha="$(git rev-parse HEAD)"
origin_main_sha="$(git rev-parse origin/main)"
if [[ "$head_sha" != "$origin_main_sha" ]]; then
  echo "ERROR: local main is not aligned with origin/main"
  echo "HEAD=$head_sha"
  echo "origin/main=$origin_main_sha"
  exit 1
fi

if git rev-parse -q --verify "refs/tags/$version" >/dev/null; then
  echo "ERROR: local tag already exists: $version"
  exit 1
fi

if git ls-remote --exit-code --tags origin "refs/tags/$version" >/dev/null 2>&1; then
  echo "ERROR: remote tag already exists: $version"
  exit 1
fi

if ! grep -Eq "^## \\[?$version\\]?( |$)" CHANGELOG.md; then
  echo "ERROR: CHANGELOG.md must contain a release heading for $version"
  exit 1
fi

release_pkg_version="$(sed -n 's/^ *Version[[:space:]]*=[[:space:]]*"\([^"]*\)".*/\1/p' pkg/testkitx/version.go | head -n1)"
if [[ -z "$release_pkg_version" ]]; then
  echo "ERROR: could not read pkg/testkitx/version.go version"
  exit 1
fi

if [[ "$release_pkg_version" != "$version" ]]; then
  echo "ERROR: pkg/testkitx/version.go version ($release_pkg_version) does not match VERSION ($version)"
  exit 1
fi

if ! grep -Eq "^  table_version: ${version}$" .repo-contract.yaml; then
  echo "ERROR: .repo-contract.yaml table_version must match VERSION ($version)"
  exit 1
fi

if ! grep -Eq "^  latest_git_tag: ${version}$" .repo-contract.yaml; then
  echo "ERROR: .repo-contract.yaml latest_git_tag must match VERSION ($version)"
  exit 1
fi

if ! grep -Eq "\"version\": \"${version}\"" release/manifest/template.json; then
  echo "ERROR: release/manifest/template.json version must match VERSION ($version)"
  exit 1
fi

for tool in golangci-lint govulncheck; do
  if ! command -v "$tool" >/dev/null 2>&1; then
    echo "ERROR: $tool not installed"
    exit 1
  fi
done

echo "release preflight metadata checks passed for $version"
