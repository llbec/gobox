package app

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"gobox/apps/cdcs-squadron/act-tool/internal/act"
	"gobox/apps/cdcs-squadron/act-tool/internal/config"
)

// App 应用主体
type App struct {
	manager *act.Manager
	config  *config.Config
	menu    *Menu
	confirm *Confirm
}

// NewApp 创建新的应用实例
func NewApp(manager any, cfg *config.Config) *App {
	// 类型转换：从 gobox/tools/act 的 ActManager 转为内部 Manager
	var internalManager *act.Manager
	if m, ok := manager.(*act.Manager); ok {
		internalManager = m
	} else {
		// 如果不是已包装的 Manager，则创建新的
		internalManager = act.NewManager(cfg.TargetPath, cfg.BackupPath)
	}

	return &App{
		manager: internalManager,
		config:  cfg,
		menu:    NewMenu(),
		confirm: NewConfirm(),
	}
}

// RunInteractive 交互式模式运行
func (a *App) RunInteractive(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		// 显示主菜单
		a.menu.ShowMainMenu()

		// 读取用户选择
		if !scanner.Scan() {
			return nil // EOF
		}

		choice := strings.TrimSpace(scanner.Text())

		// 处理特殊命令
		if choice == "q" || choice == "exit" {
			fmt.Println("再见！")
			return nil
		}

		// 分发到各个命令
		switch choice {
		case "1":
			err := a.handleArchive(ctx, scanner)
			if err != nil {
				fmt.Printf("错误: %v\n", err)
			}
		case "2":
			err := a.handleList(ctx, scanner)
			if err != nil {
				fmt.Printf("错误: %v\n", err)
			}
		case "3":
			err := a.handleResend(ctx, scanner)
			if err != nil {
				fmt.Printf("错误: %v\n", err)
			}
		case "4":
			err := a.handleModify(ctx, scanner)
			if err != nil {
				fmt.Printf("错误: %v\n", err)
			}
		case "5":
			err := a.handleRollback(ctx, scanner)
			if err != nil {
				fmt.Printf("错误: %v\n", err)
			}
		case "0":
			fmt.Println("再见！")
			return nil
		default:
			fmt.Println("无效的选择，请重试")
		}

		fmt.Println()
	}
}

// RunBatch 批处理模式运行
func (a *App) RunBatch(ctx context.Context, cmd string) error {
	switch cmd {
	case "archive":
		return fmt.Errorf("批处理模式下的 archive 命令尚未实现")
	case "list":
		return fmt.Errorf("批处理模式下的 list 命令尚未实现")
	case "resend":
		return fmt.Errorf("批处理模式下的 resend 命令尚未实现")
	case "modify":
		return fmt.Errorf("批处理模式下的 modify 命令尚未实现")
	case "rollback":
		return fmt.Errorf("批处理模式下的 rollback 命令尚未实现")
	default:
		return fmt.Errorf("未知的命令: %s", cmd)
	}
}

// 各个命令的处理函数

func (a *App) handleArchive(ctx context.Context, scanner *bufio.Scanner) error {
	menu := &Menu{}
	confirm := &Confirm{}

	fmt.Println("\n【备份操作】")
	fmt.Println("将指定时间之前的文件归档到备份目录")

	// 提示输入时间
	beforeDate, err := menu.PromptDate("请输入归档时间点 (YYYY-MM-DD): ")
	if err != nil {
		confirm.ShowError("时间输入", err)
		return nil
	}

	// 执行预览或正式操作
	result, err := a.manager.Archive(beforeDate, a.config.DryRun)
	if err != nil {
		confirm.ShowError("备份操作", err)
		return nil
	}

	// 显示预览结果
	confirm.ShowPreview("【预览】待归档的文件", result)

	// 如果启用 interactive，要求用户确认
	if a.config.Interactive && !a.config.DryRun {
		confirmed, err := menu.PromptConfirm("确认执行备份操作吗?")
		if err != nil {
			confirm.ShowError("用户确认", err)
			return nil
		}

		if confirmed {
			result, err := a.manager.Archive(beforeDate, false)
			if err != nil {
				confirm.ShowError("备份操作", err)
				return nil
			}
			confirm.ShowSuccess("备份完成", result)
		} else {
			fmt.Println("操作已取消")
		}
	} else if !a.config.DryRun {
		confirm.ShowSuccess("备份完成", result)
	}

	return nil
}

func (a *App) handleList(ctx context.Context, scanner *bufio.Scanner) error {
	menu := &Menu{}
	confirm := &Confirm{}

	fmt.Println("\n【列出文件】")
	fmt.Println("根据筛选条件列出符合要求的文件")

	// 提示输入 PatientID（可选）
	patientID, err := menu.PromptString("请输入 PatientID (可为空按 Enter 跳过): ")
	if err != nil {
		confirm.ShowError("输入 PatientID", err)
		return nil
	}

	// 提示选择时间过滤方式
	fmt.Println("\n时间过滤选项:")
	fmt.Println("  1 - 无时间限制")
	fmt.Println("  2 - 等于某个日期")
	fmt.Println("  3 - 早于某个日期")
	fmt.Println("  4 - 晚于某个日期")
	fmt.Println("  5 - 时间范围内")
	timeChoice, err := menu.PromptString("请选择 [1-5]: ")
	if err != nil {
		confirm.ShowError("选择时间过滤", err)
		return nil
	}

	params := ListParams{
		PatientID: patientID,
	}

	switch timeChoice {
	case "2":
		t, err := menu.PromptDate("请输入日期 (YYYY-MM-DD): ")
		if err != nil {
			confirm.ShowError("输入日期", err)
			return nil
		}
		params.TimeType = "equal"
		params.Time = t
	case "3":
		t, err := menu.PromptDate("请输入日期 (YYYY-MM-DD): ")
		if err != nil {
			confirm.ShowError("输入日期", err)
			return nil
		}
		params.TimeType = "before"
		params.Time = t
	case "4":
		t, err := menu.PromptDate("请输入日期 (YYYY-MM-DD): ")
		if err != nil {
			confirm.ShowError("输入日期", err)
			return nil
		}
		params.TimeType = "after"
		params.Time = t
	case "5":
		start, end, err := menu.PromptTimeRange("请输入时间范围，格式: YYYY-MM-DD YYYY-MM-DD: ")
		if err != nil {
			confirm.ShowError("输入时间范围", err)
			return nil
		}
		params.TimeType = "range"
		params.StartTime = start
		params.EndTime = end
	default:
		params.TimeType = "none"
	}

	// 转换为 Filter
	filter := ConvertListParamsToFilter(params)

	// 执行列出操作
	docs, _, err := a.manager.ListFiles(filter)
	if err != nil {
		confirm.ShowError("列出文件", err)
		return nil
	}

	// 构建输出
	output := fmt.Sprintf("共找到 %d 个符合条件的文件:\n\n", len(docs))
	output += "┌──────────────────────┬──────────────┬─────────────────────────┐\n"
	output += "│ 文件名               │ PatientID    │ 测试时间                │\n"
	output += "├──────────────────────┼──────────────┼─────────────────────────┤\n"

	for i, doc := range docs {
		if i >= 10 { // 仅显示前 10 个
			output += fmt.Sprintf("│ ... 还有 %d 个文件 (省略显示)  │\n", len(docs)-10)
			break
		}
		output += fmt.Sprintf("│ %-20s │ %-12s │ %-23s │\n",
			truncate(doc.Filename, 20),
			doc.PatientID,
			doc.TestTime.Format("2006-01-02 15:04:05"))
	}
	output += "└──────────────────────┴──────────────┴─────────────────────────┘\n"

	confirm.ShowSuccess("列出完成", output)
	return nil
}

func (a *App) handleResend(ctx context.Context, scanner *bufio.Scanner) error {
	menu := &Menu{}
	confirm := &Confirm{}

	fmt.Println("\n【重发操作】")
	fmt.Println("将指定文件重新发送（生成新文件")

	// 提示输入 PatientID（可选）
	patientID, err := menu.PromptString("请输入 PatientID (可为空按 Enter 跳过): ")
	if err != nil {
		confirm.ShowError("输入 PatientID", err)
		return nil
	}

	// 提示选择时间过滤方式
	fmt.Println("\n时间过滤选项:")
	fmt.Println("  1 - 无时间限制")
	fmt.Println("  2 - 等于某个日期")
	fmt.Println("  3 - 早于某个日期")
	fmt.Println("  4 - 晚于某个日期")
	fmt.Println("  5 - 时间范围内")
	timeChoice, err := menu.PromptString("请选择 [1-5]: ")
	if err != nil {
		confirm.ShowError("选择时间过滤", err)
		return nil
	}

	params := ListParams{
		PatientID: patientID,
	}

	switch timeChoice {
	case "2", "3", "4", "5":
		// 与 handleList 类似处理时间参数
		if timeChoice == "2" {
			t, err := menu.PromptDate("请输入日期 (YYYY-MM-DD): ")
			if err != nil {
				confirm.ShowError("输入日期", err)
				return nil
			}
			params.TimeType = "equal"
			params.Time = t
		} else if timeChoice == "3" {
			t, err := menu.PromptDate("请输入日期 (YYYY-MM-DD): ")
			if err != nil {
				confirm.ShowError("输入日期", err)
				return nil
			}
			params.TimeType = "before"
			params.Time = t
		} else if timeChoice == "4" {
			t, err := menu.PromptDate("请输入日期 (YYYY-MM-DD): ")
			if err != nil {
				confirm.ShowError("输入日期", err)
				return nil
			}
			params.TimeType = "after"
			params.Time = t
		} else if timeChoice == "5" {
			start, end, err := menu.PromptTimeRange("请输入时间范围，格式: YYYY-MM-DD YYYY-MM-DD: ")
			if err != nil {
				confirm.ShowError("输入时间范围", err)
				return nil
			}
			params.TimeType = "range"
			params.StartTime = start
			params.EndTime = end
		}
	default:
		params.TimeType = "none"
	}

	filter := ConvertListParamsToFilter(params)

	// 执行预览
	result, err := a.manager.Resend(filter, true) // 先显示预览
	if err != nil {
		confirm.ShowError("重发预览", err)
		return nil
	}

	confirm.ShowPreview("【预览】待重发的文件", result)

	// 如果启用 interactive 或 dry-run，要求用户确认
	if a.config.Interactive || a.config.DryRun {
		if a.config.DryRun {
			fmt.Println("（DRY-RUN 模式，不执行实际操作）")
			return nil
		}

		confirmed, err := menu.PromptConfirm("确认执行重发操作吗?")
		if err != nil {
			confirm.ShowError("用户确认", err)
			return nil
		}

		if confirmed {
			result, err := a.manager.Resend(filter, false)
			if err != nil {
				confirm.ShowError("重发操作", err)
				return nil
			}
			confirm.ShowSuccess("重发完成", result)
		} else {
			fmt.Println("操作已取消")
		}
	} else {
		// 直接执行
		result, err := a.manager.Resend(filter, false)
		if err != nil {
			confirm.ShowError("重发操作", err)
			return nil
		}
		confirm.ShowSuccess("重发完成", result)
	}

	return nil
}

func (a *App) handleModify(ctx context.Context, scanner *bufio.Scanner) error {
	menu := &Menu{}
	confirm := &Confirm{}

	fmt.Println("\n【修改操作】")
	fmt.Println("修改指定文件中的 PatientID 并重新发送")

	// 提示输入新 PatientID
	newPatientID, err := menu.PromptString("请输入新的 PatientID: ")
	if err != nil {
		confirm.ShowError("输入新 PatientID", err)
		return nil
	}

	if strings.TrimSpace(newPatientID) == "" {
		fmt.Println("错误: PatientID 不能为空")
		return nil
	}

	// 提示输入原 PatientID（可选）
	oldPatientID, err := menu.PromptString("请输入原 PatientID (可为空按 Enter 跳过): ")
	if err != nil {
		confirm.ShowError("输入原 PatientID", err)
		return nil
	}

	// 提示选择时间过滤方式
	fmt.Println("\n时间过滤选项:")
	fmt.Println("  1 - 无时间限制")
	fmt.Println("  2 - 等于某个日期")
	fmt.Println("  3 - 早于某个日期")
	fmt.Println("  4 - 晚于某个日期")
	fmt.Println("  5 - 时间范围内")
	timeChoice, err := menu.PromptString("请选择 [1-5]: ")
	if err != nil {
		confirm.ShowError("选择时间过滤", err)
		return nil
	}

	params := ListParams{
		PatientID: oldPatientID,
	}

	switch timeChoice {
	case "2":
		t, err := menu.PromptDate("请输入日期 (YYYY-MM-DD): ")
		if err != nil {
			confirm.ShowError("输入日期", err)
			return nil
		}
		params.TimeType = "equal"
		params.Time = t
	case "3":
		t, err := menu.PromptDate("请输入日期 (YYYY-MM-DD): ")
		if err != nil {
			confirm.ShowError("输入日期", err)
			return nil
		}
		params.TimeType = "before"
		params.Time = t
	case "4":
		t, err := menu.PromptDate("请输入日期 (YYYY-MM-DD): ")
		if err != nil {
			confirm.ShowError("输入日期", err)
			return nil
		}
		params.TimeType = "after"
		params.Time = t
	case "5":
		start, end, err := menu.PromptTimeRange("请输入时间范围，格式: YYYY-MM-DD YYYY-MM-DD: ")
		if err != nil {
			confirm.ShowError("输入时间范围", err)
			return nil
		}
		params.TimeType = "range"
		params.StartTime = start
		params.EndTime = end
	default:
		params.TimeType = "none"
	}

	filter := ConvertListParamsToFilter(params)

	// 执行预览
	result, err := a.manager.Modify(filter, newPatientID, true) // 先显示预览
	if err != nil {
		confirm.ShowError("修改预览", err)
		return nil
	}

	confirm.ShowPreview("【预览】待修改的文件", result)

	// 如果启用 interactive 或 dry-run，要求用户确认
	if a.config.Interactive || a.config.DryRun {
		if a.config.DryRun {
			fmt.Println("（DRY-RUN 模式，不执行实际操作）")
			return nil
		}

		confirmed, err := menu.PromptConfirm("确认执行修改操作吗?")
		if err != nil {
			confirm.ShowError("用户确认", err)
			return nil
		}

		if confirmed {
			result, err := a.manager.Modify(filter, newPatientID, false)
			if err != nil {
				confirm.ShowError("修改操作", err)
				return nil
			}
			confirm.ShowSuccess("修改完成", result)
		} else {
			fmt.Println("操作已取消")
		}
	} else {
		// 直接执行
		result, err := a.manager.Modify(filter, newPatientID, false)
		if err != nil {
			confirm.ShowError("修改操作", err)
			return nil
		}
		confirm.ShowSuccess("修改完成", result)
	}

	return nil
}

func (a *App) handleRollback(ctx context.Context, scanner *bufio.Scanner) error {
	menu := &Menu{}
	confirm := &Confirm{}

	fmt.Println("\n【回滚操作】")
	fmt.Println("回滚之前的重发或修改操作")

	// 提示选择操作类型
	fmt.Println("操作类型:")
	fmt.Println("  1 - 回滚重发操作 (resend)")
	fmt.Println("  2 - 回滚修改操作 (modify)")
	opChoice, err := menu.PromptString("请选择 [1-2]: ")
	if err != nil {
		confirm.ShowError("选择操作类型", err)
		return nil
	}

	var operation string
	switch opChoice {
	case "1":
		operation = "resend"
	case "2":
		operation = "modify"
	default:
		fmt.Println("无效的选择")
		return nil
	}

	// 提示输入操作索引
	index, err := menu.PromptInt("请输入操作索引 (整数): ")
	if err != nil {
		confirm.ShowError("输入索引", err)
		return nil
	}

	// 执行预览
	result, err := a.manager.Rollback(operation, index, true)
	if err != nil {
		confirm.ShowError("回滚预览", err)
		return nil
	}

	confirm.ShowPreview("【预览】待回滚的操作", result)

	// 如果启用 interactive 或 dry-run，要求用户确认
	if a.config.Interactive || a.config.DryRun {
		if a.config.DryRun {
			fmt.Println("（DRY-RUN 模式，不执行实际操作）")
			return nil
		}

		confirmed, err := menu.PromptConfirm("确认执行回滚操作吗?")
		if err != nil {
			confirm.ShowError("用户确认", err)
			return nil
		}

		if confirmed {
			result, err := a.manager.Rollback(operation, index, false)
			if err != nil {
				confirm.ShowError("回滚操作", err)
				return nil
			}
			confirm.ShowSuccess("回滚完成", result)
		} else {
			fmt.Println("操作已取消")
		}
	} else {
		// 直接执行
		result, err := a.manager.Rollback(operation, index, false)
		if err != nil {
			confirm.ShowError("回滚操作", err)
			return nil
		}
		confirm.ShowSuccess("回滚完成", result)
	}

	return nil
}

// 辅助函数：截断字符串
func truncate(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}
