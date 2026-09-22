{{/*
Before trying to contribute this file to upstream, please read below.
This helpers file contains Giant Swarm specific overrides to helpers defined
in the original upstream _helpers.tpl file.
*/}}

{{/*
Chart name and version as a valid Kubernetes label value.
The upstream `chartName` helper truncates to 63 characters and trims only a
trailing "-". A long chart version, as dev builds carry because they embed the
branch name, can leave a trailing "." or "_", which the API server rejects.
*/}}
{{- define "chartLabel" -}}
{{- regexReplaceAll "[^a-zA-Z0-9]+$" (include "chartName" .) "" -}}
{{- end -}}

{{/*
Labels that should be added on each resource
*/}}
{{- define "labels" -}}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
giantswarm.io/service-type: "managed"
application.giantswarm.io/team: {{ .Values.global.team | default "shield" | quote }}
{{- if eq (default "helm" .Values.creator) "helm" }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
helm.sh/chart: {{ include "chartLabel" . }}
{{- end -}}
{{- end -}}

{{/*
Override for original helper because Giant Swarm cert-manager chart v2 label selectors are different
*/}}
{{- define "cainjector.name" -}}
{{- printf "%s" (include "cert-manager.name" .) -}}
{{- end -}}

{{/*
Override for original helper because Giant Swarm cert-manager chart v2 label selectors are different
*/}}
{{- define "webhook.name" -}}
{{- printf "%s" (include "cert-manager.name" .) -}}
{{- end -}}
