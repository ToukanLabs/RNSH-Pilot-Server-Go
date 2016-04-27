package main

import (
	"log"
	"net/http"

	"github.com/FiviumAustralia/RNSH-Pilot-Server-Go/rnshschema"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql-go-handler"
)

func main() {

	// define GraphQL schema using relay library helpers
	schema, err := graphql.NewSchema(rnshschema.RnshSchema)

	if err != nil {
		log.Fatal(err)
	}

	h := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})

	fs := http.FileServer(http.Dir("static"))

	// serve HTTP
	http.Handle("/graphql", h)
	http.Handle("/", fs)
	err = http.ListenAndServe(":3001", nil)
	if err != nil {
		log.Fatal(err)
	}

}
