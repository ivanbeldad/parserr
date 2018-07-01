package main

import (
	"fmt"
	"log"
	"sonarr-parser-helper/api"
	"time"
)

const (
	// MaxTime Max interval to check series and clean them
	MaxTime = time.Minute * 5
	// CheckInterval Time between requests to check if rescan is completed
	CheckInterval = time.Second * 5
)

// CleanFixedShows ...
func CleanFixedShows(shows []Show) error {
	command, err := api.ExecuteCommand(api.NewRescanSeriesCommand())
	if err != nil {
		return err
	}
	log.Printf("executing rescan series")
	err = waitToFinish(command)
	if err != nil {
		return err
	}
	var fullError string
	for _, s := range shows {
		if s.HasBeenRenamed {
			err = api.DeleteQueueItem(s.QueueElement.ID)
			if err != nil {
				log.Print(err)
				fullError += ", " + err.Error()
			}
		}
	}
	if fullError != "" {
		return fmt.Errorf("%s", fullError)
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
