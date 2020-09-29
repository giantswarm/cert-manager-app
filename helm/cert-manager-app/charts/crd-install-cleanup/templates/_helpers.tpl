{{- define "CRDCleanupAnnotations" -}}
"helm.sh/hook": "post-install,post-upgrade"
"helm.sh/hook-delete-policy": "hook-succeeded,hook-failed"
{{- end -}}

