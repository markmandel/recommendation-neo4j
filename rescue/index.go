package main

import (
	"github.com/markmandel/recommendation-neo4j/models"
	"log"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	dogs, err := models.ListDogs(db)

	if err != nil {
		log.Printf("Error listing dogs: %v\n", err)
		http.Error(w, err.Error(), 500)
		return
	}

	data := map[string]interface{}{}
	data["dogs"] = dogs
	data["title"] = "Home"

	err = indexTemplate.Execute(w, data)

	if err != nil {
		log.Printf("Error rendering template: %v\n", err)
		http.Error(w, err.Error(), 500)
		return
	}
}
