#!/bin/bash

# Directory containing the template files
TEMPLATES_DIR="templates"

update_image_reference() {
    local file="$1"
    local component="$2"

    sed -i '
    /image:.*template "image".*Values.'$component'\.image/c\
          {{- with .Values.'$component'.image }}\
          image: "{{- if .registry -}}{{ .registry }}/{{- end -}}{{ .repository }}{{- if (.digest) -}} @{{ .digest }}{{- else -}}:{{ default $.Chart.AppVersion .tag }} {{- end -}}"\
          {{- end }}
    ' "$file"
}

for file in "$TEMPLATES_DIR"/*.yaml
do
    update_image_reference "$file" "cainjector"
    update_image_reference "$file" "controller"
    update_image_reference "$file" "webhook"
    update_image_reference "$file" "startupapicheck"
done

echo "Image references have been updated in template files."
