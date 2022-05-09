package govcncode

import (
	"fmt"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestNewACode(t *testing.T) {
	ac := NewGovCnCode()
	ac.document.Find("title").Each(func(i int, s *goquery.Selection) {
		fmt.Print(s.Text())
	})
}

func TestGetProvinceList(t *testing.T) {
	ac := NewGovCnCode()
	ps := ac.GetProvinceList()
	fmt.Println(len(ps))
	for i, v := range ps {
		fmt.Printf("%-4d%-20s\t%6d\n", i, v.Name, v.Code)
		for j, city := range v.Citys {
			fmt.Printf("\t%d-%d\t%-20s\t%6d\n", i, j, city.Name, city.Code)
			fmt.Printf("\t\t")
			for k, county := range city.Countys {
				fmt.Printf("(%d-%d-%d,%s,%d)", i, j, k, county.Name, county.Code)
			}
			fmt.Printf("\n")
		}
	}
}
