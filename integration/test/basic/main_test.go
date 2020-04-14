// +build k8srequired

package basic

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/giantswarm/appcatalog"
	e2esetup "github.com/giantswarm/e2esetup/chart"
	"github.com/giantswarm/e2esetup/chart/env"
	"github.com/giantswarm/e2etests/basicapp"
	"github.com/giantswarm/helmclient"
	"github.com/giantswarm/k8sclient"
	"github.com/giantswarm/micrologger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	app            = "cert-manager"
	appName        = "cert-manager-app"
	catalogURL     = "https://giantswarm.github.io/default-catalog"
	testCatalogURL = "https://giantswarm.github.io/default-test-catalog"
)

const (
	envVarTarballURL = "E2E_TARBALL_URL"
)

var (
	ba         *basicapp.BasicApp
	helmClient *helmclient.Client
	k8sSetup   *k8sclient.Setup
	l          micrologger.Logger
	tarballURL string
)

func init() {
	ctx := context.Background()
	var err error

	var latestRelease string
	{
		latestRelease, err = appcatalog.GetLatestVersion(ctx, catalogURL, appName)
		if err != nil {
			panic(err.Error())
		}
	}

	{
		version := fmt.Sprintf("%s-%s", latestRelease, env.CircleSHA())
		tarballURL, err = appcatalog.NewTarballURL(testCatalogURL, appName, version)
		if err != nil {
			panic(err.Error())
		}
	}

	{
		c := micrologger.Config{}
		l, err = micrologger.New(c)
		if err != nil {
			panic(err.Error())
		}
	}

	var k8sClients *k8sclient.Clients
	{
		c := k8sclient.ClientsConfig{
			Logger: l,

			KubeConfigPath: env.KubeConfigPath(),
		}
		k8sClients, err = k8sclient.NewClients(c)
		if err != nil {
			panic(err.Error())
		}
	}

	{
		c := k8sclient.SetupConfig{
			Logger: l,

			Clients: k8sClients,
		}
		k8sSetup, err = k8sclient.NewSetup(c)
		if err != nil {
			panic(err.Error())
		}
	}

	{
		c := helmclient.Config{
			Logger:          l,
			K8sClient:       k8sClients.K8sClient(),
			RestConfig:      k8sClients.RESTConfig(),
			TillerNamespace: "giantswarm",
		}
		helmClient, err = helmclient.New(c)
		if err != nil {
			panic(err.Error())
		}
	}

	{
		c := basicapp.Config{
			Clients:    k8sClients,
			HelmClient: helmClient,
			Logger:     l,

			App: basicapp.Chart{
				Name:      app,
				Namespace: metav1.NamespaceSystem,
				URL:       tarballURL,
			},
			ChartResources: basicapp.ChartResources{
				Deployments: []basicapp.Deployment{
					{
						Name:      app,
						Namespace: metav1.NamespaceSystem,
						DeploymentLabels: map[string]string{
							"giantswarm.io/service-type": "managed",
							"app":                        app,
						},
						MatchLabels: map[string]string{
							"app": app,
						},
						PodLabels: map[string]string{
							"giantswarm.io/service-type": "managed",
							"app":                        app,
						},
					},
				},
			},
		}
		ba, err = basicapp.New(c)
		if err != nil {
			panic(err.Error())
		}
	}
}

// TestMain allows us to have common setup and teardown steps that are run
// once for all the tests https://golang.org/pkg/testing/#hdr-Main.
func TestMain(m *testing.M) {
	ctx := context.Background()

	{
		c := e2esetup.Config{
			HelmClient: helmClient,
			Setup:      k8sSetup,
		}

		v, err := e2esetup.Setup(ctx, m, c)
		if err != nil {
			l.LogCtx(ctx, "level", "error", "message", "e2e test failed", "stack", fmt.Sprintf("%#v\n", err))
		}

		os.Exit(v)
	}
}
