/*
Seed our data set.

We need a list of dog breeds, and pictures of dogs to go with them.
We will scrape Wikipedia (Attribution: http://en.wikipedia.org/wiki/List_of_dog_breeds) to get the
list of dogs.

Then we'll hit up creative commons images on Flickr to look for cute version of the dog photos.
*/
package main

import (
	"flag"
	"github.com/markmandel/recommendation-neo4j/lib"
	"log"
)

var step string

func init() {
	flag.StringVar(&step, "step", "", "Which step do you want in the sequence of 'breeds,dogs'")
	flag.Parse()
}

func main() {
	db := lib.Connect()

	switch step {
	case "breeds":
		log.Println("Seeding Breeds")
		storeBreeds(db, generateBreeds())
	case "dogs":
		log.Println("Seeding Dogs")
		storeDogs(db, generateDogs(db))
	}
}
