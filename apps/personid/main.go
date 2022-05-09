package main

import (
	"gobox/src/idcard/govcncode"
	"gobox/src/idcard/sixstep"
)

var addrdata govcncode.AddressCode

func main() {
	addrdata = govcncode.NewGovCnCode()
	task := sixstep.NewSixStep(addrdata.GetProvinceList())
	task.Run()
}
