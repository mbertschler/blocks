package html

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"testing"
)

func TestHTMLTemplate(t *testing.T) {
	renderTemplate(false)
}

func BenchmarkHTMLTemplate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		renderTemplate(false)
	}
}

func TestBlocks(t *testing.T) {
	renderBlocks(false)
}

func BenchmarkBlocks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		renderBlocks(false)
	}
}

var t *template.Template

func renderTemplate(print bool) {
	const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
		{{range .Items}}<div>{{ . }}</div>{{else}}<div><strong>no rows</strong></div>{{end}}
	</body>
</html>`
	var err error
	if t == nil {
		t, err = template.New("webpage").Parse(tpl)
		if err != nil {
			log.Fatal(err)
		}
	}

	data := struct {
		Title string
		Items []string
	}{
		Title: "My page",
		Items: []string{
			"My photos",
			"My blog",
		},
	}
	var out = &bytes.Buffer{}
	err = t.Execute(out, data)
	if err != nil {
		log.Fatal(err)
	}

	noItems := struct {
		Title string
		Items []string
	}{
		Title: "My another page",
		Items: []string{},
	}

	err = t.Execute(out, noItems)
	if err != nil {
		log.Fatal(err)
	}
	if print {
		fmt.Println(out.String())
	}
}

func renderBlocks(print bool) {
	type Data struct {
		Title string
		Items []string
	}
	blocks := func(d Data) Block {
		var rows Blocks
		if len(d.Items) == 0 {
			rows.Add(Div(nil, Strong(nil, Text("no rows"))))
		} else {
			for _, e := range d.Items {
				rows.Add(Div(nil, Text(e)))
			}
		}

		return Blocks{
			Doctype("html"),
			Html(nil,
				Head(nil,
					Meta(Charset("UTF-8")),
					Title(nil, Text(d.Title)),
				),
				Body(nil, rows),
			),
		}
	}

	var out = ioutil.Discard
	data := Data{
		Title: "My page",
		Items: []string{
			"My photos",
			"My blog",
		},
	}
	err := Render(blocks(data), out)
	if err != nil {
		log.Fatal(err)
	}

	noItems := Data{
		Title: "My another page",
		Items: []string{},
	}

	err = Render(blocks(noItems), out)
	if err != nil {
		log.Fatal(err)
	}
	if print {
		// fmt.Println(out.String())
	}
}
