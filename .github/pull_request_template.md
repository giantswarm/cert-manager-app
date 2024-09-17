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

#### Optional app

- [ ] fresh install
- [ ] upgrade from previous version

#### Pre-installed app

- [ ] fresh install during cluster creation
- [ ] upgrade from previous version in a pre-existing cluster

#### Other testing

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
