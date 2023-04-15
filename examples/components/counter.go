package main

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mbertschler/blocks/html"
)

type Counter struct {
	*App
}

func (c *Counter) Component() *ComponentConfig {
	return &ComponentConfig{
		Name: "Counter",
		Actions: map[string]Callable{
			"Increase": c.Increase,
			"Decrease": c.Decrease,
		},
	}
}

func counterLayout(main html.Block) html.Block {
	return html.Blocks{
		html.Doctype("html"),
		html.Html(nil,
			html.Head(nil,
				html.Meta(html.Charset("utf-8")),
				html.Title(nil, html.Text("Blocks")),
				html.Link(html.Rel("stylesheet").Href("https://cdn.jsdelivr.net/npm/simpledotcss@2.2.0/simple.min.css")),
				html.Link(html.Rel("stylesheet").Href("/dist/bundle.css")),
			),
			html.Body(nil,
				main,
				html.A(html.Href("/"), html.Text("TodoMVC Example")),
				html.Script(html.Src("/dist/bundle.js")),
			),
		),
	}
}

func (c *Counter) RenderPage(ctx *gin.Context) (html.Block, error) {
	block, err := c.RenderBlock(ctx)
	if err != nil {
		return nil, err
	}
	main := html.Main(nil,
		html.H1(nil, html.Text("Blocks")),
		html.P(nil, html.Text("Blocks is a framework for building web applications in Go.")),
		block,
	)
	return counterLayout(main), nil
}

func (c *Counter) RenderBlock(ctx *gin.Context) (html.Block, error) {
	sess := sessionFromContext(ctx)
	counter, err := c.App.DB.GetCounter(sess.ID)
	if err != nil {
		return nil, err
	}
	block := html.Div(html.Id("counter"),
		html.H3(nil, html.Text("Counter")),
		html.P(html.Id("count"), html.Text(fmt.Sprintf("Current count: %d", counter.Count))),
		html.Button(html.Class("ga").Attr("ga-on", "click").Attr("ga-action", "Counter.Decrease"), html.Text("-")),
		html.Text(" "),
		html.Button(html.Class("ga").Attr("ga-on", "click").Attr("ga-action", "Counter.Increase"), html.Text("+")),
	)
	return block, nil
}

func (c *Counter) Increase(ctx *gin.Context, args json.RawMessage) (*Response, error) {
	sess := sessionFromContext(ctx)
	counter, err := c.App.DB.GetCounter(sess.ID)
	if err != nil {
		return nil, err
	}
	counter.Count++
	err = c.App.DB.SetCounter(counter)
	if err != nil {
		return nil, err
	}
	return ReplaceContent("#count", html.Text(fmt.Sprintf("Current count: %d", counter.Count)))
}

func (c *Counter) Decrease(ctx *gin.Context, args json.RawMessage) (*Response, error) {
	sess := sessionFromContext(ctx)
	counter, err := c.App.DB.GetCounter(sess.ID)
	if err != nil {
		return nil, err
	}
	counter.Count--
	err = c.App.DB.SetCounter(counter)
	if err != nil {
		return nil, err
	}
	return ReplaceContent("#count", html.Text(fmt.Sprintf("Current count: %d", counter.Count)))
}
