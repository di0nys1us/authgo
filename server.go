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

	roleType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Role",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	})

	authorityType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Authority",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"roles": &graphql.Field{
				Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(roleType))),
			},
		},
	})

	roleType.AddFieldConfig("authorities", &graphql.Field{
		Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(authorityType))),
	})

	userEventType := graphql.NewObject(graphql.ObjectConfig{
		Name: "UserEvent",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"userId": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"createdBy": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"createdAt": &graphql.Field{
				Type: graphql.NewNonNull(graphql.DateTime),
			},
			"type": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"description": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	})

	userType := graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"firstName": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"lastName": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"email": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"password": &graphql.Field{
				Type: graphql.String,
			},
			"enabled": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Boolean),
			},
			"roles": &graphql.Field{
				Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(roleType))),
			},
			"userEvents": &graphql.Field{
				Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(userEventType))),
			},
		},
	})

	fields := graphql.Fields{
		"users": &graphql.Field{
			Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(userType))),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return r.FindUsers()
			},
		},
		"user": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.ID,
				},
				"email": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if id, ok := p.Args["id"].(string); ok {
					return r.FindUser(id)
				}

				if email, ok := p.Args["email"].(string); ok {
					return r.FindUserByEmail(email)
				}

				return nil, nil
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "Query", Fields: fields}
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
		g.Method(http.MethodGet, "/authenticate", httpgo.ErrorHandlerFunc(s.GetAuthenticationForm))
		g.Method(http.MethodPost, "/authenticate", httpgo.ErrorHandlerFunc(s.Authenticate))
		g.Method(http.MethodGet, "/invalidate", httpgo.ErrorHandlerFunc(security.Invalidate))
	})

	return router
}
