#!/usr/bin/env bash
set -euo pipefail

echo "checking secrets..."

PATTERNS=(
  "(^|[^A-Za-z0-9_])password[[:space:]]*=([^=]|$)"
  "(^|[^A-Za-z0-9_])passwd[[:space:]]*=([^=]|$)"
  "(^|[^A-Za-z0-9_])secret[[:space:]]*=([^=]|$)"
  "(^|[^A-Za-z0-9_])token[[:space:]]*=([^=]|$)"
  "(^|[^A-Za-z0-9_])access_key[[:space:]]*=([^=]|$)"
  "(^|[^A-Za-z0-9_])secret_key[[:space:]]*=([^=]|$)"
  "AKIA[0-9A-Z]{16}"
  "BEGIN RSA PRIVATE KEY"
  "BEGIN OPENSSH PRIVATE KEY"
)

for pattern in "${PATTERNS[@]}"; do
  if grep -R -E "$pattern" . \
    --exclude-dir=.git \
    --exclude-dir=.omx \
    --exclude-dir=vendor \
    --exclude="*.sum" \
    --exclude="check_secrets.sh" \
    --exclude="goal.md"; then
    echo "ERROR: possible secret found: $pattern"
    exit 1
  fi
done

echo "secret check passed"
