package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/jmcvetta/neoism"
	"github.com/markmandel/recommendation-neo4j/models"
	"log"
	"net/url"
)

// Scrapes wikipedia for the list of dog names and descriptions and returns them all
func generateBreeds() []models.Breed {
	doc, err := goquery.NewDocument("http://en.wikipedia.org/wiki/List_of_dog_breeds")

	if err != nil {
		log.Fatalln("Could not hit url", err)
	}

	breeds := []models.Breed{}

	//finding initial table
	doc.Find("table").Eq(0).Find("tr").Each(func(i int, s *goquery.Selection) {
		cells := s.Find("td")

		if cells.Length() > 0 {
			breed := models.Breed{}

			first := cells.Eq(0)
			breed.Name = first.Text()

			href, exists := first.Find("a").Attr("href")

			if !exists {
				log.Fatalln("eh, what now?", breed.String())
			}

			ref, err := url.Parse("http://en.wikipedia.org" + href)
			breed.WikiURL = ref.String()

			if err != nil {
				log.Fatalln("Could not parser url: ", href)
			}

			extinct := cells.Eq(2).Text() == "Extinct"

			if !extinct {
				s, err := goquery.NewDocument(breed.WikiURL)

				if err != nil {
					log.Fatalln("Could not get breed details", breed, err)
				}

				breed.Description = s.Find("p").First().Text()

				log.Println("Breed: ", breed.String())
				breeds = append(breeds, breed)
			}
		}
	})

	return breeds
}

// inserts the breeds into neo4j
func storeBreeds(db *neoism.Database, breeds []models.Breed) {
	//make sure the dog breed name index is there
	cq := new(neoism.CypherQuery)
	cq.Statement = `CREATE
					(n:Breed
					{
						name: {name},
						description: {description},
						wikiURL: {wikiIURL}
					})`

	for _, b := range breeds {
		cq.Parameters = neoism.Props{"name": b.Name, "description": b.Description, "wikiIURL": b.WikiURL}
		err = db.Cypher(cq)

		if err != nil {
			log.Fatalln("Could not insert data", err, b, cq)
		}
	}
}
