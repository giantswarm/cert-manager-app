#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd ) ; readonly dir
cd "${dir}/.."

# Stage 1 sync
set -x
vendir sync
{ set +x; } 2>/dev/null

# Remove trailing whitespace end of lines (hack to fix vendir bug)
# find vendor/ -type f -exec sed -i 's/[[:space:]]*$//' {} \;

# Patches
./sync/patches/values/patch.sh
./sync/patches/chart/patch.sh
./sync/patches/image-registry/patch.sh
./sync/patches/pss/patch.sh
./sync/patches/webhook-pdb/patch.sh
./sync/patches/cainjector-service/patch.sh

HELM_DOCS="docker run --rm -u $(id -u) -v ${PWD}:/helm-docs -w /helm-docs jnorwood/helm-docs:v1.11.0"
$HELM_DOCS --template-files=sync/readme.gotmpl -g helm/cert-manager -f values.yaml -o README.md

# Store diffs
rm -f ./diffs/*

cp -R ./sync/charts ./helm/cert-manager
cp -R ./sync/templates/* ./helm/cert-manager/templates
cp  ./sync/.kube-linter.yaml ./helm/cert-manager/.kube-linter.yaml
cp  ./sync/.helmignore ./helm/cert-manager/.helmignore
