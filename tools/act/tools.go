package act

import (
	"context"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	TargetPath = "d:\\de"
	BackupPath = "d:\\deBackup"
)

type SetConfigInput struct {
	TargetPath string `json:"targetPath" jsonschema:"目标文件路径，要扫描的目录"`
	BackupPath string `json:"backupPath" jsonschema:"备份文件路径，归档文件将存储在此目录"`
}

type SetConfigOutput struct {
	Result string `json:"result" jsonschema:"配置结果信息"`
}

func SetConfig(ctx context.Context, req *mcp.CallToolRequest, input SetConfigInput) (
	*mcp.CallToolResult,
	SetConfigOutput,
	error,
) {
	TargetPath = input.TargetPath
	BackupPath = input.BackupPath
	return nil, SetConfigOutput{Result: fmt.Sprintf("配置已更新: TargetPath=%s, BackupPath=%s", TargetPath, BackupPath)}, nil
}

type ArchiveInput struct {
	BeforeTime string `json:"beforeTime" jsonschema:"归档时间，格式为 2006-01-02 15:04:05，早于此时间的文件将被归档"`
}

type ArchiveOutput struct {
	Result string `json:"result" jsonschema:"归档操作的结果信息"`
}

func Archive(ctx context.Context, req *mcp.CallToolRequest, input ArchiveInput) (
	*mcp.CallToolResult,
	ArchiveOutput,
	error,
) {
	if TargetPath == "" || BackupPath == "" {
		return nil, ArchiveOutput{}, fmt.Errorf("请先调用 SetConfig 设置目标路径和备份路径")
	}

	beforeTime, err := time.Parse("2006-01-02 15:04:05", input.BeforeTime)
	if err != nil {
		return nil, ArchiveOutput{}, fmt.Errorf("时间格式错误，请使用格式 2006-01-02 15:04:05: %w", err)
	}

	manager := NewActManager(TargetPath, BackupPath)
	result, err := manager.ArchiveFiles(beforeTime)
	if err != nil {
		return nil, ArchiveOutput{}, err
	}

	return nil, ArchiveOutput{Result: result}, nil
}

type FilterInput struct {
	PatientID      string `json:"patientID,omitempty" jsonschema:"患者ID筛选条件，为空则不筛选"`
	TimeFilterType string `json:"timeFilterType,omitempty" jsonschema:"时间筛选类型：equal/before/after/range"`
	Time           string `json:"time,omitempty" jsonschema:"时间，格式为 2006-01-02 15:04:05，用于 equal/before/after"`
	StartTime      string `json:"startTime,omitempty" jsonschema:"起始时间，格式为 2006-01-02 15:04:05，用于 range"`
	EndTime        string `json:"endTime,omitempty" jsonschema:"结束时间，格式为 2006-01-02 15:04:05，用于 range"`
}

type FilterOutput struct {
	Documents []DocumentInfo `json:"documents" jsonschema:"符合条件的文档列表"`
	Result    string         `json:"result" jsonschema:"筛选结果信息"`
}

func FilterFiles(ctx context.Context, req *mcp.CallToolRequest, input FilterInput) (
	*mcp.CallToolResult,
	FilterOutput,
	error,
) {
	if TargetPath == "" || BackupPath == "" {
		return nil, FilterOutput{}, fmt.Errorf("请先调用 SetConfig 设置目标路径和备份路径")
	}

	manager := NewActManager(TargetPath, BackupPath)

	var filter Filter
	filter.PatientID = input.PatientID

	if input.TimeFilterType != "" {
		timeFilter := &TimeFilter{}

		switch input.TimeFilterType {
		case "equal":
			t, err := time.Parse("2006-01-02 15:04:05", input.Time)
			if err != nil {
				return nil, FilterOutput{}, fmt.Errorf("时间格式错误: %w", err)
			}
			timeFilter.Type = TimeEqual
			timeFilter.Time = t
		case "before":
			t, err := time.Parse("2006-01-02 15:04:05", input.Time)
			if err != nil {
				return nil, FilterOutput{}, fmt.Errorf("时间格式错误: %w", err)
			}
			timeFilter.Type = TimeBefore
			timeFilter.Time = t
		case "after":
			t, err := time.Parse("2006-01-02 15:04:05", input.Time)
			if err != nil {
				return nil, FilterOutput{}, fmt.Errorf("时间格式错误: %w", err)
			}
			timeFilter.Type = TimeAfter
			timeFilter.Time = t
		case "range":
			startT, err := time.Parse("2006-01-02 15:04:05", input.StartTime)
			if err != nil {
				return nil, FilterOutput{}, fmt.Errorf("起始时间格式错误: %w", err)
			}
			endT, err := time.Parse("2006-01-02 15:04:05", input.EndTime)
			if err != nil {
				return nil, FilterOutput{}, fmt.Errorf("结束时间格式错误: %w", err)
			}
			timeFilter.Type = TimeRange
			timeFilter.StartTime = startT
			timeFilter.EndTime = endT
		default:
			return nil, FilterOutput{}, fmt.Errorf("不支持的时间筛选类型: %s", input.TimeFilterType)
		}

		filter.TimeFilter = timeFilter
	}

	docs, result, err := manager.FilterFiles(filter)
	if err != nil {
		return nil, FilterOutput{}, err
	}

	return nil, FilterOutput{Documents: docs, Result: result}, nil
}

type ResendInput struct {
	Documents []DocumentInfo `json:"documents" jsonschema:"要重传的文档列表"`
}

type ResendOutput struct {
	Result string `json:"result" jsonschema:"重传操作的结果信息"`
}

func Resend(ctx context.Context, req *mcp.CallToolRequest, input ResendInput) (
	*mcp.CallToolResult,
	ResendOutput,
	error,
) {
	if TargetPath == "" || BackupPath == "" {
		return nil, ResendOutput{}, fmt.Errorf("请先调用 SetConfig 设置目标路径和备份路径")
	}

	manager := NewActManager(TargetPath, BackupPath)
	result, err := manager.ResendFiles(input.Documents)
	if err != nil {
		return nil, ResendOutput{}, err
	}

	return nil, ResendOutput{Result: result}, nil
}

type ModifyInput struct {
	Documents    []DocumentInfo `json:"documents" jsonschema:"要修改的文档列表"`
	NewPatientID string         `json:"newPatientID" jsonschema:"新的患者ID"`
}

type ModifyOutput struct {
	Result string `json:"result" jsonschema:"修改操作的结果信息"`
}

func Modify(ctx context.Context, req *mcp.CallToolRequest, input ModifyInput) (
	*mcp.CallToolResult,
	ModifyOutput,
	error,
) {
	if TargetPath == "" || BackupPath == "" {
		return nil, ModifyOutput{}, fmt.Errorf("请先调用 SetConfig 设置目标路径和备份路径")
	}

	manager := NewActManager(TargetPath, BackupPath)
	result, err := manager.ModifyFiles(input.Documents, input.NewPatientID)
	if err != nil {
		return nil, ModifyOutput{}, err
	}

	return nil, ModifyOutput{Result: result}, nil
}

type RollbackInput struct {
	Operation string `json:"operation" jsonschema:"操作类型：resend/modify"`
	Index     int    `json:"index" jsonschema:"操作的目录编号"`
}

type RollbackOutput struct {
	Result string `json:"result" jsonschema:"回滚操作的结果信息"`
}

func Rollback(ctx context.Context, req *mcp.CallToolRequest, input RollbackInput) (
	*mcp.CallToolResult,
	RollbackOutput,
	error,
) {
	if BackupPath == "" {
		return nil, RollbackOutput{}, fmt.Errorf("请先调用 SetConfig 设置备份路径")
	}

	manager := NewActManager("", BackupPath)
	result, err := manager.Rollback(input.Operation, input.Index)
	if err != nil {
		return nil, RollbackOutput{}, err
	}

	return nil, RollbackOutput{Result: result}, nil
}
