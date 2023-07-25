package main

import (
	"encoding/binary"
	"fmt"
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
