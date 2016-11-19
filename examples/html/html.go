package main

import (
	"fmt"

	"git.exahome.net/tools/blocks/html"
)

func main() {
	root := html.Blocks{
		html.Doctype(html.Attr{{"html", nil}}),
		HeadBlock{html.Attr{{"key", "key"}, {"value", "super"}}},
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
	fmt.Print(out)
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
