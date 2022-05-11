package main

import (
	"fmt"
	"testing"
)

func Test_electricList(t *testing.T) {
	urls, err := electricList("arbitrum")
	if err != nil {
		t.Fatal(err)
	}
	for i, u := range urls {
		fmt.Println(i, u)
	}
}

func Test_weblist(t *testing.T) {
	urls, err := webList("https://portal.arbitrum.one/")
	if err != nil {
		t.Fatal(err)
	}
	for i, u := range urls {
		fmt.Println(i, u)
	}
}
