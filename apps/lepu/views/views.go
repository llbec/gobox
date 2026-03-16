package views

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/llbec/gobox/apps/lepu/config"
	"github.com/llbec/gobox/apps/lepu/logger"

	_ "github.com/go-sql-driver/mysql"
)

type CaseInfo struct {
	MonitorCaseID string     `json:"monitor_case_id"`
	PID           string     `json:"pid"`
	PatientName   string     `json:"patient_name"`
	NTestNo       *string    `json:"n_test_no"`
	NBarNo        *string    `json:"n_bar_no"`
	NAdNo         *string    `json:"n_ad_no"`
	StartTime     time.Time  `json:"start_time"`
	EndTime       *time.Time `json:"end_time"`
}

type ECGData struct {
	MonitorCaseID string    `json:"monitor_case_id"`
	DataTime      time.Time `json:"data_time"`
	HR            *int      `json:"hr"`
}

type SPO2Data struct {
	MonitorCaseID string    `json:"monitor_case_id"`
	DataTime      time.Time `json:"data_time"`
	SPO2          *int      `json:"spo2"`
	PR            *int      `json:"pr"`
}

type RespData struct {
	MonitorCaseID string    `json:"monitor_case_id"`
	DataTime      time.Time `json:"data_time"`
	RR            *int      `json:"rr"`
}

type NIBPData struct {
	MonitorCaseID string    `json:"monitor_case_id"`
	DataTime      time.Time `json:"data_time"`
	SYS           *int      `json:"sys"`
	DIA           *int      `json:"dia"`
	MAP           *int      `json:"map"`
}

type VitalData struct {
	Name  string    `json:"name"`
	Value int       `json:"value"`
	Time  time.Time `json:"time"`
}

func GetPatients(db *sql.DB, startTime time.Time) ([]CaseInfo, error) {
	query := `SELECT monitor_case_id, pid, patient_name, n_test_no, n_bar_no, n_ad_no, start_time, end_time FROM view_case_info WHERE start_time > ?`
	rows, err := db.Query(query, startTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []CaseInfo
	for rows.Next() {
		var p CaseInfo
		var startTimeStr, endTimeStr sql.NullString
		err := rows.Scan(&p.MonitorCaseID, &p.PID, &p.PatientName, &p.NTestNo, &p.NBarNo, &p.NAdNo, &startTimeStr, &endTimeStr)
		if err != nil {
			return nil, err
		}

		// Parse start_time
		if startTimeStr.Valid {
			if t, err := time.ParseInLocation("2006-01-02 15:04:05", startTimeStr.String, config.TimezoneLocation); err == nil {
				p.StartTime = t
			} else {
				// Try different format
				if t, err := time.Parse(time.RFC3339, startTimeStr.String); err == nil {
					p.StartTime = t
				}
			}
		}

		// Parse end_time
		if endTimeStr.Valid {
			if t, err := time.ParseInLocation("2006-01-02 15:04:05", endTimeStr.String, config.TimezoneLocation); err == nil {
				p.EndTime = &t
			} else {
				// Try different format
				if t, err := time.Parse(time.RFC3339, endTimeStr.String); err == nil {
					p.EndTime = &t
				}
			}
		}

		patients = append(patients, p)
	}
	return patients, nil
}

func GetVitals(db *sql.DB, caseID string, start, end time.Time) ([]VitalData, error) {
	// Query each view
	ecgData, err := GetECGData(db, caseID, start, end)
	if err != nil {
		return nil, err
	}
	spo2Data, err := GetSPO2Data(db, caseID, start, end)
	if err != nil {
		return nil, err
	}
	respData, err := GetRespData(db, caseID, start, end)
	if err != nil {
		return nil, err
	}
	nibpData, err := GetNIBPData(db, caseID, start, end)
	if err != nil {
		return nil, err
	}

	var result []VitalData
	totalRecords := 0
	filteredRecords := 0

	// Process ECG data
	for _, ecg := range ecgData {
		totalRecords++
		if ecg.HR != nil {
			result = append(result, VitalData{
				Name:  "hr",
				Value: *ecg.HR,
				Time:  ecg.DataTime,
			})
		} else {
			filteredRecords++
		}
	}

	// Process SPO2 data
	for _, spo2 := range spo2Data {
		totalRecords++
		if spo2.SPO2 != nil {
			result = append(result, VitalData{
				Name:  "spo2",
				Value: *spo2.SPO2,
				Time:  spo2.DataTime,
			})
		} else {
			filteredRecords++
		}
		totalRecords++
		if spo2.PR != nil {
			result = append(result, VitalData{
				Name:  "pr",
				Value: *spo2.PR,
				Time:  spo2.DataTime,
			})
		} else {
			filteredRecords++
		}
	}

	// Process Resp data
	for _, resp := range respData {
		totalRecords++
		if resp.RR != nil {
			result = append(result, VitalData{
				Name:  "rr",
				Value: *resp.RR,
				Time:  resp.DataTime,
			})
		} else {
			filteredRecords++
		}
	}

	// Process NIBP data
	for _, nibp := range nibpData {
		totalRecords++
		if nibp.SYS != nil {
			result = append(result, VitalData{
				Name:  "nIBPs",
				Value: *nibp.SYS,
				Time:  nibp.DataTime,
			})
		} else {
			filteredRecords++
		}
		totalRecords++
		if nibp.DIA != nil {
			result = append(result, VitalData{
				Name:  "nIBPd",
				Value: *nibp.DIA,
				Time:  nibp.DataTime,
			})
		} else {
			filteredRecords++
		}
		totalRecords++
		if nibp.MAP != nil {
			result = append(result, VitalData{
				Name:  "nIBPm",
				Value: *nibp.MAP,
				Time:  nibp.DataTime,
			})
		} else {
			filteredRecords++
		}
	}

	// Log total and filtered records
	logger.Logger.Printf("Total records processed: %d, Filtered records (--- values): %d\n", totalRecords, filteredRecords)

	return result, nil
}

func GetECGData(db *sql.DB, caseID string, start, end time.Time) ([]ECGData, error) {
	query := `SELECT monitor_case_id, data_time, hr FROM view_case_ecg WHERE monitor_case_id = ? AND data_time BETWEEN ? AND ?`
	rows, err := db.Query(query, caseID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []ECGData
	for rows.Next() {
		var d ECGData
		var dataTimeStr, hrStr string
		err := rows.Scan(&d.MonitorCaseID, &dataTimeStr, &hrStr)
		if err != nil {
			return nil, err
		}

		// Parse data_time
		if t, err := time.ParseInLocation("2006-01-02 15:04:05", dataTimeStr, config.TimezoneLocation); err == nil {
			d.DataTime = t
		} else {
			// Try different format
			if t, err := time.Parse(time.RFC3339, dataTimeStr); err == nil {
				d.DataTime = t
			}
		}

		d.HR = parseIntValue(hrStr)
		data = append(data, d)
	}
	return data, nil
}

func GetSPO2Data(db *sql.DB, caseID string, start, end time.Time) ([]SPO2Data, error) {
	query := `SELECT monitor_case_id, data_time, spo2, pr FROM view_case_spo2 WHERE monitor_case_id = ? AND data_time BETWEEN ? AND ?`
	rows, err := db.Query(query, caseID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []SPO2Data
	for rows.Next() {
		var d SPO2Data
		var dataTimeStr, spo2Str, prStr string
		err := rows.Scan(&d.MonitorCaseID, &dataTimeStr, &spo2Str, &prStr)
		if err != nil {
			return nil, err
		}

		// Parse data_time
		if t, err := time.ParseInLocation("2006-01-02 15:04:05", dataTimeStr, config.TimezoneLocation); err == nil {
			d.DataTime = t
		} else {
			// Try different format
			if t, err := time.Parse(time.RFC3339, dataTimeStr); err == nil {
				d.DataTime = t
			}
		}

		d.SPO2 = parseIntValue(spo2Str)
		d.PR = parseIntValue(prStr)
		data = append(data, d)
	}
	return data, nil
}

func GetRespData(db *sql.DB, caseID string, start, end time.Time) ([]RespData, error) {
	query := `SELECT monitor_case_id, data_time, rr FROM view_case_resp WHERE monitor_case_id = ? AND data_time BETWEEN ? AND ?`
	rows, err := db.Query(query, caseID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []RespData
	for rows.Next() {
		var d RespData
		var dataTimeStr, rrStr string
		err := rows.Scan(&d.MonitorCaseID, &dataTimeStr, &rrStr)
		if err != nil {
			return nil, err
		}

		// Parse data_time
		if t, err := time.ParseInLocation("2006-01-02 15:04:05", dataTimeStr, config.TimezoneLocation); err == nil {
			d.DataTime = t
		} else {
			// Try different format
			if t, err := time.Parse(time.RFC3339, dataTimeStr); err == nil {
				d.DataTime = t
			}
		}

		d.RR = parseIntValue(rrStr)
		data = append(data, d)
	}
	return data, nil
}

func GetNIBPData(db *sql.DB, caseID string, start, end time.Time) ([]NIBPData, error) {
	query := `SELECT monitor_case_id, data_time, sys, dia, map FROM view_case_nibp WHERE monitor_case_id = ? AND data_time BETWEEN ? AND ?`
	rows, err := db.Query(query, caseID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []NIBPData
	for rows.Next() {
		var d NIBPData
		var dataTimeStr, sysStr, diaStr, mapStr string
		err := rows.Scan(&d.MonitorCaseID, &dataTimeStr, &sysStr, &diaStr, &mapStr)
		if err != nil {
			return nil, err
		}

		// Parse data_time
		if t, err := time.ParseInLocation("2006-01-02 15:04:05", dataTimeStr, config.TimezoneLocation); err == nil {
			d.DataTime = t
		} else {
			// Try different format
			if t, err := time.Parse(time.RFC3339, dataTimeStr); err == nil {
				d.DataTime = t
			}
		}

		d.SYS = parseIntValue(sysStr)
		d.DIA = parseIntValue(diaStr)
		d.MAP = parseIntValue(mapStr)
		data = append(data, d)
	}
	return data, nil
}

func TestDBConnection(db *sql.DB) error {
	return db.Ping()
}

// parseIntValue converts string value to *int, treating "---" as nil
func parseIntValue(s string) *int {
	if s == "---" {
		return nil
	}
	// Try to parse as int
	var val int
	_, err := fmt.Sscanf(s, "%d", &val)
	if err != nil {
		return nil
	}
	return &val
}
