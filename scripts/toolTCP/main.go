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
	HostIP       string `yaml:"host_ip"`
	HostPort     int    `yaml:"host_port"`
	SendFile     string `yaml:"send_file"`
	PeriodicSend bool   `yaml:"periodic_send"`
	SendInterval int    `yaml:"send_interval"`
	SendData     bool   `yaml:"send_data"`
	WaitTime     int    `yaml:"wait_time"`
	LogFile      string `yaml:"log_file"`    // 日志文件路径
	TcpType      string `yaml:"tcp_type"`    // tcp类型，server or client
	EncodeType   string `yaml:"encode_type"` // 编码类型，hex or string
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

// recive data
// param conn net.Conn
// cfg *Config
func reciveData(conn net.Conn, cfg *Config) {
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
			if strings.ToLower(cfg.EncodeType) == "string" {
				log.Printf("[%v]-[%v]收到数据:\n%s", conn.RemoteAddr().String(), rcvCount, string(data))
			} else {
				log.Printf("[%v]-[%v]收到数据:\n%s", conn.RemoteAddr().String(), rcvCount, bytesToHexString(data))
			}
		}
	}
}

// send data
// param conn net.Conn
// cfg *Config
func sendData(conn net.Conn, cfg *Config) {
	// 等待配置的等待时间
	if cfg.WaitTime > 0 {
		time.Sleep(time.Duration(cfg.WaitTime) * time.Second)
	}

	// 发送数据
	if cfg.SendData {
		content, err := os.ReadFile(cfg.SendFile)
		if err != nil {
			log.Fatalf("读取发送文件失败: %v", err)
		}
		var sendData []byte
		if strings.ToLower(cfg.EncodeType) == "string" {
			sendData = content
			log.Printf("待发送数据:\n%s", string(sendData))
		} else {
			sendData, err = hexStringToBytes(string(content))
			if err != nil {
				log.Fatalf("解析发送数据失败: %v", err)
			}
			log.Printf("待发送数据:\n%s", bytesToHexString(sendData))
		}

		if cfg.PeriodicSend {
			ticker := time.NewTicker(time.Duration(cfg.SendInterval) * time.Millisecond)
			log.Printf("每隔 %d 毫秒发送一次数据", cfg.SendInterval)
			defer ticker.Stop()
			for {
				_, err = conn.Write(sendData)
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
			_, err = conn.Write(sendData)
			if err != nil {
				log.Printf("发送数据失败: %v", err)
			} else {
				log.Printf("已发送数据:\n%s", bytesToHexString(sendData))
			}
		}
	}
}

// tcp server
// param cfg *Config
func server(cfg *Config) {
	addr := net.JoinHostPort(cfg.HostIP, fmt.Sprintf("%d", cfg.HostPort))
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("监听失败: %v", err)
	}
	defer ln.Close()
	log.Printf("服务器已启动，监听 %s", addr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("接受连接失败: %v", err)
			continue
		}
		log.Printf("新连接来自 %s", conn.RemoteAddr().String())
		go func(c net.Conn) {
			defer c.Close()
			// 启动接收协程
			go reciveData(c, cfg)

			// 发送数据
			go sendData(c, cfg)

			select {}
		}(conn)
	}
	// 阻塞 main
	select {}
}

// tcp client
// param cfg *Config
// param sendData []byte
func client(cfg *Config) {
	addr := net.JoinHostPort(cfg.HostIP, fmt.Sprintf("%d", cfg.HostPort))
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()
	log.Printf("连接成功 %s", addr)

	// 等待配置的等待时间
	if cfg.WaitTime > 0 {
		time.Sleep(time.Duration(cfg.WaitTime) * time.Second)
	}

	// 启动接收协程
	go reciveData(conn, cfg)

	// 发送数据
	sendData(conn, cfg)

	// 阻塞 main
	select {}
}

// main
func main() {
	cfg, err := loadConfig("tcp.yml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 设置日志输出到控制台和文件
	var writers []io.Writer
	writers = append(writers, os.Stdout)

	if cfg.LogFile != "" {
		file, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("打开日志文件失败: %v", err)
		}
		defer file.Close()
		writers = append(writers, file)
	}

	log.SetOutput(io.MultiWriter(writers...))

	// 根据配置启动 server 或 client
	log.Printf("配置: %+v", cfg)
	if strings.ToLower(cfg.TcpType) == "server" {
		log.Printf("启动 TCP 服务器，监听 %s:%d", cfg.HostIP,
			cfg.HostPort)
		server(cfg)
	} else {
		log.Printf("启动 TCP 客户端，连接 %s:%d", cfg.HostIP,
			cfg.HostPort)
		client(cfg)
	}
}
