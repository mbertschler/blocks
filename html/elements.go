package html

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
	return NewElement("input", attr, children, SelfClose)
}
func Br() Block {
	return NewElement("br", nil, nil, SelfClose)
}
