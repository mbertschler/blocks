package main

import (
	"fmt"

	ht "git.exahome.net/tools/blocks/html"
)

func main() {
	root := ht.Blocks{
		// Option 1: directly add an element
		ht.Doctype(ht.Attr{{"html", nil}}),
		// Option 2: struct that implements Block interface (RenderHTML() Block)
		HeadBlock{ht.Attr{{"key", "key"}, {"value", "super"}}},
		// Option 3: function that returns a Block
		BodyBlock("Hello, world! :)"),
	}
	out, err := ht.RenderString(root)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Print(out)

	out, err = ht.RenderMinifiedString(root)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(out)
}

type HeadBlock struct {
	ht.Attr
}

func (h HeadBlock) RenderHTML() ht.Block {
	return ht.Head(ht.NoAttr,
		ht.Meta(h.Attr),
	)
}

func BodyBlock(in string) ht.Block {
	return ht.Body(ht.NoAttr,
		ht.Main(ht.Attr{{"class", "main-class"}},
			ht.H1(ht.NoAttr,
				ht.Text(in),
			),
		),
	)
}
