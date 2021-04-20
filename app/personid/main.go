package main

import (
	"gobox/idcard/govcncode"
	"gobox/idcard/sixstep"
)

var addrdata govcncode.AddressCode

func main() {
	addrdata = govcncode.NewGovCnCode()
	task := sixstep.NewSixStep(addrdata.GetProvinceList())
	task.Run()
}
