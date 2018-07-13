package parser

import (
	"fmt"
	"log"
	"parserr/api"
	"path"
	"path/filepath"
	"strings"
)

// FixFileNames Try to rename downloaded files to the original torrent name
func FixFileNames(failedFiles []*api.Media, m Mover, downloadFolder string) error {
	var errors []string
	for _, file := range failedFiles {
		err := fixFileName(file, m, downloadFolder)
		if err != nil {
			errors = append(errors, err.Error())
		}
	}
	if len(errors) == 0 {
		return nil
	}
	return fmt.Errorf("%s", strings.Join(errors, ", "))
}

func fixFileName(failedFile *api.Media, m Mover, downloadFolder string) (err error) {
	log.Printf("fixing: %s", failedFile.QueueElement.Title)
	originalFileLocation := failedFile.OriginalFileLocation
	if fileIsOnRoot(failedFile, failedFile.OriginalFilename) {
		originalFileLocation, err = moveFileToFolder(failedFile.OriginalFileLocation, m)
		if err != nil {
			log.Printf("cannot move file to a folder")
			return err
		}
	}
	newPath := path.Join(filepath.Dir(originalFileLocation), failedFile.FinalFilename+failedFile.FileExtension)
	log.Printf("moving from %s to %s", originalFileLocation, newPath)
	err = m.Move(originalFileLocation, newPath)
	if err != nil {
		return err
	}
	return nil
}

func fileIsOnRoot(failedFile *api.Media, filename string) bool {
	return failedFile.QueueElement.Title == filename
}

func moveFileToFolder(oldpath string, m Mover) (dest string, err error) {
	log.Printf("moving file to a folder with its own name")
	tmpPath := oldpath + ".tmp"
	err = m.Move(oldpath, tmpPath)
	if err != nil {
		return "", err
	}
	err = m.Mkdir(oldpath)
	if err != nil {
		m.Move(tmpPath, oldpath)
		return "", err
	}
	dest = path.Join(oldpath, filepath.Base(oldpath))
	err = m.Move(tmpPath, dest)
	if err != nil {
		return
	}
	return
}
