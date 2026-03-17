package act

import (
	"fmt"
	"time"

	"gobox/tools/act"
)

// Manager 封装和扩展 gobox/tools/act 的业务逻辑
type Manager struct {
	underlying *act.ActManager
}

// NewManager 创建一个新的 Manager 实例
func NewManager(targetPath, backupPath string) *Manager {
	return &Manager{
		underlying: act.NewActManager(targetPath, backupPath),
	}
}

// FilterFilesPreview 过滤文件并返回预览（用于 dry-run）
func (m *Manager) FilterFilesPreview(filter act.Filter) ([]act.DocumentInfo, string, error) {
	return m.underlying.FilterFiles(filter)
}

// Archive 归档操作
func (m *Manager) Archive(beforeTime time.Time, dryRun bool) (string, error) {
	if dryRun {
		// 预览模式：仅列出将要归档的文件
		filter := act.Filter{
			TimeFilter: &act.TimeFilter{
				Type: act.TimeBefore,
				Time: beforeTime,
			},
		}
		docs, _, err := m.underlying.FilterFiles(filter)
		if err != nil {
			return "", err
		}

		// 构建预览输出
		previewMsg := fmt.Sprintf("【预览模式】将归档 %d 个文件到 %s:\n", len(docs), beforeTime.Format("2006-01-02"))
		for i, doc := range docs {
			if i < 5 { // 仅显示前 5 个
				previewMsg += fmt.Sprintf("  - %s (PatientID: %s, TestTime: %s)\n", doc.Filename, doc.PatientID, doc.TestTime.Format("2006-01-02 15:04:05"))
			}
		}
		if len(docs) > 5 {
			previewMsg += fmt.Sprintf("  ... 还有 %d 个文件\n", len(docs)-5)
		}
		return previewMsg, nil
	}

	// 正常模式：执行归档
	return m.underlying.ArchiveFiles(beforeTime)
}

// Resend 重发操作
func (m *Manager) Resend(filter act.Filter, dryRun bool) (string, error) {
	if dryRun {
		// 预览模式
		docs, _, err := m.underlying.FilterFiles(filter)
		if err != nil {
			return "", err
		}

		previewMsg := fmt.Sprintf("【预览模式】将重发 %d 个文件:\n", len(docs))
		for i, doc := range docs {
			if i < 5 {
				previewMsg += fmt.Sprintf("  - %s (PatientID: %s)\n", doc.Filename, doc.PatientID)
			}
		}
		if len(docs) > 5 {
			previewMsg += fmt.Sprintf("  ... 还有 %d 个文件\n", len(docs)-5)
		}
		return previewMsg, nil
	}

	// 正常模式：执行重发
	docs, _, err := m.underlying.FilterFiles(filter)
	if err != nil {
		return "", err
	}
	return m.underlying.ResendFiles(docs)
}

// Modify 修改操作
func (m *Manager) Modify(filter act.Filter, newPatientID string, dryRun bool) (string, error) {
	if dryRun {
		// 预览模式
		docs, _, err := m.underlying.FilterFiles(filter)
		if err != nil {
			return "", err
		}

		previewMsg := fmt.Sprintf("【预览模式】将修改 %d 个文件的 PatientID 为 %s:\n", len(docs), newPatientID)
		for i, doc := range docs {
			if i < 5 {
				previewMsg += fmt.Sprintf("  - %s (%s -> %s)\n", doc.Filename, doc.PatientID, newPatientID)
			}
		}
		if len(docs) > 5 {
			previewMsg += fmt.Sprintf("  ... 还有 %d 个文件\n", len(docs)-5)
		}
		return previewMsg, nil
	}

	// 正常模式：执行修改
	docs, _, err := m.underlying.FilterFiles(filter)
	if err != nil {
		return "", err
	}
	return m.underlying.ModifyFiles(docs, newPatientID)
}

// Rollback 回滚操作
func (m *Manager) Rollback(operation string, index int, dryRun bool) (string, error) {
	if dryRun {
		return fmt.Sprintf("【预览模式】将回滚操作 %s (索引 %d)", operation, index), nil
	}

	// 正常模式：执行回滚
	return m.underlying.Rollback(operation, index)
}

// ListFiles 列出文件
func (m *Manager) ListFiles(filter act.Filter) ([]act.DocumentInfo, string, error) {
	return m.underlying.FilterFiles(filter)
}
