package main

import (
	"forum/server"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	//utils.CreateDatabase()
	server := server.NewServer()
	server.StartServer("8080")
}
