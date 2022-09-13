package main

import (
	"Solution/emodule"
	"Solution/evisiable"
	"fmt"
	"os"
	"strings"
)

func main() {
	file, _ := os.Open("./test.xml")
	tree := emodule.ReadXML(file)
	canvas := evisiable.NewCanvas(500, 500)
	for _, v := range tree.Root.Children {
		if v.Tag == "BLOCK" {
			for _, block := range v.Children {
				emodule.NewBlock(block.Tag, block.Content)
				//fmt.Println(b)
				// b.Draw(canvas)
				// canvas.Save(b.Name + ".png")
				// canvas.Clear()
			}
		}
		if v.Tag == "INSTANCE" {
			for _, instance := range v.Children {
				_, err := emodule.NewInstance(instance.Tag, instance.Content)
				if err != nil {
					fmt.Println(err)
				}
				//fmt.Println(b)
			}
		}
		if v.Tag == "NODE" {
			for _, node := range v.Children {
				b, err := emodule.NewNode(node.Tag, node.Content)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(b)
			}
		}
		if v.Tag == "FLY_LINE" {
			for _, flyLine := range v.Children {
				_, err := emodule.NewFlyLine(flyLine.Tag)
				if err != nil {
					fmt.Println(err)
				}
				//fmt.Println(b)
			}
		}
	}
	for _, v := range emodule.InstanceSet {
		//fmt.Println(i, v.Name)
		if !strings.ContainsAny(v.Name, "/") {
			v.Draw(canvas)
			canvas.Save(v.Name + ".png")
			//canvas.Clear()
		}
	}
	return
}
