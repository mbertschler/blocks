package css

import (
	"bytes"
	"fmt"
	"io"
)

var renderDebug = true
var indentation = 4

type Rule struct {
	Sel   string
	Decls []Line
}

func (Rule) RenderCSS() Block { return nil }

type Line struct {
	Key   string
	Value string
}

func New(sel string, lines ...Line) Rule {
	return Rule{Sel: sel, Decls: lines}
}

func L(key, val string) Line {
	return Line{Key: key, Value: val}
}

type Blocks []Block

func (b *Blocks) Add(block Block) {
	*b = append(*b, block)
}

func (b *Blocks) AddBlocks(blocks Blocks) {
	*b = append(*b, blocks...)
}

func (Blocks) RenderCSS() Block { return nil }

type Block interface {
	RenderCSS() Block
}

func Render(root Block, w io.Writer) error {
	err := RenderCSS(root, w, &renderCtx{})
	if err != nil {
		return err
	}
	return nil
}

func RenderMinified(root Block, w io.Writer) error {
	err := RenderCSS(root, w, &renderCtx{minified: true})
	if err != nil {
		return err
	}
	return nil
}

func RenderString(root Block) (string, error) {
	buf := bytes.Buffer{}
	err := RenderCSS(root, &buf, &renderCtx{})
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func RenderMinifiedString(root Block) (string, error) {
	buf := bytes.Buffer{}
	err := RenderCSS(root, &buf, &renderCtx{minified: true})
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

type renderCtx struct {
	item     int
	level    int
	minified bool
}

func (c *renderCtx) enter() (item int) {
	item = c.item
	c.item = 0
	c.level++
	return item
}

func (c *renderCtx) next() {
	c.item++
}

func (c *renderCtx) exit(item int) {
	c.level--
	c.item = item
}

func RenderCSS(c Block, w io.Writer, ctx *renderCtx) error {
	//var item int
	switch el := c.(type) {
	case Rule:
		if ctx.minified {
			w.Write([]byte(el.Sel + "{"))
			for _, l := range el.Decls {
				w.Write([]byte(l.Key + ":" + l.Value + ";"))
			}
			w.Write([]byte{'}'})
		} else {
			w.Write([]byte(el.Sel + " {\n"))
			for _, l := range el.Decls {
				w.Write(bytes.Repeat([]byte{' '}, (ctx.level+1)*indentation))
				w.Write([]byte(l.Key + ": " + l.Value + ";\n"))
			}
			w.Write([]byte("}\n\n"))
		}
		ctx.next()
	case Blocks:
		for _, e := range el {
			RenderCSS(e, w, ctx)
		}
	case Block:
		c := el.RenderCSS()
		RenderCSS(c, w, ctx)
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
