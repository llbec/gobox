package main

import (
	"gobox/src/website"
	"testing"
)

func TestLoadCosplay(t *testing.T) {
	loadCfg()

	resp, err := website.GetWeb(g_stURLs[0])
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var buf [1024]byte

	n, err := resp.Body.Read(buf[:])
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("read %v bytes\n", n)

	t.Log(string(buf[:]))
}
