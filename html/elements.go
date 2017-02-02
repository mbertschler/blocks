package html

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
)

var renderDebug = true

type Text string
type Comment string
type CSS template.CSS
type HTML template.HTML
type HTMLAttr template.HTMLAttr
type JS template.JS
type JSStr template.JSStr
type URL template.URL

func (Text) RenderHTML() Block       { return nil }
func (Comment) RenderHTML() Block    { return nil }
func (CSS) RenderHTML() Block        { return nil }
func (HTML) RenderHTML() Block       { return nil }
func (HTMLAttr) RenderHTML() Block   { return nil }
func (JS) RenderHTML() Block         { return nil }
func (b JS) renderString() string    { return string(b) }
func (JSStr) RenderHTML() Block      { return nil }
func (b JSStr) renderString() string { return string(b) }
func (URL) RenderHTML() Block        { return nil }

var indentation = 2

type Blocks []Block

func (b *Blocks) Add(block Block) {
	*b = append(*b, block)
}

func (b *Blocks) AddBlocks(blocks Blocks) {
	*b = append(*b, blocks...)
}

func (Blocks) RenderHTML() Block { return nil }

type Block interface {
	RenderHTML() Block
}

type stringRenderer interface {
	renderString() string
}

func Render(root Block, w io.Writer) error {
	err := renderHTML(root, w, &renderCtx{})
	if err != nil {
		return err
	}
	return nil
}

func RenderMinified(root Block, w io.Writer) error {
	err := renderHTML(root, w, &renderCtx{minified: true})
	if err != nil {
		return err
	}
	return nil
}

func RenderString(root Block) (string, error) {
	buf := bytes.Buffer{}
	err := renderHTML(root, &buf, &renderCtx{})
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func RenderMinifiedString(root Block) (string, error) {
	buf := bytes.Buffer{}
	err := renderHTML(root, &buf, &renderCtx{minified: true})
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

type renderCtx struct {
	level    int
	item     int
	minified bool
}

func (c *renderCtx) enter() (item int) {
	item = c.item
	c.level++
	c.item = 0
	return item
}

func (c *renderCtx) next() {
	c.item++
}

func (c *renderCtx) exit(item int) {
	c.level--
	c.item = item
}

func renderHTML(c Block, w io.Writer, ctx *renderCtx) error {
	//var item int
	switch el := c.(type) {
	case Text:
		if !ctx.minified {
			w.Write(bytes.Repeat([]byte{' '}, ctx.level*indentation))
		}
		w.Write([]byte(el + " "))
		if !ctx.minified {
			w.Write([]byte{'\n'})
		}
		ctx.next()
	case stringRenderer:
		if !ctx.minified {
			w.Write(bytes.Repeat([]byte{' '}, ctx.level*indentation))
		}
		w.Write([]byte(el.renderString()))
		if !ctx.minified {
			w.Write([]byte{'\n'})
		}
		ctx.next()
	case Comment:
		if !ctx.minified {
			w.Write(bytes.Repeat([]byte{' '}, ctx.level*indentation))
		}
		w.Write([]byte("<!--" + el + "-->"))
		if !ctx.minified {
			w.Write([]byte{'\n'})
		}
		ctx.next()
	case Element:
		if !ctx.minified {
			w.Write(bytes.Repeat([]byte{' '}, ctx.level*indentation))
		}
		var attr string
		for _, v := range el.Attr {
			if v.Value == nil {
				attr += " " + v.Key
				continue
			}
			attr += " " + v.Key + "=" + fmt.Sprint("\"", v.Value, "\"")
		}
		w.Write([]byte("<" + el.Type + attr))
		if el.Options&SelfClose != 0 {
			w.Write([]byte("/>"))
		} else {
			w.Write([]byte(">"))
		}
		if len(el.Children) > 0 {
			if !ctx.minified {
				w.Write([]byte{'\n'})
			}
			item := ctx.enter()
			var min bool
			if el.Options&NoWhitespace != 0 {
				min = ctx.minified
				ctx.minified = true
			}
			for _, e := range el.Children {
				renderHTML(e, w, ctx)
			}
			if el.Options&NoWhitespace != 0 {
				ctx.minified = min
			}
			ctx.exit(item)
		}
		if el.Options&Void+el.Options&SelfClose == 0 {
			if !ctx.minified {
				w.Write(bytes.Repeat([]byte{' '}, ctx.level*indentation))
			}
			w.Write([]byte("</" + el.Type + ">"))
		}
		if !ctx.minified {
			w.Write([]byte{'\n'})
		}
		ctx.next()
	case Blocks:
		for _, e := range el {
			renderHTML(e, w, ctx)
		}
	case Block:
		c := el.RenderHTML()
		renderHTML(c, w, ctx)
	default:
		if !ctx.minified {
			w.Write(bytes.Repeat([]byte{' '}, ctx.level*indentation))
		}
		fmt.Fprintf(w, "{{ ERROR value=%#v\n }}", c)
		if !ctx.minified {
			w.Write([]byte{'\n'})
		}
		ctx.next()
	}
	return nil
}

type Attr []AttrPair
type AttrPair struct {
	Key   string
	Value interface{}
}

func Id(id interface{}) Attr {
	return Attr{AttrPair{Key: "id", Value: id}}
}
func Class(class interface{}) Attr {
	return Attr{AttrPair{Key: "class", Value: class}}
}
func Href(href interface{}) Attr {
	return Attr{AttrPair{Key: "href", Value: href}}
}
func Rel(rel interface{}) Attr {
	return Attr{AttrPair{Key: "rel", Value: rel}}
}
func Name(name interface{}) Attr {
	return Attr{AttrPair{Key: "name", Value: name}}
}
func Content(name interface{}) Attr {
	return Attr{AttrPair{Key: "content", Value: name}}
}
func Defer() Attr {
	return Attr{AttrPair{Key: "defer", Value: nil}}
}
func Src(src interface{}) Attr {
	return Attr{AttrPair{Key: "src", Value: src}}
}
func Action(action interface{}) Attr {
	return Attr{AttrPair{Key: "action", Value: action}}
}
func Method(method interface{}) Attr {
	return Attr{AttrPair{Key: "method", Value: method}}
}
func Type(typ interface{}) Attr {
	return Attr{AttrPair{Key: "type", Value: typ}}
}
func For(fo interface{}) Attr {
	return Attr{AttrPair{Key: "for", Value: fo}}
}
func Value(value interface{}) Attr {
	return Attr{AttrPair{Key: "value", Value: value}}
}

func (a Attr) Id(id interface{}) Attr {
	return append(a, AttrPair{Key: "id", Value: id})
}
func (a Attr) Class(class interface{}) Attr {
	return append(a, AttrPair{Key: "class", Value: class})
}
func (a Attr) Href(href interface{}) Attr {
	return append(a, AttrPair{Key: "href", Value: href})
}
func (a Attr) Rel(rel interface{}) Attr {
	return append(a, AttrPair{Key: "rel", Value: rel})
}
func (a Attr) Name(name interface{}) Attr {
	return append(a, AttrPair{Key: "name", Value: name})
}
func (a Attr) Content(name interface{}) Attr {
	return append(a, AttrPair{Key: "content", Value: name})
}
func (a Attr) Defer() Attr {
	return append(a, AttrPair{Key: "defer", Value: nil})
}
func (a Attr) Src(src interface{}) Attr {
	return append(a, AttrPair{Key: "src", Value: src})
}
func (a Attr) Action(action interface{}) Attr {
	return append(a, AttrPair{Key: "action", Value: action})
}
func (a Attr) Method(method interface{}) Attr {
	return append(a, AttrPair{Key: "method", Value: method})
}
func (a Attr) Type(typ interface{}) Attr {
	return append(a, AttrPair{Key: "type", Value: typ})
}
func (a Attr) For(fo interface{}) Attr {
	return append(a, AttrPair{Key: "for", Value: fo})
}
func (a Attr) Value(value interface{}) Attr {
	return append(a, AttrPair{Key: "value", Value: value})
}

type Element struct {
	Type string
	Attr
	Children Blocks
	Options  Option
}

func (Element) RenderHTML() Block { return nil }

type Option int8

const (
	Void Option = 1 << iota
	SelfClose
	CSSElement
	JSElement
	NoWhitespace
)

func Doctype(arg string) Block {
	return NewElement("!DOCTYPE", Attr{{arg, nil}}, nil, Void)
}
func Html(attr Attr, children ...Block) Block {
	return NewElement("html", attr, children, 0)
}
func Head(attr Attr, children ...Block) Block {
	return NewElement("head", attr, children, 0)
}
func Noscript(attr Attr, children ...Block) Block {
	return NewElement("noscript", attr, children, 0)
}
func Iframe(attr Attr, children ...Block) Block {
	return NewElement("iframe", attr, children, 0)
}
func Link(attr Attr, children ...Block) Block {
	return NewElement("link", attr, children, SelfClose)
}
func Img(attr Attr, children ...Block) Block {
	return NewElement("img", attr, children, SelfClose)
}
func Meta(attr Attr, children ...Block) Block {
	return NewElement("meta", attr, children, Void)
}
func Title(attr Attr, children ...Block) Block {
	return NewElement("title", attr, children, 0)
}
func Body(attr Attr, children ...Block) Block {
	return NewElement("body", attr, children, 0)
}
func Button(attr Attr, children ...Block) Block {
	return NewElement("button", attr, children, 0)
}
func Style(attr Attr, children ...Block) Block {
	return NewElement("style", attr, children, CSSElement)
}
func Script(attr Attr, children ...Block) Block {
	return NewElement("script", attr, children, JSElement)
}
func Textarea(attr Attr, children ...Block) Block {
	return NewElement("textarea", attr, children, 0)
}
func Main(attr Attr, children ...Block) Block {
	return NewElement("main", attr, children, 0)
}
func Form(attr Attr, children ...Block) Block {
	return NewElement("form", attr, children, 0)
}
func Nav(attr Attr, children ...Block) Block {
	return NewElement("nav", attr, children, 0)
}
func Span(attr Attr, children ...Block) Block {
	return NewElement("span", attr, children, 0)
}
func I(attr Attr, children ...Block) Block {
	return NewElement("i", attr, children, 0)
}
func Div(attr Attr, children ...Block) Block {
	return NewElement("div", attr, children, 0)
}
func P(attr Attr, children ...Block) Block {
	return NewElement("p", attr, children, 0)
}
func Ul(attr Attr, children ...Block) Block {
	return NewElement("ul", attr, children, 0)
}
func Li(attr Attr, children ...Block) Block {
	return NewElement("li", attr, children, 0)
}
func A(attr Attr, children ...Block) Block {
	return NewElement("a", attr, children, 0)
}
func H1(attr Attr, children ...Block) Block {
	return NewElement("h1", attr, children, 0)
}
func H2(attr Attr, children ...Block) Block {
	return NewElement("h2", attr, children, 0)
}
func H3(attr Attr, children ...Block) Block {
	return NewElement("h3", attr, children, 0)
}
func Pre(attr Attr, children ...Block) Block {
	return NewElement("pre", attr, children, NoWhitespace)
}
func Label(attr Attr, children ...Block) Block {
	return NewElement("label", attr, children, 0)
}
func Input(attr Attr, children ...Block) Block {
	return NewElement("Input", attr, children, SelfClose)
}

func NewElement(el string, attr Attr, children []Block, opt Option) Block {
	if len(children) == 0 {
		return Element{el, attr, nil, opt}
	}
	if len(children) == 1 {
		return Element{el, attr, children, opt}
	}
	return Element{el, attr, Blocks(children), opt}
}
