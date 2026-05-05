# Unifi DNS Operator

[![CI](https://github.com/bpierce/unifi-dns-operator/actions/workflows/test.yaml/badge.svg)](https://github.com/bpierce/unifi-dns-operator/actions/workflows/test.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/bpierce/unifi-dns-operator)](https://goreportcard.com/report/github.com/bpierce/unifi-dns-operator)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

A Kubernetes operator for managing DNS entries on a Unifi Network Controller. This operator allows you to declaratively manage DNS policies in your Unifi Network using Kubernetes Custom Resources.

## Overview

The Unifi DNS Operator provides a `UnifiDNSEntry` Custom Resource Definition (CRD) that enables you to manage DNS entries in your Unifi Network Controller directly from Kubernetes. This is particularly useful for managing DNS records for services running in your Kubernetes cluster.

### Supported DNS Record Types

- **A_RECORD** - IPv4 address records
- **AAAA_RECORD** - IPv6 address records
- **CNAME_RECORD** - Canonical name (alias) records
- **MX_RECORD** - Mail exchange records
- **TXT_RECORD** - Text records
- **SRV_RECORD** - Service records

## Prerequisites

- Kubernetes cluster (v1.35+)
- kubectl configured to access your cluster
- Unifi Network Controller with API access
- Unifi API key and Site ID

## Quick Start

### Install from GitHub Releases (Recommended)

1. **Create the credentials secret:**

```bash
kubectl create secret generic unifi-dns-operator-credentials \
  --from-literal=UNIFI_API_URL=https://your-controller/proxy/network/integration \
  --from-literal=UNIFI_SITE_ID=your-site-id \
  --from-literal=UNIFI_API_KEY=your-api-key \
  -n unifi-dns-operator-system
```

> **Note for self-hosted controllers:** Use `https://YOUR_CONTROLLER_IP/proxy/network/integration` as the API URL.
> **Note for cloud controllers:** Use `https://api.ui.com` as the API URL.

2. **Install the operator:**

```bash
kubectl apply -f https://github.com/bpierce/unifi-dns-operator/releases/latest/download/unifi-dns-operator-latest.yaml
```

Or install a specific version:

```bash
VERSION=v0.1.0
kubectl apply -f https://github.com/bpierce/unifi-dns-operator/releases/download/${VERSION}/unifi-dns-operator-${VERSION}.yaml
```

3. **Verify installation:**

```bash
kubectl get pods -n unifi-dns-operator-system
```

### Alternative: Install from Source

If you want to build and deploy from source:

1. **Clone the repository:**

```bash
git clone https://github.com/bpierce/unifi-dns-operator.git
cd unifi-dns-operator
```

2. **Install the CRD:**

```bash
make install
```

3. **Create credentials secret** (same as above)

4. **Deploy the operator:**

```bash
make deploy IMG=ghcr.io/bpierce/unifi-dns-operator:latest
```

## Usage

### Creating DNS Entries

Create a `UnifiDNSEntry` resource to manage a DNS entry:

#### A Record Example

```yaml
apiVersion: dns.bx.network/v1alpha1
kind: UnifiDNSEntry
metadata:
  name: my-service-dns
spec:
  type: A_RECORD
  domain: myservice.local
  enabled: true
  recordData:
    ipv4Address: "192.168.1.100"
  ttlSeconds: 300
```

#### CNAME Record Example

```yaml
apiVersion: dns.bx.network/v1alpha1
kind: UnifiDNSEntry
metadata:
  name: app-alias
spec:
  type: CNAME_RECORD
  domain: app.example.local
  enabled: true
  recordData:
    targetDomain: "example.local"
  ttlSeconds: 300
```

#### AAAA Record Example

```yaml
apiVersion: dns.bx.network/v1alpha1
kind: UnifiDNSEntry
metadata:
  name: ipv6-service
spec:
  type: AAAA_RECORD
  domain: ipv6.example.local
  enabled: true
  recordData:
    ipv6Address: "2001:db8::1"
  ttlSeconds: 300
```

### Applying DNS Entries

```bash
kubectl apply -f config/samples/dns_v1alpha1_unifidnsentry.yaml
```

### Checking Status

```bash
kubectl get unifidnsentries
kubectl describe unifidnsentry unifidnsentry-sample
```

### Deleting DNS Entries

When you delete a `UnifiDNSEntry` resource, the operator will automatically remove the corresponding DNS policy from Unifi:

```bash
kubectl delete unifidnsentry unifidnsentry-sample
```

## Development

### Running Locally

1. Install the CRDs:
```bash
make install
```

2. Set environment variables:
```bash
export UNIFI_API_URL=https://api.ui.com
export UNIFI_SITE_ID=your-site-id
export UNIFI_API_KEY=your-api-key
```

3. Run the operator:
```bash
make run
```

### Building

Build the operator binary:
```bash
make build
```

Build and push the Docker image:
```bash
make docker-build docker-push IMG=your-registry/unifi-dns-operator:tag
```

### Testing

Run tests:
```bash
make test
```

## API Reference

### UnifiDNSEntry Spec

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `type` | string | Yes | DNS record type (A_RECORD, AAAA_RECORD, CNAME_RECORD, MX_RECORD, TXT_RECORD, SRV_RECORD) |
| `domain` | string | Yes | The domain name for the DNS entry |
| `enabled` | bool | Yes | Whether the DNS policy is active (default: true) |
| `recordData` | RecordData | Yes | Type-specific DNS record data |
| `ttlSeconds` | int | No | Time-To-Live for the DNS record in seconds |

### RecordData Fields

| Field | Type | Record Types | Description |
|-------|------|--------------|-------------|
| `ipv4Address` | string | A_RECORD | IPv4 address |
| `ipv6Address` | string | AAAA_RECORD | IPv6 address |
| `targetDomain` | string | CNAME_RECORD | Target domain for alias |
| `mailServerDomain` | string | MX_RECORD | Mail server domain |
| `priority` | int | MX_RECORD, SRV_RECORD | Priority value |
| `textValue` | string | TXT_RECORD | Text value |
| `target` | string | SRV_RECORD | Target server |
| `port` | int | SRV_RECORD | Port number |
| `weight` | int | SRV_RECORD | Weight for load balancing |

### UnifiDNSEntry Status

| Field | Type | Description |
|-------|------|-------------|
| `policyId` | string | UUID of the DNS policy in Unifi |
| `synced` | bool | Whether the entry is synchronized with Unifi |
| `lastSyncTime` | timestamp | Time of last successful sync |
| `conditions` | []Condition | Standard Kubernetes conditions |

## Architecture

The operator follows the standard Kubernetes operator pattern:

1. **Custom Resource Definition (CRD)**: Defines the `UnifiDNSEntry` schema
2. **Controller**: Watches for `UnifiDNSEntry` resources and reconciles them with Unifi
3. **Unifi API Client**: Handles communication with the Unifi Cloud Controller API
4. **Finalizers**: Ensures DNS policies are cleaned up when resources are deleted

### Reconciliation Flow

1. Fetch the `UnifiDNSEntry` resource
2. Check if the resource is being deleted (handle finalizer)
3. If policy ID exists in status, update the existing Unifi policy
4. If no policy ID, create a new Unifi policy
5. Update the resource status with sync information and conditions

## Contributing

Contributions are welcome! This is a basic operator that can be extended with additional features such as:

- Webhook validation for CRD fields
- Multi-site support
- Enhanced metrics and observability
- Bulk operations
- Status reporting improvements

## License

Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
