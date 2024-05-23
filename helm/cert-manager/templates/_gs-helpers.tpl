{{/*
Before trying to contribute this file to upstream, please read below.
This helpers file contains Giant Swarm specific overrides to helpers defined
in the original upstream _helpers.tpl file.
*/}}

{{/*
Labels that should be added on each resource
*/}}
{{- define "labels" -}}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
giantswarm.io/service-type: "managed"
giantswarm.io/monitoring_basic_sli: "true"
application.giantswarm.io/team: {{ index .Chart.Annotations "application.giantswarm.io/team" | quote }}
{{- if eq (default "helm" .Values.creator) "helm" }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
helm.sh/chart: {{ include "chartName" . }}
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

{{- define "issuerLabels" -}}
app.kubernetes.io/name: {{ .Values.name }}
giantswarm.io/service-type: "managed"
{{- end -}}

{{- define "issuerAnnotations" -}}
helm.sh/hook: post-install,post-upgrade
helm.sh/hook-delete-policy: before-hook-creation,hook-succeeded,hook-failed
{{- end -}}

{{- define "registry" }}
{{- $registry := .Values.image.registry -}}
{{- if and .Values.global (and .Values.global.image .Values.global.image.registry) -}}
{{- $registry = .Values.global.image.registry -}}
{{- end -}}
{{- printf "%s" $registry -}}
{{- end -}}

{{- define "clusterIssuer" }}
{{- if .Values.install }}
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-giantswarm
  labels:
    {{- include "issuerLabels" . | nindent 4 }}
  annotations:
    {{- include "issuerAnnotations" . | nindent 4 }}
spec:
  acme:
    # The ACME server URL.
    server: https://acme-v02.api.letsencrypt.org/directory
    # Email address used for ACME registration.
    email: accounts@giantswarm.io
    # Secret resource used to store the account's private key.
    privateKeySecretRef:
      name: letsencrypt-giantswarm
    # Add challenge solvers
    solvers:
    {{ if .Values.acme.dns01.cloudflare.enabled -}}
    - dns01:
        cloudflare:
          email: accounts@giantswarm.io
          apiTokenSecretRef:
            name: cloudflare-api-token-secret
            key: api-token
    {{ end }}
    {{ if .Values.acme.dns01.route53.enabled -}}
    - dns01:
        route53:
          region: {{ .Values.acme.dns01.route53.region }}
          role: {{ .Values.acme.dns01.route53.role }}
          {{- if .Values.acme.dns01.route53.hostedZoneID }}
          hostedZoneID: {{ .Values.acme.dns01.route53.hostedZoneID }}
          {{- end }}
          {{- if .Values.acme.dns01.route53.accessKeyID }}
          accessKeyID: {{ .Values.acme.dns01.route53.accessKeyID }}
          {{- end }}
          {{- if .Values.acme.dns01.route53.secretAccessKey }}
          secretAccessKeySecretRef:
            name: route53-access-key-secret
            key: secret-access-key
          {{- end }}
    {{ end }}
    {{ if .Values.acme.dns01.azureDNS.enabled -}}
    - dns01:
        azureDNS:
          hostedZoneName: {{ .Values.acme.dns01.azureDNS.zoneName }}
          resourceGroupName: {{ .Values.acme.dns01.azureDNS.resourceGroupName }}
          subscriptionID: {{ .Values.acme.dns01.azureDNS.subscriptionID }}
          {{- if .Values.acme.dns01.azureDNS.tenantID }}
          tenantID: {{ .Values.acme.dns01.azureDNS.tenantID }}
          {{- end }}
          environment: {{ .Values.acme.dns01.azureDNS.environment }}
          {{- if .Values.acme.dns01.azureDNS.identityClientID }}
          managedIdentity:
            clientID: {{ .Values.acme.dns01.azureDNS.identityClientID }}
          {{- end }}
          {{- if .Values.acme.dns01.azureDNS.clientID }}
          clientID: {{ .Values.acme.dns01.azureDNS.clientID }}
          {{- end }}
          {{- if .Values.acme.dns01.azureDNS.clientSecret }}
          clientSecretSecretRef:
            name: azuredns-config
            key: client-secret
          {{- end }}
    {{ end }}
    {{ if .Values.acme.http01.enabled -}}
    - http01:
        ingress:
          ingressClassName: {{ .Values.acme.http01.ingressClassName }}
    {{ end }}
---
{{- end }}
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: selfsigned-giantswarm
  labels:
    giantswarm.io/service-type: "managed"
spec:
  selfSigned: {}
{{- end }}
