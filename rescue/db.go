package main

import (
	"github.com/jmcvetta/neoism"
	"log"
)

//setup indexes and constraints
func indexes() {
	log.Printf("Applying indexes and contraints")
	indexes := []string{
		"CREATE CONSTRAINT ON (b:Breed) ASSERT b.name IS UNIQUE",
		"CREATE CONSTRAINT ON (d:Dog) ASSERT d.name IS UNIQUE",
		"CREATE CONSTRAINT ON (s:MuxSession) ASSERT s.name IS UNIQUE",
	}

	cq := new(neoism.CypherQuery)

	for _, index := range indexes {
		cq.Statement = index
		err := db.Cypher(cq)

		if err != nil {
			log.Printf("Warning: Could not apply unique constraint/index: %v. %#v", err, cq)
		}
	}

}
