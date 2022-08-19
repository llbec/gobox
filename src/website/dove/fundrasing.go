package dove

import (
	"os"

	"github.com/PuerkitoBio/goquery"
)

var (
	Projects []*Project
	Map      map[string]*Project
)

func init() {
	Projects = make([]*Project, 0)
}

func LoadFundrasing() error {
	fileObj, err := os.Open("dove")
	if err != nil {
		return err
	}
	doc, err := goquery.NewDocumentFromReader(fileObj)
	if err != nil {
		return err
	}
	right := doc.Find("div.dataRightPaneInnerContent.paneInnerContent").Children()
	Projects = Projects[:0]
	Map = make(map[string]*Project)

	for _, node := range right.Nodes {
		rn := goquery.NewDocumentFromNode(node)
		pro := &Project{}
		for _, cell := range rn.Find("div.cell.read").Nodes {
			cl := goquery.NewDocumentFromNode(cell)
			switch cl.AttrOr("data-columnid", "") {
			case "fldSGmVP3olLGmAID": //date
				pro.Date = cl.Find("div.truncate.css-10jy3hn").Text()
			case "fldbHI1iWcw6U912R": //amount
				pro.Amount = cl.Find("div.flex-auto.truncate.line-height-4.right-align.tabular-nums").Text()
			case "fldhntVnAppLIOUAl": //investors
				pro.Investors = cl.Find("div.foreign-key-blue.rounded.px-half.flex-none.mr-half.items-center.fit.truncate.line-height-4.text-dark.flex-inline").Text()
			case "fld8ZuXfzyuH1b9Dv": //url
				pro.Website = cl.Find("span.url").Text()
			case "fldZpYosKMqtY2nqQ": //founder
				pro.Founder = cl.Find("div.foreign-key-blue.rounded.px-half.flex-none.mr-half.items-center.fit.truncate.line-height-4.text-dark.flex-inline").Text()
			case "fldjd43zfXdpAWzaq": //category
				pro.Category = cl.Find("div.foreign-key-blue.rounded.px-half.flex-none.mr-half.items-center.fit.truncate.line-height-4.text-dark.flex-inline").Text()
			case "fldT0Fasv4hkjwbb3": //subcategories
				pro.Subcategories = cl.Find("div.foreign-key-blue.rounded.px-half.flex-none.mr-half.items-center.fit.truncate.line-height-4.text-dark.flex-inline").Text()
			case "fldJHMHegLEl2A56n": //description
				pro.Description = cl.Find("div.line-height-4.overflow-hidden.truncate").Text()
			case "fld0t86SH12Fx2aD6": //stages
				pro.Stages = cl.Find("div.foreign-key-blue.rounded.px-half.flex-none.mr-half.items-center.fit.truncate.line-height-4.text-dark.flex-inline").Text()
			case "fld05vhEx0sIhbK4B": //valuation
				pro.Valuation = cl.Find("div.flex-auto.truncate.line-height-4.right-align.tabular-nums").Text()
			case "fldruqJ51OKbiswEI": //project
				pro.Name = cl.Find("div.foreign-key-blue.rounded.px-half.flex-none.mr-half.items-center.fit.truncate.line-height-4.text-dark.flex-inline").Text()
			case "fldcDMf8A5D64Ecdf": //announcement
				pro.Announcement = cl.Find("span.url").Text()
			}
		}
		Projects = append(Projects, pro)
		Map[pro.Name] = pro
	}

	return nil
}
