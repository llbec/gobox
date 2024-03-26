package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {
	//获取命令行参数
	if len(os.Args) < 2 {
		fmt.Println("Example: ./tsBot filePath")
		return
	}
	filePath := os.Args[1]

	bs := make([]byte, 4)
	timeUnix := time.Now().Unix()
	binary.BigEndian.PutUint32(bs, uint32(timeUnix))

	//读写方式打开文件
	file, err := os.OpenFile(filePath, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("open file filed.", err)
		return
	}
	//defer关闭文件
	defer file.Close()

	//获取文件大小
	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	var size = stat.Size()
	fmt.Println("file size:", size)

	//读取文件内容到io中
	reader := bufio.NewReader(file)
	pos := int64(0)

	for {
		//读取每一行内容
		line, err := reader.ReadString('\n')
		if err != nil {
			//读到末尾
			if err == io.EOF {
				fmt.Println("File read ok!")
				break
			} else {
				fmt.Println("Read file error!", err)
				return
			}
		}
		//fmt.Println(line)

		//根据关键词覆盖当前行
		tsRegexp := regexp.MustCompile(`0x[0-9a-fA-F]+;`)
		if strings.Contains(line, "TSV_TIME_STAMP0") {
			bytes := []byte(tsRegexp.ReplaceAllString(line, fmt.Sprintf("0x%02X;", bs[0])))
			file.WriteAt(bytes, pos)
		} else if strings.Contains(line, "TSV_TIME_STAMP1") {
			bytes := []byte(tsRegexp.ReplaceAllString(line, fmt.Sprintf("0x%02X;", bs[1])))
			file.WriteAt(bytes, pos)
		} else if strings.Contains(line, "TSV_TIME_STAMP2") {
			bytes := []byte(tsRegexp.ReplaceAllString(line, fmt.Sprintf("0x%02X;", bs[2])))
			file.WriteAt(bytes, pos)
		} else if strings.Contains(line, "TSV_TIME_STAMP3") {
			bytes := []byte(tsRegexp.ReplaceAllString(line, fmt.Sprintf("0x%02X;", bs[3])))
			file.WriteAt(bytes, pos)

			fmt.Printf("Timestamp update to %v(0x%02X, 0x%02X, 0x%02X, 0x%02X)\n", time.Unix(timeUnix, 0).Format("2006-01-02 15:04:05"), bs[0], bs[1], bs[2], bs[3])

			return
		}

		//每一行读取完后记录位置
		pos += int64(len(line))
	}
}
