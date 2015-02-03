package models

import (
	"encoding/json"
	"fmt"
	"github.com/jmcvetta/neoism"
)

// A dog that is up for adoption
type Dog struct {
	Id      int    `json:-`
	Name    string `json:"name"`
	Breed   *Breed
	Adopted bool   `json:"adopted"`
	PicURL  string `json:"picURL"`
}

func (d Dog) String() string {
	return fmt.Sprintf("[ Name: %v, Adopted: %v, PicURL: %v Breed: %v ]", d.Name, d.Adopted, d.PicURL, d.Breed)
}

//update the dog. It's assumed the dog will know its ID
func UpdateDog(db *neoism.Database, d *Dog) error {
	cq := &neoism.CypherQuery{
		Statement: `
		MATCH (d:Dog)
		SET d.PicURL = {picURL},
			d.Adopted = {adopted},
			d.Name = {name}
		WHERE ID(d) = {id}`,
		Parameters: neoism.Props{"picURL": d.PicURL,
			"adopted": d.Adopted,
			"name":    d.Name},
	}

	return db.Cypher(cq)
}

func ListDogs(db *neoism.Database) (results []*Dog, err error) {
	result := []struct {
		D neoism.Node
	}{}

	cq := &neoism.CypherQuery{
		Statement: `
		MATCH (d:Dog)
		RETURN d
		ORDER BY d.name
		`,
		Result: &result,
	}

	err = db.Cypher(cq)

	if err == nil {
		results = []*Dog{}
		for _, node := range result {
			dog := new(Dog)

			node.D.Db = db
			dog.Id = node.D.Id()

			j, err := json.Marshal(node.D.Data)

			if err != nil {
				return results, err
			}

			err = json.Unmarshal(j, dog)

			if err != nil {
				return results, err
			}

			results = append(results, dog)
		}
	}

	return
}
