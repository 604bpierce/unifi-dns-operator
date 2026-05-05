# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of Unifi DNS Operator
- Support for all DNS record types (A, AAAA, CNAME, MX, TXT, SRV)
- Custom Resource Definition (CRD) for managing DNS entries
- Automatic synchronization with Unifi Network Controller
- Finalizers for proper cleanup on resource deletion
- Status conditions for monitoring sync state
- TLS skip verification for self-hosted controllers
- X-API-KEY authentication support for local Unifi controllers

### Features
- Declarative DNS management via Kubernetes manifests
- Full lifecycle management (create, update, delete)
- Multi-architecture support (amd64, arm64)
- Helm chart for easy deployment (coming soon)

## [0.1.0] - TBD

### Added
- Initial beta release
