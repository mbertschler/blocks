package html

func Attr(key string, value interface{}) Attributes {
	return Attributes{AttrPair{Key: key, Value: value}}
}
func (a Attributes) Attr(key string, id interface{}) Attributes {
	return append(a, AttrPair{Key: "id", Value: id})
}

func Id(id interface{}) Attributes {
	return Attributes{AttrPair{Key: "id", Value: id}}
}
func (a Attributes) Id(id interface{}) Attributes {
	return append(a, AttrPair{Key: "id", Value: id})
}

func Class(class interface{}) Attributes {
	return Attributes{AttrPair{Key: "class", Value: class}}
}
func (a Attributes) Class(class interface{}) Attributes {
	return append(a, AttrPair{Key: "class", Value: class})
}

func Href(href interface{}) Attributes {
	return Attributes{AttrPair{Key: "href", Value: href}}
}
func (a Attributes) Href(href interface{}) Attributes {
	return append(a, AttrPair{Key: "href", Value: href})
}

func Rel(rel interface{}) Attributes {
	return Attributes{AttrPair{Key: "rel", Value: rel}}
}
func (a Attributes) Rel(rel interface{}) Attributes {
	return append(a, AttrPair{Key: "rel", Value: rel})
}

func Name(name interface{}) Attributes {
	return Attributes{AttrPair{Key: "name", Value: name}}
}
func (a Attributes) Name(name interface{}) Attributes {
	return append(a, AttrPair{Key: "name", Value: name})
}

func (a Attributes) Content(name interface{}) Attributes {
	return append(a, AttrPair{Key: "content", Value: name})
}
func Content(name interface{}) Attributes {
	return Attributes{AttrPair{Key: "content", Value: name}}
}

func Checked() Attributes {
	return Attributes{AttrPair{Key: "checked", Value: nil}}
}
func (a Attributes) Checked() Attributes {
	return append(a, AttrPair{Key: "checked", Value: nil})
}

func Defer() Attributes {
	return Attributes{AttrPair{Key: "defer", Value: nil}}
}
func (a Attributes) Defer() Attributes {
	return append(a, AttrPair{Key: "defer", Value: nil})
}

func Src(src interface{}) Attributes {
	return Attributes{AttrPair{Key: "src", Value: src}}
}
func (a Attributes) Src(src interface{}) Attributes {
	return append(a, AttrPair{Key: "src", Value: src})
}

func Action(action interface{}) Attributes {
	return Attributes{AttrPair{Key: "action", Value: action}}
}
func (a Attributes) Action(action interface{}) Attributes {
	return append(a, AttrPair{Key: "action", Value: action})
}

func Method(method interface{}) Attributes {
	return Attributes{AttrPair{Key: "method", Value: method}}
}
func (a Attributes) Method(method interface{}) Attributes {
	return append(a, AttrPair{Key: "method", Value: method})
}

func Type(typ interface{}) Attributes {
	return Attributes{AttrPair{Key: "type", Value: typ}}
}
func (a Attributes) Type(typ interface{}) Attributes {
	return append(a, AttrPair{Key: "type", Value: typ})
}

func For(fo interface{}) Attributes {
	return Attributes{AttrPair{Key: "for", Value: fo}}
}
func (a Attributes) For(fo interface{}) Attributes {
	return append(a, AttrPair{Key: "for", Value: fo})
}

func Value(value interface{}) Attributes {
	return Attributes{AttrPair{Key: "value", Value: value}}
}
func (a Attributes) Value(value interface{}) Attributes {
	return append(a, AttrPair{Key: "value", Value: value})
}

func Data(key string, value interface{}) Attributes {
	return Attributes{AttrPair{Key: "data-" + key, Value: value}}
}
func (a Attributes) Data(key string, value interface{}) Attributes {
	return append(a, AttrPair{Key: "data-" + key, Value: value})
}

func Charset(charset interface{}) Attributes {
	return Attributes{AttrPair{Key: "charset", Value: charset}}
}
func (a Attributes) Charset(charset interface{}) Attributes {
	return append(a, AttrPair{Key: "charset", Value: charset})
}

func Styles(style string) Attributes {
	return Attributes{AttrPair{Key: "style", Value: style}}
}
func (a Attributes) Styles(style string) Attributes {
	return append(a, AttrPair{Key: "style", Value: style})
}
