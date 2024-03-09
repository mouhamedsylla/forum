package main

import (
	"forum/server"
	"forum/utils"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	utils.CreateDatabase()
	server := server.NewServer()
	server.StartServer("8080")
}
