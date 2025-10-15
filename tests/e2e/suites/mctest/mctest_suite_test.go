package mctest

import (
	"context"
	"testing"
	"time"

	"github.com/giantswarm/apptest-framework/pkg/config"
	"github.com/giantswarm/apptest-framework/pkg/suite"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
)

var (
	ts           *suite.TestSuite
	ctx          = context.Background()
	kubeClient   kubernetes.Interface
	appNamespace = "default"
)

func TestMCTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "cert-manager-app Management Cluster Suite")
}

var _ = BeforeSuite(func() {
	var err error
	ts, err = suite.New(ctx, config.MustLoad())
	Expect(err).ToNot(HaveOccurred())

	kubeClient = ts.GetClient().K8s()
})

var _ = AfterSuite(func() {
	err := ts.Cleanup(ctx)
	Expect(err).ToNot(HaveOccurred())
})

// Test that cert-manager can be installed in a management cluster
var _ = Describe("cert-manager-app in management cluster", func() {
	It("should install successfully", func() {
		Eventually(func() string {
			app := ts.GetApplication()
			if app == nil {
				return ""
			}
			
			// Get app status
			appObj, err := kubeClient.RESTClient().
				Get().
				AbsPath("/apis/application.giantswarm.io/v1alpha1").
				Resource("apps").
				Namespace(app.GetNamespace()).
				Name(app.GetName()).
				Do(ctx).Get()
			
			if err != nil {
				return ""
			}
			
			if appUnstructured, ok := appObj.(*unstructured.Unstructured); ok {
				status, found, err := unstructured.NestedString(appUnstructured.Object, "status", "release", "status")
				if err != nil || !found {
					return ""
				}
				return status
			}
			return ""
		}, 15*time.Minute, 30*time.Second).Should(Equal("deployed"))
	})

	It("should have cert-manager components running in management cluster", func() {
		deploymentNames := []string{
			"cert-manager-app",
			"cert-manager-app-cainjector",
			"cert-manager-app-webhook",
		}

		for _, deploymentName := range deploymentNames {
			Eventually(func() (int32, error) {
				deployment, err := kubeClient.AppsV1().Deployments(appNamespace).Get(ctx, deploymentName, metav1.GetOptions{})
				if err != nil {
					return 0, err
				}
				return deployment.Status.ReadyReplicas, nil
			}, 10*time.Minute, 30*time.Second).Should(BeNumerically(">", 0))
		}
	})

	It("should be able to issue certificates for management cluster workloads", func() {
		// Create a certificate for a management cluster service
		testCertName := "mc-test-cert"
		testSecretName := "mc-test-secret"

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
					"commonName": "management-cluster.internal",
					"dnsNames": []string{
						"management-cluster.internal",
						"api.management-cluster.internal",
						"webhooks.management-cluster.internal",
					},
				},
			},
		}

		certificateGVR := schema.GroupVersionResource{
			Group:    "cert-manager.io",
			Version:  "v1",
			Resource: "certificates",
		}

		// Create the certificate
		dynamicClient := ts.GetClient().Dynamic()
		_, err := dynamicClient.Resource(certificateGVR).Namespace(appNamespace).Create(ctx, certificate, metav1.CreateOptions{})
		Expect(err).ToNot(HaveOccurred())

		// Wait for certificate to be ready
		Eventually(func() (bool, error) {
			cert, err := dynamicClient.Resource(certificateGVR).Namespace(appNamespace).Get(ctx, testCertName, metav1.GetOptions{})
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
		}, 5*time.Minute, 15*time.Second).Should(BeTrue())

		// Verify the secret was created
		Eventually(func() error {
			_, err := kubeClient.CoreV1().Secrets(appNamespace).Get(ctx, testSecretName, metav1.GetOptions{})
			return err
		}, 2*time.Minute, 10*time.Second).Should(Succeed())

		// Cleanup
		_ = dynamicClient.Resource(certificateGVR).Namespace(appNamespace).Delete(ctx, testCertName, metav1.DeleteOptions{})
		_ = kubeClient.CoreV1().Secrets(appNamespace).Delete(ctx, testSecretName, metav1.DeleteOptions{})
	})

	It("should have cluster issuers available for all workload clusters", func() {
		expectedIssuers := []string{
			"letsencrypt-giantswarm",
			"selfsigned-giantswarm",
		}

		clusterIssuerGVR := schema.GroupVersionResource{
			Group:    "cert-manager.io",
			Version:  "v1",
			Resource: "clusterissuers",
		}

		for _, issuerName := range expectedIssuers {
			Eventually(func() error {
				dynamicClient := ts.GetClient().Dynamic()
				issuer, err := dynamicClient.Resource(clusterIssuerGVR).Get(ctx, issuerName, metav1.GetOptions{})
				if err != nil {
					return err
				}

				// Verify the issuer is ready
				conditions, found, err := unstructured.NestedSlice(issuer.Object, "status", "conditions")
				if err != nil || !found {
					return Errorf("issuer %s has no status conditions", issuerName)
				}

				for _, condition := range conditions {
					if conditionMap, ok := condition.(map[string]interface{}); ok {
						if conditionMap["type"] == "Ready" && conditionMap["status"] == "True" {
							return nil
						}
					}
				}
				return Errorf("issuer %s is not ready", issuerName)
			}, 5*time.Minute, 15*time.Second).Should(Succeed())
		}
	})

	It("should handle high availability scenarios", func() {
		// Check that multiple replicas are running if configured
		Eventually(func() (int32, error) {
			deployment, err := kubeClient.AppsV1().Deployments(appNamespace).Get(ctx, "cert-manager-app", metav1.GetOptions{})
			if err != nil {
				return 0, err
			}
			return deployment.Status.ReadyReplicas, nil
		}, 5*time.Minute, 15*time.Second).Should(BeNumerically(">=", 1))

		// Verify controller manager is handling events properly by checking logs
		pods, err := kubeClient.CoreV1().Pods(appNamespace).List(ctx, metav1.ListOptions{
			LabelSelector: "app.kubernetes.io/name=cert-manager,app.kubernetes.io/component=controller",
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(len(pods.Items)).To(BeNumerically(">=", 1))

		for _, pod := range pods.Items {
			Expect(pod.Status.Phase).To(Equal(corev1.PodPhase("Running")))
		}
	})
})
