package k8sclient

import (
	"github.com/giantswarm/apiextensions/pkg/clientset/versioned"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ClientsConfig struct {
	Logger micrologger.Logger

	// KubeConfigPath and RestConfig are mutually exclusive.
	KubeConfigPath string
	// RestConfig and KubeConfigPath are mutually exclusive.
	RestConfig *rest.Config
}

type Clients struct {
	logger micrologger.Logger

	dynClient  dynamic.Interface
	extClient  *apiextensionsclient.Clientset
	g8sClient  *versioned.Clientset
	k8sClient  *kubernetes.Clientset
	restClient rest.Interface
	restConfig *rest.Config
}

func NewClients(config ClientsConfig) (*Clients, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	if config.KubeConfigPath == "" && config.RestConfig == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.KubeConfigPath or %T.RestConfig must not be empty", config, config)
	}
	if config.KubeConfigPath != "" && config.RestConfig != nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.KubeConfigPath and %T.RestConfig must not be set at the same time", config, config)
	}

	var err error

	var restConfig *rest.Config
	{
		if config.RestConfig != nil {
			restConfig = config.RestConfig
		} else {
			restConfig, err = clientcmd.BuildConfigFromFlags("", config.KubeConfigPath)
			if err != nil {
				return nil, microerror.Mask(err)
			}
		}
	}

	var dynClient dynamic.Interface
	{
		c := rest.CopyConfig(restConfig)

		dynClient, err = dynamic.NewForConfig(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var extClient *apiextensionsclient.Clientset
	{
		c := rest.CopyConfig(restConfig)

		extClient, err = apiextensionsclient.NewForConfig(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var g8sClient *versioned.Clientset
	{
		c := rest.CopyConfig(restConfig)

		g8sClient, err = versioned.NewForConfig(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var k8sClient *kubernetes.Clientset
	{
		c := rest.CopyConfig(restConfig)

		k8sClient, err = kubernetes.NewForConfig(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var restClient rest.Interface
	{
		// It would be cool to use rest.RESTClientFor here but it fails
		// because GroupVersion is not configured. So underlying core
		// RESTClient is taken.
		//
		//	panic: GroupVersion is required when initializing a RESTClient
		//
		restClient = k8sClient.RESTClient()
	}

	c := &Clients{
		logger: config.Logger,

		dynClient:  dynClient,
		extClient:  extClient,
		g8sClient:  g8sClient,
		k8sClient:  k8sClient,
		restClient: restClient,
		restConfig: restConfig,
	}

	return c, nil
}

func (c *Clients) DynClient() dynamic.Interface {
	return c.dynClient
}

func (c *Clients) ExtClient() apiextensionsclient.Interface {
	return c.extClient
}

func (c *Clients) G8sClient() versioned.Interface {
	return c.g8sClient
}

func (c *Clients) K8sClient() kubernetes.Interface {
	return c.k8sClient
}

func (c *Clients) RESTClient() rest.Interface {
	return c.restClient
}

func (c *Clients) RESTConfig() *rest.Config {
	return rest.CopyConfig(c.restConfig)
}
