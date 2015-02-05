Building a Recommendation Engine with Neo4J
===========================================

## Installation
To run this code, you will need several things installed:
- [Docker](https://www.docker.com/)
- [Fig](http://www.fig.sh/)

The go binaries are committed to the repository, so they will just run.

## Run it!
```bash
make up
```


## Interesting Cypher Queries

### Sessions to PageViews to Graphs

```cypher
MATCH (s:MuxSession)-[:`HAS_VIEWED`]->(p:PageView)-[WITH_DOG]->(d:Dog) RETURN s,p,d
```