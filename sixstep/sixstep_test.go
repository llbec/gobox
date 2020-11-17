package sixstep

import (
	//"fmt"
	"testing"

	"gobox/govcncode"
	//"personid/types"
)

func TestRun(t *testing.T) {
	addrdata := govcncode.NewGovCnCode()
	task := NewSixStep(addrdata.GetProvinceList())
	task.Run()
}
