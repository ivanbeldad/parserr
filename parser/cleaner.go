package parser

import (
	"fmt"
	"log"
	"sonarr-parser-helper/api"
	"strings"
	"time"
)

const (
	// MaxTime Max interval to check series and clean them
	MaxTime = time.Minute * 5
	// CheckInterval Time between requests to check if rescan is completed
	CheckInterval = time.Second * 5
)

// CleanFixedShows ...
func CleanFixedShows(shows []*Show) error {
	command, err := api.ExecuteCommand(api.NewRescanSeriesCommand())
	if err != nil {
		return err
	}
	log.Printf("executing rescan series")
	err = waitToFinish(command)
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

func waitToFinish(command api.Command) error {
	totalWait := CheckInterval
	for totalWait <= MaxTime {
		time.Sleep(CheckInterval)
		result, err := api.GetCommand(command.ID)
		if err == nil {
			if result.State == api.CommandStateCompleted {
				log.Print("finished rescan series successfully")
				return nil
			}
		}
		totalWait += CheckInterval
	}
	return fmt.Errorf("timeout checking command rescan series, clean not completed")
}

func hasBeenDetected(s *Show) bool {
	ep, err := api.GetEpisode(s.QueueElement.Episode.ID)
	if err != nil {
		log.Printf("cannot detect if episode %s has been detected", s.QueueElement.Title)
		return false
	}
	return ep.HasFile
}
