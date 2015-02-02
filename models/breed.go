package models

import (
	"encoding/json"
	"fmt"
	"github.com/jmcvetta/neoism"
)

type Breed struct {
	Id          int    `json:-`
	Name        string `json:"name"`
	Description string `json:"description"`
	WikiURL     string `json:"wikiURL"`
}

func (b Breed) String() string {
	return fmt.Sprintf("[ Id: %v  Breed: %v, Description: %v, Wiki: %v ]", b.Id, b.Name, b.Description, b.WikiURL)
}

func GetRandomBreed(db *neoism.Database) (b *Breed, err error) {
	result := []struct {
		B neoism.Node // Column "n" gets automagically unmarshalled into field N
	}{}

	q := &neoism.CypherQuery{
		Statement: `
		MATCH (b:Breed)
		WITH b, RAND() as r
		ORDER BY r
		LIMIT 1
		RETURN b`,
		Result: &result,
	}

	err = db.Cypher(q)

	if err == nil {
		b = new(Breed)

		breedNode := result[0].B
		breedNode.Db = db

		b.Id = breedNode.Id()

		j, err := json.Marshal(breedNode.Data)

		if err != nil {
			return b, err
		}

		//now put is back into the object we want

		err = json.Unmarshal(j, b)

		if err != nil {
			return b, err
		}
	}

	return
}
