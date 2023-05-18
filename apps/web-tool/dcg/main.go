package main

import (
	"fmt"
	"gobox/src/website/dcg"
	"gobox/src/website/dove"

	"github.com/xuri/excelize/v2"
)

var (
	projects map[string]*project
)

func main() {
	err := dcg.FillPortfolio("https://dcg.co/portfolio/")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = dove.LoadFundrasing()
	if err != nil {
		fmt.Println(err)
		return
	}

	projects = make(map[string]*project)
	excel := excelize.NewFile()
	excel.NewSheet("portfolio")
	excel.NewSheet("dove")

	excel.SetCellValue("portfolio", "A1", "Sector")
	excel.SetCellValue("portfolio", "B1", "项目")
	excel.SetCellValue("portfolio", "C1", "简介")
	excel.SetCellValue("portfolio", "D1", "地点")
	excel.SetCellValue("portfolio", "E1", "描述")
	excel.SetCellValue("portfolio", "F1", "网站")
	for i, company := range dcg.Portfolio {
		id := i + 2
		excel.SetCellValue("portfolio", fmt.Sprintf("A%d", id), company.Sector)
		excel.SetCellValue("portfolio", fmt.Sprintf("B%d", id), company.Name)
		excel.SetCellValue("portfolio", fmt.Sprintf("C%d", id), company.Details)
		excel.SetCellValue("portfolio", fmt.Sprintf("D%d", id), company.Location)
		excel.SetCellValue("portfolio", fmt.Sprintf("E%d", id), company.Description)
		excel.SetCellValue("portfolio", fmt.Sprintf("F%d", id), company.Url)

		p := &project{}
		p.copyCompany(company)
		projects[company.Name] = p
	}

	excel.SetCellValue("dove", "A1", "项目")
	excel.SetCellValue("dove", "B1", "投资时机")
	excel.SetCellValue("dove", "C1", "日期")
	excel.SetCellValue("dove", "D1", "金额")
	excel.SetCellValue("dove", "E1", "投资者")
	excel.SetCellValue("dove", "F1", "官网")
	excel.SetCellValue("dove", "G1", "创始人")
	excel.SetCellValue("dove", "H1", "类别")
	excel.SetCellValue("dove", "I1", "子类别")
	excel.SetCellValue("dove", "J1", "描述")
	excel.SetCellValue("dove", "K1", "估值")
	excel.SetCellValue("dove", "L1", "公告链接")
	for i, prjt := range dove.Projects {
		id := i + 2
		excel.SetCellValue("dove", fmt.Sprintf("A%d", id), prjt.Name)
		excel.SetCellValue("dove", fmt.Sprintf("B%d", id), prjt.Stages)
		excel.SetCellValue("dove", fmt.Sprintf("C%d", id), prjt.Date)
		excel.SetCellValue("dove", fmt.Sprintf("D%d", id), prjt.Amount)
		excel.SetCellValue("dove", fmt.Sprintf("E%d", id), prjt.Investors)
		excel.SetCellValue("dove", fmt.Sprintf("F%d", id), prjt.Website)
		excel.SetCellValue("dove", fmt.Sprintf("G%d", id), prjt.Founder)
		excel.SetCellValue("dove", fmt.Sprintf("H%d", id), prjt.Category)
		excel.SetCellValue("dove", fmt.Sprintf("I%d", id), prjt.Subcategories)
		excel.SetCellValue("dove", fmt.Sprintf("J%d", id), prjt.Description)
		excel.SetCellValue("dove", fmt.Sprintf("K%d", id), prjt.Valuation)
		excel.SetCellValue("dove", fmt.Sprintf("L%d", id), prjt.Announcement)

		copyProject(prjt)
	}

	excel.SetCellValue("Sheet1", "A1", "项目")
	excel.SetCellValue("Sheet1", "B1", "投资时机")
	excel.SetCellValue("Sheet1", "C1", "日期")
	excel.SetCellValue("Sheet1", "D1", "金额")
	excel.SetCellValue("Sheet1", "E1", "估值")
	excel.SetCellValue("Sheet1", "F1", "类别")
	excel.SetCellValue("Sheet1", "G1", "子类别")
	excel.SetCellValue("Sheet1", "H1", "Sector")
	excel.SetCellValue("Sheet1", "I1", "地点")
	excel.SetCellValue("Sheet1", "J1", "创始人")
	excel.SetCellValue("Sheet1", "K1", "简介")
	excel.SetCellValue("Sheet1", "L1", "描述")
	excel.SetCellValue("Sheet1", "M1", "网站")
	excel.SetCellValue("Sheet1", "N1", "公告链接")
	i := 2
	for _, v := range projects {
		excel.SetCellValue("Sheet1", fmt.Sprintf("A%d", i), v.name)
		excel.SetCellValue("Sheet1", fmt.Sprintf("B%d", i), v.stages)
		excel.SetCellValue("Sheet1", fmt.Sprintf("C%d", i), v.date)
		excel.SetCellValue("Sheet1", fmt.Sprintf("D%d", i), v.amount)
		excel.SetCellValue("Sheet1", fmt.Sprintf("E%d", i), v.valuation)
		excel.SetCellValue("Sheet1", fmt.Sprintf("F%d", i), v.category)
		excel.SetCellValue("Sheet1", fmt.Sprintf("G%d", i), v.subcategories)
		excel.SetCellValue("Sheet1", fmt.Sprintf("H%d", i), v.sector)
		excel.SetCellValue("Sheet1", fmt.Sprintf("I%d", i), v.location)
		excel.SetCellValue("Sheet1", fmt.Sprintf("J%d", i), v.founder)
		excel.SetCellValue("Sheet1", fmt.Sprintf("K%d", i), v.details)
		excel.SetCellValue("Sheet1", fmt.Sprintf("L%d", i), v.description)
		excel.SetCellValue("Sheet1", fmt.Sprintf("M%d", i), v.url)
		excel.SetCellValue("Sheet1", fmt.Sprintf("N%d", i), v.announcement)
		i++
	}

	if err := excel.SaveAs("dcg_portfolio.xlsx"); err != nil {
		fmt.Println(err)
	}
}

func (p *project) copyCompany(cmpy *dcg.Company) {
	p.name = cmpy.Name
	p.sector = cmpy.Sector
	p.details = cmpy.Details
	p.location = cmpy.Location
	p.description = cmpy.Description
	p.url = cmpy.Url
}

func copyProject(dovepro *dove.Project) {
	_, exsit := projects[dovepro.Name]
	if !exsit {
		projects[dovepro.Name] = &project{
			name:          dovepro.Name,
			stages:        dovepro.Stages,
			date:          dovepro.Date,
			amount:        dovepro.Amount,
			investors:     dovepro.Investors,
			founder:       dovepro.Founder,
			category:      dovepro.Category,
			subcategories: dovepro.Subcategories,
			valuation:     dovepro.Valuation,
			announcement:  dovepro.Announcement,
		}
	} else {
		if projects[dovepro.Name].date != "" && projects[dovepro.Name].date != dovepro.Date {
			projects[dovepro.Name].stages = fmt.Sprintf("%s<>%s", projects[dovepro.Name].stages, dovepro.Stages)
			projects[dovepro.Name].date = fmt.Sprintf("%s<>%s", projects[dovepro.Name].date, dovepro.Date)
			projects[dovepro.Name].amount = fmt.Sprintf("%s<>%s", projects[dovepro.Name].amount, dovepro.Amount)
			projects[dovepro.Name].valuation = fmt.Sprintf("%s<>%s", projects[dovepro.Name].valuation, dovepro.Valuation)
			fmt.Printf("update project %s at %s for %s\n", dovepro.Name, dovepro.Date, dovepro.Amount)
		} else {
			projects[dovepro.Name].stages = dovepro.Stages
			projects[dovepro.Name].date = dovepro.Date
			projects[dovepro.Name].amount = dovepro.Amount
			projects[dovepro.Name].investors = dovepro.Investors
			projects[dovepro.Name].founder = dovepro.Founder
			projects[dovepro.Name].category = dovepro.Category
			projects[dovepro.Name].subcategories = dovepro.Subcategories
			projects[dovepro.Name].valuation = dovepro.Valuation
			projects[dovepro.Name].announcement = dovepro.Announcement
		}
	}
}
