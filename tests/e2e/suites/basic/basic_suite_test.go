package basic

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
	isUpgrade = false
)

func TestBasic(t *testing.T) {
	suite.New().
		WithInstallNamespace("default").
		WithIsUpgrade(isUpgrade).
		WithValuesFile("./values.yaml").
		AfterClusterReady(func() {
			// no
		}).
		Tests(func() {
			It("should have all cert-manager deployments ready", func() {
				wcClient, _ := state.GetFramework().WC(state.GetCluster().Name)
				appNamespace := "default"

				deploymentNames := []string{
					"cert-manager",
					"cert-manager-cainjector",
					"cert-manager-webhook",
				}

				for _, deploymentName := range deploymentNames {
					logger.Log("Checking deployment: %s/%s", appNamespace, deploymentName)

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

			It("should have cluster issuers ready", func() {
				wcClient, _ := state.GetFramework().WC(state.GetCluster().Name)

				clusterIssuerGVR := schema.GroupVersionKind{
					Group:   "cert-manager.io",
					Version: "v1",
					Kind:    "ClusterIssuer",
				}

				expectedIssuers := []string{
					"letsencrypt-giantswarm",
					"selfsigned-giantswarm",
				}

				for _, issuerName := range expectedIssuers {
					logger.Log("Checking ClusterIssuer: %s", issuerName)

					Eventually(func() (bool, error) {
						issuer := &unstructured.Unstructured{}
						issuer.SetGroupVersionKind(clusterIssuerGVR)

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

			It("should issue a self-signed certificate", func() {
				wcClient, _ := state.GetFramework().WC(state.GetCluster().Name)
				appNamespace := "default"
				testCertName := "test-cert"
				testSecretName := "test-cert-secret"

				certificateGVK := schema.GroupVersionKind{
					Group:   "cert-manager.io",
					Version: "v1",
					Kind:    "Certificate",
				}

				By("Creating a test certificate")
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
							"commonName": "test.giantswarm.io",
							"dnsNames":   []string{"test.giantswarm.io"},
						},
					},
				}
				certificate.SetGroupVersionKind(certificateGVK)

				err := wcClient.Create(context.Background(), certificate)
				Expect(err).ToNot(HaveOccurred())

				By("Waiting for certificate to be ready")
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
			})
		}).
		AfterSuite(func() {
			wcClient, _ := state.GetFramework().WC(state.GetCluster().Name)
			appNamespace := "default"

			certificateGVK := schema.GroupVersionKind{
				Group:   "cert-manager.io",
				Version: "v1",
				Kind:    "Certificate",
			}

			By("Cleaning up test certificate")
			cert := &unstructured.Unstructured{}
			cert.SetGroupVersionKind(certificateGVK)
			cert.SetName("test-cert")
			cert.SetNamespace(appNamespace)
			_ = wcClient.Delete(context.Background(), cert)

			By("Cleaning up test secret")
			secret := &corev1.Secret{}
			secret.Name = "test-cert-secret"
			secret.Namespace = appNamespace
			_ = wcClient.Delete(context.Background(), secret)
		}).
		Run(t, "Basic Test")
}
