package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"scandiff/models"
)

type SnapshotService struct {
	db *sql.DB
}

func NewSnapshotService(db *sql.DB) *SnapshotService {
	return &SnapshotService{db: db}
}

// CreateSnapshot stores a new snapshot
func (s *SnapshotService) CreateSnapshot(snapshot *models.Snapshot, filename string) error {
	// Serialize snapshot data to JSON
	data, err := json.Marshal(snapshot)
	if err != nil {
		return fmt.Errorf("failed to marshal snapshot: %w", err)
	}

	query := `
		INSERT INTO snapshots (ip, timestamp, filename, data)
		VALUES (?, ?, ?, ?)
	`

	_, err = s.db.Exec(query, snapshot.IP, snapshot.Timestamp, filename, string(data))
	if err != nil {
		return fmt.Errorf("failed to insert snapshot: %w", err)
	}

	return nil
}

// GetSnapshotByID retrieves a snapshot by its ID
func (s *SnapshotService) GetSnapshotByID(id int) (*models.Snapshot, error) {
	query := `
		SELECT id, ip, timestamp, filename, data, uploaded_at
		FROM snapshots
		WHERE id = ?
		LIMIT 1
	`

	var record models.SnapshotRecord
	err := s.db.QueryRow(query, id).Scan(
		&record.ID,
		&record.IP,
		&record.Timestamp,
		&record.Filename,
		&record.Data,
		&record.UploadedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("snapshot not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query snapshot: %w", err)
	}

	return record.ToSnapshot()
}

// ListSnapshotsByIP retrieves all snapshots for a given IP
func (s *SnapshotService) ListSnapshotsByIP(ip string) ([]*models.Snapshot, error) {
	query := `
		SELECT id, ip, timestamp, filename, data, uploaded_at
		FROM snapshots
		WHERE ip = ?
		ORDER BY timestamp DESC
	`

	rows, err := s.db.Query(query, ip)
	if err != nil {
		return nil, fmt.Errorf("failed to query snapshots: %w", err)
	}
	defer rows.Close()

	var snapshots []*models.Snapshot
	for rows.Next() {
		var record models.SnapshotRecord
		err := rows.Scan(
			&record.ID,
			&record.IP,
			&record.Timestamp,
			&record.Filename,
			&record.Data,
			&record.UploadedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		snapshot, err := record.ToSnapshot()
		if err != nil {
			return nil, fmt.Errorf("failed to convert record: %w", err)
		}

		snapshots = append(snapshots, snapshot)
	}

	return snapshots, nil
}

// ListAllHosts retrieves all unique host IPs
func (s *SnapshotService) ListAllHosts() ([]string, error) {
	query := `
		SELECT DISTINCT ip
		FROM snapshots
		ORDER BY ip
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query hosts: %w", err)
	}
	defer rows.Close()

	var hosts []string
	for rows.Next() {
		var ip string
		if err := rows.Scan(&ip); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		hosts = append(hosts, ip)
	}

	return hosts, nil
}
