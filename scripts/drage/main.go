package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// 配置结构体
type Config struct {
	Server struct {
		Address           string `yaml:"address"`
		HeartbeatInterval int    `yaml:"heartbeat_interval"`
		HeartbeatData     string `yaml:"heartbeat_data"`
	} `yaml:"server"`
	Commands map[string]string `yaml:"commands"`
}

// 加载配置文件
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func main() {
	cfg, err := LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("❌ 配置文件加载失败: %v", err)
	}

	fmt.Printf("✅ 连接到服务器: %s\n", cfg.Server.Address)

	conn, err := net.Dial("tcp", cfg.Server.Address)
	if err != nil {
		log.Fatalf("❌ 无法连接到服务器: %v", err)
	}
	defer conn.Close()

	fmt.Println("🚀 已连接，开始通信...")

	// 启动心跳协程
	go startHeartbeat(conn, cfg.Server.HeartbeatData, cfg.Server.HeartbeatInterval)

	// 启动接收协程
	go func() {
		reader := bufio.NewReader(conn)
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("⚠️ 读取错误或连接关闭: %v", err)
				return
			}
			message = strings.TrimSpace(message)
			if message == "" {
				continue
			}
			fmt.Printf("📩 收到消息: %s\n", message)

			// 查找回应
			response := matchResponse(message, cfg.Commands)
			if response != "" {
				_, err := conn.Write([]byte(response + "\n"))
				if err != nil {
					log.Printf("⚠️ 发送回应失败: %v", err)
				} else {
					fmt.Printf("📤 已回应: %s -> %s\n", message, response)
				}
			} else {
				fmt.Println("🤔 无匹配回应，忽略。")
			}
		}
	}()

	// 主线程等待退出
	select {}
}

// 心跳包协程
func startHeartbeat(conn net.Conn, data string, interval int) {
	if interval <= 0 {
		return
	}
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		_, err := conn.Write([]byte(data + "\n"))
		if err != nil {
			log.Printf("⚠️ 心跳发送失败: %v", err)
			return
		}
		fmt.Printf("💓 已发送心跳: %s\n", data)
	}
}

// 根据收到的消息查找回应
func matchResponse(message string, commands map[string]string) string {
	upperMsg := strings.ToUpper(strings.TrimSpace(message))
	for k, v := range commands {
		if strings.EqualFold(upperMsg, k) {
			return v
		}
	}
	return ""
}
