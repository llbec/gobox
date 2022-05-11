package electric

import (
	"fmt"
	"gobox/src/website"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	contentMap map[string]string
	//linkMap    map[string]string
	//orgMap     map[string]*Organization
	//cacheOrg   *Organization
	baseURL string
)

func init() {
	contentMap = make(map[string]string)
	//linkMap = make(map[string]string)
	//orgMap = make(map[string]*Organization)
	baseURL = "https://github.com"
}

func getContent() (map[string]string, error) {
	resp, err := website.GetWeb("https://github.com/electric-capital/crypto-ecosystems/tree/master/data/ecosystems")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	for k := range contentMap {
		delete(contentMap, k)
	}

	doc.Find("div.js-navigation-item").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}
		a := s.Find("a")
		k, ok := a.Attr("title")
		if !ok {
			return
		}
		v, ok := a.Attr("href")
		if !ok {
			return
		}
		contentMap[k] = baseURL + v
	})
	return contentMap, nil
}

func (elec *ElecInfo) getArchive(name string) (string, error) {
	name = nameFormat(name)
	arc := elec.ArchiveMap[name[0:1]]
	if arc == "" {
		return "", fmt.Errorf("have no archive about %s(%s)", name[0:1], name)
	}
	return arc, nil
}

func (elec *ElecInfo) getLink(name string) (string, error) {
	name = nameFormat(name) + ".toml"
	if elec.linkMap[name] != "" {
		return elec.linkMap[name], nil
	}
	url, err := elec.getArchive(name)
	if err != nil {
		return "", err
	}
	resp, err := website.GetWeb(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}
	doc.Find("div.js-navigation-item").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}
		a := s.Find("a")
		k, ok := a.Attr("title")
		if !ok {
			return
		}
		v, ok := a.Attr("href")
		if !ok {
			return
		}
		elec.linkMap[k] = baseURL + v
	})
	if elec.linkMap[name] == "" {
		return "", fmt.Errorf("Page(%s) without %s", url, name)
	}
	return elec.linkMap[name], nil
}

func (elec *ElecInfo) GetOrg(name string) (*Organization, error) {
	name = nameFormat(name)

	if org := elec.Orgs[name]; org != nil {
		return org, nil
	}

	link, err := elec.getLink(name)
	if err != nil {
		return nil, err
	}
	resp, err := website.GetWeb(link)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	tbody := doc.Find("tbody")
	l := len(doc.Find("tbody").Find("tr").Nodes)

	newOrg := &Organization{}
	index := 2 //title
	str := getPLsText(tbody, index)
	str = str[1 : len(str)-1]
	newOrg.SetName(nameFormat(str))

	index = 4 //subsysterms
	if v := getPLsText(tbody, index); v != "" {
		v = v[1 : len(v)-1]
		newOrg.AddSub(v, elec)
	}
	index = 5
	for ; index <= l; index++ {
		//fmt.Println(td.Find("span.pl-smi").Html())
		v := getPLsText(tbody, index)
		if v == "" {
			break
		}
		v = v[1 : len(v)-1]
		newOrg.AddSub(v, elec)
	}

	index += 1 //github_organizations
	//fmt.Println(index, "org git")
	for ; index <= l; index++ {
		if getPLsmiText(tbody, index) == "github_organizations" {
			if v := getPLsText(tbody, index); v != "" {
				v = v[1 : len(v)-1]
				newOrg.SetGithub(v)
			}
		}
	}

	elec.Orgs[name] = newOrg

	return newOrg, nil
}

func getPLsmiText(doc *goquery.Selection, i int) string {
	td := doc.Find(fmt.Sprintf("td#LC%d", i))
	s, _ := td.Find("span.pl-smi").Html()
	return s
}

func getPLsText(doc *goquery.Selection, i int) string {
	td := doc.Find(fmt.Sprintf("td#LC%d", i))
	//fmt.Println(td.Html())
	s := td.Find("span.pl-s").Text()
	//fmt.Println(i, s)
	return s
}

func nameFormat(name string) string {
	str := strings.ToLower(name)
	reg := regexp.MustCompile("[^a-z0-9]")
	return reg.ReplaceAllString(str, "-")
}
