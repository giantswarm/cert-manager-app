{{- define "issuerLabels" -}}
app.kubernetes.io/name: {{ .Values.name }}
giantswarm.io/service-type: "managed"
{{- end -}}

{{- define "issuerAnnotations" -}}
helm.sh/hook: post-install,post-upgrade
helm.sh/hook-delete-policy: before-hook-creation,hook-succeeded,hook-failed
{{- end -}}

{{- define "clusterIssuer" }}
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-giantswarm
  labels:
    giantswarm.io/service-type: "managed"
spec:
  acme:
    # The ACME server URL.
    server: https://acme-v02.api.letsencrypt.org/directory
    # Email address used for ACME registration.
    email: accounts@giantswarm.io
    # Secret resource used to store the account's private key.
    privateKeySecretRef:
      name: letsencrypt-giantswarm
    # Add a single challenge solver, HTTP01 using nginx.
    solvers:
    {{ if eq .Values.global.acmeSolver.type "dns01" -}}
    - dns01:
        {{ if eq .Values.global.acmeSolver.provider "cloudflare" -}}
        cloudflare:
          email: accounts@giantswarm.io
          apiTokenSecretRef:
            name: cloudflare-api-token-secret
            key: api-token
        {{ end -}}
        {{ if eq .Values.global.acmeSolver.provider "route53" -}}
        route53:
          region: {{ .Values.global.acmeSolver.secret.route53.region }}
          accessKeyID: {{ .Values.global.acmeSolver.secret.route53.accessKeyID }}
          secretAccessKeySecretRef:
            name: route53-access-key-secret
            key: secret-access-key
        {{ end -}}
    {{ else -}}
    - http01:
        ingress:
          class: nginx
    {{ end }}
---
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: selfsigned-giantswarm
  labels:
    giantswarm.io/service-type: "managed"
spec:
  selfSigned: {}
{{- end }}
