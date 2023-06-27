import logging
from contextlib import contextmanager
from pathlib import Path
from typing import Dict, List

import pykube
import pytest
from pytest_helm_charts.clusters import Cluster
from pytest_helm_charts.k8s.deployment import wait_for_deployments_to_run

logger = logging.getLogger(__name__)

namespace_name = "default"

timeout: int = 360


# scope "module" means this is run only once, for the first test case requesting! It might be tricky
# if you want to assert this multiple times
@pytest.fixture(scope="module")
def app_deployment(kube_cluster: Cluster) -> List[pykube.Deployment]:
    deployments = wait_for_deployments_to_run(
        kube_cluster.kube_client,
        ["cert-manager-app-cainjector", "cert-manager-app", "cert-manager-app-webhook"],
        "default",
        timeout,
    )
    return deployments


# when we start the tests on circleci, we have to wait for pods to be available, hence
# this additional delay and retries
@pytest.mark.smoke
@pytest.mark.upgrade
@pytest.mark.flaky(reruns=5, reruns_delay=10)
def test_pods_available(kube_cluster: Cluster, app_deployment: List[pykube.Deployment]):
    for d in app_deployment:
        assert int(d.obj["status"]["readyReplicas"]) > 0


@pytest.mark.smoke
@pytest.mark.flaky(reruns=5, reruns_delay=10)
def test_clusterissuer_available(kube_cluster: Cluster):
    cluster_issuers = kube_cluster.kubectl("get clusterissuers")

    cluster_issuer_names = [cl["metadata"]["name"] for cl in cluster_issuers]

    assert "letsencrypt-giantswarm" in cluster_issuer_names
    assert "selfsigned-giantswarm" in cluster_issuer_names


# Using smoke here, because redeployment takes too much time
@pytest.mark.smoke
def test_self_signed_certificates(request, kube_cluster: Cluster):
    # Request self signed certificates and check if they get Ready
    kube_cluster.kubectl("apply", filename=Path(request.fspath.dirname) / "selfsigned.yaml", output_format="")

    kube_cluster.kubectl(f"wait certificate test-ca test --for=condition=Ready", timeout="60s", output_format="")
