package basic

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	cmapi "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmapiv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	"github.com/giantswarm/apptest-framework/pkg/config"
	"github.com/giantswarm/apptest-framework/pkg/state"
	"github.com/giantswarm/apptest-framework/pkg/suite"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	isUpgrade              = false
	namespace              = "kube-system"
	testTimeout            = 10 * time.Minute
	retryInterval          = 5 * time.Second
	deploymentReadyTimeout = 5 * time.Minute
)

func TestBasic(t *testing.T) {
	suite.New(config.MustLoad("../../config.yaml")).
		WithInstallNamespace(namespace).
		WithIsUpgrade(isUpgrade).
		WithValuesFile("./values.yaml").
		Tests(func() {
			It("should have CRDs installed", func() {
				wcClient := state.GetFramework().MC()

				crds := []string{
					"certificates.cert-manager.io",
					"certificaterequests.cert-manager.io",
					"issuers.cert-manager.io",
					"clusterissuers.cert-manager.io",
					"orders.acme.cert-manager.io",
					"challenges.acme.cert-manager.io",
				}

				for _, crdName := range crds {
					By(fmt.Sprintf("Checking CRD %s", crdName))
					Eventually(func() error {
						crd := &apiextensionsv1.CustomResourceDefinition{}
						err := wcClient.Get(context.Background(), types.NamespacedName{Name: crdName}, crd)
						return err
					}).
						WithTimeout(testTimeout).
						WithPolling(retryInterval).
						ShouldNot(HaveOccurred())
				}
			})

			It("should deploy cert-manager components", func() {
				wcClient := state.GetFramework().MC()

				deployments := []string{
					"cert-manager-app",
					"cert-manager-app-webhook",
					"cert-manager-app-cainjector",
				}

				for _, deploymentName := range deployments {
					By(fmt.Sprintf("Checking deployment %s", deploymentName))
					Eventually(func() error {
						deployment := &appsv1.Deployment{}
						err := wcClient.Get(context.Background(), types.NamespacedName{
							Name:      deploymentName,
							Namespace: namespace,
						}, deployment)
						if err != nil {
							return err
						}

						if deployment.Status.AvailableReplicas < 1 {
							return fmt.Errorf("deployment %s not ready", deploymentName)
						}
						return nil
					}).
						WithTimeout(deploymentReadyTimeout).
						WithPolling(retryInterval).
						ShouldNot(HaveOccurred())
				}
			})

			It("should create required ServiceAccounts", func() {
				wcClient := state.GetFramework().MC()

				serviceAccounts := []string{
					"cert-manager-app",
					"cert-manager-app-webhook",
					"cert-manager-app-cainjector",
				}

				for _, saName := range serviceAccounts {
					By(fmt.Sprintf("Checking ServiceAccount %s", saName))
					Eventually(func() error {
						sa := &corev1.ServiceAccount{}
						err := wcClient.Get(context.Background(), types.NamespacedName{
							Name:      saName,
							Namespace: namespace,
						}, sa)
						return err
					}).
						WithTimeout(testTimeout).
						WithPolling(retryInterval).
						ShouldNot(HaveOccurred())
				}
			})

			It("should create ClusterRoles and ClusterRoleBindings", func() {
				wcClient := state.GetFramework().MC()

				clusterRoles := []string{
					"cert-manager-app-controller-issuers",
					"cert-manager-app-controller-clusterissuers",
					"cert-manager-app-controller-certificates",
					"cert-manager-app-controller-orders",
					"cert-manager-app-controller-challenges",
					"cert-manager-app-controller-ingress-shim",
					"cert-manager-app-cainjector",
					"cert-manager-app-webhook:subjectaccessreviews",
				}

				for _, roleName := range clusterRoles {
					By(fmt.Sprintf("Checking ClusterRole %s", roleName))
					Eventually(func() error {
						role := &rbacv1.ClusterRole{}
						err := wcClient.Get(context.Background(), types.NamespacedName{Name: roleName}, role)
						return err
					}).
						WithTimeout(testTimeout).
						WithPolling(retryInterval).
						ShouldNot(HaveOccurred())

					By(fmt.Sprintf("Checking ClusterRoleBinding %s", roleName))
					Eventually(func() error {
						binding := &rbacv1.ClusterRoleBinding{}
						err := wcClient.Get(context.Background(), types.NamespacedName{Name: roleName}, binding)
						return err
					}).
						WithTimeout(testTimeout).
						WithPolling(retryInterval).
						ShouldNot(HaveOccurred())
				}
			})

			It("should create default ClusterIssuers", func() {
				wcClient := state.GetFramework().MC()

				issuers := []string{
					"letsencrypt-giantswarm",
					"selfsigned-giantswarm",
				}

				for _, issuerName := range issuers {
					By(fmt.Sprintf("Checking ClusterIssuer %s", issuerName))
					Eventually(func() error {
						issuer := &cmapi.ClusterIssuer{}
						err := wcClient.Get(context.Background(), types.NamespacedName{Name: issuerName}, issuer)
						if err != nil {
							return err
						}

						// Check if issuer is ready
						for _, condition := range issuer.Status.Conditions {
							if condition.Type == "Ready" {
								if condition.Status == "True" {
									return nil
								}
								return fmt.Errorf("issuer %s not ready: %s", issuerName, condition.Message)
							}
						}
						return errors.New("ready condition not found")
					}).
						WithTimeout(testTimeout).
						WithPolling(retryInterval).
						ShouldNot(HaveOccurred())
				}
			})

			It("should create test certificate with self-signed issuer", func() {
				wcClient := state.GetFramework().MC()

				By("Creating test certificate")
				cert := &cmapiv1.Certificate{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-cert",
						Namespace: namespace,
					},
					Spec: cmapiv1.CertificateSpec{
						SecretName: "test-cert-tls",
						IssuerRef: cmmeta.ObjectReference{
							Name:  "selfsigned-giantswarm",
							Kind:  "ClusterIssuer",
							Group: "cert-manager.io",
						},
						CommonName: "test.local",
						DNSNames:   []string{"test.local"},
					},
				}

				err := wcClient.Create(context.Background(), cert)
				Expect(err).ShouldNot(HaveOccurred())

				By("Waiting for certificate to be ready")
				Eventually(func() error {
					cert := &cmapi.Certificate{}
					err := wcClient.Get(context.Background(), types.NamespacedName{
						Name:      "test-cert",
						Namespace: namespace,
					}, cert)
					if err != nil {
						return err
					}

					for _, condition := range cert.Status.Conditions {
						if condition.Type == "Ready" {
							if condition.Status == "True" {
								return nil
							}
							return fmt.Errorf("certificate not ready: %s", condition.Message)
						}
					}
					return errors.New("ready condition not found")
				}).
					WithTimeout(testTimeout).
					WithPolling(retryInterval).
					ShouldNot(HaveOccurred())

				By("Verifying the certificate secret was created")
				Eventually(func() error {
					secret := &corev1.Secret{}
					err := wcClient.Get(context.Background(), types.NamespacedName{
						Name:      "test-cert-tls",
						Namespace: namespace,
					}, secret)
					return err
				}).
					WithTimeout(testTimeout).
					WithPolling(retryInterval).
					ShouldNot(HaveOccurred())
			})
		}).
		AfterSuite(func() {
			wcClient := state.GetFramework().MC()

			By("Cleaning up test certificate")
			cert := &cmapi.Certificate{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-cert",
					Namespace: namespace,
				},
			}
			err := wcClient.Delete(context.Background(), cert)
			Expect(err).ShouldNot(HaveOccurred())

			By("Cleaning up test certificate secret")
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-cert-tls",
					Namespace: namespace,
				},
			}
			err = wcClient.Delete(context.Background(), secret)
			Expect(err).ShouldNot(HaveOccurred())
		}).
		Run(t, "Basic Test")
}
