package api

import (
	"fmt"
	"net/http"
	"scandiff/services"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type DiffHandler struct {
	snapshotService *services.SnapshotService
	diffService     *services.DiffService
}

func NewDiffHandler(snapshotService *services.SnapshotService, diffService *services.DiffService) *DiffHandler {
	return &DiffHandler{
		snapshotService: snapshotService,
		diffService:     diffService,
	}
}

// CompareSnapshots compares two snapshots and returns a diff report
func (h *DiffHandler) CompareSnapshots(c *gin.Context) {
	oldIDStr := c.Param("oldId")
	newIDStr := c.Param("newId")

	oldID, err := strconv.Atoi(oldIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid old snapshot ID"})
		return
	}

	newID, err := strconv.Atoi(newIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid new snapshot ID"})
		return
	}

	// Retrieve snapshots
	snapshot1, err := h.snapshotService.GetSnapshotByID(oldID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "First snapshot not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get first snapshot"})
		fmt.Printf("[x] failed to get first snapshot: %v", err)
		return
	}

	snapshot2, err := h.snapshotService.GetSnapshotByID(newID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Second snapshot not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get second snapshot"})
		fmt.Printf("[x] failed to get second snapshot: %v", err)
		return
	}

	// Validate that both snapshots are for the same host
	if snapshot1.IP != snapshot2.IP {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Snapshots must be from the same host"})
		return
	}

	// Determine which snapshot is older based on timestamp
	// This allows the API to work regardless of the order IDs are provided
	var oldSnapshot, newSnapshot = snapshot1, snapshot2
	if snapshot1.Timestamp > snapshot2.Timestamp {
		oldSnapshot = snapshot2
		newSnapshot = snapshot1
	}

	// Generate diff
	diff := h.diffService.CompareSnapshots(oldSnapshot, newSnapshot)

	c.JSON(http.StatusOK, diff)
}
