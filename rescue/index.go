package main

import (
	"github.com/markmandel/recommendation-neo4j/models"
	"log"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, siteSessionName)

	if err != nil {
		log.Printf("Error getting session: %v\n", err)
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

	log.Printf("Session state: %#v", session)
	err = session.Save(r, w)

	//view
	err = indexTemplate.Execute(w, data)

	if err != nil {
		log.Printf("Error rendering template: %v\n", err)
		http.Error(w, err.Error(), 500)
		return
	}

	log.Printf("The current header [v2]: %#v", w.Header())

	if err != nil {
		log.Printf("Error saving session: %v\n", err)
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("X-Testing", "HELLO?")

	log.Printf("Finis!")
}
