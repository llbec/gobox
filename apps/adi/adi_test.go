package adi

import (
	"fmt"
	"strconv"
	"testing"
)

func TestHexCC(t *testing.T) {
	cc := hexCC("1000100018F09FE50000A0E118F09FE518F09FE5BB")
	fmt.Printf("%s\n", strconv.FormatInt(int64(cc), 16))
}

func TestRow(t *testing.T) {
	row := NewHexRow("1000100018F09FE50000A0E118F09FE518F09FE5BB")
	row.I2CEnable()
	fmt.Printf("%v\n", row.Hex())
	row.I2CDisable()
	fmt.Printf("%v\n", row.Hex())
}
