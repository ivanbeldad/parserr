package main

import (
	"log"
	"os"
	"parserr/api"
	"parserr/parser"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	apis := getAPIs()
	for _, a := range apis {
		execute(a)
	}
}

func execute(a api.RRAPI) {
	parser.ExtractAll(a.GetDownloadFolder())
	a.ExecuteCommandAndWait(a.CheckFinishedDownloadsCommand(), api.DefaultRetries)
	move := parser.BasicMover{}
	files, err := parser.FailedMedia(a)
	if err != nil {
		log.Println(err)
		return
	}
	fixStrategy := parser.StrategyFactory(a, move)
	err = parser.FixMedia(files, fixStrategy)
	if err != nil {
		log.Println(err)
		return
	}
}

func getAPIs() (apis []api.RRAPI) {
	if os.Getenv(api.EnvRadarrURL) != "" {
		apis = append(apis, radarr())
	}
	if os.Getenv(api.EnvSonarrURL) != "" {
		apis = append(apis, sonarr())
	}
	return apis
}

func sonarr() api.RRAPI {
	if os.Getenv(api.EnvSonarrAPIKey) == "" {
		log.Fatal("empty sonarr apikey")
	}
	if os.Getenv(api.EnvSonarrDownloadFolder) == "" {
		log.Fatal("empty sonarr download folder")
	}
	if os.Getenv(api.EnvSonarrURL) == "" {
		log.Fatal("empty sonarr url")
	}
	log.Print("adding sonarr api")
	return api.NewSonarr(
		os.Getenv("SONARR_URL"),
		os.Getenv("SONARR_APIKEY"),
		os.Getenv("SONARR_DOWNLOAD_FOLDER"))
}

func radarr() api.RRAPI {
	if os.Getenv(api.EnvRadarrAPIKey) == "" {
		log.Fatal("empty radarr apikey")
	}
	if os.Getenv(api.EnvRadarrDownloadFolder) == "" {
		log.Fatal("empty radarr download folder")
	}
	if os.Getenv(api.EnvRadarrURL) == "" {
		log.Fatal("empty radarr url")
	}
	log.Print("adding radarr api")
	return api.NewRadarr(
		os.Getenv("RADARR_URL"),
		os.Getenv("RADARR_APIKEY"),
		os.Getenv("RADARR_DOWNLOAD_FOLDER"))
}
