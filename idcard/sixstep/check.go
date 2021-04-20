package sixstep

import (
	"fmt"
)

var (
	weighted = [17]int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	code     = map[int]string{0: "1", 1: "0", 2: "X", 3: "9", 4: "8", 5: "7", 6: "6", 7: "5", 8: "4", 9: "3", 10: "2"}
)

//CalcCode calcute the check code
func CalcCode(s17 string) (string, error) {
	if len(s17) != 17 {
		return "", fmt.Errorf("Error length invalid")
	}
	sum := 0
	for i, v := range s17 {
		sum += (int(v) - '0') * weighted[i]
	}
	return code[sum%11], nil
}
