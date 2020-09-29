{{- define "CRDCleanupAnnotations" -}}
"helm.sh/hook": "post-install,post-upgrade"
"helm.sh/hook-delete-policy": "before-hook-creation,hook-succeeded,hook-failed"
{{- end -}}

