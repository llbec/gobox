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
	fCalc bool
)

func init() {
	flag.StringVar(&fFile, "f", "", "specify file dir")
	flag.StringVar(&fWeb, "w", "", "specify website url")
	flag.StringVar(&fName, "n", "", "specify orgnazition name")
	flag.BoolVar(&fCalc, "c", false, "calc list add or sub")
}

func main() {
	flag.Parse()
	if fCalc {
		var urls []string
		for {
			a := ""
			t := ""
			v := ""
			fmt.Printf("Enter add or sub, others quit\n")
			fmt.Scanln(&a)
			if a != "add" && a != "sub" {
				break
			}
			fmt.Printf("Enter type(org,web,file) and value path(%v)\n", t)
			fmt.Scanln(&t, &v)
			urls = handleCalc(urls, a, t, v)
		}
		outputList(urls)
	}
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

func handleCalc(src []string, act, v_type, v string) (urls []string) {
	var tmp []string
	var err error
	switch v_type {
	case "web":
		tmp, err = webList(v)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "file":
		tmp, err = fileList(v)
		if err != nil {
			fmt.Println(err)
			return
		}
	default:
		tmp, err = electricList(v)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	srcmap := make(map[string]bool)
	for _, v := range src {
		v = formatUrl(v)
		srcmap[v] = true
	}
	for _, v := range tmp {
		v = formatUrl(v)
		if act == "add" {
			srcmap[v] = true
		} else {
			srcmap[v] = false
		}
	}
	for k, v := range srcmap {
		if v {
			urls = append(urls, k)
		}
	}
	fmt.Printf("src's length is %d, dest is %d\n", len(src), len(urls))
	return
}

func webList(url string) (urls []string, err error) {
	resp, err := website.GetWeb(url)
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
	doc, err := ioutil.ReadFile(fileDir)
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
	fmt.Printf("Total %d\n", len(list))
	for _, itm := range list {
		//fmt.Printf("%4d,\t%s\n", i, itm)
		fmt.Printf("%s\n", itm)
	}
	os.Exit(0)
}

func formatUrl(str string) string {
	l := len(str)
	if str[l-1:l] == "/" {
		str = str[:l-1]
	}
	return str
}
