package mctest

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

func TestMCBasic(t *testing.T) {
	suite.New().
		WithInstallNamespace("kube-system").
		WithIsUpgrade(isUpgrade).
		WithValuesFile("./values.yaml").
		AfterClusterReady(func() {
			// no
		}).
		Tests(func() {
			It("should have all cert-manager deployments ready in MC", func() {
				wcClient, _ := state.GetFramework().WC(state.GetCluster().Name)
				appNamespace := "kube-system"

				deploymentNames := []string{
					"cert-manager",
					"cert-manager-cainjector",
					"cert-manager-webhook",
				}

				for _, deploymentName := range deploymentNames {
					logger.Log("Checking MC deployment: %s/%s", appNamespace, deploymentName)

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

			It("should have cluster issuers ready in MC", func() {
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
					logger.Log("Checking MC ClusterIssuer: %s", issuerName)

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

			It("should issue certificates in MC for management workloads", func() {
				wcClient, _ := state.GetFramework().WC(state.GetCluster().Name)
				appNamespace := "kube-system"
				testCertName := "mc-test-cert"
				testSecretName := "mc-test-secret"

				certificateGVK := schema.GroupVersionKind{
					Group:   "cert-manager.io",
					Version: "v1",
					Kind:    "Certificate",
				}

				By("Creating a test certificate in MC")
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
							"commonName": "mc-test.giantswarm.io",
							"dnsNames":   []string{"mc-test.giantswarm.io"},
						},
					},
				}
				certificate.SetGroupVersionKind(certificateGVK)

				err := wcClient.Create(context.Background(), certificate)
				Expect(err).ToNot(HaveOccurred())

				By("Waiting for MC certificate to be ready")
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

			By("Cleaning up MC test certificate")
			cert := &unstructured.Unstructured{}
			cert.SetGroupVersionKind(certificateGVK)
			cert.SetName("mc-test-cert")
			cert.SetNamespace(appNamespace)
			_ = wcClient.Delete(context.Background(), cert)

			By("Cleaning up MC test secret")
			secret := &corev1.Secret{}
			secret.Name = "mc-test-secret"
			secret.Namespace = appNamespace
			_ = wcClient.Delete(context.Background(), secret)
		}).
		Run(t, "MC Basic Test")
}

