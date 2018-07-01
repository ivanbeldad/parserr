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
	// CommandStateCompleted ...
	CommandStateCompleted = "completed"
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
	ID                    int
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
	format := "QueueElement\nID: %d\nDownloadID: %s\nTitle: %s\nStatus: %s\nTrackedDownloadStatus: %s\n%s%s%s%s\n"
	return fmt.Sprintf(format, q.ID, q.DownloadID, q.Title, q.Status, q.TrackedDownloadStatus, q.Series, q.Episode, q.Quality, q.StatusMessages)
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
	ID            int
	SeasonNumber  int
	EpisodeNumber int
	HasFile       bool
}

func (e Episode) String() string {
	format := "Episode\nID: %d\nSeasonNumber: %d\nEpisodeNumber: %d\nHasFile: %v\n"
	return fmt.Sprintf(format, e.ID, e.SeasonNumber, e.EpisodeNumber, e.HasFile)
}

// Series ...
type Series struct {
	Title string
	Path  string
}

func (s Series) String() string {
	return fmt.Sprintf("Series\nTitle: %s\nPath: %s\n", s.Title, s.Path)
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
	ID    int
	Name  string
	State string
}

func (c Command) String() string {
	return fmt.Sprintf("Command\nID: %d\nName: %s\nState: %s\n", c.ID, c.Name, c.State)
}

// CommandBody ...
type CommandBody struct {
	Name string `json:"name"`
}

// NewRescanSeriesCommand Create a command instance to force to rescan series form disk
func NewRescanSeriesCommand() CommandBody {
	return CommandBody{Name: "RescanSeries"}
}