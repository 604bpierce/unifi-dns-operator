/*
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
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DNSRecordType represents the type of DNS record
// +kubebuilder:validation:Enum=A_RECORD;AAAA_RECORD;CNAME_RECORD;MX_RECORD;TXT_RECORD;SRV_RECORD
type DNSRecordType string

const (
	ARecord     DNSRecordType = "A_RECORD"
	AAAARecord  DNSRecordType = "AAAA_RECORD"
	CNAMERecord DNSRecordType = "CNAME_RECORD"
	MXRecord    DNSRecordType = "MX_RECORD"
	TXTRecord   DNSRecordType = "TXT_RECORD"
	SRVRecord   DNSRecordType = "SRV_RECORD"
)

// RecordData contains type-specific DNS record data
type RecordData struct {
	// IPv4Address is the IPv4 address for A_RECORD
	// +optional
	IPv4Address *string `json:"ipv4Address,omitempty"`

	// IPv6Address is the IPv6 address for AAAA_RECORD
	// +optional
	IPv6Address *string `json:"ipv6Address,omitempty"`

	// TargetDomain is the target domain for CNAME_RECORD
	// +optional
	TargetDomain *string `json:"targetDomain,omitempty"`

	// MailServerDomain is the mail server domain for MX_RECORD
	// +optional
	MailServerDomain *string `json:"mailServerDomain,omitempty"`

	// Priority is the priority for MX_RECORD and SRV_RECORD
	// +optional
	Priority *int `json:"priority,omitempty"`

	// TextValue is the text value for TXT_RECORD
	// +optional
	TextValue *string `json:"textValue,omitempty"`

	// Target is the target server for SRV_RECORD
	// +optional
	Target *string `json:"target,omitempty"`

	// Port is the port number for SRV_RECORD
	// +optional
	Port *int `json:"port,omitempty"`

	// Weight is the weight for SRV_RECORD
	// +optional
	Weight *int `json:"weight,omitempty"`
}

// UnifiDNSEntrySpec defines the desired state of UnifiDNSEntry
type UnifiDNSEntrySpec struct {
	// Type specifies the DNS record type
	// +kubebuilder:validation:Required
	Type DNSRecordType `json:"type"`

	// Domain is the domain name for the DNS entry
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Domain string `json:"domain"`

	// Enabled determines if the DNS policy is active
	// +kubebuilder:default=true
	Enabled bool `json:"enabled"`

	// RecordData contains type-specific DNS record information
	// +kubebuilder:validation:Required
	RecordData RecordData `json:"recordData"`

	// TTLSeconds is the Time-To-Live for the DNS record in seconds
	// +optional
	TTLSeconds *int `json:"ttlSeconds,omitempty"`
}

// UnifiDNSEntryStatus defines the observed state of UnifiDNSEntry.
type UnifiDNSEntryStatus struct {
	// PolicyID is the UUID of the DNS policy in Unifi
	// +optional
	PolicyID string `json:"policyId,omitempty"`

	// Synced indicates whether the DNS entry is synchronized with Unifi
	// +optional
	Synced bool `json:"synced,omitempty"`

	// LastSyncTime is the timestamp of the last successful sync
	// +optional
	LastSyncTime *metav1.Time `json:"lastSyncTime,omitempty"`

	// conditions represent the current state of the UnifiDNSEntry resource.
	// Each condition has a unique type and reflects the status of a specific aspect of the resource.
	//
	// Standard condition types include:
	// - "Available": the resource is fully functional
	// - "Progressing": the resource is being created or updated
	// - "Degraded": the resource failed to reach or maintain its desired state
	//
	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// UnifiDNSEntry is the Schema for the unifidnsentries API
type UnifiDNSEntry struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of UnifiDNSEntry
	// +required
	Spec UnifiDNSEntrySpec `json:"spec"`

	// status defines the observed state of UnifiDNSEntry
	// +optional
	Status UnifiDNSEntryStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// UnifiDNSEntryList contains a list of UnifiDNSEntry
type UnifiDNSEntryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []UnifiDNSEntry `json:"items"`
}

func init() {
	SchemeBuilder.Register(&UnifiDNSEntry{}, &UnifiDNSEntryList{})
}
