# cert-manager-app

![Version: 3.8.2](https://img.shields.io/badge/Version-3.8.2-informational?style=flat-square) ![AppVersion: v1.16.2](https://img.shields.io/badge/AppVersion-v1.16.2-informational?style=flat-square)

Simplifies the process of obtaining, renewing and using certificates.

**Homepage:** <https://github.com/giantswarm/cert-manager-app>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Shield |  |  |

## Source Code

* <https://github.com/cert-manager/cert-manager>

## Requirements

Kubernetes: `>=1.22.0-0`

| Repository | Name | Version |
|------------|------|---------|
|  | cert-manager | 1.16.2 |
|  | cert-manager-giantswarm-ciliumnetworkpolicies | 0.1.0 |
|  | giantSwarmClusterIssuer(cert-manager-giantswarm-clusterissuer) | 2.0.0 |
|  | cert-manager-giantswarm-netpol | 0.1.0 |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| acmesolver.image.pullPolicy | string | `"IfNotPresent"` |  |
| acmesolver.image.registry | string | `"gsoci.azurecr.io"` |  |
| acmesolver.image.repository | string | `"giantswarm/cert-manager-acmesolver"` |  |
| affinity | object | `{}` |  |
| approveSignerNames[0] | string | `"issuers.cert-manager.io/*"` |  |
| approveSignerNames[1] | string | `"clusterissuers.cert-manager.io/*"` |  |
| cainjector.affinity | object | `{}` |  |
| cainjector.config | object | `{}` |  |
| cainjector.containerSecurityContext.allowPrivilegeEscalation | bool | `false` |  |
| cainjector.containerSecurityContext.capabilities.drop[0] | string | `"ALL"` |  |
| cainjector.containerSecurityContext.readOnlyRootFilesystem | bool | `true` |  |
| cainjector.containerSecurityContext.runAsNonRoot | bool | `true` |  |
| cainjector.enableServiceLinks | bool | `false` |  |
| cainjector.enabled | bool | `true` |  |
| cainjector.extraArgs | list | `[]` |  |
| cainjector.extraEnv | list | `[]` |  |
| cainjector.featureGates | string | `""` |  |
| cainjector.image.pullPolicy | string | `"IfNotPresent"` |  |
| cainjector.image.registry | string | `"gsoci.azurecr.io"` |  |
| cainjector.image.repository | string | `"giantswarm/cert-manager-cainjector"` |  |
| cainjector.nodeSelector."kubernetes.io/os" | string | `"linux"` |  |
| cainjector.podAnnotations."cluster-autoscaler.kubernetes.io/safe-to-evict" | string | `"true"` |  |
| cainjector.podDisruptionBudget.enabled | bool | `false` |  |
| cainjector.podDisruptionBudget.minAvailable | string | `"50%"` |  |
| cainjector.podLabels | object | `{}` |  |
| cainjector.replicaCount | int | `1` |  |
| cainjector.resources.limits.cpu | string | `"100m"` |  |
| cainjector.resources.limits.memory | string | `"1Gi"` |  |
| cainjector.resources.requests.cpu | string | `"20m"` |  |
| cainjector.resources.requests.memory | string | `"64Mi"` |  |
| cainjector.securityContext.runAsGroup | int | `1000` |  |
| cainjector.securityContext.runAsNonRoot | bool | `true` |  |
| cainjector.securityContext.runAsUser | int | `1000` |  |
| cainjector.securityContext.seccompProfile.type | string | `"RuntimeDefault"` |  |
| cainjector.serviceAccount.automountServiceAccountToken | bool | `true` |  |
| cainjector.serviceAccount.create | bool | `true` |  |
| cainjector.serviceLabels | object | `{}` |  |
| cainjector.strategy | object | `{}` |  |
| cainjector.tolerations | list | `[]` |  |
| cainjector.topologySpreadConstraints | list | `[]` |  |
| cainjector.verticalPodAutoscaler.controlledValues | string | `"RequestsAndLimits"` |  |
| cainjector.verticalPodAutoscaler.enabled | bool | `true` |  |
| cainjector.verticalPodAutoscaler.mode | string | `"Auto"` |  |
| cainjector.verticalPodAutoscaler.updatePolicy.updateMode | string | `"Auto"` |  |
| cainjector.volumeMounts | list | `[]` |  |
| cainjector.volumes | list | `[]` |  |
| ciliumNetworkPolicy.enabled | bool | `false` |  |
| clusterResourceNamespace | string | `""` |  |
| config | object | `{}` |  |
| containerSecurityContext.allowPrivilegeEscalation | bool | `false` |  |
| containerSecurityContext.capabilities.drop[0] | string | `"ALL"` |  |
| containerSecurityContext.readOnlyRootFilesystem | bool | `true` |  |
| containerSecurityContext.runAsNonRoot | bool | `true` |  |
| crds.enabled | bool | `true` |  |
| crds.keep | bool | `true` |  |
| creator | string | `"helm"` |  |
| disableAutoApproval | bool | `false` |  |
| dns01RecursiveNameservers | string | `""` |  |
| dns01RecursiveNameserversOnly | bool | `false` |  |
| enableCertificateOwnerRef | bool | `false` |  |
| enableServiceLinks | bool | `false` |  |
| enabled | bool | `true` |  |
| extraArgs | list | `[]` |  |
| extraEnv | list | `[]` |  |
| extraObjects | list | `[]` |  |
| featureGates | string | `""` |  |
| fullnameOverride | string | `"cert-manager-app"` |  |
| giantswarmNetworkPolicy.enabled | bool | `true` |  |
| global.commonLabels | object | `{}` |  |
| global.imagePullSecrets | list | `[]` |  |
| global.leaderElection.namespace | string | `"kube-system"` |  |
| global.logLevel | int | `2` |  |
| global.podSecurityPolicy.enabled | bool | `false` |  |
| global.podSecurityPolicy.useAppArmor | bool | `true` |  |
| global.podSecurityStandards.enforced | bool | `false` |  |
| global.priorityClassName | string | `""` |  |
| global.rbac.aggregateClusterRoles | bool | `true` |  |
| global.rbac.create | bool | `true` |  |
| hostAliases | list | `[]` |  |
| image.pullPolicy | string | `"IfNotPresent"` |  |
| image.registry | string | `"gsoci.azurecr.io"` |  |
| image.repository | string | `"giantswarm/cert-manager-controller"` |  |
| ingressShim | object | `{}` |  |
| installCRDs | bool | `false` |  |
| livenessProbe.enabled | bool | `false` |  |
| livenessProbe.failureThreshold | int | `8` |  |
| livenessProbe.initialDelaySeconds | int | `10` |  |
| livenessProbe.periodSeconds | int | `10` |  |
| livenessProbe.successThreshold | int | `1` |  |
| livenessProbe.timeoutSeconds | int | `15` |  |
| maxConcurrentChallenges | int | `60` |  |
| namespace | string | `""` |  |
| nodeSelector."kubernetes.io/os" | string | `"linux"` |  |
| podAnnotations."cluster-autoscaler.kubernetes.io/safe-to-evict" | string | `"true"` |  |
| podDisruptionBudget.enabled | bool | `false` |  |
| podLabels | object | `{}` |  |
| prometheus.enabled | bool | `true` |  |
| prometheus.podmonitor.annotations | object | `{}` |  |
| prometheus.podmonitor.enabled | bool | `false` |  |
| prometheus.podmonitor.endpointAdditionalProperties | object | `{}` |  |
| prometheus.podmonitor.honorLabels | bool | `false` |  |
| prometheus.podmonitor.interval | string | `"60s"` |  |
| prometheus.podmonitor.labels | object | `{}` |  |
| prometheus.podmonitor.path | string | `"/metrics"` |  |
| prometheus.podmonitor.prometheusInstance | string | `"default"` |  |
| prometheus.podmonitor.scrapeTimeout | string | `"30s"` |  |
| prometheus.servicemonitor.annotations | object | `{}` |  |
| prometheus.servicemonitor.enabled | bool | `true` |  |
| prometheus.servicemonitor.endpointAdditionalProperties.relabelings[0].action | string | `"replace"` |  |
| prometheus.servicemonitor.endpointAdditionalProperties.relabelings[0].regex | string | `";(.*)"` |  |
| prometheus.servicemonitor.endpointAdditionalProperties.relabelings[0].replacement | string | `"$1"` |  |
| prometheus.servicemonitor.endpointAdditionalProperties.relabelings[0].separator | string | `";"` |  |
| prometheus.servicemonitor.endpointAdditionalProperties.relabelings[0].sourceLabels[0] | string | `"namespace"` |  |
| prometheus.servicemonitor.endpointAdditionalProperties.relabelings[0].sourceLabels[1] | string | `"__meta_kubernetes_namespace"` |  |
| prometheus.servicemonitor.endpointAdditionalProperties.relabelings[0].targetLabel | string | `"namespace"` |  |
| prometheus.servicemonitor.endpointAdditionalProperties.relabelings[1].action | string | `"replace"` |  |
| prometheus.servicemonitor.endpointAdditionalProperties.relabelings[1].sourceLabels[0] | string | `"__meta_kubernetes_pod_label_app"` |  |
| prometheus.servicemonitor.endpointAdditionalProperties.relabelings[1].targetLabel | string | `"app"` |  |
| prometheus.servicemonitor.endpointAdditionalProperties.relabelings[2].action | string | `"replace"` |  |
| prometheus.servicemonitor.endpointAdditionalProperties.relabelings[2].sourceLabels[0] | string | `"__meta_kubernetes_pod_node_name"` |  |
| prometheus.servicemonitor.endpointAdditionalProperties.relabelings[2].targetLabel | string | `"node"` |  |
| prometheus.servicemonitor.honorLabels | bool | `false` |  |
| prometheus.servicemonitor.interval | string | `"60s"` |  |
| prometheus.servicemonitor.labels | object | `{}` |  |
| prometheus.servicemonitor.path | string | `"/metrics"` |  |
| prometheus.servicemonitor.prometheusInstance | string | `"default"` |  |
| prometheus.servicemonitor.scrapeTimeout | string | `"30s"` |  |
| prometheus.servicemonitor.targetPort | int | `9402` |  |
| replicaCount | int | `1` |  |
| resources.limits.cpu | string | `"500m"` |  |
| resources.limits.memory | string | `"1Gi"` |  |
| resources.requests.cpu | string | `"50m"` |  |
| resources.requests.memory | string | `"100Mi"` |  |
| securityContext.runAsGroup | int | `1000` |  |
| securityContext.runAsNonRoot | bool | `true` |  |
| securityContext.runAsUser | int | `1000` |  |
| securityContext.seccompProfile.type | string | `"RuntimeDefault"` |  |
| serviceAccount.automountServiceAccountToken | bool | `true` |  |
| serviceAccount.create | bool | `true` |  |
| startupapicheck.affinity | object | `{}` |  |
| startupapicheck.backoffLimit | int | `4` |  |
| startupapicheck.containerSecurityContext.allowPrivilegeEscalation | bool | `false` |  |
| startupapicheck.containerSecurityContext.capabilities.drop[0] | string | `"ALL"` |  |
| startupapicheck.containerSecurityContext.readOnlyRootFilesystem | bool | `true` |  |
| startupapicheck.containerSecurityContext.runAsNonRoot | bool | `true` |  |
| startupapicheck.enableServiceLinks | bool | `false` |  |
| startupapicheck.enabled | bool | `false` |  |
| startupapicheck.extraArgs[0] | string | `"-v"` |  |
| startupapicheck.extraEnv | list | `[]` |  |
| startupapicheck.image.pullPolicy | string | `"IfNotPresent"` |  |
| startupapicheck.image.registry | string | `"gsoci.azurecr.io"` |  |
| startupapicheck.image.repository | string | `"giantswarm/cert-manager-startupapicheck"` |  |
| startupapicheck.jobAnnotations."helm.sh/hook" | string | `"post-install"` |  |
| startupapicheck.jobAnnotations."helm.sh/hook-delete-policy" | string | `"before-hook-creation,hook-succeeded"` |  |
| startupapicheck.jobAnnotations."helm.sh/hook-weight" | string | `"1"` |  |
| startupapicheck.nodeSelector."kubernetes.io/os" | string | `"linux"` |  |
| startupapicheck.podLabels | object | `{}` |  |
| startupapicheck.rbac.annotations."helm.sh/hook" | string | `"post-install"` |  |
| startupapicheck.rbac.annotations."helm.sh/hook-delete-policy" | string | `"before-hook-creation,hook-succeeded"` |  |
| startupapicheck.rbac.annotations."helm.sh/hook-weight" | string | `"-5"` |  |
| startupapicheck.resources.requests.cpu | string | `"20m"` |  |
| startupapicheck.resources.requests.memory | string | `"64Mi"` |  |
| startupapicheck.securityContext.runAsNonRoot | bool | `true` |  |
| startupapicheck.securityContext.seccompProfile.type | string | `"RuntimeDefault"` |  |
| startupapicheck.serviceAccount.annotations."helm.sh/hook" | string | `"post-install"` |  |
| startupapicheck.serviceAccount.annotations."helm.sh/hook-delete-policy" | string | `"before-hook-creation,hook-succeeded"` |  |
| startupapicheck.serviceAccount.annotations."helm.sh/hook-weight" | string | `"-5"` |  |
| startupapicheck.serviceAccount.automountServiceAccountToken | bool | `true` |  |
| startupapicheck.serviceAccount.create | bool | `true` |  |
| startupapicheck.timeout | string | `"1m"` |  |
| startupapicheck.tolerations[0].effect | string | `"NoSchedule"` |  |
| startupapicheck.tolerations[0].key | string | `"node-role.kubernetes.io/master"` |  |
| startupapicheck.tolerations[1].effect | string | `"NoSchedule"` |  |
| startupapicheck.tolerations[1].key | string | `"node-role.kubernetes.io/control-plane"` |  |
| startupapicheck.volumeMounts | list | `[]` |  |
| startupapicheck.volumes | list | `[]` |  |
| strategy | object | `{}` |  |
| tolerations[0].effect | string | `"NoSchedule"` |  |
| tolerations[0].key | string | `"node-role.kubernetes.io/master"` |  |
| tolerations[1].effect | string | `"NoSchedule"` |  |
| tolerations[1].key | string | `"node-role.kubernetes.io/control-plane"` |  |
| topologySpreadConstraints | list | `[]` |  |
| volumeMounts | list | `[]` |  |
| volumes | list | `[]` |  |
| webhook.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution[0].podAffinityTerm.labelSelector.matchExpressions[0].key | string | `"apps.giantswarm.io/affinity"` |  |
| webhook.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution[0].podAffinityTerm.labelSelector.matchExpressions[0].operator | string | `"In"` |  |
| webhook.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution[0].podAffinityTerm.labelSelector.matchExpressions[0].values[0] | string | `"cert-manager-webhook"` |  |
| webhook.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution[0].podAffinityTerm.labelSelector.matchExpressions[1].key | string | `"app.kubernetes.io/component"` |  |
| webhook.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution[0].podAffinityTerm.labelSelector.matchExpressions[1].operator | string | `"In"` |  |
| webhook.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution[0].podAffinityTerm.labelSelector.matchExpressions[1].values[0] | string | `"webhook"` |  |
| webhook.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution[0].podAffinityTerm.topologyKey | string | `"kubernetes.io/hostname"` |  |
| webhook.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution[0].weight | int | `100` |  |
| webhook.config | object | `{}` |  |
| webhook.containerSecurityContext.allowPrivilegeEscalation | bool | `false` |  |
| webhook.containerSecurityContext.capabilities.drop[0] | string | `"ALL"` |  |
| webhook.containerSecurityContext.readOnlyRootFilesystem | bool | `true` |  |
| webhook.containerSecurityContext.runAsNonRoot | bool | `true` |  |
| webhook.enableServiceLinks | bool | `false` |  |
| webhook.extraArgs | list | `[]` |  |
| webhook.extraEnv | list | `[]` |  |
| webhook.featureGates | string | `""` |  |
| webhook.hostNetwork | bool | `false` |  |
| webhook.image.pullPolicy | string | `"IfNotPresent"` |  |
| webhook.image.registry | string | `"gsoci.azurecr.io"` |  |
| webhook.image.repository | string | `"giantswarm/cert-manager-webhook"` |  |
| webhook.livenessProbe.failureThreshold | int | `3` |  |
| webhook.livenessProbe.initialDelaySeconds | int | `60` |  |
| webhook.livenessProbe.periodSeconds | int | `10` |  |
| webhook.livenessProbe.successThreshold | int | `1` |  |
| webhook.livenessProbe.timeoutSeconds | int | `1` |  |
| webhook.mutatingWebhookConfiguration.namespaceSelector | object | `{}` |  |
| webhook.networkPolicy.egress[0].ports[0].port | int | `80` |  |
| webhook.networkPolicy.egress[0].ports[0].protocol | string | `"TCP"` |  |
| webhook.networkPolicy.egress[0].ports[1].port | int | `443` |  |
| webhook.networkPolicy.egress[0].ports[1].protocol | string | `"TCP"` |  |
| webhook.networkPolicy.egress[0].ports[2].port | int | `53` |  |
| webhook.networkPolicy.egress[0].ports[2].protocol | string | `"TCP"` |  |
| webhook.networkPolicy.egress[0].ports[3].port | int | `53` |  |
| webhook.networkPolicy.egress[0].ports[3].protocol | string | `"UDP"` |  |
| webhook.networkPolicy.egress[0].ports[4].port | int | `6443` |  |
| webhook.networkPolicy.egress[0].ports[4].protocol | string | `"TCP"` |  |
| webhook.networkPolicy.egress[0].to[0].ipBlock.cidr | string | `"0.0.0.0/0"` |  |
| webhook.networkPolicy.enabled | bool | `true` |  |
| webhook.networkPolicy.ingress[0].from[0].ipBlock.cidr | string | `"0.0.0.0/0"` |  |
| webhook.nodeSelector."kubernetes.io/os" | string | `"linux"` |  |
| webhook.podDisruptionBudget.enabled | bool | `true` |  |
| webhook.podDisruptionBudget.minAvailable | string | `"50%"` |  |
| webhook.podLabels."apps.giantswarm.io/affinity" | string | `"cert-manager-webhook"` |  |
| webhook.readinessProbe.failureThreshold | int | `3` |  |
| webhook.readinessProbe.initialDelaySeconds | int | `5` |  |
| webhook.readinessProbe.periodSeconds | int | `5` |  |
| webhook.readinessProbe.successThreshold | int | `1` |  |
| webhook.readinessProbe.timeoutSeconds | int | `1` |  |
| webhook.replicaCount | int | `2` |  |
| webhook.resources.limits.cpu | string | `"100m"` |  |
| webhook.resources.limits.memory | string | `"100Mi"` |  |
| webhook.resources.requests.cpu | string | `"20m"` |  |
| webhook.resources.requests.memory | string | `"50Mi"` |  |
| webhook.securePort | int | `10250` |  |
| webhook.securityContext.runAsGroup | int | `1000` |  |
| webhook.securityContext.runAsNonRoot | bool | `true` |  |
| webhook.securityContext.runAsUser | int | `1000` |  |
| webhook.securityContext.seccompProfile.type | string | `"RuntimeDefault"` |  |
| webhook.serviceAccount.automountServiceAccountToken | bool | `true` |  |
| webhook.serviceAccount.create | bool | `true` |  |
| webhook.serviceIPFamilies | list | `[]` |  |
| webhook.serviceIPFamilyPolicy | string | `""` |  |
| webhook.serviceLabels | object | `{}` |  |
| webhook.serviceType | string | `"ClusterIP"` |  |
| webhook.strategy | object | `{}` |  |
| webhook.timeoutSeconds | int | `30` |  |
| webhook.tolerations[0].effect | string | `"NoSchedule"` |  |
| webhook.tolerations[0].key | string | `"node-role.kubernetes.io/master"` |  |
| webhook.tolerations[1].effect | string | `"NoSchedule"` |  |
| webhook.tolerations[1].key | string | `"node-role.kubernetes.io/control-plane"` |  |
| webhook.topologySpreadConstraints | list | `[]` |  |
| webhook.url | object | `{}` |  |
| webhook.validatingWebhookConfiguration.namespaceSelector.matchExpressions[0].key | string | `"cert-manager.io/disable-validation"` |  |
| webhook.validatingWebhookConfiguration.namespaceSelector.matchExpressions[0].operator | string | `"NotIn"` |  |
| webhook.validatingWebhookConfiguration.namespaceSelector.matchExpressions[0].values[0] | string | `"true"` |  |
| webhook.volumeMounts | list | `[]` |  |
| webhook.volumes | list | `[]` |  |

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.11.0](https://github.com/norwoodj/helm-docs/releases/v1.11.0)
