package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"scandiff/models"
	"scandiff/services"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type SnapshotHandler struct {
	service *services.SnapshotService
}

func NewSnapshotHandler(service *services.SnapshotService) *SnapshotHandler {
	return &SnapshotHandler{service: service}
}

// UploadSnapshot handles snapshot file upload
func (h *SnapshotHandler) UploadSnapshot(c *gin.Context) {
	// Get uploaded file
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// Validate file extension
	if !strings.HasSuffix(header.Filename, ".json") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only JSON files are allowed"})
		return
	}

	// Read file contents
	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	// Parse JSON
	var snapshot models.Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid JSON format: %v", err)})
		return
	}

	// Validate required fields
	if snapshot.IP == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IP address is required"})
		return
	}
	if snapshot.Timestamp == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Timestamp is required"})
		return
	}

	// Store snapshot
	if err := h.service.CreateSnapshot(&snapshot, header.Filename); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			c.JSON(http.StatusConflict, gin.H{"error": "Snapshot with this IP and timestamp already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to store snapshot: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Snapshot uploaded successfully",
		"ip":        snapshot.IP,
		"timestamp": snapshot.Timestamp,
	})
}

// GetSnapshot retrieves a specific snapshot by ID
func (h *SnapshotHandler) GetSnapshot(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid snapshot ID"})
		return
	}

	snapshot, err := h.service.GetSnapshotByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Snapshot not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, snapshot)
}

// ListHosts retrieves all unique host IPs
func (h *SnapshotHandler) ListHosts(c *gin.Context) {
	hosts, err := h.service.ListAllHosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"hosts": hosts})
}

// ListSnapshotsByIP retrieves all snapshots for a given IP
func (h *SnapshotHandler) ListSnapshotsByIP(c *gin.Context) {
	ip := c.Param("ip")
	if ip == "" {
		fmt.Printf("[x] ip address is required\n")
		c.JSON(http.StatusBadRequest, gin.H{"error": "IP address is required"})
		return
	}

	if !regexp.MustCompile(`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`).MatchString(ip) && !regexp.MustCompile(`^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`).MatchString(ip) {
		fmt.Printf("[x] invalid ip address: %s\n", ip)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IP address"})
		return
	}

	snapshots, err := h.service.ListSnapshotsByIP(ip)
	if err != nil {
		fmt.Printf("[x] failed to list snapshots by ip: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"snapshots": snapshots})
}
