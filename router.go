package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/di0nys1us/authgo/sqlgo"

	"github.com/di0nys1us/authgo/security"
	"github.com/di0nys1us/httpgo"
	"github.com/go-chi/chi"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

const (
	regexpUUID = "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"
)

func newRouter() http.Handler {
	db, err := newDB()

	if err != nil {
		log.Fatal(err)
	}

	s := security.New(db)
	router := chi.NewRouter()

	db2, err := sqlgo.NewDB()

	if err != nil {
		log.Fatal(err)
	}

	uh := &userHandler{db2}

	schema, err := graphql.ParseSchema(readSchema(), &rootResolver{
		&rootQuery{db},
		&rootMutation{db},
	})

	if err != nil {
		log.Fatal(err)
	}

	// Protected routes
	router.Group(func(g chi.Router) {
		g.Use(security.Authorize)

		g.Handle("/graphql", &relay.Handler{Schema: schema})
		g.Method(http.MethodGet, "/", httpgo.ErrorHandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
			tmpl, err := template.ParseFiles("./templates/graphiql.html")

			if err != nil {
				return err
			}

			return tmpl.Execute(w, nil)
		}))

		g.Method(http.MethodGet, "/users", httpgo.ErrorHandlerFunc(uh.getUsers))
		g.Method(http.MethodGet, fmt.Sprintf("/users/{userID:%s}", regexpUUID), httpgo.ErrorHandlerFunc(uh.getUser))
	})

	// Public routes
	router.Group(func(g chi.Router) {
		g.Method(http.MethodGet, "/login", httpgo.ErrorHandlerFunc(security.GetLogin))
		g.Method(http.MethodPost, "/login", httpgo.ErrorHandlerFunc(s.PostLogin))
		g.Method(http.MethodPost, "/authenticate", httpgo.ErrorHandlerFunc(s.Authenticate))
		g.Method(http.MethodGet, "/logout", httpgo.ErrorHandlerFunc(security.Logout))
	})

	return router
}
