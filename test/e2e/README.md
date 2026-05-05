# E2E Tests

## Overview

End-to-end tests validate the operator's functionality in a real Kubernetes cluster with actual Unifi controller integration.

## Prerequisites

**These tests require:**
1. A running Unifi Network Controller (self-hosted or cloud)
2. Valid API credentials (API key and site ID)
3. Kind installed locally
4. Network access to the Unifi controller

## Running E2E Tests

### Local Execution

1. **Set up credentials:**
   ```bash
   export UNIFI_API_URL="https://your-controller/proxy/network/integration"
   export UNIFI_SITE_ID="your-site-id"
   export UNIFI_API_KEY="your-api-key"
   ```

2. **Run the tests:**
   ```bash
   make test-e2e
   ```

   This will:
   - Create a Kind cluster
   - Install CertManager
   - Deploy the operator
   - Run E2E tests
   - Clean up the cluster

### GitHub Actions

E2E tests are **not run automatically** in CI because they require real Unifi credentials.

To run them manually:
1. Go to Actions → E2E Tests → Run workflow
2. Note: This will fail unless you've configured secrets in the repository

## Why Not in CI?

E2E tests are disabled in CI because:
- They require actual Unifi controller access
- API credentials can't be safely provided in public CI
- They take longer to run than unit tests
- Unit tests already cover core functionality

## What E2E Tests Validate

- Operator deployment and pod startup
- CRD installation
- Controller manager health
- Metrics endpoint
- Integration with cert-manager

## Adding E2E Tests for Unifi Integration

To add tests that actually interact with Unifi:

1. Create a test secret with credentials before deployment
2. Verify DNS entries are created in Unifi
3. Verify updates sync correctly
4. Verify deletions remove entries from Unifi

Example test structure:
```go
It("should create DNS entry in Unifi", func() {
    // Create UnifiDNSEntry resource
    // Verify it appears in Unifi controller
    // Verify status is updated with policyId
})
```

## Troubleshooting

**Pod stuck in CreateContainerConfigError:**
- This means the credentials secret doesn't exist
- Create the secret in the test cluster before running tests

**Tests timeout:**
- Verify network access to Unifi controller
- Check API credentials are valid
- Ensure Kind cluster can reach your controller

## Future Improvements

- Mock Unifi API server for testing without real controller
- Integration test suite with docker-compose Unifi setup
- Automated E2E tests with test credentials in a separate workflow
