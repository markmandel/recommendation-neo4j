package models

import (
	"github.com/gorilla/sessions"
	"github.com/jmcvetta/neoism"
)

//IncrementPageViews increments the counter for how many times this
//session has viewed this dog.
//It will also create the relationship between the session and the dog
//if it doesn't exist.
func IncrementPageViews(db *neoism.Database, s *sessions.Session, d *Dog) error {

	return nil
}
