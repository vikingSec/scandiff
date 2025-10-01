package models

import (
	"encoding/json"
	"time"
)

// Software represents software information
type Software struct {
	Vendor  string `json:"vendor,omitempty"`
	Product string `json:"product,omitempty"`
	Version string `json:"version,omitempty"`
}

// TLS represents TLS configuration
type TLS struct {
	Version               string `json:"version,omitempty"`
	Cipher                string `json:"cipher,omitempty"`
	CertFingerprintSHA256 string `json:"cert_fingerprint_sha256,omitempty"`
}

// Service represents a network service
type Service struct {
	Port            int      `json:"port"`
	Protocol        string   `json:"protocol"`
	Status          int      `json:"status,omitempty"`
	Software        Software `json:"software,omitempty"`
	TLS             *TLS     `json:"tls,omitempty"`
	Vulnerabilities []string `json:"vulnerabilities,omitempty"`
}

// Snapshot represents a host snapshot
type Snapshot struct {
	ID           int       `json:"id,omitempty"`
	IP           string    `json:"ip"`
	Timestamp    string    `json:"timestamp"`
	Filename     string    `json:"filename,omitempty"`
	Services     []Service `json:"services"`
	ServiceCount int       `json:"service_count"`
	UploadedAt   time.Time `json:"uploaded_at,omitempty"`
}

// SnapshotRecord represents a database record
type SnapshotRecord struct {
	ID         int
	IP         string
	Timestamp  string
	Filename   string
	Data       string
	UploadedAt time.Time
}

// ToSnapshot converts a SnapshotRecord to a Snapshot
func (sr *SnapshotRecord) ToSnapshot() (*Snapshot, error) {
	var snapshot Snapshot
	if err := json.Unmarshal([]byte(sr.Data), &snapshot); err != nil {
		return nil, err
	}
	snapshot.ID = sr.ID
	snapshot.Filename = sr.Filename
	snapshot.UploadedAt = sr.UploadedAt
	return &snapshot, nil
}
