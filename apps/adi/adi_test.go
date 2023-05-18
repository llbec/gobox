package adi

import (
	"fmt"
	"strconv"
	"testing"
)

func TestHexCC(t *testing.T) {
	cc := hexCC("1000100018F09FE5FFFFFFFF18F09FE518F09FE5")
	fmt.Printf("%s\n", strconv.FormatInt(int64(cc), 16))
}
