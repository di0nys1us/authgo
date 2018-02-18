package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/di0nys1us/httpgo"
	"github.com/go-chi/chi"
	"github.com/neelance/graphql-go"
	"github.com/neelance/graphql-go/relay"
)

func newRouter() http.Handler {
	db, err := newDB()

	if err != nil {
		log.Fatal(err)
	}

	security := newSecurity(db)
	router := chi.NewRouter()

	schema, err := graphql.ParseSchema(readSchema(), &rootResolver{db})

	if err != nil {
		log.Fatal(err)
	}

	// Protected routes
	router.Group(func(g chi.Router) {
		g.Use(authorize)

		g.Handle("/graphql", &relay.Handler{Schema: schema})
		g.Method(http.MethodGet, "/graphiql", httpgo.ErrorHandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
			tmpl, err := template.ParseFiles("./templates/graphiql.html")

			if err != nil {
				return err
			}

			return tmpl.Execute(w, nil)
		}))
	})

	// Public routes
	router.Group(func(g chi.Router) {
		g.Method(http.MethodGet, "/login", httpgo.ErrorHandlerFunc(security.getLogin))
		g.Method(http.MethodPost, "/login", httpgo.ErrorHandlerFunc(security.postLogin))
		g.Method(http.MethodPost, "/authenticate", httpgo.ErrorHandlerFunc(security.authenticate))
		g.Method(http.MethodGet, "/invalidate", httpgo.ErrorHandlerFunc(invalidate))
	})

	return router
}
