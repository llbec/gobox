package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/llbec/gobox/apps/lepu/config"
	"github.com/llbec/gobox/apps/lepu/http"
	"github.com/llbec/gobox/apps/lepu/logger"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Initialize logger
	logger.InitLogger()

	// Set Gin to release mode and redirect logs to file
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = logger.LogWriter

	// Load config
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to database
	db, err := sql.Open("mysql", cfg.Database.URL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Setup router
	r := http.SetupRouter(db)

	// Get port from config, default to 8080
	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}

	// Start server
	logger.Logger.Println("Server starting on :" + port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
