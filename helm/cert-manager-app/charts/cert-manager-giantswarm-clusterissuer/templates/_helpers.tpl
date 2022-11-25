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
    {{ if eq .Values.global.acmeSolver.type "dns01" }}
    - dns01:
        cloudflare:
          email: accounts@giantswarm.io
          apiTokenSecretRef:
            name: cloudflare-api-token-secret
            key: api-token
    {{ else }}
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
{{- if .Values.global.privateIssuer.enabled }}
---
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: private-giantswarm
  labels:
    giantswarm.io/service-type: "managed"
spec:
  ca:
    secretName: private-giantswarm-secret
---
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: private-giantswarm-ca
  namespace: kube-system
spec:
  isCA: true
  commonName: gigantic.internal
  secretName: private-giantswarm-secret
  privateKey:
    algorithm: ECDSA
    size: 256
  issuerRef:
    name: selfsigned-giantswarm
    kind: ClusterIssuer
    group: cert-manager.io
{{- end }}
{{- end }}
