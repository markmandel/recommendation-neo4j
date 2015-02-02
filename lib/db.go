package lib

import (
	"github.com/jmcvetta/neoism"
	"os"
)

const NEO4J_ENV_KEY = "NEO4J_HOST"

// Connect to the neo4j database.
// Looks for ENV var of NEO4J_HOST for the path to neo4j
// if it's not found, uses http://localhost:7474/db/data as the default
func Connect() (*neoism.Database, error) {
	host := os.Getenv(NEO4J_ENV_KEY)
	if len(host) == 0 {
		host = "http://localhost:7474/db/data"
	}

	return neoism.Connect(host)
}
