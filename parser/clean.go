package parser

import (
	"fmt"
	"log"
	"sonarr-parser-helper/api"
	"strings"
)

// CleanFixedShows ...
func CleanFixedShows(a api.API, mediaFiles []*api.Media) error {
	var err error
	if len(mediaFiles) == 0 {
		return nil
	}
	if mediaFiles[0].Type == api.TypeMovie {
		_, err = a.ExecuteCommandAndWait(api.NewRescanMovieCommand())
	} else {
		_, err = a.ExecuteCommandAndWait(api.NewRescanSeriesCommand())
	}
	if err != nil {
		return err
	}
	var errors []string
	for _, s := range mediaFiles {
		// TODO
		// If there is no way to rename episode
		// or it isn't been detected then
		// add to blacklist and retry download
		if s.HasBeenRenamed && s.HasBeenDetected(a) {
			err = a.DeleteQueueItem(s.QueueElement.ID)
			if err != nil {
				log.Print(err)
				errors = append(errors, err.Error())
			} else {
				log.Printf("file cleared from the queue: %s", s.QueueElement.Title)
			}
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("%s", strings.Join(errors, ", "))
	}
	return nil
}
