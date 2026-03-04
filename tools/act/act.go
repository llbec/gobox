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
	timeCache         = make(map[string]time.Time)
)

// NewActManager 创建一个新的 ActManager 实例
// targetPath: 目标文件路径
// backupPath: 备份文件路径
func NewActManager(targetPath, backupPath string) *ActManager {
	return &ActManager{
		TargetPath: targetPath,
		BackupPath: backupPath,
	}
}

// FilterFiles 根据过滤条件筛选文件
func (am *ActManager) FilterFiles(filter Filter) ([]DocumentInfo, string, error) {
	files, err := os.ReadDir(am.TargetPath)
	if err != nil {
		return nil, "", fmt.Errorf("读取目标目录失败: %w", err)
	}

	// 过滤出 XML 文件
	var xmlFiles []os.DirEntry
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !strings.HasSuffix(strings.ToLower(file.Name()), ".xml") {
			continue
		}
		xmlFiles = append(xmlFiles, file)
	}

	// 创建通道用于接收处理结果
	resultCh := make(chan DocumentInfo, len(xmlFiles))
	doneCh := make(chan struct{}, len(xmlFiles))

	// 处理每个文件
	for _, file := range xmlFiles {
		go func(f os.DirEntry) {
			defer func() {
				doneCh <- struct{}{}
			}()

			filePath := filepath.Join(am.TargetPath, f.Name())
			content, err := readFile(filePath)
			if err != nil {
				return
			}

			patientID := ""
			testTimeStr := ""

			if match := actpatientPattern.FindSubmatch(content); len(match) > 1 {
				patientID = string(match[1])
			}

			if match := acttimePattern.FindSubmatch(content); len(match) > 1 {
				testTimeStr = string(match[1])
			}

			if patientID == "" || testTimeStr == "" {
				return
			}

			testTime, err := parseTime(testTimeStr)
			if err != nil {
				return
			}

			match := filter.PatientID == "" || patientID == filter.PatientID

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
				resultCh <- DocumentInfo{
					Filename:  f.Name(),
					PatientID: patientID,
					TestTime:  testTime,
				}
			}
		}(file)
	}

	// 收集结果
	var filtered []DocumentInfo
	for i := 0; i < len(xmlFiles); i++ {
		select {
		case doc := <-resultCh:
			filtered = append(filtered, doc)
		case <-doneCh:
			// 处理完成，继续等待
		}
	}

	result := fmt.Sprintf("筛选完成，共找到 %d 个符合条件的文件", len(filtered))
	return filtered, result, nil
}

// ArchiveFiles 归档指定时间之前的文件
func (am *ActManager) ArchiveFiles(beforeTime time.Time) (string, error) {
	filter := Filter{
		TimeFilter: &TimeFilter{
			Type: TimeBefore,
			Time: beforeTime,
		},
	}

	docs, _, err := am.FilterFiles(filter)
	if err != nil {
		return "", err
	}

	// 按年月份统计文件数量
	monthlyStats := make(map[string]int)
	for _, doc := range docs {
		year := doc.TestTime.Year()
		month := int(doc.TestTime.Month())
		key := fmt.Sprintf("%d-%02d", year, month)
		monthlyStats[key]++

		archivePath := filepath.Join(am.BackupPath, fmt.Sprintf("%04d", year), fmt.Sprintf("%02d", month))

		if err := os.MkdirAll(archivePath, 0755); err != nil {
			return "", fmt.Errorf("创建归档目录失败: %w", err)
		}

		srcPath := filepath.Join(am.TargetPath, doc.Filename)
		dstPath := filepath.Join(archivePath, doc.Filename)

		if err := moveFile(srcPath, dstPath); err != nil {
			return "", fmt.Errorf("归档文件失败: %w", err)
		}
	}

	// 构建返回字符串
	result := fmt.Sprintf("归档完成，共处理 %d 个文件: \n", len(docs))

	// 按年份和月份排序并输出
	var keys []string
	for key := range monthlyStats {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		yearMonth := strings.Split(key, "-")
		year, _ := strconv.Atoi(yearMonth[0])
		month, _ := strconv.Atoi(yearMonth[1])
		count := monthlyStats[key]
		result = fmt.Sprintf("%s%d年    %d月    %d个文件\n", result, year, month, count)
	}

	return result, nil
}

// prepareOperationDir 创建操作目录并返回目录路径
func (am *ActManager) prepareOperationDir(operation string) (string, error) {
	basePath := filepath.Join(am.BackupPath, operation)
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return "", err
	}

	dirNum, err := getNextDirNumber(basePath)
	if err != nil {
		return "", err
	}

	opDir := filepath.Join(basePath, fmt.Sprintf("%d", dirNum))
	if err = os.Mkdir(opDir, 0755); err != nil {
		return "", err
	}

	return opDir, nil
}

// moveFilesToDir 将文件从目标路径移动到指定目录
func (am *ActManager) moveFilesToDir(docs []DocumentInfo, dstDir string) error {
	for _, doc := range docs {
		srcPath := filepath.Join(am.TargetPath, doc.Filename)
		dstPath := filepath.Join(dstDir, doc.Filename)
		if err := moveFile(srcPath, dstPath); err != nil {
			return err
		}
	}
	return nil
}

// processFiles 处理目录中的文件，根据操作类型生成新文件
func (am *ActManager) processFiles(opDir string, operation string, newPatientID string) (string, error) {
	files, err := os.ReadDir(opDir)
	if err != nil {
		return "", err
	}

	scanCount := 0
	matchCount := 0
	generatedCount := 0

	// 存储文件详情
	type fileDetail struct {
		index        int
		filename     string
		patientID    string
		testTime     time.Time
		originalFile string
	}
	var fileDetails []fileDetail

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		scanCount++

		time.Sleep(3 * time.Second)

		filePath := filepath.Join(opDir, file.Name())
		content, err := readFile(filePath)
		if err != nil {
			continue
		}

		testTimeStr := ""
		patientID := ""

		if match := actpatientPattern.FindSubmatch(content); len(match) > 1 {
			patientID = string(match[1])
		}

		if match := acttimePattern.FindSubmatch(content); len(match) > 1 {
			testTimeStr = string(match[1])
		}

		if patientID == "" || testTimeStr == "" {
			continue
		}

		testTime, err := parseTime(testTimeStr)
		if err != nil {
			continue
		}

		matchCount++

		var newContent []byte
		var prefix string

		switch operation {
		case "resend":
			newContent = content
			prefix = "R"
		case "modify":
			modifiedContent := actpatientPattern.ReplaceAll(content, fmt.Appendf(nil, `PatientID="%s"`, newPatientID))
			newContent = modifiedContent
			prefix = "M"
		}

		now := time.Now()
		newFilename := fmt.Sprintf("%s%s_%s.xml", prefix, now.Format("200601021504"), testTime.Format("20060102150405"))
		newFilePath := filepath.Join(am.TargetPath, newFilename)

		if err := writeFile(newFilePath, newContent); err != nil {
			continue
		}

		generatedCount++
		fileDetails = append(fileDetails, fileDetail{
			index:        generatedCount,
			filename:     newFilename,
			patientID:    patientID,
			testTime:     testTime,
			originalFile: file.Name(),
		})
	}

	// 创建结果文件
	resultFilePath := filepath.Join(opDir, "result.txt")
	resultFile, err := os.Create(resultFilePath)
	if err == nil {
		defer resultFile.Close()

		for _, detail := range fileDetails {
			line := fmt.Sprintf("%s|%s\n", detail.filename, detail.originalFile)
			_, _ = resultFile.WriteString(line)
		}
	}

	// 构建返回字符串
	result := fmt.Sprintf("处理完成，共扫描%d个文件，匹配%d个文件，生成%d个新文件：\n", scanCount, matchCount, generatedCount)
	result += "index       file-name  ID   test-time \n"

	// 添加文件详情
	for _, detail := range fileDetails {
		result = fmt.Sprintf("%s%5d  %-10s  %-5s  %s\n", result, detail.index, detail.filename, detail.patientID, detail.testTime.Format("2006-01-02 15:04:05"))
	}

	return result, nil
}

// ResendFiles 将指定文件移动到重发目录并重新发送
func (am *ActManager) ResendFiles(docs []DocumentInfo) (string, error) {
	resendDir, err := am.prepareOperationDir("resend")
	if err != nil {
		return "", err
	}

	if err := am.moveFilesToDir(docs, resendDir); err != nil {
		return "", err
	}

	return am.processFiles(resendDir, "resend", "")
}

// ModifyFiles 修改指定文件的患者ID并重新发送
func (am *ActManager) ModifyFiles(docs []DocumentInfo, newPatientID string) (string, error) {
	modifyDir, err := am.prepareOperationDir("modify")
	if err != nil {
		return "", err
	}

	if err := am.moveFilesToDir(docs, modifyDir); err != nil {
		return "", err
	}

	return am.processFiles(modifyDir, "modify", newPatientID)
}

// Rollback 回滚操作
func (am *ActManager) Rollback(operation string, index int) (string, error) {
	basePath := filepath.Join(am.BackupPath, operation)
	opDir := filepath.Join(basePath, fmt.Sprintf("%d", index))

	if _, err := os.Stat(opDir); os.IsNotExist(err) {
		return "", fmt.Errorf("操作目录不存在: %s", opDir)
	}

	resultFilePath := filepath.Join(opDir, "result.txt")
	content, err := readFile(resultFilePath)
	if err != nil {
		return "", fmt.Errorf("读取结果文件失败: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	rollbackCount := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) != 2 {
			continue
		}

		newFilename := parts[0]
		originalFilename := parts[1]

		newFilePath := filepath.Join(am.TargetPath, newFilename)
		originalFilePath := filepath.Join(opDir, originalFilename)

		if err := os.Remove(newFilePath); err != nil && !os.IsNotExist(err) {
			continue
		}

		targetOriginalPath := filepath.Join(am.TargetPath, originalFilename)
		if err := moveFile(originalFilePath, targetOriginalPath); err != nil {
			continue
		}

		rollbackCount++
	}

	return fmt.Sprintf("回滚完成，共回滚 %d 个文件", rollbackCount), nil
}

// readFile 读取文件内容
func readFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

// writeFile 写入文件内容
func writeFile(filePath string, content []byte) error {
	return os.WriteFile(filePath, content, 0644)
}

// moveFile 移动文件，优先使用 os.Rename，失败时回退到复制后删除
func moveFile(src, dst string) error {
	// 尝试使用 os.Rename 直接移动文件（同一文件系统内更高效）
	if err := os.Rename(src, dst); err == nil {
		return nil
	}

	// 如果 os.Rename 失败（可能是跨文件系统），使用复制后删除
	input, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		_ = input.Close()
	}()

	output, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		_ = output.Close()
	}()

	if _, err = io.Copy(output, input); err != nil {
		return err
	}

	return os.Remove(src)
}

// getNextDirNumber 获取下一个目录编号
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

// parseTime 解析时间字符串，使用缓存避免重复解析
func parseTime(timeStr string) (time.Time, error) {
	// 检查缓存中是否已有解析结果
	if t, ok := timeCache[timeStr]; ok {
		return t, nil
	}

	formats := []string{
		"2006-01-02 15:04:05",
		"2006/01/02 15:04:05",
		"20060102150405",
		"2006-01-02T15:04:05",
	}

	for _, format := range formats {
		t, err := time.ParseInLocation(format, timeStr, time.Local)
		if err == nil {
			// 将解析结果存入缓存
			timeCache[timeStr] = t
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time: %s", timeStr)
}
