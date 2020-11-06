{{- define "issuerAnnotations" -}}
"helm.sh/hook": "post-install"
"helm.sh/hook-delete-policy": "before-hook-creation,hook-succeeded,hook-failed"
{{- end -}}

{{- define "issuerLabels" -}}
app: "{{ .Values.name }}"
giantswarm.io/service-type: "managed"
{{- end -}}
