package website

import (
	"net/http"
	"regexp"
)

func GetGithubUrl(doc []byte) [][]byte {
	reg := regexp.MustCompile(`\"https://github.com\S+\"`)
	urls := reg.FindAll(doc, -1)
	for i, u := range urls {
		urls[i] = urls[i][1 : len(u)-1]
	}
	return urls
}

func GetWeb(url string) (resp *http.Response, err error) {
	for i := 0; i < 3; i++ {
		resp, err = http.Get(url)
		if err == nil {
			return
		}
	}
	return
}
