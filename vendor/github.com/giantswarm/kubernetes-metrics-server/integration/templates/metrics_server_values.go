// +build k8srequired

package templates

// MetricsServerValues values required by kubernetes-metrics-server-chart.
const MetricsServerValues = `---
name: metrics-server
namespace: kube-system
serviceType: managed
k8sAppLabel: metrics-server
rbac:
  create: true
serviceAccount:
  create: true
  name: metrics-server
apiService:
  create: true
  insecureSkipTLSVerify: true
image:
  repository: quay.io/giantswarm/metrics-server-amd64
  tag: v0.3.1
  pullPolicy: IfNotPresent
args:
  - --logtostderr
  - --kubelet-insecure-tls
  - --cert-dir=/tmp
  - --secure-port=10443
resources: {}
nodeSelector: {}
tolerations: []
e2e: true
`
