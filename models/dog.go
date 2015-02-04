package models

import (
	"encoding/json"
	"github.com/jmcvetta/neoism"
)

//Dog is a dog that is up for adoption
type Dog struct {
	ID      int    `json:"-"`
	Name    string `json:"name"`
	Breed   *Breed
	Adopted bool   `json:"adopted"`
	PicURL  string `json:"picURL"`
}

func (d *Dog) fromNode(n neoism.Node) error {
	d.ID = n.Id()

	j, err := json.Marshal(n.Data)

	if err != nil {
		return err
	}

	err = json.Unmarshal(j, d)

	if err != nil {
		return err
	}

	return nil
}

//UpdateDog updates the dog record. It's assumed the dog will know its ID
func UpdateDog(db *neoism.Database, d *Dog) error {
	cq := &neoism.CypherQuery{
		Statement: `
		MATCH (d:Dog)
		WHERE ID(d) = {id}
		SET d.picURL = {picURL},
			d.adopted = {adopted},
			d.name = {name}`,
		Parameters: neoism.Props{
			"id":      d.ID,
			"picURL":  d.PicURL,
			"adopted": d.Adopted,
			"name":    d.Name},
	}

	return db.Cypher(cq)
}

//ListDogs lists all dogs in neo4j currently
func ListDogs(db *neoism.Database) (results []*Dog, err error) {
	result := []struct {
		D neoism.Node
		B neoism.Node
	}{}

	cq := &neoism.CypherQuery{
		Statement: `
		MATCH (d:Dog)-[:HAS_BREED]->(b:Breed)
		RETURN d, b
		ORDER BY d.name
		`,
		Result: &result,
	}

	err = db.Cypher(cq)

	if err == nil {
		results = []*Dog{}
		for _, node := range result {
			node.D.Db = db
			node.B.Db = db

			dog := new(Dog)
			breed := new(Breed)

			err = dog.fromNode(node.D)
			err = breed.fromNode(node.B)

			dog.Breed = breed

			results = append(results, dog)
		}
	}

	return
}
