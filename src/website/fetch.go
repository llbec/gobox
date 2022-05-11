package website

import (
	"fmt"
	"net/http"
	"regexp"
)

func GetGithubUrl(doc []byte) []string {
	reg := regexp.MustCompile(`\"https://github.com\S+\"`)
	urls := reg.FindAll(doc, -1)
	var ul []string
	for i, u := range urls {
		ul = append(ul, string(urls[i][1:len(u)-1]))
	}
	return ul
}

func GetWeb(url string) (resp *http.Response, err error) {
	for i := 0; i < 3; i++ {
		resp, err = http.Get(url)
		if err == nil && resp.StatusCode == 200 {
			return
		}
	}
	if err == nil {
		err = fmt.Errorf("get error: %v", resp.Status)
	}
	return
}
