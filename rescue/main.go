package main

import (
	"github.com/gorilla/mux"
	"github.com/markmandel/recommendation-neo4j/rescue/templates"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const rescuePortEnvKey = "RESCUE_PORT"
const resourcesDirEnvKey = "RESOURCES_DIR"

var indexTemplate *template.Template
var dogTemplate *template.Template

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/dog/{id}", dogHandler)

	resourcesDir := os.Getenv(resourcesDirEnvKey)
	if len(resourcesDir) == 0 {
		var err error
		resourcesDir, err = filepath.Abs("./resources")

		if err != nil {
			log.Fatalf("Could not access resource directory: %v. %v", resourcesDir, err)
		}
	}

	//static files
	log.Printf("Static file path is: %v", resourcesDir)
	staticHandler := http.FileServer(http.Dir(resourcesDir))
	r.Handle("/resources/{dir}/{file}", http.StripPrefix("/resources/", staticHandler))

	http.Handle("/", r)

	port := os.Getenv(rescuePortEnvKey)

	if len(port) == 0 {
		port = "8080"
	}

	log.Println("Starting server on port", port)

	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func init() {
	//set up all my templates with the standard dependencies.
	indexTemplate = template.Must(template.New("index").Parse(templates.Index))
	dogTemplate = template.Must(template.New("index").Parse(templates.Dog))

	for _, t := range []*template.Template{indexTemplate, dogTemplate} {
		template.Must(t.New("header").Parse(templates.Header))
		template.Must(t.New("footer").Parse(templates.Footer))
		template.Must(t.New("disclaimer").Parse(templates.Disclaimer))
	}
}
