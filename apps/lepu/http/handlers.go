package http

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/llbec/gobox/apps/lepu/logger"
	"github.com/llbec/gobox/apps/lepu/views"

	"github.com/gin-gonic/gin"
)

// Response structs with Unix timestamps
type PatientResponse struct {
	MonitorCaseID string  `json:"monitor_case_id"`
	PID           string  `json:"pid"`
	PatientName   string  `json:"patient_name"`
	NTestNo       *string `json:"n_test_no"`
	NBarNo        *string `json:"n_bar_no"`
	NAdNo         *string `json:"n_ad_no"`
	StartTime     int64   `json:"start_time"`
	EndTime       *int64  `json:"end_time"`
}

type VitalResponse struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
	Time  int64  `json:"time"`
}

func SetupRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()

	r.POST("/health", func(c *gin.Context) {
		if err := views.TestDBConnection(db); err != nil {
			logger.Logger.Println("Database connection test failed:", err)
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "Database connection failed", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	r.POST("/patients", func(c *gin.Context) {
		var req struct {
			StartTime int64 `json:"start_time" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		startTime := time.Unix(req.StartTime, 0)

		patients, err := views.GetPatients(db, startTime)
		if err != nil {
			logger.Logger.Println("Error getting patients:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		// Convert to response format with Unix timestamps
		var response []PatientResponse
		for _, p := range patients {
			resp := PatientResponse{
				MonitorCaseID: p.MonitorCaseID,
				PID:           p.PID,
				PatientName:   p.PatientName,
				NTestNo:       p.NTestNo,
				NBarNo:        p.NBarNo,
				NAdNo:         p.NAdNo,
				StartTime:     p.StartTime.Unix(),
			}
			if p.EndTime != nil {
				endTime := p.EndTime.Unix()
				resp.EndTime = &endTime
			}
			response = append(response, resp)
		}

		c.JSON(http.StatusOK, response)
	})

	r.POST("/vitals", func(c *gin.Context) {
		var req struct {
			CaseID string `json:"case_id" binding:"required"`
			Start  int64  `json:"start" binding:"required"`
			End    int64  `json:"end" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		start := time.Unix(req.Start, 0)
		end := time.Unix(req.End, 0)

		vitals, err := views.GetVitals(db, req.CaseID, start, end)
		if err != nil {
			logger.Logger.Println("Error getting vitals:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		// Convert to response format with Unix timestamps
		var response []VitalResponse
		for _, v := range vitals {
			response = append(response, VitalResponse{
				Name:  v.Name,
				Value: v.Value,
				Time:  v.Time.Unix(),
			})
		}

		c.JSON(http.StatusOK, response)
	})

	return r
}
