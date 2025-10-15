package basic

import (
	"context"
	"testing"
	"time"

	"github.com/giantswarm/apptest-framework/pkg/config"
	"github.com/giantswarm/apptest-framework/pkg/suite"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
)

var (
	ts            *suite.TestSuite
	ctx           = context.Background()
	kubeClient    kubernetes.Interface
	appName       = "cert-manager-app"
	appNamespace  = "default"
)

func TestBasic(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "cert-manager-app Basic Suite")
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

// Test that cert-manager app can be installed successfully
var _ = Describe("cert-manager-app installation", func() {
	It("should install successfully and reach deployed status", func() {
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

	It("should have all required deployments ready", func() {
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
			}, 5*time.Minute, 15*time.Second).Should(BeNumerically(">", 0))
		}
	})

	It("should have all required cluster issuers created", func() {
		// Check for the expected cluster issuers
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
				// Use dynamic client to check for ClusterIssuer (cert-manager CRD)
				dynamicClient := ts.GetClient().Dynamic()
				_, err := dynamicClient.Resource(clusterIssuerGVR).Get(ctx, issuerName, metav1.GetOptions{})
				return err
			}, 5*time.Minute, 15*time.Second).Should(Succeed())
		}
	})

	It("should have cert-manager pods running in correct namespace", func() {
		pods, err := kubeClient.CoreV1().Pods(appNamespace).List(ctx, metav1.ListOptions{
			LabelSelector: "app.kubernetes.io/name=cert-manager",
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(len(pods.Items)).To(BeNumerically(">=", 1))

		for _, pod := range pods.Items {
			Expect(pod.Status.Phase).To(Equal(corev1.PodPhase("Running")))
		}
	})
})

// Test basic certificate functionality
var _ = Describe("cert-manager certificate functionality", func() {
	var testCertName = "test-selfsigned-cert"
	var testSecretName = "test-selfsigned-secret"

	AfterEach(func() {
		// Cleanup test certificate and secret
		certificateGVR := schema.GroupVersionResource{
			Group:    "cert-manager.io",
			Version:  "v1",
			Resource: "certificates",
		}
		_ = ts.GetClient().Dynamic().Resource(certificateGVR).Namespace(appNamespace).Delete(ctx, testCertName, metav1.DeleteOptions{})
		_ = kubeClient.CoreV1().Secrets(appNamespace).Delete(ctx, testSecretName, metav1.DeleteOptions{})
	})

	It("should be able to issue a self-signed certificate", func() {
		// Create a test certificate using the selfsigned cluster issuer
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
		}, 3*time.Minute, 10*time.Second).Should(BeTrue())

		// Verify the secret was created
		Eventually(func() error {
			_, err := kubeClient.CoreV1().Secrets(appNamespace).Get(ctx, testSecretName, metav1.GetOptions{})
			return err
		}, 1*time.Minute, 5*time.Second).Should(Succeed())
	})
})
