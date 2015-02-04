package main

import (
	"log"
	"net/http"
)

func dogHandler(w http.ResponseWriter, r *http.Request) {
	/*db, err := lib.Connect()

	if err != nil {
		log.Printf("Error getting data connection: %v\n", err)
		http.Error(w, err.Error(), 500)
		return
	}*/

	data := map[string]interface{}{}
	data["title"] = "TODO"

	err := dogTemplate.Execute(w, data)

	if err != nil {
		log.Printf("Error rendering template: %v\n", err)
		http.Error(w, err.Error(), 500)
		return
	}

}
