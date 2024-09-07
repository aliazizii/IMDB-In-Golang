package main

import (
	"github.com/aliazizii/IMDB-In-Golang/cmd"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cmd.Execute()
}
