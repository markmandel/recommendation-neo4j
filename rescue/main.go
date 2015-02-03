package main

import (
	"github.com/gorilla/mux"
	"github.com/markmandel/recommendation-neo4j/rescue/templates"
	"html/template"
	"log"
	"net/http"
	"os"
)

const RESCUE_PORT_ENV_KEY = "RESCUE_PORT"

func main() {
	r := mux.NewRouter()
	http.Handle("/", r)

	port := os.Getenv("RESCUE_PORT_ENV_KEY")
	if len(port) == 0 {
		port = "8080"
	}

	log.Println("Starting server on port", port)

	http.ListenAndServe(":"+port, nil)
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
