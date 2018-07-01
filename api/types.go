package api

import "fmt"

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

// HistoryRecord ...
type HistoryRecord struct {
	DownloadID            string
	SourceTitle           string
	Status                string
	TrackedDownloadStatus string
	Series                Series
	Episode               Episode
	Quality               Quality
}

func (h HistoryRecord) String() string {
	format := "HistoryRecord\nDownloadID: %s\nSourceTitle: %s\nStatus: %s\nTrackedDownloadStatus: %s\n%s%s%s\n"
	return fmt.Sprintf(format, h.DownloadID, h.SourceTitle, h.Status, h.TrackedDownloadStatus, h.Series, h.Episode, h.Quality)
}

// QueueElement ...
type QueueElement struct {
	DownloadID            string
	Title                 string
	Status                string
	TrackedDownloadStatus string
	Series                Series
	Episode               Episode
	Quality               Quality
	StatusMessages        []StatusMessage
}

func (q QueueElement) String() string {
	format := "QueueElement\nDownloadID: %s\nTitle: %s\nStatus: %s\nTrackedDownloadStatus: %s\n%s%s%s%s\n"
	return fmt.Sprintf(format, q.DownloadID, q.Title, q.Status, q.TrackedDownloadStatus, q.Series, q.Episode, q.Quality, q.StatusMessages)
}

// History ...
type History struct {
	Page     int
	PageSize int
	Records  []HistoryRecord
}

func (h History) String() string {
	format := "History\nPage: %d\nPageSize: %d\n%s\n"
	return fmt.Sprintf(format, h.Page, h.PageSize, h.Records)
}

// Episode ...
type Episode struct {
	SeasonNumber  int
	EpisodeNumber int
}

func (e Episode) String() string {
	format := "Episode\nSeasonNumber: %d\nEpisodeNumber: %d\n"
	return fmt.Sprintf(format, e.SeasonNumber, e.EpisodeNumber)
}

// Series ...
type Series struct {
	Title string
}

func (s Series) String() string {
	return fmt.Sprintf("Series\nTitle: %s\n", s.Title)
}

// Quality ...
type Quality struct {
	EpisodeQuality EpisodeQuality `json:"quality"`
}

func (q Quality) String() string {
	return fmt.Sprintf("Quality\n%s\n", q.EpisodeQuality)
}

// EpisodeQuality ...
type EpisodeQuality struct {
	Name string
}

func (eq EpisodeQuality) String() string {
	return fmt.Sprintf("EpisodeQuality\nName: %s\n", eq.Name)
}

// StatusMessage ...
type StatusMessage struct {
	Title string
}

func (sm StatusMessage) String() string {
	return fmt.Sprintf("StatusMessage\nTitle: %s\n", sm.Title)
}

// Command ...
type Command struct {
	Name string `json:"name"`
}

// NewRescanSeriesCommand Create a command instance to force to rescan series form disk
func NewRescanSeriesCommand() Command {
	return Command{Name: "RescanSeries"}
}
