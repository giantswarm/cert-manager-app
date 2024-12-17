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
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	isUpgrade        = false
	appReadyTimeout  = 5 * time.Minute
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
			It("should deploy cert-manager App CR successfully", func() {
				Expect(state.GetCluster()).NotTo(BeNil(), "cluster state should be initialized")
				Expect(state.GetCluster().Organization).NotTo(BeNil(), "organization should be available")

				namespace := state.GetCluster().Organization.GetNamespace()

				By("Verifying cert-manager App CR is deployed")
				for _, component := range components {
					appName := fmt.Sprintf("%s-%s", state.GetCluster().Name, component)
					Eventually(wait.IsAppDeployed(context.Background(),
						state.GetFramework().MC(),
						appName,
						namespace)).
						WithTimeout(appReadyTimeout).
						WithPolling(appReadyInterval).
						Should(BeTrue(), fmt.Sprintf("%s should be deployed", component))
				}
			})

			It("should have all components running and ready", func() {
				wcClient, err := state.GetFramework().WC(state.GetCluster().Name)
				Expect(err).NotTo(HaveOccurred(), "should get workload cluster client")

				componentConfigs := map[string]struct {
					namespace string
					kind      string
					name      string
				}{
					// cert-manager namespace
					"cert-manager-controller": {
						namespace: "kube-system",
						kind:      "Deployment",
						name:      "cert-manager",
					},
					"cert-manager-cainjector": {
						namespace: "kube-system",
						kind:      "Deployment",
						name:      "cert-manager-cainjector",
					},
					"cert-manager-webhook": {
						namespace: "kube-system",
						kind:      "Deployment",
						name:      "cert-manager-webhook",
					},

					// webhooks
					"cert-manager-webhook-configuration": {
						namespace: "",
						kind:      "ValidatingWebhookConfiguration",
						name:      "cert-manager-webhook",
					},
					"cert-manager-webhook-mutating": {
						namespace: "",
						kind:      "MutatingWebhookConfiguration",
						name:      "cert-manager-webhook",
					},
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

						case "MutatingWebhookConfiguration":
							mutatingWebhook := &admissionregistrationv1.MutatingWebhookConfiguration{}
							err := wcClient.Get(context.Background(), client.ObjectKey{Name: config.name}, mutatingWebhook)
							return err == nil

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
		}).
		Run(t, "Basic Test")
}
