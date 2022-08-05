package emodule

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

/*<CONSTRAfloat64S> 5 5 3 <CONSTRAfloat64S>
<BLOCK>
<block_A> {{0 0} {0 100} {100 100} {100 0}} <block_A>
<block_B> {{0 0} {0 50} {200 50} {200 0}} <block_B>
<block_C> {{0 0} {0 25} {120 25} {120 0}} <block_C>
<BLOCK>
<INSTANCE>
<A_1> block_A {10 10} R0 <A_1>
<A_2> block_A {210 10} MY <A_2>
<B> block_B {10 110} R0 <B>
<B/C> block_C {80 125} R0 <B/C>
<INSTANCE>
<NODE>
<A_1/N1> {60 60} <A_1/N1>
<A_2/N1> {160 60} <A_2/N1>
<B/N1> {60 135} <B/N1>
<B/C/N1> {180 135} <B/C/N1>
<NODE>
<FLY_LINE>
<A_1/N1 B/N1>
<A_1/N1 B/C/N1>
<A_2/N1 B/C/N1>
<FLY_LINE>*/

type Node struct {
	Tag      string
	HasChild bool
	Content  string
	Children []*Node
	Father   *Node
}

type XMLTree struct {
	root *Node
}

func (p *Node) addToFather(fa *Node) {
	fa.HasChild = true
	fa.Children = append(fa.Children, p)
	p.Father = fa
}

func ReadXML(reader io.Reader) XMLTree {
	root := Node{
		Tag:      "root",
		HasChild: true,
		Children: make([]*Node, 0),
		Father:   nil,
	}
	stack := make([]*Node, 0)
	stack = append(stack, &root)
	bufreader := bufio.NewReader(reader)
	contentBuilder := strings.Builder{}
	for {
		ch, _, err := bufreader.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}
		if ch == '<' {
			builder := strings.Builder{}
			for {
				nch, _, err := bufreader.ReadRune()
				if nch == '>' {
					break
				}
				if err == io.EOF {
					break
				}
				builder.WriteRune(nch)
			}
			node := Node{
				Tag: builder.String(),
			}
			if stack[len(stack)-1].Tag == node.Tag {
				node.Content = contentBuilder.String()
				node.Children = stack[len(stack)-1].Children
				stack = stack[:len(stack)-1]
				node.addToFather(stack[len(stack)-1])
				contentBuilder.Reset()
			} else {
				if stack[len(stack)-1].Tag == "FLY_LINE" {
					node.addToFather(stack[len(stack)-1])
				} else {
					stack = append(stack, &node)
				}
			}
		} else {
			contentBuilder.WriteRune(ch)
		}
	}
	return XMLTree{
		root: &root,
	}
}

func Prfloat64XML() {
	file, err := os.Open("./emodule/test.xml")
	if err != nil {
		fmt.Println(err)
	}
	xml := ReadXML(file)
	node := xml.root
	prfloat64Node(node, 0)
}
func prfloat64Node(p *Node, cnt int) {
	for i := 0; i < cnt; i++ {
		fmt.Print("    ")
	}
	fmt.Println(p.Tag + ":" + p.Content)
	fmt.Println("-------------------")
	for _, v := range p.Children {
		prfloat64Node(v, cnt+1)
	}
}

type Constrafloat64 struct {
	Height    float64
	Width     float64
	LineWidth float64
}

func NewConstrafloat64(str string) *Constrafloat64 {
	ans := Constrafloat64{}
	fmt.Sscan(str, &ans.Height, &ans.Height, &ans.Width, &ans.LineWidth)
	return &ans
}

type Block struct {
	rect     []Rect
	PointSet []Point
}

func NewBlock(str string) *Block {
	ans := Block{}
	builder := strings.Builder{}
	for _, v := range str {
		if v != '{' && v != '}' {
			builder.WriteRune(v)
		}
	}
	str = builder.String()
	numbers := make([]float64, 0)
	num := 0.0
	str = strings.TrimSpace(str)
	for _, v := range str {
		if v >= '0' && v <= '9' {
			tmp := v - '0'
			num = num * 10
			num += float64(tmp)
		} else {
			numbers = append(numbers, num)
			num = 0
		}
	}
	for i := 0; i < len(numbers); i += 2 {
		ans.PointSet = append(ans.PointSet, NewPoint(numbers[i], numbers[i+1]))
	}
	return &ans
}
