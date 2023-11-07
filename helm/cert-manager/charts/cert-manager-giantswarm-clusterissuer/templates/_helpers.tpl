{{- define "issuerLabels" -}}
app.kubernetes.io/name: {{ .Values.name }}
giantswarm.io/service-type: "managed"
{{- end -}}

{{- define "issuerAnnotations" -}}
helm.sh/hook: post-install,post-upgrade
helm.sh/hook-delete-policy: before-hook-creation,hook-failed
{{- end -}}

{{- define "registry" }}
{{- $registry := .Values.image.registry -}}
{{- if and .Values.global (and .Values.global.image .Values.global.image.registry) -}}
{{- $registry = .Values.global.image.registry -}}
{{- end -}}
{{- printf "%s" $registry -}}
{{- end -}}
