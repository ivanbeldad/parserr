package parser

import (
	"fmt"
	"log"
	"sonarr-parser-helper/api"
	"strings"
)

// Rename ...
func Rename(a api.RRAPI, mediaFiles []*api.Media) error {
	fixedIds := getIds(a, mediaFiles)
	if len(fixedIds) > 0 {
		log.Printf("renaming files with ids: %s", idsStr(fixedIds))
		cs, err := a.ExecuteCommandAndWait(a.RenameCommand(fixedIds))
		if err != nil {
			return err
		}
		log.Printf("rename files status: %s", cs.State)
	}
	return nil
}

func getIds(a api.RRAPI, mediaFiles []*api.Media) (fixedIds []int) {
	for _, file := range mediaFiles {
		if file.HasBeenRenamed && file.HasBeenDetected(a) {
			if file.Type == api.TypeMovie {
				fixedIds = append(fixedIds, file.QueueElement.Movie.ID)
			} else {
				fixedIds = append(fixedIds, file.QueueElement.Series.ID)
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
