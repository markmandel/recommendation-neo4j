package main

import (
	"fmt"
	"github.com/jmcvetta/neoism"
	"github.com/markmandel/recommendation-neo4j/models"
	"io"
	"log"
	"net/http"
	"os"
)

//DownloadDogPictures downloads dog pictures and stores them in ./resources/images
func DownloadDogPictures(db *neoism.Database) {
	dogs, err := models.ListDogs(db)

	if err != nil {
		log.Fatalf("Could not get a list of the dogs: %v", err)
	}

	for _, d := range dogs {
		log.Printf("Image: %v\n", d.PicURL)

		imagePath := fmt.Sprintf("./resources/images/%v-%v.jpg", d.ID, d.Name)
		log.Printf("Downloading to: %v\n", imagePath)
		out, err := os.Create(imagePath)

		if err != nil {
			log.Fatalf("Could not open file %v. %v", imagePath, err)
		}

		defer out.Close()

		resp, err := http.Get(d.PicURL)

		if err != nil {
			log.Fatalf("Could not retrieve image for dog %#v. Error: %v", d, err)
		}

		defer resp.Body.Close()

		_, err = io.Copy(out, resp.Body)

		if err != nil {
			log.Fatalf("Could not copy all of the file. %v", err)
		}
	}
}
