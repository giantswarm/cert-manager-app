# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project's packages adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Add `io.giantswarm.application.audience` and `io.giantswarm.application.managed` chart annotations for Backstage visibility.
- Add PodLogs for log collection.

### Fixed

- Fix `controller` Vertical Pod Autoscaler (VPA) resource syntax.

## [3.9.4] - 2025-10-17

### Changed

- Upgrade cert-manager to v1.18.2.

### Added

- Add E2E tests using apptest-framework for automated PR testing across multiple providers (CAPA, CAPV, CAPZ, CAPVCD).
  - **Basic test suite**: Validates fresh installations
  - **Upgrade test suite**: Tests upgrade scenarios and certificate reconciliation
- Add certificate issuance integration test to cluster-test-suites.

## [3.9.3] - 2025-10-07

### Changed

- Fix missing targetPort in `cainjector-service`

## [3.9.2] - 2025-09-02

### Changed

- Add `alloy` ingress rules for cainjector metrics ingestion.

## [3.9.1] - 2025-04-16

### Added

- Added Vertical Pod Autoscaler support for `controller` pods.
- Added renovate configutarion

### Removed

- Removed dependabot configuration

## [3.9.0] - 2025-01-28

### Changed

- Updates Cert-manager Chart to Upstream 1.16.2

### Added

- Adds new sync method based on Vendir to sync from upstream

## [3.8.2] - 2024-12-03

### Fix

- added the option to configure additional approveSignerNames

### Changed

- Changed ownership to team Shield

### Removed

- Get rid of label `giantswarm.io/monitoring_basic_sli` as this slo generation label is not used anymore.

## [3.8.1] - 2024-07-30

### Changed

- Bump architect-orb@5.3.1 to fix CVE-2024-24790.

## [3.8.0] - 2024-07-10

### Fix

- Improves `cainjector`'s Vertical Pod Autoscaler

## [3.7.9] - 2024-07-08

### Fix

- Remove quotes from acme-http01-solver-image argument. The quotes are used when looking up the image which causes an error.

## [3.7.8] - 2024-07-03

### Added

- Improves container security by setting `runAsGroup` and `runAsUser` greater than zero for all deployments.

## [3.7.7] - 2024-06-18

### Changed

- Changed the way registry is being parsed in helm templates
- Enable VPA by default

## [3.7.6] - 2024-06-10

### Added

- Added Vertical Pod Autoscaler support for `cainjector` pods.

## [3.7.5] - 2024-05-16

### Added

- Added annotation `helm.sh/resource-policy: keep` on CRDs to prevent them from being pruned in an unexpected rollback event.

## [3.7.4] - 2024-03-19

### Added

- Added support for `AzureDNS` integration with a `Service Principal` on `clusterIssuer` helm chart .

### Changed

- Changed `appVersion` to `v1.14.2`

## [3.7.3] - 2024-02-12

### Added

- Added readiness check for `cert-manager-app-webhook` before attempting installation on `clusterIssuers` chart

## [3.7.2] - 2024-02-05

### Added

- Added support for `azureDNS` dns01 challenge solver on `cluster-issuer` chart

## [3.7.1] - 2024-01-29

### Added

- Added `acme-solvers-networkpolicy` `NetworkPolicy` namespace to `kube-system`

## [3.7.0] - 2024-01-22

### Changed

- Update to latest `helm chart` version to `v1.13.3`

## [3.6.1] - 2023-12-06

### Changed

- Configure `gsoci.azurecr.io` as the default container image registry.

## [3.6.0] - 2023-11-22

### Added

- Allow skipping Giant Swarm specific NetworkPolicy resources with `giantswarmNetworkPolicy.enabled` value.

## [3.5.3] - 2023-11-16
### Added
- adds extra `helm chart` for the `ciliumNetworkPolicies`

### Changed
- changes the previous `netpols` `helm chart` to be used only for `networkPolicies`
- disables the `startup-api-check` job that waits for the webhookendpoints to become available

## [3.5.2] - 2023-11-09
### Changed

- Moved `acme-solvers-networkpolicy` to the NetworkPolicies Helm chart for better organization and management of network policies.

### Removed

- Removed `acme-solvers-ciliumnetworkpolicy`

## [3.5.1] - 2023-11-02

### Added
- Introduced `acme-solvers-networkpolicy` and `acme-solvers-ciliumnetworkpolicy` for enhanced network security and control.

## [3.5.0] - 2023-10-12

### Added
- 

- cert-manager-giantswarm-clusterissuer: Allow setting `hostedZoneID` for `route53` DNS01 challenge.
- cert-manager-giantswarm-clusterissuer: Make `accessKeyID` and `secretAccessKey` optional for `route53` DNS01 challenge.

## [3.4.1] - 2023-10-10

### Changed

- Move cert-manager ownership to team BigMac. ([#349](https://github.com/giantswarm/cert-manager-app/pull/349))
- Add default cpu and memory limits to controller, cainjector and webhook deployments. ([#367](https://github.com/giantswarm/cert-manager-app/pull/367))
- Change the Pod Disruption Budget (PDB) to percentage-based ([#372](https://github.com/giantswarm/cert-manager-app/pull/372))

## [3.4.0] - 2023-09-26

### Added

- Add `Values.global.podSecurityStandards.enforced` flag in preparation of PSP to PSS migration ([#359](https://github.com/giantswarm/cert-manager-app/pull/359))

### Changed

- Enable ServiceMonitor by default. ([#361](https://github.com/giantswarm/cert-manager-app/pull/361))
- Remove control plane node toleration of CA injector deployment. This caused problems on single control plane node clusters. ([#360](https://github.com/giantswarm/cert-manager-app/pull/360))
- Update container image versions to use [v1.12.4](https://github.com/cert-manager/cert-manager/releases/tag/v1.12.4) ([#363](https://github.com/giantswarm/cert-manager-app/pull/363))

## [3.3.0] - 2023-08-29

⚠️ Attention: Major release [3.0.0](#300---2023-07-26) contains breaking changes in user values! Please make yourself familiar with its changelog! ⚠️

### Added

- Add NetworkPolicies for controller and cainjector. ([#354](https://github.com/giantswarm/cert-manager-app/pull/354))

## [3.2.1] - 2023-08-29

⚠️ Attention: Major release [3.0.0](#300---2023-07-26) contains breaking changes in user values! Please make yourself familiar with its changelog! ⚠️

### Changed

- Add missing controller config ConfigMap template. ([#352](https://github.com/giantswarm/cert-manager-app/pull/352))

## [3.2.0] - 2023-08-24

⚠️ Attention: Major release [3.0.0](#300---2023-07-26) contains breaking changes in user values! Please make yourself familiar with its changelog! ⚠️

### Changed

- Make `spec.enableServiceLinks` field configurable for controller, cainjector and webhook Deployments and startupapicheck Job. ([#350](https://github.com/giantswarm/cert-manager-app/pull/350))
- Update chart from upstream. Relevant upstream PRs: [#6241](https://github.com/cert-manager/cert-manager/pull/6241), [#6156](https://github.com/cert-manager/cert-manager/pull/6156), [#6292](https://github.com/cert-manager/cert-manager/pull/6292), [#5337](https://github.com/cert-manager/cert-manager/pull/5337). ([#350](https://github.com/giantswarm/cert-manager-app/pull/350))

## [3.1.0] - 2023-07-27

⚠️ Attention: Major release [3.0.0](#300---2023-07-26) contains breaking changes in user values! Please make yourself familiar with its changelog! ⚠️

### Changed

- Update container image versions to use [v1.12.3](https://github.com/cert-manager/cert-manager/releases/tag/v1.12.3) ([#344](https://github.com/giantswarm/cert-manager-app/pull/344))
- Fix PodDisruptionBudget templates for simultaneous minAvailable and maxUnavailable null values. ([#344](https://github.com/giantswarm/cert-manager-app/pull/344))
- Make resource names less long through default values change. ([#343](https://github.com/giantswarm/cert-manager-app/pull/343))

## [3.0.1] - 2023-07-26

⚠️ Attention: Major release [3.0.0](#300---2023-07-26) contains breaking changes in user values! Please make yourself familiar with its changelog! ⚠️

### Changed

- Explicitly set `ciliumNetworkPolicy.enabled` to `false` in default values. ([#341](https://github.com/giantswarm/cert-manager-app/pull/341))

## [3.0.0] - 2023-07-26

⚠️ Attention: This major release contains breaking changes in user values! ⚠️

We decided to move the helm chart code moves closer to upstream. This means we're pulling in the [helm chart templates from the cert-manager repository](https://github.com/cert-manager/cert-manager/tree/master/deploy/charts/cert-manager) and applying changes to ensure best compatibility to the Giant Swarm clusters.
This results in some breaking changes in the chart values. Please review the [upgrade guide](./docs/upgrading.md) to see if you're affected.

### Changed

- We aligned the chart templates to the [upstream cert-manager chart v1.12.2](https://github.com/cert-manager/cert-manager/tree/master/deploy/charts/cert-manager). Please review the [upgrade guide](./docs/upgrading.md). ([#316](https://github.com/giantswarm/cert-manager-app/pull/316))

## [2.25.0] - 2023-09-26

### Changed

- Remove control plane node toleration of CA injector deployment. This caused problems on single control plane node clusters. ([#362](https://github.com/giantswarm/cert-manager-app/pull/362))
- Update container image versions to use [v1.12.4](https://github.com/cert-manager/cert-manager/releases/tag/v1.12.4)

## [2.24.1] - 2023-06-28

### Added

- Add `cluster-autoscaler safe-to-evict` annotation to `controller` and `cainjector` through newly introduced `controller.podAnnotations` and `cainjector.podAnnotations` values. ([#330](https://github.com/giantswarm/cert-manager-app/pull/330))

## [2.24.0] - 2023-06-26

### Changed

- Add helm adoption annotations to CRD templates. This change is done in preparation of the next major chart release. ([#331](https://github.com/giantswarm/cert-manager-app/pull/331))

## [2.23.2] - 2023-06-19

### Changed

- PodSecurityPolicy: Set `allowedProfileNames`. ([#328](https://github.com/giantswarm/cert-manager-app/pull/328))

## [2.23.1] - 2023-06-19

### Changed

- PodSecurityPolicy: Set `allowedProfileNames`. ([#326](https://github.com/giantswarm/cert-manager-app/pull/326))

## [2.23.0] - 2023-06-14

### Changed

- Update container image versions to use v1.12.1 ([#323](https://github.com/giantswarm/cert-manager-app/pull/323))
- Do not try to install PodSecurityPolicies if not available. This will make the Chart compatible with kubernetes >= 1.25 ([#321](https://github.com/giantswarm/cert-manager-app/pull/321))
- Change security contexts to make the chart work with PSS restricted profile ([#324](https://github.com/giantswarm/cert-manager-app/pull/324)

## [2.22.0] - 2023-06-07

### Changed

- Install `giantswarm-selfsigned` ClusterIssuer regardless of `global.giantSwarmClusterIssuer.install` value. It is required as a default component for Giant Swarm cluster installations.

## [2.21.0] - 2023-04-04

### Added

- Chart: Add `CiliumNetworkPolicy`. ([#301](https://github.com/giantswarm/cert-manager-app/pull/301))

## [2.20.3] - 2023-03-22

### Added

- Add `node-role.kubernetes.io/control-plane` key to list of tolerations

## [2.20.2] - 2023-03-16

### Changed

- Default to new IRSA role for `cert-manager-controller` that has permissions needed for the DNS01 challenge via AWS Route53

## [2.20.1] - 2023-03-16

Not released because of build failure.

## [2.20.0] - 2023-02-20

### Added

- Adds support for DNS01 challenge via AWS Route53 ([#284](https://github.com/giantswarm/cert-manager-app/pull/292))

## [2.19.0] - 2023-02-17

### Added

- Enable `route53` with `IRSA`

## [2.18.2] - 2023-02-10

### Fixed

- controller-psp to allow volumes of type projected for IRSA capability ([#286](https://github.com/giantswarm/cert-manager-app/pull/268))
- Fix indentation when specifying multiple controller extraArgs. ([#284](https://github.com/giantswarm/cert-manager-app/pull/284))

## [2.18.1] - 2023-01-27

### Changed

- Align `controller.serviceAccount` for cert-manager-controller with upstream chart for configurable `controller.serviceAccount.name` and `controller.serviceAccount.annotations`.

## [2.18.0] - 2022-11-14

### Added

- Support for running behind a proxy.
  - `HTTP_PROXY`,`HTTPS_PROXY` and `NO_PROXY` are set as environment variables in `deployment/cert-manager-cainjector`, `deployment/cert-manager-controller` and `deployment/cert-manager-webhook` if defined in `values.yaml`.
- Support for using `cluster-apps-operator` specific `cluster.proxy` values.

## [2.17.1] - 2022-10-11

### Changed

- Align `PodSecurityPolicy` for CRD & `ClusterIssuer` install jobs to actual needs.
- Fix `PodSecurityPolicy` name for CA injector.

## [2.17.0] - 2022-09-22

### Changed

- Rework hooks. ([#263](https://github.com/giantswarm/cert-manager-app/pull/263))
  - Migrate `Chart.yaml` to API version v2.
  - Rename labels.
  - Add `post-upgrade` hook.
  - Move `ClusterIssuer` CRs to helpers.
  - Remove unneccessary hook weights.
  - Refine PSP & RBAC.
  - Improve CRD installation job.
  - Simplify default issuer installation job.
  - Add `values.schema.json` for default isser chart.

## [2.16.0] - 2022-09-12

Before you upgrade to this release, make sure to read the [Upgrading from v1.7 to v1.8](https://cert-manager.io/docs/installation/upgrading/upgrading-1.7-1.8/) document.

### Changed

- Upgrade to upstream image [`v1.8.2`](https://github.com/jetstack/cert-manager/releases/tag/v1.8.2). ([#259](https://github.com/giantswarm/cert-manager-app/pull/259))

## [2.15.3] - 2022-08-22

### Added

- Webhook: Add `PodDisruptionBudget` and pod anti-affinity.
- Startup  API check: Add `NetworkPolicy`.

### Changed

- Webhook: Increase replica count to 2.

## [2.15.2] - 2022-07-27

### Fixed

- RBAC for `cmctl upgrade migrate-api-version` ([#249](https://github.com/giantswarm/cert-manager-app/pull/249)).

## [2.15.1] - 2022-07-07

### Fixed

- Automatically try to execute `cmctl upgrade migrate-api-version` in crd install job to upgrade stored apiversions ([#245](https://github.com/giantswarm/cert-manager-app/pull/245))

## [2.15.0] - 2022-06-24

### Changed

- Upgrade to upstream image [`v1.7.3`](https://github.com/jetstack/cert-manager/releases/tag/v1.7.3) which increases some hard-coded timeouts for certain ACME issuers (ZeroSSL and Sectigo) ([#243](https://github.com/giantswarm/cert-manager-app/pull/243))
- Update kubectl container version to `1.24.2` ([#243](https://github.com/giantswarm/cert-manager-app/pull/243))

## [2.14.0] - 2022-06-20

### Fixed

- Fixed broken relative URLs in the README

### Changed

- Upgrade to upstream image [`v1.7.2`](https://github.com/jetstack/cert-manager/releases/tag/v1.7.2) ([#204](https://github.com/giantswarm/cert-manager-app/pull/214)). This version completely removes cert-manager API versions `v1alpha2, v1alpha3, and v1beta1`. If you need to upgrade your resources, [this document](https://cert-manager.io/docs/installation/upgrading/remove-deprecated-apis/#upgrading-existing-cert-manager-resources) explains the process.
- Update pytest-helm-charts to version [0.7.0](https://github.com/giantswarm/pytest-helm-charts/blob/master/CHANGELOG.md) and adjust dependencies ([#239](https://github.com/giantswarm/cert-manager-app/pull/239))
- Update kubectl container version to `1.24.1` ([#204](https://github.com/giantswarm/cert-manager-app/pull/214))

## [2.13.0] - 2022-04-11

### Changed

- Use retagged container image for HTTP01 AcmeSolver ([#212](https://github.com/giantswarm/cert-manager-app/pull/212))
- Pin kubectl to 1.23.3 in crd-install and clusterissuer-install jobs ([#216](https://github.com/giantswarm/cert-manager-app/pull/216))
- Add `application.giantswarm.io/team` to default labels ([#224](https://github.com/giantswarm/cert-manager-app/pull/224)).

## [2.12.0] - 2021-12-16

### Changed

- Upgrade to upstream image [`v1.6.1`](https://github.com/jetstack/cert-manager/releases/tag/v1.6.1) ([#204](https://github.com/giantswarm/cert-manager-app/pull/204)). This version stops serving cert-manager API versions `v1alpha2, v1alpha3, and v1beta1`. If you need to upgrade your resources, [this document](https://cert-manager.io/docs/installation/upgrading/remove-deprecated-apis/#upgrading-existing-cert-manager-resources) explains the process.

## [2.11.1] - 2021-10-21

### Changed

- Use SVG icon

## [2.11.0] - 2021-10-02

### Changed

- Label default ClusterIssuers with `giantswarm.io/service-type: "managed"` ([#187](https://github.com/giantswarm/cert-manager-app/pull/187)).
- Fix startupjob PSP ([#191](https://github.com/giantswarm/cert-manager-app/pull/191))
- Upgrade to upstream image `v1.5.4` ([#191](https://github.com/giantswarm/cert-manager-app/pull/191))

## [2.10.0] - 2021-09-09

### Changed

- Upgrade to upstream `v1.5.3` ([#184](https://github.com/giantswarm/cert-manager-app/pull/184)). This is the first version compatible with Kubernetes 1.22.
- Add metadata to enable metrics scraping ([#181](https://github.com/giantswarm/cert-manager-app/pull/181)).

## [2.9.0] - 2021-08-18

### Changed

- Update to upstream `v1.4.2` ([#174](https://github.com/giantswarm/cert-manager-app/pull/174)). This deprecates `v1alpha2`, `v1alpha3` and `v1beta1` versions of `cert-manager.io` and `acme.cert-manager.io` CRDs. Further information can be found in the [upstream release notes](https://cert-manager.io/docs/release-notes/release-notes-1.4/) of cert-manager.
- Increase resource requests for the ClusterIssuer and CRD installation Jobs ([#174](https://github.com/giantswarm/cert-manager-app/pull/174)) to prevent timeouts.

## [2.8.0] - 2021-07-26

### Changed

- Label deployments with `giantswarm.io/monitoring_basic_sli: "true"`. ([#171](https://github.com/giantswarm/cert-manager-app/pull/171))
- Migrate values file structure to match `config` repo. ([#172](https://github.com/giantswarm/cert-manager-app/pull/172))

## [2.7.1] - 2021-05-20

### Changed

- Set authoritative nameserver to `coredns` when using `dns01` ACME solver. ([#162](https://github.com/giantswarm/cert-manager-app/pull/162))

## [2.7.0] - 2021-05-05

- Update to upstream `v1.3.1` ([#155](https://github.com/giantswarm/cert-manager-app/pull/155)). This mitigates failed cert-manager-app installations due to CRD conversion issues.

## [2.6.0] - 2021-04-30

### Added

- Add support for *dns01* ACME solver.

## [2.5.0] - 2021-04-19

### Changed

- Update to upstream `v1.2.0`. ([#151](https://github.com/giantswarm/cert-manager-app/pull/151))
- cert-manager-app now requires kubernetes version >=1.16.0. ([#151](https://github.com/giantswarm/cert-manager-app/pull/151))
- Switch rbac rules from `extensions` to `networking.k8s.io` for ingresses. ([#151](https://github.com/giantswarm/cert-manager-app/pull/151))

### Fixed

- Allow strings and integers in values schema for resources requests and limits. ([#150](https://github.com/giantswarm/cert-manager-app/pull/150))

## [2.4.4] - 2021-04-06

### Changed

- Rename clusterissuer subchart to match it's name in its Chart.yaml. ([#140](https://github.com/giantswarm/cert-manager-app/pull/140))
- Make pods of deployments use read-only file systems. ([#140](https://github.com/giantswarm/cert-manager-app/pull/140))
- Make pre-install/pre-upgrade hooks use server side apply. Possibly fixes upgrade timeouts. ([#140](https://github.com/giantswarm/cert-manager-app/pull/140))

## [2.4.3] - 2021-03-26

### Changed

- Set docker.io as the default registry

## [2.4.2] - 2021-01-29

### Added

- Enabled configuration of certificate Secret deletion when the parent Certificate is deleted. ([#127](https://github.com/giantswarm/cert-manager-app/pull/127))

### Changed

- Made CRD install Job backoffLimit configurable (and increased the default value). ([#129](https://github.com/giantswarm/cert-manager-app/pull/129))

## [2.4.1] - 2021-01-19

### Changed

- Made backoffLimit for clusterissuer job configurable. ([#125](https://github.com/giantswarm/cert-manager-app/pull/125))
- Updated clusterissuer subchart API groups to `cert-manager.io/v1`. ([#124](https://github.com/giantswarm/cert-manager-app/pull/124))

## [2.4.0] - 2020-12-22

### Changed

- Update to upstream `v1.1.0`. ([#119](https://github.com/giantswarm/cert-manager-app/pull/119))

## [2.3.3] - 2020-11-23

### Changed

- Schedule hook Jobs on master nodes. ([#106](https://github.com/giantswarm/cert-manager-app/pull/106))

## [2.3.2] - 2020-11-09

### Added

- Added values.schema.json for validation of default values. ([#90](https://github.com/giantswarm/cert-manager-app/pull/90))
- Made cert-manager version configurable. ([#91](https://github.com/giantswarm/cert-manager-app/pull/91))

### Changed

- Updated `cert-manager` to v1.0.4. ([#95](https://github.com/giantswarm/cert-manager-app/pull/95))

### Fixed

- Updated app version in Chart.yaml metadata to `v1.0.3`. ([#91](https://github.com/giantswarm/cert-manager-app/pull/91))

## [2.3.1] - 2020-10-29

### Changed

- Update RBAC API versions. ([#84](https://github.com/giantswarm/cert-manager-app/pull/84))
- Update `cert-manager` to `v1.0.3`. ([#86](https://github.com/giantswarm/cert-manager-app/pull/86))

## [2.3.0] - 2020-10-02

### Changed

- Update `cert-manager` to `v1.0.2` ([#69](https://github.com/giantswarm/cert-manager-app/pull/69))
- Errors from `kubectl` invocation are now surfaced correctly. ([#69](https://github.com/giantswarm/cert-manager-app/pull/69))

## [2.2.5] - 2020-09-29

### Fixed

- Fix `hook-delete-policy` to delete hook resources to make upgrades reliable. ([#76](https://github.com/giantswarm/cert-manager-app/pull/76))

## [2.2.4] - 2020-09-29

### Added

- Add an optional Kiam annotation in case that Route53 wants to be used. ([#71](https://github.com/giantswarm/cert-manager-app/pull/71))

## [2.1.4] - 2020-09-07

### Changed

- Allow clusterissuer subchart to patch clusterissuer resources. ([#65](https://github.com/giantswarm/cert-manager-app/pull/65))

## [2.1.3] - 2020-09-04

### Changed

- Drop resource limits from `ClusterIssuer` subchart to stop it running out of memory. ([#63](https://github.com/giantswarm/cert-manager-app/pull/63))

### Fixed

- Allow `crd-install` job to patch Custom Resource Definitons. ([#62](https://github.com/giantswarm/cert-manager-app/pull/62))

## [2.1.2] - 2020-09-04

### Fixed

- Fix secret name for `orders.acme.cert-manager.io` CR, allowing proper CA injection. ([#60](https://github.com/giantswarm/cert-manager-app/pull/60))

## [2.1.1] - 2020-08-19

### Changed

- Re-introduced hook to allow for CRD install during app upgrades. ([#55](https://github.com/giantswarm/cert-manager-app/pull/55))
- Use unique names for CRD install hooks to avoid naming collisions. ([#56](https://github.com/giantswarm/cert-manager-app/pull/56))

### Added

- Add sub-chart to delete any orphaned resources. ([#56](https://github.com/giantswarm/cert-manager-app/pull/56))

## [2.1.0] - 2020-08-11

### Changed

- Update cert-manager from upstream 0.15.2 to 0.16.1. ([#51](https://github.com/giantswarm/cert-manager-app/pull/51))

## [2.0.2] - 2020-07-29

### Changed

- Install CRDs during chart installation only. ([#46](https://github.com/giantswarm/cert-manager-app/pull/46))

## [2.0.1] - 2020-07-28

### Changed

- Fixed leader election namespace. ([#41](https://github.com/giantswarm/cert-manager-app/pull/41/))
- Template validatingwebhook namespace selector. ([#43](https://github.com/giantswarm/cert-manager-app/pull/43))
- Align CRD templating with the main chart. ([#42](https://github.com/giantswarm/cert-manager-app/pull/42))

### Added

- Add Github release workflow.

## [2.0.0] 2020-07-21

- Upgrade cert-manager from 0.9.0 to 0.15.2 ([#31](https://github.com/giantswarm/cert-manager-app/pull/31))
  - **This is a breaking change**. Please review the upgrade notes [here](https://github.com/giantswarm/cert-manager-app#upgrading-from-v090-giant-swarm-app-v108).
- Upgrade helmclient to 1.0.2
- Upgrade architect-orb to 0.10.0

### Added

- Webhook component to validate requests and prevent incorrect configurations.

## [1.1.0] 2020-10-1

### Added

- Add an optional Kiam annotation for route53.

## [1.0.8] 2020-04-30

### Changed

- Allowed resource requests and limits to be configured with `values.yaml`. ([#24](https://github.com/giantswarm/cert-manager-app/pull/24))

## [1.0.7] 2020-04-09

### Changed

- Fixed sub-chart resources namespace. ([#19](https://github.com/giantswarm/cert-manager-app/pull/19), [#21](https://github.com/giantswarm/cert-manager-app/pull/21))

## [1.0.6] 2020-02-28

### Changed

- Configured app icon. ([#15](https://github.com/giantswarm/cert-manager-app/pull/15))

## [1.0.5] 2020-02-19

### Changed

- Updated helm chart to use a same image registry on parent/subcharts. ([#13](https://github.com/giantswarm/cert-manager-app/pull/13))

## [1.0.4] 2020-01-15

### Changed

- Updated helm chart for clusters with restrictive network policies. ([#9](https://github.com/giantswarm/cert-manager-app/pull/9))

## [1.0.3] 2020-01-03

### Changed

- Updated manifests for Kubernetes 1.16. ([#6](https://github.com/giantswarm/cert-manager-app/pull/6))

## [1.0.2] 2019-12-27

### Changed

- Removed CPU limits. ([#5](https://github.com/giantswarm/cert-manager-app/pull/5))

## [1.0.1] 2019-12-04

### Changed

- Pushed app to both `default` and `giantswarm` catalogs. ([#3](https://github.com/giantswarm/cert-manager-app/pull/3))

## [1.0.0] 2019-10-28

### Added

- `cert-manager` upstream helm chart `v0.9.0`. ([#1](https://github.com/giantswarm/cert-manager-app/pull/1))

[Unreleased]: https://github.com/giantswarm/cert-manager-app/compare/v3.9.4...HEAD
[3.9.4]: https://github.com/giantswarm/cert-manager-app/compare/v3.9.3...v3.9.4
[3.9.3]: https://github.com/giantswarm/cert-manager-app/compare/v3.9.2...v3.9.3
[3.9.2]: https://github.com/giantswarm/cert-manager-app/compare/v3.9.1...v3.9.2
[3.9.1]: https://github.com/giantswarm/cert-manager-app/compare/v3.9.0...v3.9.1
[3.9.0]: https://github.com/giantswarm/cert-manager-app/compare/v3.8.2...v3.9.0
[3.8.2]: https://github.com/giantswarm/cert-manager-app/compare/v3.8.1...v3.8.2
[3.8.1]: https://github.com/giantswarm/cert-manager-app/compare/v3.8.0...v3.8.1
[3.8.0]: https://github.com/giantswarm/cert-manager-app/compare/v3.7.9...v3.8.0
[3.7.9]: https://github.com/giantswarm/cert-manager-app/compare/v3.7.8...v3.7.9
[3.7.8]: https://github.com/giantswarm/cert-manager-app/compare/v3.7.7...v3.7.8
[3.7.7]: https://github.com/giantswarm/cert-manager-app/compare/v3.7.6...v3.7.7
[3.7.6]: https://github.com/giantswarm/cert-manager-app/compare/v3.7.5...v3.7.6
[3.7.5]: https://github.com/giantswarm/cert-manager-app/compare/v3.7.4...v3.7.5
[3.7.4]: https://github.com/giantswarm/cert-manager-app/compare/v3.7.3...v3.7.4
[3.7.3]: https://github.com/giantswarm/cert-manager-app/compare/v3.7.2...v3.7.3
[3.7.2]: https://github.com/giantswarm/cert-manager-app/compare/v3.7.1...v3.7.2
[3.7.1]: https://github.com/giantswarm/cert-manager-app/compare/v3.7.0...v3.7.1
[3.7.0]: https://github.com/giantswarm/cert-manager-app/compare/v3.6.1...v3.7.0
[3.6.1]: https://github.com/giantswarm/cert-manager-app/compare/v3.6.0...v3.6.1
[3.6.0]: https://github.com/giantswarm/cert-manager-app/compare/v3.5.3...v3.6.0
[3.5.3]: https://github.com/giantswarm/cert-manager-app/compare/v3.5.2...v3.5.3
[3.5.2]: https://github.com/giantswarm/cert-manager-app/compare/v3.5.1...v3.5.2
[3.5.1]: https://github.com/giantswarm/cert-manager-app/compare/v3.5.0...v3.5.1
[3.5.0]: https://github.com/giantswarm/cert-manager-app/compare/v3.4.1...v3.5.0
[3.4.1]: https://github.com/giantswarm/cert-manager-app/compare/v3.4.0...v3.4.1
[3.4.0]: https://github.com/giantswarm/cert-manager-app/compare/v3.3.0...v3.4.0
[3.3.0]: https://github.com/giantswarm/cert-manager-app/compare/v3.2.1...v3.3.0
[3.2.1]: https://github.com/giantswarm/cert-manager-app/compare/v3.2.0...v3.2.1
[3.2.0]: https://github.com/giantswarm/cert-manager-app/compare/v3.1.0...v3.2.0
[3.1.0]: https://github.com/giantswarm/cert-manager-app/compare/v3.0.1...v3.1.0
[3.0.1]: https://github.com/giantswarm/cert-manager-app/compare/v3.0.0...v3.0.1
[3.0.0]: https://github.com/giantswarm/cert-manager-app/compare/v2.25.0...v3.0.0
[2.25.0]: https://github.com/giantswarm/cert-manager-app/compare/v2.24.1...v2.25.0
[2.24.1]: https://github.com/giantswarm/cert-manager-app/compare/v2.24.0...v2.24.1
[2.24.0]: https://github.com/giantswarm/cert-manager-app/compare/v2.23.2...v2.24.0
[2.23.2]: https://github.com/giantswarm/cert-manager-app/compare/v2.23.1...v2.23.2
[2.23.1]: https://github.com/giantswarm/cert-manager-app/compare/v2.23.0...v2.23.1
[2.23.0]: https://github.com/giantswarm/cert-manager-app/compare/v2.22.0...v2.23.0
[2.22.0]: https://github.com/giantswarm/cert-manager-app/compare/v2.21.0...v2.22.0
[2.21.0]: https://github.com/giantswarm/cert-manager-app/compare/v2.20.3...v2.21.0
[2.20.3]: https://github.com/giantswarm/cert-manager-app/compare/v2.20.2...v2.20.3
[2.20.2]: https://github.com/giantswarm/cert-manager-app/compare/v2.20.1...v2.20.2
[2.20.1]: https://github.com/giantswarm/cert-manager-app/compare/v2.20.0...v2.20.1
[2.20.0]: https://github.com/giantswarm/cert-manager-app/compare/v2.19.0...v2.20.0
[2.19.0]: https://github.com/giantswarm/cert-manager-app/compare/v2.18.2...v2.19.0
[2.18.2]: https://github.com/giantswarm/cert-manager-app/compare/v2.18.1...v2.18.2
[2.18.1]: https://github.com/giantswarm/cert-manager-app/compare/v2.18.0...v2.18.1
[2.18.0]: https://github.com/giantswarm/cert-manager-app/compare/v2.17.1...v2.18.0
[2.17.1]: https://github.com/giantswarm/cert-manager-app/compare/v2.17.0...v2.17.1
[2.17.0]: https://github.com/giantswarm/cert-manager-app/compare/v2.16.0...v2.17.0
[2.16.0]: https://github.com/giantswarm/cert-manager-app/compare/v2.15.3...v2.16.0
[2.15.3]: https://github.com/giantswarm/cert-manager-app/compare/v2.15.2...v2.15.3
[2.15.2]: https://github.com/giantswarm/cert-manager-app/compare/v2.15.1...v2.15.2
[2.15.1]: https://github.com/giantswarm/cert-manager-app/compare/v2.15.0...v2.15.1
[2.15.0]: https://github.com/giantswarm/cert-manager-app/compare/v2.14.0...v2.15.0
[2.14.0]: https://github.com/giantswarm/cert-manager-app/compare/v2.13.0...v2.14.0
[2.13.0]: https://github.com/giantswarm/cert-manager-app/compare/v2.12.0...v2.13.0
[2.12.0]: https://github.com/giantswarm/cert-manager-app/compare/v2.11.1...v2.12.0
[2.11.1]: https://github.com/giantswarm/cert-manager-app/compare/v2.11.0...v2.11.1
[2.11.0]: https://github.com/giantswarm/cert-manager-app/compare/v2.10.0...v2.11.0
[2.10.0]: https://github.com/giantswarm/cert-manager-app/compare/v2.9.0...v2.10.0
[2.9.0]: https://github.com/giantswarm/cert-manager-app/compare/v2.8.0...v2.9.0
[2.8.0]: https://github.com/giantswarm/cert-manager-app/compare/v2.7.1...v2.8.0
[2.7.1]: https://github.com/giantswarm/cert-manager-app/compare/v2.7.0...v2.7.1
[2.7.0]: https://github.com/giantswarm/cert-manager-app/compare/v2.6.0...v2.7.0
[2.6.0]: https://github.com/giantswarm/cert-manager-app/compare/v2.5.0...v2.6.0
[2.5.0]: https://github.com/giantswarm/cert-manager-app/compare/v2.4.4...v2.5.0
[2.4.4]: https://github.com/giantswarm/cert-manager-app/compare/v2.4.3...v2.4.4
[2.4.3]: https://github.com/giantswarm/cert-manager-app/compare/v2.4.2...v2.4.3
[2.4.2]: https://github.com/giantswarm/cert-manager-app/compare/v2.4.1...v2.4.2
[2.4.1]: https://github.com/giantswarm/cert-manager-app/compare/v2.4.0...v2.4.1
[2.4.0]: https://github.com/giantswarm/cert-manager-app/compare/v2.3.3...v2.4.0
[2.3.3]: https://github.com/giantswarm/cert-manager-app/compare/v2.3.2...v2.3.3
[2.3.2]: https://github.com/giantswarm/cert-manager-app/compare/v2.3.1...v2.3.2
[2.3.1]: https://github.com/giantswarm/cert-manager-app/compare/v2.3.0...v2.3.1
[2.3.0]: https://github.com/giantswarm/cert-manager-app/compare/v2.2.5...v2.3.0
[2.2.5]: https://github.com/giantswarm/cert-manager-app/compare/v2.2.4...v2.2.5
[2.2.4]: https://github.com/giantswarm/cert-manager-app/compare/v2.1.4...v2.2.4
[2.1.4]: https://github.com/giantswarm/cert-manager-app/compare/v2.1.3...v2.1.4
[2.1.3]: https://github.com/giantswarm/cert-manager-app/compare/v2.1.2...v2.1.3
[2.1.2]: https://github.com/giantswarm/cert-manager-app/compare/v2.1.1...v2.1.2
[2.1.1]: https://github.com/giantswarm/cert-manager-app/compare/v2.1.0...v2.1.1
[2.1.0]: https://github.com/giantswarm/cert-manager-app/compare/v2.0.2...v2.1.0
[2.0.2]: https://github.com/giantswarm/cert-manager-app/compare/v2.0.1...v2.0.2
[2.0.1]: https://github.com/giantswarm/cert-manager-app/compare/v2.0.0...v2.0.1
[2.0.0]: https://github.com/giantswarm/cert-manager-app/releases/tag/v1.0.8...v2.0.0
[1.1.0]: https://github.com/giantswarm/cert-manager-app/releases/tag/v1.0.8...v1.1.0
[1.0.8]: https://github.com/giantswarm/cert-manager-app/releases/tag/v1.0.7...v1.0.8
[1.0.7]: https://github.com/giantswarm/cert-manager-app/releases/tag/v1.0.6...v1.0.7
[1.0.6]: https://github.com/giantswarm/cert-manager-app/releases/tag/v1.0.5...v1.0.6
[1.0.5]: https://github.com/giantswarm/cert-manager-app/releases/tag/v1.0.4...v1.0.5
[1.0.4]: https://github.com/giantswarm/cert-manager-app/releases/tag/v1.0.3...v1.0.4
[1.0.3]: https://github.com/giantswarm/cert-manager-app/releases/tag/v1.0.2...v1.0.3
[1.0.2]: https://github.com/giantswarm/cert-manager-app/releases/tag/v1.0.1...v1.0.2
[1.0.1]: https://github.com/giantswarm/cert-manager-app/releases/tag/v1.0.0...v1.0.1
[1.0.0]: https://github.com/giantswarm/cert-manager-app/releases/tag/v1.0.0
