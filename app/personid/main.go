package main

import (
	"gobox/govcncode"
	"gobox/sixstep"
)

var addrdata govcncode.AddressCode

func main() {
	addrdata = govcncode.NewGovCnCode()
	task := sixstep.NewSixStep(addrdata.GetProvinceList())
	task.Run()
}
