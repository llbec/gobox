package slideplayer

import (
	"fmt"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func Test_SlideDoc(t *testing.T) {
	fileObj, err := os.Open("test.html")
	if err != nil {
		t.Fatal(err)
	}
	doc, err := goquery.NewDocumentFromReader(fileObj)
	if err != nil {
		t.Fatal(err)
	}

	/*
		<div id="player_layout" style="width: 800px; height: 600px;">
			<div id="slide_1" class="">
				<img src="/89/14366206/slides/slide_1.jpg" width="0" height="0" style="width: 800px; height: 600px;">
			</div>
			<div id="slide_2" class="hidden_slides">
				<img src="/89/14366206/slides/slide_2.jpg" width="0" height="0" style="width: 800px; height: 600px;">
			</div>
		</div>
	*/

	//layout := doc.Find("div#player_layout") //. class; # id
	//pageNum := len(layout.Children().Nodes)

	slide := doc.Find("div#slide_1")
	fmt.Println(slide.Html())

	fmt.Println(slide.Find("img").Attr("src"))
}

func Test_ImageList(t *testing.T) {
	LoadDocuments("test.html")
	list, err := AutoImageList()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Print(list)
}

func Test_Download(t *testing.T) {
	LoadDocuments("test.html")
	list, err := AutoImageList()
	if err != nil {
		t.Fatal(err)
	}
	for i, u := range list {
		err := SaveImg(u, fmt.Sprintf("ppt_%d", i+1))
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("[%4d]Download img %s \n", i, u)
	}
}
