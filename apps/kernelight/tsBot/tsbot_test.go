package main

import (
	"encoding/binary"
	"fmt"
	"regexp"
	"testing"
	"time"
)

func Test_Ts(t *testing.T) {
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, uint32(time.Now().Unix()))
	for i := 0; i < len(bs); i++ {
		fmt.Printf("%2X\n", bs[i])
	}
}

func Test_Reg(t *testing.T) {
	source := "    module_map[MOD_TIME_STAMP0] = 0x64;"
	tsRegexp := regexp.MustCompile(`0x[0-9a-fA-F]+;`)

	dest := tsRegexp.ReplaceAllString(source, fmt.Sprintf("0x%2X;", 17))
	fmt.Println(dest)

	if dest != "    module_map[MOD_TIME_STAMP0] = 0x11;" {
		t.Fatalf(fmt.Sprintf("expexted:       module_map[MOD_TIME_STAMP0] = 0x11;\n \nUnexpected: %v\n", dest))
	}
}
