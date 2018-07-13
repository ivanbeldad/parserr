package api

import (
	"fmt"
	"log"
	"os"
	"parserr/helpers"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	// TypeMovie ...
	TypeMovie = "movie"
	// TypeShow ...
	TypeShow = "show"
)

// Media ...
type Media struct {
	HistoryRecord        HistoryRecord
	QueueElement         QueueElement
	OriginalFileLocation string
	FinalFileLocation    string
	OriginalFilename     string
	FinalFilename        string
	Type                 string
	FileExtension        string
}

// NewMedia Generate a new Media struct with correct type and names
func NewMedia(a RRAPI, hr HistoryRecord, qe QueueElement) (m Media, err error) {
	if qe.Movie.Title != "" {
		m.Type = TypeMovie
	}
	if qe.Series.Title != "" {
		m.Type = TypeShow
	}
	m.HistoryRecord = hr
	m.QueueElement = qe
	filename, err := m.guessOriginalFilename()
	if err != nil {
		return
	}
	m.OriginalFilename = filename
	finalname, err := m.guessFinalFilename()
	if err != nil {
		return
	}
	m.FinalFilename = finalname
	location, err := helpers.FindFile(a.GetDownloadFolder(), m.OriginalFilename)
	if err != nil {
		return
	}
	m.OriginalFileLocation = location
	m.FileExtension = filepath.Ext(m.OriginalFilename)
	return
}

// IsBroken ...
func (m Media) IsBroken() bool {
	return m.HistoryRecord.TrackedDownloadStatus == TrackedDownloadStatusWarning
}

// HasBeenDetected Return true if the show has been detected,
// false otherwise (including errors)
func (m Media) HasBeenDetected(a RRAPI) bool {
	if m.Type == TypeMovie {
		movie, err := a.GetMovie(m.QueueElement.Movie.ID)
		if err != nil {
			log.Printf("cannot detect if movie %s has been detected", m.QueueElement.Title)
			return false
		}
		return movie.HasFile
	}
	if m.Type == TypeShow {
		ep, err := a.GetEpisode(m.QueueElement.Episode.ID)
		if err != nil {
			log.Printf("cannot detect if episode %s has been detected", m.QueueElement.Title)
			return false
		}
		return ep.HasFile
	}
	return false
}

// DeleteFile Removes the file wherever the show is located
func (m Media) DeleteFile() error {
	if m.FinalFileLocation == "" {
		return fmt.Errorf("cannot delete %s because destiny path is empty", m.QueueElement.Title)
	}
	err := os.Remove(m.FinalFileLocation)
	if err != nil {
		log.Printf("cannot delete %s from %s", m.QueueElement.Title, m.FinalFileLocation)
	}
	return err
}

// GuessFileName ...
func (m Media) guessOriginalFilename() (string, error) {
	if m.Type == TypeMovie {
		return guessMovieFileName(m)
	}
	if m.Type == TypeShow {
		return guessShowFileName(m)
	}
	return "", fmt.Errorf("cannot guess filename of unrecognized media type: %s", m.Type)
}

func guessShowFileName(m Media) (string, error) {
	episode := m.QueueElement.Episode
	regexString := fmt.Sprintf("%d.{0,4}%d", episode.SeasonNumber, episode.EpisodeNumber)
	regex := regexp.MustCompile(regexString)
	for _, message := range m.QueueElement.StatusMessages {
		if regex.MatchString(message.Title) {
			extension := filepath.Ext(message.Title)
			validExtensions := map[string]bool{".mkv": true, ".mp4": true, ".avi": true}
			if validExtensions[extension] {
				return message.Title, nil
			}
			log.Printf("is not a valid file, skipping: %s\n", message.Title)
		}
	}
	return "", fmt.Errorf("impossible to guess file name for %s", m.QueueElement.Title)
}

func guessMovieFileName(m Media) (string, error) {
	for _, message := range m.QueueElement.StatusMessages {
		extension := filepath.Ext(message.Title)
		validExtensions := map[string]bool{".mkv": true, ".mp4": true, ".avi": true}
		if validExtensions[extension] {
			return message.Title, nil
		}
		log.Printf("is not a valid file, skipping: %s\n", message.Title)
	}
	return "", fmt.Errorf("impossible to guess file name for %s", m.QueueElement.Title)
}

// GuessFinalName ...
func (m Media) guessFinalFilename() (string, error) {
	if m.Type == TypeMovie {
		return m.guessMovieFinalName()
	}
	if m.Type == TypeShow {
		return m.guessShowFinalName()
	}
	return "", fmt.Errorf("cannot guess finalname of file with type %q", m.Type)
}

func (m Media) guessMovieFinalName() (string, error) {
	finalTitle := m.HistoryRecord.SourceTitle
	if len(m.QueueElement.StatusMessages) == 1 {
		return finalTitle, nil
	}
	episode := m.QueueElement.Episode
	regexString := fmt.Sprintf("[.\\-_ ]([\\-_0-9sSeExX]{2,10})[.\\-_ ]")
	regex := regexp.MustCompile(regexString)
	if !regex.MatchString(finalTitle) {
		return "", fmt.Errorf("unable to guess final episode name of %s", m.OriginalFilename)
	}
	match := regex.FindString(finalTitle)
	new := fmt.Sprintf(".S%.2dE%.2d.", episode.SeasonNumber, episode.EpisodeNumber)
	finalTitle = strings.Replace(finalTitle, match, new, 1)
	return finalTitle, nil
}

func (m Media) guessShowFinalName() (string, error) {
	finalTitle := m.HistoryRecord.SourceTitle
	if len(m.QueueElement.StatusMessages) == 1 {
		return finalTitle, nil
	}
	episode := m.QueueElement.Episode
	regexString := fmt.Sprintf("[.\\-_ ]([\\-_0-9sSeExX]{2,10})[.\\-_ ]")
	regex := regexp.MustCompile(regexString)
	if !regex.MatchString(finalTitle) {
		return "", fmt.Errorf("unable to guess final episode name of %s", m.OriginalFilename)
	}
	match := regex.FindString(finalTitle)
	new := fmt.Sprintf(".S%.2dE%.2d.", episode.SeasonNumber, episode.EpisodeNumber)
	finalTitle = strings.Replace(finalTitle, match, new, 1)
	return finalTitle, nil
}
