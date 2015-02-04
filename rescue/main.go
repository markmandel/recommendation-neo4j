package main

import (
	"github.com/gorilla/mux"
	"github.com/markmandel/recommendation-neo4j/rescue/templates"
	"html/template"
	"log"
	"net/http"
	"os"
)

const rescuePortEnvKey = "RESCUE_PORT"

var indexTemplate *template.Template
var dogTemplate *template.Template

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/dog/{id}", dogHandler)
	http.Handle("/", r)

	port := os.Getenv(rescuePortEnvKey)
	if len(port) == 0 {
		port = "8080"
	}

	log.Println("Starting server on port", port)

	http.ListenAndServe(":"+port, nil)
}

func init() {
	indexTemplate = template.Must(template.New("index").Parse(templates.Index))
	dogTemplate = template.Must(template.New("index").Parse(templates.Dog))

	for _, t := range []*template.Template{indexTemplate, dogTemplate} {
		template.Must(t.New("header").Parse(templates.Header))
		template.Must(t.New("footer").Parse(templates.Footer))
	}
}
