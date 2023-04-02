package main

import (
	"log"
	"net/http"

	"github.com/mbertschler/blocks/html"
)

// TODO:
// -[ ] cookie based session
// -[ ] guiapi
// -[ ] counter demo
//
// maybe later:
// -[ ] text field for your name
// -[ ] esbuild integration for more complex JS and importing

func rootHandler(w http.ResponseWriter, r *http.Request) {
	page := html.Blocks{
		html.Doctype("html"),
		html.Html(nil,
			html.Head(nil,
				html.Meta(html.Charset("utf-8")),
				html.Title(nil, html.Text("Blocks")),
				html.Link(html.Rel("stylesheet").Href("https://cdn.jsdelivr.net/npm/simpledotcss@2.2.0/simple.min.css")),
				html.Link(html.Rel("stylesheet").Href("/css/main.css")),
			),
			html.Body(nil,
				html.Main(nil,
					html.H1(nil, html.Text("Blocks")),
					html.P(nil, html.Text("Blocks is a framework for building web applications in Go.")),
					html.H3(nil, html.Text("Counter")),
					html.P(nil, html.Text("Current count: 0")),
					html.Button(nil, html.Text("-")),
					html.Text(" "),
					html.Button(nil, html.Text("+")),
				),
				html.Script(html.Src("/js/main.js")),
			),
		),
	}
	err := html.RenderMinified(w, page)
	if err != nil {
		log.Println("error during page rendering:", err)
	}
}

func main() {
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.HandleFunc("/", rootHandler)
	log.Println("listening on localhost:8000")
	err := http.ListenAndServe("localhost:8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
