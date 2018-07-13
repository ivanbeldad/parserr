package parser

import (
	"fmt"
	"parserr/api"
	"strings"
)

// FixMedia Try to rename downloaded files to the original torrent name
func FixMedia(failedMediaFiles []*api.Media, s FixStrategy) error {
	var errors []string
	for _, file := range failedMediaFiles {
		err := s.Fix(file)
		if err != nil {
			errors = append(errors, err.Error())
		}
	}
	if len(errors) == 0 {
		return nil
	}
	return fmt.Errorf("%s", strings.Join(errors, ", "))
}
