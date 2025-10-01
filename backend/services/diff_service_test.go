package services

import (
	"fmt"
	"scandiff/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompareSnapshots_NoChanges(t *testing.T) {
	diffService := NewDiffService()

	snapshot1 := &models.Snapshot{
		IP:        "192.168.1.1",
		Timestamp: "2025-09-10T03:00:00Z",
		Services: []models.Service{
			{
				Port:     80,
				Protocol: "HTTP",
				Status:   200,
				Software: models.Software{
					Vendor:  "nginx",
					Product: "nginx",
					Version: "1.22.1",
				},
			},
		},
	}

	snapshot2 := &models.Snapshot{
		IP:        "192.168.1.1",
		Timestamp: "2025-09-15T03:00:00Z",
		Services: []models.Service{
			{
				Port:     80,
				Protocol: "HTTP",
				Status:   200,
				Software: models.Software{
					Vendor:  "nginx",
					Product: "nginx",
					Version: "1.22.1",
				},
			},
		},
	}

	diff := diffService.CompareSnapshots(snapshot1, snapshot2)

	assert.False(t, diff.HasChanges)
	assert.Empty(t, diff.PortsAdded)
	assert.Empty(t, diff.PortsRemoved)
	assert.Empty(t, diff.ServicesChanged)
}

func TestCompareSnapshots_PortAdded(t *testing.T) {
	diffService := NewDiffService()

	snapshot1 := &models.Snapshot{
		IP:        "192.168.1.1",
		Timestamp: "2025-09-10T03:00:00Z",
		Services: []models.Service{
			{
				Port:     80,
				Protocol: "HTTP",
			},
		},
	}

	snapshot2 := &models.Snapshot{
		IP:        "192.168.1.1",
		Timestamp: "2025-09-15T03:00:00Z",
		Services: []models.Service{
			{
				Port:     80,
				Protocol: "HTTP",
			},
			{
				Port:     443,
				Protocol: "HTTPS",
			},
		},
	}

	diff := diffService.CompareSnapshots(snapshot1, snapshot2)

	assert.True(t, diff.HasChanges)
	assert.Len(t, diff.PortsAdded, 1)
	assert.Equal(t, 443, diff.PortsAdded[0].Port)
	assert.Equal(t, models.ChangeTypeAdded, diff.PortsAdded[0].ChangeType)
}

func TestCompareSnapshots_PortRemoved(t *testing.T) {
	diffService := NewDiffService()

	snapshot1 := &models.Snapshot{
		IP:        "192.168.1.1",
		Timestamp: "2025-09-10T03:00:00Z",
		Services: []models.Service{
			{
				Port:     80,
				Protocol: "HTTP",
			},
			{
				Port:     443,
				Protocol: "HTTPS",
			},
		},
	}

	snapshot2 := &models.Snapshot{
		IP:        "192.168.1.1",
		Timestamp: "2025-09-15T03:00:00Z",
		Services: []models.Service{
			{
				Port:     80,
				Protocol: "HTTP",
			},
		},
	}

	diff := diffService.CompareSnapshots(snapshot1, snapshot2)

	assert.True(t, diff.HasChanges)
	assert.Len(t, diff.PortsRemoved, 1)
	assert.Equal(t, 443, diff.PortsRemoved[0].Port)
	assert.Equal(t, models.ChangeTypeRemoved, diff.PortsRemoved[0].ChangeType)
}

func TestCompareSnapshots_SoftwareVersionChanged(t *testing.T) {
	diffService := NewDiffService()

	snapshot1 := &models.Snapshot{
		IP:        "192.168.1.1",
		Timestamp: "2025-09-10T03:00:00Z",
		Services: []models.Service{
			{
				Port:     80,
				Protocol: "HTTP",
				Software: models.Software{
					Vendor:  "nginx",
					Product: "nginx",
					Version: "1.22.1",
				},
			},
		},
	}

	snapshot2 := &models.Snapshot{
		IP:        "192.168.1.1",
		Timestamp: "2025-09-15T03:00:00Z",
		Services: []models.Service{
			{
				Port:     80,
				Protocol: "HTTP",
				Software: models.Software{
					Vendor:  "nginx",
					Product: "nginx",
					Version: "1.23.0",
				},
			},
		},
	}

	diff := diffService.CompareSnapshots(snapshot1, snapshot2)

	assert.True(t, diff.HasChanges)
	assert.Len(t, diff.ServicesChanged, 1)
	assert.True(t, diff.ServicesChanged[0].SoftwareChanged)
	assert.Equal(t, "1.22.1", diff.ServicesChanged[0].OldSoftware.Version)
	assert.Equal(t, "1.23.0", diff.ServicesChanged[0].NewSoftware.Version)
}

func TestCompareSnapshots_VulnerabilitiesAdded(t *testing.T) {
	diffService := NewDiffService()

	snapshot1 := &models.Snapshot{
		IP:        "192.168.1.1",
		Timestamp: "2025-09-10T03:00:00Z",
		Services: []models.Service{
			{
				Port:            22,
				Protocol:        "SSH",
				Vulnerabilities: []string{},
			},
		},
	}

	snapshot2 := &models.Snapshot{
		IP:        "192.168.1.1",
		Timestamp: "2025-09-15T03:00:00Z",
		Services: []models.Service{
			{
				Port:            22,
				Protocol:        "SSH",
				Vulnerabilities: []string{"CVE-2023-12345"},
			},
		},
	}

	diff := diffService.CompareSnapshots(snapshot1, snapshot2)

	assert.True(t, diff.HasChanges)
	assert.Len(t, diff.ServicesChanged, 1)
	assert.Contains(t, diff.ServicesChanged[0].VulnerabilitiesAdded, "CVE-2023-12345")
}

func TestCompareSnapshots_VulnerabilitiesFixed(t *testing.T) {
	diffService := NewDiffService()

	snapshot1 := &models.Snapshot{
		IP:        "192.168.1.1",
		Timestamp: "2025-09-10T03:00:00Z",
		Services: []models.Service{
			{
				Port:            22,
				Protocol:        "SSH",
				Vulnerabilities: []string{"CVE-2023-12345"},
			},
		},
	}

	snapshot2 := &models.Snapshot{
		IP:        "192.168.1.1",
		Timestamp: "2025-09-15T03:00:00Z",
		Services: []models.Service{
			{
				Port:            22,
				Protocol:        "SSH",
				Vulnerabilities: []string{},
			},
		},
	}

	diff := diffService.CompareSnapshots(snapshot1, snapshot2)

	assert.True(t, diff.HasChanges)
	assert.Len(t, diff.ServicesChanged, 1)
	assert.Contains(t, diff.ServicesChanged[0].VulnerabilitiesFixed, "CVE-2023-12345")
}

func TestCompareSnapshots_TLSChanged(t *testing.T) {
	diffService := NewDiffService()

	snapshot1 := &models.Snapshot{
		IP:        "192.168.1.1",
		Timestamp: "2025-09-10T03:00:00Z",
		Services: []models.Service{
			{
				Port:     80,
				Protocol: "HTTP",
				TLS: &models.TLS{
					Version:               "TLS 1.2",
					Cipher:                "AES-256-GCM-SHA384",
					CertFingerprintSHA256: "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
				},
			},
		},
	}

	snapshot2 := &models.Snapshot{
		IP:        "192.168.1.1",
		Timestamp: "2025-09-15T03:00:00Z",
		Services: []models.Service{
			{
				Port:     80,
				Protocol: "HTTP",
				TLS: &models.TLS{
					Version:               "TLS 1.3",
					Cipher:                "AES-256-GCM-SHA384",
					CertFingerprintSHA256: "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
				},
			},
		},
	}

	diff := diffService.CompareSnapshots(snapshot1, snapshot2)

	assert.True(t, diff.HasChanges)
	assert.Len(t, diff.ServicesChanged, 1)
	assert.True(t, diff.ServicesChanged[0].TLSChanged)
	assert.Equal(t, "TLS 1.2", diff.ServicesChanged[0].OldTLS.Version)
	assert.Equal(t, "TLS 1.3", diff.ServicesChanged[0].NewTLS.Version)
}

// test ipv6 address
func TestCompareSnapshots_IPv6(t *testing.T) {
	diffService := NewDiffService()

	snapshot1 := &models.Snapshot{
		IP:        "2001:db8::1",
		Timestamp: "2025-09-10T03:00:00Z",
		Services: []models.Service{
			{
				Port:     80,
				Protocol: "TCP",
			},
		},
	}

	snapshot2 := &models.Snapshot{
		IP:        "2001:db8::1",
		Timestamp: "2025-09-15T03:00:00Z",
		Services: []models.Service{
			{
				Port:     80,
				Protocol: "HTTP",
			},
		},
	}

	diff := diffService.CompareSnapshots(snapshot1, snapshot2)
	fmt.Printf("diff: %#v\n", diff)
	assert.True(t, diff.HasChanges)
	assert.Len(t, diff.ServicesChanged, 1)
	assert.Equal(t, "TCP", diff.ServicesChanged[0].OldProtocol)
	assert.Equal(t, "HTTP", diff.ServicesChanged[0].NewProtocol)
}
