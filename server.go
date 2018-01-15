package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/graphql-go/graphql"

	"github.com/di0nys1us/httpgo"

	"github.com/di0nys1us/authgo/repository"
	"github.com/di0nys1us/authgo/security"
	"github.com/go-chi/chi"
	"github.com/graphql-go/handler"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

const (
	defaultPort = "3000"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	addr := fmt.Sprintf(":%s", defaultPort)

	if port, ok := os.LookupEnv("AUTHGO_PORT"); ok {
		addr = fmt.Sprintf(":%s", port)
	}

	log.Fatal(http.ListenAndServe(addr, createRouter()))
}

func createRouter() *chi.Mux {
	r, err := repository.NewRepository()

	if err != nil {
		log.Fatal(err)
	}

	//defer r.Close()

	s := security.NewSecurity(r.FindUserByEmail)

	router := chi.NewRouter()

	fields := graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "world", nil
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)

	if err != nil {
		log.Fatal(err)
	}

	handler := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	// Protected routes
	router.Group(func(g chi.Router) {
		g.Use(security.Authorize)

		g.Handle("/graphql", handler)
	})

	// Public routes
	router.Group(func(g chi.Router) {
		g.Method(http.MethodPost, "/authenticate", httpgo.ErrorHandlerFunc(s.Authenticate))
		g.Method(http.MethodGet, "/invalidate", httpgo.ErrorHandlerFunc(security.Invalidate))
	})

	return router
}
