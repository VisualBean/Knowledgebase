package main

import (
	"fmt"
	"knowledgebase/api"
	"os"
)

var server api.Server

func main() {
	server.Initialize("user", "12345678", "127.0.0.1")
	server.Start(fmt.Sprint(":%s", os.Getenv("PORT")))
}
