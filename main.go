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
	// check()
	testRadarr()
}

func createAPI() api.API {
	return api.NewAPI(
		os.Getenv("RADARR_URL"),
		os.Getenv("RADARR_APIKEY"),
		os.Getenv("RADARR_DOWNLOAD_FOLDER"))
}

func testRadarr() {
	a := createAPI()
	move := parser.DiskMove{}
	_, err := parser.FixFailedShows(a, move)
	if err != nil {
		log.Fatal(err)
	}
}

func check() {
	a := createAPI()
	parser.ExtractAll(os.Getenv(api.EnvSonarrDownloadFolder))
	shows, err := parser.FixFailedShows(a, parser.FakeMove{})
	if err != nil {
		log.Printf("error fixing shows: %s", err.Error())
		return
	}
	if len(shows) == 0 {
		log.Print("no failed episodes")
		return
	}
	err = parser.CleanFixedShows(a, shows)
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
