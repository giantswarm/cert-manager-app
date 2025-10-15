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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
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
			// Hook for actions after cluster is ready
		}).
		Tests(func() {
			It("should have all required deployments ready", func() {
				deploymentNames := []string{
					"cert-manager-app",
					"cert-manager-app-cainjector",
					"cert-manager-app-webhook",
				}

				wcClient, _ := state.GetFramework().WC(state.GetCluster().Name)
				appNamespace := "default"

				for _, deploymentName := range deploymentNames {
					logger.Log("Checking deployment: %s/%s", appNamespace, deploymentName)

					Eventually(func() (int32, error) {
						deployment, err := wcClient.AppsV1().Deployments(appNamespace).Get(context.Background(), deploymentName, metav1.GetOptions{})
						if err != nil {
							return 0, err
						}
						return deployment.Status.ReadyReplicas, nil
					}, 10*time.Minute, 15*time.Second).Should(BeNumerically(">", 0))
				}
			})

			It("should have all required cluster issuers created", func() {
				expectedIssuers := []string{
					"letsencrypt-giantswarm",
					"selfsigned-giantswarm",
				}

				wcClient, _ := state.GetFramework().WC(state.GetCluster().Name)

				clusterIssuerGVR := schema.GroupVersionResource{
					Group:    "cert-manager.io",
					Version:  "v1",
					Resource: "clusterissuers",
				}

				for _, issuerName := range expectedIssuers {
					logger.Log("Checking ClusterIssuer: %s", issuerName)

					Eventually(func() error {
						_, err := wcClient.Dynamic().Resource(clusterIssuerGVR).Get(context.Background(), issuerName, metav1.GetOptions{})
						return err
					}, 5*time.Minute, 15*time.Second).Should(Succeed())
				}
			})

			It("should have cert-manager pods running", func() {
				wcClient, _ := state.GetFramework().WC(state.GetCluster().Name)
				appNamespace := "default"

				logger.Log("Checking cert-manager pods in namespace: %s", appNamespace)

				pods, err := wcClient.CoreV1().Pods(appNamespace).List(context.Background(), metav1.ListOptions{
					LabelSelector: "app.kubernetes.io/name=cert-manager",
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(len(pods.Items)).To(BeNumerically(">=", 1))

				for _, pod := range pods.Items {
					Expect(pod.Status.Phase).To(Equal(corev1.PodPhase("Running")))
				}
			})

			It("should be able to issue a self-signed certificate", func() {
				wcClient, _ := state.GetFramework().WC(state.GetCluster().Name)
				appNamespace := "default"
				testCertName := "test-selfsigned-cert"
				testSecretName := "test-selfsigned-secret"

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
							"commonName": "test.example.com",
							"dnsNames":   []string{"test.example.com", "www.test.example.com"},
						},
					},
				}

				certificateGVR := schema.GroupVersionResource{
					Group:    "cert-manager.io",
					Version:  "v1",
					Resource: "certificates",
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
				}, 3*time.Minute, 10*time.Second).Should(BeTrue())

				By("Verifying the secret was created")
				Eventually(func() error {
					_, err := wcClient.CoreV1().Secrets(appNamespace).Get(context.Background(), testSecretName, metav1.GetOptions{})
					return err
				}, 1*time.Minute, 5*time.Second).Should(Succeed())

				By("Cleaning up test resources")
				_ = wcClient.Dynamic().Resource(certificateGVR).Namespace(appNamespace).Delete(context.Background(), testCertName, metav1.DeleteOptions{})
				_ = wcClient.CoreV1().Secrets(appNamespace).Delete(context.Background(), testSecretName, metav1.DeleteOptions{})
			})
		}).
		Run(t, "cert-manager-app Basic Suite")
}
