#!/usr/bin/env bash

set -euo pipefail

if [[ $# -ne 1 || -z "$1" ]]; then
  echo "usage: $0 <controller-image>" >&2
  exit 2
fi

readonly image="$1"
readonly repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
readonly temp_dir="$(mktemp -d "${TMPDIR:-/tmp}/deptrack-operator-kustomize.XXXXXX")"
trap 'rm -rf "$temp_dir"' EXIT

cp -R "$repo_root/config" "$temp_dir/config"
(
  cd "$temp_dir/config/manager"
  kustomize edit set image "controller=${image}"
)
kustomize build "$temp_dir/config/default"
