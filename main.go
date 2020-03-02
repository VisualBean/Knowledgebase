package main

import (
	"fmt"
	"knowledgebase/api"
	"os"
)

var server api.Server

func main() {

	server.Initialize(
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"))

	server.Start(fmt.Sprint(":%s", os.Getenv("PORT")))
}
