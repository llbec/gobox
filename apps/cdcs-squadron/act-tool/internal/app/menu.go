package app

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Menu 菜单管理
type Menu struct{}

// NewMenu 创建菜单实例
func NewMenu() *Menu {
	return &Menu{}
}

// ShowMainMenu 显示主菜单
func (m *Menu) ShowMainMenu() {
	fmt.Println("┌─────────────────────────────────────┐")
	fmt.Println("│     act-tool (交互式模式)           │")
	fmt.Println("├─────────────────────────────────────┤")
	fmt.Println("│  1) 备份 (archive)                  │")
	fmt.Println("│  2) 列出 (list)                     │")
	fmt.Println("│  3) 重发 (resend)                   │")
	fmt.Println("│  4) 修改 (modify)                   │")
	fmt.Println("│  5) 回滚 (rollback)                 │")
	fmt.Println("│  0) 退出                            │")
	fmt.Println("└─────────────────────────────────────┘")
	fmt.Print("请选择操作 [0-5]: ")
}

// PromptString 提示用户输入字符串
func (m *Menu) PromptString(message string) (string, error) {
	fmt.Print(message)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text()), nil
	}
	return "", scanner.Err()
}

// PromptInt 提示用户输入整数
func (m *Menu) PromptInt(message string) (int, error) {
	fmt.Print(message)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		val, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
		return val, err
	}
	return 0, scanner.Err()
}

// PromptDate 提示用户输入日期 (YYYY-MM-DD)
func (m *Menu) PromptDate(message string) (time.Time, error) {
	fmt.Print(message)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		dateStr := strings.TrimSpace(scanner.Text())
		return time.ParseInLocation("2006-01-02", dateStr, time.Local)
	}
	return time.Time{}, scanner.Err()
}

// PromptTimeRange 提示用户输入时间范围
func (m *Menu) PromptTimeRange(message string) (time.Time, time.Time, error) {
	fmt.Print(message)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		parts := strings.Fields(strings.TrimSpace(scanner.Text()))
		if len(parts) != 2 {
			return time.Time{}, time.Time{}, fmt.Errorf("请输入两个日期，格式: YYYY-MM-DD YYYY-MM-DD")
		}

		start, err := time.ParseInLocation("2006-01-02", parts[0], time.Local)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}

		end, err := time.ParseInLocation("2006-01-02", parts[1], time.Local)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}

		return start, end, nil
	}
	return time.Time{}, time.Time{}, scanner.Err()
}

// PromptConfirm 提示用户确认（Y/N）
func (m *Menu) PromptConfirm(message string) (bool, error) {
	fmt.Print(message + " (Y/n): ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input := strings.ToLower(strings.TrimSpace(scanner.Text()))
		return input == "y" || input == "yes" || input == "", nil
	}
	return false, scanner.Err()
}
