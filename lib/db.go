package lib

import (
	"github.com/jmcvetta/neoism"
	"log"
	"os"
)

const NEO4J_ENV_KEY = "NEO4J_HOST"

func Connect() *neoism.Database {
	host := os.Getenv(NEO4J_ENV_KEY)
	if len(host) == 0 {
		host = "http://localhost:7474/db/data"
	}

	db, err := neoism.Connect(host)

	if err != nil {
		log.Fatalln("Could not connect to db", err)
	}

	return db
}
