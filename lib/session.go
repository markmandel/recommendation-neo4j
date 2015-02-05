package lib

import (
	"encoding/base32"
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/jmcvetta/neoism"
	"net/http"
	"strings"
)

const identPropertyKey = "ident"

// NewNeo4JStore returns a new NewNeo4JStore.
// Only support storing single level values of map[string]interface{}
//
// See sessions.NewCookieStore() for a description of the other parameters.
func NewNeo4JStore(db *neoism.Database, keyPairs ...[]byte) *Neo4JStore {
	return &Neo4JStore{
		Codecs: securecookie.CodecsFromPairs(keyPairs...),
		Options: &sessions.Options{
			Path:   "/",
			MaxAge: 86400 * 30,
		},
		db: db,
	}
}

//Neo4JStore is storage for Gorilla's session storage mechanism
type Neo4JStore struct {
	Codecs  []securecookie.Codec
	Options *sessions.Options // default configuration
	db      *neoism.Database
}

// Get returns a session for the given name after adding it to the registry.
//
// See CookieStore.Get().
func (n *Neo4JStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(n, name)
}

// New returns a session for the given name without adding it to the registry.
//
// See CookieStore.New().
func (n *Neo4JStore) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(n, name)
	opts := *n.Options
	session.Options = &opts
	session.IsNew = true
	var err error
	if c, errCookie := r.Cookie(name); errCookie == nil {
		err = securecookie.DecodeMulti(name, c.Value, &session.ID, n.Codecs...)
		if err == nil {
			err = n.load(session)
			if err == nil {
				session.IsNew = false
			}
		}
	}
	return session, err
}

// Save adds a single session to the response.
func (n *Neo4JStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	if session.ID == "" {
		session.ID = strings.TrimRight(
			base32.StdEncoding.EncodeToString(
				securecookie.GenerateRandomKey(32)), "=")
	}
	if err := n.save(session); err != nil {
		return err
	}
	encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, n.Codecs...)
	if err != nil {
		return err
	}
	http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	return nil
}

//load loads the property data into the session
func (n *Neo4JStore) load(s *sessions.Session) error {
	result := []struct {
		S neoism.Node
	}{}

	cq := &neoism.CypherQuery{
		Statement: `
		MATCH (s:MuxSession { ident:{ident} })
		RETURN s
		`,
		Parameters: neoism.Props{identPropertyKey: s.ID},
		Result:     &result,
	}

	err := n.db.Cypher(cq)

	if err != nil {
		return err
	}

	if len(result) == 0 {
		return fmt.Errorf("Could not find session for ident: %v", s.ID)
	}

	s.Values = stringMapToInterfaceMap(result[0].S.Data)

	return nil
}

func (n *Neo4JStore) save(s *sessions.Session) error {
	props := interfaceMapToStringMap(s.Values)
	props[identPropertyKey] = s.ID

	_, _, err := n.db.GetOrCreateNode("MuxSession", identPropertyKey, props)

	return err
}

func stringMapToInterfaceMap(src map[string]interface{}) map[interface{}]interface{} {
	dst := map[interface{}]interface{}{}

	for k, v := range src {
		dst[k] = v
	}

	return dst
}

func interfaceMapToStringMap(src map[interface{}]interface{}) map[string]interface{} {
	dst := map[string]interface{}{}

	for k, v := range src {
		dst[fmt.Sprintf("%v", k)] = v
	}

	return dst
}
