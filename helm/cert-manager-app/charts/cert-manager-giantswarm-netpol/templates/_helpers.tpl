{{- define "commonLabels" -}}
app.kubernetes.io/name: {{ .Values.name }}
giantswarm.io/service-type: "managed"
{{- end -}}

