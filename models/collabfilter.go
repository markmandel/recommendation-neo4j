package models

/*
Cypher!

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
RETURN SUM(stepDerivative) as derivative //combine them all up
*/
