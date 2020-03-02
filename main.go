package main

import (
	"knowledgebase/api"
	"os"
)

var (
	server   api.Server
	user     = getEnvOrDefault("DATABASE_USER", "user")
	password = getEnvOrDefault("DATABASE_PASSWORD", "12345678")
	host     = getEnvOrDefault("DATABASE_HOST", "127.0.0.1")
	db_port  = getEnvOrDefault("DATABASE_PORT", "3306")
	db_name  = getEnvOrDefault("DATABASE_NAME", "KB")
	port     = getEnvOrDefault("PORT", "8080")
)

func main() {
	server.Initialize(user, password, host, db_port, db_name)
	server.Start(":" + port)
}

func getEnvOrDefault(key string, value string) string {
	output := os.Getenv(key)
	if output == "" {
		return value
	}

	return output
}
