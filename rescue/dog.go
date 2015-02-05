package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/markmandel/recommendation-neo4j/models"
	"log"
	"net/http"
	"strconv"
)

func dogHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		log.Printf("Error converting var id to int: %v\n", err)
		http.Error(w, err.Error(), 500)
		return
	}

	dog, err := models.GetDog(db, id)

	if err != nil {
		log.Printf("Error retrieving dog: %v\n", err)
		http.Error(w, err.Error(), 500)
		return
	}

	data := map[string]interface{}{}
	data["dog"] = dog
	data["title"] = fmt.Sprintf("%v - %v", dog.Name, dog.Breed.Name)

	err = dogTemplate.Execute(w, data)

	if err != nil {
		log.Printf("Error rendering template: %v\n", err)
		http.Error(w, err.Error(), 500)
		return
	}

}
