package main

import (
	"fmt"

	"github.com/mbertschler/blocks/html"
	"github.com/mbertschler/blocks/html/attr"
)

func main() {
	root := html.Blocks{
		// Option 1: directly add an element
		html.Doctype("html"),
		html.Html(nil,
			// Option 2: struct that implements Block interface (RenderHTML() Block)
			HeadBlock{attr.Name("key").Content("super")},
			// Option 3: function that returns a Block
			BodyBlock("Hello, world! :) <br>"),
		),
	}
	out, err := html.RenderString(root)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(out)

	out, err = html.RenderMinifiedString(root)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(out)
}

type HeadBlock struct {
	attr.Attributes
}

func (h HeadBlock) RenderHTML() html.Block {
	return html.Head(nil,
		html.Meta(h.Attributes),
	)
}

func BodyBlock(in string) html.Block {
	return html.Body(nil,
		html.Main(attr.Class("main-class\" href=\"/evil/link"),
			html.H1(nil,
				html.Text(in),
				html.Br(nil),
				html.UnsafeString(in),
			),
		),
	)
}
