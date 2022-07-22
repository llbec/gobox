package main

import (
	"fmt"
	"gobox/src/website/dcg"

	"github.com/xuri/excelize/v2"
)

func main() {
	err := dcg.FillPortfolio("https://dcg.co/portfolio/")
	if err != nil {
		fmt.Println(err)
		return
	}
	f := excelize.NewFile()
	for i, company := range dcg.Portfolio {
		id := i + 1
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", id), company.Sector)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", id), company.Name)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", id), company.Details)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", id), company.Location)
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", id), company.Description)
		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", id), company.Url)
	}
	if err := f.SaveAs("dcg_portfolio.xlsx"); err != nil {
		fmt.Println(err)
	}
}
