#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
report_dir="$repo_root/release/docker"
report="$report_dir/toolchain-check.md"
mkdir -p "$report_dir"

write_report() {
  local status="$1"
  local detail="$2"

  {
    echo "# Docker Toolchain Check"
    echo
    echo "- status: $status"
    echo "- detail: $detail"
  } > "$report"
}

allow_missing="${XLIB_DOCKER_ALLOW_MISSING:-0}"

if ! command -v docker >/dev/null 2>&1; then
  if [[ "$allow_missing" == "1" ]]; then
    write_report "unavailable_allowed" "docker command not found"
    echo "docker command not found; wrote $report"
    exit 0
  fi

  write_report "failed" "docker command not found"
  echo "ERROR: docker command not found" >&2
  exit 1
fi

if ! docker version >/dev/null 2>&1; then
  if [[ "$allow_missing" == "1" ]]; then
    write_report "unavailable_allowed" "docker daemon unavailable"
    echo "docker daemon unavailable; wrote $report"
    exit 0
  fi

  write_report "failed" "docker daemon unavailable"
  echo "ERROR: docker daemon unavailable" >&2
  exit 1
fi

if ! docker buildx version >/dev/null 2>&1; then
  write_report "failed" "docker buildx unavailable"
  echo "ERROR: docker buildx unavailable" >&2
  exit 1
fi

compose_detail="not installed"
if docker compose version >/dev/null 2>&1; then
  compose_detail="$(docker compose version)"
fi

write_report "passed" "$(docker --version); $(docker buildx version); docker compose: $compose_detail"
echo "docker toolchain check passed; wrote $report"
