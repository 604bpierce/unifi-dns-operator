package unifi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is a client for the Unifi Cloud Controller API
type Client struct {
	baseURL    string
	siteID     string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Unifi API client
func NewClient(baseURL, siteID, apiKey string) *Client {
	// Skip TLS verification for self-signed certificates (common with local Unifi controllers)
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	return &Client{
		baseURL: baseURL,
		siteID:  siteID,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout:   30 * time.Second,
			Transport: transport,
		},
	}
}

// CreateDNSPolicy creates a new DNS policy in Unifi
func (c *Client) CreateDNSPolicy(policy DNSPolicy) (string, error) {
	url := fmt.Sprintf("%s/v1/sites/%s/dns/policies", c.baseURL, c.siteID)

	body, err := json.Marshal(policy)
	if err != nil {
		return "", fmt.Errorf("failed to marshal policy: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var createdPolicy DNSPolicy
	if err := json.NewDecoder(resp.Body).Decode(&createdPolicy); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return createdPolicy.ID, nil
}

// GetDNSPolicy retrieves a DNS policy from Unifi
func (c *Client) GetDNSPolicy(policyID string) (*DNSPolicy, error) {
	url := fmt.Sprintf("%s/v1/sites/%s/dns/policies/%s", c.baseURL, c.siteID, policyID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-API-KEY", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var policy DNSPolicy
	if err := json.NewDecoder(resp.Body).Decode(&policy); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &policy, nil
}

// UpdateDNSPolicy updates an existing DNS policy in Unifi
func (c *Client) UpdateDNSPolicy(policyID string, policy DNSPolicy) error {
	url := fmt.Sprintf("%s/v1/sites/%s/dns/policies/%s", c.baseURL, c.siteID, policyID)

	body, err := json.Marshal(policy)
	if err != nil {
		return fmt.Errorf("failed to marshal policy: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// DeleteDNSPolicy deletes a DNS policy from Unifi
func (c *Client) DeleteDNSPolicy(policyID string) error {
	url := fmt.Sprintf("%s/v1/sites/%s/dns/policies/%s", c.baseURL, c.siteID, policyID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-API-KEY", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
