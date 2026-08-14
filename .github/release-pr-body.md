<!--
Used verbatim as the body of release PRs created by giantswarm/github-workflows
(see .github/workflows/zz_generated.create_release_pr.yaml).

The `/run` line makes Tekton start the E2E tests as soon as the release PR is
opened. app-test-suites runs every suite for each configured provider by
default, so the full set runs before the release is tagged.

See https://github.com/giantswarm/roadmap/issues/4334
-->

The full set of E2E test suites for this app runs automatically on this release PR.

/run app-test-suites
