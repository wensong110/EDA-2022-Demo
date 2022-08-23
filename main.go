package main

import (
	"Solution/emodule"
	"fmt"
	"os"
)

func main() {
	file, _ := os.Open("./test.xml")
	tree := emodule.ReadXML(file)
	for _, v := range tree.Root.Children {
		if v.Tag == "BLOCK" {
			for _, block := range v.Children {
				b := emodule.NewBlock(block.Tag, block.Content)
				fmt.Println(b)
			}
		}

	}
	return
}
