package parser

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sonarr-parser-helper/api"
	"strings"
)

// Show ...
type Show struct {
	HistoryRecord  api.HistoryRecord
	QueueElement   api.QueueElement
	HasBeenRenamed bool
	FileLocation   string
}

// IsBroken ...
func (s Show) IsBroken() bool {
	return s.HistoryRecord.TrackedDownloadStatus == api.TrackedDownloadStatusWarning
}

// FixNaming Try to rename downloaded files to the original
// torrent name.
func (s *Show) FixNaming(m Move) error {
	filename, err := s.guessFileName()
	if err != nil {
		return err
	}
	oldPath, err := locationOfFile(os.Getenv(api.EnvSonarrDownloadFolder), filename)
	if err != nil {
		return err
	}
	finalName, err := s.guessFinalName(filename)
	if err != nil {
		return err
	}
	newPath := path.Join(s.QueueElement.Series.Path, finalName+filepath.Ext(oldPath))
	log.Printf("renaming %s to %s", oldPath, newPath)
	err = m.Move(oldPath, newPath)
	if err != nil {
		return err
	}
	s.HasBeenRenamed = true
	return nil
}

// HasBeenDetected Return true if the show has been detected,
// false otherwise (including errors)
func (s Show) HasBeenDetected() bool {
	ep, err := api.GetEpisode(s.QueueElement.Episode.ID)
	if err != nil {
		log.Printf("cannot detect if episode %s has been detected", s.QueueElement.Title)
		return false
	}
	return ep.HasFile
}

// DeleteFile Removes the file wherever the show is located
func (s Show) DeleteFile() error {
	if s.FileLocation == "" {
		return fmt.Errorf("cannot delete %s because destiny path is empty", s.QueueElement.Title)
	}
	err := os.Remove(s.FileLocation)
	if err != nil {
		log.Printf("cannot delete %s from %s", s.QueueElement.Title, s.FileLocation)
	}
	return err
}

func (s Show) guessFileName() (string, error) {
	episode := s.QueueElement.Episode
	regexString := fmt.Sprintf("%d.{0,4}%d", episode.SeasonNumber, episode.EpisodeNumber)
	regex := regexp.MustCompile(regexString)
	for _, message := range s.QueueElement.StatusMessages {
		if regex.MatchString(message.Title) {
			return message.Title, nil
		}
	}
	return "", fmt.Errorf("imposible to guess file name for %s", s.QueueElement.Title)
}

func (s Show) guessFinalName(filename string) (string, error) {
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
