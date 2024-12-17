package basic

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/giantswarm/apptest-framework/pkg/config"
	"github.com/giantswarm/apptest-framework/pkg/state"
	"github.com/giantswarm/apptest-framework/pkg/suite"
	"github.com/giantswarm/clustertest/pkg/wait"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	isUpgrade        = true
	appReadyTimeout  = 10 * time.Minute
	appReadyInterval = 5 * time.Second
)

var components = []string{
	"cert-manager",
}

func TestBasic(t *testing.T) {
	suite.New(config.MustLoad("../../config.yaml")).
		WithIsUpgrade(isUpgrade).
		WithValuesFile("./values.yaml").
		Tests(func() {
			It("should have cert-manager App CR upgraded successfully", func() {
				Expect(state.GetCluster()).NotTo(BeNil(), "cluster state should be initialized")
				Expect(state.GetCluster().Organization).NotTo(BeNil(), "organization should be available")

				namespace := state.GetCluster().Organization.GetNamespace()

				By("Verifying cert-manager App CR is upgraded")
				for _, component := range components {
					appName := fmt.Sprintf("%s-%s", state.GetCluster().Name, component)
					Eventually(wait.IsAppDeployed(context.Background(),
						state.GetFramework().MC(),
						appName,
						namespace)).
						WithTimeout(appReadyTimeout).
						WithPolling(appReadyInterval).
						Should(BeTrue(), fmt.Sprintf("%s should be upgraded and deployed", component))
				}
			})

			It("should verify all upgraded components are running and ready", func() {
				wcClient, err := state.GetFramework().WC(state.GetCluster().Name)
				Expect(err).NotTo(HaveOccurred(), "should get workload cluster client")

				componentConfigs := map[string]struct {
					namespace string
					kind      string
					name      string
				}{
					// Core components in cert-manager namespace
					"cert-manager":            {namespace: "cert-manager", kind: "Deployment", name: "cert-manager"},
					"cert-manager-webhook":    {namespace: "cert-manager", kind: "Deployment", name: "cert-manager-webhook"},
					"cert-manager-cainjector": {namespace: "cert-manager", kind: "Deployment", name: "cert-manager-cainjector"},

					// Webhook configurations
					"cert-manager-webhook-config": {namespace: "", kind: "ValidatingWebhookConfiguration", name: "cert-manager-webhook"},
				}

				for component, config := range componentConfigs {
					By(fmt.Sprintf("Checking %s %s", component, config.kind))
					Eventually(func() bool {
						var ready, replicas int32
						switch config.kind {
						case "Deployment":
							deployment := &appsv1.Deployment{}
							err := wcClient.Get(context.Background(), client.ObjectKey{Namespace: config.namespace, Name: config.name}, deployment)
							if err != nil {
								return false
							}
							ready = deployment.Status.ReadyReplicas
							replicas = deployment.Status.Replicas
						case "ValidatingWebhookConfiguration":
							validatingWebhook := &admissionregistrationv1.ValidatingWebhookConfiguration{}
							err := wcClient.Get(context.Background(), client.ObjectKey{Name: config.name}, validatingWebhook)
							return err == nil
						}
						return ready == replicas && replicas > 0
					}).
						WithTimeout(appReadyTimeout).
						WithPolling(appReadyInterval).
						Should(BeTrue(), fmt.Sprintf("%s %s should be ready", component, config.kind))
				}
			})

			It("should verify CRDs are available and established after upgrade", func() {
				wcClient, err := state.GetFramework().WC(state.GetCluster().Name)
				Expect(err).NotTo(HaveOccurred(), "should get workload cluster client")

				crds := []string{
					"certificates.cert-manager.io",
					"certificaterequests.cert-manager.io",
					"issuers.cert-manager.io",
					"clusterissuers.cert-manager.io",
					"orders.acme.cert-manager.io",
					"challenges.acme.cert-manager.io",
				}

				for _, crdName := range crds {
					By(fmt.Sprintf("Verifying CRD %s after upgrade", crdName))
					Eventually(func() bool {
						crd := &apiextensionsv1.CustomResourceDefinition{}
						err := wcClient.Get(context.Background(), types.NamespacedName{Name: crdName}, crd)
						if err != nil {
							return false
						}

						for _, condition := range crd.Status.Conditions {
							if condition.Type == "Established" && condition.Status == "True" {
								return true
							}
						}
						return false
					}).
						WithTimeout(appReadyTimeout).
						WithPolling(appReadyInterval).
						Should(BeTrue(), fmt.Sprintf("CRD %s should be established after upgrade", crdName))
				}
			})

			It("should verify existing certificates remain valid after upgrade", func() {
				wcClient, err := state.GetFramework().WC(state.GetCluster().Name)
				Expect(err).NotTo(HaveOccurred(), "should get workload cluster client")

				By("Verifying existing ClusterIssuers are still valid")
				issuers := []string{
					"letsencrypt-giantswarm",
					"selfsigned-giantswarm",
				}

				for _, issuerName := range issuers {
					Eventually(func() error {
						issuer := &apiextensionsv1.CustomResourceDefinition{}
						return wcClient.Get(context.Background(), client.ObjectKey{Name: issuerName}, issuer)
					}).
						WithTimeout(appReadyTimeout).
						WithPolling(appReadyInterval).
						Should(Succeed(), fmt.Sprintf("ClusterIssuer %s should remain valid after upgrade", issuerName))
				}
			})
		}).
		Run(t, "Basic Test")
}
