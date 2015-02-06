/*
Binary to calculate all the deviations between each adopted dog, so we can do
personalised recommendations for people.
*/
package main

import (
	"github.com/markmandel/recommendation-neo4j/lib"
	"github.com/markmandel/recommendation-neo4j/models"
	"log"
)

//main runs the deviation function for all dogs
func main() {
	db, err := lib.Connect()

	if err != nil {
		log.Fatalf("Could not connect to db: %v", err)
	}

	log.Printf("Starting to calculate derivatives...")

	err = models.CalculateDerivatives(db)

	if err != nil {
		log.Fatalf("Error calculating derivatives: %v", err)
	}

	log.Printf("...Finished!")
}
