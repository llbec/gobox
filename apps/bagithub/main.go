package main

import (
	"flag"
	"fmt"
	"gobox/src/website"
	electric "gobox/src/website/electricCapital"
	"io/ioutil"
	"os"
)

var (
	fFile string
	fWeb  string
	fName string
)

func init() {
	flag.StringVar(&fFile, "f", "", "specify file dir")
	flag.StringVar(&fWeb, "w", "", "specify website url")
	flag.StringVar(&fName, "n", "", "specify orgnazition name")
}

func main() {
	flag.Parse()
	if fWeb != "" {
		fmt.Println(fWeb, "start!")
		urls, err := webList(fWeb)
		if err != nil {
			fmt.Println(err)
			return
		}
		outputList(urls)
	}
	if fFile != "" {
		fmt.Println(fFile, "start!")
		urls, err := fileList(fFile)
		if err != nil {
			fmt.Println(err)
			return
		}
		outputList(urls)
	}
	if fName != "" {
		fmt.Println(fName, "start!")
		urls, err := electricList(fName)
		if err != nil {
			fmt.Println(err)
			return
		}
		outputList(urls)
	}
	flag.Usage()
}

func webList(url string) (urls []string, err error) {
	resp, err := website.GetWeb(fWeb)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	urls = website.GetGithubUrl(body)
	return
}

func fileList(fileDir string) (urls []string, err error) {
	doc, err := ioutil.ReadFile(fFile)
	if err != nil {
		return
	}
	urls = website.GetGithubUrl(doc)
	return
}

func electricList(name string) (urls []string, err error) {
	elec := electric.NewElecInfo()

	org, err := elec.GetOrg(name)
	if err != nil {
		return nil, err
	}
	subs := org.GetSubs()
	for _, s := range subs {
		url := s.GetGithub()
		if url != "" {
			urls = append(urls, url)
		}
	}
	return
}

func outputList(list []string) {
	for i, itm := range list {
		fmt.Printf("%4d,\t%s\n", i, itm)
	}
	os.Exit(0)
}
