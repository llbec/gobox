package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ServerIP     string `yaml:"server_ip"`
	ServerPort   int    `yaml:"server_port"`
	SendFile     string `yaml:"send_file"`
	PeriodicSend bool   `yaml:"periodic_send"`
	SendInterval int    `yaml:"send_interval"`
	SendData     bool   `yaml:"send_data"`
	WaitTime     int    `yaml:"wait_time"`
	LogFile      string `yaml:"log_file"` // 日志文件路径
}

var (
	sendCount uint64 = 0
	rcvCount  uint64 = 0
)

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	return &cfg, err
}

// hex 字符串（带空格） -> byte 数组
func hexStringToBytes(s string) ([]byte, error) {
	parts := strings.Fields(s)
	res := make([]byte, len(parts))
	for i, p := range parts {
		v, err := strconv.ParseUint(p, 16, 8)
		if err != nil {
			return nil, fmt.Errorf("解析 hex %s 失败: %v", p, err)
		}
		res[i] = byte(v)
	}
	return res, nil
}

// byte 数组 -> 大写 hex，每16个换行
func bytesToHexString(data []byte) string {
	var sb strings.Builder
	for i, b := range data {
		sb.WriteString(fmt.Sprintf("%02X", b))
		if i != len(data)-1 {
			sb.WriteByte(' ')
		}
		if (i+1)%16 == 0 && i != len(data)-1 {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func main() {
	cfg, err := loadConfig("client.yml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 设置日志输出到控制台和文件
	var writers []io.Writer
	writers = append(writers, os.Stdout)

	if cfg.LogFile != "" {
		file, err1 := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err1 != nil {
			log.Fatalf("打开日志文件失败: %v", err1)
		}
		defer file.Close()
		writers = append(writers, file)
	}

	log.SetOutput(io.MultiWriter(writers...))

	addr := net.JoinHostPort(cfg.ServerIP, fmt.Sprintf("%d", cfg.ServerPort))
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()
	log.Printf("连接成功 %s", addr)

	// 读取发送文件
	sendData := []byte{}
	if cfg.SendData {
		content, err := os.ReadFile(cfg.SendFile)
		if err != nil {
			log.Fatalf("读取发送文件失败: %v", err)
		}
		sendData, err = hexStringToBytes(string(content))
		if err != nil {
			log.Fatalf("解析发送数据失败: %v", err)
		}
		log.Printf("待发送数据:\n%s", bytesToHexString(sendData))
	}

	// 等待配置的等待时间
	if cfg.WaitTime > 0 {
		time.Sleep(time.Duration(cfg.WaitTime) * time.Second)
	}

	// 启动接收协程
	go func() {
		reader := bufio.NewReader(conn)
		buf := make([]byte, 1024*4)
		for {
			n, err := reader.Read(buf)
			if err != nil {
				log.Printf("读取数据失败: %v", err)
				return
			}
			if n > 0 {
				rcvCount++
				data := buf[:n]
				log.Printf("[%v]-[%v]收到数据:\n%s", conn.RemoteAddr().String(), rcvCount, bytesToHexString(data))
			}
		}
	}()

	// 发送数据
	if cfg.SendData {
		if cfg.PeriodicSend {
			ticker := time.NewTicker(time.Duration(cfg.SendInterval) * time.Millisecond)
			defer ticker.Stop()
			for {
				_, err := conn.Write(sendData)
				if err != nil {
					log.Printf("发送数据失败: %v", err)
					break
				}
				//log.Printf("已发送数据:\n%s", bytesToHexString(sendData))
				sendCount++
				log.Printf("[%v]已发送数据%v\n", conn.RemoteAddr().String(), sendCount)
				<-ticker.C
			}
		} else {
			_, err := conn.Write(sendData)
			if err != nil {
				log.Printf("发送数据失败: %v", err)
			} else {
				log.Printf("已发送数据:\n%s", bytesToHexString(sendData))
			}
		}
	}

	// 阻塞 main
	select {}
}
