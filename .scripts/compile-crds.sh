#!/bin/bash

SRC_DIR="cert-manager/deploy/crds"
DEST_FILE="helm/cert-manager/templates/crds.yaml"

# shellcheck disable=SC2046
mkdir -p $(dirname "$DEST_FILE")

# shellcheck disable=SC2188
> "$DEST_FILE"


for file in "$SRC_DIR"/*.yaml
do
    echo "---" >> "$DEST_FILE"
    cat "$file" >> "$DEST_FILE"
done

echo "CRDs have been compiled into $DEST_FILE"
