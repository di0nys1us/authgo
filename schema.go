package authgo

import (
	"io/ioutil"
	"log"
)

func readSchema() string {
	data, err := ioutil.ReadFile("../schema.graphql")

	if err != nil {
		log.Fatal(err)
	}

	return string(data)
}
