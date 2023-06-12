package html

import "github.com/mbertschler/blocks/html/attr"

func Doctype(arg string) Block {
	return newElement("!DOCTYPE", attr.Attributes{{Key: arg}}, nil, Void)
}
