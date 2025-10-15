# cert-manager-app E2E Tests

This directory contains end-to-end tests for the cert-manager-app using the [Giant Swarm apptest-framework](https://github.com/giantswarm/apptest-framework).

## Test Structure

The tests are organized into multiple suites to cover different scenarios:

```
tests/e2e/
├── config.yaml                    # Global test configuration
├── go.mod                         # Go dependencies
├── go.sum                         # Go dependency checksums  
├── suites/
│   ├── basic/                      # Basic functionality tests
│   │   ├── config.yaml             # Suite-specific config
│   │   ├── basic_suite_test.go
│   │   └── values.yaml
│   ├── advanced/                   # Advanced certificate scenarios
│   │   ├── config.yaml             # Suite-specific config
│   │   ├── advanced_suite_test.go
│   │   └── values.yaml
│   └── mctest/                     # Management cluster tests
│       ├── config.yaml             # MC-specific config (isMCTest: true)
│       ├── mctest_suite_test.go
│       └── values.yaml
└── README.md                       # This file
```

## Test Suites

### Basic Suite (`suites/basic/`)
- Tests core cert-manager deployment and functionality
- Validates all deployments are ready (controller, webhook, cainjector)
- Checks ClusterIssuer creation (letsencrypt-giantswarm, selfsigned-giantswarm)
- Tests basic certificate issuance using self-signed issuer

### Advanced Suite (`suites/advanced/`)
- Tests advanced certificate scenarios
- Namespace-scoped Issuer creation
- Certificate renewal scenarios
- Multi-DNS name certificates
- Webhook validation testing

### Management Cluster Suite (`suites/mctest/`)
- Tests cert-manager deployment in management clusters
- High availability scenarios
- Management cluster specific certificate issuance
- Cross-cluster issuer availability

## Configuration

### Suite Configuration (`config.yaml`)

Each test suite requires its own `config.yaml` file in its directory. This allows per-suite customization of providers and test type.

**Workload Cluster Suites** (`basic/`, `advanced/`):
```yaml
appName: cert-manager-app
repoName: cert-manager-app  
appCatalog: giantswarm
providers:
  - capa
  - capv
isMCTest: false
```

**Management Cluster Suite** (`mctest/`):
```yaml
appName: cert-manager-app
repoName: cert-manager-app  
appCatalog: giantswarm
providers:
  - capa  # MC tests currently only support CAPA
isMCTest: true  # This flag tells the framework to use an MC
```

### Test Suite Values
Each test suite includes a `values.yaml` file with specific configuration:
- **Basic**: Minimal resource configuration for core functionality testing
- **Advanced**: Higher resources and additional features enabled
- **MC Test**: High availability configuration optimized for management clusters

## Running Tests

### Prerequisites

1. **Install Ginkgo**:
   ```bash
   go install github.com/onsi/ginkgo/v2/ginkgo@latest
   ```

2. **Set Environment Variables**:
   ```bash
   export E2E_KUBECONFIG="./kube/e2e.yaml"          # Path to management cluster kubeconfig
   export E2E_KUBECONFIG_CONTEXT="capa"             # Kubeconfig context to use
   export E2E_APP_VERSION="3.9.3"                  # Version of cert-manager-app to test
   ```

3. **Optional - Reuse Existing Workload Cluster**:
   ```bash
   export E2E_WC_NAME="test-cluster"                # Existing workload cluster name
   export E2E_WC_NAMESPACE="org-example"           # Workload cluster namespace
   export E2E_WC_KEEP="true"                       # Skip deleting cluster after tests
   ```

### Running Individual Test Suites

Navigate to the e2e directory:
```bash
cd tests/e2e
```

Run specific test suites:
```bash
# Basic functionality tests
ginkgo --timeout 4h -v -r ./suites/basic/

# Advanced certificate scenarios
ginkgo --timeout 4h -v -r ./suites/advanced/

# Management cluster tests
ginkgo --timeout 4h -v -r ./suites/mctest/
```

### Running All Test Suites

Run all workload cluster tests:
```bash
ginkgo --timeout 4h -v -r ./suites/basic/ ./suites/advanced/
```

Run management cluster tests separately:
```bash  
ginkgo --timeout 4h -v -r ./suites/mctest/
```

### Running Tests with Local apptest-framework Changes

If you need to test with local changes to apptest-framework:

1. Add a replace directive to `go.mod`:
   ```go
   replace github.com/giantswarm/apptest-framework => /path/to/your/apptest-framework
   ```

2. Run `go mod tidy` to update dependencies

3. Run tests as normal

## Running Tests in CI/CD

### Triggering Tests on Pull Requests

Add a comment to any PR to trigger tests:

```
/run app-test-suites
```

Run tests for specific provider:
```
/run app-test-suites-single PROVIDER=capa
```

Run specific test suites:
```
/run app-test-suites-single PROVIDER=capa TARGET_SUITES=basic,advanced
```

Run management cluster tests:
```
/run app-test-suites-single PROVIDER=capa TARGET_SUITES=mctest
```

### Test Configuration

The CI pipeline uses the configuration from `config.yaml` and will:
- Create ephemeral workload clusters for workload cluster tests
- Create ephemeral management clusters for management cluster tests (when `isMCTest: true`)
- Install the specified version of cert-manager-app
- Run the configured test suites
- Clean up resources after testing

## Test Development

### Adding New Test Cases

1. **Choose the appropriate suite** or create a new one:
   - `basic/` - Core functionality that should always work
   - `advanced/` - Complex scenarios or edge cases
   - `mctest/` - Management cluster specific tests

2. **Follow the existing patterns**:
   ```go
   var _ = Describe("Feature description", func() {
       It("should behave as expected", func() {
           // Test implementation
       })
   })
   ```

3. **Use proper cleanup**:
   ```go
   AfterEach(func() {
       // Clean up test resources
   })
   ```

### Best Practices

1. **Use descriptive test names** that explain what is being tested
2. **Include proper timeouts** for Eventually/Consistently assertions
3. **Clean up resources** created during tests
4. **Use dynamic clients** for cert-manager CRDs (Certificate, Issuer, etc.)
5. **Test both positive and negative scenarios**
6. **Include proper error handling** and meaningful error messages

### Testing cert-manager Specific Resources

Use the Kubernetes dynamic client to interact with cert-manager CRDs:

```go
// Define the GVR for certificates
certificateGVR := schema.GroupVersionResource{
    Group:    "cert-manager.io",
    Version:  "v1", 
    Resource: "certificates",
}

// Create a certificate
dynamicClient := ts.GetClient().Dynamic()
_, err := dynamicClient.Resource(certificateGVR).Namespace(namespace).Create(ctx, certificate, metav1.CreateOptions{})
```

## Troubleshooting

### Common Issues

1. **"No test suites found" in CI**:
   - **Cause**: Missing `config.yaml` in test suite directory
   - **Solution**: Ensure each suite has its own `config.yaml` with provider configuration
   - **Example**: See `suites/basic/config.yaml`, `suites/advanced/config.yaml`, `suites/mctest/config.yaml`

2. **"Skipping [suite] as not configured for provider"**:
   - **Cause**: Test suite's `config.yaml` doesn't include the provider being tested
   - **Solution**: Add the provider to the `providers` list in the suite's `config.yaml`
   - **Example**: Add `- capz` to the providers array to enable Azure testing

3. **Timeout errors**: Increase timeouts in Eventually/Consistently assertions

4. **Resource cleanup**: Ensure AfterEach blocks properly clean up test resources

5. **CRD availability**: Wait for cert-manager CRDs to be available before creating resources

6. **Webhook readiness**: Ensure webhook is ready before creating cert-manager resources

### Debugging

1. **Check app status**:
   ```bash
   kubectl get apps -A
   kubectl describe app cert-manager-app -n default
   ```

2. **Check cert-manager pods**:
   ```bash
   kubectl get pods -n default -l app.kubernetes.io/name=cert-manager
   kubectl logs -n default deployment/cert-manager-app
   ```

3. **Check certificates and issuers**:
   ```bash
   kubectl get certificates,clusterissuers -A
   kubectl describe certificate test-cert -n default
   ```

4. **Verbose test output**:
   ```bash
   ginkgo -v --progress ./suites/basic/
   ```

## Resources

- [apptest-framework documentation](https://github.com/giantswarm/apptest-framework)  
- [Ginkgo testing framework](https://onsi.github.io/ginkgo/)
- [cert-manager documentation](https://cert-manager.io/docs/)
- [Giant Swarm cluster-standup-teardown](https://github.com/giantswarm/cluster-standup-teardown)
