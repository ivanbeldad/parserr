package parser

import (
	"fmt"
	"log"
	"sonarr-parser-helper/api"
	"strings"
)

// CleanFixedShows ...
func CleanFixedShows(shows []*Show) error {
	log.Printf("executing rescan series")
	_, err := api.ExecuteCommandAndWait(api.NewRescanSeriesCommand())
	if err != nil {
		return err
	}
	var errors []string
	for _, s := range shows {
		if s.HasBeenRenamed && hasBeenDetected(s) {
			err = api.DeleteQueueItem(s.QueueElement.ID)
			if err != nil {
				log.Print(err)
				errors = append(errors, err.Error())
			} else {
				log.Printf("episode cleared from the queue: %s", s.QueueElement.Title)
			}
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("%s", strings.Join(errors, ", "))
	}
	return nil
}

func hasBeenDetected(s *Show) bool {
	ep, err := api.GetEpisode(s.QueueElement.Episode.ID)
	if err != nil {
		log.Printf("cannot detect if episode %s has been detected", s.QueueElement.Title)
		return false
	}
	return ep.HasFile
}
