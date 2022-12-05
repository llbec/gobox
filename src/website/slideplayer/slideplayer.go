package slideplayer

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

var (
	doc     *goquery.Document
	baseurl = "https://slidesplayer.com/slide/14366206"
)

func LoadDocuments(file string) (err error) {
	fileObj, err := os.Open("test.html")
	if err != nil {
		return
	}
	if d, err := goquery.NewDocumentFromReader(fileObj); err == nil {
		doc = d
	}
	return
}

func AutoImageList() (list []string, err error) {
	layout := doc.Find("div#player_layout")
	imgNum := len(layout.Children().Nodes)
	slide := doc.Find("div#slide_1")
	src, ok := slide.Find("img").Attr("src")
	if !ok {
		return list, fmt.Errorf("div id=slide_1 img not find")
	}

	///89/14366206/slides/slide_1.jpg delete 1.jpg
	imgBase := baseurl + src[:len(src)-5]

	for i := 0; i < imgNum; i++ {
		list = append(list, fmt.Sprintf("%s%d.jpg", imgBase, i+1))
	}
	return
}

func ManualImgList(headUrl string, offset int, num int) (list []string) {
	base := headUrl[0 : len(headUrl)-offset]

	for i := 0; i < num; i++ {
		list = append(list, fmt.Sprintf("%s%d.jpg", base, i+1))
	}
	return
}

func SaveImg(url string, name string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 创建一个文件用于保存
	out, err := os.Create(name)
	if err != nil {
		return err
	}
	defer out.Close()

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
