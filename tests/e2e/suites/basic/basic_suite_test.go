package basic

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/giantswarm/apptest-framework/pkg/state"
	"github.com/giantswarm/apptest-framework/pkg/suite"
	"github.com/giantswarm/clustertest/pkg/logger"

	helmv2beta2 "github.com/fluxcd/helm-controller/api/v2beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
			It("should deploy the HelmRelease", func() {
				Eventually(func() (bool, error) {
					appNamespace := state.GetCluster().Organization.GetNamespace()
					appName := fmt.Sprintf("%s-cert-manager-app", state.GetCluster().Name)

					mcKubeClient := state.GetFramework().MC()

					logger.Log("HelmRelease: %s/%s", appNamespace, appName)

					release := &helmv2beta2.HelmRelease{}
					err := mcKubeClient.Get(state.GetContext(), types.NamespacedName{Name: appName, Namespace: appNamespace}, release)
					if err != nil {
						return false, err
					}

					for _, c := range release.Status.Conditions {
						if c.Type == "Ready" {
							if c.Status == "True" {
								return true, nil
							} else {
								return false, errors.New(fmt.Sprintf("HelmRelease not ready [%s]: %s", c.Reason, c.Message))
							}
						}
					}

					return false, errors.New("HelmRelease not ready")
				}).
					WithTimeout(10 * time.Minute).
					WithPolling(15 * time.Second).
					Should(BeTrue())
			})

			It("should have all cert-manager deployments ready", func() {
				wcClient, _ := state.GetFramework().WC(state.GetCluster().Name)
				appNamespace := "default"

				deploymentNames := []string{
					"cert-manager-app",
					"cert-manager-app-cainjector",
					"cert-manager-app-webhook",
				}

				for _, deploymentName := range deploymentNames {
					logger.Log("Checking deployment: %s/%s", appNamespace, deploymentName)

					Eventually(func() (int32, error) {
						deployment, err := wcClient.AppsV1().Deployments(appNamespace).Get(context.Background(), deploymentName, metav1.GetOptions{})
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

				clusterIssuerGVR := schema.GroupVersionResource{
					Group:    "cert-manager.io",
					Version:  "v1",
					Resource: "clusterissuers",
				}

				expectedIssuers := []string{
					"letsencrypt-giantswarm",
					"selfsigned-giantswarm",
				}

				for _, issuerName := range expectedIssuers {
					logger.Log("Checking ClusterIssuer: %s", issuerName)

					Eventually(func() (bool, error) {
						issuer, err := wcClient.Dynamic().Resource(clusterIssuerGVR).Get(context.Background(), issuerName, metav1.GetOptions{})
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

				certificateGVR := schema.GroupVersionResource{
					Group:    "cert-manager.io",
					Version:  "v1",
					Resource: "certificates",
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

				_, err := wcClient.Dynamic().Resource(certificateGVR).Namespace(appNamespace).Create(context.Background(), certificate, metav1.CreateOptions{})
				Expect(err).ToNot(HaveOccurred())

				By("Waiting for certificate to be ready")
				Eventually(func() (bool, error) {
					cert, err := wcClient.Dynamic().Resource(certificateGVR).Namespace(appNamespace).Get(context.Background(), testCertName, metav1.GetOptions{})
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
					_, err := wcClient.CoreV1().Secrets(appNamespace).Get(context.Background(), testSecretName, metav1.GetOptions{})
					return err
				}).
					WithTimeout(1 * time.Minute).
					WithPolling(5 * time.Second).
					Should(Succeed())
			})
		}).
		AfterSuite(func() {
			wcClient, _ := state.GetFramework().WC(state.GetCluster().Name)
			appNamespace := "default"

			certificateGVR := schema.GroupVersionResource{
				Group:    "cert-manager.io",
				Version:  "v1",
				Resource: "certificates",
			}

			By("Cleaning up test certificate")
			_ = wcClient.Dynamic().Resource(certificateGVR).Namespace(appNamespace).Delete(context.Background(), "test-cert", metav1.DeleteOptions{})

			By("Cleaning up test secret")
			_ = wcClient.CoreV1().Secrets(appNamespace).Delete(context.Background(), "test-cert-secret", metav1.DeleteOptions{})
		}).
		Run(t, "Basic Test")
}
