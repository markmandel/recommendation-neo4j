package models

import (
	"github.com/gorilla/sessions"
	"github.com/jmcvetta/neoism"
	"log"
)

/*
Calculating the initial deviation

MATCH (leftD:Dog),(rightD:Dog)
WHERE ID(leftD) = 547 AND ID(rightD) = 492
MERGE (leftD)-[:L_DEVIATION]->(deviation:SlopeOneDeviation)<-[:R_DEVIATION]-(rightD)
RETURN deviation

MATCH (leftS:MuxSession)-[:HAS_VIEWED]->(:PageView)-[:WITH_DOG]->(leftD:Dog)
WHERE ID(leftD) = 547

WITH leftS
MATCH (rightS:MuxSession)-[:HAS_VIEWED]->(:PageView)-[:WITH_DOG]->(rightD:Dog)
WHERE
	ID(rightD) = 492
	AND
	leftS = rightS

WITH COUNT(DISTINCT leftS) as totalSessions //count how many sessions we have
MATCH (leftD:Dog)-[:L_DEVIATION]->(d:SlopeOneDeviation)<-[:R_DEVIATION]-(rightD:Dog)
WHERE ID(leftD) = 547 AND ID(rightD) = 492
SET d.totalSessions = totalSessions //set the value of how many sessions we have. Need this for actual recommendations.

WITH totalSessions
MATCH (leftS:MuxSession)-[:HAS_VIEWED]->(leftP:PageView)-[:WITH_DOG]->(leftD:Dog)
WHERE ID(leftD) = 547

WITH leftD, leftS, COUNT(DISTINCT leftP) as leftTotal, totalSessions
MATCH (rightS:MuxSession)-[:HAS_VIEWED]->(rightP:PageView)-[:WITH_DOG]->(rightD:Dog)
WHERE ID(rightD) = 492
AND leftS = rightS

WITH leftS, (TOFLOAT((leftTotal - COUNT(DISTINCT rightP)))/totalSessions) as stepDeviation //get the average deviation for each session

WITH SUM(stepDeviation) as deviation //add them all up
MATCH (leftD:Dog)-[:L_DEVIATION]->(d:SlopeOneDeviation)<-[:R_DEVIATION]-(rightD:Dog)
WHERE ID(leftD) = 547 AND ID(rightD) = 492
SET d.deviation = deviation //finally also set the actual deviation value.
RETURN d
*/

//CalculateWeightedSlopeOneDeviation calculates and stores all the deviations for all
//the combinations of Dogs.
func CalculateWeightedSlopeOneDeviation(db *neoism.Database) error {
	dogs, err := ListDogs(db)

	if err != nil {
		return err
	}

	total := len(dogs)

	//for every combination of i and j
	for counter, i := range dogs {
		log.Printf("Calculating deviation for block %v/%v", counter+1, total)

		qs := []*neoism.CypherQuery{}

		for _, j := range dogs {
			part := calculateSingleDeviation(db, i, j)
			qs = append(qs, part...)
		}

		//batch up each 100
		tx, err := db.Begin(qs)

		if err != nil {
			log.Printf("Error attempting to store deviation values: %v. %#v.", err, qs)

			err := tx.Rollback()

			if err != nil {
				return err
			}

			return err
		}

		err = tx.Commit()

		if err != nil {
			return err
		}

		//reset array
		qs = nil
	}

	return nil
}

//calculateSingleDeviation returns the Cypher query to calculates the deviation for a
//left and right dog
func calculateSingleDeviation(db *neoism.Database, l *Dog, r *Dog) []*neoism.CypherQuery {
	props := neoism.Props{"left": l.ID, "right": r.ID}
	return []*neoism.CypherQuery{
		//First create the relationship, if it's not there already
		&neoism.CypherQuery{
			Statement: `
			MATCH (leftD:Dog),(rightD:Dog)
			WHERE ID(leftD) = {left} AND ID(rightD) = {right}
			MERGE (leftD)-[:L_DEVIATION]->(deviation:SlopeOneDeviation)<-[:R_DEVIATION]-(rightD)
			RETURN deviation
			`,
			Parameters: props,
		},
		&neoism.CypherQuery{
			Statement: `
			MATCH (leftS:MuxSession)-[:HAS_VIEWED]->(:PageView)-[:WITH_DOG]->(leftD:Dog)
			WHERE ID(leftD) = {left}
			WITH leftS

			MATCH (rightS:MuxSession)-[:HAS_VIEWED]->(:PageView)-[:WITH_DOG]->(rightD:Dog)
			WHERE
				ID(rightD) = {right}
				AND
				leftS = rightS
			WITH COUNT(DISTINCT leftS) as totalSessions //count how many sessions we have

			MATCH (leftD:Dog)-[:L_DEVIATION]->(d:SlopeOneDeviation)<-[:R_DEVIATION]-(rightD:Dog)
			WHERE ID(leftD) = {left} AND ID(rightD) = {right}
			SET d.totalSessions = totalSessions //set the value of how many sessions we have. Need this for actual recommendations.
			WITH totalSessions

			MATCH (leftS:MuxSession)-[:HAS_VIEWED]->(leftP:PageView)-[:WITH_DOG]->(leftD:Dog)
			WHERE ID(leftD) = {left}
			WITH leftD, leftS, COUNT(DISTINCT leftP) as leftTotal, totalSessions

			MATCH (rightS:MuxSession)-[:HAS_VIEWED]->(rightP:PageView)-[:WITH_DOG]->(rightD:Dog)
			WHERE ID(rightD) = {right}
			AND leftS = rightS
			WITH leftS, (TOFLOAT((leftTotal - COUNT(DISTINCT rightP)))/totalSessions) as stepDeviation //get the average deviation for each session

			WITH SUM(stepDeviation) as deviation //add them all up

			MATCH (leftD:Dog)-[:L_DEVIATION]->(d:SlopeOneDeviation)<-[:R_DEVIATION]-(rightD:Dog)
			WHERE ID(leftD) = {left} AND ID(rightD) = {right}
			SET d.deviation = deviation //finally also set the actual deviation value.
			RETURN d
			`,
			Parameters: props,
		},
	}
}

/*
Determining the actual recommendation.

//all dogs that have been 'rated'(viewed) for this session, with their view count
MATCH(:MuxSession {ident: 'GP5OJLWCLPI7CVSVJQGFN7TSUIW7Z6Q5IKZ7DSO6Z2F2IJK5RJAQ'})-[:HAS_VIEWED]->(view:PageView)-[:WITH_DOG]->(viewedDog:Dog)
WITH viewedDog, COUNT(DISTINCT view) as pageViews

//all dogs this session that have not been viewed, (and aren't adopted)
MATCH (recommendation:Dog { adopted: false })-[:HAS_BREED]->(breed:Breed)
WHERE NOT (:MuxSession {ident: 'GP5OJLWCLPI7CVSVJQGFN7TSUIW7Z6Q5IKZ7DSO6Z2F2IJK5RJAQ'})-[:HAS_VIEWED]->(:PageView)-[:WITH_DOG]->(recommendation)
WITH DISTINCT recommendation, breed, viewedDog, pageViews

//for each dog that has been viewed, add the number of views to the average deviation from recommendation->viewedDog
MATCH (recommendation)-[:L_DEVIATION]->(deviation:SlopeOneDeviation)<-[:R_DEVIATION]-(viewedDog)
WITH ((deviation.deviation + pageViews) * deviation.totalSessions) as score, deviation.totalSessions as totalSessions, recommendation, breed

//SUM all the new scores per recommendation for the numerators, and the SUM of the totalSessions for the denominator
WITH SUM(score) as numerator, SUM(totalSessions) as denominator, recommendation, breed
WHERE denominator > 0

//Wrap it up in a bow, and hand it off
RETURN (numerator/denominator) as expectedViews, recommendation, breed
ORDER BY expectedViews DESC
LIMIT 6

*/

//WeightedSlopeOneRecommendation returns a recommendation based on the
//Weighted Slope One alogorithm.
func WeightedSlopeOneRecommendation(db *neoism.Database, s *sessions.Session) (results []*Dog, err error) {
	result := []struct {
		Recommendation neoism.Node
		Breed          neoism.Node
	}{}

	cq := &neoism.CypherQuery{
		Statement: `
		//all dogs that have been 'rated'(viewed) for this session, with their view count
		MATCH(:MuxSession {ident: {ident}})-[:HAS_VIEWED]->(view:PageView)-[:WITH_DOG]->(viewedDog:Dog)
		WITH viewedDog, COUNT(DISTINCT view) as pageViews

		//all dogs this session that have not been viewed, (and aren't adopted)
		MATCH (recommendation:Dog { adopted: false })-[:HAS_BREED]->(breed:Breed)
		WHERE NOT (:MuxSession {ident: {ident}})-[:HAS_VIEWED]->(:PageView)-[:WITH_DOG]->(recommendation)
		WITH DISTINCT recommendation, breed, viewedDog, pageViews

		//for each dog that has been viewed, add the number of views to the average deviation from recommendation->viewedDog
		MATCH (recommendation)-[:L_DEVIATION]->(deviation:SlopeOneDeviation)<-[:R_DEVIATION]-(viewedDog)
		WITH ((deviation.deviation + pageViews) * deviation.totalSessions) as score, deviation.totalSessions as totalSessions, recommendation, breed

		//SUM all the new scores per recommendation for the numerators, and the SUM of the totalSessions for the denominator
		WITH SUM(score) as numerator, SUM(totalSessions) as denominator, recommendation, breed
		WHERE denominator > 0

		//Wrap it up in a bow, and hand it off
		RETURN (numerator/denominator) as expectedViews, recommendation, breed
		ORDER BY expectedViews DESC
		LIMIT 6
		`,
		Parameters: neoism.Props{"ident": s.ID},
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
