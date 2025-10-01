package models

// ChangeType represents the type of change
type ChangeType string

const (
	ChangeTypeAdded    ChangeType = "added"
	ChangeTypeRemoved  ChangeType = "removed"
	ChangeTypeModified ChangeType = "modified"
)

// PortChange represents a change in ports
type PortChange struct {
	Port       int        `json:"port"`
	ChangeType ChangeType `json:"change_type"`
	OldService *Service   `json:"old_service,omitempty"`
	NewService *Service   `json:"new_service,omitempty"`
}

// ServiceChange represents changes in a specific service
type ServiceChange struct {
	Port                 int       `json:"port"`
	Protocol             string    `json:"protocol"`
	StatusChanged        bool      `json:"status_changed"`
	OldStatus            int       `json:"old_status,omitempty"`
	NewStatus            int       `json:"new_status,omitempty"`
	SoftwareChanged      bool      `json:"software_changed"`
	OldSoftware          *Software `json:"old_software,omitempty"`
	NewSoftware          *Software `json:"new_software,omitempty"`
	TLSChanged           bool      `json:"tls_changed"`
	OldTLS               *TLS      `json:"old_tls,omitempty"`
	NewTLS               *TLS      `json:"new_tls,omitempty"`
	VulnerabilitiesAdded []string  `json:"vulnerabilities_added,omitempty"`
	VulnerabilitiesFixed []string  `json:"vulnerabilities_fixed,omitempty"`
	ProtocolChanged      bool      `json:"protocol_changed"`
	OldProtocol          string    `json:"old_protocol,omitempty"`
	NewProtocol          string    `json:"new_protocol,omitempty"`
}

// DiffReport represents the comparison between two snapshots
type DiffReport struct {
	OldSnapshot     *Snapshot       `json:"old_snapshot"`
	NewSnapshot     *Snapshot       `json:"new_snapshot"`
	PortsAdded      []PortChange    `json:"ports_added"`
	PortsRemoved    []PortChange    `json:"ports_removed"`
	ServicesChanged []ServiceChange `json:"services_changed"`
	HasChanges      bool            `json:"has_changes"`
}
