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
	session, err := sessionStore.Get(r, siteSessionName)

	if err != nil {
		log.Printf("Error getting session: %v\n", err)
		http.Error(w, err.Error(), 500)
		return
	}

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

	err = models.InsertPageView(db, session, dog)

	if err != nil {
		log.Printf("Error inserting page view for dog. %v", err)
	}

	data := map[string]interface{}{}
	data["dog"] = dog
	data["anchor"] = "#" + dog.Name
	data["title"] = fmt.Sprintf("%v - %v", dog.Name, dog.Breed.Name)

	err = session.Save(r, w)

	if err != nil {
		log.Printf("Error saving session: %v\n", err)
		http.Error(w, err.Error(), 500)
		return
	}

	//view
	err = dogTemplate.Execute(w, data)

	if err != nil {
		log.Printf("Error rendering template: %v\n", err)
		http.Error(w, err.Error(), 500)
		return
	}
}
