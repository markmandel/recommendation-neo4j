package main

import (
	"github.com/jmcvetta/neoism"
	"github.com/markmandel/recommendation-neo4j/models"
	"log"
)

//put pictures on any dogs that don't have them
func AddDogPicturesFlickr(db *neoism.Database) {
	dogs, err := models.ListDogs(db)

	if err != nil {
		log.Fatalln("Could not list dogs", err)
	}

	log.Println("dogs: ", dogs)
}
