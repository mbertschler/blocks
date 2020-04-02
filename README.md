Blocks rendering library
========================
![status stable](https://img.shields.io/badge/status-stable-green?style=flat-square) [![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/mbertschler/blocks/html?tab=doc)


Blocks is a Go library for writing HTML similar to React components.

Example
-------
Applications can rendern their user interface like in this example [examples/html/html.go](./examples/html/html.go). You can run it by running `cd examples/html && go run html.go`. Full code of this demo: 

```go
package main

import (
	"fmt"

	"github.com/mbertschler/blocks/html"
)

func main() {
	root := html.Blocks{
		// Option 1: directly add an element
		html.Doctype("html"),
		html.Html(nil,
			// Option 2: struct that implements Block interface (RenderHTML() Block)
			HeadBlock{html.Name("key").Content("super")},
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
	html.Attributes
}

func (h HeadBlock) RenderHTML() html.Block {
	return html.Head(nil,
		html.Meta(h.Attributes),
	)
}

func BodyBlock(in string) html.Block {
	return html.Body(nil,
		html.Main(html.Class("main-class\" href=\"/evil/link"),
			html.H1(nil,
				html.Text(in),
				html.Br(),
				html.UnsafeString(in),
			),
		),
	)
}
```

License
-------
Blocks is released under the Apache 2.0 license. See [LICENSE](LICENSE).
