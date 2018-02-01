package authgo

import (
	"log"
	"net/http"

	"github.com/graphql-go/graphql"

	"github.com/di0nys1us/httpgo"

	"github.com/go-chi/chi"
	"github.com/graphql-go/handler"
)

func NewRouter() http.Handler {
	db, err := connect()

	if err != nil {
		log.Fatal(err)
	}

	sec := newSecurity(db)

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
		},
	})

	roleType.AddFieldConfig("authorities", &graphql.Field{
		Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(authorityType))),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return nil, nil
		},
	})
	authorityType.AddFieldConfig("roles", &graphql.Field{
		Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(roleType))),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return nil, nil
		},
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
				Type:    graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(roleType))),
				Resolve: nil,
			},
			"events": &graphql.Field{
				Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(userEventType))),
			},
		},
	})

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"users": &graphql.Field{
				Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(userType))),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return nil, nil
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
						return findUserByID(nil, id)
					}

					if email, ok := p.Args["email"].(string); ok {
						return findUserByEmail(nil, email)
					}

					return nil, nil
				},
			},
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})

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
		g.Use(authorize)

		g.Handle("/graphql", handler)
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
