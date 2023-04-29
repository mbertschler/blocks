package main

import (
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/mbertschler/blocks/html"
)

func main() {
	generateElements()
	generateAttributes()
}

func generateElements() {
	elementsTemplate, err := template.New("elements").Parse(elementsFileTemplate)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create("html/elements_gen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data := []elementTemplateData{}
	for _, element := range elements {
		option := elementOptions[element.Option]
		if option == "" {
			log.Fatalf("unknown option %d for element %s", element.Option, element.Name)
		}
		data = append(data, elementTemplateData{
			FuncName:   funcName(element.Name),
			TagName:    element.Name,
			Option:     option,
			NoChildren: element.Option == html.Void,
		})
	}
	err = elementsTemplate.Execute(file, data)
	if err != nil {
		log.Fatal(err)
	}
}

var elementOptions = map[html.ElementOption]string{
	0:                 "0",
	html.Void:         "Void",
	html.CSSElement:   "CSSElement",
	html.JSElement:    "JSElement",
	html.NoWhitespace: "NoWhitespace",
}

type elementTemplateData struct {
	FuncName   string
	TagName    string
	Option     string
	NoChildren bool
}

func generateAttributes() {
	attributesTemplate, err := template.New("attributes").Parse(attributesFileTemplate)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create("html/attr/attributes_gen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data := []attributeTemplateData{}
	for _, attribute := range attributes {
		data = append(data, attributeTemplateData{
			FuncName: funcName(attribute.Name),
			AttrName: attribute.Name,
		})
	}
	err = attributesTemplate.Execute(file, data)
	if err != nil {
		log.Fatal(err)
	}
}

func funcName(name string) string {
	title := strings.Title(name)
	return strings.ReplaceAll(title, "-", "")
}

type attributeTemplateData struct {
	FuncName string
	AttrName string
}
