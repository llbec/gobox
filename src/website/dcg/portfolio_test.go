package dcg

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func Test_loadPortfolio(t *testing.T) {
	left, err := loadPortfolio("https://dcg.co/portfolio/")
	if err != nil {
		t.Fatal(err)
	}
	for i, n := range left.Children().Nodes {
		q := goquery.NewDocumentFromNode(n)
		if q.Is("h2.header") {
			fmt.Println(i, q.Text())
		} else {
			name := strings.TrimSpace(q.Find("div.name").Text())
			ds := strings.Split(q.Find("div.details").Text(), "/")
			details := strings.TrimSpace(ds[0])
			location := strings.TrimSpace(ds[1])
			dp := q.Find("div.description")
			description := strings.TrimSpace(dp.Find("p").Text())
			weburl := strings.TrimSpace(dp.Find("a").Text())
			fmt.Printf("%d,\nname:%s\ndetails:%s\nlocation:%s\ndescription:%s\nurl:%s\n",
				i, name, details, location, description, weburl,
			)
		}
		//s, e := q.Html()
		//fmt.Println(i, s, e)
	}
}
