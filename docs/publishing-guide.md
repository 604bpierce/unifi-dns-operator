# Publishing Guide

This guide explains how to publish releases of the Unifi DNS Operator to GitHub Container Registry.

## Prerequisites

Before you can publish releases, ensure you have:

1. **GitHub Repository** - Your code pushed to GitHub
2. **GitHub Actions enabled** - Actions should be enabled for your repository
3. **Package permissions** - Ensure GitHub Actions has permission to write packages:
   - Go to Settings → Actions → General
   - Under "Workflow permissions", select "Read and write permissions"

## Publishing Your First Release

### Step 1: Push Your Code to GitHub

```bash
cd /Users/bpierce/Projects/bx-network/unifi-dns-operator
git init
git add .
git commit -m "Initial commit of Unifi DNS Operator"
git branch -M main
git remote add origin https://github.com/bpierce/unifi-dns-operator.git
git push -u origin main
```

### Step 2: Create Your First Release

There are two ways to create a release:

#### Option A: Via GitHub Actions (Recommended)

1. Go to your repository on GitHub
2. Click on "Actions" tab
3. Select "Release" workflow
4. Click "Run workflow"
5. Enter version: `v0.1.0`
6. Click "Run workflow"

The automation will:
- Create and push the git tag
- Build Docker images for amd64 and arm64
- Push images to ghcr.io/bpierce/unifi-dns-operator
- Generate deployment manifests
- Create a GitHub release with downloadable assets

#### Option B: Manual Tag Creation

```bash
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
```

The build-and-publish workflow will automatically trigger and create the release.

### Step 3: Verify the Release

1. **Check the GitHub Release:**
   - Go to https://github.com/bpierce/unifi-dns-operator/releases
   - Verify the release was created with deployment manifests

2. **Check the Container Image:**
   - Go to https://github.com/bpierce?tab=packages
   - Find `unifi-dns-operator`
   - Verify the image was published

3. **Test the Installation:**
   ```bash
   kubectl apply -f https://github.com/bpierce/unifi-dns-operator/releases/download/v0.1.0/unifi-dns-operator-v0.1.0.yaml
   ```

## Continuous Deployment

Every push to the `main` branch will:
- Run tests
- Build and push a `:latest` tagged image
- Build and push a `:main` tagged image

Every tagged release (v*.*.*) will:
- Run tests
- Build and push versioned images (`:v0.1.0`, `:v0.1`, `:v0`)
- Generate deployment manifests
- Create a GitHub release

## Container Image Tags

Each release creates multiple tags:

- `ghcr.io/bpierce/unifi-dns-operator:v0.1.0` - Exact version
- `ghcr.io/bpierce/unifi-dns-operator:v0.1` - Minor version
- `ghcr.io/bpierce/unifi-dns-operator:v0` - Major version
- `ghcr.io/bpierce/unifi-dns-operator:latest` - Latest stable release
- `ghcr.io/bpierce/unifi-dns-operator:main` - Latest commit to main branch

## Deployment Manifests

Each release includes a deployment manifest file:
- `unifi-dns-operator-v0.1.0.yaml` - Complete installation manifest

Users can install directly from GitHub releases:
```bash
kubectl apply -f https://github.com/bpierce/unifi-dns-operator/releases/download/v0.1.0/unifi-dns-operator-v0.1.0.yaml
```

## Making Package Public

By default, GitHub packages are private. To make your operator publicly accessible:

1. Go to https://github.com/bpierce?tab=packages
2. Click on `unifi-dns-operator`
3. Click "Package settings"
4. Scroll to "Danger Zone"
5. Click "Change visibility"
6. Select "Public"
7. Confirm by typing the package name

## Troubleshooting

### Build Fails

Check the GitHub Actions logs:
1. Go to Actions tab
2. Click on the failed workflow run
3. Expand the failed step to see error details

### Image Not Found

Ensure:
- The package is set to public
- The tag was created correctly
- The workflow completed successfully

### Permission Errors

Verify:
- Repository settings → Actions → General → Workflow permissions is set to "Read and write"
- Your GitHub token has package write permissions

## Best Practices

1. **Update CHANGELOG.md** before each release
2. **Follow Semantic Versioning:**
   - MAJOR: Breaking changes
   - MINOR: New features (backward compatible)
   - PATCH: Bug fixes
3. **Test locally** before creating a release
4. **Tag releases** from the main branch only
5. **Document breaking changes** clearly in release notes
