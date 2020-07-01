{{/* vim: set filetype=mustache: */}}

{{/* Create a default fully qualified app name. Truncated to meet DNS naming spec. */}}
{{- define "certManager.name" -}}
{{- default .Chart.Name .Values.name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/* Create chart name and version as used by the chart label. */}}
{{- define "certManager.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "certManager.defaultLabels" -}}
app: {{ template "certManager.name" . }}
app.kubernetes.io/name: {{ template "certManager.name" . }}
app.kubernetes.io/instance: "{{ .Release.Name }}"
app.kubernetes.io/managed-by: "{{ .Release.Service }}"
helm.sh/chart: {{ template "certManager.chart" . }}
giantswarm.io/service-type: "managed"
{{- end -}}

{{- define "certManager.selectorLabels" -}}
app.kubernetes.io/name: {{ template "certManager.name" . }}
app.kubernetes.io/instance: "{{ .Release.Name }}"
{{- end -}}
