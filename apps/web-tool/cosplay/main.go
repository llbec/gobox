package main

import (
	"fmt"
	"gobox/src/website"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	g_count int
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func getTitle(doc *goquery.Document) string {
	fmt.Println(doc.Find("div.gm").Find("gd2").Html())
	return doc.Find("div.gm").Find("gd2").Find("h1#gn").Text()
}

func getImagesNumber(doc *goquery.Document) int {
	// <div class="gtb">
	div_gtb := doc.Find("div.gtb").Text()
	//fmt.Println(div_gtb)

	// 000 images
	reg := regexp.MustCompile("[0-9]+\\simages")
	if reg == nil {
		fmt.Println("regexp.MustCompile failed !!!")
	}
	n, err := strconv.Atoi(strings.Split(reg.FindStringSubmatch(div_gtb)[0], " ")[0])
	if err != nil {
		fmt.Println(err)
	}
	return n
}

func getPagesURL(doc *goquery.Document) (string, int) {
	// <div class="gtb">
	dov_gtb := doc.Find("div.gtb")
	// <td
	table := dov_gtb.Find("table.ptt").Find("tbody").Find("tr").Find("td")
	/*table.Each(func(i int, s *goquery.Selection) {
		fmt.Println(i)
		fmt.Println(s.Html())
	})*/
	/*for _, node := range table.Nodes {
		d := goquery.NewDocumentFromNode(node)
		fmt.Println(d.Html())
	}*/

	if len(table.Nodes) == 0 {
		fmt.Printf("page table length is %v, node number is %v\n", table.Length(), len(table.Nodes))
		return "", 0
	}

	tb_end := goquery.NewDocumentFromNode(table.Nodes[table.Length()-2])
	href, exist := tb_end.Find("a").Attr("href")
	if exist != true {
		fmt.Println("find last table data failed !!!")
	}

	//fmt.Println(strings.Split(href, "=")[0], tb_end.Text())
	num, err := strconv.Atoi(tb_end.Text())
	if err != nil {
		fmt.Println(err)
	}
	return strings.Split(href, "=")[0], num
}

func getImageBaseURL(doc *goquery.Document) string {
	// <div id="gdt">
	// div class="gdtm"
	// <a
	// href=""
	href, exist := doc.Find("div#gdt").Find("div.gdtm").Find("a").Attr("href") //Find("img").Attr("src")
	if exist != true {
		fmt.Println("find href failed !!!")
	}

	strs := strings.Split(href, "-")

	return strs[0] + "-" + strs[1]
}

func main() {
	loadCfg()
	g_count = 1

	resp, err := website.GetWeb(g_stURLs[0])
	if err != nil {
		myPanic(err)
	}
	defer resp.Body.Close()

	/*var buf [102400]byte

	n, err := resp.Body.Read(buf[:])
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("read %v bytes\n", n)

	fmt.Println(string(buf[:]))*/

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		myPanic(err)
	}

	//fmt.Printf("doc size is %v\n", doc.Size())
	//fmt.Println(doc.Html())

	//title := getTitle(doc)
	downpath := fmt.Sprintf("%s\\%s", g_path, g_title)
	fmt.Println(downpath)

	ok, err := PathExists(downpath)
	if err != nil {
		myPanic(err)
	}
	if !ok {
		os.MkdirAll(downpath, os.ModePerm)
	}

	// find out how many pages
	//n := getImagesNumber(doc)
	//fmt.Println(n)
	pageBaseURL, num := getPagesURL(doc)
	if num == 0 {
		fmt.Println(doc.Html())
		return
	}
	for i := 0; i < num; i++ {
		var page *goquery.Document
		if i == 0 {
			page = doc
		} else {
			respTM, err := website.GetWeb(fmt.Sprintf("%s=%d", pageBaseURL, i))
			if err != nil {
				myPanic(err)
			}
			defer respTM.Body.Close()
			page, err = goquery.NewDocumentFromReader(respTM.Body)
			if err != nil {
				myPanic(err)
			}
		}
		download(downpath, page)
	}
}

func download(downpath string, doc *goquery.Document) {
	doc.Find("div#gdt").Find("div.gdtm").Each(func(i int, s *goquery.Selection) {
		if g_count < g_start {
			return
		}
		href, exist := s.Find("a").Attr("href")
		if !exist {
			myPanic("not find href")
		}
		resp, err := website.GetWeb(href)
		if err != nil {
			myPanic(err)
		}
		defer resp.Body.Close()

		imgDoc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			myPanic(err)
		}

		// <div id="i3">
		src, exist := imgDoc.Find("div#i1").Find("div#i3").Find("a").Find("img").Attr("src")
		if exist != true {
			myPanic("find href failed !!!")
		}

		// Get the data
		respSrc, err := website.GetWeb(src)
		if err != nil {
			myPanic(err)
		}
		defer respSrc.Body.Close()

		// 创建一个文件用于保存
		out, err := os.Create(fmt.Sprintf("%s\\%s%d.jpg", downpath, g_title, g_count))
		if err != nil {
			myPanic(err)
		}
		defer out.Close()

		// 然后将响应流和文件流对接起来
		_, err = io.Copy(out, respSrc.Body)
		if err != nil {
			myPanic(err)
		}
		fmt.Println("download image", g_count)
		g_count++
	})
}

/*
	// find out image links
	baseURL := getImageBaseURL(doc)
	fmt.Println("baseurl ", baseURL)
	for i := 1; i <= n; i++ {
		respImage, err := website.GetWeb(fmt.Sprintf("%s-%d", baseURL, i))
		if err != nil {
			myPanic(err)
		}
		defer respImage.Body.Close()
		imgDoc, err := goquery.NewDocumentFromReader(respImage.Body)
		if err != nil {
			myPanic(err)
		}

		// <div id="i3">
		src, exist := imgDoc.Find("div#i1").Find("div#i3").Find("a").Find("img").Attr("src")
		if exist != true {
			myPanic("find href failed !!!")
		}

		// Get the data
		respSrc, err := website.GetWeb(src)
		if err != nil {
			myPanic(err)
		}
		defer respSrc.Body.Close()

		// 创建一个文件用于保存
		out, err := os.Create(fmt.Sprintf("%s\\%s%d.jpg", downpath, g_title, i))
		if err != nil {
			myPanic(err)
		}
		defer out.Close()

		// 然后将响应流和文件流对接起来
		_, err = io.Copy(out, respSrc.Body)
		if err != nil {
			myPanic(err)
		}
	}
*/
