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

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", indexHandler)
	http.Handle("/", r)

	port := os.Getenv(rescuePortEnvKey)
	if len(port) == 0 {
		port = "8080"
	}

	log.Println("Starting server on port", port)

	http.ListenAndServe(":"+port, nil)
}

func init() {
	indexTemplate = templateMust(template.New("index").Parse(templates.Index))
}

func templateMust(t *template.Template, err error) *template.Template {
	if err != nil {
		log.Fatalln("Could not create template.", t.Name(), err)
	}

	return t
}
