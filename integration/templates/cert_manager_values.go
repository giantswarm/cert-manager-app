// +build k8srequired

package templates

// CertManagerValues values required by cert-manager-app chart.
const CertManagerValues = `---
name: cert-manager
namespace: kube-system

userID: 1000
groupID: 1000

image:
  registry: quay.io
  name: giantswarm/cert-manager-controller
  tag: v0.9.0
e2e: true
`
