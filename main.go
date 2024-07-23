package main

import (
	"fmt"
	"os"
)

func main() {
	// Locally running postgres instance with default values
	host := "localhost"
	port := 5432
	user := "postgres"
	name := "postgres"
	password := "password"

	database := GetDatabase(host, user, name, password, port)

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . [command]: Use 'populate' to fill database, or 'create' to create puzzles.")
		return
	}

	command := os.Args[1]

	switch command {
	case "populate":
		// not provided: crosswords directory and common_words.txt
		PopulateDatabase(database, "./data/crosswords", "./data/common_words.txt")
	case "create":
		// must create output directory before use
		CreatePuzzles(database, 12, "output")
	default:
		fmt.Println("Usage: go run . [command]: Use 'populate' to fill database, or 'create' to create puzzles.")
	}
}
