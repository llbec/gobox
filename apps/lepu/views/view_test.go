package views_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/llbec/gobox/apps/lepu/config"
	"github.com/llbec/gobox/apps/lepu/views"

	_ "github.com/go-sql-driver/mysql"
)

var testDB *sql.DB

var (
	testCaseID = "VOn_v8gq--1dyM3lWNKWniLJ"
	testStart  = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	testEnd    = time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC)
)

func TestMain(m *testing.M) {
	// Load config
	cfg, err := config.LoadConfig("../config.yaml")
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	// Connect to database
	testDB, err = sql.Open("mysql", cfg.Database.URL)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Test connection
	if err := testDB.Ping(); err != nil {
		panic("Failed to ping database: " + err.Error())
	}

	// Run tests
	m.Run()

	// Close connection
	testDB.Close()
}

func TestGetPatients(t *testing.T) {
	patients, err := views.GetPatients(testDB, testStart)
	if err != nil {
		t.Fatalf("GetPatients failed: %v", err)
	}

	t.Logf("Found %d patients", len(patients))

	// Basic validation
	for i, patient := range patients {
		t.Logf("Patient %d: MonitorCaseID=%s, PID=%s, PatientName=%s, StartTime=%v, EndTime=%v",
			i, patient.MonitorCaseID, patient.PID, patient.PatientName, patient.StartTime, patient.EndTime)
		if patient.NTestNo != nil {
			t.Logf("  NTestNo: %s", *patient.NTestNo)
		} else {
			t.Logf("  NTestNo: null")
		}
		if patient.NBarNo != nil {
			t.Logf("  NBarNo: %s", *patient.NBarNo)
		} else {
			t.Logf("  NBarNo: null")
		}
		if patient.NAdNo != nil {
			t.Logf("  NAdNo: %s", *patient.NAdNo)
		} else {
			t.Logf("  NAdNo: null")
		}

		if patient.MonitorCaseID == "" {
			t.Errorf("Patient %d has empty MonitorCaseID", i)
		}
		if patient.PID == "" {
			t.Errorf("Patient %d has empty PID", i)
		}
		if patient.PatientName == "" {
			t.Errorf("Patient %d has empty PatientName", i)
		}
		// StartTime should be valid
		if patient.StartTime.IsZero() {
			t.Errorf("Patient %d has zero StartTime", i)
		}
	}
}

func TestGetECGData(t *testing.T) {
	data, err := views.GetECGData(testDB, testCaseID, testStart, testEnd)
	if err != nil {
		t.Fatalf("getECGData failed: %v", err)
	}

	t.Logf("Found %d ECG records for case %s", len(data), testCaseID)

	// Basic validation
	for i, record := range data {
		t.Logf("ECG Record %d: MonitorCaseID=%s, DataTime=%v, HR=%v", i, record.MonitorCaseID, record.DataTime, record.HR)

		if record.MonitorCaseID != testCaseID {
			t.Errorf("Record %d has wrong MonitorCaseID: expected %s, got %s", i, testCaseID, record.MonitorCaseID)
		}
		if record.DataTime.Before(testStart) || record.DataTime.After(testEnd) {
			t.Errorf("Record %d DataTime %v is outside range [%v, %v]", i, record.DataTime, testStart, testEnd)
		}
	}
}

func TestGetSPO2Data(t *testing.T) {
	data, err := views.GetSPO2Data(testDB, testCaseID, testStart, testEnd)
	if err != nil {
		t.Fatalf("getSPO2Data failed: %v", err)
	}

	t.Logf("Found %d SPO2 records for case %s", len(data), testCaseID)

	// Basic validation
	for i, record := range data {
		t.Logf("SPO2 Record %d: MonitorCaseID=%s, DataTime=%v, SPO2=%v, PR=%v", i, record.MonitorCaseID, record.DataTime, record.SPO2, record.PR)

		if record.MonitorCaseID != testCaseID {
			t.Errorf("Record %d has wrong MonitorCaseID: expected %s, got %s", i, testCaseID, record.MonitorCaseID)
		}
		if record.DataTime.Before(testStart) || record.DataTime.After(testEnd) {
			t.Errorf("Record %d DataTime %v is outside range [%v, %v]", i, record.DataTime, testStart, testEnd)
		}
	}
}

func TestGetRespData(t *testing.T) {
	data, err := views.GetRespData(testDB, testCaseID, testStart, testEnd)
	if err != nil {
		t.Fatalf("getRespData failed: %v", err)
	}

	t.Logf("Found %d Resp records for case %s", len(data), testCaseID)

	// Basic validation
	for i, record := range data {
		t.Logf("Resp Record %d: MonitorCaseID=%s, DataTime=%v, RR=%v", i, record.MonitorCaseID, record.DataTime, record.RR)

		if record.MonitorCaseID != testCaseID {
			t.Errorf("Record %d has wrong MonitorCaseID: expected %s, got %s", i, testCaseID, record.MonitorCaseID)
		}
		if record.DataTime.Before(testStart) || record.DataTime.After(testEnd) {
			t.Errorf("Record %d DataTime %v is outside range [%v, %v]", i, record.DataTime, testStart, testEnd)
		}
	}
}

func TestGetNIBPData(t *testing.T) {
	data, err := views.GetNIBPData(testDB, testCaseID, testStart, testEnd)
	if err != nil {
		t.Fatalf("getNIBPData failed: %v", err)
	}

	t.Logf("Found %d NIBP records for case %s", len(data), testCaseID)

	// Basic validation
	for i, record := range data {
		sysVal := "nil"
		if record.SYS != nil {
			sysVal = fmt.Sprintf("%d", *record.SYS)
		}
		diaVal := "nil"
		if record.DIA != nil {
			diaVal = fmt.Sprintf("%d", *record.DIA)
		}
		mapVal := "nil"
		if record.MAP != nil {
			mapVal = fmt.Sprintf("%d", *record.MAP)
		}
		t.Logf("NIBP Record %d: MonitorCaseID=%s, DataTime=%v, SYS=%s, DIA=%s, MAP=%s", i, record.MonitorCaseID, record.DataTime, sysVal, diaVal, mapVal)

		if record.MonitorCaseID != testCaseID {
			t.Errorf("Record %d has wrong MonitorCaseID: expected %s, got %s", i, testCaseID, record.MonitorCaseID)
		}
		if record.DataTime.Before(testStart) || record.DataTime.After(testEnd) {
			t.Errorf("Record %d DataTime %v is outside range [%v, %v]", i, record.DataTime, testStart, testEnd)
		}
	}
}

func TestGetVitals(t *testing.T) {
	vitals, err := views.GetVitals(testDB, testCaseID, testStart, testEnd)
	if err != nil {
		t.Fatalf("GetVitals failed: %v", err)
	}

	t.Logf("Found %d vital records for case %s", len(vitals), testCaseID)

	// Basic validation - check that all records have valid data
	for i, record := range vitals {
		t.Logf("Vital Record %d: Name=%s, Value=%d, Time=%v", i, record.Name, record.Value, record.Time)

		if record.Name == "" {
			t.Errorf("Record %d has empty name", i)
		}
		if record.Time.IsZero() {
			t.Errorf("Record %d has zero time", i)
		}
		if record.Time.Before(testStart) || record.Time.After(testEnd) {
			t.Errorf("Record %d Time %v is outside range [%v, %v]", i, record.Time, testStart, testEnd)
		}
		// Value should be positive for valid records
		if record.Value <= 0 {
			t.Errorf("Record %d has invalid value: %d", i, record.Value)
		}
	}
}
