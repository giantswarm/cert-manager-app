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

# Patches. vendir syncs the pristine upstream chart into
# helm/cert-manager/charts/cert-manager; these patches only tweak upstream-rendered
# objects that cannot live in the parent chart's templates.
./sync/patches/image-registry/patch.sh
./sync/patches/webhook-pdb/patch.sh
./sync/patches/cainjector-service/patch.sh

# Giant Swarm subcharts (clusterissuer, netpol, ciliumnetworkpolicies) live alongside
# the vendored cert-manager subchart.
cp -R ./sync/charts/* ./helm/cert-manager/charts/

# Giant-Swarm-only templates and the _gs-helpers.tpl overlay go into the parent chart.
cp -R ./sync/templates/* ./helm/cert-manager/templates/

cp  ./sync/.kube-linter.yaml ./helm/cert-manager/.kube-linter.yaml
cp  ./sync/.helmignore ./helm/cert-manager/.helmignore

HELM_DOCS="docker run --rm -u $(id -u) -v ${PWD}:/helm-docs -w /helm-docs jnorwood/helm-docs:v1.11.0"
$HELM_DOCS --template-files=sync/readme.gotmpl -g helm/cert-manager -f values.yaml -o README.md
