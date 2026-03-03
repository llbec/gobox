package act

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type DocumentInfo struct {
	Filename  string
	PatientID string
	TestTime  time.Time
}

type TimeFilterType int

const (
	TimeEqual TimeFilterType = iota
	TimeBefore
	TimeAfter
	TimeRange
)

type TimeFilter struct {
	Type      TimeFilterType
	Time      time.Time
	StartTime time.Time
	EndTime   time.Time
}

type Filter struct {
	TimeFilter *TimeFilter
	PatientID  string
}

type ActManager struct {
	TargetPath string
	BackupPath string
}

var (
	acttimePattern    = regexp.MustCompile(`TestDateTime="([^"]+)"`)
	actpatientPattern = regexp.MustCompile(`PatientID="([^"]+)"`)
)

func NewActManager(targetPath, backupPath string) *ActManager {
	return &ActManager{
		TargetPath: targetPath,
		BackupPath: backupPath,
	}
}

func (am *ActManager) FilterFiles(filter Filter) ([]DocumentInfo, error) {
	var filtered []DocumentInfo

	files, err := os.ReadDir(am.TargetPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !strings.HasSuffix(strings.ToLower(file.Name()), ".xml") {
			continue
		}

		filePath := filepath.Join(am.TargetPath, file.Name())
		content, err := readFile(filePath)
		if err != nil {
			continue
		}

		patientID := ""
		testTimeStr := ""

		if match := actpatientPattern.FindStringSubmatch(string(content)); len(match) > 1 {
			patientID = match[1]
		}

		if match := acttimePattern.FindStringSubmatch(string(content)); len(match) > 1 {
			testTimeStr = match[1]
		}

		if patientID == "" || testTimeStr == "" {
			continue
		}

		testTime, err := parseTime(testTimeStr)
		if err != nil {
			continue
		}

		match := true

		if filter.PatientID != "" && patientID != filter.PatientID {
			match = false
		}

		if match && filter.TimeFilter != nil {
			switch filter.TimeFilter.Type {
			case TimeEqual:
				match = testTime.Equal(filter.TimeFilter.Time)
			case TimeBefore:
				match = testTime.Before(filter.TimeFilter.Time)
			case TimeAfter:
				match = testTime.After(filter.TimeFilter.Time)
			case TimeRange:
				match = !testTime.Before(filter.TimeFilter.StartTime) && !testTime.After(filter.TimeFilter.EndTime)
			}
		}

		if match {
			filtered = append(filtered, DocumentInfo{
				Filename:  file.Name(),
				PatientID: patientID,
				TestTime:  testTime,
			})
		}
	}

	return filtered, nil
}

func (am *ActManager) ArchiveFiles(beforeTime time.Time) error {
	filter := Filter{
		TimeFilter: &TimeFilter{
			Type: TimeBefore,
			Time: beforeTime,
		},
	}

	docs, err := am.FilterFiles(filter)
	if err != nil {
		return err
	}

	for _, doc := range docs {
		year := doc.TestTime.Year()
		month := int(doc.TestTime.Month())
		archivePath := filepath.Join(am.BackupPath, fmt.Sprintf("%04d", year), fmt.Sprintf("%02d", month))

		if err := os.MkdirAll(archivePath, 0755); err != nil {
			return err
		}

		srcPath := filepath.Join(am.TargetPath, doc.Filename)
		dstPath := filepath.Join(archivePath, doc.Filename)

		if err := moveFile(srcPath, dstPath); err != nil {
			return err
		}
	}

	return nil
}

func (am *ActManager) ResendFiles(filenames []string) error {
	resendBasePath := filepath.Join(am.BackupPath, "resend")
	if err := os.MkdirAll(resendBasePath, 0755); err != nil {
		return err
	}

	dirNum, err := getNextDirNumber(resendBasePath)
	if err != nil {
		return err
	}

	resendDir := filepath.Join(resendBasePath, fmt.Sprintf("%d", dirNum))
	if err = os.Mkdir(resendDir, 0755); err != nil {
		return err
	}

	for _, filename := range filenames {
		srcPath := filepath.Join(am.TargetPath, filename)
		dstPath := filepath.Join(resendDir, filename)
		if err = moveFile(srcPath, dstPath); err != nil {
			return err
		}
	}

	files, err := os.ReadDir(resendDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		time.Sleep(3 * time.Second)

		filePath := filepath.Join(resendDir, file.Name())
		content, err := readFile(filePath)
		if err != nil {
			continue
		}

		testTimeStr := ""
		if match := acttimePattern.FindStringSubmatch(string(content)); len(match) > 1 {
			testTimeStr = match[1]
		}

		testTime, err := parseTime(testTimeStr)
		if err != nil {
			continue
		}

		now := time.Now()
		newFilename := fmt.Sprintf("R%s_%s.xml", now.Format("200601021504"), testTime.Format("20060102150405"))
		newFilePath := filepath.Join(am.TargetPath, newFilename)

		if err := writeFile(newFilePath, content); err != nil {
			continue
		}
	}

	return nil
}

func (am *ActManager) ModifyFiles(filenames []string, newPatientID string) error {
	modifyBasePath := filepath.Join(am.BackupPath, "modify")
	if err := os.MkdirAll(modifyBasePath, 0755); err != nil {
		return err
	}

	dirNum, err := getNextDirNumber(modifyBasePath)
	if err != nil {
		return err
	}

	modifyDir := filepath.Join(modifyBasePath, fmt.Sprintf("%d", dirNum))
	if err = os.Mkdir(modifyDir, 0755); err != nil {
		return err
	}

	for _, filename := range filenames {
		srcPath := filepath.Join(am.TargetPath, filename)
		dstPath := filepath.Join(modifyDir, filename)
		if err = moveFile(srcPath, dstPath); err != nil {
			return err
		}
	}

	files, err := os.ReadDir(modifyDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		time.Sleep(3 * time.Second)

		filePath := filepath.Join(modifyDir, file.Name())
		content, err := readFile(filePath)
		if err != nil {
			continue
		}

		testTimeStr := ""
		if match := acttimePattern.FindStringSubmatch(string(content)); len(match) > 1 {
			testTimeStr = match[1]
		}

		testTime, err := parseTime(testTimeStr)
		if err != nil {
			continue
		}

		modifiedContent := actpatientPattern.ReplaceAllString(string(content), fmt.Sprintf(`PatientID="%s"`, newPatientID))

		now := time.Now()
		newFilename := fmt.Sprintf("M%s_%s.xml", now.Format("200601021504"), testTime.Format("20060102150405"))
		newFilePath := filepath.Join(am.TargetPath, newFilename)

		if err := writeFile(newFilePath, []byte(modifiedContent)); err != nil {
			continue
		}
	}

	return nil
}

func readFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

func writeFile(filePath string, content []byte) error {
	return os.WriteFile(filePath, content, 0644)
}

func moveFile(src, dst string) error {
	input, err := os.Open(src)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer output.Close()

	if _, err = io.Copy(output, input); err != nil {
		return err
	}

	if err = input.Close(); err != nil {
		return err
	}

	return os.Remove(src)
}

func getNextDirNumber(basePath string) (int, error) {
	files, err := os.ReadDir(basePath)
	if err != nil {
		return 0, err
	}

	var dirNums []int
	for _, file := range files {
		if file.IsDir() {
			num, err := strconv.Atoi(file.Name())
			if err == nil {
				dirNums = append(dirNums, num)
			}
		}
	}

	if len(dirNums) == 0 {
		return 1, nil
	}

	sort.Ints(dirNums)
	return dirNums[len(dirNums)-1] + 1, nil
}

func parseTime(timeStr string) (time.Time, error) {
	formats := []string{
		"2006-01-02 15:04:05",
		"2006/01/02 15:04:05",
		"20060102150405",
		"2006-01-02T15:04:05",
	}

	for _, format := range formats {
		t, err := time.ParseInLocation(format, timeStr, time.Local)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time: %s", timeStr)
}
