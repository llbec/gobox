package getwebsite

import (
	"fmt"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

var (
	testFile = "test.html"
)

func Test_Fetch(t *testing.T) {
	fileObj, err := os.Open(testFile)
	if err != nil {
		t.Fatal(err)
	}
	doc, err := goquery.NewDocumentFromReader(fileObj)
	if err != nil {
		t.Fatal(err)
	}
	//client-box-out shown uk-first-column uk-scrollspy-inview uk-animation-fade
	doc.Find("div.uk-animation-fade").Each(func(i int, selection *goquery.Selection) {
		box := selection.Find("div.text-box")
		in := box.Find("div.text-box-in")
		name, _ := in.Find(".client-name").Html()
		//objType, _ := in.Find(".client-type").Html()
		//socials := in.Find("div.socials>a")
		socials := in.Find("div.socials")
		github := socials.Find("a.github")
		url, _ := github.Attr("href")
		/*giturl := "nil"
		for _, s := range socials.Nodes {
			if strings.Contains(s.Attr[0].Val, "github") {
				giturl = s.Attr[1].Val
			}
		}*/
		fmt.Println(i, name, url)
	})
}
