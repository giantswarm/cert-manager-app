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
{{ include "certManager.selectorLabels" . }}
app.kubernetes.io/managed-by: "{{ .Release.Service }}"
helm.sh/chart: "{{ template "certManager.chart" . }}"
giantswarm.io/service-type: "managed"
application.giantswarm.io/team: {{ index .Chart.Annotations "application.giantswarm.io/team" | quote }}
{{- end -}}

{{- define "certManager.CRDInstallAnnotations" -}}
helm.sh/hook: pre-install,pre-upgrade
helm.sh/hook-delete-policy: before-hook-creation,hook-succeeded,hook-failed
{{- end -}}

{{- define "certManager.CRDLabels" -}}
app: "{{ template "certManager.name" . }}"
app.kubernetes.io/name: "{{ template "certManager.name" . }}"
app.kubernetes.io/instance: "{{ template "certManager.name" . }}"
app.kubernetes.io/managed-by: "{{ .Release.Service }}"
helm.sh/chart: "{{ template "certManager.chart" . }}"
giantswarm.io/service-type: "managed"
cert-manager.io/disable-validation: "true"
{{- end -}}

{{- define "certManager.selectorLabels" -}}
app.kubernetes.io/name: "{{ template "certManager.name" . }}"
app.kubernetes.io/instance: "{{ template "certManager.name" . }}"
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

{{/*
Create the name of the service account to use
*/}}
{{- define "certManager.controller.serviceAccountName" -}}
{{- if .Values.controller.serviceAccount.create -}}
    {{ default (include "certManager.name.controller" .) .Values.controller.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.controller.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{/*
Set the role name for IRSA
*/}}
{{- define "aws.iam.role" -}}
{{- if .Values.controller.aws.role }}
{{- printf "%s" .Values.controller.aws.role }}
{{- else }}
{{- printf "%s-CertManager-Role" .Values.clusterID }}
{{- end }}
{{- end }}

{{/*
Set Giant Swarm serviceAccountAnnotations.
*/}}
{{- define "giantswarm.serviceAccountAnnotations" -}}
{{- if and (eq .Values.provider "aws") (eq .Values.controller.aws.irsa "true") (not (hasKey .Values.controller.serviceAccount.annotations "eks.amazonaws.com/role-arn")) }}
{{- $_ := set .Values.controller.serviceAccount.annotations "eks.amazonaws.com/role-arn" (tpl "arn:aws:iam::{{ .Values.aws.accountID }}:role/{{ template \"aws.iam.role\" . }}" .) }}
{{- else if and (eq .Values.provider "aws") (or (eq .Values.region "cn-north-1") (eq .Values.region "cn-northwest-1")) (eq .Values.aws.irsa "true") (not (hasKey .Values.controller.serviceAccount.annotations "eks.amazonaws.com/role-arn")) }}
{{- $_ := set .Values.controller.serviceAccount.annotations "eks.amazonaws.com/role-arn" (tpl "arn:aws-cn:iam::{{ .Values.aws.accountID }}:role/{{ .Values.clusterID }}-Route53Manager-Role" .) }}
{{- else if and (eq .Values.provider "aws") (eq .Values.aws.irsa "true") (not (hasKey .Values.controller.serviceAccount.annotations "eks.amazonaws.com/role-arn")) }}
{{- $_ := set .Values.controller.serviceAccount.annotations "eks.amazonaws.com/role-arn" (tpl "arn:aws:iam::{{ .Values.aws.accountID }}:role/{{ .Values.clusterID }}-Route53Manager-Role" .) }}
{{- else if and (eq .Values.provider "capa") (not (hasKey .Values.controller.serviceAccount.annotations "eks.amazonaws.com/role-arn")) }}
{{- $_ := set .Values.controller.serviceAccount.annotations "eks.amazonaws.com/role-arn" (include "aws.iam.role" .) }}
{{- end }}
{{- end -}}
