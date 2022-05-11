package electric

import (
	"fmt"
	"testing"
)

func Test_Ecosystems(t *testing.T) {
	m, e := getContent()
	if e != nil {
		t.Fatal(e)
	}
	for k, v := range m {
		fmt.Println(k, v)
	}
	fmt.Println(m["5"])
}

func Test_link(t *testing.T) {
	el := NewElecInfo()
	link, err := el.getLink("arbirise")
	if err != nil {
		t.Fatalf(err.Error())
	}
	//fmt.Println(link)
	if link != "https://github.com/electric-capital/crypto-ecosystems/blob/master/data/ecosystems/a/arbirise.toml" {
		t.Fatalf("expect\nhttps://github.com/electric-capital/crypto-ecosystems/blob/master/data/ecosystems/a/arbirise.toml\nget\n%s\n", link)
	}
}

func Test_getOrg(t *testing.T) {
	//contentMap["a"] = "https://github.com/electric-capital/crypto-ecosystems/blob/master/data/ecosystems/a/arbitrum.toml"
	//linkMap["arbitrum.toml"] = "https://github.com/electric-capital/crypto-ecosystems/blob/master/data/ecosystems/a/arbitrum.toml"
	el := NewElecInfo()
	arb, err := el.GetOrg("arbirise")
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(arb)
	//fmt.Println(linkMap)
}

func Test_nameformat(t *testing.T) {
	fmt.Println(nameFormat("Arbi's Finance"))
}
