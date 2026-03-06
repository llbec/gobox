package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Uptime    string    `json:"uptime"`
	GoVersion string    `json:"go_version"`
	Hostname  string    `json:"hostname"`
}

var startTime time.Time

// healthHandler 处理健康检查请求
func healthHandler(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(startTime)
	hostname, _ := os.Hostname()

	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Uptime:    uptime.String(),
		GoVersion: runtime.Version(),
		Hostname:  hostname,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Record start time
	startTime = time.Now()

	// Create a server with tools.
	server := createMcpServer()

	// Create HTTP handler for MCP
	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, nil)

	// Set up HTTP routes
	http.Handle("/mcp", handler)
	http.HandleFunc("/health", healthHandler)

	// Start HTTP server
	port := 8080
	addr := fmt.Sprintf(":%d", port)
	log.Printf("MCP HTTP server starting on port %d...", port)
	log.Printf("Access MCP API at: http://localhost:%d/mcp", port)
	log.Printf("Health check endpoint: http://localhost:%d/health", port)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
