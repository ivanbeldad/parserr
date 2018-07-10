package parser

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sonarr-parser-helper/api"
)

// FixFailedShows ...
func FixFailedShows(m Move) ([]*Show, error) {
	shows, err := loadFailedShows()
	if err != nil {
		return nil, err
	}
	for _, s := range shows {
		err = s.FixNaming(m)
		if err != nil {
			log.Printf("error fixing show %s: %s", s.QueueElement.Title, err.Error())
		}
	}
	return shows, nil
}

// loadFailedShows ...
func loadFailedShows() ([]*Show, error) {
	shows := make([]*Show, 0)
	queue, err := api.GetQueue()
	if err != nil {
		return nil, err
	}
	history, err := api.GetHistory(1)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(queue); i++ {
		isNotCompleted := queue[i].Status != api.StatusCompleted
		isNotFailed := queue[i].TrackedDownloadStatus != api.TrackedDownloadStatusWarning
		if isNotCompleted || isNotFailed {
			continue
		}
		found := false
		for _, he := range history.Records {
			sameDownloadID := queue[i].DownloadID == he.DownloadID
			sameEpisode := queue[i].Episode.EpisodeNumber == he.Episode.EpisodeNumber
			sameSeason := queue[i].Episode.SeasonNumber == he.Episode.SeasonNumber
			if sameDownloadID && sameSeason && sameEpisode {
				found = true
				newShow := Show{HistoryRecord: he, QueueElement: queue[i]}
				shows = append(shows, &newShow)
				log.Printf("failed show detected: %s", queue[i].Title)
			}
		}
		if !found {
			history, err = addPageToHistory(history)
			if err != nil {
				return nil, fmt.Errorf("%s, imposible to guess failed file", err)
			}
			i--
		}
	}
	return shows, nil
}

func addPageToHistory(h api.History) (api.History, error) {
	newPage := h.Page + 1
	newHistory, err := api.GetHistory(newPage)
	if err != nil {
		return h, err
	}
	h.Records = append(h.Records, newHistory.Records...)
	h.Page = newPage
	return h, nil
}

// locationOfFile Search recursively on root for a file with filename
// and return its path
func locationOfFile(root, filename string) (string, error) {
	var location string
	var err error
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.Name() == filename {
			location = path
			return fmt.Errorf("ok")
		}
		return nil
	})
	if err != nil && err.Error() == "ok" {
		err = nil
	}
	if location == "" {
		err = fmt.Errorf("%s doesn't exists inside %s", filename, root)
	}
	return location, err
}
