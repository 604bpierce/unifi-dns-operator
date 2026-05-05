# Contributing to Unifi DNS Operator

Thank you for your interest in contributing to the Unifi DNS Operator!

## Development Setup

### Prerequisites

- Go 1.25+
- Docker
- kubectl
- Kubebuilder
- Access to a Kubernetes cluster
- Unifi Network Controller (for testing)

### Local Development

1. Clone the repository:
```bash
git clone https://github.com/bpierce/unifi-dns-operator.git
cd unifi-dns-operator
```

2. Install dependencies:
```bash
go mod download
```

3. Run tests:
```bash
make test
```

4. Build the operator:
```bash
make build
```

5. Run locally:
```bash
export UNIFI_API_URL=https://your-controller/proxy/network/integration
export UNIFI_SITE_ID=your-site-id
export UNIFI_API_KEY=your-api-key
make run
```

## Making Changes

1. Create a new branch:
```bash
git checkout -b feature/your-feature-name
```

2. Make your changes and add tests

3. Run tests and linting:
```bash
make test
make vet
make fmt
```

4. Commit your changes:
```bash
git add .
git commit -m "Description of your changes"
```

5. Push to your fork and create a pull request

## Pull Request Guidelines

- Write clear, descriptive commit messages
- Include tests for new functionality
- Update documentation as needed
- Ensure all tests pass
- Keep PRs focused on a single feature or fix

## Creating a Release

Releases are automated via GitHub Actions. To create a new release:

1. **Update the CHANGELOG.md** with your changes

2. **Create a release via GitHub Actions:**
   - Go to Actions → Release → Run workflow
   - Enter the version (e.g., `v0.1.0`)
   - Click "Run workflow"

3. **The automation will:**
   - Create and push the git tag
   - Build multi-arch Docker images (amd64, arm64)
   - Push images to GitHub Container Registry
   - Generate deployment manifests
   - Create a GitHub release with artifacts

4. **Verify the release:**
   - Check the [Releases page](https://github.com/bpierce/unifi-dns-operator/releases)
   - Verify the container image at [ghcr.io](https://github.com/bpierce/unifi-dns-operator/pkgs/container/unifi-dns-operator)

## Release Checklist

- [ ] All tests passing
- [ ] CHANGELOG.md updated
- [ ] Documentation updated
- [ ] Version number follows [Semantic Versioning](https://semver.org/)
- [ ] GitHub Actions workflow runs successfully
- [ ] Container image builds for both amd64 and arm64
- [ ] Release artifacts generated correctly

## Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Run `go vet` to catch common issues
- Add comments for exported functions and types

## Testing

### Unit Tests

Unit tests run in CI on every push and PR:
```bash
make test
```

These tests validate:
- CRD creation and validation
- Controller reconciliation logic
- Status updates and conditions

### E2E Tests

End-to-end tests require a real Unifi controller and are **not run in CI**.

To run E2E tests locally:
```bash
# Set your Unifi credentials
export UNIFI_API_URL="https://your-controller/proxy/network/integration"
export UNIFI_SITE_ID="your-site-id"
export UNIFI_API_KEY="your-api-key"

# Run E2E tests (requires Kind installed)
make test-e2e
```

See [test/e2e/README.md](test/e2e/README.md) for details.

### Testing Guidelines

- Write unit tests for new functionality
- Test with actual Unifi controller when possible
- Ensure backward compatibility
- Don't commit credentials or sensitive data

## Questions?

Feel free to open an issue for:
- Bug reports
- Feature requests
- Questions about development
- Documentation improvements

## License

By contributing, you agree that your contributions will be licensed under the Apache 2.0 License.
