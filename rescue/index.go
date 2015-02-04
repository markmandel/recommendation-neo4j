package main

import (
	"github.com/markmandel/recommendation-neo4j/lib"
	"github.com/markmandel/recommendation-neo4j/models"
	"log"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	db, err := lib.Connect()

	if err != nil {
		log.Printf("Error getting data connection: %v\n", err)
		http.Error(w, err.Error(), 500)
		return
	}

	dogs, err := models.ListDogs(db)

	if err != nil {
		log.Printf("Error listing dogs: %v\n", err)
		http.Error(w, err.Error(), 500)
		return
	}

	data := map[string]interface{}{}
	data["dogs"] = dogs
	data["title"] = "Home"

	indexTemplate.Execute(w, data)
}
