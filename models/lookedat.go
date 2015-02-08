package models

import (
	"github.com/gorilla/sessions"
	"github.com/jmcvetta/neoism"
)

//PeopleWhoLookedAtDogAlsoLookedAt sends back all dogs that were also looked at
//by other sessions, in a count descending order.
func PeopleWhoLookedAtDogAlsoLookedAt(db *neoism.Database, d *Dog, s *sessions.Session) (results []*Dog, err error) {
	result := []struct {
		Recommendation neoism.Node
		Breed          neoism.Node
	}{}

	cq := &neoism.CypherQuery{
		Statement: `
		MATCH (origin:Dog)<-[:WITH_DOG]-(:PageView)<-[:HAS_VIEWED]-(session:MuxSession)-[:HAS_VIEWED]->(view:PageView)-[:WITH_DOG]->(recommendation:Dog)-[:HAS_BREED]->(breed:Breed)
		WHERE
			ID(origin) = {id}
			AND recommendation <> origin
			AND recommendation.adopted = false
			AND session.ident <> {ident}
		RETURN COUNT(DISTINCT view) as total, MAX(view.timestamp) as latest, recommendation, breed
		ORDER BY total DESC, latest DESC
		LIMIT 6
		`,
		Parameters: neoism.Props{"id": d.ID, "ident": s.ID},
		Result:     &result,
	}

	err = db.Cypher(cq)

	if err == nil {
		results = []*Dog{}
		for _, node := range result {
			var dog *Dog
			dog, err = createDogFromResult(db, node.Recommendation, node.Breed)
			results = append(results, dog)
		}
	}

	return
}
