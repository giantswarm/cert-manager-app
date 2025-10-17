<!--
Not all PRs will require all tests to be carried out. Refer to the
testing doc below and delete where appropriate.

https://intranet.giantswarm.io/docs/dev-and-releng/app-developer-processes/cert-manager/
-->

<!--
@team-shield will be automatically requested for review once
this PR has been submitted.
-->

This PR:

- adds/changes/removes etc

### Testing

#### Automated E2E Tests

Run automated tests across all providers:
```
/run app-test-suites
```

This will test both **fresh install** and **upgrade** scenarios on CAPA, CAPV, CAPZ, and CAPVCD.

**Test Suites**:
- [ ] `basic` - Fresh install validation (deployments, ClusterIssuers, certificate issuance)
- [ ] `upgrade` - Upgrade validation (certificate persistence, post-upgrade issuance)

See [tests/e2e/README.md](../tests/e2e/README.md) for details.

#### Manual Testing (Optional)

##### Optional app

- [ ] fresh install
- [ ] upgrade from previous version

##### Pre-installed app

- [ ] fresh install during cluster creation
- [ ] upgrade from previous version in a pre-existing cluster

##### Other testing

<!--
Install ingress-nginx and hello-world to obtain a certificate,
then upgrade the cert-manager-app and ensure the CRs are still reconciled after the upgrade.
-->

- [ ] check reconciliation of existing resources after upgrading

<!--
Changelog must always be updated.
-->

### Checklist

- [ ] Update changelog in CHANGELOG.md.
