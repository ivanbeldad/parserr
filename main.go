package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	loadEnvOrFatal()
	check()
}

func check() {
	shows, err := FixFailedShows()
	if err != nil {
		log.Printf("error fixing shows: %s", err.Error())
		return
	}
	err = CleanFixedShows(shows)
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
