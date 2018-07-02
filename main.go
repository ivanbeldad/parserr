package main

import (
	"log"
	"os"
	"sonarr-parser-helper/api"
	"sonarr-parser-helper/parser"

	"github.com/joho/godotenv"
)

func main() {
	loadEnvOrFatal()
	check()
}

func check() {
	parser.ExtractAll(os.Getenv(api.EnvSonarrDownloadFolder))
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
	godotenv.Load()
	if os.Getenv(api.EnvSonarrAPIKey) == "" {
		log.Fatal("empty apikey")
	}
	if os.Getenv(api.EnvSonarrDownloadFolder) == "" {
		log.Fatal("empty download folder")
	}
	if os.Getenv(api.EnvSonarrURL) == "" {
		log.Fatal("empty sonarr url")
	}
}
