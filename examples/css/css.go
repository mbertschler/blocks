package main

import (
	"fmt"

	"github.com/mbertschler/blocks/css"
)

func main() {
	root := css.Blocks{
		css.New("body",
			css.L("position", "relative"),
			css.L("display", "block"),
		),
		HeadBlock{size: 32},
		ParagraphBlock("sans-serif"),
	}
	out, err := css.RenderString(root)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(out)

	out, err = css.RenderMinifiedString(root)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(out)
}

type HeadBlock struct {
	size int
}

func (h HeadBlock) RenderCSS() css.Block {
	return css.New("h1",
		css.L("font-size", fmt.Sprint(h.size, "px")),
	)
}

func ParagraphBlock(font string) css.Block {
	return css.New("p",
		css.L("font-family", font),
	)

}
