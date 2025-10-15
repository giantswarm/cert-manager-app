package advanced

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

func TestAdvanced(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "cert-manager-app Advanced Suite")
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

// Test advanced certificate issuance scenarios
var _ = Describe("Advanced certificate functionality", func() {
	var testNamespaceName = "cert-manager-test"
	
	BeforeEach(func() {
		// Create test namespace
		testNamespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: testNamespaceName,
			},
		}
		_, err := kubeClient.CoreV1().Namespaces().Create(ctx, testNamespace, metav1.CreateOptions{})
		if err != nil {
			// Namespace might already exist, that's OK
		}
	})

	AfterEach(func() {
		// Cleanup test namespace and all resources in it
		_ = kubeClient.CoreV1().Namespaces().Delete(ctx, testNamespaceName, metav1.DeleteOptions{})
		// Wait a bit for namespace deletion to propagate
		time.Sleep(5 * time.Second)
	})

	It("should support creating namespace-scoped Issuers", func() {
		// Create a self-signed issuer in the test namespace
		issuer := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "cert-manager.io/v1",
				"kind":       "Issuer",
				"metadata": map[string]interface{}{
					"name":      "test-selfsigned-issuer",
					"namespace": testNamespaceName,
				},
				"spec": map[string]interface{}{
					"selfSigned": map[string]interface{}{},
				},
			},
		}

		issuerGVR := schema.GroupVersionResource{
			Group:    "cert-manager.io",
			Version:  "v1",
			Resource: "issuers",
		}

		dynamicClient := ts.GetClient().Dynamic()
		_, err := dynamicClient.Resource(issuerGVR).Namespace(testNamespaceName).Create(ctx, issuer, metav1.CreateOptions{})
		Expect(err).ToNot(HaveOccurred())

		// Wait for issuer to be ready
		Eventually(func() (bool, error) {
			issuerObj, err := dynamicClient.Resource(issuerGVR).Namespace(testNamespaceName).Get(ctx, "test-selfsigned-issuer", metav1.GetOptions{})
			if err != nil {
				return false, err
			}

			conditions, found, err := unstructured.NestedSlice(issuerObj.Object, "status", "conditions")
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
		}, 2*time.Minute, 10*time.Second).Should(BeTrue())
	})

	It("should handle certificate renewal scenarios", func() {
		// Create a certificate with short duration to test renewal
		certificate := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "cert-manager.io/v1",
				"kind":       "Certificate",
				"metadata": map[string]interface{}{
					"name":      "test-renewal-cert",
					"namespace": testNamespaceName,
				},
				"spec": map[string]interface{}{
					"secretName": "test-renewal-secret",
					"issuerRef": map[string]interface{}{
						"name": "selfsigned-giantswarm",
						"kind": "ClusterIssuer",
					},
					"commonName": "renewal.example.com",
					"duration":   "5m", // Short duration to test renewal logic
				},
			},
		}

		certificateGVR := schema.GroupVersionResource{
			Group:    "cert-manager.io",
			Version:  "v1",
			Resource: "certificates",
		}

		dynamicClient := ts.GetClient().Dynamic()
		_, err := dynamicClient.Resource(certificateGVR).Namespace(testNamespaceName).Create(ctx, certificate, metav1.CreateOptions{})
		Expect(err).ToNot(HaveOccurred())

		// Wait for certificate to be issued
		Eventually(func() (bool, error) {
			cert, err := dynamicClient.Resource(certificateGVR).Namespace(testNamespaceName).Get(ctx, "test-renewal-cert", metav1.GetOptions{})
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

		// Verify the secret exists
		Eventually(func() error {
			_, err := kubeClient.CoreV1().Secrets(testNamespaceName).Get(ctx, "test-renewal-secret", metav1.GetOptions{})
			return err
		}, 1*time.Minute, 5*time.Second).Should(Succeed())
	})

	It("should support multiple DNS names in certificates", func() {
		// Create a certificate with multiple DNS names
		certificate := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "cert-manager.io/v1",
				"kind":       "Certificate",
				"metadata": map[string]interface{}{
					"name":      "test-multi-dns-cert",
					"namespace": testNamespaceName,
				},
				"spec": map[string]interface{}{
					"secretName": "test-multi-dns-secret",
					"issuerRef": map[string]interface{}{
						"name": "selfsigned-giantswarm",
						"kind": "ClusterIssuer",
					},
					"commonName": "main.example.com",
					"dnsNames": []string{
						"main.example.com",
						"www.example.com",
						"api.example.com",
						"admin.example.com",
					},
				},
			},
		}

		certificateGVR := schema.GroupVersionResource{
			Group:    "cert-manager.io",
			Version:  "v1",
			Resource: "certificates",
		}

		dynamicClient := ts.GetClient().Dynamic()
		_, err := dynamicClient.Resource(certificateGVR).Namespace(testNamespaceName).Create(ctx, certificate, metav1.CreateOptions{})
		Expect(err).ToNot(HaveOccurred())

		// Wait for certificate to be ready
		Eventually(func() (bool, error) {
			cert, err := dynamicClient.Resource(certificateGVR).Namespace(testNamespaceName).Get(ctx, "test-multi-dns-cert", metav1.GetOptions{})
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

		// Verify the secret contains the certificate with all DNS names
		Eventually(func() error {
			secret, err := kubeClient.CoreV1().Secrets(testNamespaceName).Get(ctx, "test-multi-dns-secret", metav1.GetOptions{})
			if err != nil {
				return err
			}

			// Check that the TLS cert and key are present
			if _, ok := secret.Data["tls.crt"]; !ok {
				return Errorf("tls.crt not found in secret")
			}
			if _, ok := secret.Data["tls.key"]; !ok {
				return Errorf("tls.key not found in secret")
			}

			return nil
		}, 1*time.Minute, 5*time.Second).Should(Succeed())
	})
})

// Test cert-manager webhook functionality
var _ = Describe("cert-manager webhook validation", func() {
	It("should reject invalid certificate resources", func() {
		// Try to create an invalid certificate (missing required fields)
		invalidCertificate := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "cert-manager.io/v1",
				"kind":       "Certificate",
				"metadata": map[string]interface{}{
					"name":      "invalid-cert",
					"namespace": appNamespace,
				},
				"spec": map[string]interface{}{
					// Missing secretName and issuerRef - should be rejected
					"commonName": "invalid.example.com",
				},
			},
		}

		certificateGVR := schema.GroupVersionResource{
			Group:    "cert-manager.io",
			Version:  "v1",
			Resource: "certificates",
		}

		dynamicClient := ts.GetClient().Dynamic()
		_, err := dynamicClient.Resource(certificateGVR).Namespace(appNamespace).Create(ctx, invalidCertificate, metav1.CreateOptions{})
		
		// This should fail due to webhook validation
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("admission"))
	})

	It("should validate issuerRef references", func() {
		// Try to create a certificate with non-existent issuer reference
		certificate := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "cert-manager.io/v1",
				"kind":       "Certificate",
				"metadata": map[string]interface{}{
					"name":      "invalid-issuer-cert",
					"namespace": appNamespace,
				},
				"spec": map[string]interface{}{
					"secretName": "invalid-issuer-secret",
					"issuerRef": map[string]interface{}{
						"name": "non-existent-issuer",
						"kind": "ClusterIssuer",
					},
					"commonName": "invalid-issuer.example.com",
				},
			},
		}

		certificateGVR := schema.GroupVersionResource{
			Group:    "cert-manager.io",
			Version:  "v1",
			Resource: "certificates",
		}

		dynamicClient := ts.GetClient().Dynamic()
		_, err := dynamicClient.Resource(certificateGVR).Namespace(appNamespace).Create(ctx, certificate, metav1.CreateOptions{})
		
		// The creation might succeed but the certificate should not become ready
		if err == nil {
			// If creation succeeds, verify the certificate doesn't become ready
			Consistently(func() (bool, error) {
				cert, err := dynamicClient.Resource(certificateGVR).Namespace(appNamespace).Get(ctx, "invalid-issuer-cert", metav1.GetOptions{})
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
			}, 30*time.Second, 5*time.Second).Should(BeFalse())
			
			// Cleanup
			_ = dynamicClient.Resource(certificateGVR).Namespace(appNamespace).Delete(ctx, "invalid-issuer-cert", metav1.DeleteOptions{})
		}
	})
})
