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

## Upgrading from v0.9.0 (Giant Swarm app v1.0.8)

If you are using a version of the app prior to `v1.0.8` then please upgrade to `v1.0.8` before carrying out the following steps.

1: First cordon the Chart custom resource. This ensures that `chart-operator` doesn't try and replace the app until the following steps are complete.

```bash
kubectl -n giantswarm annotate chart cert-manager 'chart-operator.giantswarm.io/cordon-reason'='Update in progress'
kubectl -n giantswarm annotate chart cert-manager 'chart-operator.giantswarm.io/cordon-until'='2020-07-20T16:00:00'
```

Where the app is named `cert-manager` and `2020-07-20T16:00:00` is the date and time when reconcilliation of the Chart will be resumed. Ensure you allow yourself enough time to complete the following steps.

2: Back up all Kubernetes secrets of type `kubernetes.io/tls`.

```bash
#!/bin/bash

mkdir kubernetes-secrets-backup

for ns in `kubectl get ns -o custom-columns=NAME:.metadata.name --no-headers=true`; do
  for secret in `kubectl get secrets -n $ns -o custom-columns=NAME:.metadata.name --no-headers=true --field-selector type="kubernetes.io/tls"`; do
    kubectl get secret $secret -n $ns -o json | jq 'del(.metadata.resourceVersion,.metadata.uid) | .metadata.creationTimestamp=null' > kubernetes-secrets-backup/"${ns}_${secret}".json
  done
done
```

This bash snippet requires [jq](https://github.com/stedolan/jq) to be installed.

Note that this will strip resource-specific information from the manifests - if this is not desirable then omit the pipe to `jq`.

3: Uninstall the Helm release.

```bash
helm --tiller-namespace giantswarm delete --purge cert-manager
```

Where `cert-manager` is the name of the release. This requires Helm v2.

4: Upgrade the app to `v2.0.0` (which contains cert-manager `v0.15.2`) via Happa or the API.

5: Uncordon the Chart.

```bash
kubectl -n giantswarm annotate chart cert-manager chart-operator.giantswarm.io/cordon-reason-
kubectl -n giantswarm annotate chart cert-manager chart-operator.giantswarm.io/cordon-until-
```

Where `cert-manager` is the name of the Chart.

The app will be updated when `chart-operator` next reconciles the Chart resource.

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
