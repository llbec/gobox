package drage

import (
	"encoding/hex"
	"fmt"
	"strings"
)

const SYNC_BYTE = 0xA5

var IdValueMapping = map[string]string{
	"01": "hr",
	"04": "pvc",
	"5D": "nIBPm",
	"5B": "nIBPs",
	"5C": "nIBPd",
	"64": "spo2",
	"65": "pr",
	"60": "rr",
	"81": "temp1",
	"82": "temp2",
	"63": "iBPm",
	"61": "iBPs",
	"62": "iBPd",
	"19": "bis",
	"8B": "cco",
	"8C": "cci",
	"22": "artm",
	"20": "arts",
	"21": "artd",
	"25": "iBPm",
	"23": "iBPs",
	"24": "iBPd",
	"28": "iBPs",
	"26": "iBPs",
	"27": "iBPd",
	"2B": "LAP_mean",
	"29": "LAP_sys",
	"2A": "LAP_dia",
	"2E": "cvp",
	"2C": "RAP",
	"2D": "LAP",
	"2F": "ICP",
	"66": "awRR",
	"67": "etCo2",
	"69": "fiCO2",
	"92": "fio2",
	"A4": "etSev",
	"A3": "fiSev",
	"5E": "PCWP",
}

type Result struct {
	Offset int
	ID     string
	Name   string
	Value  string
}

func hexStringToBytes(s string) ([]byte, error) {

	s = strings.ReplaceAll(s, " ", "")
	return hex.DecodeString(s)
}

func decode(bytes []byte) string {

	var sb strings.Builder

	if len(bytes) < 1 {
		return "消息编码：-；状态编码：-；子包数量：0；"
	}

	if bytes[0] != SYNC_BYTE {
		return "消息编码：-；状态编码：-；子包数量：0；"
	}

	if len(bytes) < 4 {
		return "消息编码：-；状态编码：-；子包数量：0；"
	}

	transactionCode := bytes[3]

	if transactionCode != 0x57 && transactionCode != 0x77 {
		return "消息编码：-；状态编码：-；子包数量：0；"
	}

	if len(bytes) < 6 {
		return "消息编码：-；状态编码：-；子包数量：0；"
	}

	serverStatus := bytes[5]

	if len(bytes) <= 24 {
		return "消息编码：-；状态编码：-；子包数量：0；"
	}

	subPackets := int(bytes[23])

	sb.WriteString(fmt.Sprintf("消息编码 %02X；状态编码 %02X；子包数量 %d；\n", transactionCode, serverStatus, subPackets))

	startIndex := 24

	for i := 0; i < subPackets; i++ {

		if startIndex >= len(bytes) {
			break
		}

		packetLen := int(bytes[startIndex])

		if startIndex+1 >= len(bytes) {
			break
		}

		paras := int(bytes[startIndex+1] & 0x07)

		packetStart := startIndex
		packetEnd := startIndex + packetLen
		if packetEnd > len(bytes) {
			packetEnd = len(bytes)
		}
		packetRaw := bytes[packetStart:packetEnd]
		packetHex := hex.EncodeToString(packetRaw)

		sb.WriteString(fmt.Sprintf("子包%d：%s\n", i+1, formatHexString(packetHex)))
		sb.WriteString(fmt.Sprintf("    长度 %d，检测项数量 %d\n", packetLen, paras))

		paraIndex := startIndex + 8

		for k := 0; k < paras; k++ {

			if paraIndex >= len(bytes) {
				break
			}

			startRawIndex := paraIndex

			var paraID string

			if transactionCode == 0x57 {

				if paraIndex >= len(bytes) {
					break
				}

				iParaID := int(bytes[paraIndex] & 0xFF)

				if iParaID >= 0x37 && iParaID <= 0x5A {

					switch (iParaID - 0x37) % 3 {

					case 0:
						paraID = "61"
					case 1:
						paraID = "62"
					case 2:
						paraID = "63"
					}

				} else {

					paraID = fmt.Sprintf("%02X", bytes[paraIndex])
				}

				paraIndex++

			} else {

				if paraIndex+1 >= len(bytes) {
					break
				}

				paraID = fmt.Sprintf("%02X%02X", bytes[paraIndex], bytes[paraIndex+1])
				paraIndex += 2
			}

			paraID = strings.ToUpper(paraID)

			if paraIndex >= len(bytes) {
				break
			}

			paraStatus := bytes[paraIndex]
			_ = paraStatus
			paraIndex++

			var value strings.Builder

			for j := 0; j < 5; j++ {

				if paraIndex+j >= len(bytes) {
					break
				}

				if bytes[paraIndex+j] == 0x00 {

					paraIndex += j + 1
					break
				}

				value.WriteByte(bytes[paraIndex+j])
			}

			valueEndIndex := paraIndex
			paraVal := value.String()

			// 多体温通道
			if paraID == "80" || paraID == "84" || paraID == "87" {
				paraID = "81"
			}
			if paraID == "85" || paraID == "88" {
				paraID = "82"
			}

			if paraID == "96" {
				paraID = "67"
			}

			if paraID == "95" {
				paraID = "69"
			}

			if paraID == "68" {
				paraID = "66"
			}

			name, ok := IdValueMapping[paraID]
			if !ok {
				name = "unknown"
			}

			if paraVal != "" && !strings.HasPrefix(paraVal, "^^") && !strings.HasPrefix(paraVal, "?") {

				rawBytes := bytes[startRawIndex:valueEndIndex]
				rawHex := hex.EncodeToString(rawBytes)

				sb.WriteString(fmt.Sprintf("    检测项%d-%d：ID = %s，name = %s，Value = %s，raw data = %s\n",
					i+1, k+1, paraID, name, paraVal, formatHexString(rawHex)))
			}

		}

		startIndex += packetLen

	}

	return sb.String()
}

func formatHexString(hexStr string) string {
	var sb strings.Builder
	for i := 0; i < len(hexStr); i += 2 {
		if i > 0 {
			sb.WriteString(" ")
		}
		if i+2 <= len(hexStr) {
			sb.WriteString(hexStr[i : i+2])
		} else {
			sb.WriteString(hexStr[i:])
		}
	}
	return sb.String()
}

func printResult(results []Result) {

	fmt.Printf("检测项长度：%d\n\n", len(results))

	fmt.Printf("%-6s %-10s %-6s %-10s %-10s\n", "序号", "偏移量", "ID", "Name", "Value")

	for i, r := range results {

		fmt.Printf("%-6d %-10d %-6s %-10s %-10s\n",
			i+1,
			r.Offset,
			r.ID,
			r.Name,
			r.Value)
	}

}
