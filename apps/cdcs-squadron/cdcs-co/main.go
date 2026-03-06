package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

// 与服务器相同的输入输出结构
type SetConfigInput struct {
	TargetPath string `json:"targetPath"`
	BackupPath string `json:"backupPath"`
}

type SetConfigOutput struct {
	Result string `json:"result"`
}

type ArchiveInput struct {
	BeforeTime string `json:"beforeTime"`
}

type ArchiveOutput struct {
	Result string `json:"result"`
}

type DocumentInfo struct {
	FileName  string    `json:"fileName"`
	PatientID string    `json:"patientID"`
	TestTime  time.Time `json:"testTime"`
}

type FilterInput struct {
	PatientID      string `json:"patientID,omitempty"`
	TimeFilterType string `json:"timeFilterType,omitempty"`
	Time           string `json:"time,omitempty"`
	StartTime      string `json:"startTime,omitempty"`
	EndTime        string `json:"endTime,omitempty"`
}

type FilterOutput struct {
	Documents []DocumentInfo `json:"documents"`
	Result    string         `json:"result"`
}

type ResendInput struct {
	Documents []DocumentInfo `json:"documents"`
}

type ResendOutput struct {
	Result string `json:"result"`
}

type ModifyInput struct {
	Documents    []DocumentInfo `json:"documents"`
	NewPatientID string         `json:"newPatientID"`
}

type ModifyOutput struct {
	Result string `json:"result"`
}

type RollbackInput struct {
	Operation string `json:"operation"`
	Index     int    `json:"index"`
}

type RollbackOutput struct {
	Result string `json:"result"`
}

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Uptime    string    `json:"uptime"`
	GoVersion string    `json:"go_version"`
	Hostname  string    `json:"hostname"`
}

// MCP请求结构
type MCPRequest struct {
	Tool  string      `json:"tool"`
	Input interface{} `json:"input"`
}

// MCP响应结构
type MCPResponse struct {
	Output interface{} `json:"output"`
	Error  string      `json:"error,omitempty"`
}

// 服务器配置
type ServerConfig struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Enabled bool   `json:"enabled"`
}

// 应用配置
type AppConfig struct {
	Servers       []ServerConfig `json:"servers"`
	HTTPPort      int            `json:"httpPort"`
	CheckInterval int            `json:"checkInterval"` // 健康检查间隔（秒）
}

// 配置文件路径
const configFile = "config.json"

// 全局配置
var appConfig AppConfig

// 加载配置
func loadConfig() error {
	// 检查配置文件是否存在
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// 创建默认配置
		appConfig = AppConfig{
			Servers: []ServerConfig{
				{
					Name:    "default",
					URL:     "http://localhost:8080",
					Enabled: true,
				},
			},
			HTTPPort:      8081,
			CheckInterval: 30,
		}
		// 保存默认配置
		return saveConfig()
	}

	// 读取配置文件
	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析配置
	if err := json.Unmarshal(data, &appConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 确保至少有一个默认服务器
	if len(appConfig.Servers) == 0 {
		appConfig.Servers = append(appConfig.Servers, ServerConfig{
			Name:    "default",
			URL:     "http://localhost:8080",
			Enabled: true,
		})
		return saveConfig()
	}

	return nil
}

// 保存配置
func saveConfig() error {
	data, err := json.MarshalIndent(appConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

// 获取服务器配置
func getServerConfig(name string) (ServerConfig, error) {
	for _, server := range appConfig.Servers {
		if server.Name == name {
			return server, nil
		}
	}

	// 如果找不到指定服务器，返回默认服务器
	if len(appConfig.Servers) > 0 {
		return appConfig.Servers[0], nil
	}

	return ServerConfig{}, fmt.Errorf("没有配置服务器")
}

// 添加服务器
func addServer(name, url string) error {
	// 检查服务器是否已存在
	for _, server := range appConfig.Servers {
		if server.Name == name {
			return fmt.Errorf("服务器名称已存在")
		}
	}

	// 添加新服务器
	appConfig.Servers = append(appConfig.Servers, ServerConfig{
		Name:    name,
		URL:     url,
		Enabled: true,
	})

	return saveConfig()
}

// 删除服务器
func deleteServer(name string) error {
	if len(appConfig.Servers) <= 1 {
		return fmt.Errorf("至少需要保留一个服务器")
	}

	newServers := []ServerConfig{}
	for _, server := range appConfig.Servers {
		if server.Name != name {
			newServers = append(newServers, server)
		}
	}

	if len(newServers) == len(appConfig.Servers) {
		return fmt.Errorf("服务器不存在")
	}

	appConfig.Servers = newServers
	return saveConfig()
}

// 启用/禁用服务器
func toggleServer(name string, enabled bool) error {
	for i, server := range appConfig.Servers {
		if server.Name == name {
			appConfig.Servers[i].Enabled = enabled
			return saveConfig()
		}
	}

	return fmt.Errorf("服务器不存在")
}

// 发送MCP请求
func sendMCPRequest(tool string, input any, output any, serverName string) error {
	mcpRequest := MCPRequest{
		Tool:  tool,
		Input: input,
	}

	data, err := json.Marshal(mcpRequest)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %w", err)
	}

	// 获取服务器配置
	server, err := getServerConfig(serverName)
	if err != nil {
		return err
	}

	// 构建MCP URL
	mcpURL := server.URL + "/mcp"

	resp, err := http.Post(mcpURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	var mcpResponse MCPResponse
	mcpResponse.Output = output

	if err := json.NewDecoder(resp.Body).Decode(&mcpResponse); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	if mcpResponse.Error != "" {
		return fmt.Errorf("MCP错误: %s", mcpResponse.Error)
	}

	return nil
}

// 健康检查
func checkHealth(serverName string) (HealthResponse, error) {
	// 获取服务器配置
	server, err := getServerConfig(serverName)
	if err != nil {
		return HealthResponse{}, err
	}

	// 构建健康检查URL
	healthURL := server.URL + "/health"

	resp, err := http.Get(healthURL)
	if err != nil {
		return HealthResponse{}, fmt.Errorf("健康检查失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return HealthResponse{}, fmt.Errorf("健康检查返回错误状态码: %d", resp.StatusCode)
	}

	var healthResp HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&healthResp); err != nil {
		return HealthResponse{}, fmt.Errorf("解析健康检查响应失败: %w", err)
	}

	return healthResp, nil
}

// 打印健康检查结果
func printHealthCheckResult(serverName string, healthResp HealthResponse, err error) {
	fmt.Printf("服务器 %s 健康检查结果:\n", serverName)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}
	fmt.Printf("状态: %s\n", healthResp.Status)
	fmt.Printf("时间戳: %s\n", healthResp.Timestamp)
	fmt.Printf("运行时间: %s\n", healthResp.Uptime)
	fmt.Printf("Go版本: %s\n", healthResp.GoVersion)
	fmt.Printf("主机名: %s\n", healthResp.Hostname)
}

// 设置配置
func setConfig(serverName string) error {
	var input SetConfigInput
	fmt.Print("请输入目标路径 (默认: d:\\de): ")
	fmt.Scanln(&input.TargetPath)
	if input.TargetPath == "" {
		input.TargetPath = "d:\\de"
	}

	fmt.Print("请输入备份路径 (默认: d:\\deBackup): ")
	fmt.Scanln(&input.BackupPath)
	if input.BackupPath == "" {
		input.BackupPath = "d:\\deBackup"
	}

	var output SetConfigOutput
	err := sendMCPRequest("set_config", input, &output, serverName)
	if err != nil {
		return fmt.Errorf("设置配置失败: %w", err)
	}

	fmt.Printf("设置配置成功: %s\n", output.Result)
	return nil
}

// 归档文件
func archiveFiles(serverName string) error {
	var input ArchiveInput
	fmt.Print("请输入归档时间 (格式: 2006-01-02 15:04:05): ")
	fmt.Scanln(&input.BeforeTime)

	var output ArchiveOutput
	err := sendMCPRequest("archive_files", input, &output, serverName)
	if err != nil {
		return fmt.Errorf("归档文件失败: %w", err)
	}

	fmt.Printf("归档文件成功: %s\n", output.Result)
	return nil
}

// 筛选文件
func filterFiles(serverName string) error {
	var input FilterInput
	fmt.Print("请输入患者ID (可选，留空则不筛选): ")
	fmt.Scanln(&input.PatientID)

	fmt.Print("请输入时间筛选类型 (equal/before/after/range，可选，留空则不筛选): ")
	fmt.Scanln(&input.TimeFilterType)

	switch input.TimeFilterType {
	case "equal", "before", "after":
		fmt.Print("请输入时间 (格式: 2006-01-02 15:04:05): ")
		fmt.Scanln(&input.Time)
	case "range":
		fmt.Print("请输入起始时间 (格式: 2006-01-02 15:04:05): ")
		fmt.Scanln(&input.StartTime)
		fmt.Print("请输入结束时间 (格式: 2006-01-02 15:04:05): ")
		fmt.Scanln(&input.EndTime)
	}

	var output FilterOutput
	err := sendMCPRequest("filter_files", input, &output, serverName)
	if err != nil {
		return fmt.Errorf("筛选文件失败: %w", err)
	}

	fmt.Printf("筛选结果: %s\n", output.Result)
	fmt.Println("符合条件的文档:")
	for i, doc := range output.Documents {
		fmt.Printf("%d. 文件名: %s, 患者ID: %s, 测试时间: %s\n", i+1, doc.FileName, doc.PatientID, doc.TestTime)
	}

	return nil
}

// 重传文件
func resendFiles(serverName string) error {
	// 先筛选文件获取文档列表
	var filterInput FilterInput
	fmt.Print("请输入患者ID (可选，留空则不筛选): ")
	fmt.Scanln(&filterInput.PatientID)

	var filterOutput FilterOutput
	err := sendMCPRequest("filter_files", filterInput, &filterOutput, serverName)
	if err != nil {
		return fmt.Errorf("筛选文件失败: %w", err)
	}

	if len(filterOutput.Documents) == 0 {
		fmt.Println("没有符合条件的文档")
		return nil
	}

	fmt.Println("可用的文档:")
	for i, doc := range filterOutput.Documents {
		fmt.Printf("%d. 文件名: %s, 患者ID: %s, 测试时间: %s\n", i+1, doc.FileName, doc.PatientID, doc.TestTime)
	}

	fmt.Print("请输入要重传的文档编号 (多个编号用空格分隔): ")
	var indices []int
	for {
		var idx int
		_, err := fmt.Scan(&idx)
		if err != nil {
			break
		}
		indices = append(indices, idx)
	}

	// 清空输入缓冲区
	var c byte
	for {
		_, err = fmt.Scanf("%c", &c)
		if err != nil || c == '\n' {
			break
		}
	}

	if len(indices) == 0 {
		fmt.Println("未选择任何文档")
		return nil
	}

	// 构建要重传的文档列表
	var documents []DocumentInfo
	for _, idx := range indices {
		if idx >= 1 && idx <= len(filterOutput.Documents) {
			documents = append(documents, filterOutput.Documents[idx-1])
		}
	}

	if len(documents) == 0 {
		fmt.Println("没有有效的文档编号")
		return nil
	}

	input := ResendInput{Documents: documents}
	var output ResendOutput
	err = sendMCPRequest("resend_files", input, &output, serverName)
	if err != nil {
		return fmt.Errorf("重传文件失败: %w", err)
	}

	fmt.Printf("重传文件成功: %s\n", output.Result)
	return nil
}

// 修改文件
func modifyFiles(serverName string) error {
	// 先筛选文件获取文档列表
	var filterInput FilterInput
	fmt.Print("请输入患者ID (可选，留空则不筛选): ")
	fmt.Scanln(&filterInput.PatientID)

	var filterOutput FilterOutput
	err := sendMCPRequest("filter_files", filterInput, &filterOutput, serverName)
	if err != nil {
		return fmt.Errorf("筛选文件失败: %w", err)
	}

	if len(filterOutput.Documents) == 0 {
		fmt.Println("没有符合条件的文档")
		return nil
	}

	fmt.Println("可用的文档:")
	for i, doc := range filterOutput.Documents {
		fmt.Printf("%d. 文件名: %s, 患者ID: %s, 测试时间: %s\n", i+1, doc.FileName, doc.PatientID, doc.TestTime)
	}

	fmt.Print("请输入要修改的文档编号 (多个编号用空格分隔): ")
	var indices []int
	for {
		var idx int
		_, err = fmt.Scan(&idx)
		if err != nil {
			break
		}
		indices = append(indices, idx)
	}

	// 清空输入缓冲区
	var c byte
	for {
		_, err = fmt.Scanf("%c", &c)
		if err != nil || c == '\n' {
			break
		}
	}

	if len(indices) == 0 {
		fmt.Println("未选择任何文档")
		return nil
	}

	// 构建要修改的文档列表
	var documents []DocumentInfo
	for _, idx := range indices {
		if idx >= 1 && idx <= len(filterOutput.Documents) {
			documents = append(documents, filterOutput.Documents[idx-1])
		}
	}

	if len(documents) == 0 {
		fmt.Println("没有有效的文档编号")
		return nil
	}

	var input ModifyInput
	input.Documents = documents
	fmt.Print("请输入新的患者ID: ")
	fmt.Scanln(&input.NewPatientID)

	var output ModifyOutput
	err = sendMCPRequest("modify_files", input, &output, serverName)
	if err != nil {
		return fmt.Errorf("修改文件失败: %w", err)
	}

	fmt.Printf("修改文件成功: %s\n", output.Result)
	return nil
}

// 回滚操作
func rollbackOperation(serverName string) error {
	var input RollbackInput
	fmt.Print("请输入操作类型 (resend/modify): ")
	fmt.Scanln(&input.Operation)

	fmt.Print("请输入操作的目录编号: ")
	fmt.Scanln(&input.Index)

	var output RollbackOutput
	err := sendMCPRequest("rollback_operation", input, &output, serverName)
	if err != nil {
		return fmt.Errorf("回滚操作失败: %w", err)
	}

	fmt.Printf("回滚操作成功: %s\n", output.Result)
	return nil
}

// 显示菜单
func showMenu() {
	fmt.Println("=====================================")
	fmt.Println("CDCS Squadron 客户端")
	fmt.Println("=====================================")
	fmt.Println("1. 健康检查")
	fmt.Println("2. 设置配置")
	fmt.Println("3. 归档文件")
	fmt.Println("4. 筛选文件")
	fmt.Println("5. 重传文件")
	fmt.Println("6. 修改文件")
	fmt.Println("7. 回滚操作")
	fmt.Println("8. 服务器管理")
	fmt.Println("9. 启动HTTP服务")
	fmt.Println("0. 退出")
	fmt.Println("=====================================")
	fmt.Print("请选择操作: ")
}

// 显示服务器管理菜单
func showServerMenu() {
	fmt.Println("=====================================")
	fmt.Println("服务器管理")
	fmt.Println("=====================================")
	fmt.Println("1. 查看服务器列表")
	fmt.Println("2. 添加服务器")
	fmt.Println("3. 删除服务器")
	fmt.Println("4. 启用/禁用服务器")
	fmt.Println("0. 返回主菜单")
	fmt.Println("=====================================")
	fmt.Print("请选择操作: ")
}

// 服务器管理
func serverManagement() error {
	for {
		showServerMenu()
		var choice int
		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("输入错误，请重新输入")
			// 清空输入缓冲区
			var c byte
			for {
				_, err := fmt.Scanf("%c", &c)
				if err != nil || c == '\n' {
					break
				}
			}
			continue
		}

		switch choice {
		case 0:
			return nil
		case 1:
			// 查看服务器列表
			fmt.Println("服务器列表:")
			for i, server := range appConfig.Servers {
				status := "启用"
				if !server.Enabled {
					status = "禁用"
				}
				fmt.Printf("%d. %s - %s (%s)\n", i+1, server.Name, server.URL, status)
			}
		case 2:
			// 添加服务器
			var name, url string
			fmt.Print("请输入服务器名称: ")
			fmt.Scanln(&name)
			fmt.Print("请输入服务器URL: ")
			fmt.Scanln(&url)
			err := addServer(name, url)
			if err != nil {
				fmt.Printf("错误: %v\n", err)
			} else {
				fmt.Println("服务器添加成功")
			}
		case 3:
			// 删除服务器
			var name string
			fmt.Print("请输入要删除的服务器名称: ")
			fmt.Scanln(&name)
			err := deleteServer(name)
			if err != nil {
				fmt.Printf("错误: %v\n", err)
			} else {
				fmt.Println("服务器删除成功")
			}
		case 4:
			// 启用/禁用服务器
			var name string
			var enabled bool
			fmt.Print("请输入服务器名称: ")
			fmt.Scanln(&name)
			fmt.Print("请输入状态 (true=启用, false=禁用): ")
			fmt.Scanln(&enabled)
			err := toggleServer(name, enabled)
			if err != nil {
				fmt.Printf("错误: %v\n", err)
			} else {
				fmt.Println("服务器状态更新成功")
			}
		default:
			fmt.Println("无效的选择，请重新输入")
		}

		fmt.Println("按回车键继续...")
		// 清空输入缓冲区
		var c byte
		for {
			_, err := fmt.Scanf("%c", &c)
			if err != nil || c == '\n' {
				break
			}
		}
	}
}

// 启动HTTP服务
func startHTTPServer() error {
	// 定义HTTP路由
	http.HandleFunc("/health", httpHealthHandler)
	http.HandleFunc("/api/health", httpApiHealthHandler)
	http.HandleFunc("/api/servers", httpApiServersHandler)
	http.HandleFunc("/api/config", httpApiConfigHandler)
	http.HandleFunc("/api/archive", httpApiArchiveHandler)
	http.HandleFunc("/api/filter", httpApiFilterHandler)
	http.HandleFunc("/api/resend", httpApiResendHandler)
	http.HandleFunc("/api/modify", httpApiModifyHandler)
	http.HandleFunc("/api/rollback", httpApiRollbackHandler)

	// 静态文件服务
	http.Handle("/", http.FileServer(http.Dir("./ui")))

	// 启动健康检查轮询
	//go startHealthCheckPolling()

	// 启动HTTP服务器
	port := appConfig.HTTPPort
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("HTTP服务启动在端口 %d...\n", port)
	fmt.Printf("健康检查端点: http://localhost:%d/health\n", port)
	fmt.Printf("API端点: http://localhost:%d/api\n", port)
	fmt.Printf("前端页面: http://localhost:%d/\n", port)

	return http.ListenAndServe(addr, nil)
}

// HTTP健康检查处理
func httpHealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 检查所有启用的服务器
	results := make(map[string]any)
	for _, server := range appConfig.Servers {
		if server.Enabled {
			healthResp, err := checkHealth(server.Name)
			if err != nil {
				results[server.Name] = map[string]string{"status": "error", "message": err.Error()}
			} else {
				results[server.Name] = healthResp
			}
		}
	}

	response := map[string]any{
		"status":    "healthy",
		"timestamp": time.Now(),
		"servers":   results,
	}

	json.NewEncoder(w).Encode(response)
}

// API健康检查处理
func httpApiHealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	serverName := r.URL.Query().Get("server")
	healthResp, err := checkHealth(serverName)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(healthResp)
}

// API服务器列表处理
func httpApiServersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appConfig.Servers)
}

// API配置处理
func httpApiConfigHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {
		var input SetConfigInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "无效的请求数据"})
			return
		}

		serverName := r.URL.Query().Get("server")
		var output SetConfigOutput
		err := sendMCPRequest("set_config", input, &output, serverName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		json.NewEncoder(w).Encode(output)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "只支持POST方法"})
	}
}

// API归档处理
func httpApiArchiveHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {
		var input ArchiveInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "无效的请求数据"})
			return
		}

		serverName := r.URL.Query().Get("server")
		var output ArchiveOutput
		err := sendMCPRequest("archive_files", input, &output, serverName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		json.NewEncoder(w).Encode(output)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "只支持POST方法"})
	}
}

// API筛选处理
func httpApiFilterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {
		var input FilterInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "无效的请求数据"})
			return
		}

		serverName := r.URL.Query().Get("server")
		var output FilterOutput
		err := sendMCPRequest("filter_files", input, &output, serverName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		json.NewEncoder(w).Encode(output)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "只支持POST方法"})
	}
}

// API重传处理
func httpApiResendHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {
		var input ResendInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "无效的请求数据"})
			return
		}

		serverName := r.URL.Query().Get("server")
		var output ResendOutput
		err := sendMCPRequest("resend_files", input, &output, serverName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		json.NewEncoder(w).Encode(output)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "只支持POST方法"})
	}
}

// API修改处理
func httpApiModifyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {
		var input ModifyInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "无效的请求数据"})
			return
		}

		serverName := r.URL.Query().Get("server")
		var output ModifyOutput
		err := sendMCPRequest("modify_files", input, &output, serverName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		json.NewEncoder(w).Encode(output)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "只支持POST方法"})
	}
}

// API回滚处理
func httpApiRollbackHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {
		var input RollbackInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "无效的请求数据"})
			return
		}

		serverName := r.URL.Query().Get("server")
		var output RollbackOutput
		err := sendMCPRequest("rollback_operation", input, &output, serverName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		json.NewEncoder(w).Encode(output)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "只支持POST方法"})
	}
}

// 健康检查轮询
func startHealthCheckPolling() {
	ticker := time.NewTicker(time.Duration(appConfig.CheckInterval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		var wg sync.WaitGroup

		for _, server := range appConfig.Servers {
			if !server.Enabled {
				continue
			}

			wg.Add(1)

			go func(s ServerConfig) {
				defer wg.Done()

				_, err := checkHealth(s.Name)
				if err != nil {
					fmt.Printf("服务器 %s 健康检查失败: %v\n", s.Name, err)
					return
				}

				fmt.Printf("服务器 %s 健康检查成功\n", s.Name)
			}(server)
		}

		wg.Wait()
	}
}

func main() {
	// 加载配置
	if err := loadConfig(); err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		return
	}

	for {
		showMenu()
		var choice int
		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("输入错误，请重新输入")
			// 清空输入缓冲区
			var c byte
			for {
				_, err := fmt.Scanf("%c", &c)
				if err != nil || c == '\n' {
					break
				}
			}
			continue
		}

		switch choice {
		case 0:
			fmt.Println("退出客户端...")
			return
		case 1:
			// 健康检查
			var serverName string
			fmt.Print("请输入服务器名称 (留空使用默认服务器): ")
			fmt.Scanln(&serverName)
			healthResp, err := checkHealth(serverName)
			printHealthCheckResult(serverName, healthResp, err)
		case 2:
			// 设置配置
			var serverName string
			fmt.Print("请输入服务器名称 (留空使用默认服务器): ")
			fmt.Scanln(&serverName)
			err := setConfig(serverName)
			if err != nil {
				fmt.Printf("错误: %v\n", err)
			}
		case 3:
			// 归档文件
			var serverName string
			fmt.Print("请输入服务器名称 (留空使用默认服务器): ")
			fmt.Scanln(&serverName)
			err := archiveFiles(serverName)
			if err != nil {
				fmt.Printf("错误: %v\n", err)
			}
		case 4:
			// 筛选文件
			var serverName string
			fmt.Print("请输入服务器名称 (留空使用默认服务器): ")
			fmt.Scanln(&serverName)
			err := filterFiles(serverName)
			if err != nil {
				fmt.Printf("错误: %v\n", err)
			}
		case 5:
			// 重传文件
			var serverName string
			fmt.Print("请输入服务器名称 (留空使用默认服务器): ")
			fmt.Scanln(&serverName)
			err := resendFiles(serverName)
			if err != nil {
				fmt.Printf("错误: %v\n", err)
			}
		case 6:
			// 修改文件
			var serverName string
			fmt.Print("请输入服务器名称 (留空使用默认服务器): ")
			fmt.Scanln(&serverName)
			err := modifyFiles(serverName)
			if err != nil {
				fmt.Printf("错误: %v\n", err)
			}
		case 7:
			// 回滚操作
			var serverName string
			fmt.Print("请输入服务器名称 (留空使用默认服务器): ")
			fmt.Scanln(&serverName)
			err := rollbackOperation(serverName)
			if err != nil {
				fmt.Printf("错误: %v\n", err)
			}
		case 8:
			// 服务器管理
			err := serverManagement()
			if err != nil {
				fmt.Printf("错误: %v\n", err)
			}
		case 9:
			// 启动HTTP服务
			fmt.Println("启动HTTP服务...")
			go func() {
				if err := startHTTPServer(); err != nil {
					fmt.Printf("HTTP服务启动失败: %v\n", err)
				}
			}()
			fmt.Println("HTTP服务已在后台启动")
		default:
			fmt.Println("无效的选择，请重新输入")
		}

		fmt.Println("按回车键继续...")
		// 清空输入缓冲区
		var c byte
		for {
			_, err := fmt.Scanf("%c", &c)
			if err != nil || c == '\n' {
				break
			}
		}
	}
}
