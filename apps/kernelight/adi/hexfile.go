package adi

import (
	"fmt"
	"strconv"
	"strings"
)

/*
LL MM MM TT DD...DD CC
length memory_offset type data CC
cc = 100 - (LL + MM + MM + TT + DD + ... + DD)%100
*/
type HexRow struct {
	len    int
	offset int
	style  int
	data   []byte
	cc     int
}

func NewHexRow(hex string) *HexRow {
	cc := hexCC(hex)

	idx := 0
	ll, e := strconv.ParseInt(string([]byte(hex)[idx:idx+2]), 16, 32)
	if e != nil {
		panic(e)
	}
	idx += 2

	mm, e := strconv.ParseInt(string([]byte(hex)[idx:idx+4]), 16, 32)
	if e != nil {
		panic(e)
	}
	idx += 4

	tt, e := strconv.ParseInt(string([]byte(hex)[idx:idx+2]), 16, 32)
	if e != nil {
		panic(e)
	}
	idx += 2

	var buf []byte
	for i := 0; i < int(ll); i++ {
		dd, e := strconv.ParseInt(string([]byte(hex)[idx:idx+2]), 16, 32)
		if e != nil {
			panic(e)
		}
		idx += 2
		buf = append(buf, byte(dd))
	}

	c, e := strconv.ParseInt(string([]byte(hex)[idx:idx+2]), 16, 32)
	if e != nil {
		panic(e)
	}
	if int(c) != cc {
		panic(fmt.Errorf("check code invalid %v != %v", c, cc))
	}

	return &HexRow{
		len:    int(ll),
		offset: int(mm),
		style:  int(tt),
		data:   buf,
		cc:     int(c),
	}
}

func (row *HexRow) Hex() string {
	//hex := strconv.FormatInt(int64(row.len), 16)
	hex := fmt.Sprintf("%02x", row.len)
	hex += fmt.Sprintf("%04x", row.offset)
	hex += fmt.Sprintf("%02x", row.style)
	for i := 0; i < row.len; i++ {
		//hex += strconv.FormatInt(int64(row.data[i]), 16)
		hex += fmt.Sprintf("%02x", row.data[i])
	}
	hex += fmt.Sprintf("%02x", row.cc)
	return strings.ToUpper(hex)
}

func (row *HexRow) UpdateData(buf []byte, m int, l int) {
	if l > len(buf) || m+l > row.len {
		panic(fmt.Errorf("invalid para!"))
	}
	for i := 0; i < l; i++ {
		row.data[m+i] = buf[i]
	}
	row.cc = hexCC(row.Hex())
}

func (row *HexRow) I2CEnable() {
	if row.len != 0x10 || row.offset != 0x10 || row.style != 0 {
		panic(fmt.Errorf("invalid string"))
	}
	row.UpdateData([]byte{0xFF, 0xFF, 0xFF, 0xFF}, 4, 4)
}

func (row *HexRow) I2CDisable() {
	if row.len != 0x10 || row.offset != 0x10 || row.style != 0 {
		panic(fmt.Errorf("invalid string"))
	}
	row.UpdateData([]byte{0x12, 0x34, 0x56, 0x78}, 4, 4)
}

/* str = ll mm mm tt dd...dd */
func hexCC(str string) (cc int) {
	ll, e := strconv.ParseInt(string([]byte(str)[0:2]), 16, 32)
	if e != nil {
		panic(e)
	}
	lowLen := int(ll + 4)
	if lowLen*2 > len(str) {
		panic(fmt.Errorf("Invalid(%v) %d > %d", str, lowLen*2, len(str)))
	}
	sum := 0
	for i := 0; i < lowLen; i++ {
		v, e := strconv.ParseInt(string([]byte(str)[i*2:i*2+2]), 16, 32)
		if e != nil {
			panic(e)
		}
		sum += int(v)
	}
	return 0x100 - (sum % 0x100)
}
