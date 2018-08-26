package html

func Doctype(arg string) Block {
	return NewElement("!DOCTYPE", Attributes{{arg, nil}}, nil, Void)
}
func Html(attr Attributes, children ...Block) Block {
	return NewElement("html", attr, children, 0)
}
func Head(attr Attributes, children ...Block) Block {
	return NewElement("head", attr, children, 0)
}
func Noscript(attr Attributes, children ...Block) Block {
	return NewElement("noscript", attr, children, 0)
}
func Iframe(attr Attributes, children ...Block) Block {
	return NewElement("iframe", attr, children, 0)
}
func Link(attr Attributes, children ...Block) Block {
	return NewElement("link", attr, children, SelfClose)
}
func Img(attr Attributes, children ...Block) Block {
	return NewElement("img", attr, children, SelfClose)
}
func Meta(attr Attributes, children ...Block) Block {
	return NewElement("meta", attr, children, Void)
}
func Title(attr Attributes, children ...Block) Block {
	return NewElement("title", attr, children, 0)
}
func Body(attr Attributes, children ...Block) Block {
	return NewElement("body", attr, children, 0)
}
func Button(attr Attributes, children ...Block) Block {
	return NewElement("button", attr, children, 0)
}
func Style(attr Attributes, children ...Block) Block {
	return NewElement("style", attr, children, CSSElement)
}
func Script(attr Attributes, children ...Block) Block {
	return NewElement("script", attr, children, JSElement)
}
func Textarea(attr Attributes, children ...Block) Block {
	return NewElement("textarea", attr, children, NoWhitespace)
}
func Main(attr Attributes, children ...Block) Block {
	return NewElement("main", attr, children, 0)
}
func Form(attr Attributes, children ...Block) Block {
	return NewElement("form", attr, children, 0)
}
func Nav(attr Attributes, children ...Block) Block {
	return NewElement("nav", attr, children, 0)
}
func Span(attr Attributes, children ...Block) Block {
	return NewElement("span", attr, children, 0)
}
func I(attr Attributes, children ...Block) Block {
	return NewElement("i", attr, children, 0)
}
func Div(attr Attributes, children ...Block) Block {
	return NewElement("div", attr, children, 0)
}
func P(attr Attributes, children ...Block) Block {
	return NewElement("p", attr, children, 0)
}
func Ul(attr Attributes, children ...Block) Block {
	return NewElement("ul", attr, children, 0)
}
func Li(attr Attributes, children ...Block) Block {
	return NewElement("li", attr, children, 0)
}
func A(attr Attributes, children ...Block) Block {
	return NewElement("a", attr, children, 0)
}
func H1(attr Attributes, children ...Block) Block {
	return NewElement("h1", attr, children, 0)
}
func H2(attr Attributes, children ...Block) Block {
	return NewElement("h2", attr, children, 0)
}
func H3(attr Attributes, children ...Block) Block {
	return NewElement("h3", attr, children, 0)
}
func H4(attr Attributes, children ...Block) Block {
	return NewElement("h4", attr, children, 0)
}
func H5(attr Attributes, children ...Block) Block {
	return NewElement("h5", attr, children, 0)
}
func H6(attr Attributes, children ...Block) Block {
	return NewElement("h6", attr, children, 0)
}
func Pre(attr Attributes, children ...Block) Block {
	return NewElement("pre", attr, children, NoWhitespace)
}
func Label(attr Attributes, children ...Block) Block {
	return NewElement("label", attr, children, 0)
}
func Strong(attr Attributes, children ...Block) Block {
	return NewElement("strong", attr, children, 0)
}
func Input(attr Attributes, children ...Block) Block {
	return NewElement("input", attr, children, SelfClose)
}
func Br() Block {
	return NewElement("br", nil, nil, SelfClose)
}
func Hr(attr Attributes) Block {
	return NewElement("hr", attr, nil, SelfClose)
}
