package parser

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sonarr-parser-helper/api"
)

// MoveFailedShows ...
func MoveFailedShows(a api.API, m Mover) ([]*api.Media, error) {
	a.ExecuteCommandAndWait(api.NewRescanMovieCommand())
	a.ExecuteCommandAndWait(api.NewRescanSeriesCommand())
	mediaFiles, err := loadFailedMediaFiles(a)
	if err != nil {
		return nil, err
	}
	for _, s := range mediaFiles {
		err = fixNaming(s, m, a.DownloadFolder)
		if err != nil {
			log.Printf("error fixing file %s: %s", s.QueueElement.Title, err.Error())
		}
	}
	return mediaFiles, nil
}

// loadFailedMediaFiles ...
func loadFailedMediaFiles(a api.API) ([]*api.Media, error) {
	mediaFiles := make([]*api.Media, 0)
	queue, err := a.GetQueue()
	if err != nil {
		return nil, err
	}
	history, err := a.GetHistory(1)
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
		for _, hr := range history.Records {
			if itsTheSame(queue[i], hr) {
				found = true
				newMediaFile := api.NewMedia(hr, queue[i])
				mediaFiles = append(mediaFiles, &newMediaFile)
				log.Printf("failed media file detected: %s", queue[i].Title)
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
	return mediaFiles, nil
}

func itsTheSame(qe api.QueueElement, hr api.HistoryRecord) bool {
	sameDownloadID := qe.DownloadID == hr.DownloadID
	sameEpisode := qe.Episode.EpisodeNumber == hr.Episode.EpisodeNumber
	sameSeason := qe.Episode.SeasonNumber == hr.Episode.SeasonNumber
	return sameDownloadID && sameSeason && sameEpisode
}

// fixNaming Try to rename downloaded files to the original
// torrent name.
func fixNaming(mediaFile *api.Media, m Mover, downloadFolder string) error {
	filename, err := mediaFile.GuessFileName()
	if err != nil {
		return err
	}
	oldPath, err := locationOfFile(downloadFolder, filename)
	if err != nil {
		return err
	}
	finalName, err := mediaFile.GuessFinalName(filename)
	if err != nil {
		return err
	}
	newPath := path.Join(mediaFile.QueueElement.Path(), finalName+filepath.Ext(oldPath))
	log.Printf("renaming %s to %s", oldPath, newPath)
	err = m.Move(oldPath, newPath)
	if err != nil {
		return err
	}
	mediaFile.HasBeenRenamed = true
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
