Building a Recommendation Engine with Neo4J
===========================================

## Installation
To run this code, you will need several things installed:
- [Docker](https://www.docker.com/)
- [Fig](http://www.fig.sh/)

The go binaries are committed to the repository, so it is not required to compile them to run this project.

## Run it!
```bash
make neo4j-restore
make up
```


## Interesting Cypher Queries

### Sessions to PageViews to Dogs

```cypher
MATCH (s:MuxSession)-[:`HAS_VIEWED`]->(p:PageView)-[WITH_DOG]->(d:Dog) RETURN s,p,d
```