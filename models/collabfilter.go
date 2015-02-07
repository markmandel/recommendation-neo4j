package models

import (
	"github.com/jmcvetta/neoism"
	"log"
)

/*
Calculating the initial deviation

MATCH (leftD:Dog),(rightD:Dog)
WHERE ID(leftD) = 547 AND ID(rightD) = 492
MERGE (leftD)-[:L_DERIVATIVE]->(derivative:SlopeOneDerivative)<-[:R_DERIVATIVE]-(rightD)
RETURN derivative

MATCH (leftS:MuxSession)-[:HAS_VIEWED]->(:PageView)-[:WITH_DOG]->(leftD:Dog)
WHERE ID(leftD) = 547

WITH leftS
MATCH (rightS:MuxSession)-[:HAS_VIEWED]->(:PageView)-[:WITH_DOG]->(rightD:Dog)
WHERE
	ID(rightD) = 492
	AND
	leftS = rightS

WITH COUNT(DISTINCT leftS) as totalSessions //count how many sessions we have
MATCH (leftD:Dog)-[:L_DERIVATIVE]->(d:SlopeOneDerivative)<-[:R_DERIVATIVE]-(rightD:Dog)
WHERE ID(leftD) = 547 AND ID(rightD) = 492
SET d.totalSessions = totalSessions //set the value of how many sessions we have. Need this for actual recommendations.

WITH totalSessions
MATCH (leftS:MuxSession)-[:HAS_VIEWED]->(leftP:PageView)-[:WITH_DOG]->(leftD:Dog)
WHERE ID(leftD) = 547

WITH leftD, leftS, COUNT(DISTINCT leftP) as leftTotal, totalSessions
MATCH (rightS:MuxSession)-[:HAS_VIEWED]->(rightP:PageView)-[:WITH_DOG]->(rightD:Dog)
WHERE ID(rightD) = 492
AND leftS = rightS

WITH leftS, (TOFLOAT((leftTotal - COUNT(DISTINCT rightP)))/totalSessions) as stepDerivative //get the average deviation for each session

WITH SUM(stepDerivative) as derivative //add them all up
MATCH (leftD:Dog)-[:L_DERIVATIVE]->(d:SlopeOneDerivative)<-[:R_DERIVATIVE]-(rightD:Dog)
WHERE ID(leftD) = 547 AND ID(rightD) = 492
SET d.derivative = derivative //finally also set the actual deviation value.
RETURN d
*/

//CalculateDerivatives calculates and stores all the derivatives for all
//the combinations of Dogs.
func CalculateDerivatives(db *neoism.Database) error {
	dogs, err := ListDogs(db)

	if err != nil {
		return err
	}

	total := len(dogs)

	//for every combination of i and j
	for counter, i := range dogs {
		log.Printf("Calculating derivatives for block %v/%v", counter+1, total)

		qs := []*neoism.CypherQuery{}

		for _, j := range dogs {
			part := calculateSingleDerivative(db, i, j)
			qs = append(qs, part...)
		}

		//batch up each 100
		tx, err := db.Begin(qs)

		if err != nil {
			log.Printf("Error attempting to store derivative values: %v. %#v.", err, qs)

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

//calculateSingleDerivative returns the Cypher query to calculates the derivative for a
//left and right dog
func calculateSingleDerivative(db *neoism.Database, l *Dog, r *Dog) []*neoism.CypherQuery {
	props := neoism.Props{"left": l.ID, "right": r.ID}
	return []*neoism.CypherQuery{
		//First create the relationship, if it's not there already
		&neoism.CypherQuery{
			Statement: `
			MATCH (leftD:Dog),(rightD:Dog)
			WHERE ID(leftD) = {left} AND ID(rightD) = {right}
			MERGE (leftD)-[:L_DERIVATIVE]->(derivative:SlopeOneDerivative)<-[:R_DERIVATIVE]-(rightD)
			RETURN derivative
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

			MATCH (leftD:Dog)-[:L_DERIVATIVE]->(d:SlopeOneDerivative)<-[:R_DERIVATIVE]-(rightD:Dog)
			WHERE ID(leftD) = {left} AND ID(rightD) = {right}
			SET d.totalSessions = totalSessions //set the value of how many sessions we have. Need this for actual recommendations.
			WITH totalSessions

			MATCH (leftS:MuxSession)-[:HAS_VIEWED]->(leftP:PageView)-[:WITH_DOG]->(leftD:Dog)
			WHERE ID(leftD) = {left}
			WITH leftD, leftS, COUNT(DISTINCT leftP) as leftTotal, totalSessions

			MATCH (rightS:MuxSession)-[:HAS_VIEWED]->(rightP:PageView)-[:WITH_DOG]->(rightD:Dog)
			WHERE ID(rightD) = {right}
			AND leftS = rightS
			WITH leftS, (TOFLOAT((leftTotal - COUNT(DISTINCT rightP)))/totalSessions) as stepDerivative //get the average deviation for each session

			WITH SUM(stepDerivative) as derivative //add them all up

			MATCH (leftD:Dog)-[:L_DERIVATIVE]->(d:SlopeOneDerivative)<-[:R_DERIVATIVE]-(rightD:Dog)
			WHERE ID(leftD) = {left} AND ID(rightD) = {right}
			SET d.derivative = derivative //finally also set the actual deviation value.
			RETURN d
			`,
			Parameters: props,
		},
	}
}

/*
Determining the actual recommendation.

//all dogs that have been 'rated'(viewed) for this session, with their view count

//score for a particular dog.
MATCH(s:MuxSession {ident: 'GP5OJLWCLPI7CVSVJQGFN7TSUIW7Z6Q5IKZ7DSO6Z2F2IJK5RJAQ'})-[:HAS_VIEWED]->(p:PageView)-[:WITH_DOG]->(rightD:Dog)
WITH rightD, COUNT(DISTINCT p) as views

MATCH (leftD)-[:L_DERIVATIVE]->(derivative:SlopeOneDerivative)<-[:R_DERIVATIVE]-(rightD)
WHERE ID(leftD) = 554
WITH ((derivative.derivative + views) * derivative.totalSessions) as score

WITH SUM(score) as numerator

MATCH(s:MuxSession {ident: 'GP5OJLWCLPI7CVSVJQGFN7TSUIW7Z6Q5IKZ7DSO6Z2F2IJK5RJAQ'})-[:HAS_VIEWED]->(p:PageView)-[:WITH_DOG]->(rightD:Dog)
WITH  rightD, COUNT(DISTINCT p) as views, numerator

MATCH (leftD)-[:L_DERIVATIVE]->(derivative:SlopeOneDerivative)<-[:R_DERIVATIVE]-(rightD)
WHERE ID(leftD) = 554
RETURN TOFLOAT(numerator) / SUM(derivative.totalSessions)


//numerators for all dogs not in the current session's viewed set
MATCH(s:MuxSession {ident: 'GP5OJLWCLPI7CVSVJQGFN7TSUIW7Z6Q5IKZ7DSO6Z2F2IJK5RJAQ'})-[:HAS_VIEWED]->(p:PageView)-[:WITH_DOG]->(rightD:Dog)
WITH rightD, COUNT(DISTINCT p) as rightViews

MATCH (leftD:Dog)-[:HAS_BREED]->(leftBreed:Breed)
WHERE NOT (:MuxSession {ident: 'GP5OJLWCLPI7CVSVJQGFN7TSUIW7Z6Q5IKZ7DSO6Z2F2IJK5RJAQ'})-[:HAS_VIEWED]->(:PageView)-[:WITH_DOG]->(leftD)
WITH DISTINCT leftD, leftBreed, rightD, rightViews

MATCH (leftD)-[:L_DERIVATIVE]->(derivative:SlopeOneDerivative)<-[:R_DERIVATIVE]-(rightD)
WITH ((derivative.derivative + rightViews) * derivative.totalSessions) as score, leftD, leftBreed

RETURN SUM(score) as numerator, leftD, leftBreed, ID(leftD) as sort
ORDER BY sort

//denominators for all dogs not in the current session's viewed set
MATCH(s:MuxSession {ident: 'GP5OJLWCLPI7CVSVJQGFN7TSUIW7Z6Q5IKZ7DSO6Z2F2IJK5RJAQ'})-[:HAS_VIEWED]->(p:PageView)-[:WITH_DOG]->(rightD:Dog)
WITH rightD, COUNT(DISTINCT p) as rightViews

MATCH (leftD:Dog)
WHERE NOT (:MuxSession {ident: 'GP5OJLWCLPI7CVSVJQGFN7TSUIW7Z6Q5IKZ7DSO6Z2F2IJK5RJAQ'})-[:HAS_VIEWED]->(:PageView)-[:WITH_DOG]->(leftD)
WITH DISTINCT leftD, rightD, rightViews

MATCH (leftD)-[:L_DERIVATIVE]->(derivative:SlopeOneDerivative)<-[:R_DERIVATIVE]-(rightD)
RETURN SUM(derivative.totalSessions), ID(leftD) as sort
ORDER BY sort

*/
