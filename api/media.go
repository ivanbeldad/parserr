package api

import (
	"fmt"
	"log"
	"os"
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
	HistoryRecord  HistoryRecord
	QueueElement   QueueElement
	HasBeenRenamed bool
	FileLocation   string
	Type           string
}

// NewMedia Generate a new Media struct with correct type
func NewMedia(hr HistoryRecord, qe QueueElement) (m Media) {
	if qe.Movie.Title != "" {
		m.Type = TypeMovie
	}
	if qe.Series.Title != "" {
		m.Type = TypeShow
	}
	m.HistoryRecord = hr
	m.QueueElement = qe
	return m
}

// IsBroken ...
func (s Media) IsBroken() bool {
	return s.HistoryRecord.TrackedDownloadStatus == TrackedDownloadStatusWarning
}

// HasBeenDetected Return true if the show has been detected,
// false otherwise (including errors)
func (s Media) HasBeenDetected(a RRAPI) bool {
	if s.Type == TypeMovie {
		movie, err := a.GetMovie(s.QueueElement.Movie.ID)
		if err != nil {
			log.Printf("cannot detect if movie %s has been detected", s.QueueElement.Title)
			return false
		}
		return movie.HasFile
	}
	if s.Type == TypeShow {
		ep, err := a.GetEpisode(s.QueueElement.Episode.ID)
		if err != nil {
			log.Printf("cannot detect if episode %s has been detected", s.QueueElement.Title)
			return false
		}
		return ep.HasFile
	}
	return false
}

// DeleteFile Removes the file wherever the show is located
func (s Media) DeleteFile() error {
	if s.FileLocation == "" {
		return fmt.Errorf("cannot delete %s because destiny path is empty", s.QueueElement.Title)
	}
	err := os.Remove(s.FileLocation)
	if err != nil {
		log.Printf("cannot delete %s from %s", s.QueueElement.Title, s.FileLocation)
	}
	return err
}

// GuessFileName ...
func (s Media) GuessFileName() (string, error) {
	if s.Type == TypeMovie {
		return guessMovieFileName(s)
	}
	if s.Type == TypeShow {
		return guessShowFileName(s)
	}
	return "", fmt.Errorf("cannot guess filename of unrecognized media type: %s", s.Type)
}

func guessShowFileName(s Media) (string, error) {
	episode := s.QueueElement.Episode
	regexString := fmt.Sprintf("%d.{0,4}%d", episode.SeasonNumber, episode.EpisodeNumber)
	regex := regexp.MustCompile(regexString)
	for _, message := range s.QueueElement.StatusMessages {
		if regex.MatchString(message.Title) {
			extension := filepath.Ext(message.Title)
			validExtensions := map[string]bool{".mkv": true, ".mp4": true, ".avi": true}
			if validExtensions[extension] {
				return message.Title, nil
			}
			log.Printf("is not a valid file, skipping: %s\n", message.Title)
		}
	}
	return "", fmt.Errorf("imposible to guess file name for %s", s.QueueElement.Title)
}

func guessMovieFileName(s Media) (string, error) {
	for _, message := range s.QueueElement.StatusMessages {
		extension := filepath.Ext(message.Title)
		validExtensions := map[string]bool{".mkv": true, ".mp4": true, ".avi": true}
		if validExtensions[extension] {
			return message.Title, nil
		}
		log.Printf("is not a valid file, skipping: %s\n", message.Title)
	}
	return "", fmt.Errorf("imposible to guess file name for %s", s.QueueElement.Title)
}

// GuessFinalName ...
func (s Media) GuessFinalName(filename string) (string, error) {
	if s.Type == TypeMovie {
		return s.guessMovieFinalName(filename)
	}
	finalTitle := s.HistoryRecord.SourceTitle
	if len(s.QueueElement.StatusMessages) == 1 {
		return finalTitle, nil
	}
	episode := s.QueueElement.Episode
	regexString := fmt.Sprintf("[.\\-_ ]([\\-_0-9sSeExX]{2,10})[.\\-_ ]")
	regex := regexp.MustCompile(regexString)
	if !regex.MatchString(finalTitle) {
		return "", fmt.Errorf("unable to guess final episode name of %s", filename)
	}
	match := regex.FindString(finalTitle)
	new := fmt.Sprintf(".S%.2dE%.2d.", episode.SeasonNumber, episode.EpisodeNumber)
	finalTitle = strings.Replace(finalTitle, match, new, 1)
	return finalTitle, nil
}

func (s Media) guessMovieFinalName(filename string) (string, error) {
	finalTitle := s.HistoryRecord.SourceTitle
	if len(s.QueueElement.StatusMessages) == 1 {
		return finalTitle, nil
	}
	episode := s.QueueElement.Episode
	regexString := fmt.Sprintf("[.\\-_ ]([\\-_0-9sSeExX]{2,10})[.\\-_ ]")
	regex := regexp.MustCompile(regexString)
	if !regex.MatchString(finalTitle) {
		return "", fmt.Errorf("unable to guess final episode name of %s", filename)
	}
	match := regex.FindString(finalTitle)
	new := fmt.Sprintf(".S%.2dE%.2d.", episode.SeasonNumber, episode.EpisodeNumber)
	finalTitle = strings.Replace(finalTitle, match, new, 1)
	return finalTitle, nil
}
