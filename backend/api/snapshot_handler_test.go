package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"scandiff/db"
	"scandiff/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *sql.DB {
	testDB, err := db.InitDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}
	return testDB
}

func TestUploadSnapshot_Success(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()

	router := SetupRouter(testDB)

	// Create a test snapshot file
	snapshotData := models.Snapshot{
		IP:           "192.168.1.1",
		Timestamp:    "2025-09-10T03:00:00Z",
		Services:     []models.Service{},
		ServiceCount: 0,
	}

	jsonData, _ := json.Marshal(snapshotData)

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test_snapshot.json")
	io.Copy(part, bytes.NewReader(jsonData))
	writer.Close()

	// Create request
	req, _ := http.NewRequest("POST", "/api/snapshots/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Execute request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Snapshot
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "192.168.1.1", response.IP)
	assert.Equal(t, "2025-09-10T03:00:00Z", response.Timestamp)
}

func TestUploadSnapshot_InvalidJSON(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()

	router := SetupRouter(testDB)

	// Create an invalid JSON file
	invalidData := []byte("{invalid json")

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test_snapshot.json")
	io.Copy(part, bytes.NewReader(invalidData))
	writer.Close()

	// Create request
	req, _ := http.NewRequest("POST", "/api/snapshots/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Execute request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUploadSnapshot_MissingRequiredField(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()

	router := SetupRouter(testDB)

	// Create snapshot without required IP field
	snapshotData := models.Snapshot{
		Timestamp:    "2025-09-10T03:00:00Z",
		Services:     []models.Service{},
		ServiceCount: 0,
	}

	jsonData, _ := json.Marshal(snapshotData)

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test_snapshot.json")
	io.Copy(part, bytes.NewReader(jsonData))
	writer.Close()

	// Create request
	req, _ := http.NewRequest("POST", "/api/snapshots/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Execute request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListHosts_Empty(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()

	router := SetupRouter(testDB)

	// Create request
	req, _ := http.NewRequest("GET", "/api/hosts", nil)

	// Execute request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response []*models.Snapshot
	json.Unmarshal(w.Body.Bytes(), &response)
	hosts := response
	assert.Nil(t, hosts)
}

func TestHealthCheck(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()

	router := SetupRouter(testDB)

	// Create request
	req, _ := http.NewRequest("GET", "/health", nil)

	// Execute request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	type responseType struct {
		Status string `json:"status"`
	}
	var response responseType
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "ok", response.Status)
}

func TestCompareDiff_OrderIndependent(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()

	router := SetupRouter(testDB)

	// Upload first snapshot (earlier timestamp)
	snapshot1Data := models.Snapshot{
		IP:           "192.168.1.1",
		Timestamp:    "2025-09-10T03:00:00Z",
		Services:     []models.Service{},
		ServiceCount: 0,
	}
	jsonData1, _ := json.Marshal(snapshot1Data)

	body1 := &bytes.Buffer{}
	writer1 := multipart.NewWriter(body1)
	part1, _ := writer1.CreateFormFile("file", "test1.json")
	io.Copy(part1, bytes.NewReader(jsonData1))
	writer1.Close()

	req1, _ := http.NewRequest("POST", "/api/snapshots/upload", body1)
	req1.Header.Set("Content-Type", writer1.FormDataContentType())
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	// Upload second snapshot (later timestamp)
	snapshot2Data := models.Snapshot{
		IP:           "192.168.1.1",
		Timestamp:    "2025-09-15T03:00:00Z",
		Services:     []models.Service{},
		ServiceCount: 0,
	}
	jsonData2, _ := json.Marshal(snapshot2Data)

	body2 := &bytes.Buffer{}
	writer2 := multipart.NewWriter(body2)
	part2, _ := writer2.CreateFormFile("file", "test2.json")
	io.Copy(part2, bytes.NewReader(jsonData2))
	writer2.Close()

	req2, _ := http.NewRequest("POST", "/api/snapshots/upload", body2)
	req2.Header.Set("Content-Type", writer2.FormDataContentType())
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusCreated, w2.Code)

	// Test diff in correct order (1, 2)
	reqDiff1, _ := http.NewRequest("GET", "/api/diff/1/2", nil)
	wDiff1 := httptest.NewRecorder()
	router.ServeHTTP(wDiff1, reqDiff1)
	assert.Equal(t, http.StatusOK, wDiff1.Code)

	var diff1 models.DiffReport
	json.Unmarshal(wDiff1.Body.Bytes(), &diff1)

	// Test diff in reverse order (2, 1)
	reqDiff2, _ := http.NewRequest("GET", "/api/diff/2/1", nil)
	wDiff2 := httptest.NewRecorder()
	router.ServeHTTP(wDiff2, reqDiff2)
	assert.Equal(t, http.StatusOK, wDiff2.Code)

	var diff2 models.DiffReport
	json.Unmarshal(wDiff2.Body.Bytes(), &diff2)

	// Both should have the same old and new snapshot IDs regardless of request order
	oldSnap1 := diff1.OldSnapshot
	newSnap1 := diff1.NewSnapshot
	oldSnap2 := diff2.OldSnapshot
	newSnap2 := diff2.NewSnapshot

	assert.Equal(t, oldSnap1.Timestamp, oldSnap2.Timestamp)
	assert.Equal(t, newSnap1.Timestamp, newSnap2.Timestamp)
	assert.Equal(t, "2025-09-10T03:00:00Z", oldSnap1.Timestamp)
	assert.Equal(t, "2025-09-15T03:00:00Z", newSnap1.Timestamp)
}

// test ipv6 address
func TestListSnapshotsByIP_IPv6(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()

	router := SetupRouter(testDB)

	// create a test snapshot
	snapshot := models.Snapshot{
		IP:           "2001:db8::1",
		Timestamp:    "2025-09-10T03:00:00Z",
		Services:     []models.Service{},
		ServiceCount: 0,
	}

	jsonData, _ := json.Marshal(snapshot)

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test_snapshot.json")
	io.Copy(part, bytes.NewReader(jsonData))
	writer.Close()

	// Create request
	req, _ := http.NewRequest("POST", "/api/snapshots/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Execute request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	req, _ = http.NewRequest("GET", "/api/hosts/2001:db8::1/snapshots", nil)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	type responseType struct {
		Snapshots []*models.Snapshot `json:"snapshots"`
	}
	var response responseType
	json.Unmarshal(w.Body.Bytes(), &response)
	fmt.Printf("response: %#v\n", response)
	assert.Equal(t, "2001:db8::1", response.Snapshots[0].IP)
}
