package main

import (
	"knowledgebase/api"
	"os"
)

var server api.Server

func main() {
	server.Initialize(os.Getenv("user"), os.Getenv("password"), os.Getenv("127.0.0.1"))
	server.Start(":8080")
}
