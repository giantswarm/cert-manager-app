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

HELM_DOCS="docker run --rm -u $(id -u) -v ${PWD}:/helm-docs -w /helm-docs jnorwood/helm-docs:v1.11.0"
$HELM_DOCS --template-files=sync/readme.gotmpl -g helm/cert-manager -f values.yaml -o README.md

# Store diffs
rm -f ./diffs/*
for f in $(git --no-pager diff --no-exit-code --no-color --no-index vendor/cert-manager helm/cert-manager --name-only) ; do
        # Skip helm/cert-manager/Chart.yaml; as we take it as our own.
        [[ "$f" == "helm/cert-manager/Chart.yaml" ]] && continue
        # Skip helm/cert-manager/README.md; as it's autogenerated.
        [[ "$f" == "helm/cert-manager/README.md" ]] && continue

        base_file="vendor/cert-manager/${f#"helm/cert-manager/"}"
        [[ ! -e $base_file ]] && base_file="/dev/null"

        set +e
        set -x

        git --no-pager diff --no-exit-code --no-color --no-index "$base_file" "${f}" \
                > "./diffs/${f//\//__}.patch" # ${f//\//__} replaces all "/" with "__"

        ret=$?
        { set +x; } 2>/dev/null
        set -e
        if [ $ret -ne 0 ] && [ $ret -ne 1 ] ; then
                exit $ret
        fi
done
cp -R ./sync/charts ./helm/cert-manager
