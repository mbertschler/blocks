package html

func Id(id interface{}) Attr {
	return Attr{AttrPair{Key: "id", Value: id}}
}
func (a Attr) Id(id interface{}) Attr {
	return append(a, AttrPair{Key: "id", Value: id})
}

func Class(class interface{}) Attr {
	return Attr{AttrPair{Key: "class", Value: class}}
}
func (a Attr) Class(class interface{}) Attr {
	return append(a, AttrPair{Key: "class", Value: class})
}

func Href(href interface{}) Attr {
	return Attr{AttrPair{Key: "href", Value: href}}
}
func (a Attr) Href(href interface{}) Attr {
	return append(a, AttrPair{Key: "href", Value: href})
}

func Rel(rel interface{}) Attr {
	return Attr{AttrPair{Key: "rel", Value: rel}}
}
func (a Attr) Rel(rel interface{}) Attr {
	return append(a, AttrPair{Key: "rel", Value: rel})
}

func Name(name interface{}) Attr {
	return Attr{AttrPair{Key: "name", Value: name}}
}
func (a Attr) Name(name interface{}) Attr {
	return append(a, AttrPair{Key: "name", Value: name})
}

func (a Attr) Content(name interface{}) Attr {
	return append(a, AttrPair{Key: "content", Value: name})
}
func Content(name interface{}) Attr {
	return Attr{AttrPair{Key: "content", Value: name}}
}

func Defer() Attr {
	return Attr{AttrPair{Key: "defer", Value: nil}}
}
func (a Attr) Defer() Attr {
	return append(a, AttrPair{Key: "defer", Value: nil})
}

func Src(src interface{}) Attr {
	return Attr{AttrPair{Key: "src", Value: src}}
}
func (a Attr) Src(src interface{}) Attr {
	return append(a, AttrPair{Key: "src", Value: src})
}

func Action(action interface{}) Attr {
	return Attr{AttrPair{Key: "action", Value: action}}
}
func (a Attr) Action(action interface{}) Attr {
	return append(a, AttrPair{Key: "action", Value: action})
}

func Method(method interface{}) Attr {
	return Attr{AttrPair{Key: "method", Value: method}}
}
func (a Attr) Method(method interface{}) Attr {
	return append(a, AttrPair{Key: "method", Value: method})
}

func Type(typ interface{}) Attr {
	return Attr{AttrPair{Key: "type", Value: typ}}
}
func (a Attr) Type(typ interface{}) Attr {
	return append(a, AttrPair{Key: "type", Value: typ})
}

func For(fo interface{}) Attr {
	return Attr{AttrPair{Key: "for", Value: fo}}
}
func (a Attr) For(fo interface{}) Attr {
	return append(a, AttrPair{Key: "for", Value: fo})
}

func Value(value interface{}) Attr {
	return Attr{AttrPair{Key: "value", Value: value}}
}
func (a Attr) Value(value interface{}) Attr {
	return append(a, AttrPair{Key: "value", Value: value})
}

func Data(key string, value interface{}) Attr {
	return Attr{AttrPair{Key: "data-" + key, Value: value}}
}
func (a Attr) Data(key string, value interface{}) Attr {
	return append(a, AttrPair{Key: "data-" + key, Value: value})
}
