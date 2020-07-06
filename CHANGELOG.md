# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project's packages adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Upgrade helmclient to 1.0.2
- Upgrade architect-orb to 0.10.0

## [v1.0.8] 2020-04-30

### Changed

- Allowed resource requests and limits to be configured with `values.yaml`. ([#24](https://github.com/giantswarm/cert-manager-app/pull/24))

## [v1.0.7] 2020-04-09

### Changed

- Fixed sub-chart resources namespace. ([#19](https://github.com/giantswarm/cert-manager-app/pull/19), [#21](https://github.com/giantswarm/cert-manager-app/pull/21))

## [v1.0.6] 2020-02-28

### Changed

- Configured app icon. ([#15](https://github.com/giantswarm/cert-manager-app/pull/15))

## [v1.0.5] 2020-02-19

### Changed

- Updated helm chart to use a same image registry on parent/subcharts. ([#13](https://github.com/giantswarm/cert-manager-app/pull/13))

## [v1.0.4] 2020-01-15

### Changed

- Updated helm chart for clusters with restrictive network policies. ([#9](https://github.com/giantswarm/cert-manager-app/pull/9))

## [v1.0.3] 2020-01-03

### Changed

- Updated manifests for Kubernetes 1.16. ([#6](https://github.com/giantswarm/cert-manager-app/pull/6))

## [v1.0.2] 2019-12-27

### Changed

- Removed CPU limits. ([#5](https://github.com/giantswarm/cert-manager-app/pull/5))

## [v1.0.1] 2019-12-04

### Changed

- Pushed app to both `default` and `giantswarm` catalogs. ([#3](https://github.com/giantswarm/cert-manager-app/pull/3))

## [v1.0.0] 2019-10-28

### Added

- `cert-manager` upstream helm chart `v0.9.0`. ([#1](https://github.com/giantswarm/cert-manager-app/pull/1))

[Unreleased]: https://github.com/giantswarm/cert-manager-app/compare/v1.0.8...master
[v1.0.8]: https://github.com/giantswarm/cert-manager-app/releases/tag/v1.0.8
[v1.0.7]: https://github.com/giantswarm/cert-manager-app/releases/tag/v1.0.7
[v1.0.6]: https://github.com/giantswarm/cert-manager-app/releases/tag/v1.0.6
[v1.0.5]: https://github.com/giantswarm/cert-manager-app/releases/tag/v1.0.5
[v1.0.4]: https://github.com/giantswarm/cert-manager-app/releases/tag/v1.0.4
[v1.0.3]: https://github.com/giantswarm/cert-manager-app/releases/tag/v1.0.3
[v1.0.2]: https://github.com/giantswarm/cert-manager-app/releases/tag/v1.0.2
[v1.0.1]: https://github.com/giantswarm/cert-manager-app/releases/tag/v1.0.1
[v1.0.0]: https://github.com/giantswarm/cert-manager-app/releases/tag/v1.0.0
