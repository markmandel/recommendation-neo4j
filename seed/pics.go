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

const FLICKR_ENV_KEY = "FLICKR_KEY"
const FLICKR_SECRET_ENV_KEY = "FLICKR_SECRET"

//put pictures on any dogs that don't have them
func AddDogPicturesFlickr(db *neoism.Database) {
	dogs, err := models.ListDogs(db)

	if err != nil {
		log.Fatalln("Could not list dogs", err)
	}

	fc := flickgo.New(os.Getenv(FLICKR_ENV_KEY), os.Getenv(FLICKR_SECRET_ENV_KEY), http.DefaultClient)

	//try this for dog[0]

	dog := dogs[0]
	log.Printf("Processing: %v", dog)

	q := map[string]string{"text": dog.Breed.Name, "license": "1,2,3,4,5,6,7,8", "sort": "relevence", "media": "photos", "per_page": "10"}
	log.Println("Search:", q)

	photos, err := fc.Search(q)

	if err != nil {
		log.Fatalln("Could not retrieve photos: ", err)
	}

	log.Println("Photo URL: ", photos.Photos[0].URL(flickgo.SizeMedium640))
	log.Println("Creative Commons URL: ", fmt.Sprintf("https://www.flickr.com/photos/%v/%v", photos.Photos[0].Owner, photos.Photos[0].ID))
}
