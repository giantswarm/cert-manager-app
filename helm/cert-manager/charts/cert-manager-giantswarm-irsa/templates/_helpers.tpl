{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "name" -}}
{{- .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "labels.common" -}}
app.kubernetes.io/name: {{ include "name" . | quote }}
app.kubernetes.io/instance: {{ .Release.Name | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service | quote }}
helm.sh/chart: {{ include "chart" . | quote }}
{{- end -}}

{{/*
Get list of all provided OIDC domains
*/}}
{{- define "oidcDomains" -}}
{{- $oidcDomains := list .Values.oidcDomain -}}
{{- if .Values.oidcDomains -}}
{{- $oidcDomains = concat $oidcDomains .Values.oidcDomains -}}
{{- end -}}
{{- compact $oidcDomains | uniq | toJson -}}
{{- end -}}
