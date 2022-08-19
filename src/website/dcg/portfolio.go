package dcg

import (
	"gobox/src/website"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	Portfolio []*Company
	Map       map[string]*Company
)

func init() {
	Portfolio = make([]*Company, 0)
}

func loadPortfolio(url string) (*goquery.Selection, error) {
	resp, err := website.GetWeb(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	return doc.Find("div.sector").Find("div.left"), nil
}

func FillPortfolio(url string) error {
	left, err := loadPortfolio(url)
	if err != nil {
		return err
	}
	sector := ""
	Portfolio = Portfolio[:0]
	Map = make(map[string]*Company)
	for _, n := range left.Children().Nodes {
		q := goquery.NewDocumentFromNode(n)
		if q.Is("h2.header") {
			sector = q.Text()
		} else {
			comp := getCompany(sector, q)
			Portfolio = append(Portfolio, comp)
			Map[comp.Name] = comp
		}
	}
	return nil
}

func getCompany(sector string, company *goquery.Document) *Company {
	name := strings.TrimSpace(company.Find("div.name").Text())
	ds := strings.Split(company.Find("div.details").Text(), "/")
	details := strings.TrimSpace(ds[0])
	location := strings.TrimSpace(ds[1])
	dp := company.Find("div.description")
	description := strings.TrimSpace(dp.Find("p").Text())
	weburl := strings.TrimSpace(dp.Find("a").Text())
	return &Company{
		Name:        name,
		Details:     details,
		Location:    location,
		Description: description,
		Url:         weburl,
		Sector:      sector,
	}
}
