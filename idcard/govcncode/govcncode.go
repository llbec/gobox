package govcncode

import (
	//"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//GovCnCode supply china's code
type GovCnCode struct {
	document *goquery.Document
}

var _ AddressCode = (*GovCnCode)(nil)

//NewGovCnCode create a new china area code object
func NewGovCnCode() *GovCnCode {
	doc, err := goquery.NewDocument("http://www.gov.cn/test/2011-08/22/content_1930111.htm")
	if err != nil {
		panic(err)
	}
	ac := new(GovCnCode)
	ac.document = doc
	return ac
}

type record struct {
	code  int
	level int
	name  string
}

//GetProvinceList get a string list about all Provinces
func (addrCode *GovCnCode) GetProvinceList() []*Province {
	var records []*Province
	addrCode.document.Find("span").Each(func(i int, s *goquery.Selection) {
		//fmt.Printf(s.Text())
		texts := strings.Split(s.Text(), "\n")
		//fmt.Println(len(texts))
		reCode := regexp.MustCompile("^[\\d]{6}")
		reName := regexp.MustCompile("[\u4e00-\u9fa5]+$")
		var (
			province *Province
			city     *City
			county   *County
		)
		for _, v := range texts {
			r := new(record)
			r.level = getlevel(v)
			r.code, _ = strconv.Atoi(reCode.FindString(v))
			r.name = reName.FindString(v)
			//fmt.Println(v, r)
			switch r.level {
			case 1:
				if province != nil {
					records = append(records, province)
				}
				province = new(Province)
				province.Code = r.code
				province.Name = r.name
				//province.Citys = new([]*City)
			case 3:
				city = new(City)
				city.Code = r.code
				city.Name = r.name
				city.Parent = province
				//city.Countys = new([]*County)
				province.Citys = append(province.Citys, city)
			case 5:
				county = new(County)
				county.Code = r.code
				county.Name = r.name
				county.Parent = city
				city.Countys = append(city.Countys, county)
			default:
			}
		}
		records = append(records, province)
		//fmt.Println(len(records))
		return
	})
	return records
}

func getlevel(s string) int {
	c := 0
	for _, v := range s {
		if v == 160 {
			c++
		}
	}
	return c
}
