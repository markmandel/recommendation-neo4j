package main

import (
	"fmt"
	"github.com/jmcvetta/neoism"
	"github.com/manki/flickgo"
	"github.com/markmandel/recommendation-neo4j/models"
	"log"
	"net/http"
	"os"
)

const flickrEnvKey = "FLICKR_KEY"
const flickrSecretEnvKey = "FLICKR_SECRET"

//AddDogPicturesFlickr puts pictures on any dogs that don't have them
func AddDogPicturesFlickr(db *neoism.Database) {
	dogs, err := models.ListDogs(db)

	if err != nil {
		log.Fatalln("Could not list dogs", err)
	}

	fc := flickgo.New(os.Getenv(flickrEnvKey), os.Getenv(flickrSecretEnvKey), http.DefaultClient)

	for _, dog := range dogs {
		log.Printf("Processing: %v. Breed: %v", dog.Name, dog.Breed.Name)

		q := map[string]string{"text": dog.Breed.Name, "license": "1,2,3,4,5,6,7,8", "sort": "relevence", "media": "photos", "per_page": "10"}
		log.Println("Search:", q)

		photos, err := fc.Search(q)

		if err != nil {
			log.Fatalln("Could not retrieve photos: ", err)
		}

		if len(photos.Photos) > 0 {
			dog.PicURL = photos.Photos[0].URL(flickgo.SizeMedium640)
			err = models.UpdateDog(db, dog)

			if err != nil {
				log.Fatalln("Could not save dog picture", err, dog)
			}

			log.Println("Creative Commons URL: ", fmt.Sprintf("https://www.flickr.com/photos/%v/%v", photos.Photos[0].Owner, photos.Photos[0].ID))
		}
	}
}
