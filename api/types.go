package api

import "fmt"

const (
	// EnvSonarrURL ...
	EnvSonarrURL = "SONARR_URL"
	// EnvSonarrAPIKey ...
	EnvSonarrAPIKey = "SONARR_APIKEY"
	// EnvSonarrDownloadFolder ...
	EnvSonarrDownloadFolder = "SONARR_DOWNLOAD_FOLDER"
	// EnvRadarrURL ...
	EnvRadarrURL = "RADARR_URL"
	// EnvRadarrAPIKey ...
	EnvRadarrAPIKey = "RADARR_APIKEY"
	// EnvRadarrDownloadFolder ...
	EnvRadarrDownloadFolder = "RADARR_DOWNLOAD_FOLDER"
	// StatusWarning ...
	StatusWarning = "Warning"
	// CommandStateCompleted ...
	CommandStateCompleted = "completed"
)

// HistoryRec ...
type HistoryRec struct {
	DownloadID            string
	SourceTitle           string
	Status                string
	TrackedDownloadStatus string
	Movie                 Movie
	Series                Series
	Episode               Episode
	Quality               Quality
}

func (h HistoryRec) String() string {
	format := "HistoryRecord\nDownloadID: %s\nSourceTitle: %s\nStatus: %s\nTrackedDownloadStatus: %s\n%s%s%s%s\n"
	return fmt.Sprintf(format, h.DownloadID, h.SourceTitle, h.Status, h.TrackedDownloadStatus, h.Movie, h.Series, h.Episode, h.Quality)
}

// Path Return the path of the movie / show
func (h HistoryRec) Path() string {
	if h.Series.Path != "" {
		return h.Series.Path
	}
	return h.Movie.Path
}

// QueueElem ...
type QueueElem struct {
	ID                    int
	DownloadID            string
	Title                 string
	Status                string
	TrackedDownloadStatus string
	Movie                 Movie
	Series                Series
	Episode               Episode
	Quality               Quality
	StatusMessages        []StatusMessage
}

func (q QueueElem) String() string {
	format := "QueueElement\nID: %d\nDownloadID: %s\nTitle: %s\nStatus: %s\nTrackedDownloadStatus: %s\n%s%s%s%s%s\n"
	return fmt.Sprintf(format, q.ID, q.DownloadID, q.Title, q.Status, q.TrackedDownloadStatus, q.Movie, q.Series, q.Episode, q.Quality, q.StatusMessages)
}

// Path Return the path of the movie / show
func (q QueueElem) Path() string {
	if q.Series.Path != "" {
		return q.Series.Path
	}
	return q.Movie.Path
}

// History ...
type History struct {
	Page     int
	PageSize int
	Records  []HistoryRec
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
	ID    int
	Title string
	Path  string
}

func (s Series) String() string {
	return fmt.Sprintf("Series\nID: %d\nTitle: %s\nPath: %s\n", s.ID, s.Title, s.Path)
}

// Movie ...
type Movie struct {
	ID      int
	Title   string
	Path    string
	HasFile bool
}

func (m Movie) String() string {
	format := "Movie\nID: %d\nTitle: %s\nPath: %s\nHasFile: %v\n"
	return fmt.Sprintf(format, m.ID, m.Title, m.Path, m.HasFile)
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
	ID   int
	Name string
}

func (c Command) String() string {
	return fmt.Sprintf("Command\nID: %d\nName: %s\n", c.ID, c.Name)
}

// CommandStatus ...
type CommandStatus struct {
	Command
	State string `json:"state"`
}

func (c CommandStatus) String() string {
	return fmt.Sprintf("Command\nID: %d\nName: %s\nState: %s\n", c.ID, c.Name, c.State)
}

// CommandBody ...
type CommandBody struct {
	Name      string `json:"name"`
	Path      string `json:"path,omitempty"`
	SeriesIds []int  `json:"seriesIds,omitempty"`
	MovieIds  []int  `json:"movieIds,omitempty"`
}

func (c CommandBody) String() string {
	format := "Command\nName: %s\nSeriesIds: %s\nMovieIds: %s\n"
	return fmt.Sprintf(format, c.Name, c.SeriesIds, c.MovieIds)
}
