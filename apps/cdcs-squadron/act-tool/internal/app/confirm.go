package app

import (
	"fmt"
)

// Confirm 交互确认管理
type Confirm struct{}

// NewConfirm 创建 Confirm 实例
func NewConfirm() *Confirm {
	return &Confirm{}
}

// ShowPreview 显示待执行的操作预览
func (c *Confirm) ShowPreview(title string, message string) {
	fmt.Println()
	fmt.Println("┌─────────────────────────────────────┐")
	fmt.Printf("│  %s\n", title)
	fmt.Println("├─────────────────────────────────────┤")
	fmt.Println(message)
	fmt.Println("└─────────────────────────────────────┘")
	fmt.Println()
}

// ConfirmExecution 要求用户确认是否执行
func (c *Confirm) ConfirmExecution(operation string) (bool, error) {
	menu := &Menu{}
	prompt := fmt.Sprintf("确认执行 %s 操作吗?", operation)
	return menu.PromptConfirm(prompt)
}

// ShowSummary 显示操作总结
func (c *Confirm) ShowSummary(title string, count int, details string) {
	fmt.Println()
	fmt.Println("┌─────────────────────────────────────┐")
	fmt.Printf("│  操作总结: %s\n", title)
	fmt.Println("├─────────────────────────────────────┤")
	fmt.Printf("│  处理文件数: %d\n", count)
	fmt.Println("├─────────────────────────────────────┤")
	fmt.Println(details)
	fmt.Println("└─────────────────────────────────────┘")
	fmt.Println()
}

// ShowError 显示错误信息
func (c *Confirm) ShowError(title string, err error) {
	fmt.Println()
	fmt.Println("┌─────────────────────────────────────┐")
	fmt.Printf("│  ❌ 错误: %s\n", title)
	fmt.Println("├─────────────────────────────────────┤")
	fmt.Printf("│  %v\n", err)
	fmt.Println("└─────────────────────────────────────┘")
	fmt.Println()
}

// ShowSuccess 显示成功信息
func (c *Confirm) ShowSuccess(title string, message string) {
	fmt.Println()
	fmt.Println("┌─────────────────────────────────────┐")
	fmt.Printf("│  ✓ 成功: %s\n", title)
	fmt.Println("├─────────────────────────────────────┤")
	fmt.Println(message)
	fmt.Println("└─────────────────────────────────────┘")
	fmt.Println()
}
