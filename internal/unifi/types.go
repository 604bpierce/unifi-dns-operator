package unifi

// DNSPolicyType represents the type of DNS policy in Unifi API
type DNSPolicyType string

const (
	ARecordType     DNSPolicyType = "A_RECORD"
	AAAARecordType  DNSPolicyType = "AAAA_RECORD"
	CNAMERecordType DNSPolicyType = "CNAME_RECORD"
	MXRecordType    DNSPolicyType = "MX_RECORD"
	TXTRecordType   DNSPolicyType = "TXT_RECORD"
	SRVRecordType   DNSPolicyType = "SRV_RECORD"
)

// DNSPolicyMetadata contains metadata about the DNS policy
type DNSPolicyMetadata struct {
	Origin string `json:"origin"`
}

// DNSPolicy represents a DNS policy in the Unifi API
type DNSPolicy struct {
	ID       string             `json:"id,omitempty"`
	Type     DNSPolicyType      `json:"type"`
	Enabled  bool               `json:"enabled"`
	Domain   string             `json:"domain,omitempty"`
	Metadata *DNSPolicyMetadata `json:"metadata,omitempty"`

	// A_RECORD fields
	IPv4Address *string `json:"ipv4Address,omitempty"`

	// AAAA_RECORD fields
	IPv6Address *string `json:"ipv6Address,omitempty"`

	// CNAME_RECORD fields
	TargetDomain *string `json:"targetDomain,omitempty"`

	// MX_RECORD fields
	MailServerDomain *string `json:"mailServerDomain,omitempty"`
	Priority         *int    `json:"priority,omitempty"`

	// TXT_RECORD fields
	TextValue *string `json:"textValue,omitempty"`

	// SRV_RECORD fields
	Target *string `json:"target,omitempty"`
	Port   *int    `json:"port,omitempty"`
	Weight *int    `json:"weight,omitempty"`

	// Common optional field
	TTLSeconds *int `json:"ttlSeconds,omitempty"`
}

// ErrorResponse represents an error response from the Unifi API
type ErrorResponse struct {
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}
