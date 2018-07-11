package parser

import (
	"fmt"
	"log"
	"sonarr-parser-helper/api"
	"strings"
)

// Rename ...
func Rename(a api.API, mediaFiles []*api.Media) error {
	fixedMoviesIds := getMoviesIds(a, mediaFiles)
	fixedShowsIds := getShowsIds(a, mediaFiles)
	if len(fixedMoviesIds) > 0 {
		log.Printf("renaming movies with ids: %s", idsStr(fixedMoviesIds))
		cs, err := a.ExecuteCommandAndWait(api.NewRenameMoviesCommand(fixedMoviesIds))
		if err != nil {
			return err
		}
		log.Printf("rename movies status: %s", cs.State)
	}
	if len(fixedShowsIds) > 0 {
		log.Printf("renaming series with ids: %s", idsStr(fixedShowsIds))
		cs, err := a.ExecuteCommandAndWait(api.NewRenameSeriesCommand(fixedShowsIds))
		if err != nil {
			return err
		}
		log.Printf("rename series status: %s", cs.State)
	}
	return nil
}

func getMoviesIds(a api.API, mediaFiles []*api.Media) (fixedMoviesIds []int) {
	for _, file := range mediaFiles {
		if file.HasBeenRenamed && file.HasBeenDetected(a) {
			if file.Type == api.TypeMovie {
				fixedMoviesIds = append(fixedMoviesIds, file.QueueElement.Movie.ID)
			}
		}
	}
	return
}

func getShowsIds(a api.API, mediaFiles []*api.Media) (fixedShowsIds []int) {
	for _, file := range mediaFiles {
		if file.HasBeenRenamed && file.HasBeenDetected(a) {
			if file.Type == api.TypeMovie {
				fixedShowsIds = append(fixedShowsIds, file.QueueElement.Series.ID)
			}
		}
	}
	return
}

func idsStr(ids []int) string {
	arrStr := fmt.Sprint(ids)
	fields := strings.Fields(arrStr)
	return strings.Trim(strings.Join(fields, ", "), "[]")
}
