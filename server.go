package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/di0nys1us/httpgo"
	"github.com/pkg/errors"

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

	h := handlers.NewHandler(r)
	s := security.NewSecurity(func(email string) (security.Subject, error) {
		user, err := r.FindUserByEmail(email)

		if err != nil {
			return nil, errors.Wrap(err, "authgo: error when finding user")
		}

		if user == nil {
			return nil, nil
		}

		return user, nil
	})

	router := chi.NewRouter()

	// Protected routes
	router.Group(func(g chi.Router) {
		g.Use(security.Authorize)

		g.Method(http.MethodGet, "/users", httpgo.ErrorHandlerFunc(h.GetUsers))
		g.Method(http.MethodGet, "/users/{id:[0-9]+}", httpgo.ErrorHandlerFunc(h.GetUser))
		g.Method(http.MethodPost, "/users", httpgo.ErrorHandlerFunc(h.PostUser))
		g.Method(http.MethodPut, "/users/{id:[0-9]+}", httpgo.ErrorHandlerFunc(h.PutUser))
	})

	// Public routes
	router.Group(func(g chi.Router) {
		g.Method(http.MethodPost, "/authenticate", httpgo.ErrorHandlerFunc(s.Authenticate))
		g.Method(http.MethodGet, "/invalidate", httpgo.ErrorHandlerFunc(security.Invalidate))
	})

	return router
}
