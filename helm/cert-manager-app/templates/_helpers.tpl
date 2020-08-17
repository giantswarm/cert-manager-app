{{/* vim: set filetype=mustache: */}}

{{/* Create a default fully qualified app name. Truncated to meet DNS naming spec. */}}
{{- define "certManager.name" -}}
{{- default .Chart.Name .Values.global.name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/* Generate a string containing the current Unix time for appending to CRD install resource names. */}}
{{- define "unixTime.now" -}}
{{- $time := now | unixEpoch }}
{{- printf "%s" $time -}}
{{- end -}}

{{/* Create names for each component to avoid repetition. */}}
{{- define "certManager.name.cainjector" -}}
{{- printf "%s-%s" (include "certManager.name" . ) "cainjector" | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "certManager.name.controller" -}}
{{- printf "%s-%s" ( include "certManager.name" . ) "controller" | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{*/
Create a unique name for crd-install hook resource names. We don't truncate
it as we don't want to lose the end of the epoch time - this is what makes
the name unique.
*/}}
{{- define "certManager.name.crdInstall" -}}
{{- printf "%s-%s-%s" ( include "certManager.name" . ) "crd-install" ( include "unixTime.now" . ) | replace "+" "_" | trimSuffix "-" -}}
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
"helm.sh/hook-delete-policy": "hook-succeeded,hook-failed"
{{- end -}}

{{- define "certManager.CRDLabels" -}}
app: "{{ template "certManager.name" . }}"
app.kubernetes.io/name: "{{ template "certManager.name" . }}"
app.kubernetes.io/instance: "{{ template "certManager.name" . }}"
app.kubernetes.io/managed-by: "{{ .Release.Service }}"
helm.sh/chart: "{{ .Chart.Name }}-{{ .Chart.Version }}"
giantswarm.io/service-type: "managed"
{{- end -}}

{{- define "certManager.selectorLabels" -}}
app.kubernetes.io/name: "{{ template "certManager.name" . }}"
app.kubernetes.io/instance: "{{ template "certManager.name" . }}"
{{- end -}}
