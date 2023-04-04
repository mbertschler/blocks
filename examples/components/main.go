package main

import (
	"log"
	"net/http"

	"github.com/mbertschler/blocks/html"
)

// TODO:
// -[x] cookie based session
// -[x] guiapi
// -[x] make the storage swappable
// -[x] counter demo
// -[x] turn it into a reusable component
// -[ ] extract page registration from component
// -[ ] add a second page
// -[ ] TODO MVC example?
// -[ ] https://github.com/tastejs/todomvc/blob/master/app-spec.md
//
// maybe later:
// -[ ] offer cockroachdb as storage
// -[ ] text field for your name
// -[ ] esbuild integration for more complex JS and importing

type App struct {
	DB     *DB
	Server *Server
}

func NewApp() *App {
	app := &App{}
	app.DB = NewDB()
	app.Server = NewServer(app.DB)
	return app
}

func pageLayout(main html.Block) html.Block {
	return html.Blocks{
		html.Doctype("html"),
		html.Html(nil,
			html.Head(nil,
				html.Meta(html.Charset("utf-8")),
				html.Title(nil, html.Text("Blocks")),
				html.Link(html.Rel("stylesheet").Href("https://cdn.jsdelivr.net/npm/simpledotcss@2.2.0/simple.min.css")),
				html.Link(html.Rel("stylesheet").Href("/css/main.css")),
			),
			html.Body(nil,
				main,
				html.Script(html.Src("/js/guiapi.js")),
				html.Script(html.Src("/js/main.js")),
			),
		),
	}
}

func main() {
	app := NewApp()
	app.Server.RegisterComponent(&Counter{App: app})
	app.Server.Static("/js/", "./js")
	app.Server.Static("/css/", "./css")

	log.Println("listening on localhost:8000")
	err := http.ListenAndServe("localhost:8000", app.Server.Handler())
	if err != nil {
		log.Fatal(err)
	}
}
