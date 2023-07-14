# Upgrading

## Upgrading from >= v2.24.0 to v3.0.0

In Chart version 3.0.0 and above the `values.yaml` schema and deplyment names have been changed.

If you are using custom configuration with this chart, some intervention may be required when upgrading. If you're using the default values, the upgrade should work without intervention.

### Values

- Additional arguments for the controller container moved and must be configured by setting `extraArgs` instead of `controller.extraArgs`
- Controller image pull policy settings must be configured by setting `image.pullPolicy` instead of `controller.image.pullPolicy`.
- Support of changing the `logLevel` of `cainjector`, `controller` and `webhook` individually has been removed. Use `global.logLevel` instead.
- Proxy configuration moved to top-level `http_proxy`, `https_proxy` and `no_proxy` values instead of `proxy.http`, `proxy.https` and `proxy.noProxy` or `cluster.proxy.http`, `cluster.proxy.https` and `cluster.proxy.noProxy`.
- Certificate owner ref can be enabled using `global.enableCertificateOwnerRef` instead of `global.enableCertOwnerRef`.
- The number of replicas can now be changed by setting `replicaCount`, `cainjector.replicaCount` and `webhook.replicaCount` instead of `controller.replicas`, `cainjector.replicas` and `webhook.replicas`
- To configure AWS IAM role pod annotations, set `iam.amazonaws.com/role: your-role` in `podAnnotations` instead of `controller.aws.role`.
- Provider specific tempaltes and config values have been removed. To configure IRSA for AWS set `eks.amazonaws.com/role-arn: <role-reference>` in `serviceAccount.annotations` instead of `controller.aws.irsa`.
- Provider specific arguments `provider` and `clusterID` have been removed. These fields were part of AWS IAM and IRSA configuration.
- ACMESolver settings previously configured through values `global.acmeSolver` moved into subchart config field `giantSwarmClusterIssuer`. See the subcharts [values.yaml](https://github.com/giantswarm/cert-manager-app/blob/main/helm/cert-manager-app/charts/cert-manager-giantswarm-clusterissuer/values.yaml) for details.
- For configuring DNS01 recursive nameservers, set `dns01RecursiveNameserversOnly` to `true` and set your nameservers using `dns01RecursiveNameservers` (host and port). Previously, this was done by setting `global.acmeSolver.DNSServer`.
- The container security context is now configured separately using keys `containerSecurityContext`, `webhook.containerSecurityContext` and `cainjector.containerSecurityContext` instead of `global.containerSecurityContext`.
- The pod security context is now configured separately using keys `securityContext`, `webhook.securityContext` and `cainjector.securityContext` instead of `global.securityContext`.
- The default issuer is now configured by setting `ingressShim.defaultIssuerName`, `ingressShim.defaultIssuerKind` and `ingressShim.defaultIssuerGroup` instead of `controller.defaultIssuer.name`, `controller.defaultIssuer.kind` and `controller.defaultIssuer.group`

## Upgrading from v0.9.0 (Giant Swarm app v1.0.8 to 2.x.x)

If you are using a version of the App prior to `v1.0.8`, please upgrade to `v1.0.8` first.

From `v1.0.8`, the upgrade path is as follows:

`v1.0.8 (cert-manager 0.9.0) > v2.0.2 (cert-manager 0.15.2) > v2.1.0 (cert-manager 0.16.1)`

### v2.0.2 > v2.1.0

No manual intervention is required, and the App will be upgraded in place.

### v1.0.8 > v2.0.2

The procedure below must be followed when upgrading from `v1.0.8` to `v2.0.2`; this is due to breaking changes introduced in the `cert-manager` API.

To assist with the upgrade, a [migration script](../files/migrate-v090-to-v200.sh) is available in the `files/` directory of this repository. **Read the help text thoroughly before using it**.

**Note:** The upgrade process involves **removing the existing App**. This will also remove the Custom Resource Definitions it provides, which will in turn remove any related Custom Resources.
This will mean all Custom Resources of the following types **will be removed**:

- Issuer
- ClusterIssuer
- Certificate
- CertificateRequest

The [migration script](../files/migrate-v090-to-v200.sh) can be used to backup the Custom Resources.

**1: First cordon the Charts Custom Resource.** This ensures that `chart-operator` doesn't try to replace the App until the following steps are completed.

```bash
kubectl -n giantswarm annotate chart cert-manager 'chart-operator.giantswarm.io/cordon-reason'='Update in progress'
kubectl -n giantswarm annotate chart cert-manager 'chart-operator.giantswarm.io/cordon-until'='2020-07-20T16:00:00'
```

Above `cert-manager` is the name of the App (this can differ, depending on your configuration) and `2020-07-20T16:00:00` is the date and time when reconcilliation of the Chart will be resumed. Ensure you allow yourself enough time to complete the following steps.

As an additional safety step, also scale down `chart-operator`:

```bash
kubectl -n giantswarm scale deploy/chart-operator --replicas=0
```

**2: Back up the following resources.**

#### all namespaces:

- Secret (of type 'kubernetes.io/tls', with deprecated labels/annotations)
- Ingress (where '.spec.tls' is set)
- Issuer
- Certificate
- CertificateRequest

#### cluster-scoped

- ClusterIssuer

Note: the provided [migration script](../files/migrate-v090-to-v200.sh) can be used for this.

**3: Uninstall the Helm release.**

```bash
helm --tiller-namespace giantswarm delete --purge cert-manager
```

Where `cert-manager` is the name of the Helm release. This requires Helm v2 cli to be installed.

**4: Upgrade the App.**

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

**5: Allow the Chart to be reconciled again.**

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

**6: Update annotations and labels on Ingresses and Secrets** (of type `kubernetes.io/tls`) to reflect the new API group.

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

Note: the provided [migration script](../files/migrate-v090-to-v200.sh) can be used for this.

**7: Remove deprecated annotations and labels from Ingresses and Secrets** which were updated previously.

Note: the provided [migration script](../files/migrate-v090-to-v200.sh) can be used for this.
