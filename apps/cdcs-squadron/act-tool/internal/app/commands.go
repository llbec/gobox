package app

import (
	"time"

	"gobox/tools/act"
)

// ArchiveParams 备份命令的参数
type ArchiveParams struct {
	BeforeTime time.Time
}

// ListParams 列出命令的参数
type ListParams struct {
	PatientID  string
	TimeType   string    // equal, before, after, range
	Time       time.Time // 用于 equal/before/after
	StartTime  time.Time // 用于 range
	EndTime    time.Time // 用于 range
}

// ResendParams 重发命令的参数
type ResendParams struct {
	Filter act.Filter
}

// ModifyParams 修改命令的参数
type ModifyParams struct {
	Filter       act.Filter
	NewPatientID string
}

// RollbackParams 回滚命令的参数
type RollbackParams struct {
	Operation string // resend 或 modify
	Index     int
}

// ConvertListParamsToFilter 将 ListParams 转换为 act.Filter
func ConvertListParamsToFilter(params ListParams) act.Filter {
	filter := act.Filter{
		PatientID: params.PatientID,
	}

	if params.TimeType != "" && params.TimeType != "none" {
		filter.TimeFilter = &act.TimeFilter{}

		switch params.TimeType {
		case "equal":
			filter.TimeFilter.Type = act.TimeEqual
			filter.TimeFilter.Time = params.Time
		case "before":
			filter.TimeFilter.Type = act.TimeBefore
			filter.TimeFilter.Time = params.Time
		case "after":
			filter.TimeFilter.Type = act.TimeAfter
			filter.TimeFilter.Time = params.Time
		case "range":
			filter.TimeFilter.Type = act.TimeRange
			filter.TimeFilter.StartTime = params.StartTime
			filter.TimeFilter.EndTime = params.EndTime
		}
	}

	return filter
}
