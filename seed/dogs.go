package main

import (
	"github.com/jmcvetta/neoism"
	"github.com/markmandel/recommendation-neo4j/models"
	"log"
)

// Returns a list of dogs, each one with a random breed
func generateDogs(db *neoism.Database) []models.Dog {
	result := []models.Dog{}

	for _, name := range DOG_NAMES {
		log.Println("Generating dog: ", name)
		dog := models.Dog{Name: name, Adopted: false}

		breed, err := models.GetRandomBreed(db)

		if err != nil {
			log.Fatalln("Could not get random breed: ", err)
		}

		dog.Breed = breed

		log.Printf("Dog! %v", dog)

		result = append(result, dog)
	}

	return result
}

// Insert the dogs in neo4j.
func storeDogs(db *neoism.Database, dogs []models.Dog) {
	for _, dog := range dogs {
		cq := &neoism.CypherQuery{
			Statement: `
				MATCH (b:Breed)
				WHERE ID(b) = {id}
				CREATE (d:Dog {
					name: {name},
					adopted: {adopted},
					picURL: {picURL}
				})-[:HAS_BREED]->(b)
			`,
			Parameters: neoism.Props{
				"id":      dog.Breed.Id,
				"name":    dog.Name,
				"adopted": dog.Adopted,
				"picURL":  dog.PicURL,
			},
		}

		err := db.Cypher(cq)

		if err != nil {
			log.Fatalln("Could not insert Dog, ", err, dog, cq)
		}

	}
}
