package api

const (
	// EnvSonarrURL ...
	EnvSonarrURL = "SONARR_URL"
	// EnvSonarrAPIKey ...
	EnvSonarrAPIKey = "SONARR_APIKEY"
	// EnvSonarrDownloadFolder ...
	EnvSonarrDownloadFolder = "SONARR_DOWNLOAD_FOLDER"
	// StatusWarning ...
	StatusWarning = "Warning"
)

// QueueElement ...
type QueueElement struct {
	Title                 string
	Status                string
	Episode               Episode
	Series                Series
	TrackedDownloadStatus string
}

// Episode ...
type Episode struct {
	SeasonNumber  int
	EpisodeNumber int
}

// Series ...
type Series struct {
	Title string
}

// EpisodeQuality ...
type EpisodeQuality struct {
	Name    string
	Quality Quality
}

// Quality ...
type Quality struct {
	Name string
}

// Command ...
type Command struct {
	Name string
	Path string
}
