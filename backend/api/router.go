package api

import (
	"database/sql"
	"scandiff/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:5173"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(config))

	// Initialize services
	snapshotService := services.NewSnapshotService(db)
	diffService := services.NewDiffService()

	// Initialize handlers
	snapshotHandler := NewSnapshotHandler(snapshotService)
	diffHandler := NewDiffHandler(snapshotService, diffService)

	// API routes
	api := router.Group("/api")
	{
		// Snapshot routes
		api.POST("/snapshots/upload", snapshotHandler.UploadSnapshot)
		api.GET("/snapshots/:id", snapshotHandler.GetSnapshot)
		api.GET("/hosts", snapshotHandler.ListHosts)
		api.GET("/hosts/:ip/snapshots", snapshotHandler.ListSnapshotsByIP)

		// Diff routes
		api.GET("/diff/:oldId/:newId", diffHandler.CompareSnapshots)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}
