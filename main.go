package main

import (
	"log"
	"os"
	"sonarr-parser-helper/api"
	"sonarr-parser-helper/parser"
)

func main() {
	apis := getAPIs()
	for _, a := range apis {
		execute(a)
	}
}

func execute(a api.API) {
	parser.ExtractAll(a.DownloadFolder)
	move := parser.DiskMover{}
	files, err := parser.MoveFailedShows(a, move)
	if err != nil {
		log.Println(err)
	}
	err = parser.CleanFixedShows(a, files)
	if err != nil {
		log.Println(err)
	}
	err = parser.Rename(a, files)
	if err != nil {
		log.Println(err)
	}
}

func getAPIs() (apis []api.API) {
	if os.Getenv(api.EnvRadarrURL) != "" {
		apis = append(apis, radarr())
	}
	if os.Getenv(api.EnvSonarrURL) != "" {
		apis = append(apis, sonarr())
	}
	return apis
}

func sonarr() api.API {
	if os.Getenv(api.EnvSonarrAPIKey) == "" {
		log.Fatal("empty sonarr apikey")
	}
	if os.Getenv(api.EnvSonarrDownloadFolder) == "" {
		log.Fatal("empty sonarr download folder")
	}
	if os.Getenv(api.EnvSonarrURL) == "" {
		log.Fatal("empty sonarr url")
	}
	return api.NewAPI(
		os.Getenv("SONARR_URL"),
		os.Getenv("SONARR_APIKEY"),
		os.Getenv("SONARR_DOWNLOAD_FOLDER"),
		api.TypeShow)
}

func radarr() api.API {
	if os.Getenv(api.EnvRadarrAPIKey) == "" {
		log.Fatal("empty radarr apikey")
	}
	if os.Getenv(api.EnvRadarrDownloadFolder) == "" {
		log.Fatal("empty radarr download folder")
	}
	if os.Getenv(api.EnvRadarrURL) == "" {
		log.Fatal("empty radarr url")
	}
	return api.NewAPI(
		os.Getenv("SONARR_URL"),
		os.Getenv("SONARR_APIKEY"),
		os.Getenv("SONARR_DOWNLOAD_FOLDER"),
		api.TypeMovie)
}
