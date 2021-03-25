# cert-manager

This chart installs cert-manager as a managed application. `cert-manager` is a native Kubernetes certificate management controller. It can help with issuing certificates from a variety of sources, such as Let’s Encrypt, HashiCorp Vault, Venafi, a simple signing keypair, or self-signed.


## Configuration

The following table lists the configurable parameters of the cert-manager chart, its dependencies and default values.

| Parameter                                | Description                                         | Default                                |
| ---------------------------------------- | --------------------------------------------------- | -------------------------------------- |
| `cainjector.extraArgs`                   | Additional args to pass to the cainjector container | `[]`                                   |
| `cainjector.image.pullPolicy`            | Cainjector image pull policy                        | `"IfNotPresent"`                       |
| `cainjector.logLevel`                    | Cainjector log level                                | `2`                                    |
| `cainjector.replicas`                    | Cainjector replica count                            | `1`                                    |
| `cainjector.resources.requests.cpu`      | Cainjector CPU request                              | `"10m"`                                |
| `cainjector.resources.requests.memory`   | Cainjector memory request                           | `"32Mi"`                               |
| `controller.defaultIssuer.group`         | Default Issuer group                                | `"cert-manager.io"`                    |
| `controller.defaultIssuer.kind`          | Default Issuer kind                                 | `"ClusterIssuer"`                      |
| `controller.defaultIssuer.name`          | Default Issuer name                                 | `"letsencrypt-giantswarm"`             |
| `controller.extraArgs`                   | Additional args to pass to the controller container | `[]`                                   |
| `controller.image.pullPolicy`            | Controller image pull policy                        | `"IfNotPresent"`                       |
| `controller.logLevel`                    | Controller log level                                | `2`                                    |
| `controller.replicas`                    | Controller replica count                            | `1`                                    |
| `controller.resources.requests.cpu`      | Controller CPU request                              | `"50m"`                                |
| `controller.resources.requests.memory`   | Controller memory request                           | `"100Mi"`                              |
| `crds.image.pullPolicy`                  | CRD job image pull policy                           | `"IfNotPresent"`                       |
| `crds.install`                           | Enable CRD installation                             | `true`                                 |
| `crds.resources.requests.cpu`            | CRD job CPU request                                 | `"50m"`                                |
| `crds.resources.requests.memory`         | CRD job memory request                              | `"100Mi"`                              |
| `global.giantSwarmClusterIssuer.install` | Install Giant Swarm ClusterIssuer                   | `true`                                 |
| `global.image.registry`                  | Image registry                                      | `"quay.io"`                            |
| `global.name`                            | Application name                                    | `"cert-manager"`                       |
| `global.securityContext.groupID`         | Group ID to run containers as                       | `1000`                                 |
| `global.securityContext.userID`          | User ID to run containers as                        | `1000`                                 |
| `prometheus.enabled`                     | Enable Prometheus endpoint                          | `true`                                 |
| `webhook.extraArgs`                      | Additional args to pass to the webhook container    | `[]`                                   |
| `webhook.image.pullPolicy`               | Webhook image pull policy                           | `"IfNotPresent"`                       |
| `webhook.logLevel`                       | Webhook log level                                   | `2`                                    |
| `webhook.replicas`                       | Webhook replica count                               | `1`                                    |
| `webhook.resources.requests.cpu`         | Webhook CPU request                                 | `"20m"`                                |
| `webhook.resources.requests.memory`      | Webhook memory request                              | `"50Mi"`                               |
| `webhook.securePort`                     | Webhook container listen port                       | `10250`                                |
