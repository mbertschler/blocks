package html

func Doctype(arg string) Block {
	return newElement("!DOCTYPE", Attributes{{arg, nil}}, nil, Void)
}
func Html(attr Attributes, children ...Block) Block {
	return newElement("html", attr, children, 0)
}
func Head(attr Attributes, children ...Block) Block {
	return newElement("head", attr, children, 0)
}
func Noscript(attr Attributes, children ...Block) Block {
	return newElement("noscript", attr, children, 0)
}
func Iframe(attr Attributes, children ...Block) Block {
	return newElement("iframe", attr, children, 0)
}
func Link(attr Attributes, children ...Block) Block {
	return newElement("link", attr, children, Void)
}
func Img(attr Attributes, children ...Block) Block {
	return newElement("img", attr, children, Void)
}
func Meta(attr Attributes, children ...Block) Block {
	return newElement("meta", attr, children, Void)
}
func Title(attr Attributes, children ...Block) Block {
	return newElement("title", attr, children, 0)
}
func Body(attr Attributes, children ...Block) Block {
	return newElement("body", attr, children, 0)
}
func Button(attr Attributes, children ...Block) Block {
	return newElement("button", attr, children, 0)
}
func Style(attr Attributes, children ...Block) Block {
	return newElement("style", attr, children, CSSElement)
}
func Script(attr Attributes, children ...Block) Block {
	return newElement("script", attr, children, JSElement)
}
func Textarea(attr Attributes, children ...Block) Block {
	return newElement("textarea", attr, children, NoWhitespace)
}
func Main(attr Attributes, children ...Block) Block {
	return newElement("main", attr, children, 0)
}
func Form(attr Attributes, children ...Block) Block {
	return newElement("form", attr, children, 0)
}
func Nav(attr Attributes, children ...Block) Block {
	return newElement("nav", attr, children, 0)
}
func Span(attr Attributes, children ...Block) Block {
	return newElement("span", attr, children, 0)
}
func I(attr Attributes, children ...Block) Block {
	return newElement("i", attr, children, 0)
}
func Div(attr Attributes, children ...Block) Block {
	return newElement("div", attr, children, 0)
}
func P(attr Attributes, children ...Block) Block {
	return newElement("p", attr, children, 0)
}
func Ul(attr Attributes, children ...Block) Block {
	return newElement("ul", attr, children, 0)
}
func Li(attr Attributes, children ...Block) Block {
	return newElement("li", attr, children, 0)
}
func A(attr Attributes, children ...Block) Block {
	return newElement("a", attr, children, 0)
}
func H1(attr Attributes, children ...Block) Block {
	return newElement("h1", attr, children, 0)
}
func H2(attr Attributes, children ...Block) Block {
	return newElement("h2", attr, children, 0)
}
func H3(attr Attributes, children ...Block) Block {
	return newElement("h3", attr, children, 0)
}
func H4(attr Attributes, children ...Block) Block {
	return newElement("h4", attr, children, 0)
}
func H5(attr Attributes, children ...Block) Block {
	return newElement("h5", attr, children, 0)
}
func H6(attr Attributes, children ...Block) Block {
	return newElement("h6", attr, children, 0)
}
func Pre(attr Attributes, children ...Block) Block {
	return newElement("pre", attr, children, NoWhitespace)
}
func Label(attr Attributes, children ...Block) Block {
	return newElement("label", attr, children, 0)
}
func Strong(attr Attributes, children ...Block) Block {
	return newElement("strong", attr, children, 0)
}
func Input(attr Attributes, children ...Block) Block {
	return newElement("input", attr, children, Void)
}
func Br() Block {
	return newElement("br", nil, nil, Void)
}
func Hr(attr Attributes) Block {
	return newElement("hr", attr, nil, Void)
}
