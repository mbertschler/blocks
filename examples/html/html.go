package main

import (
	"fmt"

	"git.exahome.net/tools/blocks/html"
)

func main() {
	root := html.Blocks{
		// Option 1: directly add an element
		html.Doctype(html.Attr{{"html", nil}}),
		// Option 2: struct that implements Block interface (RenderHTML() Block)
		HeadBlock{html.Attr{{"key", "key"}, {"value", "super"}}},
		// Option 3: function that returns a Block
		BodyBlock("Hello, world! :)"),
	}
	out, err := html.RenderString(root)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Print(out)

	out, err = html.RenderMinifiedString(root)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(out)
}

type HeadBlock struct {
	html.Attr
}

func (h HeadBlock) RenderHTML() html.Block {
	return html.Head(html.NoAttr,
		html.Meta(h.Attr),
	)
}

func BodyBlock(in string) html.Block {
	return html.Body(html.NoAttr,
		html.Main(html.Attr{{"class", "main-class"}},
			html.H1(html.NoAttr,
				html.Text(in),
			),
		),
	)
}
