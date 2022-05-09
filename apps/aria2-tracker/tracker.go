package main

import (
	"bufio"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

//var listUrl = "https://ngosang.github.io/trackerslist/trackers_best.txt"
var listUrl = "https://ngosang.github.io/trackerslist/master/trackers_all.txt"

var tmpl = `bt-tracker={{range $k, $v := .}}{{if eq $k 0}}{{$v}}{{else}},{{$v}}{{end}}{{end}}`

//enable-dht=true
//bt-enable-lpd=true
//enable-peer-exchange=true

func downloadList() ([]string, error) {
	resp, err := http.Get(listUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	rd := bufio.NewReader(resp.Body)
	res := []string{}
	for {
		line1, _, err := rd.ReadLine()
		if err != nil {
			break
		}
		//fmt.Println("url:", string(line1), len(line1))
		if len(line1) > 0 {
			res = append(res, string(line1))
		}
	}
	return res, nil
}

func writeConf(data []string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dir1 := filepath.Join(home, ".aria2")
	os.Mkdir(dir1, os.ModePerm)
	confPath := filepath.Join(home, ".aria2", "track.conf")
	fp, err := os.Create(confPath)
	if err != nil {
		return err
	}
	defer fp.Close()
	t := template.New("")
	t, err = t.Parse(tmpl)
	if err != nil {
		return err
	}
	return t.Execute(fp, data)
}

func main() {
	/*flag.Parse()
	if len(flag.Arg(0)) > 0 {
		listUrl = flag.Arg(0)
	} else {
		fmt.Print("Usage:\n\t./aria2-trackers\nOr:\n\t./aria2-trackers [url of list]\n")
	}*/
	list1, err := downloadList()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(list1)
	err = writeConf(list1)
	if err != nil {
		panic(err.Error())
	}
}
