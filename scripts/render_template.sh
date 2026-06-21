#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
Usage:
  scripts/render_template.sh --module-name NAME --module-path PATH --package-name NAME --out DIR

Renders testkitx into a concrete base library by copying the repository,
moving pkg/testkitx to pkg/<package>, and replacing template identifiers.
USAGE
}

module_name=""
module_path=""
package_name=""
out_dir=""

source_module_name="testkitx"
source_module_path="github.com/ZoneCNH/testkitx"
source_package_name="testkitx"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --module-name)
      module_name="${2:-}"
      shift 2
      ;;
    --module-path)
      module_path="${2:-}"
      shift 2
      ;;
    --package-name)
      package_name="${2:-}"
      shift 2
      ;;
    --out)
      out_dir="${2:-}"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "ERROR: unknown argument: $1" >&2
      usage >&2
      exit 2
      ;;
  esac
done

if [[ -z "$module_name" || -z "$module_path" || -z "$package_name" || -z "$out_dir" ]]; then
  echo "ERROR: --module-name, --module-path, --package-name and --out are required" >&2
  usage >&2
  exit 2
fi

if [[ "$package_name" =~ [^a-zA-Z0-9_] || "$package_name" =~ ^[0-9] ]]; then
  echo "ERROR: --package-name must be a valid Go package identifier" >&2
  exit 2
fi

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
mkdir -p "$out_dir"
if find "$out_dir" -mindepth 1 -maxdepth 1 | read -r _; then
  echo "ERROR: output directory must be empty: $out_dir" >&2
  exit 2
fi

(
  cd "$repo_root"
  tar \
    --exclude='./.git' \
    --exclude='./.omc' \
    --exclude='./.omx' \
    --exclude='./.worktree' \
    --exclude='./*.out' \
    --exclude='./coverage.out' \
    --exclude='./release/manifest/latest.json' \
    --exclude='./release/manifest/latest.json.sha256' \
    -cf - .
) | (
  cd "$out_dir"
  tar -xf -
)

if [[ "$package_name" != "$source_package_name" ]]; then
  mkdir -p "$out_dir/pkg"
  mv "$out_dir/pkg/$source_package_name" "$out_dir/pkg/$package_name"
fi

replace_in_text_files() {
  local find_text="$1"
  local replace_text="$2"

  while IFS= read -r -d '' file; do
    FIND_TEXT="$find_text" REPLACE_TEXT="$replace_text" perl -0pi -e 's/\Q$ENV{FIND_TEXT}\E/$ENV{REPLACE_TEXT}/g' "$file"
  done < <(
    find "$out_dir" -type f \( \
      -name '*.go' -o \
      -name '*.md' -o \
      -name '*.json' -o \
      -name '*.sh' -o \
      -name '*.yml' -o \
      -name '*.yaml' -o \
      -name 'Makefile' -o \
      -name 'go.mod' \
    \) -print0
  )
}

package_token="__TESTKITX_PACKAGE_NAME__"
if [[ "$source_module_name" == "$source_package_name" && "$module_name" != "$package_name" ]]; then
  replace_in_text_files "pkg/$source_package_name" "pkg/$package_token"
  replace_in_text_files "package $source_package_name" "package $package_token"
  replace_in_text_files "Package $source_package_name" "Package $package_token"
  replace_in_text_files "$source_package_name." "$package_token."
  replace_in_text_files "source_package_name=\"$source_package_name\"" "source_package_name=\"$package_token\""
  replace_in_text_files "\$package_name\" != \"$source_package_name\"" "\$package_name\" != \"$package_token\""
  replace_in_text_files "\\b$source_package_name\\b" "\\b$package_token\\b"
fi

replace_in_text_files "$source_module_path" "$module_path"
replace_in_text_files "$source_module_name" "$module_name"
replace_in_text_files "$package_token" "$package_name"
if [[ "$source_module_name" != "$source_package_name" || "$module_name" == "$package_name" ]]; then
  replace_in_text_files "$source_package_name" "$package_name"
fi

(
  cd "$out_dir"
  gofmt -w ./pkg ./internal ./contracts ./examples ./testkit
)

echo "rendered $module_name at $out_dir"
