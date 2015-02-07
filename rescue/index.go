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

	recommendations, err := models.WeightedSlopeOneRecommendation(db, session)

	if err != nil {
		log.Printf("Error getting recommendations: %v\n", err)
		http.Error(w, err.Error(), 500)
		return
	}

	//save the session
	err = session.Save(r, w)

	if err != nil {
		log.Printf("Error saving session: %v\n", err)
		http.Error(w, err.Error(), 500)
		return
	}

	//data for the view
	data := map[string]interface{}{}
	data["dogs"] = dogs
	data["title"] = "Home"
	data["recommendations"] = recommendations

	//view
	err = indexTemplate.Execute(w, data)

	if err != nil {
		log.Printf("Error rendering template: %v\n", err)
		http.Error(w, err.Error(), 500)
		return
	}
}
