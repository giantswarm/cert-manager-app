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

If you are using a version of the App prior to `v1.0.8` then please upgrade to `v1.0.8` first.

From `v1.0.8`, the upgrade path is as follows:

`v1.0.8 (cert-manager 0.9.0) > v2.0.2 (cert-manager 0.15.2) > v2.1.0 (cert-manager 0.16.1)`

### v2.0.2 > v2.1.0

No manual intervention is required, and the App will be upgraded in place.

### v1.0.8 > v2.0.2

The procedure below must be followed when upgrading from `v1.0.8` to `v2.0.2`,; this is due to breaking changes introduced in `cert-manager`'s API.

To assist with the upgrade, a [migration script](files/migrate-v090-to-v200.sh) is provided in the `files/` directory of this repository. If you use it, please read the help text thoroughly.

**Note:** The upgrade process involves **removing the existing App**. This will also remove the Custom Resource Definitions it provides, which will in turn remove any related Custom Resources.
This will mean all Custom Resources of the following types **will be removed**:

- Issuer
- ClusterIssuer
- Certificate
- CertificateRequest

The [migration script](files/migrate-v090-to-v200.sh) can be used to ensure that these are backed up.

1: First cordon the Chart custom resource. This ensures that `chart-operator` doesn't try and replace the App until the following steps are complete.

```bash
kubectl -n giantswarm annotate chart cert-manager 'chart-operator.giantswarm.io/cordon-reason'='Update in progress'
kubectl -n giantswarm annotate chart cert-manager 'chart-operator.giantswarm.io/cordon-until'='2020-07-20T16:00:00'
```

Where the App is named `cert-manager` and `2020-07-20T16:00:00` is the date and time when reconcilliation of the Chart will be resumed. Ensure you allow yourself enough time to complete the following steps.

As an additional safety step, also scale down `chart-operator`:

```bash
kubectl -n giantswarm scale deploy/chart-operator --replicas=0
```

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

4: Upgrade the App.

The upgrade process for the App depends on whether the App is optional or pre-installed. To determine this, use `gsctl` to inspect the cluster's release.
Assuming the cluster version is `v11.5.0`, run the following command:

```bash
gsctl show release 11.5.0 | grep cert-manager
 cert-manager: 0.9.0
```

If the command returns a version of `cert-manager`, then it is pre-installed. If nothing is returned, the App is optional.

#### optional App

If the App is optional, it can now be upgraded to `v2.0.2` via Happa or the API. `v2.0.2` is preferred as it fixes some minor bugs present in `2.0.0` and `2.0.1`.

#### pre-installed

If the App is part of a Giant Swarm Release, the cluster should now be upgraded via Happa or the API. Currently, this only applies to AWS clusters using Release `v10.0.0` or later.

5: Allow the Chart to be reconciled again.

Where `cert-manager` is the name of the Chart:

```bash
kubectl -n giantswarm annotate chart cert-manager chart-operator.giantswarm.io/cordon-reason-
kubectl -n giantswarm annotate chart cert-manager chart-operator.giantswarm.io/cordon-until-
```

And also scale `chart-operator` back up again:

```bash
kubectl -n giantswarm scale deploy/chart-operator --replicas=1
```

The App will be updated when `chart-operator` next reconciles the Chart resource.

6: Update annotations and labels on Ingresses and Secrets (of type `kubernetes.io/tls`) to reflect the new API group.

**IMPORTANT:** All references to the API group `certmanager.k8s.io` must be changed to `cert-manager.io`. These are used by `cert-manager` to indicate which resources it should interact with, and if they are left unchanged, `cert-manager` will no longer reconcile them after the App has been upgraded.

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

[app-operator]: https://github.com/giantswarm/app-operator
[cluster-operator]: https://github.com/giantswarm/cluster-operator
[default-catalog]: https://github.com/giantswarm/default-catalog
[default-test-catalog]: https://github.com/giantswarm/default-test-catalog
[cert-manager]: https:cert-manager//github.com/kubernetes-incubator/
