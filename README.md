[![CircleCI](https://circleci.com/gh/giantswarm/cert-manager-app.svg?style=shield)](https://circleci.com/gh/giantswarm/cert-manager-app)

# cert-manager-app

Helm chart for the [cert-manager](https://cert-manager.io/) app running in Giant Swarm clusters.

cert-manager adds certificates and certificate issuers (e.g. [Let's Encrypt](https://letsencrypt.org/docs/) (ACME)) as resource types in Kubernetes clusters, and simplifies the process of obtaining, renewing and using those certificates.

## Index
- [Installing](#installing)
- [Configuration](#configuration)
- [Upgrading](#upgrading)
- [For Developers](#for-developers)
  - [Installing the Chart Locally](#installing-the-chart-locally)
  - [Release Process](#release-process)
- [Contributing & Reporting Bugs](#contributing--reporting-bugs)

## Installing

There are 3 ways to install this app onto a workload cluster.

1. [Using our web interface](https://docs.giantswarm.io/ui-api/web/app-platform/#installing-an-app)
2. [Using our API](https://docs.giantswarm.io/api/#operation/createClusterAppV5)
3. Directly creating the [App custom resource](https://docs.giantswarm.io/getting-started/app-platform/deploy-app/#creating-an-app-cr) on the management cluster.

### Sample App CR for the Management Cluster

If you have access to the Kubernetes API on the management cluster, you can create the App CR directly. More information about App CRDs can be found [here](https://docs.giantswarm.io/use-the-api/management-api/crd/apps.application.giantswarm.io/)

Here is an example that would install cert-manager to workload cluster `abc12`:


```bash
# Create appCRD.yaml
$ kubectl gs template app \
>   --catalog giantswarm \
>   --name cert-manager \
>   --version 2.20.2 \
>   --cluster-name abc12 \
>   --target-namespace my-org \
>   > appCR.yaml
```
```yaml
# appCRD.yaml
---
apiVersion: application.giantswarm.io/v1alpha1
kind: App
metadata:
  # workload cluster resources live in a namespace with the same ID as the workload cluster.
  name: cert-manager
  namespace: abc12
spec:
  catalog: giantswarm
  kubeConfig:
    inCluster: false
  name: cert-manager
  namespace: my-org
  version: 2.20.2
```
If you called this file `appCRD.yaml`, you can use the command: `kubectl apply -f appCRD.yaml` to deploy this app to a workload cluster with the ID `abc12`.

### Issuing Certificates

First, an [Issuer](https://cert-manager.io/docs/configuration/) should be configured. There are several ways to [issue certificates](https://cert-manager.io/docs/usage/) through cert-manager. Check upstream documentation for your use case.

## Configuration

Configuration options are documented in [Configuration.md](https://github.com/giantswarm/cert-manager-app/blob/master/helm/cert-manager-app/Configuration.md) document.

## Upgrading

Before upgrading, please check [Upgrading.md](https://github.com/giantswarm/cert-manager-app/blob/master/docs/upgrading.md).

## For Developers

### Installing the Chart Locally

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

### Release Process

* Ensure CHANGELOG.md is up to date.
* Create a new branch to trigger the release workflow as either a patch, minor, or major. E.g. to release a patch, create a branch from master called `release#patch` and push it. Automation will create a release PR.
* Merging the release PR will push a new git tag and trigger a new tarball to be pushed to the
  [giantswarm-catalog].
* Test and verify the cert-manager release across supported environments in a new or existing WIP platform release.

## Contributing & Reporting Bugs
If you have suggestions for how `cert-manager` could be improved, or want to report a bug, open an issue! We'd love all and any contributions.

Check out the [Contributing Guide](https://github.com/giantswarm/cert-manager-app/blob/main/CONTRIBUTING.md) for details on the contribution workflow, submitting patches, and reporting bugs.

---

[app-operator]: https://github.com/giantswarm/app-operator
[cluster-operator]: https://github.com/giantswarm/cluster-operator
[cert-manager]: https://github.com/cert-manager/cert-manager
[default-catalog]: https://github.com/giantswarm/default-catalog
[default-test-catalog]: https://github.com/giantswarm/default-test-catalog
[giantswarm-catalog]: https://github.com/giantswarm/giantswarm-catalog
