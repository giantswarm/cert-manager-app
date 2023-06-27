# How to contribute

This repository contains the source of the cert-manager app for the Giant Swarm App platform.
It is Apache 2.0 licensed and accepts contributions via GitHub pull requests.

This document outlines some of the conventions on commit message formatting, contact points for developers and other resources to make getting your contribution into cert-manager easier.

## Getting started

This repository is managed using [vendir][vendir]. This means all files located in and below `helm/cert-manager` are sourced from our [upstream fork](https://github.com/giantswarm/cert-manager-upstream). Exceptions are subcharts located in `helm/cert-manager/charts`.

If you want to make a change to the main helm chart, you can either

- Contribute to the upstream helm chart directly at [https://github.com/cert-manager/cert-manager](https://github.com/cert-manager/cert-manager). We'll copy in the changes once merged there.
- Contribute to the Giant Swarm fork [https://github.com/giantswarm/cert-manager-upstream][cert-manager-upstream]. Do this if you think your change is urgent and needs to be included in this repository immediately.

This chart can be tested locally on a [kind][kind] cluster using [app-build-suite](https://github.com/giantswarm/app-build-suite) and [app-test-suite](https://github.com/giantswarm/app-test-suite).

## Development workflow

Files in `helm/cert-manager` are managed using [vendir][vendir]. Exceptions are subcharts located in `helm/cert-manager/charts`. Any direct changes to these files will be lost the next time the chart will be synced to upstream.

To make changes to these files, follow the instructions in the ["Contributions" section](https://github.com/giantswarm/cert-manager-upstream#contributions) in [the upstream fork][cert-manager-upstream].

To test your changes locally, you can use a [kind][kind] cluster, build the chart using [app-build-suite][app-build-suite] then test it using [app-test-suite][app-test-suite].

- Install [kind][kind]
- Download [`dabs.sh`](https://github.com/giantswarm/app-build-suite/releases/download/v1.1.4/dabs.sh)
- Download [`dats.sh`](https://github.com/giantswarm/app-test-suite/releases/download/v0.4.1/dats.sh)

Executing `./dabs.sh -c helm/cert-manager` in the root of the repository will launch a container executing app-build-suite. The result will be a `.tgz` file inside of a `build` directory. It means the charts syntax is valid.

Copy the `.tgz` file from inside the `build` directory into the root of the repository.

Execute `./dats.sh <name-of-the-chart-archive.tgz>` will spin up a temporary kind cluster, then execute the tests located in `tests/ats` against the installed chart.

Check out `.ats/main.yaml` for instructions for running against externally created clusters. This can help debugging failed tests.

Once you create PR, app-build-suite and app-test-suite will be executed again as part of CI.

## Update from upstream

- Prepare upstream fork
  - Clone repository https://github.com/giantswarm/cert-manager-upstream and follow ["Update from upstream" instructions in README](https://github.com/giantswarm/cert-manager-upstream#update-from-upstream)
- In this repository
  - Switch to a fresh branch
  - Run `APPLICATION=helm/cert-manger make sync-chart`
  - Review changes in upstream repository, then decide if there are additional changes to `values.yaml` or `Chart.yaml` are required. (The container image tag is derived from the `appVersion` in `Chart.yaml`)
  - Make sure the chart still works as intended (See [Development workflow](#development-workflow))
  - Commit & push changes
  - Create a PR

## Reporting Bugs and Creating Issues

Reporting bugs is one of the best ways to contribute. If you find bugs or documentation mistakes in the cert-manager project, please let us know by [opening an issue](https://github.com/giantswarm/cert-manager-app/issues/new). We treat bugs and mistakes very seriously and believe no issue is too small. Before creating a bug report, please check there that one does not already exist.

To make your bug report accurate and easy to understand, please try to create bug reports that are:

- Specific. Include as much details as possible: which version, what environment, what configuration etc. You can also attach logs.

- Reproducible. Include the steps to reproduce the problem. We understand some issues might be hard to reproduce, please includes the steps that might lead to the problem. If applicable, you can also attach affected data dir(s) and a stack trace to the bug report.

- Isolated. Please try to isolate and reproduce the bug with minimum dependencies. It would significantly slow down the speed to fix a bug if too many dependencies are involved in a bug report. Debugging external systems that rely on cert-manager is out of scope, but we are happy to point you in the right direction or help you interact with cert-manager in the correct manner.

- Unique. Do not duplicate existing bug reports.

- Scoped. One bug per report. Do not follow up with another bug inside one report.

You might also want to read [Elika Etemad’s article on filing good bug reports](http://fantasai.inkedblade.net/style/talks/filing-good-bugs/) before creating a bug report.

We might ask you for further information to locate a bug. A duplicated bug report will be closed.

## Contribution flow

This is a rough outline of what a contributor's workflow looks like:

- Create a feature branch from where you want to base your work. This is usually main.
- Make commits of logical units.
- Make sure your commit messages are in the proper format (see below).
- Push your changes to a topic branch in your fork of the repository.
- Submit a pull request to giantswarm/cert-manager-app.
- Adding unit tests will greatly improve the chance for getting a quick review and your PR accepted.
- Your PR must receive a LGTM from a maintainer.
- We merge your PR and release a new version eventually

Thanks for your contributions!

### Format of the Commit Message

We follow a rough convention for commit messages that is designed to answer two
questions: what changed and why. The subject line should feature the what and
the body of the commit should describe the why.

[vendir]: https://carvel.dev/vendir/
[cert-manager-upstream]: https://github.com/giantswarm/cert-manager-upstream
[kind]: https://kind.sigs.k8s.io/
[app-build-suite]: https://github.com/giantswarm/app-build-suite
[app-test-suite]: https://github.com/giantswarm/app-test-suite
