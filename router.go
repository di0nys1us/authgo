package authgo

import (
	"log"
	"net/http"

	"github.com/di0nys1us/httpgo"
	"github.com/go-chi/chi"
	"github.com/neelance/graphql-go"
	"github.com/neelance/graphql-go/relay"
)

// NewRouter TODO
func NewRouter() http.Handler {
	db, err := newDB()

	if err != nil {
		log.Fatal(err)
	}

	sec := newSecurity(db)

	router := chi.NewRouter()

	schema, err := graphql.ParseSchema(readSchema(), &resolver{db})

	if err != nil {
		log.Fatal(err)
	}

	graphiqlHandlerFunc := func(w http.ResponseWriter, r *http.Request) error {
		return writeString(w, templateGraphiQL)
	}

	// Protected routes
	router.Group(func(g chi.Router) {
		g.Use(authorize)

		g.Handle("/graphql", &relay.Handler{Schema: schema})
		g.Method(http.MethodGet, "/graphiql", httpgo.ErrorHandlerFunc(graphiqlHandlerFunc))
	})

	// Public routes
	router.Group(func(g chi.Router) {
		g.Method(http.MethodGet, "/login", httpgo.ErrorHandlerFunc(sec.getLogin))
		g.Method(http.MethodPost, "/login", httpgo.ErrorHandlerFunc(sec.postLogin))
		g.Method(http.MethodPost, "/authenticate", httpgo.ErrorHandlerFunc(sec.authenticate))
		g.Method(http.MethodGet, "/invalidate", httpgo.ErrorHandlerFunc(invalidate))
	})

	return router
}
