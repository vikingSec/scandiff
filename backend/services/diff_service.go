package services

import (
	"reflect"
	"scandiff/models"
)

type DiffService struct{}

func NewDiffService() *DiffService {
	return &DiffService{}
}

// CompareSnapshots generates a diff report between two snapshots
func (d *DiffService) CompareSnapshots(oldSnapshot, newSnapshot *models.Snapshot) *models.DiffReport {
	report := &models.DiffReport{
		OldSnapshot:     oldSnapshot,
		NewSnapshot:     newSnapshot,
		PortsAdded:      []models.PortChange{},
		PortsRemoved:    []models.PortChange{},
		ServicesChanged: []models.ServiceChange{},
	}

	// Create maps for quick lookup
	oldServices := make(map[int]*models.Service)
	newServices := make(map[int]*models.Service)

	for i := range oldSnapshot.Services {
		oldServices[oldSnapshot.Services[i].Port] = &oldSnapshot.Services[i]
	}

	for i := range newSnapshot.Services {
		newServices[newSnapshot.Services[i].Port] = &newSnapshot.Services[i]
	}

	// Find removed ports
	for port, service := range oldServices {
		if _, exists := newServices[port]; !exists {
			report.PortsRemoved = append(report.PortsRemoved, models.PortChange{
				Port:       port,
				ChangeType: models.ChangeTypeRemoved,
				OldService: service,
			})
		}
	}

	// Find added ports and modified services
	for port, newService := range newServices {
		if oldService, exists := oldServices[port]; !exists {
			report.PortsAdded = append(report.PortsAdded, models.PortChange{
				Port:       port,
				ChangeType: models.ChangeTypeAdded,
				NewService: newService,
			})
		} else {
			// Check for changes in the service
			if change := d.compareServices(oldService, newService); change != nil {
				report.ServicesChanged = append(report.ServicesChanged, *change)
			}
		}
	}

	report.HasChanges = len(report.PortsAdded) > 0 ||
		len(report.PortsRemoved) > 0 ||
		len(report.ServicesChanged) > 0

	return report
}

// compareServices compares two services and returns changes if any
func (d *DiffService) compareServices(oldService, newService *models.Service) *models.ServiceChange {
	change := &models.ServiceChange{
		Port:     oldService.Port,
		Protocol: oldService.Protocol,
	}

	hasChanges := false

	// Check status change
	if oldService.Status != newService.Status {
		change.StatusChanged = true
		change.OldStatus = oldService.Status
		change.NewStatus = newService.Status
		hasChanges = true
	}
	if oldService.Protocol != newService.Protocol {
		change.ProtocolChanged = true
		change.OldProtocol = oldService.Protocol
		change.NewProtocol = newService.Protocol
		hasChanges = true
	}

	// Check software change
	if !reflect.DeepEqual(oldService.Software, newService.Software) {
		change.SoftwareChanged = true
		change.OldSoftware = &oldService.Software
		change.NewSoftware = &newService.Software
		hasChanges = true
	}

	// Check TLS change
	if !d.compareTLS(oldService.TLS, newService.TLS) {
		change.TLSChanged = true
		change.OldTLS = oldService.TLS
		change.NewTLS = newService.TLS
		hasChanges = true
	}

	// Check vulnerability changes
	vulnsAdded, vulnsFixed := d.compareVulnerabilities(
		oldService.Vulnerabilities,
		newService.Vulnerabilities,
	)

	if len(vulnsAdded) > 0 {
		change.VulnerabilitiesAdded = vulnsAdded
		hasChanges = true
	}

	if len(vulnsFixed) > 0 {
		change.VulnerabilitiesFixed = vulnsFixed
		hasChanges = true
	}

	if !hasChanges {
		return nil
	}

	return change
}

// compareTLS compares two TLS configurations
func (d *DiffService) compareTLS(old, new *models.TLS) bool {
	if old == nil && new == nil {
		return true
	}
	if old == nil || new == nil {
		return false
	}
	return reflect.DeepEqual(old, new)
}

// compareVulnerabilities finds added and fixed vulnerabilities
func (d *DiffService) compareVulnerabilities(old, new []string) (added, fixed []string) {
	oldMap := make(map[string]bool)
	newMap := make(map[string]bool)

	for _, v := range old {
		oldMap[v] = true
	}

	for _, v := range new {
		newMap[v] = true
	}

	// Find added vulnerabilities
	for _, v := range new {
		if !oldMap[v] {
			added = append(added, v)
		}
	}

	// Find fixed vulnerabilities
	for _, v := range old {
		if !newMap[v] {
			fixed = append(fixed, v)
		}
	}

	return added, fixed
}
