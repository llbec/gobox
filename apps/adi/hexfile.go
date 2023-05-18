package adi

import (
	"fmt"
	"strconv"
)

/*
LL MM MM TT DD...DD CC
length memory_offset type data CC
cc = 100 - (LL + MM + MM + TT + DD + ... + DD)%100
*/

/* str = ll mm mm tt dd...dd */
func hexCC(str string) (cc int) {
	ll, e := strconv.ParseInt(string([]byte(str)[0:2]), 16, 32)
	if e != nil {
		panic(e)
	}
	lowLen := int(ll + 4)
	if lowLen*2 != len(str) {
		panic(fmt.Errorf("Invalid %d != %d", lowLen*2, len(str)))
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
