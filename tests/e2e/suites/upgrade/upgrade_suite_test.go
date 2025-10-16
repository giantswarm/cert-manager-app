package upgrade

import (
	"context"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/giantswarm/apptest-framework/pkg/state"
	"github.com/giantswarm/apptest-framework/pkg/suite"
	"github.com/giantswarm/clustertest/pkg/logger"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

const (
	isUpgrade = true
)

func TestUpgrade(t *testing.T) {
	var preUpgradeCertName = "pre-upgrade-test-cert"
	var preUpgradeCertSecret = "pre-upgrade-test-secret"

	suite.New().
		WithInstallNamespace("kube-system").
		WithIsUpgrade(isUpgrade).
		WithValuesFile("./values.yaml").
		AfterClusterReady(func() {
			// no
		}).
		BeforeUpgrade(func() {
			// E.g. ensure that the initial install has completed and has settled before upgrading
		}).
		Tests(func() {
			It("should create a pre-upgrade certificate", func() {
				wcClient, _ := state.GetFramework().WC(state.GetCluster().Name)
				appNamespace := "kube-system"

				certificateGVK := schema.GroupVersionKind{
					Group:   "cert-manager.io",
					Version: "v1",
					Kind:    "Certificate",
				}

				By("Creating pre-upgrade certificate")
				certificate := &unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": "cert-manager.io/v1",
						"kind":       "Certificate",
						"metadata": map[string]interface{}{
							"name":      preUpgradeCertName,
							"namespace": appNamespace,
						},
						"spec": map[string]interface{}{
							"secretName": preUpgradeCertSecret,
							"issuerRef": map[string]interface{}{
								"name": "selfsigned-giantswarm",
								"kind": "ClusterIssuer",
							},
							"commonName": "pre-upgrade.giantswarm.io",
							"dnsNames":   []string{"pre-upgrade.giantswarm.io"},
						},
					},
				}
				certificate.SetGroupVersionKind(certificateGVK)

				err := wcClient.Create(context.Background(), certificate)
				Expect(err).ToNot(HaveOccurred())

				By("Waiting for pre-upgrade certificate to be ready")
				Eventually(func() (bool, error) {
					cert := &unstructured.Unstructured{}
					cert.SetGroupVersionKind(certificateGVK)

					err := wcClient.Get(state.GetContext(), types.NamespacedName{Name: preUpgradeCertName, Namespace: appNamespace}, cert)
					if err != nil {
						return false, err
					}

					conditions, found, err := unstructured.NestedSlice(cert.Object, "status", "conditions")
					if err != nil || !found {
						return false, nil
					}

					for _, condition := range conditions {
						if conditionMap, ok := condition.(map[string]interface{}); ok {
							if conditionMap["type"] == "Ready" && conditionMap["status"] == "True" {
								return true, nil
							}
						}
					}
					return false, nil
				}).
					WithTimeout(3 * time.Minute).
					WithPolling(10 * time.Second).
					Should(BeTrue())

				logger.Log("Pre-upgrade certificate is ready")
			})

			It("should successfully upgrade cert-manager", func() {
				wcClient, _ := state.GetFramework().WC(state.GetCluster().Name)
				appNamespace := "kube-system"

				By("Verifying all deployments are ready after upgrade")
				deploymentNames := []string{
					"cert-manager",
					"cert-manager-cainjector",
					"cert-manager-webhook",
				}

				for _, deploymentName := range deploymentNames {
					logger.Log("Checking deployment after upgrade: %s/%s", appNamespace, deploymentName)

					Eventually(func() (int32, error) {
						deployment := &appsv1.Deployment{}
						err := wcClient.Get(state.GetContext(), types.NamespacedName{Name: deploymentName, Namespace: appNamespace}, deployment)
						if err != nil {
							return 0, err
						}
						return deployment.Status.ReadyReplicas, nil
					}).
						WithTimeout(10 * time.Minute).
						WithPolling(15 * time.Second).
						Should(BeNumerically(">", 0))
				}
			})

			It("should have cluster issuers still ready after upgrade", func() {
				wcClient, _ := state.GetFramework().WC(state.GetCluster().Name)

				clusterIssuerGVK := schema.GroupVersionKind{
					Group:   "cert-manager.io",
					Version: "v1",
					Kind:    "ClusterIssuer",
				}

				expectedIssuers := []string{
					"letsencrypt-giantswarm",
					"selfsigned-giantswarm",
				}

				for _, issuerName := range expectedIssuers {
					logger.Log("Checking ClusterIssuer after upgrade: %s", issuerName)

					Eventually(func() (bool, error) {
						issuer := &unstructured.Unstructured{}
						issuer.SetGroupVersionKind(clusterIssuerGVK)

						err := wcClient.Get(state.GetContext(), types.NamespacedName{Name: issuerName}, issuer)
						if err != nil {
							return false, err
						}

						// Check if ready
						conditions, found, err := unstructured.NestedSlice(issuer.Object, "status", "conditions")
						if err != nil || !found {
							return false, nil
						}

						for _, condition := range conditions {
							if conditionMap, ok := condition.(map[string]interface{}); ok {
								if conditionMap["type"] == "Ready" && conditionMap["status"] == "True" {
									return true, nil
								}
							}
						}
						return false, nil
					}).
						WithTimeout(5 * time.Minute).
						WithPolling(15 * time.Second).
						Should(BeTrue())
				}
			})

			It("should reconcile existing certificates after upgrade", func() {
				wcClient, _ := state.GetFramework().WC(state.GetCluster().Name)
				appNamespace := "kube-system"

				certificateGVK := schema.GroupVersionKind{
					Group:   "cert-manager.io",
					Version: "v1",
					Kind:    "Certificate",
				}

				By("Verifying pre-upgrade certificate still exists and is ready")
				Eventually(func() (bool, error) {
					cert := &unstructured.Unstructured{}
					cert.SetGroupVersionKind(certificateGVK)

					err := wcClient.Get(state.GetContext(), types.NamespacedName{Name: preUpgradeCertName, Namespace: appNamespace}, cert)
					if err != nil {
						return false, err
					}

					conditions, found, err := unstructured.NestedSlice(cert.Object, "status", "conditions")
					if err != nil || !found {
						return false, nil
					}

					for _, condition := range conditions {
						if conditionMap, ok := condition.(map[string]interface{}); ok {
							if conditionMap["type"] == "Ready" && conditionMap["status"] == "True" {
								return true, nil
							}
						}
					}
					return false, nil
				}).
					WithTimeout(3 * time.Minute).
					WithPolling(10 * time.Second).
					Should(BeTrue())

				By("Verifying the secret still exists")
				Eventually(func() error {
					secret := &corev1.Secret{}
					return wcClient.Get(state.GetContext(), types.NamespacedName{Name: preUpgradeCertSecret, Namespace: appNamespace}, secret)
				}).
					WithTimeout(1 * time.Minute).
					WithPolling(5 * time.Second).
					Should(Succeed())

				logger.Log("Pre-upgrade certificate is still valid after upgrade")
			})

			It("should issue new certificates after upgrade", func() {
				wcClient, _ := state.GetFramework().WC(state.GetCluster().Name)
				appNamespace := "kube-system"
				testCertName := "post-upgrade-cert"
				testSecretName := "post-upgrade-secret"

				certificateGVK := schema.GroupVersionKind{
					Group:   "cert-manager.io",
					Version: "v1",
					Kind:    "Certificate",
				}

				By("Creating a new certificate after upgrade")
				certificate := &unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": "cert-manager.io/v1",
						"kind":       "Certificate",
						"metadata": map[string]interface{}{
							"name":      testCertName,
							"namespace": appNamespace,
						},
						"spec": map[string]interface{}{
							"secretName": testSecretName,
							"issuerRef": map[string]interface{}{
								"name": "selfsigned-giantswarm",
								"kind": "ClusterIssuer",
							},
							"commonName": "post-upgrade.giantswarm.io",
							"dnsNames":   []string{"post-upgrade.giantswarm.io"},
						},
					},
				}
				certificate.SetGroupVersionKind(certificateGVK)

				err := wcClient.Create(context.Background(), certificate)
				Expect(err).ToNot(HaveOccurred())

				By("Waiting for post-upgrade certificate to be ready")
				Eventually(func() (bool, error) {
					cert := &unstructured.Unstructured{}
					cert.SetGroupVersionKind(certificateGVK)

					err := wcClient.Get(state.GetContext(), types.NamespacedName{Name: testCertName, Namespace: appNamespace}, cert)
					if err != nil {
						return false, err
					}

					conditions, found, err := unstructured.NestedSlice(cert.Object, "status", "conditions")
					if err != nil || !found {
						return false, nil
					}

					for _, condition := range conditions {
						if conditionMap, ok := condition.(map[string]interface{}); ok {
							if conditionMap["type"] == "Ready" && conditionMap["status"] == "True" {
								return true, nil
							}
						}
					}
					return false, nil
				}).
					WithTimeout(3 * time.Minute).
					WithPolling(10 * time.Second).
					Should(BeTrue())

				By("Verifying the secret was created")
				Eventually(func() error {
					secret := &corev1.Secret{}
					return wcClient.Get(state.GetContext(), types.NamespacedName{Name: testSecretName, Namespace: appNamespace}, secret)
				}).
					WithTimeout(1 * time.Minute).
					WithPolling(5 * time.Second).
					Should(Succeed())

				logger.Log("Successfully issued new certificate after upgrade")
			})
		}).
		AfterSuite(func() {
			wcClient, _ := state.GetFramework().WC(state.GetCluster().Name)
			appNamespace := "kube-system"

			certificateGVK := schema.GroupVersionKind{
				Group:   "cert-manager.io",
				Version: "v1",
				Kind:    "Certificate",
			}

			By("Cleaning up pre-upgrade certificate")
			cert := &unstructured.Unstructured{}
			cert.SetGroupVersionKind(certificateGVK)
			cert.SetName(preUpgradeCertName)
			cert.SetNamespace(appNamespace)
			_ = wcClient.Delete(context.Background(), cert)

			By("Cleaning up post-upgrade certificate")
			cert2 := &unstructured.Unstructured{}
			cert2.SetGroupVersionKind(certificateGVK)
			cert2.SetName("post-upgrade-cert")
			cert2.SetNamespace(appNamespace)
			_ = wcClient.Delete(context.Background(), cert2)

			By("Cleaning up secrets")
			secret1 := &corev1.Secret{}
			secret1.Name = preUpgradeCertSecret
			secret1.Namespace = appNamespace
			_ = wcClient.Delete(context.Background(), secret1)

			secret2 := &corev1.Secret{}
			secret2.Name = "post-upgrade-secret"
			secret2.Namespace = appNamespace
			_ = wcClient.Delete(context.Background(), secret2)
		}).
		Run(t, "Upgrade Test")
}

