package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/jmcvetta/neoism"
	"github.com/markmandel/recommendation-neo4j/models"
	"log"
	"net/url"
)

func GenerateBreeds() []models.Breed {
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

			var err error
			breed.WikiURL, err = url.Parse("http://en.wikipedia.org" + href)

			if err != nil {
				log.Fatalln("Could not parser url: ", href)
			}

			extinct := cells.Eq(2).Text() == "Extinct"

			if !extinct {
				s, err := goquery.NewDocument(breed.WikiURL.String())

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

func storeBreeds(breeds []models.Breed) {
	db, err := neoism.Connect("http://localhost:7474/db/data")

	if err != nil {
		log.Fatalln("Could not connect to db", err)
	}

	//make sure the dog breed name index is there
	cq := new(neoism.CypherQuery)
	cq.Statement = `CREATE CONSTRAINT ON (b:Breed) ASSERT b.Name IS UNIQUE`
	err = db.Cypher(cq)

	if err != nil {
		log.Fatalln("Could not create unique constraint", cq, err)
	}

	cq.Statement = `CREATE
					(n:Breed
					{
						name: {name},
						description: {description},
						wikiURL: {wikiIURL}
					})`

	for _, b := range breeds {
		cq.Parameters = neoism.Props{"name": b.Name, "description": b.Description, "wikiIURL": b.WikiURL.String()}
		err = db.Cypher(cq)

		if err != nil {
			log.Fatalln("Could not insert data", b, cq, err)
		}
	}
}
