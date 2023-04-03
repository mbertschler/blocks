package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mbertschler/blocks/html"
)

// TODO:
// -[x] cookie based session
// -[x] guiapi
// -[x] make the storage swappable
// -[ ] offer cockroachdb as storage
// -[x] counter demo
//
// maybe later:
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

func (a *App) CounterPage(c *gin.Context) (html.Block, error) {
	sess := sessionFromContext(c)
	counter, err := a.DB.GetCounter(sess.ID)
	if err != nil {
		return nil, err
	}
	main := html.Main(nil,
		html.H1(nil, html.Text("Blocks")),
		html.P(nil, html.Text("Blocks is a framework for building web applications in Go.")),
		html.Div(html.Id("counter"),
			html.H3(nil, html.Text("Counter")),
			html.P(html.Id("count"), html.Text(fmt.Sprintf("Current count: %d", counter.Count))),
			html.Button(html.Class("ga").Attr("ga-click", "counterDecrease"), html.Text("-")),
			html.Text(" "),
			html.Button(html.Class("ga").Attr("ga-click", "counterIncrease"), html.Text("+")),
		),
	)
	return pageLayout(main), nil
}

func (a *App) CounterIncrease(c *gin.Context, args json.RawMessage) (*Response, error) {
	sess := sessionFromContext(c)
	counter, err := a.DB.GetCounter(sess.ID)
	if err != nil {
		return nil, err
	}
	counter.Count++
	err = a.DB.SetCounter(counter)
	if err != nil {
		return nil, err
	}
	return ReplaceContent("#count", html.Text(fmt.Sprintf("Current count: %d", counter.Count)))
}

func (a *App) CounterDecrease(c *gin.Context, args json.RawMessage) (*Response, error) {
	sess := sessionFromContext(c)
	counter, err := a.DB.GetCounter(sess.ID)
	if err != nil {
		return nil, err
	}
	counter.Count--
	err = a.DB.SetCounter(counter)
	if err != nil {
		return nil, err
	}
	return ReplaceContent("#count", html.Text(fmt.Sprintf("Current count: %d", counter.Count)))
}

func main() {
	app := NewApp()

	app.Server.Page("/", app.CounterPage)
	app.Server.SetFunc("counterIncrease", app.CounterIncrease)
	app.Server.SetFunc("counterDecrease", app.CounterDecrease)

	app.Server.Static("/js/", "./js")
	app.Server.Static("/css/", "./css")

	log.Println("listening on localhost:8000")
	err := http.ListenAndServe("localhost:8000", app.Server.Handler())
	if err != nil {
		log.Fatal(err)
	}
}
