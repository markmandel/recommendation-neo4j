package models

import (
	"github.com/gorilla/sessions"
	"github.com/jmcvetta/neoism"
	"time"
)

//InsertPageView inserts a record of a page view between a session
//and a dog.
func InsertPageView(db *neoism.Database, s *sessions.Session, d *Dog) error {
	now := time.Now().UTC().Unix()
	cq := &neoism.CypherQuery{
		Statement: `
		MATCH (s:MuxSession { ident: {ident} }), (d:Dog { name: {name} })
		CREATE (s)-[:HAS_VIEWED]->(p:PageView { timestamp: {timestamp} })-[:WITH_DOG]->(d)
		`,
		Parameters: neoism.Props{"ident": s.ID, "name": d.Name, "timestamp": now},
	}

	return db.Cypher(cq)
}
