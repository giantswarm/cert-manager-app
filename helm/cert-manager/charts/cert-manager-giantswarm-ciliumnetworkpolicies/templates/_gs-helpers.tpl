{{/*
Before trying to contribute this file to upstream, please read below.
This helpers file contains Giant Swarm specific overrides to helpers defined
in the original upstream _helpers.tpl file.
*/}}

{{/*
Labels that should be added on each resource
*/}}
{{- define "labels" -}}
giantswarm.io/service-type: "managed"
giantswarm.io/monitoring_basic_sli: "true"
application.giantswarm.io/team: {{ index .Chart.Annotations "application.giantswarm.io/team" | quote }}
{{- if eq (default "helm" .Values.creator) "helm" }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}
{{- end -}}
{{/*
Define 'webhook.fullname' template
*/}}
{{- define "webhook.fullname" -}}
{{- printf "%s-webhook" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}