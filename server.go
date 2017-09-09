package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/di0nys1us/httpgo"

	"github.com/di0nys1us/authgo/handlers"
	"github.com/di0nys1us/authgo/repository"
	"github.com/di0nys1us/authgo/security"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

const (
	defaultPort = "3000"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalln(err)
	}

	addr := fmt.Sprintf(":%s", defaultPort)

	if port, ok := os.LookupEnv("AUTHGO_PORT"); ok {
		addr = fmt.Sprintf(":%s", port)
	}

	log.Fatalln(http.ListenAndServe(addr, createRouter()))
}

func createRouter() *chi.Mux {
	r, err := repository.NewRepository()

	if err != nil {
		log.Fatalln(err)
	}

	a := handlers.NewAPI(r)
	s := security.NewSecurity(func(email string) (security.Subject, error) {
		return r.FindUserByEmail(email)
	})

	router := chi.NewRouter()

	// Protected routes
	router.Group(func(g chi.Router) {
		g.Use(security.ValidateJWT)

		g.Method(http.MethodGet, "/users", httpgo.ResponseHandlerFunc(a.GetUsers))
		g.Method(http.MethodGet, "/users/{id:[0-9]+}", httpgo.ResponseHandlerFunc(a.GetUser))
		g.Method(http.MethodPost, "/users", httpgo.ResponseHandlerFunc(a.PostUser))
		g.Method(http.MethodPut, "/users/{id:[0-9]+}", httpgo.ResponseHandlerFunc(a.PutUser))
	})

	// Public routes
	router.Group(func(g chi.Router) {
		g.Method(http.MethodPost, "/authenticate", httpgo.ResponseHandlerFunc(s.Authenticate))
	})

	return router
}
