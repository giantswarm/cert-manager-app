package mcupgrade

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

func TestMCUpgrade(t *testing.T) {
	var preUpgradeCertName = "mc-pre-upgrade-cert"
	var preUpgradeCertSecret = "mc-pre-upgrade-secret"

	suite.New().
		WithInstallNamespace("kube-system").
		WithIsUpgrade(isUpgrade).
		WithValuesFile("./values.yaml").
		AfterClusterReady(func() {
			// Create a certificate before the upgrade to ensure it persists
			wcClient, _ := state.GetFramework().WC(state.GetCluster().Name)
			appNamespace := "kube-system"

			logger.Log("Creating pre-upgrade certificate in MC for upgrade testing")

			certificateGVK := schema.GroupVersionKind{
				Group:   "cert-manager.io",
				Version: "v1",
				Kind:    "Certificate",
			}

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
						"commonName": "mc-pre-upgrade.giantswarm.io",
						"dnsNames":   []string{"mc-pre-upgrade.giantswarm.io"},
					},
				},
			}
			certificate.SetGroupVersionKind(certificateGVK)

			err := wcClient.Create(context.Background(), certificate)
			if err != nil {
				logger.Log("Warning: Failed to create pre-upgrade certificate in MC: %v", err)
			} else {
				logger.Log("Pre-upgrade certificate created in MC successfully")

				// Wait for it to be ready
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

				logger.Log("Pre-upgrade certificate in MC is ready")
			}
		}).
		BeforeUpgrade(func() {
			wcClient, _ := state.GetFramework().WC(state.GetCluster().Name)
			appNamespace := "kube-system"

			logger.Log("Verifying MC cert-manager is stable before upgrade")

			deploymentNames := []string{
				"cert-manager",
				"cert-manager-cainjector",
				"cert-manager-webhook",
			}

			for _, deploymentName := range deploymentNames {
				deployment := &appsv1.Deployment{}
				err := wcClient.Get(state.GetContext(), types.NamespacedName{Name: deploymentName, Namespace: appNamespace}, deployment)
				Expect(err).ToNot(HaveOccurred())
				Expect(deployment.Status.ReadyReplicas).To(BeNumerically(">", 0))
				logger.Log("MC deployment %s is ready with %d replicas", deploymentName, deployment.Status.ReadyReplicas)
			}

			logger.Log("All MC cert-manager components are stable, ready to upgrade")
		}).
		Tests(func() {
			It("should successfully upgrade cert-manager in MC", func() {
				wcClient, _ := state.GetFramework().WC(state.GetCluster().Name)
				appNamespace := "kube-system"

				By("Verifying all deployments are ready after MC upgrade")
				deploymentNames := []string{
					"cert-manager",
					"cert-manager-cainjector",
					"cert-manager-webhook",
				}

				for _, deploymentName := range deploymentNames {
					logger.Log("Checking MC deployment after upgrade: %s/%s", appNamespace, deploymentName)

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

			It("should have cluster issuers still ready in MC after upgrade", func() {
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
					logger.Log("Checking MC ClusterIssuer after upgrade: %s", issuerName)

					Eventually(func() (bool, error) {
						issuer := &unstructured.Unstructured{}
						issuer.SetGroupVersionKind(clusterIssuerGVK)

						err := wcClient.Get(state.GetContext(), types.NamespacedName{Name: issuerName}, issuer)
						if err != nil {
							return false, err
						}

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

			It("should reconcile existing MC certificates after upgrade", func() {
				wcClient, _ := state.GetFramework().WC(state.GetCluster().Name)
				appNamespace := "kube-system"

				certificateGVK := schema.GroupVersionKind{
					Group:   "cert-manager.io",
					Version: "v1",
					Kind:    "Certificate",
				}

				By("Verifying pre-upgrade MC certificate still exists and is ready")
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

				By("Verifying the MC secret still exists")
				Eventually(func() error {
					secret := &corev1.Secret{}
					return wcClient.Get(state.GetContext(), types.NamespacedName{Name: preUpgradeCertSecret, Namespace: appNamespace}, secret)
				}).
					WithTimeout(1 * time.Minute).
					WithPolling(5 * time.Second).
					Should(Succeed())

				logger.Log("Pre-upgrade MC certificate is still valid after upgrade")
			})

			It("should issue new certificates in MC after upgrade", func() {
				wcClient, _ := state.GetFramework().WC(state.GetCluster().Name)
				appNamespace := "kube-system"
				testCertName := "mc-post-upgrade-cert"
				testSecretName := "mc-post-upgrade-secret"

				certificateGVK := schema.GroupVersionKind{
					Group:   "cert-manager.io",
					Version: "v1",
					Kind:    "Certificate",
				}

				By("Creating a new MC certificate after upgrade")
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
							"commonName": "mc-post-upgrade.giantswarm.io",
							"dnsNames":   []string{"mc-post-upgrade.giantswarm.io"},
						},
					},
				}
				certificate.SetGroupVersionKind(certificateGVK)

				err := wcClient.Create(context.Background(), certificate)
				Expect(err).ToNot(HaveOccurred())

				By("Waiting for post-upgrade MC certificate to be ready")
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

				By("Verifying the secret was created in MC")
				Eventually(func() error {
					secret := &corev1.Secret{}
					return wcClient.Get(state.GetContext(), types.NamespacedName{Name: testSecretName, Namespace: appNamespace}, secret)
				}).
					WithTimeout(1 * time.Minute).
					WithPolling(5 * time.Second).
					Should(Succeed())

				logger.Log("Successfully issued new certificate in MC after upgrade")
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

			By("Cleaning up pre-upgrade MC certificate")
			cert := &unstructured.Unstructured{}
			cert.SetGroupVersionKind(certificateGVK)
			cert.SetName(preUpgradeCertName)
			cert.SetNamespace(appNamespace)
			_ = wcClient.Delete(context.Background(), cert)

			By("Cleaning up post-upgrade MC certificate")
			cert2 := &unstructured.Unstructured{}
			cert2.SetGroupVersionKind(certificateGVK)
			cert2.SetName("mc-post-upgrade-cert")
			cert2.SetNamespace(appNamespace)
			_ = wcClient.Delete(context.Background(), cert2)

			By("Cleaning up secrets")
			secret1 := &corev1.Secret{}
			secret1.Name = preUpgradeCertSecret
			secret1.Namespace = appNamespace
			_ = wcClient.Delete(context.Background(), secret1)

			secret2 := &corev1.Secret{}
			secret2.Name = "mc-post-upgrade-secret"
			secret2.Namespace = appNamespace
			_ = wcClient.Delete(context.Background(), secret2)
		}).
		Run(t, "MC Upgrade Test")
}

