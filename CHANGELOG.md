# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project's packages adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Update `cert-manager` to `v1.0.3`. ([#86](https://github.com/giantswarm/cert-manager-app/pull/86))

## [2.3.0] - 2020-10-02

### Changed

- Update `cert-manager` to `v1.0.2` ([#69](https://github.com/giantswarm/cert-manager-app/pull/69))
- Errors from `kubectl` invocation are now surfaced correctly. ([#69](https://github.com/giantswarm/cert-manager-app/pull/69))

## [2.2.5] - 2020-09-29

### Fixed

- Fix `hook-delete-policy` to delete hook resources to make upgrades reliable.

## [2.2.4] - 2020-09-29

### Added

- Add an optional Kiam annotation in case that Route53 wants to be used.

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

[Unreleased]: https://github.com/giantswarm/cert-manager-app/compare/v2.3.0...HEAD
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
[1.0.8]: https://github.com/giantswarm/cert-manager-app/releases/tag/v1.0.7...v1.0.8
[1.0.7]: https://github.com/giantswarm/cert-manager-app/releases/tag/v1.0.6...v1.0.7
[1.0.6]: https://github.com/giantswarm/cert-manager-app/releases/tag/v1.0.5...v1.0.6
[1.0.5]: https://github.com/giantswarm/cert-manager-app/releases/tag/v1.0.4...v1.0.5
[1.0.4]: https://github.com/giantswarm/cert-manager-app/releases/tag/v1.0.3...v1.0.4
[1.0.3]: https://github.com/giantswarm/cert-manager-app/releases/tag/v1.0.2...v1.0.3
[1.0.2]: https://github.com/giantswarm/cert-manager-app/releases/tag/v1.0.1...v1.0.2
[1.0.1]: https://github.com/giantswarm/cert-manager-app/releases/tag/v1.0.0...v1.0.1
[1.0.0]: https://github.com/giantswarm/cert-manager-app/releases/tag/v1.0.0
