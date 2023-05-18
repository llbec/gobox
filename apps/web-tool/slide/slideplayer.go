package main

import (
	"fmt"
	"gobox/src/website/slideplayer"
	"os"
)

func main() {
	dir, _ := os.Getwd()
	dir += "\\download\\"

	list := slideplayer.ManualImgList("https://player.slidesplayer.com/89/14366206/slides/slide_93.jpg", 6, 129)

	for i, u := range list {
		err := slideplayer.SaveImg(u, fmt.Sprintf("%sppt_%d.jpg", dir, i+1))
		if err != nil {
			err := slideplayer.SaveImg(u, fmt.Sprintf("%sppt_%d.jpg", dir, i+1))
			if err != nil {
				err := slideplayer.SaveImg(u, fmt.Sprintf("%sppt_%d.jpg", dir, i+1))
				if err != nil {
					fmt.Printf("[%4d] download %s failed: %v\n", i+1, u, err)
					continue
				}
			}
		}
		fmt.Printf("[%4d] download %s\n", i+1, u)
	}
}
