[![CircleCI](https://circleci.com/gh/giantswarm/cert-manager-app.svg?style=shield)](https://circleci.com/gh/giantswarm/cert-manager-app)

# cert-manager-app

Helm chart for the [cert-manager](https://cert-manager.io/) app running in Giant Swarm clusters.

## Installing the Chart

To install the chart locally:

```bash
$ git clone https://github.com/giantswarm/cert-manager-app.git
$ cd cert-manager-app
$ helm install helm/cert-manager-app
```

Provide a custom `values.yaml`:

```bash
$ helm install cert-manager-app -f values.yaml
```

Deployment to Workload Clusters is handled by [app-operator](https://github.com/giantswarm/app-operator).

## Configuration

Configuration options are documented in [Configuration.md](https://github.com/giantswarm/cert-manager-app/blob/master/helm/cert-manager-app/Configuration.md) document.

## Upgrading

Before upgrading, please check [upgrading.md](https://github.com/giantswarm/cert-manager-app/blob/master/docs/upgrading.md).

---

[app-operator]: https://github.com/giantswarm/app-operator
[cluster-operator]: https://github.com/giantswarm/cluster-operator
[default-catalog]: https://github.com/giantswarm/default-catalog
[default-test-catalog]: https://github.com/giantswarm/default-test-catalog
[cert-manager]: https://github.com/jetstack/cert-manager
