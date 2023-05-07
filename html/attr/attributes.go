package attr

type Attributes []AttrPair
type AttrPair struct {
	Key   string
	Value any
}

func Attr(key string, value interface{}) Attributes {
	return Attributes{AttrPair{Key: key, Value: value}}
}
func (a Attributes) Attr(key string, id interface{}) Attributes {
	return append(a, AttrPair{Key: key, Value: id})
}

func DataAttr(key string, value interface{}) Attributes {
	return Attributes{AttrPair{Key: "data-" + key, Value: value}}
}
func (a Attributes) DataAttr(key string, value interface{}) Attributes {
	return append(a, AttrPair{Key: "data-" + key, Value: value})
}
