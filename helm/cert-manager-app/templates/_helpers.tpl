{{/* vim: set filetype=mustache: */}}

{{/* Create a default fully qualified app name. Truncated to meet DNS naming spec. */}}
{{- define "certManager.name" -}}
{{- default .Chart.Name .Values.global.name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/* Create names for each component to avoid repetition. */}}
{{- define "certManager.name.cainjector" -}}
{{- printf "%s-%s" (include "certManager.name" . ) "cainjector" | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "certManager.name.controller" -}}
{{- printf "%s-%s" ( include "certManager.name" . ) "controller" | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "certManager.name.crdInstall" -}}
{{- printf "%s-%s" ( include "certManager.name" . ) "crd-install" | replace "+" "_" | trimSuffix "-" -}}
{{- end -}}

{{- define "certManager.name.webhook" -}}
{{- printf "%s-%s" ( include "certManager.name" . ) "webhook" | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/* Create chart name and version as used by the chart label. */}}
{{- define "certManager.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "certManager.defaultLabels" -}}
app: "{{ template "certManager.name" . }}"
app.kubernetes.io/name: "{{ template "certManager.name" . }}"
app.kubernetes.io/instance: "{{ template "certManager.name" . }}"
app.kubernetes.io/managed-by: "{{ .Release.Service }}"
helm.sh/chart: "{{ template "certManager.chart" . }}"
giantswarm.io/service-type: "managed"
{{- end -}}

{{- define "certManager.CRDInstallAnnotations" -}}
"helm.sh/hook": "pre-install,pre-upgrade"
"helm.sh/hook-delete-policy": "before-hook-creation,hook-succeeded,hook-failed"
{{- end -}}

{{- define "certManager.CRDLabels" -}}
app: "{{ template "certManager.name" . }}"
app.kubernetes.io/name: "{{ template "certManager.name" . }}"
app.kubernetes.io/instance: "{{ template "certManager.name" . }}"
app.kubernetes.io/managed-by: "{{ .Release.Service }}"
helm.sh/chart: "{{ .Chart.Name }}-{{ .Chart.Version }}"
giantswarm.io/service-type: "managed"
cert-manager.io/disable-validation: "true"
{{- end -}}

{{- define "certManager.selectorLabels" -}}
app.kubernetes.io/name: "{{ template "certManager.name" . }}"
app.kubernetes.io/instance: "{{ template "certManager.name" . }}"
{{- end -}}

{{/* Create a label which can be used to select any orphaned crd-install hook resources */}}
{{- define "certManager.CRDInstallSelector" -}}
{{- printf "%s" "crd-install-hook" -}}
{{- end -}}

{{/*
startupapicheck templates
*/}}

{{/*
Expand the name of the chart.
Manually fix the 'app' and 'name' labels to 'startupapicheck' to maintain
compatibility with the v0.9 deployment selector.
*/}}
{{- define "startupapicheck.name" -}}
{{- printf "startupapicheck" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "startupapicheck.fullname" -}}
{{- $trimmedName := printf "%s" (include "certManager.name" .) | trunc 52 | trimSuffix "-" -}}
{{- printf "%s-startupapicheck" $trimmedName | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create the name of the service account to use
*/}}
{{- define "startupapicheck.serviceAccountName" -}}
{{- if .Values.startupapicheck.serviceAccount.create -}}
    {{ default (include "startupapicheck.fullname" .) .Values.startupapicheck.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.startupapicheck.serviceAccount.name }}
{{- end -}}
{{- end -}}
