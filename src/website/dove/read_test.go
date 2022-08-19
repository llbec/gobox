package dove

import (
	"fmt"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestRead(t *testing.T) {
	fileObj, err := os.Open("dove")
	if err != nil {
		t.Fatal(err)
	}
	doc, err := goquery.NewDocumentFromReader(fileObj)
	if err != nil {
		t.Fatal(err)
	}
	left := doc.Find("div.dataLeftPaneInnerContent.paneInnerContent").Children()
	right := doc.Find("div.dataRightPaneInnerContent.paneInnerContent").Children()
	fmt.Printf("\n%v ?= %v\n", len(left.Nodes), len(right.Nodes))
	for i, node := range right.Nodes {
		rn := goquery.NewDocumentFromNode(node)
		//fundRound := goquery.NewDocumentFromNode(left.Nodes[i])
		//rounds := strings.Split(fundRound.Find("div.line-height-4.overflow-hidden.truncate").Text(), "-")
		//round := strings.TrimSpace(rounds[len(rounds)-1])
		var (
			date          string
			amount        string
			investors     string
			weburl        string
			founder       string
			category      string
			subcategories string
			description   string
			stages        string
			valuation     string
			project       string
			announcement  string
		)
		for _, cell := range rn.Find("div.cell.read").Nodes {
			cl := goquery.NewDocumentFromNode(cell)
			switch cl.AttrOr("data-columnid", "") {
			case "fldSGmVP3olLGmAID": //date
				date = cl.Find("div.truncate.css-10jy3hn").Text()
			case "fldbHI1iWcw6U912R": //amount
				amount = cl.Find("div.flex-auto.truncate.line-height-4.right-align.tabular-nums").Text()
			case "fldhntVnAppLIOUAl": //investors
				investors = cl.Find("div.foreign-key-blue.rounded.px-half.flex-none.mr-half.items-center.fit.truncate.line-height-4.text-dark.flex-inline").Text()
			case "fld8ZuXfzyuH1b9Dv": //url
				weburl = cl.Find("span.url").Text()
			case "fldZpYosKMqtY2nqQ": //founder
				founder = cl.Find("div.foreign-key-blue.rounded.px-half.flex-none.mr-half.items-center.fit.truncate.line-height-4.text-dark.flex-inline").Text()
			case "fldjd43zfXdpAWzaq": //category
				category = cl.Find("div.foreign-key-blue.rounded.px-half.flex-none.mr-half.items-center.fit.truncate.line-height-4.text-dark.flex-inline").Text()
			case "fldT0Fasv4hkjwbb3": //subcategories
				subcategories = cl.Find("div.foreign-key-blue.rounded.px-half.flex-none.mr-half.items-center.fit.truncate.line-height-4.text-dark.flex-inline").Text()
			case "fldJHMHegLEl2A56n": //description
				description = cl.Find("div.line-height-4.overflow-hidden.truncate").Text()
			case "fld0t86SH12Fx2aD6": //stages
				stages = cl.Find("div.foreign-key-blue.rounded.px-half.flex-none.mr-half.items-center.fit.truncate.line-height-4.text-dark.flex-inline").Text()
			case "fld05vhEx0sIhbK4B": //valuation
				valuation = cl.Find("div.flex-auto.truncate.line-height-4.right-align.tabular-nums").Text()
			case "fldruqJ51OKbiswEI": //project
				project = cl.Find("div.foreign-key-blue.rounded.px-half.flex-none.mr-half.items-center.fit.truncate.line-height-4.text-dark.flex-inline").Text()
			case "fldcDMf8A5D64Ecdf": //announcement
				announcement = cl.Find("span.url").Text()
			}
		}
		//fmt.Printf("%d\t%s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s\n", i, round, date, amount, investors, weburl, founder, category, subcategories, description, stages, valuation, project, announcement)
		fmt.Printf("%d\t%s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s\n", i, project, date, amount, investors, weburl, founder, category, subcategories, description, stages, valuation, announcement)
	}
}
