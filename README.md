[![CircleCI](https://circleci.com/gh/giantswarm/cert-manager-app.svg?style=shield)](https://circleci.com/gh/giantswarm/cert-manager-app)

# cert-manager-app

Helm chart for the [cert-manager](https://cert-manager.io/) app running in Giant Swarm clusters.

cert-manager adds certificates and certificate issuers (e.g. [Let's Encrypt](https://letsencrypt.org/docs/) (ACME)) as resource types in Kubernetes clusters, and simplifies the process of obtaining, renewing and using those certificates.

## Index
- [Installing the Chart](#installing-the-chart)
- [Configuration](#configuration)
- [Upgrading](#upgrading)
- [Release Process](#release-process)
- [Contributing & Reporting Bugs](#contributing--reporting-bugs)

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

## Release Process

* Ensure CHANGELOG.md is up to date.
* Create a new branch to trigger the release workflow as either a patch, minor, or major. E.g. to release a patch, create a branch from master called master#release#patch and push it. Automation will create a release PR.
* Merging the release PR will push a new git tag and trigger a new tarball to be pushed to the
  [giantswarm-catalog].
* Test and verify the cert-manager release across supported environments in a new or existing WIP platform release.

## Contributing & Reporting Bugs
If you have suggestions for how `cert-manager` could be improved, or want to report a bug, open an issue! We'd love all and any contributions.

Check out the [Contributing Guide](https://github.com/giantswarm/cert-manager-app/blob/main/CONTRIBUTING.md) for details on the contribution workflow, submitting patches, and reporting bugs.


---

[app-operator]: https://github.com/giantswarm/app-operator
[cluster-operator]: https://github.com/giantswarm/cluster-operator
[cert-manager]: https://github.com/jetstack/cert-manager
[default-catalog]: https://github.com/giantswarm/default-catalog
[default-test-catalog]: https://github.com/giantswarm/default-test-catalog
[giantswarm-catalog]: https://github.com/giantswarm/giantswarm-catalog
