[![CircleCI](https://circleci.com/gh/giantswarm/cert-manager-app.svg?style=shield)](https://circleci.com/gh/giantswarm/cert-manager-app)

# cert-manager-app
Helm chart for cert-manager service running in tenant clusters

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

Deployment to Tenant Clusters is handled by [app-operator](https://github.com/giantswarm/app-operator).

## Configuration

Configuration options are documented in [Configuration.md](helm/cert-manager-app/Configuration.md) document.

## Upgrading from v0.9.0 (Giant Swarm app v1.0.8 to 2.0.x)

If you are using a version of the app prior to `v1.0.8` then please upgrade to `v1.0.8` before carrying out the following steps.

A [migration script](files/migrate-v090-to-v200.sh) is provided in the `files/` directory of this repository. If you use it, please read the help text thoroughly.

1: First cordon the Chart custom resource. This ensures that `chart-operator` doesn't try and replace the app until the following steps are complete.

```bash
kubectl -n giantswarm annotate chart cert-manager 'chart-operator.giantswarm.io/cordon-reason'='Update in progress'
kubectl -n giantswarm annotate chart cert-manager 'chart-operator.giantswarm.io/cordon-until'='2020-07-20T16:00:00'
```

Where the app is named `cert-manager` and `2020-07-20T16:00:00` is the date and time when reconcilliation of the Chart will be resumed. Ensure you allow yourself enough time to complete the following steps.

2: Back up the following resources.

#### all namespaces:

- Secret (of type 'kubernetes.io/tls', with deprecated labels/annotations)
- Ingress (where '.spec.tls' is set)
- Issuer
- Certificate
- CertificateRequest

#### cluster-scoped

- ClusterIssuer

Note: the provided [migration script](files/migrate-v090-to-v200.sh) can be used for this.

3: Uninstall the Helm release.

```bash
helm --tiller-namespace giantswarm delete --purge cert-manager
```

Where `cert-manager` is the name of the release. This requires Helm v2.

4: Upgrade the app to `v2.0.2` (the latest version, which fixes some minor bugs present in `2.0.0` and `2.0.1`) via Happa or the API.

5: Uncordon the Chart.

```bash
kubectl -n giantswarm annotate chart cert-manager chart-operator.giantswarm.io/cordon-reason-
kubectl -n giantswarm annotate chart cert-manager chart-operator.giantswarm.io/cordon-until-
```

Where `cert-manager` is the name of the Chart.

The app will be updated when `chart-operator` next reconciles the Chart resource.

6: Update annotations and labels on Ingresses and Secrets (of type `kubernetes.io/tls`) to reflect the new API group.

**IMPORTANT:** All references to the API group `certmanager.k8s.io` must be changed to `cert-manager.io`. If left unchanged, `cert-manager` will no longer reconcile them.

An example secret **before** being updated:

```yaml
kind: Secret
metadata:
  annotations:
    certmanager.k8s.io/alt-names: helloworld.sag8c.k8s.gauss.eu-central-1.aws.gigantic.io
    certmanager.k8s.io/common-name: helloworld.sag8c.k8s.gauss.eu-central-1.aws.gigantic.io
    certmanager.k8s.io/ip-sans: ""
    certmanager.k8s.io/issuer-kind: ClusterIssuer
    certmanager.k8s.io/issuer-name: letsencrypt-giantswarm
  creationTimestamp: "2020-07-16T11:02:54Z"
  labels:
    certmanager.k8s.io/certificate-name: helloworld-tls
```

The same secret updated to match the new API group:

```yaml
kind: Secret
metadata:
  annotations:
    cert-manager.io/alt-names: helloworld.sag8c.k8s.gauss.eu-central-1.aws.gigantic.io
    cert-manager.io/common-name: helloworld.sag8c.k8s.gauss.eu-central-1.aws.gigantic.io
    cert-manager.io/ip-sans: ""
    cert-manager.io/issuer-kind: ClusterIssuer
    cert-manager.io/issuer-name: letsencrypt-giantswarm
  creationTimestamp: "2020-07-16T11:02:54Z"
  labels:
    cert-manager.io/certificate-name: helloworld-tls
```

Note: the provided [migration script](files/migrate-v090-to-v200.sh) can be used for this.

7: Remove deprecated annotations and labels from Ingresses and Secrets which were updated previously.

Note: the provided [migration script](files/migrate-v090-to-v200.sh) can be used for this.

## Release Process

* Ensure CHANGELOG.md is up to date.
* Create a new GitHub release with the version e.g. `v0.1.0` and link the
changelog entry.
* This will push a new git tag and trigger a new tarball to be pushed to the
[default-catalog].
* Update [cluster-operator] with the new version.

[app-operator]: https://github.com/giantswarm/app-operator
[cluster-operator]: https://github.com/giantswarm/cluster-operator
[default-catalog]: https://github.com/giantswarm/default-catalog
[default-test-catalog]: https://github.com/giantswarm/default-test-catalog
[cert-manager]: https:cert-manager//github.com/kubernetes-incubator/
