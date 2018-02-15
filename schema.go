package main

import (
	"io/ioutil"
	"log"
)

func readSchema() string {
	data, err := ioutil.ReadFile("./schemas/schema.graphql")

	if err != nil {
		log.Fatal(err)
	}

	return string(data)
}
