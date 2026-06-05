#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$repo_root"

required_files=(
  "Dockerfile"
  "docker-compose.yml"
  ".dockerignore"
  ".devcontainer/devcontainer.json"
  "scripts/docker/check_toolchain.sh"
  "scripts/docker/docker_gate.sh"
  "scripts/docker/check_contract.sh"
  "contracts/docker-toolchain.schema.json"
  "contracts/downstream-adoption-proof.schema.json"
)

required_targets=(
  "docker-toolchain-check"
  "docker-build"
  "docker-build-check"
  "docker-shell"
  "docker-ci"
  "docker-release-check"
  "docker-release-final-check"
  "docker-goalcli"
  "docker-goalcli-image"
  "docker-goalcli-version"
  "docker-runtime-check"
  "docker-drift-check"
  "docker-contract"
)

for file in "${required_files[@]}"; do
  if [[ ! -f "$file" ]]; then
    echo "ERROR: missing Docker toolchain contract file: $file" >&2
    exit 1
  fi
done

for target in "${required_targets[@]}"; do
  if ! grep -Eq "^\\.PHONY: .*${target}([[:space:]]|$)|^${target}:" Makefile; then
    echo "ERROR: missing Makefile Docker target: $target" >&2
    exit 1
  fi
done

go_version="$(awk '/^go / {print $2; exit}' go.mod)"
if [[ -z "$go_version" ]]; then
  echo "ERROR: unable to read Go version from go.mod" >&2
  exit 1
fi

if ! grep -q "GO_VERSION=${go_version}" Dockerfile; then
  echo "ERROR: Dockerfile must default to the go.mod major/minor toolchain" >&2
  exit 1
fi

echo "docker contract check passed"
