package parser

import (
	"fmt"
	"log"
	"parserr/api"
)

// FailedMedia ...
func FailedMedia(a api.RRAPI) ([]*api.Media, error) {
	mediaFiles := make([]*api.Media, 0)
	queue, err := a.GetQueue()
	if err != nil {
		return nil, err
	}
	history := api.History{Page: 0, PageSize: 10}
	for _, qe := range queue {
		if isNotCompletedOrFailed(qe) {
			continue
		}
		found := false
		var err error
		for !found && err == nil {
			found = false
			for _, hr := range history.Records {
				if itsNotTheSame(qe, hr) {
					continue
				}
				found = true
				newMediaFile, fileErr := api.NewMedia(a, hr, qe)
				if fileErr == nil {
					mediaFiles = append(mediaFiles, &newMediaFile)
					log.Printf("add failed media file correctly: %s", qe.Title)
				} else {
					log.Printf("cannot add failed media file: %s", fileErr.Error())
				}
				break
			}
			if !found {
				err = addPageToHistory(a, &history)
			}
		}
	}
	return mediaFiles, nil
}

func addPageToHistory(a api.RRAPI, h *api.History) error {
	h.Page = h.Page + 1
	newHistory, err := a.GetHistory(h.Page)
	if err != nil {
		return err
	}
	if len(newHistory.Records) == 0 {
		return fmt.Errorf("no more pages in history")
	}
	h.Records = append(h.Records, newHistory.Records...)
	return nil
}

func isNotCompletedOrFailed(qe api.QueueElem) bool {
	isNotCompleted := qe.Status != api.StatusCompleted
	isNotFailed := qe.TrackedDownloadStatus != api.TrackedDownloadStatusWarning
	return isNotCompleted || isNotFailed
}

func itsNotTheSame(qe api.QueueElem, hr api.HistoryRec) bool {
	sameDownloadID := qe.DownloadID == hr.DownloadID
	sameEpisode := qe.Episode.EpisodeNumber == hr.Episode.EpisodeNumber
	sameSeason := qe.Episode.SeasonNumber == hr.Episode.SeasonNumber
	itsTheSame := sameDownloadID && sameSeason && sameEpisode
	return !itsTheSame
}
