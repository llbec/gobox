package main

import (
	"io"
	"log"
	"net"
	"os"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ServerIP   string `yaml:"server_ip"`
	ServerPort int    `yaml:"server_port"`
	LogFile    string `yaml:"log_file"` // 日志文件路径
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	return &cfg, err
}

func main() {
	cfg, err := loadConfig("config.yml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 设置日志输出到控制台和文件
	var writers []io.Writer
	writers = append(writers, os.Stdout)

	if cfg.LogFile != "" {
		file, errOS := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if errOS != nil {
			log.Fatalf("打开日志文件失败: %v", errOS)
		}
		defer file.Close()
		writers = append(writers, file)
	}

	log.SetOutput(io.MultiWriter(writers...))

	addr := net.UDPAddr{
		IP:   net.ParseIP(cfg.ServerIP), // 监听所有地址
		Port: cfg.ServerPort,
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Printf("Error listening: %v\n", err)
		return
	}
	defer conn.Close()

	log.Printf("UDP server (multi-threaded) listening on port %d...", cfg.ServerPort)

	buffer := make([]byte, 1024)

	for {
		// 读取数据
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from UDP: %v\n", err)
			continue
		}

		// 拷贝消息内容，避免在 goroutine 里被覆盖
		message := strings.TrimSpace(string(buffer[:n]))
		client := *clientAddr // 拷贝地址

		// 每个请求开一个 goroutine 处理
		go handleRequest(conn, client, message)
	}
}

func handleRequest(conn *net.UDPConn, clientAddr net.UDPAddr, message string) {
	log.Printf("Received from %s: [%s]\n", clientAddr.String(), message)

	if EqualIgnoreSpace(message, "Who is center server?", true) || EqualIgnoreSpace(message, "Who is DCS server?", true) {
		reply := "It is me."
		_, err := conn.WriteToUDP([]byte(reply), &clientAddr)
		if err != nil {
			log.Printf("Error writing to UDP: %v\n", err)
			return
		}
		log.Printf("Replied to %s: %s\n", clientAddr.String(), reply)
	} else {
		log.Printf("Unknown query from %s: %s\n", clientAddr.String(), message)
		for i, r := range message {
			log.Printf("char[%d] = '%c' (U+%04X)\n", i, r, r)
		}
	}
}

func EqualIgnoreSpace(s1, s2 string, ignoreCase bool) bool {
	normalize := func(s string) string {
		var builder strings.Builder
		for _, r := range s {
			// 跳过空白和控制字符
			if unicode.IsSpace(r) || unicode.IsControl(r) {
				continue
			}
			builder.WriteRune(r)
		}
		res := builder.String()
		if ignoreCase {
			return strings.ToLower(res)
		}
		return res
	}
	return normalize(s1) == normalize(s2)
}
