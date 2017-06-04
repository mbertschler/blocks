package html

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"io"
)

var renderDebug = true

type UnsafeString string
type Text string
type Comment string
type CSS template.CSS
type HTML template.HTML
type HTMLAttr template.HTMLAttr
type JS template.JS
type JSStr template.JSStr
type URL template.URL

func (UnsafeString) RenderHTML() Block { return nil }
func (Text) RenderHTML() Block         { return nil }
func (Comment) RenderHTML() Block      { return nil }
func (CSS) RenderHTML() Block          { return nil }
func (HTML) RenderHTML() Block         { return nil }
func (HTMLAttr) RenderHTML() Block     { return nil }
func (JS) RenderHTML() Block           { return nil }
func (b JS) renderString() string      { return string(b) }
func (JSStr) RenderHTML() Block        { return nil }
func (b JSStr) renderString() string   { return string(b) }
func (URL) RenderHTML() Block          { return nil }

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
	case nil:
		return nil
	case UnsafeString:
		if !ctx.minified {
			w.Write(bytes.Repeat([]byte{' '}, ctx.level*indentation))
		}
		w.Write([]byte(el))
		if !ctx.minified {
			w.Write([]byte{'\n'})
		}
		ctx.next()
	case Text:
		if !ctx.minified {
			w.Write(bytes.Repeat([]byte{' '}, ctx.level*indentation))
		}
		w.Write([]byte(template.HTMLEscapeString(string(el))))
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
			attr += " " + v.Key + "=\"" + html.EscapeString(fmt.Sprint(v.Value)) + "\""
		}
		w.Write([]byte("<" + el.Type + attr))
		if el.Options&SelfClose != 0 {
			w.Write([]byte("/>"))
		} else {
			w.Write([]byte(">"))
		}
		if len(el.Children) > 0 {
			if !ctx.minified && el.Options&NoWhitespace == 0 {
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
			if !ctx.minified && el.Options&NoWhitespace == 0 {
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

func NewElement(el string, attr Attr, children []Block, opt Option) Block {
	if len(children) == 0 {
		return Element{el, attr, nil, opt}
	}
	if len(children) == 1 {
		return Element{el, attr, children, opt}
	}
	return Element{el, attr, Blocks(children), opt}
}
