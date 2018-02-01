package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/di0nys1us/authgo"
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

	log.Fatal(http.ListenAndServe(addr, authgo.NewRouter()))
}
