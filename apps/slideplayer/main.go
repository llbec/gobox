package main

import (
	"fmt"
	"gobox/src/website/slideplayer"
)

func main() {
	list := slideplayer.ManualImgList("https://slidesplayer.com/slide/14366206//89/14366206/slides/slide_2.jpg", 5, 129)
	for i, u := range list {
		err := slideplayer.SaveImg(u, fmt.Sprintf("ppt_%d", i+1))
		if err != nil {
			fmt.Printf("[%4d]Download %s failed: %v\n", i, u, err)
		}
		fmt.Printf("[%4d]Download img %s \n", i, u)
	}
}
