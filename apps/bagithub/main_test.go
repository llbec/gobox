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
