package main

import (
	"log"
	"sonarr-parser-helper/parser"

	"github.com/joho/godotenv"
)

func main() {
	loadEnvOrFatal()
	check()
}

func check() {
	shows, err := parser.FixFailedShows()
	if err != nil {
		log.Printf("error fixing shows: %s", err.Error())
		return
	}
	if len(shows) == 0 {
		log.Print("no failed episodes")
		return
	}
	err = parser.CleanFixedShows(shows)
	if err != nil {
		log.Printf("error cleaning shows: %s", err.Error())
		return
	}
}

func loadEnvOrFatal() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}
