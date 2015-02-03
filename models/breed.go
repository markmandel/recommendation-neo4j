package models

import (
	"encoding/json"
	"fmt"
	"github.com/jmcvetta/neoism"
)

//Breed represents a dog breed
type Breed struct {
	ID          int    `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
	WikiURL     string `json:"wikiURL"`
}

func (b Breed) String() string {
	return fmt.Sprintf("[ ID: %v  Breed: %v, Description: %v, Wiki: %v ]", b.ID, b.Name, b.Description, b.WikiURL)
}

func (b *Breed) fromNode(n neoism.Node) error {
	b.ID = n.Id()

	j, err := json.Marshal(n.Data)

	if err != nil {
		return err
	}

	err = json.Unmarshal(j, b)

	if err != nil {
		return err
	}

	return nil
}

// GetRandomBreed returns a random dog breed
func GetRandomBreed(db *neoism.Database) (b *Breed, err error) {
	result := []struct {
		B neoism.Node // Column "b" gets automagically unmarshalled into field B
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

		err = b.fromNode(breedNode)

		if err != nil {
			return b, err
		}
	}

	return
}
