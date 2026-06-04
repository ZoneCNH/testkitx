#!/usr/bin/env bash
set -euo pipefail

IMPORT_PATTERN='"github.com/ZoneCNH/testkitx/pkg/testkitx'

count_test_imports() {
  local repo="$1"
  local count=0
  while IFS= read -r -d '' file; do
    if grep -Fq "$IMPORT_PATTERN" "$file"; then
      count=$((count + 1))
    fi
  done < <(find "$repo" -path '*/vendor/*' -prune -o -type f -name '*_test.go' -print0)
  printf '%s\n' "$count"
}

count_production_imports() {
  local repo="$1"
  local count=0
  while IFS= read -r -d '' file; do
    if grep -Fq "$IMPORT_PATTERN" "$file"; then
      count=$((count + 1))
    fi
  done < <(find "$repo" -path '*/vendor/*' -prune -o -type f -name '*.go' ! -name '*_test.go' -print0)
  printf '%s\n' "$count"
}

check_downstream_repo() {
  local repo="$1"
  if [[ ! -d "$repo" ]]; then
    echo "ERROR: DOWNSTREAM_REPO does not exist: $repo" >&2
    exit 1
  fi
  if [[ ! -f "$repo/go.mod" ]]; then
    echo "ERROR: DOWNSTREAM_REPO must point to a Go module with go.mod: $repo" >&2
    exit 1
  fi

  local test_imports
  local production_imports
  test_imports="$(count_test_imports "$repo")"
  production_imports="$(count_production_imports "$repo")"

  if [[ "$test_imports" -eq 0 ]]; then
    echo "ERROR: downstream repo has no *_test.go import of testkitx" >&2
    exit 1
  fi
  if [[ "$production_imports" -ne 0 ]]; then
    echo "ERROR: downstream repo imports testkitx from production Go files" >&2
    exit 1
  fi

  if [[ "${DOWNSTREAM_SKIP_CI:-0}" != "1" ]]; then
    (
      cd "$repo"
      GOWORK=off make ci
    )
  fi

  echo "downstream adoption check passed: test_imports=$test_imports production_imports=$production_imports"
}

check_evidence_file() {
  local evidence="$1"
  if [[ ! -f "$evidence" ]]; then
    echo "ERROR: DOWNSTREAM_ADOPTION_EVIDENCE does not exist: $evidence" >&2
    exit 1
  fi
  go run ./internal/tools/downstreamadoption -verify "$evidence"
}

if [[ -n "${DOWNSTREAM_REPO:-}" ]]; then
  check_downstream_repo "$DOWNSTREAM_REPO"
  exit 0
fi

if [[ -n "${DOWNSTREAM_ADOPTION_EVIDENCE:-}" ]]; then
  check_evidence_file "$DOWNSTREAM_ADOPTION_EVIDENCE"
  exit 0
fi

echo "ERROR: set DOWNSTREAM_REPO=/path/to/downstream or DOWNSTREAM_ADOPTION_EVIDENCE=/path/to/evidence.json" >&2
exit 1
