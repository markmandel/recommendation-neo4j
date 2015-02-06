package models

import (
	"github.com/jmcvetta/neoism"
	"log"
)

/*
MATCH (leftD:Dog),(rightD:Dog)
WHERE ID(leftD) = 547 AND ID(rightD) = 492
MERGE (leftD)-[:L_DERIVATIVE]->(derivative:SlopeOneDerivative)<-[:R_DERIVATIVE]-(rightD)
RETURN derivative

MATCH (leftS:MuxSession)-[:HAS_VIEWED]->(:PageView)-[:WITH_DOG]->(leftD:Dog)
WHERE ID(leftD) = 547
WITH leftS
MATCH (rightS:MuxSession)-[:HAS_VIEWED]->(:PageView)-[:WITH_DOG]->(rightD:Dog)
WHERE ID(rightD) = 492
AND leftS = rightS
WITH COUNT(DISTINCT leftS) as elemCount //count how many sessions we have
MATCH (leftS:MuxSession)-[:HAS_VIEWED]->(leftP:PageView)-[:WITH_DOG]->(leftD:Dog)
WHERE ID(leftD) = 547
WITH leftD, leftS, COUNT(DISTINCT leftP) as leftTotal, elemCount
MATCH (rightS:MuxSession)-[:HAS_VIEWED]->(rightP:PageView)-[:WITH_DOG]->(rightD:Dog)
WHERE ID(rightD) = 492
AND leftS = rightS
WITH leftS, (TOFLOAT((leftTotal - COUNT(DISTINCT rightP)))/elemCount) as stepDerivative //get each value for each session
WITH SUM(stepDerivative) as derivative //combine them all up
MATCH (leftD:Dog)-[:L_DERIVATIVE]->(d:SlopeOneDerivative)<-[:R_DERIVATIVE]-(rightD:Dog)
WHERE ID(leftD) = 547 AND ID(rightD) = 492
SET d.derivative = derivative
*/

//CalculateDerivatives calculates and stores all the derivatives for all
//the combinations of Dogs.
func CalculateDerivatives(db *neoism.Database) error {
	dogs, err := ListDogs(db)

	if err != nil {
		return err
	}

	//for every combination of i and j
	for _, i := range dogs {
		for _, j := range dogs {
			err := calculateSingleDerivative(db, i, j)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

//calculateSingleDerivative calculates the derivative for a
//left and right dog
func calculateSingleDerivative(db *neoism.Database, l *Dog, r *Dog) error {
	props := neoism.Props{"left": l.ID, "right": r.ID}
	qs := []*neoism.CypherQuery{
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
			Statement: `MATCH (leftS:MuxSession)-[:HAS_VIEWED]->(:PageView)-[:WITH_DOG]->(leftD:Dog)
			WHERE ID(leftD) = {left}
			WITH leftS
			MATCH (rightS:MuxSession)-[:HAS_VIEWED]->(:PageView)-[:WITH_DOG]->(rightD:Dog)
			WHERE ID(rightD) = {right}
			AND leftS = rightS
			WITH COUNT(DISTINCT leftS) as elemCount //count how many sessions we have
			MATCH (leftS:MuxSession)-[:HAS_VIEWED]->(leftP:PageView)-[:WITH_DOG]->(leftD:Dog)
			WHERE ID(leftD) = {left}
			WITH leftD, leftS, COUNT(DISTINCT leftP) as leftTotal, elemCount
			MATCH (rightS:MuxSession)-[:HAS_VIEWED]->(rightP:PageView)-[:WITH_DOG]->(rightD:Dog)
			WHERE ID(rightD) = {right}
			AND leftS = rightS
			WITH leftS, (TOFLOAT((leftTotal - COUNT(DISTINCT rightP)))/elemCount) as stepDerivative //get each value for each session
			WITH SUM(stepDerivative) as derivative //combine them all up
			MATCH (leftD:Dog)-[:L_DERIVATIVE]->(d:SlopeOneDerivative)<-[:R_DERIVATIVE]-(rightD:Dog)
			WHERE ID(leftD) = {left} AND ID(rightD) = {right}
			SET d.derivative = derivative
			`,
			Parameters: props,
		},
	}

	tx, err := db.Begin(qs)

	if err != nil {
		log.Printf("Error attempting to store derivative values: %v. %#v.", err, qs)

		err := tx.Rollback()

		if err != nil {
			return err
		}

		return err
	}

	return tx.Commit()
}
