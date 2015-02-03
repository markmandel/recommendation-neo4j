package main

import (
	"fmt"
	"github.com/markmandel/recommendation-neo4j/rescue/templates"
	"html/template"
	"log"
)

func main() {
	fmt.Printf("Hello world!")
}

var indexTemplate *template.Template

func init() {
	indexTemplate = templateMust(template.New("index").Parse(templates.INDEX))
}

func templateMust(t *template.Template, err error) *template.Template {
	if err != nil {
		log.Fatalln("Could not create template.", t.Name(), err)
	}

	return t
}
