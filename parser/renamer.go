package parser

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sonarr-parser-helper/api"
)

// FixFailedShows ...
func FixFailedShows(a api.API, m Move) ([]*api.Media, error) {
	shows, err := loadFailedShows(a)
	if err != nil {
		return nil, err
	}
	for _, s := range shows {
		err = FixNaming(s, m)
		if err != nil {
			log.Printf("error fixing show %s: %s", s.QueueElement.Title, err.Error())
		}
	}
	return shows, nil
}

// loadFailedShows ...
func loadFailedShows(a api.API) ([]*api.Media, error) {
	shows := make([]*api.Media, 0)
	queue, err := a.GetQueue()
	if err != nil {
		return nil, err
	}
	history, err := a.GetHistory(1)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s", queue)
	fmt.Printf("%s", history)
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
				newShow := api.Media{HistoryRecord: he, QueueElement: queue[i]}
				shows = append(shows, &newShow)
				log.Printf("failed show detected: %s", queue[i].Title)
			}
		}
		if !found {
			history, err = addPageToHistory(a, history)
			if err != nil {
				return nil, fmt.Errorf("%s, imposible to guess failed file", err)
			}
			i--
		}
	}
	return shows, nil
}

// FixNaming Try to rename downloaded files to the original
// torrent name.
func FixNaming(s *api.Media, m Move) error {
	filename, err := s.GuessFileName()
	if err != nil {
		return err
	}
	oldPath, err := locationOfFile(os.Getenv(api.EnvSonarrDownloadFolder), filename)
	if err != nil {
		return err
	}
	finalName, err := s.GuessFinalName(filename)
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

func addPageToHistory(a api.API, h api.History) (api.History, error) {
	newPage := h.Page + 1
	newHistory, err := a.GetHistory(newPage)
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
